package enclave

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"

	"github.com/obscuronet/go-obscuro/go/common/profiler"

	"github.com/obscuronet/go-obscuro/go/common/log"

	obscurocrypto "github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"

	"github.com/obscuronet/go-obscuro/go/enclave/rpcencryptionmanager"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/config"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

// StatsCollector Todo - replace with a proper framework
type StatsCollector interface {
	// L2Recalc registers when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id gethcommon.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	config               config.EnclaveConfig
	nodeShortID          uint64
	storage              db.Storage
	blockResolver        db.BlockResolver
	mempool              mempool.Manager
	statsCollector       StatsCollector
	l1Blockchain         *core.BlockChain
	rpcEncryptionManager rpcencryptionmanager.RPCEncryptionManager
	bridge               *bridge.Bridge

	chain *rollupchain.RollupChain

	txCh          chan *common.L2Tx
	roundWinnerCh chan *obscurocore.Rollup
	exitCh        chan bool

	// Todo - disabled temporarily until TN1 is released
	// speculativeWorkInCh  chan bool
	// speculativeWorkOutCh chan speculativeWork

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	erc20ContractLib    erc20contractlib.ERC20ContractLib
	attestationProvider AttestationProvider // interface for producing attestation reports and verifying them

	enclaveKey    *ecdsa.PrivateKey // this is a key specific to this enclave, which is included in the Attestation. Used for signing rollups and for encryption of the shared secret.
	enclavePubKey []byte            // the public key of the above

	transactionBlobCrypto obscurocrypto.TransactionBlobCrypto
	profiler              *profiler.Profiler
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(
	config config.EnclaveConfig,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	collector StatsCollector,
) common.Enclave {
	if len(config.ERC20ContractAddresses) < 2 {
		log.Panic("failed to initialise enclave. At least two ERC20 contract addresses are required - the OBX " +
			"ERC20 address and the ETH ERC20 address")
	}

	// todo - add the delay: N hashes

	nodeShortID := common.ShortAddress(config.HostID)

	var prof *profiler.Profiler
	// don't run a profiler on an attested enclave
	if !config.WillAttest && config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultEnclavePort)
		err := prof.Start()
		if err != nil {
			log.Panic("unable to start the profiler: %s", err)
		}
	}

	// Initialise the database
	backingDB, err := db.CreateDBFromConfig(nodeShortID, config)
	if err != nil {
		log.Panic("Failed to connect to backing database - %s", err)
	}
	chainConfig := params.ChainConfig{
		ChainID:             big.NewInt(config.ObscuroChainID),
		HomesteadBlock:      gethcommon.Big0,
		DAOForkBlock:        gethcommon.Big0,
		EIP150Block:         gethcommon.Big0,
		EIP155Block:         gethcommon.Big0,
		EIP158Block:         gethcommon.Big0,
		ByzantiumBlock:      gethcommon.Big0,
		ConstantinopleBlock: gethcommon.Big0,
		PetersburgBlock:     gethcommon.Big0,
		IstanbulBlock:       gethcommon.Big0,
		MuirGlacierBlock:    gethcommon.Big0,
		BerlinBlock:         gethcommon.Big0,
		LondonBlock:         gethcommon.Big0,
	}
	storage := db.NewStorage(backingDB, nodeShortID, &chainConfig)

	// Initialise the Ethereum "Blockchain" structure that will allow us to validate incoming blocks
	// Todo - check the minimum difficulty parameter
	var l1Blockchain *core.BlockChain
	if config.ValidateL1Blocks {
		if config.GenesisJSON == nil {
			log.Panic("enclave is configured to validate blocks, but genesis JSON is nil")
		}
		l1Blockchain = rollupchain.NewL1Blockchain(config.GenesisJSON)
	} else {
		common.LogWithID(common.ShortAddress(config.HostID), "validateBlocks is set to false. L1 blocks will not be validated.")
	}

	// Todo- make sure the enclave cannot be started in production with WillAttest=false
	var attestationProvider AttestationProvider
	if config.WillAttest {
		attestationProvider = &EgoAttestationProvider{}
	} else {
		common.LogWithID(nodeShortID, "WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &DummyAttestationProvider{}
	}

	// todo - this has to be read from the database when the node restarts.
	// first time the node starts we derive the obscuro key from the master seed received after the shared secret exchange
	common.LogWithID(nodeShortID, "Generating the Obscuro key")

	// todo - save this to the db
	enclaveKey, err := crypto.GenerateKey()
	if err != nil {
		log.Panic("Failed to generate enclave key. Cause: %s", err)
	}
	serializedEnclavePubKey := crypto.CompressPubkey(&enclaveKey.PublicKey)
	common.LogWithID(nodeShortID, "Generated public key %s", gethcommon.Bytes2Hex(serializedEnclavePubKey))

	obscuroKey := obscurocrypto.GetObscuroKey()
	rpcem := rpcencryptionmanager.NewRPCEncryptionManager(config.ViewingKeysEnabled, ecies.ImportECDSA(obscuroKey))

	transactionBlobCrypto := obscurocrypto.NewTransactionBlobCryptoImpl()

	obscuroBridge := bridge.New(
		config.ERC20ContractAddresses[0],
		config.ERC20ContractAddresses[1],
		mgmtContractLib,
		erc20ContractLib,
		nodeShortID,
		transactionBlobCrypto,
		config.ObscuroChainID,
		config.L1ChainID,
	)
	memp := mempool.New(config.ObscuroChainID)

	chain := rollupchain.New(nodeShortID, config.HostID, storage, l1Blockchain, obscuroBridge, transactionBlobCrypto, memp, rpcem, enclaveKey, config.L1ChainID, &chainConfig)

	return &enclaveImpl{
		config:                config,
		nodeShortID:           nodeShortID,
		storage:               storage,
		blockResolver:         storage,
		mempool:               memp,
		statsCollector:        collector,
		l1Blockchain:          l1Blockchain,
		rpcEncryptionManager:  rpcem,
		bridge:                obscuroBridge,
		chain:                 chain,
		txCh:                  make(chan *common.L2Tx),
		roundWinnerCh:         make(chan *obscurocore.Rollup),
		exitCh:                make(chan bool),
		mgmtContractLib:       mgmtContractLib,
		erc20ContractLib:      erc20ContractLib,
		attestationProvider:   attestationProvider,
		enclaveKey:            enclaveKey,
		enclavePubKey:         serializedEnclavePubKey,
		transactionBlobCrypto: transactionBlobCrypto,
		profiler:              prof,
	}
}

