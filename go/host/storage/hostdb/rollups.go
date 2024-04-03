package hostdb

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	selectExtRollup    = "SELECT ext_rollup from rollup_host r"
	selectLatestRollup = "SELECT ext_rollup FROM rollup_host ORDER BY time_stamp DESC LIMIT 1"
)

// AddRollup adds a rollup to the DB
func AddRollup(dbtx *dbTransaction, rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error {
	extRollup, err := rlp.EncodeToBytes(rollup)
	if err != nil {
		return fmt.Errorf("could not encode rollup: %w", err)
	}
	_, err = dbtx.GetDB().Exec(dbtx.GetSQLStatements().InsertRollup,
		truncTo16(rollup.Header.Hash()),      // short hash
		metadata.FirstBatchSequence.Uint64(), // first batch sequence
		rollup.Header.LastBatchSeqNo,         // last batch sequence
		metadata.StartTime,                   // timestamp
		extRollup,                            // rollup blob
		block.Hash(),                         // l1 block hash
	)
	if err != nil {
		return fmt.Errorf("could not insert rollup. Cause: %w", err)
	}

	return nil
}

// GetRollupListing returns latest rollups given a pagination.
// For example, offset 1, size 10 will return the latest 11-20 rollups.
func GetRollupListing(dbtx *dbTransaction, pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	rows, err := dbtx.GetDB().Query(dbtx.GetSQLStatements().SelectRollups, pagination.Size, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rollups []common.PublicRollup

	for rows.Next() {
		var id, startSeq, endSeq, timeStamp int
		var hash, extRollup, compressionBlock []byte

		var rollup common.PublicRollup
		err = rows.Scan(&id, &hash, &startSeq, &endSeq, &timeStamp, &extRollup, &compressionBlock)
		if err != nil {
			return nil, err
		}

		extRollupDecoded := new(common.ExtRollup)
		if err := rlp.DecodeBytes(extRollup, extRollupDecoded); err != nil {
			return nil, fmt.Errorf("could not decode rollup header. Cause: %w", err)
		}

		rollup = common.PublicRollup{
			ID:        big.NewInt(int64(id)),
			Hash:      hash,
			FirstSeq:  big.NewInt(int64(startSeq)),
			LastSeq:   big.NewInt(int64(endSeq)),
			Timestamp: uint64(timeStamp),
			Header:    extRollupDecoded.Header,
			L1Hash:    compressionBlock,
		}
		rollups = append(rollups, rollup)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &common.RollupListingResponse{
		RollupsData: rollups,
		Total:       uint64(len(rollups)),
	}, nil
}

func GetExtRollup(dbtx *dbTransaction, hash gethcommon.Hash) (*common.ExtRollup, error) {
	whereQuery := " WHERE r.hash=" + dbtx.GetSQLStatements().Placeholder
	return fetchExtRollup(dbtx.GetDB(), whereQuery, truncTo16(hash))
}

// GetRollupHeader returns the rollup with the given hash.
func GetRollupHeader(dbtx *dbTransaction, hash gethcommon.Hash) (*common.RollupHeader, error) {
	whereQuery := " WHERE r.hash=" + dbtx.GetSQLStatements().Placeholder
	return fetchRollupHeader(dbtx.GetDB(), whereQuery, truncTo16(hash))
}

// GetRollupHeaderByBlock returns the rollup for the given block
func GetRollupHeaderByBlock(dbtx *dbTransaction, blockHash gethcommon.Hash) (*common.RollupHeader, error) {
	whereQuery := " WHERE r.compression_block=" + dbtx.GetSQLStatements().Placeholder
	return fetchRollupHeader(dbtx.GetDB(), whereQuery, blockHash)
}

// GetLatestRollup returns the latest rollup ordered by timestamp
func GetLatestRollup(db *sql.DB) (*common.RollupHeader, error) {
	extRollup, err := fetchHeadRollup(db)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch head rollup: %w", err)
	}
	return extRollup.Header, nil
}

func fetchRollupHeader(db *sql.DB, whereQuery string, args ...any) (*common.RollupHeader, error) {
	rollup, err := fetchExtRollup(db, whereQuery, args...)
	if err != nil {
		return nil, err
	}
	return rollup.Header, nil
}

func fetchExtRollup(db *sql.DB, whereQuery string, args ...any) (*common.ExtRollup, error) {
	var rollupBlob []byte
	query := selectExtRollup + whereQuery
	var err error
	if len(args) > 0 {
		err = db.QueryRow(query, args...).Scan(&rollupBlob)
	} else {
		err = db.QueryRow(query).Scan(&rollupBlob)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(rollupBlob, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}

func fetchHeadRollup(db *sql.DB) (*common.ExtRollup, error) {
	var extRollup []byte
	err := db.QueryRow(selectLatestRollup).Scan(&extRollup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch rollup by hash: %w", err)
	}
	var rollup common.ExtRollup
	err = rlp.DecodeBytes(extRollup, &rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rollup: %w", err)
	}

	return &rollup, nil
}
