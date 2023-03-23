package db

import (
	"bytes"
	"crypto/ecdsa"
	sql2 "database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/eth/filters"

	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"
)

// ErrNoRollups is returned if no rollups have been published yet in the history of the network
// Note: this is not just "not found", we cache at every L1 block what rollup we are up to so we also record that we haven't seen one yet
var ErrNoRollups = errors.New("no rollups have been published")

// TODO - Consistency around whether we assert the secret is available or not.

type storageImpl struct {
	db          *sql.EnclaveDB
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorage(backingDB *sql.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	return &storageImpl{
		db:          backingDB,
		stateDB:     state.NewDatabase(backingDB),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	headHash, err := obscurorawdb.ReadL2HeadBatch(s.db)
	if err != nil {
		return nil, err
	}
	if (bytes.Equal(headHash.Bytes(), gethcommon.Hash{}.Bytes())) {
		return nil, errutil.ErrNotFound
	}
	return s.FetchBatch(*headHash)
}

func (s *storageImpl) FetchBatch(hash common.L2RootHash) (*core.Batch, error) {
	s.assertSecretAvailable()
	batch, err := obscurorawdb.ReadBatch(s.db, hash)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	hash, err := obscurorawdb.ReadCanonicalBatchHash(s.db, height)
	if err != nil {
		return nil, err
	}
	return s.FetchBatch(*hash)
}

func (s *storageImpl) StoreBlock(b *types.Block) {
	s.assertSecretAvailable()
	rawdb.WriteBlock(s.db, b)
}

func (s *storageImpl) FetchBlock(blockHash common.L1RootHash) (*types.Block, error) {
	s.assertSecretAvailable()
	height := rawdb.ReadHeaderNumber(s.db, blockHash)
	if height == nil {
		return nil, errutil.ErrNotFound
	}
	b := rawdb.ReadBlock(s.db, blockHash, *height)
	if b == nil {
		return nil, errutil.ErrNotFound
	}
	return b, nil
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	s.assertSecretAvailable()
	block, err := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	return obscurorawdb.WriteSharedSecret(s.db, secret)
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	return obscurorawdb.ReadSharedSecret(s.db)
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(block.ParentHash())
	if err != nil {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	headBatch, err := s.FetchHeadBatch()
	if err != nil {
		s.logger.Error("unable to HealthCheck storage", log.ErrKey, err)
		return false, err
	}
	return headBatch != nil, nil
}

func (s *storageImpl) assertSecretAvailable() {
	// TODO uncomment this
	//if s.FetchSecret() == nil {
	//	panic("Enclave not initialized")
	//}
}

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1RootHash) (*core.Batch, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadBatchForBlock(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head batch for block. Cause: %w", err)
	}
	return obscurorawdb.ReadBatch(s.db, *l2HeadBatch)
}

func (s *storageImpl) FetchHeadRollupForBlock(blockHash *common.L1RootHash) (*core.Rollup, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadRollup(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head rollup for block. Cause: %w", err)
	}
	if *l2HeadBatch == (gethcommon.Hash{}) { // empty hash ==> no rollups yet up to this block
		return nil, ErrNoRollups
	}
	return obscurorawdb.ReadRollup(s.db, *l2HeadBatch)
}

func (s *storageImpl) FetchLogs(blockHash common.L1RootHash) ([]*types.Log, error) {
	logs, err := obscurorawdb.ReadBlockLogs(s.db, blockHash)
	if err != nil {
		// TODO - Return the error itself, once we move from `errutil.ErrNotFound` to `ethereum.NotFound`
		return nil, errutil.ErrNotFound
	}
	return logs, nil
}

func (s *storageImpl) UpdateHeadBatch(l1Head common.L1RootHash, l2Head *core.Batch, receipts []*types.Receipt) error {
	dbBatch := s.db.NewSqlBatch()

	if err := obscurorawdb.SetL2HeadBatch(dbBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := obscurorawdb.WriteL1ToL2BatchMapping(dbBatch, l1Head, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.WriteCanonicalHash(dbBatch, l2Head); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}

	stateDB, err := s.CreateStateDB(*l2Head.Hash())
	if err != nil {
		return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
	}

	// We update the block's logs, based on the batch's logs.
	var logs []*types.Log
	for _, receipt := range receipts {
		logs = append(logs, receipt.Logs...)
		// todo write the blocks differently
		for _, _log := range receipt.Logs {
			var t0, t1, t2, t3, t4 *gethcommon.Hash
			var addr1, addr2, addr3, addr4 *gethcommon.Address
			isLifecycle := true

			n := len(_log.Topics)
			if n > 0 {
				t0 = &_log.Topics[0]
			}

			// for every indexed topic, check whether it is an end user account
			// if yes, then mark it as relevant for that account
			if n > 1 {
				t1 = &_log.Topics[1]
				if isEndUserAccount(*t1, stateDB) {
					isLifecycle = false
					a := gethcommon.BytesToAddress(t1.Bytes())
					addr1 = &a
				}
			}
			if n > 2 {
				t2 = &_log.Topics[2]
				if isEndUserAccount(*t2, stateDB) {
					isLifecycle = false
					a := gethcommon.BytesToAddress(t2.Bytes())
					addr2 = &a
				}
			}
			if n > 3 {
				t3 = &_log.Topics[3]
				if isEndUserAccount(*t3, stateDB) {
					isLifecycle = false
					a := gethcommon.BytesToAddress(t3.Bytes())
					addr3 = &a
				}
			}
			if n > 4 {
				t4 = &_log.Topics[4]
				if isEndUserAccount(*t4, stateDB) {
					isLifecycle = false
					a := gethcommon.BytesToAddress(t4.Bytes())
					addr4 = &a
				}
			}
			dbBatch.ExecuteSQL("insert into events values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
				t0, t1, t2, t3, t4,
				_log.Data, _log.BlockHash, _log.BlockNumber, _log.TxHash, _log.TxIndex, _log.Index, _log.Address,
				isLifecycle, addr1, addr2, addr3, addr4,
			)
		}
	}
	if err := obscurorawdb.WriteBlockLogs(dbBatch, l1Head, logs); err != nil {
		return fmt.Errorf("could not write block logs. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// Of the log's topics, returns those that are (potentially) user addresses. A topic is considered a user address if:
//   - It has 12 leading zero bytes (since addresses are 20 bytes long, while hashes are 32)
//   - It has a non-zero nonce (to prevent accidental or malicious creation of the address matching a given topic,
//     forcing its events to become permanently private
//   - It does not have associated code (meaning it's a smart-contract address)
func isEndUserAccount(topic gethcommon.Hash, db *state.StateDB) bool {
	potentialAddr := gethcommon.HexToAddress(topic.Hex())

	if topic.Hex()[2:len(zeroBytesHex)+2] != zeroBytesHex {
		return false
	}

	// A user address must have a non-zero nonce. This prevents accidental or malicious sending of funds to an
	// address matching a topic, forcing its events to become permanently private.
	if db.GetNonce(potentialAddr) != 0 {
		// If the address has code, it's a smart contract address instead.
		if db.GetCode(potentialAddr) == nil {
			return true
		}
	}
	return false
}

func (s *storageImpl) SetHeadBatchPointer(l2Head *core.Batch) error {
	dbBatch := s.db.NewBatch()

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.SetL2HeadBatch(dbBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateHeadRollup(l1Head *common.L1RootHash, l2Head *common.L2RootHash) error {
	dbBatch := s.db.NewBatch()
	if err := obscurorawdb.WriteL2HeadRollup(dbBatch, l1Head, l2Head); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateL1Head(l1Head common.L1RootHash) error {
	dbBatch := s.db.NewBatch()
	rawdb.WriteHeadHeaderHash(dbBatch, l1Head)
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new L1 head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateStateDB(hash common.L2RootHash) (*state.StateDB, error) {
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	// todo - snapshots?
	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}

	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	statedb, err := state.New(gethcommon.BigToHash(big.NewInt(0)), s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByHash(hash gethcommon.Hash) (types.Receipts, error) {
	number, err := obscurorawdb.ReadBatchNumber(s.db, hash)
	if err != nil {
		return nil, err
	}
	return obscurorawdb.ReadReceipts(s.db, hash, *number, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index, err := obscurorawdb.ReadTransaction(s.db, txHash)
	if err != nil {
		return nil, gethcommon.Hash{}, 0, 0, err
	}
	return tx, blockHash, blockNumber, index, nil
}

func (s *storageImpl) GetSender(txHash gethcommon.Hash) (gethcommon.Address, error) {
	tx, _, _, _, err := s.GetTransaction(txHash) //nolint:dogsled
	if err != nil {
		return gethcommon.Address{}, err
	}
	// todo - make the signer a field of the rollup chain
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not convert transaction to message to retrieve sender address in eth_getTransactionReceipt request. Cause: %w", err)
	}
	return msg.From(), nil
}

func (s *storageImpl) GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error) {
	return obscurorawdb.ReadContractTransaction(s.db, address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	_, blockHash, _, index, err := s.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	receipts, err := s.GetReceiptsByHash(blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve receipts for transaction. Cause: %w", err)
	}

	if len(receipts) <= int(index) {
		return nil, fmt.Errorf("receipt index not matching the transactions in block: %s", blockHash.Hex())
	}
	receipt := receipts[index]

	return receipt, nil
}

func (s *storageImpl) FetchAttestedKey(aggregator gethcommon.Address) (*ecdsa.PublicKey, error) {
	return obscurorawdb.ReadAttestationKey(s.db, aggregator)
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	return obscurorawdb.WriteAttestationKey(s.db, aggregator, key)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(dbBatch, *batch.Hash(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTxs(dbBatch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write batch to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1RootHash, messages common.CrossChainMessages) error {
	return obscurorawdb.StoreL1Messages(s.db, blockHash, messages, s.logger)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1RootHash) (common.CrossChainMessages, error) {
	return obscurorawdb.GetL1Messages(s.db, blockHash, s.logger)
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteRollup(dbBatch, rollup); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FilterLogs(requestingAccount *gethcommon.Address, filter *filters.FilterCriteria) ([]*types.Log, error) {
	result := []*types.Log{}
	queryParams := []any{}
	query := "select topic0, topic1, topic2, topic3, topic4, data, blockHash, blockNumber, txHash, txIdx, logIdx, address from events where 1==1 "
	if filter.BlockHash != nil {
		query += " AND blockHash = ?"
		queryParams = append(queryParams, filter.BlockHash)
	}
	if filter.FromBlock != nil {
		query += " AND blockNumber >= ?"
		queryParams = append(queryParams, filter.FromBlock)
	}
	if filter.ToBlock != nil {
		query += " AND blockNumber < ?"
		queryParams = append(queryParams, filter.ToBlock)
	}
	if len(filter.Addresses) > 0 {
		// todo ?
		query += " AND address in (?" + strings.Repeat(",?", len(filter.Addresses)-1) + ")"
		for _, address := range filter.Addresses {
			queryParams = append(queryParams, address)
		}
	}
	if len(filter.Topics) > 4 {
		return nil, fmt.Errorf("invalid filter. Too many topics")
	}
	if len(filter.Topics) > 0 {
		for i, sub := range filter.Topics {
			// empty rule set == wildcard
			if len(sub) > 0 {
				column := fmt.Sprintf("topic%d", i)
				query += " AND " + column + " in (?" + strings.Repeat(",?", len(sub)-1) + ")"
				for _, topic := range sub {
					queryParams = append(queryParams, topic)
				}
			}
		}
	}

	// Add relevancy rules
	//  An event is considered relevant to all account owners whose addresses are used as topics in the event.
	//	In case there are no account addresses in an event's topics, then the event is considered relevant to everyone (known as a "lifecycle event").
	query += " AND (lifecycleEvent OR (relAddress1=? OR relAddress2=? OR relAddress3=? OR relAddress4=?))"
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())
	queryParams = append(queryParams, requestingAccount.Bytes())

	rows, err := s.db.GetSQLDB().Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		l := types.Log{
			Topics: []gethcommon.Hash{},
		}
		var t0, t1, t2, t3, t4 sql2.NullString
		err := rows.Scan(&t0, &t1, &t2, &t3, &t4, &l.Data, &l.BlockHash, &l.BlockNumber, &l.TxHash, &l.TxIndex, &l.Index, &l.Address)
		if err != nil {
			return nil, fmt.Errorf("could not load log entry from db: %w", err)
		}
		if t0.Valid {
			l.Topics = append(l.Topics, hash(t0))
		}
		if t1.Valid {
			l.Topics = append(l.Topics, hash(t1))
		}
		if t2.Valid {
			l.Topics = append(l.Topics, hash(t2))
		}
		if t3.Valid {
			l.Topics = append(l.Topics, hash(t3))
		}
		if t4.Valid {
			l.Topics = append(l.Topics, hash(t4))
		}

		result = append(result, &l)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func hash(ns sql2.NullString) gethcommon.Hash {
	value, err := ns.Value()
	if err != nil {
		return [32]byte{}
	}
	s := value.(string)
	result := gethcommon.Hash{}
	result.SetBytes([]byte(s))
	return result
}
