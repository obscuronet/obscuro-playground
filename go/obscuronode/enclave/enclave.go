package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

// SubmitBlockResponse is the response sent from the enclave back to the node after ingesting a block
type SubmitBlockResponse struct {
	L1Hash      obscurocommon.L1RootHash // The Header Hash of the ingested block
	L1Height    uint                     // The L1 Height of the ingested block
	L1Parent    obscurocommon.L2RootHash // The L1 ParentBlock of the ingested block
	L2Hash      obscurocommon.L2RootHash // The Rollup Hash in the ingested block
	L2Height    uint                     // The Rollup Height in the ingested block
	L2Parent    obscurocommon.L2RootHash // The Rollup Hash ParentBlock inside the ingested block
	Withdrawals []nodecommon.Withdrawal  // The Withdrawals available in Rollup of the ingested block

	ProducedRollup    nodecommon.ExtRollup // The new Rollup when ingesting the block produces a new Rollup
	IngestedBlock     bool                 // Whether the block was ingested or discarded
	IngestedNewRollup bool                 // Whether the block had a new Rollup and the enclave has ingested it
}

// Enclave - The actual implementation of this interface will call an rpc service
type Enclave interface {
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() obscurocommon.AttestationReport

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret

	// FetchSecret - return the shared secret encrypted with the key from the attestation
	FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret

	// Init - initialise an enclave with a seed received by another enclave
	Init(secret obscurocommon.EncryptedSharedEnclaveSecret)

	// IsInitialised - true if the shared secret is avaible
	IsInitialised() bool

	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis() SubmitBlockResponse

	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(blocks []*types.Block)

	// Start - start speculative execution
	Start(block types.Block)

	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(block types.Block) SubmitBlockResponse

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup nodecommon.ExtRollup)

	// SubmitTx - user transactions
	SubmitTx(tx nodecommon.EncryptedTx) error

	// Balance - returns the balance of an address with a block delay
	Balance(address common.Address) uint64

	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool)

	// Stop gracefully stops the enclave
	Stop()

	// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
	GetTransaction(txHash common.Hash) (*L2Tx, bool)
}

type enclaveImpl struct {
	node           common.Address
	mining         bool
	storage        Storage
	blockResolver  BlockResolver
	statsCollector StatsCollector

	txCh                 chan L2Tx
	roundWinnerCh        chan *Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork
}

func (e *enclaveImpl) Start(block types.Block) {
	headerHash := block.Hash()
	s, f := e.storage.FetchBlockState(headerHash)
	if !f {
		panic("State should be calculated")
	}

	currentHead := s.head
	currentState := newProcessedState(e.storage.FetchRollupState(currentHead.Hash()))
	var currentProcessedTxs []L2Tx
	currentProcessedTxsMap := make(map[common.Hash]L2Tx)

	// start the speculative rollup execution loop
	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:

			currentHead = winnerRollup
			currentState = newProcessedState(e.storage.FetchRollupState(winnerRollup.Hash()))

			// determine the transactions that were not yet included
			currentProcessedTxs = currentTxs(winnerRollup, e.storage.FetchMempoolTxs(), e.storage)
			currentProcessedTxsMap = makeMap(currentProcessedTxs)

			// calculate the State after executing them
			currentState = executeTransactions(currentProcessedTxs, currentState)

		case tx := <-e.txCh:
			_, found := currentProcessedTxsMap[tx.Hash()]
			if !found {
				currentProcessedTxsMap[tx.Hash()] = tx
				currentProcessedTxs = append(currentProcessedTxs, tx)
				executeTx(&currentState, tx)
			}

		case <-e.speculativeWorkInCh:
			b := make([]L2Tx, 0, len(currentProcessedTxs))
			b = append(b, currentProcessedTxs...)
			state := copyProcessedState(currentState)
			e.speculativeWorkOutCh <- speculativeWork{
				r:   currentHead,
				s:   &state,
				txs: b,
			}

		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) ProduceGenesis() SubmitBlockResponse {
	return SubmitBlockResponse{
		L2Hash:         GenesisRollup.Header.Hash(),
		L1Hash:         obscurocommon.GenesisHash,
		ProducedRollup: GenesisRollup.ToExtRollup(),
		IngestedBlock:  true,
	}
}

func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) {
	for _, block := range blocks {
		e.storage.StoreBlock(block)
		updateState(block, e.storage, e.blockResolver)
	}
}

