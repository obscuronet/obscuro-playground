package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/host/storage"
)

const (
	selectTxCount         = "SELECT total FROM transaction_count WHERE id = 1"
	selectBatch           = "SELECT sequence, full_hash, hash, height, ext_batch FROM batch_host"
	selectExtBatch        = "SELECT ext_batch FROM batch_host"
	selectLatestBatch     = "SELECT sequence, full_hash, hash, height, ext_batch FROM batch_host ORDER BY sequence DESC LIMIT 1"
	selectTxsAndBatch     = "SELECT t.hash FROM transactions_host t JOIN batch_host b ON t.b_sequence = b.sequence WHERE b.full_hash = ?"
	selectBatchSeqByTx    = "SELECT b_sequence FROM transactions_host WHERE hash = ?"
	selectTxBySeq         = "SELECT hash FROM transactions_host WHERE b_sequence = ?"
	insertBatch           = "INSERT INTO batch_host (sequence, full_hash, hash, height, ext_batch) VALUES (?, ?, ?, ?, ?)"
	insertTransactions    = "REPLACE INTO transactions_host (hash, b_sequence) VALUES (?, ?)"
	insertTxCountMariaDB  = "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON DUPLICATE KEY UPDATE total = VALUES(total)"
	insertTxCountSqliteDB = "INSERT INTO transaction_count (id, total) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET total = excluded.total;"
)

// AddBatch adds a batch and its header to the DB
func AddBatch(hostDB *storage.HostDB, batch *common.ExtBatch) error {
	db := hostDB.DB
	// mariadb tx context gets lost without this
	if !hostDB.InMem {
		useDbCmd := fmt.Sprintf("USE %s", hostDB.DBName)
		_, err := db.Exec(useDbCmd)
		if err != nil {
			return fmt.Errorf("failed to select database %s: %w", hostDB.DBName, err)
		}
	}

	// Check if the Batch is already stored
	_, err := GetBatchHeader(db, batch.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	extBatch, err := rlp.EncodeToBytes(batch)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions: %w", err)
	}

	_, err = tx.Exec(insertBatch,
		batch.SeqNo().Uint64(),       // sequence
		batch.Hash(),                 // full hash
		truncTo16(batch.Hash()),      // shortened hash
		batch.Header.Number.Uint64(), // height
		extBatch,                     // ext_batch
	)
	if err != nil {
		return fmt.Errorf("host failed to insert batch: %w", err)
	}

	if len(batch.TxHashes) > 0 {
		for _, transaction := range batch.TxHashes {
			_, err = tx.Exec(insertTransactions, transaction.Bytes(), batch.SeqNo().Uint64())
			if err != nil {
				return fmt.Errorf("failed to insert transaction with hash: %d", err)
			}
		}
	}

	var currentTotal int
	err = tx.QueryRow(selectTxCount).Scan(&currentTotal)
	if err != nil {
		return fmt.Errorf("failed to query transaction count: %w", err)
	}

	var insertTxCount string

	if hostDB.InMem {
		insertTxCount = insertTxCountSqliteDB
	} else {
		insertTxCount = insertTxCountMariaDB
	}
	newTotal := currentTotal + len(batch.TxHashes)
	_, err = tx.Exec(insertTxCount, 1, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update transaction count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch tx: %w", err)
	}

	return nil
}

// GetBatchListing returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListing(db *sql.DB, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	headBatch, err := GetCurrentHeadBatch(db)
	if err != nil {
		return nil, err
	}
	batchesFrom := headBatch.SequencerOrderNo.Uint64() - pagination.Offset
	batchesTo := int(batchesFrom) - int(pagination.Size) + 1

	if batchesTo <= 0 {
		batchesTo = 1
	}

	var batches []common.PublicBatch
	for i := batchesFrom; i >= uint64(batchesTo); i-- {
		batch, err := GetPublicBatchBySequenceNumber(db, i)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if batch != nil {
			batches = append(batches, *batch)
		}
	}

	return &common.BatchListingResponse{
		BatchesData: batches,
		Total:       uint64(len(batches)),
	}, nil
}

