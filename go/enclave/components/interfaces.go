package components

import (
	"errors"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
)

var ErrDuplicateRollup = errors.New("duplicate rollup received")

type BlockIngestionType struct {
	// PreGenesis is true if there is no stored L1 head block.
	// (L1 head is only stored when there is an L2 state to associate it with. Soon we will start consuming from the
	// genesis block and then, we should only see one block ingested in a 'PreGenesis' state)
	PreGenesis bool

	// ChainFork contains information about the status of the new block in the chain
	ChainFork *common.ChainFork
}

func (bit *BlockIngestionType) IsFork() bool {
	if bit.ChainFork == nil {
		return false
	}
	return bit.ChainFork.IsFork()
}

type L1BlockProcessor interface {
	Process(br *common.BlockAndReceipts) (*BlockIngestionType, error)
	GetHead() (*common.L1Block, error)
	GetCrossChainContractAddress() *gethcommon.Address
}

// BatchExecutionContext - Contains all of the data that each batch depends on
type BatchExecutionContext struct {
	BlockPtr     common.L1BlockHash // Block is needed for the cross chain messages
	ParentPtr    common.L2BatchHash
	Transactions common.L2Transactions
	AtTime       uint64
	Creator      gethcommon.Address
	ChainConfig  *params.ChainConfig
	SequencerNo  *big.Int
}

// ComputedBatch - a structure representing the result of a batch
// computation where `Batch` is the newly computed batch and `Receipts`
// are the receipts for the executed transactions inside this batch.
// The `Commit` function allows for committing the stateDB resulting from
// the computation of the batch. One might not want to commit in case the
// resulting batch differs than what is being validated for example.
type ComputedBatch struct {
	Batch    *core.Batch
	Receipts types.Receipts
	Commit   func(bool) (gethcommon.Hash, error)
}

type BatchExecutor interface {
	// ComputeBatch - a more primitive ExecuteBatch
	// Call with same BatchContext should always produce identical extBatch - idempotent
	// Should be safe to call in parallel
	ComputeBatch(*BatchExecutionContext) (*ComputedBatch, error)

	// ExecuteBatch - executes the transactions and xchain messages, returns the receipts, and updates the stateDB
	ExecuteBatch(*core.Batch) (types.Receipts, error)

	// CreateGenesisState - will create and commit the genesis state in the stateDB for the given block hash,
	// and uint64 timestamp representing the time now. In this genesis state is where one can
	// find preallocated funds for faucet. TODO - make this an option
	CreateGenesisState(common.L1BlockHash, uint64) (*core.Batch, *types.Transaction, error)
}

type BatchRegistry interface {
	// BatchesAfter - Given a hash, will return batches following it until the head batch
	BatchesAfter(batchSeqNo uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, error)

	// GetBatchStateAtHeight - creates a stateDB that represents the state committed when
	// the batch with height matching the blockNumber was created and stored.
	GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error)

	// GetBatchAtHeight - same as `GetBatchStateAtHeight`, but instead returns the full batch
	// rather than its stateDB only.
	GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error)

	// SubscribeForBatches - register a callback for new batches
	SubscribeForBatches(func(*core.Batch))
	UnsubscribeFromBatches()

	NotifySubscribers(batch *core.Batch)

	// HasGenesisBatch - returns if genesis batch is available yet or not, or error in case
	// the function is unable to determine.
	HasGenesisBatch() (bool, error)
}

type RollupProducer interface {
	// CreateRollup - creates a rollup starting from the end of the last rollup
	// that has been stored and continues it towards what we consider the current L2 head.
	CreateRollup(fromBatchNo uint64, limiter limiters.RollupLimiter) (*core.Rollup, error)
}

type RollupConsumer interface {
	// ProcessRollupsInBlock - extracts the rollup from the block's transactions
	// and verifies its integrity, saving and processing any batches that have
	// not been seen previously.
	ProcessRollupsInBlock(b *common.BlockAndReceipts) error
}
