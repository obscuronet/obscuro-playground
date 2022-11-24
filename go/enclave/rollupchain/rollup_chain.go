package rollupchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/events"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
)

const (
	msgNoRollup = "could not fetch rollup"
)

var (
	errBlockAlreadyProcessed = errors.New("block already processed")
	errBlockAncestorNotFound = errors.New("block ancestor not found")
)

// RollupChain represents the canonical chain, and manages the state.
type RollupChain struct {
	hostID          gethcommon.Address
	nodeType        common.NodeType
	ethereumChainID int64
	chainConfig     *params.ChainConfig

	storage               db.Storage
	l1Blockchain          *core.BlockChain
	bridge                *bridge.Bridge
	transactionBlobCrypto crypto.TransactionBlobCrypto // todo - remove
	rpcEncryptionManager  rpc.EncryptionManager
	mempool               mempool.Manager
	faucet                Faucet
	subscriptionManager   *events.SubscriptionManager
	crossChainProcessors  *crosschain.Processors

	enclavePrivateKey    *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	blockProcessingMutex sync.Mutex
	logger               gethlog.Logger

	// Gas usage values
	// TODO use the ethconfig.Config instead
	GlobalGasCap uint64
	BaseFee      *big.Int
}

func New(hostID gethcommon.Address, nodeType common.NodeType, storage db.Storage, l1Blockchain *core.BlockChain, bridge *bridge.Bridge, subscriptionManager *events.SubscriptionManager, crossChainProcessors *crosschain.Processors, txCrypto crypto.TransactionBlobCrypto, mempool mempool.Manager, rpcem rpc.EncryptionManager, privateKey *ecdsa.PrivateKey, ethereumChainID int64, chainConfig *params.ChainConfig, logger gethlog.Logger) *RollupChain {
	return &RollupChain{
		hostID:                hostID,
		nodeType:              nodeType,
		storage:               storage,
		l1Blockchain:          l1Blockchain,
		bridge:                bridge,
		transactionBlobCrypto: txCrypto,
		mempool:               mempool,
		faucet:                NewFaucet(storage),
		subscriptionManager:   subscriptionManager,
		crossChainProcessors:  crossChainProcessors,
		enclavePrivateKey:     privateKey,
		rpcEncryptionManager:  rpcem,
		ethereumChainID:       ethereumChainID,
		chainConfig:           chainConfig,
		blockProcessingMutex:  sync.Mutex{},
		logger:                logger,
		GlobalGasCap:          5_000_000_000,
		BaseFee:               gethcommon.Big0,
	}
}

func (rc *RollupChain) ProduceGenesis(blkHash gethcommon.Hash) (*obscurocore.Rollup, *types.Block) {
	b, err := rc.storage.FetchBlock(blkHash)
	if err != nil {
		rc.logger.Crit("Could not retrieve the block used as proof for the genesis rollup.", log.ErrKey, err)
	}

	rolGenesis := obscurocore.NewRollup(
		blkHash,
		nil,
		common.L2GenesisHeight,
		gethcommon.HexToAddress("0x0"),
		[]*common.L2Tx{},
		[]common.Withdrawal{},
		common.GenerateNonce(),
		rc.faucet.GetGenesisRoot(rc.storage),
	)

	// todo::
	// Probably not the best place to put this, but ...
	deployTx, err := rc.crossChainProcessors.LocalManager.GenerateMessageBusDeployTx()
	if err != nil {
		rc.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	if err := rc.mempool.AddMempoolTx(deployTx); err != nil {
		rc.logger.Crit("Cannot create synthetic transaction for deploying the message bus contract on :|")
	}

	rc.signRollup(rolGenesis)

	return rolGenesis, b
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block
// note: this method shouldn't be called for blocks we've seen before
func (rc *RollupChain) insertBlockIntoL1Chain(block *types.Block, isLatest bool) (*blockIngestionType, error) {
	if rc.l1Blockchain != nil {
		_, err := rc.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			return nil, fmt.Errorf("block was invalid: %w", err)
		}
	}
	// todo: this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := rc.storage.FetchHeadBlock()

	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo: we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &blockIngestionType{latest: isLatest, fork: false, preGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)

		// we do a basic sanity check, comparing the received block to the head block on the chain
	} else if block.ParentHash() != prevL1Head.Hash() {
		lcaBlock, err := gethutil.LCA(block, prevL1Head, rc.storage)
		if err != nil {
			return nil, errBlockAncestorNotFound
		}
		rc.logger.Trace("parent not found",
			"blkHeight", block.NumberU64(), "blkHash", block.Hash(),
			"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
			"lcaHeight", lcaBlock.NumberU64(), "lcaHash", lcaBlock.Hash(),
		)
		if lcaBlock.NumberU64() >= prevL1Head.NumberU64() {
			// This is an unexpected error scenario (a bug) because if:
			// lca == prevL1Head:
			//   if prev L1 head is (e.g) a grandfather of ingested block, and block's parent has been seen (else LCA would error),
			//   then why is ingested block's parent not the prev l1 head
			// lca > prevL1Head:
			//   this would imply ingested block is earlier on the same branch as l1 head, but ingested block should not have been seen before
			rc.logger.Error("unexpected blockchain state, incoming block is not child of L1 head and not an earlier fork of L1 head",
				"blkHeight", block.NumberU64(), "blkHash", block.Hash(),
				"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
				"lcaHeight", lcaBlock.NumberU64(), "lcaHash", lcaBlock.Hash(),
			)
			return nil, errors.New("unexpected blockchain state")
		}

		// ingested block is on a different branch to the previously ingested block - we may have to rewind L2 state
		return &blockIngestionType{latest: isLatest, fork: true, preGenesis: false}, nil
	}

	// this is the typical, happy-path case. The ingested block's parent was the previously ingested block.
	return &blockIngestionType{latest: isLatest, fork: false, preGenesis: false}, nil
}

