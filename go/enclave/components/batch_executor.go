package components

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
)

// batchExecutor - the component responsible for executing batches
type batchExecutor struct {
	storage              storage.Storage
	crossChainProcessors *crosschain.Processors
	genesis              *genesis.Genesis
	logger               gethlog.Logger
	chainConfig          *params.ChainConfig

	// stateDBMutex - used to protect calls to stateDB.Commit as it is not safe for async access.
	stateDBMutex sync.Mutex
}

func NewBatchExecutor(storage storage.Storage, cc *crosschain.Processors, genesis *genesis.Genesis, chainConfig *params.ChainConfig, logger gethlog.Logger) BatchExecutor {
	return &batchExecutor{
		storage:              storage,
		crossChainProcessors: cc,
		genesis:              genesis,
		chainConfig:          chainConfig,
		logger:               logger,
		stateDBMutex:         sync.Mutex{},
	}
}

func (executor *batchExecutor) ComputeBatch(context *BatchExecutionContext) (*ComputedBatch, error) {
	defer executor.logger.Info("Batch context processed", log.DurationKey, measure.NewStopwatch())

	// Block is loaded first since if its missing this batch might be based on l1 fork we dont know about
	// and we want to filter out all fork batches based on not knowing the l1 block
	block, err := executor.storage.FetchBlock(context.BlockPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		return nil, errutil.ErrBlockForBatchNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %s for batch. Cause: %w", context.BlockPtr, err)
	}

	// These variables will be used to create the new batch
	parent, err := executor.storage.FetchBatch(context.ParentPtr)
	if errors.Is(err, errutil.ErrNotFound) {
		return nil, errutil.ErrAncestorBatchNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve parent batch %s. Cause: %w", context.ParentPtr, err)
	}

	parentBlock := block
	if parent.Header.L1Proof != block.Hash() {
		var err error
		parentBlock, err = executor.storage.FetchBlock(parent.Header.L1Proof)
		if err != nil {
			executor.logger.Crit(fmt.Sprintf("Could not retrieve a proof for batch %s", parent.Hash()), log.ErrKey, err)
		}
	}

	// Create a new batch based on the fromBlock of inclusion of the previous, including all new transactions
	batch := core.DeterministicEmptyBatch(parent.Header, block, context.AtTime, context.SequencerNo)

	stateDB, err := executor.storage.CreateStateDB(batch.Header.ParentHash)
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	var messages common.CrossChainMessages
	// Cross chain data is not accessible until one after the genesis batch
	if context.SequencerNo.Int64() > int64(common.L2GenesisSeqNo+1) {
		messages = executor.crossChainProcessors.Local.RetrieveInboundMessages(parentBlock, block, stateDB)
	}
	crossChainTransactions := executor.crossChainProcessors.Local.CreateSyntheticTransactions(messages, stateDB)

	successfulTxs, txReceipts, err := executor.processTransactions(batch, 0, context.Transactions, stateDB, context.ChainConfig)
	if err != nil {
		return nil, fmt.Errorf("could not process transactions. Cause: %w", err)
	}

	ccSuccessfulTxs, ccReceipts, err := executor.processTransactions(batch, len(successfulTxs), crossChainTransactions, stateDB, context.ChainConfig)
	if err != nil {
		return nil, err
	}

	if err = executor.verifyInboundCrossChainTransactions(crossChainTransactions, ccSuccessfulTxs, ccReceipts); err != nil {
		return nil, fmt.Errorf("batch computation failed due to cross chain messages. Cause: %w", err)
	}

	// we need to copy the batch to reset the internal hash cache
	copyBatch := *batch
	copyBatch.Header.Root = stateDB.IntermediateRoot(false)
	copyBatch.Transactions = successfulTxs
	copyBatch.ResetHash()

	if err = executor.populateOutboundCrossChainData(&copyBatch, block, txReceipts); err != nil {
		return nil, fmt.Errorf("failed adding cross chain data to batch. Cause: %w", err)
	}

	executor.populateHeader(&copyBatch, allReceipts(txReceipts, ccReceipts))

	// the logs and receipts produced by the EVM have the wrong hash which must be adjusted
	for _, receipt := range txReceipts {
		receipt.BlockHash = copyBatch.Hash()
		for _, l := range receipt.Logs {
			l.BlockHash = copyBatch.Hash()
		}
	}

	return &ComputedBatch{
		Batch:    &copyBatch,
		Receipts: txReceipts,
		Commit: func(deleteEmptyObjects bool) (gethcommon.Hash, error) {
			executor.stateDBMutex.Lock()
			defer executor.stateDBMutex.Unlock()
			h, err := stateDB.Commit(deleteEmptyObjects)
			if err != nil {
				return gethcommon.Hash{}, err
			}
			trieDB := executor.storage.TrieDB()
			err = trieDB.Commit(h, true, nil)
			return h, err
		},
	}, nil
}

