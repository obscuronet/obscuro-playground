package components

import (
	"errors"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/crosschain"
)

type l1BlockProcessor struct {
	storage              storage.Storage
	logger               gethlog.Logger
	crossChainProcessors *crosschain.Processors
}

func NewBlockProcessor(storage storage.Storage, cc *crosschain.Processors, logger gethlog.Logger) L1BlockProcessor {
	return &l1BlockProcessor{
		storage:              storage,
		logger:               logger,
		crossChainProcessors: cc,
	}
}

func (bp *l1BlockProcessor) Process(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, []common.L1BlockHash, error) {
	defer bp.logger.Info("L1 block processed", log.BlockHashKey, br.Block.Hash(), log.DurationKey, measure.NewStopwatch())

	ingestion, nonCanonicalPath, err := bp.tryAndInsertBlock(br, isLatest)
	if err != nil {
		return nil, nil, err
	}

	if !ingestion.PreGenesis {
		// This requires block to be stored first ... but can permanently fail a block
		err = bp.crossChainProcessors.Remote.StoreCrossChainMessages(br.Block, *br.Receipts)
		if err != nil {
			return nil, nil, errors.New("failed to process cross chain messages")
		}
	}

	return ingestion, nonCanonicalPath, nil
}

func (bp *l1BlockProcessor) tryAndInsertBlock(br *common.BlockAndReceipts, isLatest bool) (*BlockIngestionType, []common.L1BlockHash, error) {
	block := br.Block

	_, err := bp.storage.FetchBlock(block.Hash())
	if err == nil {
		return nil, nil, errutil.ErrBlockAlreadyProcessed
	}

	if !errors.Is(err, errutil.ErrNotFound) {
		return nil, nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	// We insert the block into the L1 chain and store it.
	ingestionType, canonical, nonCanonical, err := bp.ingestBlock(block, isLatest)
	if err != nil {
		// Do not store the block if the L1 chain insertion failed
		return nil, nil, err
	}
	bp.logger.Trace("block inserted successfully",
		log.BlockHeightKey, block.NumberU64(), log.BlockHashKey, block.Hash(), "ingestionType", ingestionType)

	err = bp.storage.StoreBlock(block, canonical, nonCanonical)
	if err != nil {
		return nil, nil, fmt.Errorf("could not store block. Cause: %w", err)
	}

	return ingestionType, nonCanonical, nil
}

func (bp *l1BlockProcessor) ingestBlock(block *common.L1Block, isLatest bool) (*BlockIngestionType, []common.L1BlockHash, []common.L1BlockHash, error) {
	// todo (#1056) - this is minimal L1 tracking/validation, and should be removed when we are using geth's blockchain or lightchain structures for validation
	prevL1Head, err := bp.storage.FetchHeadBlock()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// todo (@matt) - we should enforce that this block is a configured hash (e.g. the L1 management contract deployment block)
			return &BlockIngestionType{IsLatest: isLatest, Fork: false, PreGenesis: true}, nil, nil, nil
		}
		return nil, nil, nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	isFork := false
	// we do a basic sanity check, comparing the received block to the head block on the chain
	if block.ParentHash() != prevL1Head.Hash() {
		lcaBlock, newCanonicalChain, newNonCanonicalChain, err := gethutil.LCA(block, prevL1Head, bp.storage)
		if err != nil {
			bp.logger.Trace("parent not found",
				"blkHeight", block.NumberU64(), log.BlockHashKey, block.Hash(),
				"l1HeadHeight", prevL1Head.NumberU64(), "l1HeadHash", prevL1Head.Hash(),
			)
			return nil, nil, nil, errutil.ErrBlockAncestorNotFound
		}

		// fork - least common ancestor for this block and l1 head is before the l1 head.
		isFork = lcaBlock.NumberU64() < prevL1Head.NumberU64()
		if isFork {
			bp.logger.Info("Fork detected in the l1 chain", "can", lcaBlock.Hash().Hex(), "noncan", prevL1Head.Hash().Hex())
		}
		return &BlockIngestionType{IsLatest: isLatest, Fork: isFork, PreGenesis: false}, newCanonicalChain, newNonCanonicalChain, nil
	}
	return &BlockIngestionType{IsLatest: isLatest, Fork: isFork, PreGenesis: false}, nil, nil, nil
}

func (bp *l1BlockProcessor) GetHead() (*common.L1Block, error) {
	return bp.storage.FetchHeadBlock()
}

func (bp *l1BlockProcessor) GetCrossChainContractAddress() *gethcommon.Address {
	return bp.crossChainProcessors.Remote.GetBusAddress()
}