func (rc *RollupChain) noBlockStateBlockSubmissionResponse(block *types.Block) *common.BlockSubmissionResponse {
	return &common.BlockSubmissionResponse{
		BlockHeader:  block.Header(),
		FoundNewHead: false,
	}
}

func (rc *RollupChain) newBlockSubmissionResponse(bs *obscurocore.BlockState, rollup common.ExtRollup, logs map[gethrpc.ID][]byte) *common.BlockSubmissionResponse {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		rc.logger.Crit(msgNoRollup)
	}

	headBlock, err := rc.storage.FetchBlock(bs.Block)
	if err != nil {
		rc.logger.Crit("could not fetch block", log.ErrKey, err)
	}

	var head *common.Header
	if bs.FoundNewRollup {
		head = headRollup.Header
	}
	return &common.BlockSubmissionResponse{
		BlockHeader:          headBlock.Header(),
		ProducedRollup:       rollup,
		FoundNewHead:         bs.FoundNewRollup,
		IngestedRollupHeader: head,
		SubscribedLogs:       logs,
	}
}

//  STATE

// Recursively calculates and stores the block state, receipts and logs for the given block.
func (rc *RollupChain) updateState(b *types.Block) (*obscurocore.BlockState, error) {
	// This method is called recursively in case of re-orgs. Stop when state was calculated already.
	blockState, err := rc.storage.FetchBlockState(b.Hash())
	if err == nil {
		return blockState, nil
	}

	// If we get an error other than `ErrNotFound`, we return the error.
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block state. Cause: %w", err)
	}

	rollups := rc.bridge.ExtractRollups(b, rc.storage)
	genesisRollup, err := rc.storage.FetchGenesisRollup()
	// processing blocks before genesis, so there is nothing to do
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			if len(rollups) == 0 {
				return nil, nil //nolint:nilnil
			}
		} else {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
	}

	// Detect if the incoming block contains the genesis rollup, and generate an updated state.
	// Handles the case of the block containing the genesis being processed multiple times.
	genesisState, isGenesis := rc.handleGenesisRollup(b, rollups, genesisRollup)
	if isGenesis {
		return genesisState, nil
	}

	// To calculate the state after the current block, we need the state after the parent.
	// If this point is reached, there is a parent state guaranteed, because the genesis is handled above
	parentState, err := rc.storage.FetchBlockState(b.ParentHash())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not retrieve parent block state. Cause: %w", err)
		}
		// go back and calculate the Root of the Parent
		parent, err := rc.storage.FetchBlock(b.ParentHash())
		if err != nil {
			rc.logger.Crit("Could not retrieve parent block when calculating block state.", log.ErrKey, err)
		}
		parentState, err = rc.updateState(parent)
		if err != nil {
			return nil, err
		}
	}

	if parentState == nil {
		rc.logger.Crit(fmt.Sprintf("Could not calculate parent block state when calculating block state. BlockNum=%d. \n Block: %d, Block Parent: %d  ",
			b.Number(),
			common.ShortHash(b.Hash()),
			common.ShortHash(b.Header().ParentHash),
		))
	}

	blockState, head, receipts := rc.calculateBlockState(b, parentState, rollups)
	rc.logger.Trace(fmt.Sprintf("Calc block state b_%d: Found: %t - r_%d, ",
		common.ShortHash(b.Hash()),
		blockState.FoundNewRollup,
		common.ShortHash(blockState.HeadRollup)))

	logs := []*types.Log{}
	for _, receipt := range receipts {
		logs = append(logs, receipt.Logs...)
	}
	err = rc.storage.StoreNewHead(blockState, head, receipts, logs)
	if err != nil {
		rc.logger.Crit("Could not store new head.", log.ErrKey, err)
	}

	return blockState, nil
}

