package l2chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupextractor"
	"github.com/status-im/keycard-go/hexutils"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// L2Chain represents the canonical chain, and manages the state.
type L2Chain struct {
	hostID      gethcommon.Address
	nodeType    common.NodeType
	chainConfig *params.ChainConfig
	sequencerID gethcommon.Address

	storage              db.Storage
	l1Blockchain         *gethcore.BlockChain
	rollupExtractor      *rollupextractor.RollupExtractor
	mempool              mempool.Manager
	genesis              *genesis.Genesis
	crossChainProcessors *crosschain.Processors

	enclavePrivateKey    *ecdsa.PrivateKey // this is a key known only to the current enclave, and the public key was shared with everyone during attestation
	blockProcessingMutex sync.Mutex
	logger               gethlog.Logger

	// Gas usage values
	// TODO use the ethconfig.Config instead
	GlobalGasCap uint64
	BaseFee      *big.Int
}

func New(
	hostID gethcommon.Address,
	nodeType common.NodeType,
	storage db.Storage,
	l1Blockchain *gethcore.BlockChain,
	bridge *rollupextractor.RollupExtractor,
	crossChainProcessors *crosschain.Processors,
	mempool mempool.Manager,
	privateKey *ecdsa.PrivateKey,
	chainConfig *params.ChainConfig,
	sequencerID gethcommon.Address,
	genesis *genesis.Genesis,
	logger gethlog.Logger,
) *L2Chain {
	return &L2Chain{
		hostID:               hostID,
		nodeType:             nodeType,
		storage:              storage,
		l1Blockchain:         l1Blockchain,
		rollupExtractor:      bridge,
		mempool:              mempool,
		crossChainProcessors: crossChainProcessors,
		enclavePrivateKey:    privateKey,
		chainConfig:          chainConfig,
		blockProcessingMutex: sync.Mutex{},
		logger:               logger,
		GlobalGasCap:         5_000_000_000,
		BaseFee:              gethcommon.Big0,
		sequencerID:          sequencerID,
		genesis:              genesis,
	}
}

// ProcessL1Block is used to update the enclave with an additional L1 block.
func (lc *L2Chain) ProcessL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.L2RootHash, *core.Batch, *core.Rollup, error) {
	lc.blockProcessingMutex.Lock()
	defer lc.blockProcessingMutex.Unlock()

	// We update the L1 chain state.
	err := lc.updateL1State(block, receipts, isLatest)
	if err != nil {
		return nil, nil, nil, err
	}

	// We update the L1 and L2 chain heads.
	headBatchHash, producedBatch, producedRollup, err := lc.updateL1AndL2Heads(&block, isLatest)
	if err != nil {
		return nil, nil, nil, err
	}

	return headBatchHash, producedBatch, producedRollup, nil
}

// UpdateL2Chain updates the L2 chain based on the received batch.
func (lc *L2Chain) UpdateL2Chain(batch *core.Batch) error {
	lc.blockProcessingMutex.Lock()
	defer lc.blockProcessingMutex.Unlock()

	if err := lc.checkAndStoreBatch(batch); err != nil {
		return err
	}

	// If this is the genesis batch, we commit the genesis state.
	if batch.IsGenesis() {
		if err := lc.genesis.CommitGenesisState(lc.storage); err != nil {
			return fmt.Errorf("could not apply genesis state. Cause: %w", err)
		}
	}

	return nil
}

func (lc *L2Chain) GetBalance(accountAddress gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, *hexutil.Big, error) {
	// get account balance at certain block/height
	balance, err := lc.GetBalanceAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// check if account is a contract
	isAddrContract, err := lc.isAccountContractAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// Decide which address to encrypt the result with
	address := accountAddress
	// If the accountAddress is a contract, encrypt with the address of the contract owner
	if isAddrContract {
		txHash, err := lc.storage.GetContractCreationTx(accountAddress)
		if err != nil {
			return nil, nil, err
		}
		transaction, _, _, _, err := lc.storage.GetTransaction(*txHash)
		if err != nil {
			return nil, nil, err
		}
		signer := types.NewLondonSigner(lc.chainConfig.ChainID)

		sender, err := signer.Sender(transaction)
		if err != nil {
			return nil, nil, err
		}
		address = sender
	}

	return &address, balance, nil
}

