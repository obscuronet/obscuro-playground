package hostdb

import (
	"github.com/ten-protocol/go-ten/go/common"
	"math/big"
	"testing"
	"time"
)

func TestCanStoreAndRetrieveRollup(t *testing.T) {
	db, err := createSQLiteDB(t)
	if err != nil {
		t.Fatalf("unable to initialise test db: %s", err)
	}

	metadata := createRollupMetadata(batchNumber - 10)
	rollup := createRollup(batchNumber)
	block := common.L1Block{}

	err = AddRollupHeader(db.DB, &rollup, &metadata, &block)
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	extRollup, err := GetExtRollup(db.DB, rollup.Header.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve ext rollup. Cause: %s", err)
	}

	rollupHeader, err := GetRollupHeader(db.DB, rollup.Header.Hash())
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
	block := common.L1Block{}

	err = AddRollupHeader(db.DB, &rollup, &metadata, &block)
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	rollupHeader, err := GetRollupHeaderByBlock(db.DB, block.Hash())
	if err != nil {
		t.Errorf("stored rollup but could not retrieve header. Cause: %s", err)
	}
	if big.NewInt(int64(rollupHeader.LastBatchSeqNo)).Cmp(big.NewInt(batchNumber)) != 0 {
		t.Errorf("rollup header was not stored correctly")
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
	block := common.L1Block{}

	err = AddRollupHeader(db.DB, &rollup1, &metadata1, &block)
	if err != nil {
		t.Errorf("could not store rollup. Cause: %s", err)
	}

	rollup2FirstSeq := int64(batchNumber + 1)
	rollup2LastSeq := int64(batchNumber + 10)
	metadata2 := createRollupMetadata(rollup2FirstSeq)
	rollup2 := createRollup(rollup2LastSeq)

	err = AddRollupHeader(db.DB, &rollup2, &metadata2, &block)
	if err != nil {
		t.Errorf("could not store rollup 2. Cause: %s", err)
	}

	rollup3FirstSeq := int64(batchNumber + 11)
	rollup3LastSeq := int64(batchNumber + 20)
	metadata3 := createRollupMetadata(rollup3FirstSeq)
	rollup3 := createRollup(rollup3LastSeq)
	err = AddRollupHeader(db.DB, &rollup3, &metadata3, &block)
	if err != nil {
		t.Errorf("could not store rollup 3. Cause: %s", err)
	}

	// page 1, size 2
	rollupListing, err := GetRollupListing(db.DB, &common.QueryPagination{Offset: 1, Size: 2})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be two elements
	if big.NewInt(int64(rollupListing.Total)).Cmp(big.NewInt(2)) != 0 {
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
	rollupListing1, err := GetRollupListing(db.DB, &common.QueryPagination{Offset: 0, Size: 3})
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
	rollupListing2, err := GetRollupListing(db.DB, &common.QueryPagination{Offset: 0, Size: 4})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be 3 elements
	if big.NewInt(int64(rollupListing2.Total)).Cmp(big.NewInt(3)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
	}

	// page 5, size 1
	rollupListing3, err := GetRollupListing(db.DB, &common.QueryPagination{Offset: 5, Size: 1})
	if err != nil {
		t.Errorf("could not get rollup listing. Cause: %s", err)
	}

	// should be 0 elements
	if big.NewInt(int64(rollupListing3.Total)).Cmp(big.NewInt(0)) != 0 {
		t.Errorf("rollup listing was not paginated correctly")
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