func (rc *RollupChain) handleGenesisRollup(b *types.Block, rollups []*obscurocore.Rollup, genesisRollup *obscurocore.Rollup) (genesisState *obscurocore.BlockState, isGenesis bool) {
	// the incoming block holds the genesis rollup
	// calculate and return the new block state
	// todo change this to an hardcoded hash on testnet/mainnet
	if genesisRollup == nil && len(rollups) == 1 {
		rc.logger.Info("Found genesis rollup", "l1Height", b.NumberU64(), "l1Hash", b.Hash())

		genesis := rollups[0]
		rc.storage.StoreGenesisRollup(genesis)

		// The genesis rollup is part of the canonical chain and will be included in an L1 block by the first Aggregator.
		bs := obscurocore.BlockState{
			Block:          b.Hash(),
			HeadRollup:     genesis.Hash(),
			FoundNewRollup: true,
		}
		err := rc.storage.StoreNewHead(&bs, genesis, nil, []*types.Log{})
		if err != nil {
			return nil, false
		}

		err = rc.faucet.CalculateGenesisState(rc.storage)
		if err != nil {
			return nil, false
		}

		return &bs, true
	}

	// Re-processing the block that contains the rollup. This can happen as blocks can be fed to the enclave multiple times.
	// In this case we don't update the state and move on.
	if genesisRollup != nil && len(rollups) == 1 && bytes.Equal(rollups[0].Header.Hash().Bytes(), genesisRollup.Hash().Bytes()) {
		return nil, true
	}
	return nil, false
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }

// This is where transactions are executed and the state is calculated.
// Obscuro includes a bridge embedded in the platform, and this method is responsible for processing deposits as well.
// The rollup can be a final rollup as received from peers or the rollup under construction.
func (rc *RollupChain) processState(rollup *obscurocore.Rollup, txs []*common.L2Tx, stateDB *state.StateDB) (gethcommon.Hash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, rollup.Header, rc.storage, rc.chainConfig, 0, rc.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			rc.logger.Crit("There should be an entry for each transaction ")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			rc.logger.Info(fmt.Sprintf("Excluding transaction %s from rollup r_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(rollup.Hash()), result))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	parent, found := rc.storage.ParentRollup(rollup)
	if !found {
		rc.logger.Crit("Sanity check. Rollup has no parent.")
	}

	parentProof, err := rc.storage.Proof(parent)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for rollup %s", rollup.Hash()), log.ErrKey, err)
	}
	rollupProof, err := rc.storage.Proof(rollup)
	if err != nil {
		rc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for rollup %s", rollup.Hash()), log.ErrKey, err)
	}

	depositTxs := rc.bridge.ExtractDeposits(
		parentProof,
		rollupProof,
		rc.storage,
		stateDB,
	)

	depositResponses := evm.ExecuteTransactions(depositTxs, stateDB, rollup.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
	depositReceipts := make([]*types.Receipt, len(depositResponses))
	if len(depositResponses) != len(depositTxs) {
		rc.logger.Crit("Sanity check. Some deposit transactions failed.")
	}
	i := 0
	for _, resp := range depositResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok {
			// TODO - Handle the case of an error (e.g. insufficient funds).
			rc.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
		}
		depositReceipts[i] = rec
		i++
	}

	onChainCallFunc := func(transactions common.L2Transactions) crosschain.OnChainEVMExecutorResponse {
		return evm.ExecuteTransactions(transactions, stateDB, rollup.Header, rc.storage, rc.chainConfig, len(executedTransactions), rc.logger)
	}

	offChainCallFunc := func(msg types.Message) (*core.ExecutionResult, error) {
		// Should this be deleted somehow?
		clonedDB := stateDB.Copy()
		return evm.ExecuteOffChainCall(&msg, clonedDB, rollup.Header, rc.storage, rc.chainConfig, rc.logger)
	}

	err = rc.crossChainProcessors.LocalManager.SubmitRemoteMessagesLocally(
		parentProof,
		rollupProof,
		stateDB,
		onChainCallFunc,
		offChainCallFunc,
	)
	if err != nil {
		rc.logger.Crit("[CrossChain] SubmitRemoteMessagesLocally failed!", log.ErrKey, err)
	}

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		rc.logger.Crit("could not commit to state DB. ", log.ErrKey, err)
	}

	sort.Sort(sortByTxIndex(txReceipts))
	sort.Sort(sortByTxIndex(depositReceipts))

	// todo - handle the tx execution logs
	return rootHash, executedTransactions, txReceipts, depositReceipts
}

