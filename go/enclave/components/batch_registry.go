package components

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/state"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common/async"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"
)

type batchRegistry struct {
	storage      storage.Storage
	logger       gethlog.Logger
	headBatchSeq *big.Int // keep track of the last executed batch to optimise db access

	batchesCallback   func(*core.Batch, types.Receipts)
	callbackMutex     sync.RWMutex
	healthTimeout     time.Duration
	lastExecutedBatch *async.Timestamp
}

func NewBatchRegistry(storage storage.Storage, logger gethlog.Logger) BatchRegistry {
	var headBatchSeq *big.Int
	headBatch, err := storage.FetchHeadBatch()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			headBatchSeq = nil
		} else {
			logger.Crit("Could not create batch registry", log.ErrKey, err)
			return nil
		}
	} else {
		headBatchSeq = headBatch.SeqNo()
	}
	return &batchRegistry{
		storage:           storage,
		headBatchSeq:      headBatchSeq,
		logger:            logger,
		healthTimeout:     time.Minute,
		lastExecutedBatch: async.NewAsyncTimestamp(time.Now().Add(-time.Minute)),
	}
}

func (br *batchRegistry) HeadBatchSeq() *big.Int {
	return br.headBatchSeq
}

func (br *batchRegistry) SubscribeForExecutedBatches(callback func(*core.Batch, types.Receipts)) {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()
	br.batchesCallback = callback
}

func (br *batchRegistry) UnsubscribeFromBatches() {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()

	br.batchesCallback = nil
}

func (br *batchRegistry) OnBatchExecuted(batch *core.Batch, receipts types.Receipts) {
	br.callbackMutex.RLock()
	defer br.callbackMutex.RUnlock()

	defer core.LogMethodDuration(br.logger, measure.NewStopwatch(), "Sending batch and events", log.BatchHashKey, batch.Hash())

	br.headBatchSeq = batch.SeqNo()
	if br.batchesCallback != nil {
		br.batchesCallback(batch, receipts)
	}

	br.lastExecutedBatch.Mark()
}

func (br *batchRegistry) HasGenesisBatch() (bool, error) {
	return br.headBatchSeq != nil, nil
}

func (br *batchRegistry) BatchesAfter(batchSeqNo uint64, upToL1Height uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, []*types.Block, error) {
	// sanity check
	headBatch, err := br.storage.FetchBatchBySeqNo(br.headBatchSeq.Uint64())
	if err != nil {
		return nil, nil, err
	}

	if headBatch.SeqNo().Uint64() < batchSeqNo {
		return nil, nil, fmt.Errorf("head batch height %d is in the past compared to requested batch %d", headBatch.SeqNo().Uint64(), batchSeqNo)
	}

	resultBatches := make([]*core.Batch, 0)
	resultBlocks := make([]*types.Block, 0)

	currentBatchSeq := batchSeqNo
	var currentBlock *types.Block
	for currentBatchSeq <= headBatch.SeqNo().Uint64() {
		batch, err := br.storage.FetchBatchBySeqNo(currentBatchSeq)
		if err != nil {
			return nil, nil, fmt.Errorf("could not retrieve batch by sequence number %d. Cause: %w", currentBatchSeq, err)
		}

		// check the block height
		// if it's the same block as the previous batch there is no reason to check
		if currentBlock == nil || currentBlock.Hash() != batch.Header.L1Proof {
			block, err := br.storage.FetchBlock(batch.Header.L1Proof)
			if err != nil {
				return nil, nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
			}
			currentBlock = block
			if block.NumberU64() > upToL1Height {
				break
			}
			resultBlocks = append(resultBlocks, block)
		}

		// check the limiter
		didAcceptBatch, err := rollupLimiter.AcceptBatch(batch)
		if err != nil {
			return nil, nil, err
		}
		if !didAcceptBatch {
			break
		}

		resultBatches = append(resultBatches, batch)
		br.logger.Info("Added batch to rollup", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SeqNo(), log.BatchHeightKey, batch.Number(), "l1_proof", batch.Header.L1Proof)

		currentBatchSeq++
	}

	if len(resultBatches) > 0 {
		// Sanity check that the rollup includes consecutive batches (according to the seqNo)
		current := resultBatches[0].SeqNo().Uint64()
		for i, b := range resultBatches {
			if current+uint64(i) != b.SeqNo().Uint64() {
				return nil, nil, fmt.Errorf("created invalid rollup with batches out of sequence")
			}
		}
	}

	return resultBatches, resultBlocks, nil
}

func (br *batchRegistry) GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the batch of interest.
	batch, err := br.GetBatchAtHeight(*blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := br.storage.CreateStateDB(batch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for batch %s", batch.Hash().Hex())
	}

	return blockchainState, err
}

func (br *batchRegistry) GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error) {
	if br.headBatchSeq == nil {
		return nil, fmt.Errorf("chain not initialised")
	}
	var batch *core.Batch
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisBatch, err := br.storage.FetchBatchByHeight(0)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		batch = genesisBatch
	// note: our API currently treats all these block statuses the same for obscuro batches
	case gethrpc.SafeBlockNumber, gethrpc.FinalizedBlockNumber, gethrpc.LatestBlockNumber, gethrpc.PendingBlockNumber:
		headBatch, err := br.storage.FetchBatchBySeqNo(br.headBatchSeq.Uint64())
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d was not found. Cause: %w", height, err)
		}
		batch = headBatch
	default:
		maybeBatch, err := br.storage.FetchBatchByHeight(uint64(height))
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		batch = maybeBatch
	}
	return batch, nil
}

// HealthCheck checks if the last executed batch was more than healthTimeout ago
func (br *batchRegistry) HealthCheck() (bool, error) {
	lastExecutedBatchTime := br.lastExecutedBatch.LastTimestamp()
	if time.Now().After(lastExecutedBatchTime.Add(br.healthTimeout)) {
		return false, fmt.Errorf("last executed batch was %s ago", time.Since(lastExecutedBatchTime))
	}

	return true, nil
}