// GetBalanceAtBlock returns the balance of an account at a certain height
func (lc *L2Chain) GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := lc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr)), nil
}

// ExecuteOffChainTransaction executes non-state changing transactions at a given block height (eth_call)
func (lc *L2Chain) ExecuteOffChainTransaction(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	result, err := lc.ExecuteOffChainTransactionAtBlock(apiArgs, blockNumber)
	if err != nil {
		lc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, err.Error())
		return nil, err
	}

	// the execution might have succeeded (err == nil) but the evm contract logic might have failed (result.Failed() == true)
	if result.Failed() {
		lc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", apiArgs.To), log.ErrKey, result.Err)
		return nil, result.Err
	}

	lc.logger.Trace(fmt.Sprintf("!OffChain result: %s", hexutils.BytesToHex(result.ReturnData)))

	return result, nil
}

func (lc *L2Chain) ExecuteOffChainTransactionAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	// TODO review this during gas mechanics implementation
	callMsg, err := apiArgs.ToMessage(lc.GlobalGasCap, lc.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	// fetch the chain state at given batch
	blockState, err := lc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	batch, err := lc.getBatch(*blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state batch. Cause: %w", err)
	}

	lc.logger.Trace(
		fmt.Sprintf("!OffChain call: contractAddress=%s, from=%s, data=%s, batch=b_%d, state=%s",
			callMsg.To(),
			callMsg.From(),
			hexutils.BytesToHex(callMsg.Data()),
			common.ShortHash(*batch.Hash()),
			batch.Header.Root.Hex()),
	)

	result, err := evm.ExecuteOffChainCall(&callMsg, blockState, batch.Header, lc.storage, lc.chainConfig, lc.logger)
	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		return result, err
	}

	// the execution outcome was unsuccessful, but it was able to execute the call
	if result.Failed() {
		// do not return an error
		// the result object should be evaluated upstream
		lc.logger.Error(fmt.Sprintf("!OffChain: Failed to execute contract %s.", callMsg.To()), log.ErrKey, result.Err)
	}

	return result, nil
}