func (rc *RollupChain) validateRollup(rollup *obscurocore.Rollup, rootHash gethcommon.Hash, txReceipts []*types.Receipt, depositReceipts []*types.Receipt, stateDB *state.StateDB) bool {
	h := rollup.Header
	if !bytes.Equal(rootHash.Bytes(), h.Root.Bytes()) {
		dump := strings.Replace(string(stateDB.Dump(&state.DumpConfig{})), "\n", "", -1)
		rollupBlock, err := rc.storage.Proof(rollup)
		if err != nil {
			rc.logger.Crit("Failed to extract rollup proof", log.ErrKey, err)
		}

		rc.logger.Error(fmt.Sprintf("[rollupBlock=%d] Verify rollup r_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(rollupBlock.Hash()), common.ShortHash(rollup.Hash()), rootHash, h.Root, h.Number, obscurocore.PrintTxs(rollup.Transactions), dump, depositReceipts))
		return false
	}

	//  check that the withdrawals in the header match the withdrawals as calculated
	withdrawals := rc.bridge.RollupPostProcessingWithdrawals(rollup, stateDB, toReceiptMap(txReceipts))
	for i, w := range withdrawals {
		hw := h.Withdrawals[i]
		if hw.Amount.Cmp(w.Amount) != 0 || hw.Recipient != w.Recipient || hw.Contract != w.Contract {
			rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Withdrawals don't match", common.ShortHash(rollup.Hash())))
			return false
		}
	}

	rec := allReceipts(txReceipts, depositReceipts)
	rbloom := types.CreateBloom(rec)
	if !bytes.Equal(rbloom.Bytes(), h.Bloom.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(rollup.Hash()), h.Bloom, rbloom))
		return false
	}

	receiptSha := types.DeriveSha(rec, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), h.ReceiptHash.Bytes()) {
		rc.logger.Error(fmt.Sprintf("Verify rollup r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(rollup.Hash()), h.ReceiptHash, receiptSha))
		return false
	}

	return true
}

// given an L1 block, and the State as it was in the Parent block, calculates the State after the current block.
func (rc *RollupChain) calculateBlockState(b *types.Block, parentState *obscurocore.BlockState, rollups []*obscurocore.Rollup) (*obscurocore.BlockState, *obscurocore.Rollup, []*types.Receipt) {
	currentHead, found := rc.storage.FetchRollup(parentState.HeadRollup)
	if !found {
		rc.logger.Crit("could not fetch parent rollup")
	}
	newHeadRollup, found := FindNextRollup(currentHead, rollups, rc.storage)
	var rollupTxReceipts []*types.Receipt
	var err error
	// only change the state if there is a new l2 HeadRollup in the current block
	if found {
		rollupTxReceipts, _, err = rc.checkRollup(newHeadRollup)
		// todo - this error needs to be surfaced to be used for the challenge
		if err != nil {
			rc.logger.Crit("Failed to check rollup", log.ErrKey, err)
			return nil, nil, nil
		}
	} else {
		newHeadRollup = currentHead
	}

	bs := obscurocore.BlockState{
		Block:          b.Hash(),
		HeadRollup:     newHeadRollup.Hash(),
		FoundNewRollup: found,
	}
	return &bs, newHeadRollup, rollupTxReceipts
}

// verifies that the headers of the rollup match the results of executing the transactions
func (rc *RollupChain) checkRollup(r *obscurocore.Rollup) ([]*types.Receipt, []*types.Receipt, error) { //nolint
	stateDB, err := rc.storage.CreateStateDB(r.Header.ParentHash)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the Rollup
	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(r, r.Transactions, stateDB)
	if len(successfulTxs) != len(r.Transactions) {
		return nil, nil, fmt.Errorf("all transactions that are included in a rollup must be executed")
	}

	isValid := rc.validateRollup(r, rootHash, txReceipts, depositReceipts, stateDB)
	if !isValid {
		return nil, nil, fmt.Errorf("invalid rollup")
	}

	// todo - check that the transactions hash to the header.txHash

	// verify the signature
	isValid = rc.verifySig(r)
	if !isValid {
		return nil, nil, fmt.Errorf("invalid signature")
	}

	return txReceipts, depositReceipts, nil
}

func toReceiptMap(txReceipts []*types.Receipt) map[gethcommon.Hash]*types.Receipt {
	result := make(map[gethcommon.Hash]*types.Receipt, 0)
	for _, r := range txReceipts {
		result[r.TxHash] = r
	}
	return result
}

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	receipts := make([]*types.Receipt, 0)
	receipts = append(receipts, txReceipts...)
	receipts = append(receipts, depositReceipts...)
	return receipts
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (rc *RollupChain) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, error) {
	rc.blockProcessingMutex.Lock()
	defer rc.blockProcessingMutex.Unlock()

	_, err := rc.storage.FetchBlock(block.Hash())
	if err == nil {
		return nil, rc.rejectBlockErr(errBlockAlreadyProcessed)
	}
	if !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	if rc.crossChainProcessors.Enabled() && !crosschain.VerifyReceiptHash(&block, receipts) {
		return nil, rc.rejectBlockErr(errors.New("receipts do not match receipt_root in block"))
	}

	ingestionType, err := rc.insertBlockIntoL1Chain(&block, isLatest)
	if err != nil {
		return nil, rc.rejectBlockErr(err)
	}
	rc.logger.Trace("block inserted successfully",
		"height", block.NumberU64(),
		"hash", block.Hash(),
		"ingestionType", ingestionType)

	// Only store the block if the L1 chain insertion succeeded
	rc.storage.StoreBlock(&block)

	// This requires block to be stored first ... but can permanently fail a block
	err = rc.crossChainProcessors.RemoteManager.ProcessCrossChainMessages(&block, receipts)
	if err != nil {
		return nil, rc.rejectBlockErr(errors.New("failed to process cross chain messages"))
	}

	rc.logger.Trace(fmt.Sprintf("Update state: b_%d", common.ShortHash(block.Hash())))
	// TODO - Handle error properly.
	blockState, _ := rc.updateState(&block)
	if blockState == nil {
		// not an error state, we ingested a block but no rollup head found
		return rc.noBlockStateBlockSubmissionResponse(&block), nil
	}

	logs := []*types.Log{}
	fetchedLogs, err := rc.storage.FetchLogs(block.Hash())
	if err == nil {
		logs = fetchedLogs
	} else {
		rc.logger.Error("Could not retrieve logs for stored block state; returning no logs. Cause: %w", err)
	}

	encryptedLogs, err := rc.subscriptionManager.GetSubscribedLogsEncrypted(logs, blockState.HeadRollup)
	if err != nil {
		rc.logger.Crit("Could not get subscribed logs in encrypted form. ", log.ErrKey, err)
	}

	// We do not produce a rollup if we're not an aggregator, or we're behind the L1 (in which case the rollup will be outdated).
	if !isLatest || rc.nodeType != common.Aggregator {
		return rc.newBlockSubmissionResponse(blockState, common.ExtRollup{}, encryptedLogs), nil
	}

	// As an aggregator on the head L1 block, we produce a rollup.
	r := rc.produceRollup(&block, blockState)
	rc.signRollup(r)
	// Sanity check the produced rollup
	_, _, err = rc.checkRollup(r)
	if err != nil {
		return nil, err
	}

	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	rc.storage.StoreRollup(r)
	rc.logger.Trace(fmt.Sprintf("Processed block: b_%d(%d). Produced rollup r_%d", common.ShortHash(block.Hash()), block.NumberU64(), common.ShortHash(r.Hash())))

	return rc.newBlockSubmissionResponse(blockState, r.ToExtRollup(rc.transactionBlobCrypto), encryptedLogs), nil
}

