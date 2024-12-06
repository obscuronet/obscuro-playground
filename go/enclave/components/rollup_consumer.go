package components

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/ethadapter"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

type rollupConsumerImpl struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	rollupCompression *RollupCompression
	batchRegistry     BatchRegistry

	logger gethlog.Logger

	storage      storage.Storage
	sigValidator *SignatureValidator
}

func NewRollupConsumer(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	batchRegistry BatchRegistry,
	rollupCompression *RollupCompression,
	storage storage.Storage,
	logger gethlog.Logger,
	verifier *SignatureValidator,
) RollupConsumer {
	return &rollupConsumerImpl{
		MgmtContractLib:   mgmtContractLib,
		batchRegistry:     batchRegistry,
		rollupCompression: rollupCompression,
		logger:            logger,
		storage:           storage,
		sigValidator:      verifier,
	}
}

// ProcessBlobsInBlock - processes the blobs in a block, extracts the rollups, verifies the rollups and stores them
func (rc *rollupConsumerImpl) ProcessBlobsInBlock(ctx context.Context, processed *common.ProcessedL1Data) error {
	defer core.LogMethodDuration(rc.logger, measure.NewStopwatch(), "Rollup consumer processed blobs", log.BlockHashKey, processed.BlockHeader.Hash())

	block := processed.BlockHeader
	rollups, err := rc.extractAndVerifyRollups(processed)
	if err != nil {
		rc.logger.Error("Failed to extract rollups from block", log.BlockHashKey, block.Hash(), log.ErrKey, err)
		return err
	}
	if len(rollups) == 0 {
		rc.logger.Trace("No rollups found in block", log.BlockHashKey, block.Hash())
		return nil
	}

	rollups, err = rc.getSignedRollup(rollups)
	if err != nil {
		return err
	}

	if len(rollups) > 1 {
		// todo - we need to sort this out
		rc.logger.Warn(fmt.Sprintf("Multiple rollups %d in block %s", len(rollups), block.Hash()))
	}

	for _, rollup := range rollups {
		l1CompressionBlock, err := rc.storage.FetchBlock(ctx, rollup.Header.CompressionL1Head)
		if err != nil {
			rc.logger.Warn("Can't process rollup because the l1 block used for compression is not available", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		canonicalBlockByHeight, err := rc.storage.FetchCanonicaBlockByHeight(ctx, l1CompressionBlock.Number)
		if err != nil {
			return err
		}
		if canonicalBlockByHeight.Hash() != l1CompressionBlock.Hash() {
			rc.logger.Warn("Skipping rollup because it was compressed on top of a non-canonical rollup", "block_hash", rollup.Header.CompressionL1Head, log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			continue
		}
		// read batch data from rollup, verify and store it
		internalHeader, err := rc.rollupCompression.ProcessExtRollup(ctx, rollup)
		if err != nil {
			rc.logger.Error("Failed processing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			// todo - issue challenge as a validator
			return err
		}
		if err := rc.storage.StoreRollup(ctx, rollup, internalHeader); err != nil {
			rc.logger.Error("Failed storing rollup", log.RollupHashKey, rollup.Hash(), log.ErrKey, err)
			return err
		}
	}

	return nil
}

func (rc *rollupConsumerImpl) getSignedRollup(rollups []*common.ExtRollup) ([]*common.ExtRollup, error) {
	signedRollup := make([]*common.ExtRollup, 0)

	// loop through the rollups, find the one that is signed, verify the signature, make sure it's the only one
	for _, rollup := range rollups {
		if err := rc.sigValidator.CheckSequencerSignature(rollup.Hash(), rollup.Header.Signature); err != nil {
			return nil, fmt.Errorf("rollup signature was invalid. Cause: %w", err)
		}

		signedRollup = append(signedRollup, rollup)
	}
	return signedRollup, nil
}

// todo - when processing the rollup, instead of looking up batches one by one, compare the last sequence number from the db with the ones in the rollup
// extractAndVerifyRollups returns a list of the rollups published in this block
// It processes each transaction, attempting to extract and verify rollups
// If a transaction is not a rollup or fails verification, it's skipped
// The function only returns an error if there's a critical failure in rollup reconstruction
func (rc *rollupConsumerImpl) extractAndVerifyRollups(processed *common.ProcessedL1Data) ([]*common.ExtRollup, error) {
	rollupTxs := processed.GetEvents(common.RollupTx)
	rollups := make([]*common.ExtRollup, 0, len(rollupTxs))

	blobs, blobHashes, err := rc.extractBlobsAndHashes(rollupTxs)
	if err != nil {
		return nil, err
	}

	for i, tx := range rollupTxs {
		t := rc.MgmtContractLib.DecodeTx(tx.Transaction)
		if t == nil {
			continue
		}

		rollupHashes, ok := t.(*common.L1RollupHashes)
		if !ok {
			continue
		}

		if err := verifyBlobHashes(rollupHashes, blobHashes); err != nil {
			rc.logger.Warn(fmt.Sprintf("blob hashes in rollup at index %d do not match the rollup blob hashes. Cause: %s", i, err))
			continue // Blob hashes don't match, skip this rollup
		}

		r, err := ethadapter.ReconstructRollup(blobs)
		if err != nil {
			// This is a critical error because we've already verified the blob hashes
			// If we can't reconstruct the rollup at this point, something is seriously wrong
			return nil, fmt.Errorf("could not recreate rollup from blobs. Cause: %w", err)
		}

		rollups = append(rollups, r)

		rc.logger.Info("Extracted rollup from block", log.RollupHashKey, r.Hash(), log.BlockHashKey, processed.BlockHeader.Hash())
	}

	return rollups, nil
}

// there may be many rollups in one block so the blobHashes array, so it is possible that the rollupHashes array is a
// subset of the blobHashes array
func verifyBlobHashes(rollupHashes *common.L1RollupHashes, blobHashes []gethcommon.Hash) error {
	// more efficient lookup
	blobHashSet := make(map[gethcommon.Hash]struct{}, len(blobHashes))
	for _, h := range blobHashes {
		blobHashSet[h] = struct{}{}
	}

	for i, rollupHash := range rollupHashes.BlobHashes {
		if _, exists := blobHashSet[rollupHash]; !exists {
			return fmt.Errorf(
				"rollupHash at index %d (%s) not found in blobHashes",
				i,
				rollupHash.Hex(),
			)
		}
	}
	return nil
}

func (rc *rollupConsumerImpl) extractBlobsAndHashes(rollupTxs []*common.L1TxData) ([]*kzg4844.Blob, []gethcommon.Hash, error) {
	blobs := make([]*kzg4844.Blob, 0)
	for _, tx := range rollupTxs {
		blobs = append(blobs, tx.Blobs...)
	}

	_, blobHashes, err := ethadapter.MakeSidecar(blobs)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create blob sidecar and blob hashes. Cause: %w", err)
	}

	return blobs, blobHashes, nil
}
