package hostdb

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/common"
)

func TestCanStoreAndRetrieveRollup(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	metadata := createRollupMetadata(batchNumber - 10)
	rollup := createRollup(batchNumber)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup, &metadata, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	dbtx.Write()

	extRollup, err := GetExtRollup(db, rollup.Header.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve ext rollup. Cause: %s", err)
	}

	rollupHeader, err := GetRollupHeader(db, rollup.Header.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve header. Cause: %s", err)
	}
	if big.NewInt(int64(rollupHeader.LastBatchSeqNo)).Cmp(big.NewInt(batchNumber)) != 0 {
		t.Errorf("rollup header was not stored correctly")
	}

	if rollup.Hash() != extRollup.Hash() {
		t.Errorf("rollup was not stored correctly")
	}
}

func TestGetRollupByBlockHash(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	metadata := createRollupMetadata(batchNumber - 10)
	rollup := createRollup(batchNumber)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup, &metadata, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	dbtx.Write()
	rollupHeader, err := GetRollupHeaderByBlock(db, block.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve header. Cause: %s", err)
	}
	if big.NewInt(int64(rollupHeader.LastBatchSeqNo)).Cmp(big.NewInt(batchNumber)) != 0 {
		t.Errorf("rollup header was not stored correctly")
	}
}

func TestGetLatestRollup(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	rollup1FirstSeq := int64(batchNumber - 10)
	rollup1LastSeq := int64(batchNumber)
	metadata1 := createRollupMetadata(rollup1FirstSeq)
	rollup1 := createRollup(rollup1LastSeq)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup1, &metadata1, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	// Needed to increment the timestamp
	time.Sleep(1 * time.Second)

	rollup2FirstSeq := int64(batchNumber + 1)
	rollup2LastSeq := int64(batchNumber + 10)
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup2, &metadata2, block.Header())
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}
	dbtx.Write()

	latestHeader, err := GetLatestRollup(db)
	if err != nil {
		t.Errorf("could not get latest rollup. Cause: %s", err)
	}

	if latestHeader.LastBatchSeqNo != uint64(rollup2LastSeq) {
		t.Errorf("latest rollup was not updated correctly")
	}
}

func TestGetRollupBySeqNo(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	rollup1FirstSeq := int64(batchNumber - 10)
	rollup1LastSeq := int64(batchNumber)
	metadata1 := createRollupMetadata(rollup1FirstSeq)
	rollup1 := createRollup(rollup1LastSeq)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup1, &metadata1, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}
	// Needed to increment the timestamp
	time.Sleep(1 * time.Second)

	rollup2FirstSeq := int64(batchNumber + 1) // 778
	rollup2LastSeq := int64(batchNumber + 10) // 787
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup2, &metadata2, block.Header())
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}
	dbtx.Write()

	rollup, err := GetRollupBySeqNo(db, batchNumber+5)
	if err != nil {
		t.Errorf("could not get latest rollup. Cause: %s", err)
	}

	// should fetch the second rollup
	if rollup.LastSeq.Cmp(big.NewInt(int64(rollup2.Header.LastBatchSeqNo))) != 0 {
		t.Errorf("rollup was not fetched correctly")
	}

	rollup, err = GetRollupBySeqNo(db, batchNumber-5)
	if err != nil {
		t.Errorf("could not get latest rollup. Cause: %s", err)
	}
	// should fetch the first rollup
	if rollup.LastSeq.Cmp(big.NewInt(int64(rollup1.Header.LastBatchSeqNo))) != 0 {
		t.Errorf("rollup was not fetched correctly")
	}
}