func (rc *RollupChain) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := rc.storage.FetchRollup(bs.HeadRollup)
	if !f {
		rc.logger.Crit(msgNoRollup)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs []*common.L2Tx
	var newRollupState *state.StateDB

	/*
			speculativeExecutionSucceeded := false
		   todo - reenable
			if e.speculativeExecutionEnabled {
				// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
				e.speculativeWorkInCh <- true
				speculativeRollup := <-e.speculativeWorkOutCh

				newRollupTxs = speculativeRollup.txs
				newRollupState = speculativeRollup.s
				newRollupHeader = speculativeRollup.h

				// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
				// or speculative execution is disabled
				speculativeExecutionSucceeded = speculativeRollup.found && (speculativeRollup.r.Hash() == bs.HeadRollup)

				if !speculativeExecutionSucceeded && speculativeRollup.r != nil {
					nodecommon.LogWithID(e.nodeShortID, "Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
						obscurocommon.ShortHash(speculativeRollup.r.Hash()),
						speculativeRollup.r.Header.Number,
						obscurocommon.ShortHash(bs.HeadRollup),
						headRollup.Header.Number)
					if e.statsCollector != nil {
						e.statsCollector.L2Recalc(e.nodeID)
					}
				}
			}
	*/

	// if !speculativeExecutionSucceeded {
	// In case the speculative execution thread has not succeeded in producing a valid rollup
	// we have to create a new one from the mempool transactions
	// Create a new rollup based on the fromBlock of inclusion of the previous, including all new transactions
	nonce := common.GenerateNonce()
	r, err := obscurocore.EmptyRollup(rc.hostID, headRollup.Header, b.Hash(), nonce)
	if err != nil {
		rc.logger.Crit("could not create rollup", log.ErrKey, err)
		return nil
	}

	newRollupTxs = rc.mempool.CurrentTxs(headRollup, rc.storage)
	newRollupState, err = rc.storage.CreateStateDB(r.Header.ParentHash)
	if err != nil {
		rc.logger.Crit("could not create stateDB", log.ErrKey, err)
		return nil
	}

	rootHash, successfulTxs, txReceipts, depositReceipts := rc.processState(r, newRollupTxs, newRollupState)

	r.Header.Root = rootHash
	r.Transactions = successfulTxs

	// Postprocessing - withdrawals
	txReceiptsMap := toReceiptMap(txReceipts)
	r.Header.Withdrawals = rc.bridge.RollupPostProcessingWithdrawals(r, newRollupState, txReceiptsMap)
	crossChainMessages, err := rc.crossChainProcessors.LocalManager.ExtractLocalMessages(txReceipts)
	if err != nil {
		rc.logger.Crit("[CrossChain] Extracting messages L2->L1 failed", err)
	}

	r.Header.CrossChainMessages = crossChainMessages

	rc.logger.Trace(fmt.Sprintf("[CrossChain] Added %d cross chain messages to rollup. Equivalent withdrawals in header - %d", len(r.Header.CrossChainMessages), len(r.Header.Withdrawals)))

	crossChainBind, err := rc.storage.Proof(r)
	if err != nil {
		rc.logger.Crit("Failed to extract rollup proof that should exist!")
	}

	r.Header.LatestInboudCrossChainHash = crossChainBind.Hash()
	r.Header.LatestInboundCrossChainHeight = crossChainBind.Number()

	receipts := allReceipts(txReceipts, depositReceipts)
	if len(receipts) == 0 {
		r.Header.ReceiptHash = types.EmptyRootHash
	} else {
		r.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		r.Header.Bloom = types.CreateBloom(receipts)
	}

	if len(successfulTxs) == 0 {
		r.Header.TxHash = types.EmptyRootHash
	} else {
		r.Header.TxHash = types.DeriveSha(types.Transactions(successfulTxs), trie.NewStackTrie(nil))
	}

	rc.logger.Trace("Create rollup.",
		"State", gethlog.Lazy{Fn: func() string {
			return strings.Replace(string(newRollupState.Dump(&state.DumpConfig{})), "\n", "", -1)
		}},
	)

	return r
}

