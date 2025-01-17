package nodetype

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

type validator struct {
	blockProcessor components.L1BlockProcessor
	batchExecutor  components.BatchExecutor
	batchRegistry  components.BatchRegistry

	chainConfig *params.ChainConfig

	storage      storage.Storage
	sigValidator components.SequencerSignatureVerifier
	mempool      *components.TxPool

	logger gethlog.Logger
}

func NewValidator(
	consumer components.L1BlockProcessor,
	batchExecutor components.BatchExecutor,
	registry components.BatchRegistry,
	chainConfig *params.ChainConfig,
	storage storage.Storage,
	sigValidator components.SequencerSignatureVerifier,
	mempool *components.TxPool,
	logger gethlog.Logger,
) Validator {
	return &validator{
		blockProcessor: consumer,
		batchExecutor:  batchExecutor,
		batchRegistry:  registry,
		chainConfig:    chainConfig,
		storage:        storage,
		sigValidator:   sigValidator,
		mempool:        mempool,
		logger:         logger,
	}
}

func (val *validator) OnL1Fork(ctx context.Context, fork *common.ChainFork) error {
	// nothing to do
	return nil
}

func (val *validator) VerifySequencerSignature(b *core.Batch) error {
	return val.sigValidator.CheckSequencerSignature(b.Hash(), b.Header.Signature)
}

func (val *validator) ExecuteStoredBatches(ctx context.Context) error {
	val.logger.Trace("Executing stored batches")
	headBatchSeq := val.batchRegistry.HeadBatchSeq()
	if headBatchSeq == nil {
		headBatchSeq = big.NewInt(int64(common.L2GenesisSeqNo))
	}
	batches, err := val.storage.FetchCanonicalUnexecutedBatches(ctx, headBatchSeq)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		return err
	}

	for _, batchHeader := range batches {
		if batchHeader.IsGenesis() {
			if err = val.handleGenesis(ctx, batchHeader); err != nil {
				return err
			}
		}

		val.logger.Trace("Executing stored batch", log.BatchSeqNoKey, batchHeader.SequencerOrderNo)

		// check batchHeader execution prerequisites
		canExecute, err := val.executionPrerequisites(ctx, batchHeader)
		if err != nil {
			return fmt.Errorf("could not determine the execution prerequisites for batchHeader %s. Cause: %w", batchHeader.Hash(), err)
		}
		val.logger.Trace("Can execute stored batch", log.BatchSeqNoKey, batchHeader.SequencerOrderNo, "can", canExecute)

		if canExecute {
			txs, err := val.storage.FetchBatchTransactionsBySeq(ctx, batchHeader.SequencerOrderNo.Uint64())
			if err != nil {
				return fmt.Errorf("could not get txs for batch %s. Cause: %w", batchHeader.Hash(), err)
			}

			batch := &core.Batch{
				Header:       batchHeader,
				Transactions: txs,
			}

			txResults, err := val.batchExecutor.ExecuteBatch(ctx, batch)
			if err != nil {
				return fmt.Errorf("could not execute batch %s. Cause: %w", batchHeader.Hash(), err)
			}
			err = val.storage.StoreExecutedBatch(ctx, batch, txResults)
			if err != nil {
				return fmt.Errorf("could not store executed batch %s. Cause: %w", batchHeader.Hash(), err)
			}
			err = val.batchRegistry.OnBatchExecuted(batchHeader, txResults)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (val *validator) executionPrerequisites(ctx context.Context, batch *common.BatchHeader) (bool, error) {
	// 1.l1 block exists
	block, err := val.storage.FetchBlock(ctx, batch.L1Proof)
	if err != nil && errors.Is(err, errutil.ErrNotFound) {
		val.logger.Warn("Error fetching block", log.BlockHashKey, batch.L1Proof, log.ErrKey, err)
		return false, err
	}
	val.logger.Trace("l1 block exists", log.BatchSeqNoKey, batch.SequencerOrderNo)
	// 2. parent was executed
	parentExecuted, err := val.storage.BatchWasExecuted(ctx, batch.ParentHash)
	if err != nil {
		val.logger.Info("Error reading execution status of batch", log.BatchHashKey, batch.ParentHash, log.ErrKey, err)
		return false, err
	}
	val.logger.Trace("parentExecuted", log.BatchSeqNoKey, batch.SequencerOrderNo, "val", parentExecuted)

	return block != nil && parentExecuted, nil
}

func (val *validator) handleGenesis(ctx context.Context, batch *common.BatchHeader) error {
	genBatch, _, err := val.batchExecutor.CreateGenesisState(ctx, batch.L1Proof, batch.Time, batch.Coinbase, batch.BaseFee)
	if err != nil {
		return err
	}

	if genBatch.Hash() != batch.Hash() {
		return fmt.Errorf("received invalid genesis batch")
	}

	err = val.storage.StoreExecutedBatch(ctx, genBatch, nil)
	if err != nil {
		return err
	}
	err = val.batchRegistry.OnBatchExecuted(batch, nil)
	if err != nil {
		return err
	}
	return nil
}

func (val *validator) OnL1Block(ctx context.Context, block *types.Header, result *components.BlockIngestionType) error {
	return val.ExecuteStoredBatches(ctx)
}

func (val *validator) Close() error {
	return val.mempool.Close()
}
