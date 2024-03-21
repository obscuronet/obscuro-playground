package hostdb

import (
	"errors"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common"
)

// An arbitrary number to put in the header, to check that the header is retrieved correctly from the DB.

func TestCanStoreAndRetrieveBatchHeader(t *testing.T) {
	db, _ := createSQLiteDB(t)
	batch := createBatch(batchNumber, []common.L2TxHash{})
	err := AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetBatchHeader(db.DB, batch.Header.Hash())
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch header was not stored correctly")
	}
}

func TestUnknownBatchHeaderReturnsNotFound(t *testing.T) {
	db, _ := createSQLiteDB(t)
	header := types.Header{}

	_, err := GetBatchHeader(db.DB, header.Hash())
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch header but was able to retrieve it")
	}
}

func TestHigherNumberBatchBecomesBatchHeader(t *testing.T) { //nolint:dupl
	db, _ := createSQLiteDB(t)
	batchOne := createBatch(batchNumber, []common.L2TxHash{})
	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo := createBatch(batchNumber+1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetHeadBatchHeader(db.DB)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchTwo.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestLowerNumberBatchDoesNotBecomeBatchHeader(t *testing.T) { //nolint:dupl
	db, _ := createSQLiteDB(t)
	batchOne := createBatch(batchNumber, []common.L2TxHash{})
	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTwo := createBatch(batchNumber-1, []common.L2TxHash{})
	if err != nil {
		t.Errorf("could not create batch. Cause: %s", err)
	}

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHeader, err := GetHeadBatchHeader(db.DB)
	if err != nil {
		t.Errorf("stored batch but could not retrieve header. Cause: %s", err)
	}
	if batchHeader.Number.Cmp(batchOne.Header.Number) != 0 {
		t.Errorf("head batch was not set correctly")
	}
}

func TestHeadBatchHeaderIsNotSetInitially(t *testing.T) {
	db, _ := createSQLiteDB(t)
	_, err := GetHeadBatchHeader(db.DB)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("head batch was set, but no batchs had been written")
	}
}

func TestCanRetrieveBatchHashByNumber(t *testing.T) {
	db, _ := createSQLiteDB(t)
	batch := createBatch(batchNumber, []common.L2TxHash{})
	err := AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchHash, err := GetBatchHashByNumber(db.DB, batch.Header.Number)
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers hash by number. Cause: %s", err)
	}
	if *batchHash != batch.Header.Hash() {
		t.Errorf("batch hash was not stored correctly against number")
	}
}

func TestUnknownBatchNumberReturnsNotFound(t *testing.T) {
	db, _ := createSQLiteDB(t)
	header := types.Header{Number: big.NewInt(10)}

	_, err := GetBatchHashByNumber(db.DB, header.Number)
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch hash but was able to retrieve it")
	}
}

func TestCanRetrieveBatchNumberByTxHash(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHash := gethcommon.BytesToHash([]byte("magicString"))
	batch := createBatch(batchNumber, []common.L2TxHash{txHash})
	err := AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	extBatch, err := GetBatchByTx(db.DB, txHash)
	if err != nil {
		t.Errorf("stored batch but could not retrieve batch by transaction hash. Cause: %s", err)
	}
	if extBatch.Header.Number.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
	batchNumber, err := GetBatchNumber(db.DB, txHash)
	if err != nil {
		t.Errorf("stored batch but could not retrieve number by transaction hash. Cause: %s", err)
	}
	if batchNumber.Cmp(batch.Header.Number) != 0 {
		t.Errorf("batch number was not stored correctly against transaction hash")
	}
}

