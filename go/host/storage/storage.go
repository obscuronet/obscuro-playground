package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
)

type storageImpl struct {
	db     hostdb.HostDB
	logger gethlog.Logger
	io.Closer
}

func (s *storageImpl) AddBatch(batch *common.ExtBatch) error {
	// Check if the Batch is already stored
	_, err := hostdb.GetBatchHeader(s.db, batch.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}

	if err := hostdb.AddBatch(dbtx, s.db.GetSQLStatement(), batch); err != nil {
		if err1 := dbtx.Rollback(); err1 != nil {
			return err1
		}
		if errors.Is(err, errutil.ErrAlreadyExists) {
			return err
		}
		return fmt.Errorf("could not add batch to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit batch tx. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *types.Header) error {
	// Check if the Header is already stored
	_, err := hostdb.GetRollupHeader(s.db, rollup.Header.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}

	if err := hostdb.AddRollup(dbtx, s.db.GetSQLStatement(), rollup, metadata, block); err != nil {
		if err := dbtx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("could not add rollup to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit rollup tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) AddBlock(b *types.Header) error {
	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbtx.Rollback()

	_, err = hostdb.GetBlockId(dbtx.Tx, s.db.GetSQLStatement(), b.Hash())
	switch {
	case err == nil:
		// Block already exists
		s.logger.Debug("Block already exists", "hash", b.Hash().Hex())
		return nil
	case !errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("error checking block existence: %w", err)
	}

	if err := hostdb.AddBlock(dbtx.Tx, s.db.GetSQLStatement(), b); err != nil {
		if IsConstraintError(err) {
			s.logger.Debug("Block already exists",
				"hash", b.Hash().Hex(),
				"error", err)
			return nil
		}
		return fmt.Errorf("could not add block to host: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit block tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error) {
	return hostdb.GetBatchBySequenceNumber(s.db, seqNum)
}

func (s *storageImpl) FetchBatchHashByHeight(number *big.Int) (*gethcommon.Hash, error) {
	return hostdb.GetBatchHashByNumber(s.db, number)
}

func (s *storageImpl) FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	return hostdb.GetBatchHeader(s.db, hash)
}

func (s *storageImpl) FetchHeadBatchHeader() (*common.BatchHeader, error) {
	return hostdb.GetHeadBatchHeader(s.db)
}

func (s *storageImpl) FetchPublicBatchByHash(batchHash common.L2BatchHash) (*common.PublicBatch, error) {
	return hostdb.GetPublicBatch(s.db, batchHash)
}

func (s *storageImpl) FetchBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	return hostdb.GetBatchByHash(s.db, batchHash)
}

func (s *storageImpl) FetchBatchByTx(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	return hostdb.GetBatchByTx(s.db, txHash)
}

func (s *storageImpl) FetchLatestBatch() (*common.BatchHeader, error) {
	return hostdb.GetLatestBatch(s.db)
}

func (s *storageImpl) FetchBatchHeaderByHeight(height *big.Int) (*common.BatchHeader, error) {
	return hostdb.GetBatchHeaderByHeight(s.db, height)
}

func (s *storageImpl) FetchBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	return hostdb.GetBatchByHeight(s.db, height)
}

func (s *storageImpl) FetchBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return hostdb.GetBatchListing(s.db, pagination)
}

func (s *storageImpl) FetchBatchListingDeprecated(pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	return hostdb.GetBatchListingDeprecated(s.db, pagination)
}

func (s *storageImpl) FetchLatestRollupHeader() (*common.RollupHeader, error) {
	return hostdb.GetLatestRollup(s.db)
}

func (s *storageImpl) FetchRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	return hostdb.GetRollupListing(s.db, pagination)
}

func (s *storageImpl) FetchBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	return hostdb.GetBlockListing(s.db, pagination)
}

func (s *storageImpl) FetchTotalTxCount() (*big.Int, error) {
	return hostdb.GetTotalTxCount(s.db)
}

func (s *storageImpl) FetchTotalTxsQuery() (*big.Int, error) {
	return hostdb.GetTotalTxsQuery(s.db)
}

func (s *storageImpl) FetchTransaction(hash gethcommon.Hash) (*common.PublicTransaction, error) {
	return hostdb.GetTransaction(s.db, hash)
}

func (s *storageImpl) FetchRollupByHash(rollupHash gethcommon.Hash) (*common.PublicRollup, error) {
	return hostdb.GetRollupByHash(s.db, rollupHash)
}

func (s *storageImpl) FetchRollupBySeqNo(seqNo uint64) (*common.PublicRollup, error) {
	return hostdb.GetRollupBySeqNo(s.db, seqNo)
}

func (s *storageImpl) FetchRollupBatches(rollupHash gethcommon.Hash) (*common.BatchListingResponse, error) {
	return hostdb.GetRollupBatches(s.db, rollupHash)
}

func (s *storageImpl) FetchBatchTransactions(batchHash gethcommon.Hash) (*common.TransactionListingResponse, error) {
	return hostdb.GetBatchTransactions(s.db, batchHash)
}

func (s *storageImpl) FetchTransactionListing(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return hostdb.GetTransactionListing(s.db, pagination)
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func NewHostStorageFromConfig(config *hostconfig.HostConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, logger)
}

func NewStorage(backingDB hostdb.HostDB, logger gethlog.Logger) Storage {
	return &storageImpl{
		db:     backingDB,
		logger: logger,
	}
}

// SQLite constraint error messages
const (
	ErrUniqueBlockHash = "UNIQUE constraint failed: block_host.hash"
	ErrForeignKey      = "FOREIGN KEY constraint failed"
)

// IsConstraintError returns true if the error is a known constraint error
func IsConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, ErrUniqueBlockHash) ||
		strings.Contains(errMsg, ErrForeignKey)
}
