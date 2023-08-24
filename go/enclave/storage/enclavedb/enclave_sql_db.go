package enclavedb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// enclaveDB - Implements the key-value ethdb.Database and also exposes the underlying sql database
// should not be used directly outside the db package
type enclaveDB struct {
	sqldb  *sql.DB
	logger gethlog.Logger
}

func (sqlDB *enclaveDB) Tail() (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) TruncateHead(uint64) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) TruncateTail(uint64) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) MigrateTable(string, func([]byte) ([]byte, error)) error {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) NewBatchWithSize(int) ethdb.Batch {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) AncientDatadir() (string, error) {
	// TODO implement me
	panic("implement me")
}

func (sqlDB *enclaveDB) NewSnapshot() (ethdb.Snapshot, error) {
	// TODO implement me
	panic("implement me")
}

func NewEnclaveDB(db *sql.DB, logger gethlog.Logger) (EnclaveDB, error) {
	return &enclaveDB{sqldb: db, logger: logger}, nil
}

func (sqlDB *enclaveDB) GetSQLDB() *sql.DB {
	return sqlDB.sqldb
}

func (sqlDB *enclaveDB) BeginTx() (*sql.Tx, error) {
	return sqlDB.sqldb.Begin()
}

func (sqlDB *enclaveDB) Has(key []byte) (bool, error) {
	return Has(sqlDB.sqldb, key)
}

func (sqlDB *enclaveDB) Get(key []byte) ([]byte, error) {
	return Get(sqlDB.sqldb, key)
}

func (sqlDB *enclaveDB) Put(key []byte, value []byte) error {
	return Put(sqlDB.sqldb, key, value)
}

func (sqlDB *enclaveDB) Delete(key []byte) error {
	return Delete(sqlDB.sqldb, key)
}

func (sqlDB *enclaveDB) Close() error {
	if err := sqlDB.sqldb.Close(); err != nil {
		return fmt.Errorf("failed to close sql db - %w", err)
	}
	return nil
}

func (sqlDB *enclaveDB) NewDBTransaction() *dbTransaction {
	return &dbTransaction{
		db: sqlDB,
	}
}

func (sqlDB *enclaveDB) NewBatch() ethdb.Batch {
	return &dbTransaction{
		db: sqlDB,
	}
}

func (sqlDB *enclaveDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return NewIterator(sqlDB.sqldb, prefix, start)
}

func (sqlDB *enclaveDB) Stat(_ string) (string, error) {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return "", nil
}

func (sqlDB *enclaveDB) Compact(_ []byte, _ []byte) error {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return nil
}

// no-freeze! Copied from the geth in-memory sqldb implementation these ancient method implementations disable the 'freezer'

// errNotSupported is returned if the database doesn't support the required operation.
var errNotSupported = errors.New("this operation is not supported")

// HasAncient returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) HasAncient(_ string, _ uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Ancient(_ string, _ uint64) ([]byte, error) {
	return nil, errNotSupported
}

// AncientRange returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) AncientRange(_ string, _, _, _ uint64) ([][]byte, error) {
	return nil, errNotSupported
}

// Ancients returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) AncientSize(_ string) (uint64, error) {
	return 0, errNotSupported
}

// ModifyAncients is not supported.
func (sqlDB *enclaveDB) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) TruncateAncients(_ uint64) error {
	return errNotSupported
}

// Sync returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Sync() error {
	return errNotSupported
}

func (sqlDB *enclaveDB) ReadAncients(fn func(reader ethdb.AncientReaderOp) error) (err error) {
	// Unlike other ancient-related methods, this method does not return
	// errNotSupported when invoked.
	// The reason for this is that the caller might want to do several things:
	// 1. Check if something is in freezer,
	// 2. If not, check leveldb.
	//
	// This will work, since the ancient-checks inside 'fn' will return errors,
	// and the leveldb work will continue.
	//
	// If we instead were to return errNotSupported here, then the caller would
	// have to explicitly check for that, having an extra clause to do the
	// non-ancient operations.
	return fn(sqlDB)
}