func (lc *L2Chain) updateL1State(block types.Block, receipts types.Receipts, isLatest bool) error {
	// We check whether we've already processed the block.
	_, err := lc.storage.FetchBlock(block.Hash())
	if err == nil {
		return common.ErrBlockAlreadyProcessed
	}
	if !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// Reject block if not provided with matching receipts.
	// This needs to happen before saving the block as otherwise it will be considered as processed.
	if lc.crossChainProcessors.Enabled() && !crosschain.VerifyReceiptHash(&block, receipts) {
		return errors.New("receipts do not match receipt_root in block")
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, err := lc.insertBlockIntoL1Chain(&block, isLatest)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return err
	}
	lc.logger.Trace("block inserted successfully",
		"height", block.NumberU64(), "hash", block.Hash(), "ingestionType", ingestionType)

	lc.storage.StoreBlock(&block)

	// This requires block to be stored first ... but can permanently fail a block
	err = lc.crossChainProcessors.Remote.StoreCrossChainMessages(&block, receipts)
	if err != nil {
		return errors.New("failed to process cross chain messages")
	}

	return nil
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block
// note: this method shouldn't be called for blocks we've seen before
func (lc *L2Chain) insertBlockIntoL1Chain(block *types.Block, isLatest bool) (*blockIngestionType, error) {
	if lc.l1Blockchain != nil {
		_, err := lc.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			return nil, fmt.Errorf("block was invalid: %w", err)
		}
	}
	// todo: this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := lc.storage.FetchHeadBlock()

	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo: we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &blockIngestionType{latest: isLatest, fork: false, preGenesis: true}, nil
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)

		// we do a basic sanity check, comparing the received block to the head block on the chain
	} else if block.ParentHash() != prevL1Head.Hash() {
		lcaBlock, err := gethutil.LCA(block, prevL1Head, lc.storage)
		if err != nil {
			return nil, common.ErrBlockAncestorNotFound
		}
		lc.logger.Trace("parent not found",
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
			lc.logger.Error("unexpected blockchain state, incoming block is not child of L1 head and not an earlier fork of L1 head",
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

// Updates the L1 and L2 chain heads, and returns the new head batch hash and the produced batch and rollup (if applicable).
func (lc *L2Chain) updateL1AndL2Heads(block *types.Block, isLatestBlock bool) (*common.L2RootHash, *core.Batch, *core.Rollup, error) {
	// We process the rollups, updating the head rollup associated with the L1 block as we go.
	if err := lc.processRollups(block); err != nil {
		// TODO - #718 - Determine correct course of action if one or more rollups are invalid.
		lc.logger.Error("could not process rollups", log.ErrKey, err)
	}

	// We determine whether we have produced a genesis batch yet.
	genesisBatchStored := true
	headBatch, err := lc.storage.FetchHeadBatchForBlock(block.ParentHash())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, nil, nil, fmt.Errorf("could not retrieve current head batch. Cause: %w", err)
		}
		genesisBatchStored = false
	}

	// If there is an L2 head, we retrieve its stored receipts.
	var l2HeadTxReceipts types.Receipts
	if genesisBatchStored {
		if l2HeadTxReceipts, err = lc.storage.GetReceiptsByHash(*headBatch.Hash()); err != nil {
			return nil, nil, nil, fmt.Errorf("could not fetch batch receipts. Cause: %w", err)
		}
	}

	// If we're the sequencer and we're on the latest block, we produce a new L2 head to replace the old one.
	var producedBatch *core.Batch
	var producedRollup *core.Rollup
	if lc.nodeType == common.Sequencer && isLatestBlock {
		producedBatch, l2HeadTxReceipts, err = lc.produceAndStoreBatch(block, genesisBatchStored)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not produce and store new batch. Cause: %w", err)
		}
		headBatch = producedBatch

		if producedRollup, err = lc.produceAndStoreRollup(block, producedBatch); err != nil {
			return nil, nil, nil, fmt.Errorf("could not produce and store new rollup. Cause: %w", err)
		}
	}

	// We update the L1 and L2 chain heads.
	if headBatch != nil {
		if err = lc.storage.UpdateHeadBatch(block.Hash(), headBatch, l2HeadTxReceipts); err != nil {
			return nil, nil, nil, fmt.Errorf("could not store new head. Cause: %w", err)
		}
		if err = lc.storage.UpdateL1Head(block.Hash()); err != nil {
			return nil, nil, nil, fmt.Errorf("could not store new L1 head. Cause: %w", err)
		}
	}

	var headBatchHash *gethcommon.Hash
	if headBatch != nil {
		headBatchHash = headBatch.Hash()
	}
	return headBatchHash, producedBatch, producedRollup, nil
}

// Produces a new batch, signs it and stores it.
func (lc *L2Chain) produceAndStoreBatch(block *common.L1Block, genesisBatchStored bool) (*core.Batch, types.Receipts, error) {
	l2Head, err := lc.produceBatch(block, genesisBatchStored)
	if err != nil {
		return nil, nil, fmt.Errorf("could not produce batch. Cause: %w", err)
	}

	if err = lc.signBatch(l2Head); err != nil {
		return nil, nil, fmt.Errorf("could not sign batch. Cause: %w", err)
	}

	l2HeadTxReceipts, err := lc.getTxReceipts(l2Head)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get batch transaction receipts. Cause: %w", err)
	}
	if err = lc.storage.StoreBatch(l2Head, l2HeadTxReceipts); err != nil {
		return nil, nil, fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	return l2Head, l2HeadTxReceipts, nil
}