func TestGetRollupListing(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	rollup1FirstSeq := int64(batchNumber - 10)
	rollup1LastSeq := int64(batchNumber)
	metadata1 := createRollupMetadata(rollup1FirstSeq)
	rollup1 := createRollup(rollup1LastSeq)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup1, &metadata1, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	rollup2FirstSeq := int64(batchNumber + 1)
	rollup2LastSeq := int64(batchNumber + 10)
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup2, &metadata2, block.Header())
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}

	rollup3FirstSeq := int64(batchNumber + 11)
	rollup3LastSeq := int64(batchNumber + 20)
	metadata3 := createRollupMetadata(rollup3FirstSeq)
	rollup3 := createRollup(rollup3LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup3, &metadata3, block.Header())
	dbtx.Write()
	if err != nil {
		t.Errorf("could not store rollup 3. Cause: %s", err)
	}

	// page 1, size 2
	rollupListing, err := GetRollupListing(db, &common.QueryPagination{Offset: 1, Size: 2})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(rollupListing.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// First element should be the second rollup
	if rollupListing.RollupsData[0].LastSeq.Cmp(big.NewInt(rollup2LastSeq)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}
	if rollupListing.RollupsData[0].FirstSeq.Cmp(big.NewInt(rollup2FirstSeq)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 0, size 3
	rollupListing1, err := GetRollupListing(db, &common.QueryPagination{Offset: 0, Size: 3})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// First element should be the most recent rollup since they're in descending order
	if rollupListing1.RollupsData[0].LastSeq.Cmp(big.NewInt(rollup3LastSeq)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}
	if rollupListing1.RollupsData[0].FirstSeq.Cmp(big.NewInt(rollup3FirstSeq)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// should be 3 elements
	if big.NewInt(int64(rollupListing1.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 0, size 4
	rollupListing2, err := GetRollupListing(db, &common.QueryPagination{Offset: 0, Size: 4})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(rollupListing2.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 5, size 1
	rollupListing3, err := GetRollupListing(db, &common.QueryPagination{Offset: 5, Size: 1})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(rollupListing3.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}
}

func TestGetRollupByHash(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	rollup1FirstSeq := int64(batchNumber - 10)
	rollup1LastSeq := int64(batchNumber)
	metadata1 := createRollupMetadata(rollup1FirstSeq)
	rollup1 := createRollup(rollup1LastSeq)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err = AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup1, &metadata1, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	rollup2FirstSeq := int64(batchNumber + 1)
	rollup2LastSeq := int64(batchNumber + 10)
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup2, &metadata2, block.Header())
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}
	dbtx.Write()

	publicRollup, err := GetRollupByHash(db, rollup2.Header.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve public rollup. Cause: %s", err)
	}

	if publicRollup.FirstSeq.Cmp(big.NewInt(batchNumber+1)) != 0 {
		t.Errorf("rollup was not stored correctly")
	}

	if publicRollup.LastSeq.Cmp(big.NewInt(batchNumber+10)) != 0 {
		t.Errorf("rollup was not stored correctly")
	}
}

func TestGetRollupBatches(t *testing.T) {
	db, _ := createSQLiteDB(t)
	txHashesOne := []common.L2TxHash{gethcommon.BytesToHash([]byte("magicStringOne")), gethcommon.BytesToHash([]byte("magicStringTwo"))}
	batchOne := createBatch(batchNumber, txHashesOne)
	block := types.NewBlock(&types.Header{}, nil, nil, nil)
	dbtx, _ := db.NewDBTransaction()
	err := AddBlock(dbtx.Tx, db.GetSQLStatement(), block.Header())
	if err != nil {
		t.Errorf("could not store block. Cause: %s", err)
	}
	dbtx.Write()
	dbtx, _ = db.NewDBTransaction()
	err = AddBatch(dbtx, db.GetSQLStatement(), &batchOne)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesTwo := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringThree")), gethcommon.BytesToHash([]byte("magicStringFour"))}
	batchTwo := createBatch(batchNumber+1, txHashesTwo)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchTwo)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesThree := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringFive")), gethcommon.BytesToHash([]byte("magicStringSix"))}
	batchThree := createBatch(batchNumber+2, txHashesThree)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchThree)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	txHashesFour := []gethcommon.Hash{gethcommon.BytesToHash([]byte("magicStringSeven")), gethcommon.BytesToHash([]byte("magicStringEight"))}
	batchFour := createBatch(batchNumber+3, txHashesFour)

	err = AddBatch(dbtx, db.GetSQLStatement(), &batchFour)
	if err != nil {
		t.Errorf("could not store batch. Cause: %s", err)
	}

	rollup1FirstSeq := int64(batchNumber)
	rollup1LastSeq := int64(batchNumber + 1)
	metadata1 := createRollupMetadata(rollup1FirstSeq)
	rollup1 := createRollup(rollup1LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup1, &metadata1, block.Header())
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	rollup2FirstSeq := int64(batchNumber + 2)
	rollup2LastSeq := int64(batchNumber + 3)
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)
	err = AddRollup(dbtx, db.GetSQLStatement(), &rollup2, &metadata2, block.Header())
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}
	dbtx.Write()

	// rollup one contains batches 1 & 2
	batchListing, err := GetRollupBatches(db, rollup1.Hash())
	if err != nil {
		t.Errorf("could not get rollup batches. Cause: %s", err)
	}

	// should be two elements
	if big.NewInt(int64(batchListing.Total)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("batch listing was not calculated correctly")
	}

	// second element should be batch 1 as we're ordering by height descending
	if batchListing.BatchesData[1].Header.SequencerOrderNo.Cmp(batchOne.SeqNo()) != 0 {
		t.Errorf("batch listing was not returned correctly")
	}

	// rollup one contains batches 3 & 4
	batchListing1, err := GetRollupBatches(db, rollup2.Hash())
	if err != nil {
		t.Errorf("could not get rollup batches. Cause: %s", err)
	}

	// should be two elements
	if big.NewInt(int64(batchListing1.Total)).Cmp(big.NewInt(2)) != 0 {
		t.Errorf("batch listing was not calculated correctly")
	}
	// second element should be batch 3 as we're ordering by height descending
	if batchListing1.BatchesData[1].Header.SequencerOrderNo.Cmp(batchThree.SeqNo()) != 0 {
		t.Errorf("batch listing was not returned correctly")
	}
}

func createRollup(lastBatch int64) common.ExtRollup {
	header := common.RollupHeader{
		LastBatchSeqNo: uint64(lastBatch),
	}

	rollup := common.ExtRollup{
		Header: &header,
	}

	return rollup
}

func createRollupMetadata(firstBatch int64) common.PublicRollupMetadata {
	return common.PublicRollupMetadata{
		FirstBatchSequence: big.NewInt(firstBatch),
		StartTime:          uint64(time.Now().Unix()),
	}
}