func TestUnknownBatchTxHashReturnsNotFound(t *testing.T) {
	db, _ := createSQLiteDB(t)

	_, err := GetBatchNumber(db.DB, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveBatchTransactions(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashes := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batch := createBatch(batchNumber, txHashes)
	err := AddBatch(db, &batch)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batchTxs, err := GetBatchTxs(db.DB, batch.Header.Hash())
	if err != nil {
		t.Errorf("stored batch but could not retrieve headers transactions. Cause: %s", err)
	}
	if len(batchTxs) != len(txHashes) {
		t.Errorf("batch transactions were not stored correctly")
	}
	for idx, batchTx := range batchTxs {
		if batchTx != txHashes[idx] {
			t.Errorf("batch transactions were not stored correctly")
		}
	}
}

func TestTransactionsForUnknownBatchReturnsNotFound(t *testing.T) {
	db, _ := createSQLiteDB(t)

	_, err := GetBatchNumber(db.DB, gethcommon.BytesToHash([]byte("magicString")))
	if !errors.Is(err, errutil.ErrNotFound) {
		t.Errorf("did not store batch number but was able to retrieve it")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)

	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	totalTxs, err := GetTotalTxCount(db.DB)
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of batch transactions was not stored correctly")
	}
}

func TestGetLatestBatch(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)

	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	batch, err := GetLatestBatch(db.DB)
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(batch.SequencerOrderNo.Uint64()) != int(batchTwo.SeqNo().Uint64()) {
		t.Errorf("latest batch was not retrieved correctly")
	}
}

func TestGetBatchListing(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)

	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesThree := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringFive")), gethcommon.BytesToHash([]byte("magicStringSix"))}
	batchThree := createBatch(batchNumber+2, txHashesThree)

	err = AddBatch(db, &batchThree)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	// page 1, size 2
	batchListing, err := GetBatchListing(db.DB, &common.QueryPagination{Offset: 1, Size: 2})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be two elements
	if big.NewInt(int64(batchListing.Total)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// first element should be the second batch
	if batchListing.BatchesData[0].SequencerOrderNo.Cmp(batchTwo.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 0, size 3
	batchListing1, err := GetBatchListing(db.DB, &common.QueryPagination{Offset: 0, Size: 3})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// first element should be the most recent batch since they're in descending order
	if batchListing1.BatchesData[0].SequencerOrderNo.Cmp(batchThree.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// should be 3 elements
	if big.NewInt(int64(batchListing1.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 0, size 4
	batchListing2, err := GetBatchListing(db.DB, &common.QueryPagination{Offset: 0, Size: 4})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(batchListing2.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 5, size 1
	rollupListing3, err := GetBatchListing(db.DB, &common.QueryPagination{Offset: 5, Size: 1})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 0 elements
	if big.NewInt(int64(rollupListing3.Total)).Cmp(big.NewInt(0)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}
}

func TestGetBatchListingDeprecated(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)

	err := AddBatch(db, &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(db, &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesThree := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringFive")), gethcommon.BytesToHash([]byte("magicStringSix"))}
	batchThree := createBatch(batchNumber+2, txHashesThree)

	err = AddBatch(db, &batchThree)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	// page 1, size 2
	batchListing, err := GetBatchListingDeprecated(db.DB, &common.QueryPagination{Offset: 1, Size: 2})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be two elements
	if big.NewInt(int64(batchListing.Total)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// first element should be the second batch
	if batchListing.BatchesData[0].BatchHeader.SequencerOrderNo.Cmp(batchTwo.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 0, size 3
	batchListing1, err := GetBatchListingDeprecated(db.DB, &common.QueryPagination{Offset: 0, Size: 3})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// first element should be the most recent batch since they're in descending order
	if batchListing1.BatchesData[0].BatchHeader.SequencerOrderNo.Cmp(batchThree.SeqNo()) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// should be 3 elements
	if big.NewInt(int64(batchListing1.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}

	// page 0, size 4
	batchListing2, err := GetBatchListingDeprecated(db.DB, &common.QueryPagination{Offset: 0, Size: 4})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(batchListing2.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 5, size 1
	rollupListing3, err := GetBatchListing(db.DB, &common.QueryPagination{Offset: 5, Size: 1})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 0 elements
	if big.NewInt(int64(rollupListing3.Total)).Cmp(big.NewInt(0)) != 0 {
		t.Errorf("batch listing was not paginated correctly")
	}
}

// TODO Get Batch by height
// TODO Get Batch by TX hash
// TODO Duplicate TX hash test
func createBatch(batchNum int64, txHashes []common.L2BatchHash) common.ExtBatch {
	header := common.BatchHeader{
		SequencerOrderNo: big.NewInt(batchNum),
		Number:           big.NewInt(batchNum),
	}
	batch := common.ExtBatch{
		Header:   &header,
		TxHashes: txHashes,
	}

	return batch
}