// IsReady is only implemented by the RPC wrapper
func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() error {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	// todo - reinstate after TN1
	/*	if e.config.SpeculativeExecution {
			//start the speculative rollup execution loop on its own go routine
			go e.start(block)
		}
	*/
}

func (e *enclaveImpl) ProduceGenesis(blkHash gethcommon.Hash) common.BlockSubmissionResponse {
	rolGenesis, b := e.chain.ProduceGenesis(blkHash)
	return common.BlockSubmissionResponse{
		ProducedRollup: e.transactionBlobCrypto.ToExtRollup(rolGenesis),
		BlockHeader:    b.Header(),
		IngestedBlock:  true,
	}
}

// IngestBlocks is used to update the enclave with the full history of the L1 chain to date.
func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []common.BlockSubmissionResponse {
	result := make([]common.BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		response := e.chain.IngestBlock(block)
		result[i] = response
		if !response.IngestedBlock {
			return result // We return early, as all descendant blocks will also fail verification.
		}
	}
	return result
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitBlock(block types.Block) common.BlockSubmissionResponse {
	bsr := e.chain.SubmitBlock(block)

	if bsr.RollupHead != nil {
		hr, f := e.storage.FetchRollup(bsr.RollupHead.Hash())
		if !f {
			log.Panic("This should not happen because this rollup was just processed.")
		}
		e.mempool.RemoveMempoolTxs(hr, e.storage)
	}

	return bsr
}

func (e *enclaveImpl) SubmitRollup(rollup common.ExtRollup) {
	r := e.transactionBlobCrypto.ToEnclaveRollup(rollup.ToRollup())

	// only store if the parent exists
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(r)
	} else {
		common.LogWithID(e.nodeShortID, "Received rollup with no parent: r_%d", common.ShortHash(r.Hash()))
	}
}

func (e *enclaveImpl) SubmitTx(tx common.EncryptedTx) (common.EncryptedResponseSendRawTx, error) {
	decryptedTx, err := e.rpcEncryptionManager.ExtractTxFromBinary(tx)
	if err != nil {
		log.Info(fmt.Sprintf("could not decrypt transaction. Cause: %s", err))
		return nil, fmt.Errorf("could not decrypt transaction. Cause: %w", err)
	}
	err = e.mempool.AddMempoolTx(decryptedTx)
	if err != nil {
		return nil, err
	}

	if e.config.SpeculativeExecution {
		e.txCh <- decryptedTx
	}

	// TODO - Once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`.
	signer := types.NewLondonSigner(big.NewInt(e.config.ObscuroChainID))
	sender, err := signer.Sender(decryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not recover sender to encrypt eth_sendRawTransaction response. Cause: %w", err)
	}

	txHashBytes := []byte(decryptedTx.Hash().Hex())
	encryptedResult, err := e.rpcEncryptionManager.EncryptWithViewingKey(sender, txHashBytes)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_sendRawTransaction request. Cause: %w", err)
	}

	return encryptedResult, nil
}