func (executor *batchExecutor) ExecuteBatch(batch *core.Batch) (types.Receipts, error) {
	defer executor.logger.Info("Executed batch", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	// Validators recompute the entire batch using the same batch context
	// if they have all necessary prerequisites like having the l1 block processed
	// and the parent hash. This recomputed batch is then checked against the incoming batch.
	// If the sequencer has tampered with something the hash will not add up and validation will
	// produce an error.
	cb, err := executor.ComputeBatch(&BatchExecutionContext{
		BlockPtr:     batch.Header.L1Proof,
		ParentPtr:    batch.Header.ParentHash,
		Transactions: batch.Transactions,
		AtTime:       batch.Header.Time,
		ChainConfig:  executor.chainConfig,
		SequencerNo:  batch.Header.SequencerOrderNo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed computing batch %s. Cause: %w", batch.Hash(), err)
	}

	if cb.Batch.Hash() != batch.Hash() {
		// todo @stefan - generate a validator challenge here and return it
		executor.logger.Error(fmt.Sprintf("Error validating batch. Calculated: %+v\n Incoming: %+v\n", cb.Batch.Header, batch.Header))
		return nil, fmt.Errorf("batch is in invalid state. Incoming hash: %s  Computed hash: %s", batch.Hash(), cb.Batch.Hash())
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("cannot commit stateDB for incoming valid batch %s. Cause: %w", batch.Hash(), err)
	}

	return cb.Receipts, nil
}

func (executor *batchExecutor) CreateGenesisState(blkHash common.L1BlockHash, timeNow uint64) (*core.Batch, *types.Transaction, error) {
	preFundGenesisState, err := executor.genesis.GetGenesisRoot(executor.storage)
	if err != nil {
		return nil, nil, err
	}

	genesisBatch := &core.Batch{
		Header: &common.BatchHeader{
			ParentHash:       common.L2BatchHash{},
			L1Proof:          blkHash,
			Root:             *preFundGenesisState,
			TxHash:           types.EmptyRootHash,
			Number:           big.NewInt(int64(0)),
			SequencerOrderNo: big.NewInt(int64(common.L2GenesisSeqNo)), // genesis batch has seq number 1
			ReceiptHash:      types.EmptyRootHash,
			Time:             timeNow,
		},
		Transactions: []*common.L2Tx{},
	}

	// todo (#1577) - figure out a better way to bootstrap the system contracts
	deployTx, err := executor.crossChainProcessors.Local.GenerateMessageBusDeployTx()
	if err != nil {
		executor.logger.Crit("Could not create message bus deployment transaction", "Error", err)
	}

	if err = executor.genesis.CommitGenesisState(executor.storage); err != nil {
		return nil, nil, fmt.Errorf("could not apply genesis preallocation. Cause: %w", err)
	}
	return genesisBatch, deployTx, nil
}

func (executor *batchExecutor) populateOutboundCrossChainData(batch *core.Batch, block *types.Block, receipts types.Receipts) error {
	crossChainMessages, err := executor.crossChainProcessors.Local.ExtractOutboundMessages(receipts)
	if err != nil {
		executor.logger.Error("Failed extracting L2->L1 messages", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
		return fmt.Errorf("could not extract cross chain messages. Cause: %w", err)
	}

	batch.Header.CrossChainMessages = crossChainMessages

	executor.logger.Trace(fmt.Sprintf("Added %d cross chain messages to batch.",
		len(batch.Header.CrossChainMessages)), log.CmpKey, log.CrossChainCmp)

	batch.Header.LatestInboundCrossChainHash = block.Hash()
	batch.Header.LatestInboundCrossChainHeight = block.Number()

	return nil
}

func (executor *batchExecutor) populateHeader(batch *core.Batch, receipts types.Receipts) {
	if len(receipts) == 0 {
		batch.Header.ReceiptHash = types.EmptyRootHash
	} else {
		batch.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
	}

	if len(batch.Transactions) == 0 {
		batch.Header.TxHash = types.EmptyRootHash
	} else {
		batch.Header.TxHash = types.DeriveSha(types.Transactions(batch.Transactions), trie.NewStackTrie(nil))
	}
}

func (executor *batchExecutor) verifyInboundCrossChainTransactions(transactions types.Transactions, executedTxs types.Transactions, receipts types.Receipts) error {
	if transactions.Len() != executedTxs.Len() {
		return fmt.Errorf("some synthetic transactions have not been executed")
	}

	for _, rec := range receipts {
		if rec.Status == 1 {
			continue
		}
		return fmt.Errorf("found a failed receipt for a synthetic transaction: %s", rec.TxHash.Hex())
	}
	return nil
}

func (executor *batchExecutor) processTransactions(batch *core.Batch, tCount int, txs []*common.L2Tx, stateDB *state.StateDB, cc *params.ChainConfig) ([]*common.L2Tx, []*types.Receipt, error) {
	var executedTransactions []*common.L2Tx
	var txReceipts []*types.Receipt

	txResults := evm.ExecuteTransactions(txs, stateDB, batch.Header, executor.storage, cc, tCount, executor.logger)
	for _, tx := range txs {
		result, f := txResults[tx.Hash()]
		if !f {
			return nil, nil, fmt.Errorf("there should be an entry for each transaction")
		}
		rec, foundReceipt := result.(*types.Receipt)
		if foundReceipt {
			executedTransactions = append(executedTransactions, tx)
			txReceipts = append(txReceipts, rec)
		} else {
			// Exclude all errors
			executor.logger.Info("Excluding transaction from batch", log.TxKey, tx.Hash(), log.BatchHashKey, batch.Hash(), "cause", result)
		}
	}
	sort.Sort(sortByTxIndex(txReceipts))

	return executedTransactions, txReceipts, nil
}

func allReceipts(txReceipts []*types.Receipt, depositReceipts []*types.Receipt) types.Receipts {
	return append(txReceipts, depositReceipts...)
}

type sortByTxIndex []*types.Receipt

func (c sortByTxIndex) Len() int           { return len(c) }
func (c sortByTxIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByTxIndex) Less(i, j int) bool { return c[i].TransactionIndex < c[j].TransactionIndex }