func (rc *RollupChain) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	paramBytes, err := rc.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_call request. Cause: %w", err)
	}

	// extract params from byte slice to array of strings
	var paramList []interface{}
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("unable to decode eth_call params - %w", err)
	}

	// params are [callMsg, block number (optional) ]
	if len(paramList) < 1 {
		return nil, fmt.Errorf("required at least 1 params, but received %d", len(paramList))
	}

	apiArgs, err := gethencoding.ExtractEthCall(paramList[0])
	if err != nil {
		return nil, fmt.Errorf("unable to decode EthCall Params - %w", err)
	}

	// TODO Hook up the blockNumber
	result, err := rc.ExecuteOffChainTransactionAtBlock(apiArgs, gethrpc.BlockNumber(0))
	if err != nil {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, err.Error())
		return nil, err
	}
	// the execution might have succeeded but the evm contract logic might have failed
	if result.Failed() {
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, result.Err)
		return nil, result.Err
	}

	rc.logger.Trace(fmt.Sprintf("!OffChain result: %s", hexutils.BytesToHex(result.ReturnData)))

	var encodedResult string
	if len(result.ReturnData) != 0 {
		encodedResult = hexutil.Encode(result.ReturnData)
	}
	encryptedResult, err := rc.rpcEncryptionManager.EncryptWithViewingKey(*apiArgs.From, []byte(encodedResult))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_call request. Cause: %w", err)
	}

	return encryptedResult, nil
}