// Produces a new rollup, signs it and stores it.
func (lc *L2Chain) produceAndStoreRollup(block *types.Block, producedBatch *core.Batch) (*core.Rollup, error) {
	parentRollup, err := lc.storage.FetchHeadRollupForBlock(&block.Header().ParentHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, fmt.Errorf("could not retrieve head rollup. Cause: %w", err)
	}
	var parentRollupHeader *common.RollupHeader
	if parentRollup != nil {
		parentRollupHeader = parentRollup.Header
	}

	producedRollup := &core.Rollup{
		Header:  producedBatch.Header.ToRollupHeader(parentRollupHeader),
		Batches: []*core.Batch{producedBatch},
	}

	if err = lc.signRollup(producedRollup); err != nil {
		return nil, fmt.Errorf("could not sign produced rollup. Cause: %w", err)
	}
	if err = lc.storage.StoreRollup(producedRollup); err != nil {
		return nil, fmt.Errorf("could not store produced rollup. Cause: %w", err)
	}
	blockHash := block.Hash()
	if err = lc.storage.UpdateHeadRollup(&blockHash, producedRollup.Hash()); err != nil {
		return nil, fmt.Errorf("could not update head rollup for produced rollup. Cause: %w", err)
	}

	return producedRollup, nil
}

// Creates a genesis batch linked to the provided L1 block and signs it.
func (lc *L2Chain) produceGenesisBatch(blkHash common.L1RootHash) (*core.Batch, error) {
	preFundGenesisState, err := lc.genesis.GetGenesisRoot(lc.storage)
	if err != nil {
		return nil, err
	}

	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			Agg:         lc.hostID,
			ParentHash:  common.L2RootHash{},
			L1Proof:     blkHash,
			Root:        *preFundGenesisState,
			TxHash:      types.EmptyRootHash,
			Number:      big.NewInt(int64(0)),
			ReceiptHash: types.EmptyRootHash,
			Time:        uint64(time.Now().Unix()),
		},
		Transactions: []*common.L2Tx{},
	}

	// TODO: Figure out a better way to bootstrap the system contracts.
	deployTx, err := lc.crossChainProcessors.Local.GenerateMessageBusDeployTx()
	if err != nil {
		lc.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	// Add transaction to mempool so it gets processed when it can.
	// Should be the first transaction to be processed.
	if err := lc.mempool.AddMempoolTx(deployTx); err != nil {
		lc.logger.Crit("Cannot create synthetic transaction for deploying the message bus contract on :|")
	}

	if err = lc.genesis.CommitGenesisState(lc.storage); err != nil {
		return nil, fmt.Errorf("could not apply genesis preallocation. Cause: %w", err)
	}
	return genesisBatch, nil
}

