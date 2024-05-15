package hostdb

import (
	"fmt"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
)

func TestGetTransactionListing(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHash12 := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHash12)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db.GetSQLStatement(), &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHash34 := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHash34)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHash56 := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringFive")), gethcommon.BytesToHash([]byte("magicStringSix"))}
	batchThree := createBatch(batchNumber+2, txHash56)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchThree)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	// offset 0, size 3
	txListing, err := GetTransactionListing(db, &common.QueryPagination{Offset: 0, Size: 3})
	if err != nil {
		t.Errorf("could not get tx listing. Cause: %s", err)
	}

	// should be three elements
	if big.NewInt(int64(txListing.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// second element should be in the third batch as they're descending
	if txListing.TransactionsData[1].BatchHeight.Cmp(batchThree.Header.Number) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// third element should be in the second batch
	if txListing.TransactionsData[2].BatchHeight.Cmp(batchTwo.Header.Number) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// page 1, size 3
	txListing1, err := GetTransactionListing(db, &common.QueryPagination{Offset: 3, Size: 3})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(txListing1.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// first element should be in the second batch
	if txListing1.TransactionsData[0].BatchHeight.Cmp(batchTwo.Header.Number) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// third element should be in the first batch
	if txListing1.TransactionsData[2].BatchHeight.Cmp(batchOne.Header.Number) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	// size overflow, only 6 elements
	txListing2, err := GetTransactionListing(db, &common.QueryPagination{Offset: 0, Size: 7})
	if err != nil {
		t.Errorf("could not get batch listing. Cause: %s", err)
	}

	// should be 6 elements
	if big.NewInt(int64(txListing2.Total)).Cmp(big.NewInt(6)) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}
}

func TestGetTransaction(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHash1 := gethcommon.BytesToHash([]byte("magicStringOne"))
	txHash2 := gethcommon.BytesToHash([]byte("magicStringOne"))
	txHashes := []common.L2TxHash{txHash1, txHash2}
	batchOne := createBatch(batchNumber, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db.GetSQLStatement(), &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	tx, err := GetTransaction(db, txHash2)
	if err != nil {
		t.Errorf("was not able to get transaction. Cause: %s", err)
	}

	if tx.BatchHeight.Cmp(big.NewInt(batchNumber)) != 0 {
		t.Errorf("tx batch height was not retrieved correctly")
	}
	if tx.TransactionHash.Cmp(txHash2) != 0 {
		t.Errorf("tx hash was not retrieved correctly")
	}
}

func TestCanRetrieveTotalNumberOfTransactions(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db.GetSQLStatement(), &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}
	dbtx.Write()

	totalTxs, err := GetTotalTxCount(db)
	if err != nil {
		t.Errorf("was not able to read total number of transactions. Cause: %s", err)
	}

	if int(totalTxs.Int64()) != len(txHashesOne)+len(txHashesTwo) {
		t.Errorf("total number of batch transactions was not stored correctly")
	}
}

func TestGetTransactionListingLarge(t *testing.T) {
	db, _ := createSQLiteDB(t)
	var txHashes []common.L2TxHash
	for i := 0; i < 1000; i++ {
		txHash := gethcommon.BytesToHash([]byte(fmt.Sprintf("magicString%d", i+1)))
		txHashes = append(txHashes, txHash)
	}
	batchOne := createBatch(batchNumber, txHashes)
	dbtx, _ := db.NewDBTransaction()
	err := AddBatch(dbtx, db.GetSQLStatement(), &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	// offset 1000, size 2
	txListing, err := GetTransactionListing(db, &common.QueryPagination{Offset: 1000, Size: 2})
	if err != nil {
		t.Errorf("could not get tx listing. Cause: %s", err)
	}

	if big.NewInt(int64(txListing.Total)).Cmp(big.NewInt(0)) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}

	txListing1, err := GetTransactionListing(db, &common.QueryPagination{Offset: 995, Size: 5})
	if err != nil {
		t.Errorf("could not get tx listing. Cause: %s", err)
	}

	if big.NewInt(int64(txListing1.Total)).Cmp(big.NewInt(5)) != 0 {
		t.Errorf("tx listing was not paginated correctly")
	}
}