func (e *enclaveImpl) SubmitBlock(block types.Block) SubmitBlockResponse {
	// Todo - investigate further why this is needed.
	// So far this seems to recover correctly
	defer func() {
		if r := recover(); r != nil {
			log.Log(fmt.Sprintf("Agg%d Panic %s", obscurocommon.ShortAddress(e.node), r))
		}
	}()

	_, foundBlock := e.storage.FetchBlock(block.Hash())
	if foundBlock {
		return SubmitBlockResponse{IngestedBlock: false}
	}

	e.storage.StoreBlock(&block)
	// this is where much more will actually happen.
	// the "blockchain" logic from geth has to be executed here,
	// to determine the total proof of work, to verify some key aspects, etc

	_, f := e.storage.FetchBlock(block.Header().ParentHash)
	if !f && e.storage.HeightBlock(&block) > obscurocommon.L1GenesisHeight {
		return SubmitBlockResponse{IngestedBlock: false}
	}
	blockState := updateState(&block, e.storage, e.blockResolver)

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	e.storage.RemoveMempoolTxs(historicTxs(blockState.head, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	return SubmitBlockResponse{
		L1Hash:      block.Hash(),
		L1Height:    uint(e.blockResolver.HeightBlock(&block)),
		L1Parent:    blockState.block.Header().ParentHash,
		L2Hash:      blockState.head.Hash(),
		L2Height:    uint(blockState.head.Height.Load().(int)),
		L2Parent:    blockState.head.Header.ParentHash,
		Withdrawals: blockState.head.Header.Withdrawals,

		ProducedRollup:    r.ToExtRollup(),
		IngestedBlock:     true,
		IngestedNewRollup: blockState.foundNewRollup,
	}
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := Rollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Txs),
	}

	// only store if the parent exists
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(&r)
	} else {
		log.Log(fmt.Sprintf("Agg%d:> Received rollup with no parent: r_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(r.Hash())))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := DecryptTx(tx)
	err := verifySignature(&decryptedTx)
	if err != nil {
		return err
	}
	e.storage.AddMempoolTx(decryptedTx)
	e.txCh <- decryptedTx
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", parent))
	}

	rollupsReceivedFromPeers := e.storage.FetchRollups(e.storage.HeightRollup(head) + 1)
	// filter out rollups with a different ParentBlock
	var usefulRollups []*Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.FetchRollupState(head.Hash())
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState, e.storage, e.blockResolver)

	e.storage.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := winnerRollup.Proof(e.blockResolver)
		w := e.storage.ParentRollup(winnerRollup)
		log.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. FetchRollupTxs: %v. State=%v.",
			obscurocommon.ShortAddress(e.node),
			obscurocommon.ShortHash(winnerRollup.Hash()), e.storage.HeightRollup(winnerRollup),
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State),
		)
		return winnerRollup.ToExtRollup(), true
	}
	return nodecommon.ExtRollup{}, false
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *Rollup) {
	//if atomic.LoadInt32(e.interrupt) == 1 {
	//	return
	//}
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	return e.storage.FetchHeadState().state[address]
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs blockState) *Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := *speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	// if true {
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.head.Hash()) {
		if speculativeRollup.r != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				obscurocommon.ShortAddress(e.node),
				obscurocommon.ShortHash(speculativeRollup.r.Hash()),
				e.storage.HeightRollup(speculativeRollup.r),
				obscurocommon.ShortHash(bs.head.Hash()),
				e.storage.HeightRollup(bs.head)),
			)
			e.statsCollector.L2Recalc(e.node)
		}

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.head, e.storage.FetchMempoolTxs(), e.storage)
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.state))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.head.Proof(e.blockResolver)
	newRollupState = processDeposits(proof, b, copyProcessedState(newRollupState), e.blockResolver)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := NewRollup(b, bs.head, e.node, newRollupTxs, newRollupState.w, obscurocommon.GenerateNonce(), serialize(newRollupState.s))
	// h := r.Height(e.storage)
	// fmt.Printf("h:=%d\n", h)
	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) (*L2Tx, bool) {
	// todo add some sort of cache
	rollup := e.storage.FetchHeadState().head

	var found bool
	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx, true
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		rollup, found = e.storage.FetchRollup(rollup.Hash())
		if !found {
			panic(fmt.Sprintf("Could not find rollup: r_%s", rollup.Hash()))
		}
		if rollup.Height.Load() == obscurocommon.L2GenesisHeight {
			return nil, false
		}
	}
}

func (e *enclaveImpl) Stop() {
	e.exitCh <- true
}

func (e *enclaveImpl) Attestation() obscurocommon.AttestationReport {
	// Todo
	return obscurocommon.AttestationReport{Owner: e.node}
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		panic(fmt.Sprintf("Could not generate secret: %s", err))
	}
	e.storage.StoreSecret(secret)
	return encryptSecret(secret)
}

// Init - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) Init(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.storage.StoreSecret(decryptSecret(secret))
}

func (e *enclaveImpl) FetchSecret(obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return encryptSecret(e.storage.FetchSecret())
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.storage.FetchSecret() != nil
}

// Todo - implement with crypto
func decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) SharedEnclaveSecret {
	return SharedEnclaveSecret(secret)
}

// Todo - implement with crypto
func encryptSecret(secret SharedEnclaveSecret) obscurocommon.EncryptedSharedEnclaveSecret {
	return obscurocommon.EncryptedSharedEnclaveSecret(secret)
}

// internal structure to pass information.
type speculativeWork struct {
	r   *Rollup
	s   *RollupState
	txs []L2Tx
}

func NewEnclave(id common.Address, mining bool, collector StatsCollector) Enclave {
	storage := NewStorage()
	return &enclaveImpl{
		node:                 id,
		storage:              storage,
		blockResolver:        storage,
		mining:               mining,
		txCh:                 make(chan L2Tx),
		roundWinnerCh:        make(chan *Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		statsCollector:       collector,
	}
}