// GetBatchListingDeprecated returns latest batches given a pagination.
// For example, page 0, size 10 will return the latest 10 batches.
func GetBatchListingDeprecated(db *sql.DB, pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	headBatch, err := GetCurrentHeadBatch(db)
	if err != nil {
		return nil, err
	}
	batchesFrom := headBatch.SequencerOrderNo.Uint64() - pagination.Offset
	batchesTo := int(batchesFrom) - int(pagination.Size) + 1

	if batchesTo <= 0 {
		batchesTo = 1
	}

	var batches []common.PublicBatchDeprecated
	var txHashes []common.TxHash
	for i := batchesFrom; i >= uint64(batchesTo); i-- {
		batch, err := GetPublicBatchBySequenceNumber(db, i)
		if batch == nil {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get batch by seq no: %w", err)
		}

		txHashes, err = GetTxsBySequenceNumber(db, batch.Header.SequencerOrderNo.Uint64())
		if err != nil {
			return nil, fmt.Errorf("failed to get tx hashes by seq no: %w", err)
		}
		if batch == nil || batch.Header == nil {
			return nil, fmt.Errorf("batch or batch header is nil")
		} else {
			publicBatchDeprecated := common.PublicBatchDeprecated{
				BatchHeader: *batch.Header,
				TxHashes:    txHashes,
			}
			batches = append(batches, publicBatchDeprecated)
		}
	}

	return &common.BatchListingResponseDeprecated{
		BatchesData: batches,
		Total:       uint64(len(batches)),
	}, nil
}

// GetPublicBatchBySequenceNumber returns the batch with the given sequence number.
func GetPublicBatchBySequenceNumber(db *sql.DB, seqNo uint64) (*common.PublicBatch, error) {
	return fetchPublicBatch(db, " WHERE sequence=?", seqNo)
}

// GetTxsBySequenceNumber returns the transaction hashes with sequence number.
func GetTxsBySequenceNumber(db *sql.DB, seqNo uint64) ([]common.TxHash, error) {
	return fetchTx(db, seqNo)
}

// GetBatchBySequenceNumber returns the ext batch for a given sequence number.
func GetBatchBySequenceNumber(db *sql.DB, seqNo uint64) (*common.ExtBatch, error) {
	return fetchFullBatch(db, " WHERE sequence=?", seqNo)
}

// GetCurrentHeadBatch retrieves the current head batch with the largest sequence number (or height).
func GetCurrentHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	return fetchHeadBatch(db)
}

// GetBatchHeader returns the batch header given the hash.
func GetBatchHeader(db *sql.DB, hash gethcommon.Hash) (*common.BatchHeader, error) {
	return fetchBatchHeader(db, " where hash=?", truncTo16(hash))
}

// GetBatchHashByNumber returns the hash of a batch given its number.
func GetBatchHashByNumber(db *sql.DB, number *big.Int) (*gethcommon.Hash, error) {
	batch, err := fetchBatchHeader(db, " where height=?", number.Uint64())
	if err != nil {
		return nil, err
	}
	l2BatchHash := batch.Hash()
	return &l2BatchHash, nil
}

// GetHeadBatchHeader returns the latest batch header.
func GetHeadBatchHeader(db *sql.DB) (*common.BatchHeader, error) {
	batch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, err
	}
	return batch.Header, nil
}

// GetBatchNumber returns the height of the batch containing the given transaction hash.
func GetBatchNumber(db *sql.DB, txHash gethcommon.Hash) (*big.Int, error) {
	txBytes := txHash.Bytes()
	batchHeight, err := fetchBatchNumber(db, txBytes)
	if err != nil {
		return nil, err
	}
	return batchHeight, nil
}

// GetBatchTxs returns the transaction hashes of the batch with the given hash.
func GetBatchTxs(db *sql.DB, batchHash gethcommon.Hash) ([]gethcommon.Hash, error) {
	rows, err := db.Query(selectTxsAndBatch, batchHash)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var transactions []gethcommon.Hash
	for rows.Next() {
		var txHashBytes []byte
		if err := rows.Scan(&txHashBytes); err != nil {
			return nil, fmt.Errorf("failed to scan transaction hash: %w", err)
		}
		txHash := gethcommon.BytesToHash(txHashBytes)
		transactions = append(transactions, txHash)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error looping through transacion rows: %w", err)
	}

	return transactions, nil
}

// GetTotalTxCount returns the total number of batched transactions.
func GetTotalTxCount(db *sql.DB) (*big.Int, error) {
	var totalCount int
	err := db.QueryRow(selectTxCount).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve total transaction count: %w", err)
	}
	return big.NewInt(int64(totalCount)), nil
}

// GetPublicBatch returns the batch with the given hash.
func GetPublicBatch(db *sql.DB, hash common.L2BatchHash) (*common.PublicBatch, error) {
	return fetchPublicBatch(db, " where b.hash=?", truncTo16(hash))
}