// This is where transactions are executed and the state is calculated.
// Obscuro includes a message bus embedded in the platform, and this method is responsible for transferring messages as well.
// The batch can be a final batch as received from peers or the batch under construction.
func (lc *L2Chain) processState(batch *core.Batch, txs []*common.L2Tx, stateDB *state.StateDB) (common.L2RootHash, []*common.L2Tx, []*types.Receipt, []*types.Receipt) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, batch.Header, lc.storage, lc.chainConfig, 0, lc.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			lc.logger.Crit("There should be an entry for each transaction ")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			lc.logger.Info(fmt.Sprintf("Excluding transaction %s from batch b_%d. Cause: %s", tx.Hash().Hex(), common.ShortHash(*batch.Hash()), result))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the fromBlock of the parent to the current block (which is the fromBlock of the new rollup)
	parent, err := lc.storage.FetchBatch(batch.Header.ParentHash)
	if err != nil {
		lc.logger.Crit("Sanity check. Rollup has no parent.", log.ErrKey, err)
	}

	parentProof, err := lc.storage.FetchBlock(parent.Header.L1Proof)
	if err != nil {
		lc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}
	batchProof, err := lc.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil {
		lc.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", batch.Hash()), log.ErrKey, err)
	}

	messages := lc.crossChainProcessors.Local.RetrieveInboundMessages(parentProof, batchProof, stateDB)
	transactions := lc.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)
	syntheticTransactionsResponses := evm.ExecuteTransactions(transactions, stateDB, batch.Header, lc.storage, lc.chainConfig, len(executedTransactions), lc.logger)
	synthReceipts := make([]*types.Receipt, len(syntheticTransactionsResponses))
	if len(syntheticTransactionsResponses) != len(transactions) {
		lc.logger.Crit("Sanity check. Some synthetic transactions failed.")
	}

	i := 0
	for _, resp := range syntheticTransactionsResponses {
		rec, ok := resp.(*types.Receipt)
		if !ok { // Еxtract reason for failing deposit.
			// TODO - Handle the case of an error (e.g. insufficient funds).
			lc.logger.Crit("Sanity check. Expected a receipt", log.ErrKey, resp)
		}

		if rec.Status == 0 { // Synthetic transactions should not fail. In case of failure get the revert reason.
			failingTx := transactions[i]
			txCallMessage := types.NewMessage(
				lc.crossChainProcessors.Local.GetOwner(),
				failingTx.To(),
				stateDB.GetNonce(lc.crossChainProcessors.Local.GetOwner()),
				failingTx.Value(),
				failingTx.Gas(),
				gethcommon.Big0,
				gethcommon.Big0,
				gethcommon.Big0,
				failingTx.Data(),
				failingTx.AccessList(),
				false)

			clonedDB := stateDB.Copy()
			res, err := evm.ExecuteOffChainCall(&txCallMessage, clonedDB, batch.Header, lc.storage, lc.chainConfig, lc.logger)
			lc.logger.Crit("Synthetic transaction failed!", log.ErrKey, err, "result", res)
		}

		synthReceipts[i] = rec
		i++
	}

	rootHash, err := stateDB.Commit(true)
	if err != nil {
		lc.logger.Crit("could not commit to state DB. ", log.ErrKey, err)
	}

	sort.Sort(sortByTxIndex(txReceipts))

	// todo - handle the tx execution logs
	return rootHash, executedTransactions, txReceipts, synthReceipts
}

// Checks the internal validity of the batch.
func (lc *L2Chain) isInternallyValidBatch(batch *core.Batch) (types.Receipts, error) {
	stateDB, err := lc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the batch
	rootHash, executedTxs, txReceipts, depositReceipts := lc.processState(batch, batch.Transactions, stateDB)
	if len(executedTxs) != len(batch.Transactions) {
		return nil, fmt.Errorf("all transactions that are included in a batch must be executed")
	}

	// Check that the root hash in the header matches the root hash as calculated.
	if !bytes.Equal(rootHash.Bytes(), batch.Header.Root.Bytes()) {
		dump := strings.Replace(string(stateDB.Dump(&state.DumpConfig{})), "\n", "", -1)
		return nil, fmt.Errorf("verify batch b_%d: Calculated a different state. This should not happen as there are no malicious actors yet. \nGot: %s\nExp: %s\nHeight:%d\nTxs:%v\nState: %s.\nDeposits: %+v",
			common.ShortHash(*batch.Hash()), rootHash, batch.Header.Root, batch.Header.Number, core.PrintTxs(batch.Transactions), dump, depositReceipts)
	}

	// Check that the receipts bloom in the header matches the receipts bloom as calculated.
	receipts := allReceipts(txReceipts, depositReceipts)
	receiptBloom := types.CreateBloom(receipts)
	if !bytes.Equal(receiptBloom.Bytes(), batch.Header.Bloom.Bytes()) {
		return nil, fmt.Errorf("verify batch r_%d: Invalid bloom (remote: %x  local: %x)", common.ShortHash(*batch.Hash()), batch.Header.Bloom, receiptBloom)
	}

	// Check that the receipts SHA in the header matches the receipts SHA as calculated.
	receiptSha := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	if !bytes.Equal(receiptSha.Bytes(), batch.Header.ReceiptHash.Bytes()) {
		return nil, fmt.Errorf("verify batch r_%d: invalid receipt root hash (remote: %x local: %x)", common.ShortHash(*batch.Hash()), batch.Header.ReceiptHash, receiptSha)
	}

	// Check that the signature is valid.
	if err = lc.checkSequencerSignature(batch.Hash(), &batch.Header.Agg, batch.Header.R, batch.Header.S); err != nil {
		return nil, fmt.Errorf("verify batch r_%d: invalid signature. Cause: %w", common.ShortHash(*batch.Hash()), err)
	}

	// todo - check that the transactions hash to the header.txHash

	return txReceipts, nil
}