func (rc *RollupChain) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	// We decrypt the request.
	paramBytes, err := rc.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_getBalance request. Cause: %w", err)
	}

	// We extract the params from the request.
	var paramList []string
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
	}
	if len(paramList) != 2 {
		return nil, fmt.Errorf("required exactly two params, but received %d", len(paramList))
	}
	// TODO - Replace all usages of `HexToAddress` with a `SafeHexToAddress` that checks that the string does not exceed 20 bytes.
	accountAddress := gethcommon.HexToAddress(paramList[0])
	blockNumber := gethrpc.BlockNumber(0)
	err = blockNumber.UnmarshalJSON([]byte(paramList[1]))
	if err != nil {
		return nil, fmt.Errorf("could not parse requested rollup number - %w", err)
	}

	// get account balance at certain block/height
	balance, err := rc.GetBalanceAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch balance for account - %w", err)
	}

	// check if account is a contract
	isAddrContract, err := rc.IsAccountContractAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch balance for account - %w", err)
	}

	// We encrypt the result.
	address := accountAddress
	// If the accountAddress is a contract, encrypt with the address of the contract owner
	if isAddrContract {
		txHash, err := rc.storage.GetContractCreationTx(accountAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve tx that created contract %s. Cause %w", accountAddress.Hex(), err)
		}
		transaction, _, _, _, err := rc.storage.GetTransaction(*txHash)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve tx that created contract %s. Cause %w", accountAddress.Hex(), err)
		}
		signer := types.NewLondonSigner(rc.chainConfig.ChainID)

		sender, err := signer.Sender(transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to verify signature. Cause %w", err)
		}
		address = sender
	}

	encryptedBalance, err := rc.rpcEncryptionManager.EncryptWithViewingKey(address, []byte(balance.String()))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getBalance request. Cause: %w", err)
	}
	return encryptedBalance, nil
}

// GetChainStateAtBlock is a helper function that returns the state of the chain at height
// TODO make this cacheable
func (rc *RollupChain) GetChainStateAtBlock(blockNumber gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the rollup of interest.
	rollup, err := rc.getRollup(blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := rc.storage.CreateStateDB(rollup.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for rollup %s", rollup.Hash().Hex())
	}

	return blockchainState, err
}

// GetBalanceAtBlock returns the balance of an account at a certain height
func (rc *RollupChain) GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := rc.GetChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr)), nil
}