// GetBatchByTx returns the batch with the given hash.
func GetBatchByTx(db *sql.DB, txHash gethcommon.Hash) (*common.ExtBatch, error) {
	var seqNo uint64
	err := db.QueryRow(selectBatchSeqByTx, txHash).Scan(&seqNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return GetBatchBySequenceNumber(db, seqNo)
}

// GetBatchByHash returns the batch with the given hash.
func GetBatchByHash(db *sql.DB, hash common.L2BatchHash) (*common.ExtBatch, error) {
	return fetchFullBatch(db, " where hash=?", truncTo16(hash))
}

// GetLatestBatch returns the head batch header
func GetLatestBatch(db *sql.DB) (*common.BatchHeader, error) {
	headBatch, err := fetchHeadBatch(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head batch: %w", err)
	}
	return headBatch.Header, nil
}

// GetBatchByHeight returns the batch header given the height
func GetBatchByHeight(db *sql.DB, height *big.Int) (*common.BatchHeader, error) {
	headBatch, err := fetchBatchHeader(db, "where height=?", height.Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to batch header: %w", err)
	}
	return headBatch, nil
}

func fetchBatchHeader(db *sql.DB, whereQuery string, args ...any) (*common.BatchHeader, error) {
	var extBatch []byte
	query := selectExtBatch + " " + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&extBatch)
	} else {
		err = db.QueryRow(query).Scan(&extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	// Decode batch
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode batch header. Cause: %w", err)
	}
	return b.Header, nil
}

func fetchBatchNumber(db *sql.DB, args ...any) (*big.Int, error) {
	var seqNo uint64
	query := selectBatchSeqByTx
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&seqNo)
	} else {
		err = db.QueryRow(query).Scan(&seqNo)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	batch, err := GetPublicBatchBySequenceNumber(db, seqNo)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch by seq no. Cause: %w", err)
	}
	return batch.Height, nil
}

func fetchPublicBatch(db *sql.DB, whereQuery string, args ...any) (*common.PublicBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + " " + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	batch := &common.PublicBatch{
		SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
		Hash:             hash,
		FullHash:         fullHash,
		Height:           new(big.Int).SetInt64(int64(heightInt64)),
		TxCount:          new(big.Int).SetInt64(int64(len(b.TxHashes))),
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
	}

	return batch, nil
}

func fetchFullBatch(db *sql.DB, whereQuery string, args ...any) (*common.ExtBatch, error) {
	var sequenceInt64 uint64
	var fullHash common.TxHash
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	query := selectBatch + whereQuery

	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	} else {
		err = db.QueryRow(query).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	return &b, nil
}

func fetchHeadBatch(db *sql.DB) (*common.PublicBatch, error) {
	var sequenceInt64 int
	var fullHash gethcommon.Hash // common.Hash
	var hash []byte
	var heightInt64 int
	var extBatch []byte

	err := db.QueryRow(selectLatestBatch).Scan(&sequenceInt64, &fullHash, &hash, &heightInt64, &extBatch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch current head batch: %w", err)
	}

	var b common.ExtBatch
	err = rlp.DecodeBytes(extBatch, &b)
	if err != nil {
		return nil, fmt.Errorf("could not decode ext batch. Cause: %w", err)
	}

	batch := &common.PublicBatch{
		SequencerOrderNo: new(big.Int).SetInt64(int64(sequenceInt64)),
		Hash:             hash,
		FullHash:         fullHash,
		Height:           new(big.Int).SetInt64(int64(heightInt64)),
		TxCount:          new(big.Int).SetInt64(int64(len(b.TxHashes))),
		Header:           b.Header,
		EncryptedTxBlob:  b.EncryptedTxBlob,
	}

	return batch, nil
}

func fetchTx(db *sql.DB, seqNo uint64) ([]common.TxHash, error) {
	rows, err := db.Query(selectTxBySeq, seqNo)
	if err != nil {
		return nil, fmt.Errorf("query execution for select txs failed: %w", err)
	}
	defer rows.Close()

	var transactions []gethcommon.Hash
	for rows.Next() {
		var txHashBytes []byte
		if err := rows.Scan(&txHashBytes); err != nil {
			return nil, fmt.Errorf("failed to scan transaction hash: %w", err)
		}
		txHash := gethcommon.BytesToHash(txHashBytes)
		transactions = append(transactions, txHash)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error looping through transacion rows: %w", err)
	}

	return transactions, nil
}