// Returns the receipts for the transactions in the batch.
func (lc *L2Chain) getTxReceipts(batch *core.Batch) ([]*types.Receipt, error) {
	if batch.IsGenesis() {
		return nil, nil
	}

	stateDB, err := lc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	// calculate the state to compare with what is in the batch
	_, _, txReceipts, _ := lc.processState(batch, batch.Transactions, stateDB) //nolint:dogsled
	return txReceipts, nil
}

func (lc *L2Chain) signBatch(batch *core.Batch) error {
	var err error
	h := batch.Hash()
	batch.Header.R, batch.Header.S, err = ecdsa.Sign(rand.Reader, lc.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign batch. Cause: %w", err)
	}
	return nil
}

func (lc *L2Chain) signRollup(rollup *core.Rollup) error {
	var err error
	h := rollup.Hash()
	rollup.Header.R, rollup.Header.S, err = ecdsa.Sign(rand.Reader, lc.enclavePrivateKey, h[:])
	if err != nil {
		return fmt.Errorf("could not sign rollup. Cause: %w", err)
	}
	return nil
}

// Checks that the header is signed validly by the sequencer.
func (lc *L2Chain) checkSequencerSignature(headerHash *gethcommon.Hash, aggregator *gethcommon.Address, sigR *big.Int, sigS *big.Int) error {
	// Batches and rollups should only be produced by the sequencer.
	// TODO - #718 - Sequencer identities should be retrieved from the L1 management contract.
	if !bytes.Equal(aggregator.Bytes(), lc.sequencerID.Bytes()) {
		return fmt.Errorf("expected batch to be produced by sequencer %s, but was produced by %s", lc.sequencerID.Hex(), aggregator.Hex())
	}

	if sigR == nil || sigS == nil {
		return fmt.Errorf("missing signature on batch")
	}

	pubKey, err := lc.storage.FetchAttestedKey(*aggregator)
	if err != nil {
		return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", aggregator, err)
	}

	if !ecdsa.Verify(pubKey, headerHash.Bytes(), sigR, sigS) {
		return fmt.Errorf("could not verify ECDSA signature")
	}
	return nil
}

// Retrieves the batch with the given height, with special handling for earliest/latest/pending .
func (lc *L2Chain) getBatch(height gethrpc.BlockNumber) (*core.Batch, error) {
	var batch *core.Batch
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisBatch, err := lc.storage.FetchBatchByHeight(0)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		batch = genesisBatch
	case gethrpc.PendingBlockNumber:
		// TODO - Depends on the current pending rollup; leaving it for a different iteration as it will need more thought.
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		headBatch, err := lc.storage.FetchHeadBatch()
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d was not found. Cause: %w", height, err)
		}
		batch = headBatch
	default:
		maybeBatch, err := lc.storage.FetchBatchByHeight(uint64(height))
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		batch = maybeBatch
	}
	return batch, nil
}

