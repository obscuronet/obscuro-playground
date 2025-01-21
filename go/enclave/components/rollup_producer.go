package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/merkle"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

type rollupProducerImpl struct {
	enclaveID     gethcommon.Address
	storage       storage.Storage
	batchRegistry BatchRegistry
	logger        gethlog.Logger
}

func NewRollupProducer(enclaveID gethcommon.Address, storage storage.Storage, batchRegistry BatchRegistry, logger gethlog.Logger) RollupProducer {
	return &rollupProducerImpl{
		enclaveID:     enclaveID,
		logger:        logger,
		batchRegistry: batchRegistry,
		storage:       storage,
	}
}

func (re *rollupProducerImpl) CreateInternalRollup(ctx context.Context, fromBatchNo uint64, upToL1Height uint64, limiter limiters.RollupLimiter) (*core.Rollup, error) {
	batches, blocks, err := re.batchRegistry.BatchesAfter(ctx, fromBatchNo, upToL1Height, limiter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch 'from' batch (seqNo=%d) for rollup: %w", fromBatchNo, err)
	}

	hasBatches := len(batches) != 0

	if !hasBatches {
		return nil, fmt.Errorf("no batches for rollup")
	}

	block, err := re.storage.FetchHeadBlock(ctx)
	if err != nil {
		return nil, err
	}

	rh := common.RollupHeader{}
	rh.CompressionL1Head = block.Hash()
	rh.CompressionL1Number = block.Number

	lastBatch := batches[len(batches)-1]
	rh.LastBatchSeqNo = lastBatch.SeqNo().Uint64()

	blockMap := map[common.L1BlockHash]*types.Header{}
	for _, b := range blocks {
		blockMap[b.Hash()] = b
	}

	exportedCrossChainRoot, err := exportCrossChainData(ctx, re.storage, batches[0].SeqNo().Uint64(), rh.LastBatchSeqNo)
	if err != nil {
		return nil, err
	}

	rh.CrossChainRoot = *exportedCrossChainRoot

	newRollup := &core.Rollup{
		Header:  &rh,
		Blocks:  blockMap,
		Batches: batches,
	}

	re.logger.Info(fmt.Sprintf("Created new rollup %s with %d batches. From %d to %d", newRollup.Hash(), len(newRollup.Batches), batches[0].SeqNo(), rh.LastBatchSeqNo))

	return newRollup, nil
}

func exportCrossChainData(ctx context.Context, storage storage.Storage, fromSeqNo uint64, toSeqNo uint64) (*gethcommon.Hash, error) {
	canonicalBatches, err := storage.FetchCanonicalBatchesBetween((ctx), fromSeqNo, toSeqNo)
	if err != nil {
		return nil, err
	}

	root, _, err := merkle.ComputeCrossChainRootFromBatches(canonicalBatches)
	if err != nil {
		return nil, err
	}

	return &root, nil
}
