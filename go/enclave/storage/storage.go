package storage

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/rlp"

	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/obscuronet/go-obscuro/go/common/syserr"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// todo - this will require a dedicated table when updates are implemented
const masterSeedCfg = "MASTER_SEED"

type storageImpl struct {
	db          enclavedb.EnclaveDB
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorageFromConfig(config *config.EnclaveConfig, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, chainConfig, logger)
}

func NewStorage(backingDB enclavedb.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	cacheConfig := &gethcore.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
		SnapshotLimit:  256,
		SnapshotWait:   true,
	}

	return &storageImpl{
		db: backingDB,
		stateDB: state.NewDatabaseWithConfig(backingDB, &trie.Config{
			Cache:     cacheConfig.TrieCleanLimit,
			Journal:   cacheConfig.TrieCleanJournal,
			Preimages: cacheConfig.Preimages,
		}),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) TrieDB() *trie.Database {
	return s.stateDB.TrieDB()
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	return enclavedb.ReadCurrentHeadBatch(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo() (*big.Int, error) {
	return enclavedb.ReadCurrentSequencerNo(s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(hash common.L2BatchHash) (*core.Batch, error) {
	return enclavedb.ReadBatchByHash(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error) {
	return enclavedb.ReadBatchHeader(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	return enclavedb.ReadCanonicalBatchByHeight(s.db.GetSQLDB(), height)
}

func (s *storageImpl) StoreBlock(b *types.Block, chainFork *common.ChainFork) error {
	dbBatch := s.db.NewDBTransaction()
	if chainFork != nil && chainFork.IsFork() {
		s.logger.Info(fmt.Sprintf("Fork. %+v.", chainFork))
		enclavedb.UpdateCanonicalBlocks(dbBatch, chainFork.CanonicalPath, chainFork.NonCanonicalPath)
	} else {
		enclavedb.UpdateCanonicalBlocks(dbBatch, nil, nil)
	}

	if err := enclavedb.WriteBlock(dbBatch, b.Header()); err != nil {
		return fmt.Errorf("could not store block %s. Cause: %w", b.Hash(), err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not store block %s. Cause: %w", b.Hash(), err)
	}
	return nil
}

func (s *storageImpl) FetchBlock(blockHash common.L1BlockHash) (*types.Block, error) {
	return enclavedb.FetchBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	return enclavedb.FetchHeadBlock(s.db.GetSQLDB())
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	enc, err := rlp.EncodeToBytes(secret)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	_, err = enclavedb.WriteConfig(s.db.GetSQLDB(), masterSeedCfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	var ss crypto.SharedEnclaveSecret

	cfg, err := enclavedb.FetchConfig(s.db.GetSQLDB(), masterSeedCfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(block.ParentHash())
	if err != nil {
		s.logger.Warn("Could not find block with hash", log.BlockHashKey, block.ParentHash(), log.ErrKey, err)
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
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

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error) {
	return enclavedb.ReadHeadBatchForBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error) {
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("could not create state DB. Cause: %w", err))
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

// GetReceiptsByBatchHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByBatchHash(hash gethcommon.Hash) (types.Receipts, error) {
	return enclavedb.ReadReceiptsByBatchHash(s.db.GetSQLDB(), hash, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	return enclavedb.ReadTransaction(s.db.GetSQLDB(), txHash)
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
	return enclavedb.GetContractCreationTx(s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	return enclavedb.ReadReceipt(s.db.GetSQLDB(), txHash, s.chainConfig)
}

func (s *storageImpl) FetchAttestedKey(address gethcommon.Address) (*ecdsa.PublicKey, error) {
	key, err := enclavedb.FetchAttKey(s.db.GetSQLDB(), address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := gethcrypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	_, err := enclavedb.WriteAttKey(s.db.GetSQLDB(), aggregator, gethcrypto.CompressPubkey(key))
	return err
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	return enclavedb.ReadBatchBySeqNo(s.db.GetSQLDB(), seqNum)
}

func (s *storageImpl) FetchBatchesByBlock(block common.L1BlockHash) ([]*core.Batch, error) {
	return enclavedb.ReadBatchesByBlock(s.db.GetSQLDB(), block)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt) error {
	// sanity check that this is not overlapping
	prev, err := s.FetchBatchBySeqNo(batch.SeqNo().Uint64())
	if err == nil && !bytes.Equal(prev.Hash().Bytes(), batch.Hash().Bytes()) {
		return fmt.Errorf("a different batch with same sequence number already exists: %d", batch.SeqNo())
	}

	dbTx := s.db.NewDBTransaction()
	s.logger.Trace("write batch", "hash", batch.Hash(), "l1_proof", batch.Header.L1Proof)
	if err := enclavedb.WriteBatchAndTransactions(dbTx, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	if len(receipts) > 0 {
		err := s.storeReceipts(batch, receipts, dbTx)
		if err != nil {
			return err
		}
	}

	err = dbTx.Write()
	if err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}
	return nil
}

func (s *storageImpl) storeReceipts(batch *core.Batch, receipts []*types.Receipt, dbTx enclavedb.DBTransaction) error {
	for _, receipt := range receipts {
		s.logger.Trace("store receipt", "txHash", receipt.TxHash, "batch", receipt.BlockHash)
	}
	if err := enclavedb.WriteReceipts(dbTx, receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}

	if batch.Number().Int64() > 1 {
		stateDB, err := s.CreateStateDB(batch.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
		}

		err2 := enclavedb.StoreEventLogs(dbTx, receipts, stateDB)
		if err2 != nil {
			return fmt.Errorf("could not save logs %w", err2)
		}
	}
	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	return enclavedb.WriteL1Messages(s.db.GetSQLDB(), blockHash, messages)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	return enclavedb.FetchL1Messages(s.db.GetSQLDB(), blockHash)
}

const enclaveKeyKey = "ek"

func (s *storageImpl) StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error {
	if enclaveKey == nil {
		return errors.New("enclaveKey cannot be nil")
	}
	keyBytes := gethcrypto.FromECDSA(enclaveKey)

	_, err := enclavedb.WriteConfig(s.db.GetSQLDB(), enclaveKeyKey, keyBytes)
	return err
}

func (s *storageImpl) GetEnclaveKey() (*ecdsa.PrivateKey, error) {
	keyBytes, err := enclavedb.FetchConfig(s.db.GetSQLDB(), enclaveKeyKey)
	if err != nil {
		return nil, err
	}
	enclaveKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to construct ECDSA private key from enclave key bytes - %w", err)
	}
	return enclaveKey, nil
}

func (s *storageImpl) StoreRollup(rollup *common.ExtRollup) error {
	dbBatch := s.db.NewDBTransaction()

	if err := enclavedb.WriteRollup(dbBatch, rollup.Header); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	return enclavedb.DebugGetLogs(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogs(
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	return enclavedb.FilterLogs(s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics)
}

func (s *storageImpl) GetContractCount() (*big.Int, error) {
	return enclavedb.ReadContractCreationCount(s.db.GetSQLDB())
}

func (s *storageImpl) GetPublicTxsBySender(address *gethcommon.Address) ([]common.PublicTxData, error) {
	return enclavedb.ReadPublicTxsBySender(s.db.GetSQLDB(), address)
}