// Creates either a genesis or regular (i.e. post-genesis) batch.
func (lc *L2Chain) produceBatch(block *types.Block, genesisBatchStored bool) (*core.Batch, error) {
	// We handle producing the genesis batch as a special case.
	if !genesisBatchStored {
		return lc.produceGenesisBatch(block.Hash())
	}

	headBatch, err := lc.storage.FetchHeadBatchForBlock(block.ParentHash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)
	}

	// These variables will be used to create the new batch
	var newBatchTxs []*common.L2Tx
	var newBatchState *state.StateDB

	// Create a new batch based on the fromBlock of inclusion of the previous, including all new transactions
	batch, err := core.EmptyBatch(lc.hostID, headBatch.Header, block.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create batch. Cause: %w", err)
	}

	newBatchTxs, err = lc.mempool.CurrentTxs(headBatch, lc.storage)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve current transactions. Cause: %w", err)
	}

	newBatchState, err = lc.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	rootHash, successfulTxs, txReceipts, depositReceipts := lc.processState(batch, newBatchTxs, newBatchState)

	batch.Header.Root = rootHash
	batch.Transactions = successfulTxs

	crossChainMessages, err := lc.crossChainProcessors.Local.ExtractOutboundMessages(txReceipts)
	if err != nil {
		lc.logger.Crit("Extracting messages L2->L1 failed", err, log.CmpKey, log.CrossChainCmp)
	}

	batch.Header.CrossChainMessages = crossChainMessages

	lc.logger.Trace(fmt.Sprintf("Added %d cross chain messages to batch.",
		len(batch.Header.CrossChainMessages)), log.CmpKey, log.CrossChainCmp)

	crossChainBind, err := lc.storage.FetchBlock(batch.Header.L1Proof)
	if err != nil {
		lc.logger.Crit("Failed to extract batch proof that should exist!")
	}

	batch.Header.LatestInboundCrossChainHash = crossChainBind.Hash()
	batch.Header.LatestInboundCrossChainHeight = crossChainBind.Number()

	receipts := allReceipts(txReceipts, depositReceipts)
	if len(receipts) == 0 {
		batch.Header.ReceiptHash = types.EmptyRootHash
	} else {
		batch.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		batch.Header.Bloom = types.CreateBloom(receipts)
	}

	if len(successfulTxs) == 0 {
		batch.Header.TxHash = types.EmptyRootHash
	} else {
		batch.Header.TxHash = types.DeriveSha(types.Transactions(successfulTxs), trie.NewStackTrie(nil))
	}

	lc.logger.Trace("Create batch.",
		"State", gethlog.Lazy{Fn: func() string {
			return strings.Replace(string(newBatchState.Dump(&state.DumpConfig{})), "\n", "", -1)
		}},
	)

	return batch, nil
}

