package hostdb

import (
	"database/sql"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestCanStoreAndRetrieveBlock(t *testing.T) {
	db, _ := createSQLiteDB(t)
	block1 := createBlock(batchNumber)
	block2 := createBlock(batchNumber + 1)

	// verify we get ErrNoRows for a non-existent block
	randomHash := gethcommon.Hash{}
	randomHash.SetBytes(make([]byte, 32)) // 32 bytes for appropriate length
	dbtx, _ := db.NewDBTransaction()
	statements := db.GetSQLStatement()
	_, err := GetBlockId(dbtx.Tx, statements, randomHash)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows for non-existent block, got: %v", err)
	}
	dbtx.Rollback()

	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, &block1)
	if err != nil {
		t.Errorf("could not store block1: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block1: %s", err)
	}

	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, &block2)
	if err != nil {
		t.Errorf("could not store block2: %s", err)
	}
	err = dbtx.Write()
	if err != nil {
		t.Errorf("could not commit block2: %s", err)
	}

	dbtx, _ = db.NewDBTransaction()
	blockId, err := GetBlockId(dbtx.Tx, statements, block2.Hash())
	if err != nil {
		t.Errorf("stored block but could not retrieve block ID: %s", err)
	}
	if *blockId != 2 {
		t.Errorf("expected block ID 2, got %d", *blockId)
	}
	dbtx.Rollback()
}

func TestAddBlockWithForeignKeyConstraint(t *testing.T) {
	db, _ := createSQLiteDB(t)
	dbtx, _ := db.NewDBTransaction()
	statements := db.GetSQLStatement()
	metadata := createRollupMetadata(batchNumber - 10)
	rollup := createRollup(batchNumber)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)

	// add block
	err := AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()

	// add rollup referencing block
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup, &metadata, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	dbtx.Write()

	// this should fail due to the UNIQUE constraint on block_host.hash
	dbtx, _ = db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, statements, block.Header())
	if !strings.Contains(err.Error(), "UNIQUE constraint failed: block_host.hash") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}

	// verify the block still exists
	_, err = GetBlockId(dbtx.Tx, statements, block.Hash())
	if err != nil {
		t.Fatalf("failed to get block id after duplicate insert: %v", err)
	}

	dbtx.Rollback()
}

func createBlock(blockNum int64) types.Header {
	return types.Header{
		Number: big.NewInt(blockNum),
		Time:   uint64(time.Now().Unix()),
	}
}