func (e *enclaveImpl) RoundWinner(parent common.L2RootHash) (common.ExtRollup, bool, error) {
	return e.chain.RoundWinner(parent)
}

func (e *enclaveImpl) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	resp, err := e.chain.ExecuteOffChainTransaction(encryptedParams)
	if err != nil {
		common.LogWithID(e.nodeShortID, "Could not execute off chain call. Cause: %s", err)
	}
	return resp, err
}

func (e *enclaveImpl) Nonce(address gethcommon.Address) uint64 {
	// todo user encryption
	hs := e.storage.FetchHeadState()
	if hs == nil {
		return 0
	}
	s := e.storage.CreateStateDB(hs.HeadRollup)
	return s.GetNonce(address)
}

func (e *enclaveImpl) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (common.EncryptedResponseGetTxByHash, error) {
	hashBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt encrypted RPC request params. Cause: %w", err)
	}
	var paramList []string
	err = json.Unmarshal(hashBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
	}
	txHash := gethcommon.HexToHash(paramList[0])

	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index := e.storage.GetTransaction(txHash)

	if tx == nil {
		// If there's no transaction, there's no `from` field we can use to determine which key to use to encrypt the response.
		return nil, fmt.Errorf("transaction does not exist")
	}

	// Unlike in the Geth impl, we hardcode the use of a London signer.
	signer := types.NewLondonSigner(tx.ChainId())
	rpcTx := newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer)

	txBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, fmt.Errorf("could not marshal transaction to JSON. Cause: %w", err)
	}
	return e.rpcEncryptionManager.EncryptWithViewingKey(rpcTx.From, txBytes)
}