// IsAccountContractAtBlock returns the whether the account is a contract or not at a certain height
func (rc *RollupChain) IsAccountContractAtBlock(accountAddr gethcommon.Address, blockNumber gethrpc.BlockNumber) (bool, error) {
	chainState, err := rc.GetChainStateAtBlock(blockNumber)
	if err != nil {
		return false, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return len(chainState.GetCode(accountAddr)) > 0, nil
}

func (rc *RollupChain) ExecuteOffChainTransactionAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber gethrpc.BlockNumber) (*core.ExecutionResult, error) {
	// TODO review this during gas mechanics implementation
	callMsg, err := apiArgs.ToMessage(rc.GlobalGasCap, rc.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	hs, err := rc.storage.FetchHeadState()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state. Cause: %w", err)
	}
	// todo - get the parent
	r, f := rc.storage.FetchRollup(hs.HeadRollup)
	if !f {
		return nil, fmt.Errorf("unable to fetch head state rollup")
	}

	rc.logger.Trace(fmt.Sprintf("!OffChain call: contractAddress=%s, from=%s, data=%s, rollup=r_%d, state=%s", callMsg.To(), callMsg.From(), hexutils.BytesToHex(callMsg.Data()), common.ShortHash(r.Hash()), r.Header.Root.Hex()))
	s, err := rc.storage.CreateStateDB(hs.HeadRollup)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	result, err := evm.ExecuteOffChainCall(&callMsg, s, r.Header, rc.storage, rc.chainConfig, rc.logger)
	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		return result, err
	}

	// the execution outcome was unsuccessful, but it was able to execute the call
	if result.Failed() {
		// do not return an error
		// the result object should be evaluated upstream
		rc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", callMsg.To()), log.ErrKey, result.Err)
	}

	return result, nil
}

func (rc *RollupChain) signRollup(r *obscurocore.Rollup) {
	var err error
	h := r.Hash()
	r.Header.R, r.Header.S, err = ecdsa.Sign(rand.Reader, rc.enclavePrivateKey, h[:])
	if err != nil {
		rc.logger.Crit("Could not sign rollup. ", log.ErrKey, err)
	}
}

func (rc *RollupChain) verifySig(r *obscurocore.Rollup) bool {
	// If this rollup is generated by the current enclave skip the sig verification
	if bytes.Equal(r.Header.Agg.Bytes(), rc.hostID.Bytes()) {
		return true
	}

	h := r.Hash()
	if r.Header.R == nil || r.Header.S == nil {
		rc.logger.Error("Missing signature on rollup")
		return false
	}

	pubKey, err := rc.storage.FetchAttestedKey(r.Header.Agg)
	if err != nil {
		rc.logger.Error("Could not retrieve attested key for aggregator %s. Cause: %w", r.Header.Agg, err)
		return false
	}

	return ecdsa.Verify(pubKey, h[:], r.Header.R, r.Header.S)
}

// Retrieves the rollup with the given height, with special handling for earliest/latest/pending .
func (rc *RollupChain) getRollup(height gethrpc.BlockNumber) (*obscurocore.Rollup, error) {
	var rollup *obscurocore.Rollup
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisRollup, err := rc.storage.FetchGenesisRollup()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		rollup = genesisRollup
	case gethrpc.PendingBlockNumber:
		// TODO - Depends on the current pending rollup; leaving it for a different iteration as it will need more thought.
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		headState, err := rc.storage.FetchHeadState()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve head state. Cause: %w", err)
		}
		var found bool
		rollup, found = rc.storage.FetchRollup(headState.HeadRollup)
		if !found {
			return nil, fmt.Errorf("rollup with requested height %d was not found", height)
		}
	default:
		maybeRollup, found := rc.storage.FetchRollupByHeight(uint64(height))
		if !found {
			return nil, fmt.Errorf("rollup with requested height %d was not found", height)
		}
		rollup = maybeRollup
	}
	return rollup, nil
}

func (rc *RollupChain) rejectBlockErr(cause error) *common.BlockRejectError {
	var hash gethcommon.Hash
	l1Head, err := rc.storage.FetchHeadBlock()
	// TODO - Handle error.
	if err == nil {
		hash = l1Head.Hash()
	}
	return &common.BlockRejectError{
		L1Head:  hash,
		Wrapped: cause,
	}
}
