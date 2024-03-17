package components

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"

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

func (re *rollupProducerImpl) CreateInternalRollup(fromBatchNo uint64, upToL1Height uint64, limiter limiters.RollupLimiter) (*core.Rollup, error) {
	batches, blocks, err := re.batchRegistry.BatchesAfter(fromBatchNo, upToL1Height, limiter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch 'from' batch (seqNo=%d) for rollup: %w", fromBatchNo, err)
	}

	hasBatches := len(batches) != 0

	if !hasBatches {
		return nil, fmt.Errorf("no batches for rollup")
	}

	block, err := re.storage.FetchCanonicaBlockByHeight(big.NewInt(int64(upToL1Height)))
	if err != nil {
		return nil, err
	}

	rh := common.RollupHeader{}
	rh.CompressionL1Head = block.Hash()

	rh.CrossChainMessages = make([]MessageBus.StructsCrossChainMessage, 0)
	for _, b := range batches {
		rh.CrossChainMessages = append(rh.CrossChainMessages, b.Header.CrossChainMessages...)
	}

	lastBatch := batches[len(batches)-1]
	rh.LastBatchSeqNo = lastBatch.SeqNo().Uint64()

	blockMap := map[common.L1BlockHash]*types.Block{}
	for _, b := range blocks {
		blockMap[b.Hash()] = b
	}

	newRollup := &core.Rollup{
		Header:  &rh,
		Blocks:  blockMap,
		Batches: batches,
	}

	re.logger.Info(fmt.Sprintf("Created new rollup %s with %d batches. From %d to %d", newRollup.Hash(), len(newRollup.Batches), batches[0].SeqNo(), rh.LastBatchSeqNo))

	return newRollup, nil
}