func (e *enclaveImpl) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (common.EncryptedResponseGetTxReceipt, error) {
	txHash, err := e.rpcEncryptionManager.ExtractTxHash(encryptedParams)
	if err != nil {
		return nil, err
	}

	viewingKeyAddress, err := e.storage.GetSender(txHash)
	if err != nil {
		if errors.Is(err, db.ErrTxNotFound) {
			return nil, nil
		}
		return nil, err
	}

	txReceipt, err := e.storage.GetTransactionReceipt(txHash)
	if err != nil {
		if errors.Is(err, db.ErrTxNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	encryptedTxReceipt, err := e.rpcEncryptionManager.EncryptTxReceiptWithViewingKey(viewingKeyAddress, txReceipt)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getTransactionReceipt request. Cause: %w", err)
	}

	return encryptedTxReceipt, nil
}

func (e *enclaveImpl) GetRollup(rollupHash common.L2RootHash) (*common.ExtRollup, error) {
	rollup, found := e.storage.FetchRollup(rollupHash)
	if found {
		extRollup := e.transactionBlobCrypto.ToExtRollup(rollup)
		return &extRollup, nil
	}
	return nil, fmt.Errorf("rollup with hash %s could not be found", rollupHash.Hex())
}

func (e *enclaveImpl) GetRollupByHeight(rollupHeight int64) (*common.ExtRollup, error) {
	// TODO - Consider improving efficiency by directly fetching rollup by number.
	rollup := e.storage.FetchHeadRollup()
	maxRollupHeight := rollup.NumberU64()

	// -1 is used by Ethereum to indicate that we should fetch the head.
	if rollupHeight != -1 {
		for {
			if rollup == nil {
				// We've reached the head of the chain without finding the block.
				return nil, fmt.Errorf("rollup with height %d could not be found. Max rollup height was %d", rollupHeight, maxRollupHeight)
			}
			if rollup.Number().Int64() == rollupHeight {
				// We have found the block.
				break
			}
			if rollup.Number().Int64() < rollupHeight {
				// The current block number is below the sought number. Continuing to walk up the chain is pointless.
				return nil, fmt.Errorf("rollup with height %d could not be found. Max rollup height was %d", rollupHeight, maxRollupHeight)
			}

			// We grab the next rollup and loop.
			rollup = e.storage.ParentRollup(rollup)
		}
	}

	extRollup := e.transactionBlobCrypto.ToExtRollup(rollup)
	return &extRollup, nil
}

func (e *enclaveImpl) Attestation() *common.AttestationReport {
	if e.enclavePubKey == nil {
		log.Panic("public key not initialized, we can't produce the attestation report")
	}
	report, err := e.attestationProvider.GetReport(e.enclavePubKey, e.config.HostID, e.config.HostAddress)
	if err != nil {
		log.Panic("Failed to produce remote report.")
	}
	return report
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() common.EncryptedSharedEnclaveSecret {
	secret := obscurocrypto.GenerateEntropy()
	e.storage.StoreSecret(secret)
	encSec, err := obscurocrypto.EncryptSecret(e.enclavePubKey, secret, e.nodeShortID)
	if err != nil {
		log.Panic("failed to encrypt secret. Cause: %s", err)
	}
	return encSec
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s common.EncryptedSharedEnclaveSecret) error {
	secret, err := obscurocrypto.DecryptSecret(s, e.enclaveKey)
	if err != nil {
		return err
	}
	e.storage.StoreSecret(*secret)
	common.TraceWithID(e.nodeShortID, "Secret decrypted and stored. Secret: %v", secret)
	return nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (e *enclaveImpl) ShareSecret(att *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := e.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, err
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = VerifyIdentity(data, att); err != nil {
		return nil, err
	}
	common.LogWithID(e.nodeShortID, "Successfully verified attestation and identity. Owner: %s", att.Owner)

	secret := e.storage.FetchSecret()
	if secret == nil {
		return nil, errors.New("secret was nil, no secret to share - this shouldn't happen")
	}
	return obscurocrypto.EncryptSecret(att.PubKey, *secret, e.nodeShortID)
}

func (e *enclaveImpl) AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error {
	return e.rpcEncryptionManager.AddViewingKey(encryptedViewingKeyBytes, signature)
}

func (e *enclaveImpl) StoreAttestation(att *common.AttestationReport) error {
	common.LogWithID(e.nodeShortID, "Store attestation. Owner: %s", att.Owner)
	// Store the attestation
	key, err := crypto.DecompressPubkey(att.PubKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key %w", err)
	}
	e.storage.StoreAttestedKey(att.Owner, key)
	return nil
}

func (e *enclaveImpl) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	return e.chain.GetBalance(encryptedParams)
}

func (e *enclaveImpl) GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, error) {
	return e.storage.CreateStateDB(*rollupHash).GetCode(address), nil
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.storage.FetchSecret() != nil
}

func (e *enclaveImpl) Stop() error {
	if e.config.SpeculativeExecution {
		e.exitCh <- true
	}

	if e.profiler != nil {
		return e.profiler.Stop()
	}

	return nil
}

// Todo - reinstate speculative execution after TN1
/*
// internal structure to pass information.
type speculativeWork struct {
	found bool
	r     *obscurocore.Rollup
	s     *state.StateDB
	h     *nodecommon.Header
	txs   []*nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup              // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header               // the header of the new rollup
	processedTxs    []*nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]*nodecommon.L2Tx // structure used to prevent duplicates
	state           *state.StateDB                   // the state as calculated from the previous rollup and the processed transactions
}
*/
/*
func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]*nodecommon.L2Tx)}
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		env.headRollup, _ = e.storage.FetchRollup(blockState.HeadRollup)
		if env.headRollup != nil {
			env.state = e.storage.CreateStateDB(env.headRollup.Hash())
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:
			hash := winnerRollup.Hash()
			env.header = obscurocore.NewHeader(&hash, winnerRollup.Header.Number.Uint64()+1, e.config.HostID)
			env.headRollup = winnerRollup
			env.state = e.storage.CreateStateDB(winnerRollup.Hash())
			common.TraceWithID(e.nodeShortID, "Create new speculative env  r_%d(%d).",
				obscurocommon.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			))

			// determine the transactions that were not yet included
			env.processedTxs = currentTxs(winnerRollup, e.mempool.FetchMempoolTxs(), e.storage)
			env.processedTxsMap = obscurocore.MakeMap(env.processedTxs)

			// calculate the State after executing them
			evm.ExecuteTransactions(env.processedTxs, env.state, env.headRollup.Header, e.storage, e.config.ObscuroChainID, 0)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					evm.ExecuteTransactions([]*nodecommon.L2Tx{tx}, env.state, env.header, e.storage, e.config.ObscuroChainID, 0)
				}
			}

		case <-e.speculativeWorkInCh:
			if env.header == nil {
				e.speculativeWorkOutCh <- speculativeWork{found: false}
			} else {
				b := make([]*nodecommon.L2Tx, 0, len(env.processedTxs))
				b = append(b, env.processedTxs...)
				e.speculativeWorkOutCh <- speculativeWork{
					found: true,
					r:     env.headRollup,
					s:     env.state,
					h:     env.header,
					txs:   b,
				}
			}

		case <-e.exitCh:
			return
		}
	}
}
*/