// Returns the state of the chain at height
// TODO make this cacheable
func (lc *L2Chain) getChainStateAtBlock(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the batch of interest.
	batch, err := lc.getBatch(*blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := lc.storage.CreateStateDB(*batch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for batch %s", batch.Hash().Hex())
	}

	return blockchainState, err
}

// Returns the whether the account is a contract or not at a certain height
func (lc *L2Chain) isAccountContractAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (bool, error) {
	chainState, err := lc.getChainStateAtBlock(blockNumber)
	if err != nil {
		return false, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return len(chainState.GetCode(accountAddr)) > 0, nil
}

// Validates and stores the rollup in a given block, and returns the new head rollup's hash.
// TODO - #718 - Design a mechanism to detect a case where the rollups never contain any batches (despite batches arriving via P2P).
func (lc *L2Chain) processRollups(block *common.L1Block) error {
	l1ParentHash := block.ParentHash()
	currentHeadRollup, err := lc.storage.FetchHeadRollupForBlock(&l1ParentHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not fetch current L2 head rollup")
	}

	rollups := lc.rollupExtractor.ExtractRollups(block, lc.storage)
	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.Number.Cmp(rollups[j].Header.Number) < 0
	})

	// If this is the first rollup we've ever received, we check that it's the genesis rollup.
	if currentHeadRollup == nil && len(rollups) != 0 && !rollups[0].IsGenesis() {
		return fmt.Errorf("received rollup with number %d but no genesis rollup is stored", rollups[0].Number())
	}

	for idx, rollup := range rollups {
		if err = lc.checkSequencerSignature(rollup.Hash(), &rollup.Header.Agg, rollup.Header.R, rollup.Header.S); err != nil {
			return fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		if !rollup.IsGenesis() {
			previousRollup := currentHeadRollup
			if idx != 0 {
				previousRollup = rollups[idx-1]
			}
			if err = lc.checkRollupsCorrectlyChained(rollup, previousRollup); err != nil {
				return err
			}
		}

		for _, batch := range rollup.Batches {
			if err = lc.checkAndStoreBatch(batch); err != nil {
				return fmt.Errorf("could not store batch. Cause: %w", err)
			}
		}

		if err = lc.storage.StoreRollup(rollup); err != nil {
			return fmt.Errorf("could not store rollup. Cause: %w", err)
		}
	}

	newHeadRollup := currentHeadRollup
	if len(rollups) > 0 {
		newHeadRollup = rollups[len(rollups)-1]
	}
	if newHeadRollup != nil && lc.nodeType != common.Sequencer {
		l1Head := block.Hash()
		if err = lc.storage.UpdateHeadRollup(&l1Head, newHeadRollup.Hash()); err != nil {
			return fmt.Errorf("could not update L2 head rollup. Cause: %w", err)
		}
	}

	return nil
}

// Checks that the rollup:
//   - Has a number exactly 1 higher than the previous rollup
//   - Links to the previous rollup by hash
//   - Has a first batch whose parent is the head batch of the previous rollup
func (lc *L2Chain) checkRollupsCorrectlyChained(rollup *core.Rollup, previousRollup *core.Rollup) error {
	if big.NewInt(0).Sub(rollup.Header.Number, previousRollup.Header.Number).Cmp(big.NewInt(1)) != 0 {
		return fmt.Errorf("found gap in rollups between rollup %d and rollup %d",
			previousRollup.Header.Number, rollup.Header.Number)
	}

	if rollup.Header.ParentHash != *previousRollup.Hash() {
		return fmt.Errorf("found gap in rollups. Rollup %d did not reference rollup %d by hash",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	if len(rollup.Batches) != 0 && previousRollup.Header.HeadBatchHash != rollup.Batches[0].Header.ParentHash {
		return fmt.Errorf("found gap in rollup batches. Batches in rollup %d did not chain to batches in rollup %d",
			rollup.Header.Number, previousRollup.Header.Number)
	}

	return nil
}

// Checks the batch. If we've not seen a batch at this height before, we store it. If we have seen a batch at this
// height before, we validate it against the other received batch at the same height.
func (lc *L2Chain) checkAndStoreBatch(batch *core.Batch) error {
	// We check the batch.
	var txReceipts types.Receipts
	// TODO - #718 - Determine what level of checking we should perform on the genesis batch.
	if !batch.IsGenesis() {
		var err error
		txReceipts, err = lc.isInternallyValidBatch(batch)
		if err != nil {
			return fmt.Errorf("batch was invalid. Cause: %w", err)
		}

		// We check that we've stored the batch's parent.
		if _, err = lc.storage.FetchBatch(batch.Header.ParentHash); err != nil {
			return fmt.Errorf("could not retrieve parent batch. Cause: %w", err)
		}
	}

	// If we've stored a batch at this height before, we ensure that it has the same transactions.
	_, err := lc.storage.FetchBatchByHeight(batch.NumberU64())
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("could not fetch batch. Cause: %w", err)
	}
	// TODO - #718 - Once the sequencer includes transactions deterministically (i.e. a batch of a given height always
	//  contains the same transactions, regardless of reorgs), uncomment this check.
	//if err == nil && batch.Header.TxHash != storedBatch.Header.TxHash {
	//	return fmt.Errorf("two batches at same height did not have the same transactions")
	//}

	// If we haven't stored the batch before, we store it and update the head batch for that L1 block.
	// TODO - FetchBatch should return errutil.ErrNotFound for unstored batches, so we can handle that type of error
	//  separately.
	if _, err = lc.storage.FetchBatch(*batch.Hash()); err != nil {
		if err = lc.storage.StoreBatch(batch, txReceipts); err != nil {
			return fmt.Errorf("failed to store batch. Cause: %w", err)
		}
		if err = lc.storage.UpdateHeadBatch(batch.Header.L1Proof, batch, txReceipts); err != nil {
			return fmt.Errorf("could not store new L2 head. Cause: %w", err)
		}
	}

	return nil
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	return append(txReceipts, depositReceipts...)
}
