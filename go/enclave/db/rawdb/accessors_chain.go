package rawdb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

func ReadRollup(db ethdb.KeyValueReader, hash gethcommon.Hash) (*core.Rollup, error) {
	header, err := readHeader(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}

	body, err := ReadBody(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read body. Cause: %w", err)
	}

	return &core.Rollup{
		Header:       header,
		Transactions: body,
	}, nil
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash gethcommon.Hash) (*uint64, error) {
	data, err := db.Get(headerNumberKey(hash))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	if len(data) != 8 {
		return nil, fmt.Errorf("header number bytes had wrong length")
	}
	number := binary.BigEndian.Uint64(data)
	return &number, nil
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *core.Rollup) error {
	if err := writeHeader(db, rollup.Header); err != nil {
		return fmt.Errorf("could not write header. Cause: %w", err)
	}
	if err := writeBody(db, rollup.Hash(), rollup.Transactions); err != nil {
		return fmt.Errorf("could not write body. Cause: %w", err)
	}
	return nil
}

// Stores a rollup header into the database and also stores the hash-to-number mapping.
func writeHeader(db ethdb.KeyValueWriter, header *common.Header) error {
	hash := header.Hash()
	number := header.Number.Uint64()

	// Write the hash -> number mapping
	err := writeHeaderNumber(db, hash, number)
	if err != nil {
		return fmt.Errorf("could not write header number. Cause: %w", err)
	}

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		return fmt.Errorf("could not encode rollup header. Cause: %w", err)
	}
	key := headerKey(hash)
	if err = db.Put(key, data); err != nil {
		return fmt.Errorf("could not put header in DB. Cause: %w", err)
	}
	return nil
}

// Stores the hash->number mapping.
func writeHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) error {
	key := headerNumberKey(hash)
	enc := encodeRollupNumber(number)
	if err := db.Put(key, enc); err != nil {
		return fmt.Errorf("could not put header number in DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the rollup header corresponding to the hash.
func readHeader(db ethdb.KeyValueReader, hash gethcommon.Hash) (*common.Header, error) {
	data, err := readHeaderRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read header. Cause: %w", err)
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		return nil, fmt.Errorf("could not decode rollup header. Cause: %w", err)
	}
	return header, nil
}

// Retrieves a block header in its raw RLP database encoding.
func readHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash) (rlp.RawValue, error) {
	data, err := db.Get(headerKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block header. Cause: %w", err)
	}
	return data, nil
}

func writeBody(db ethdb.KeyValueWriter, hash gethcommon.Hash, body []*common.L2Tx) error {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		return fmt.Errorf("could not encode L2 transactions. Cause: %w", err)
	}
	if err = writeBodyRLP(db, hash, data); err != nil {
		return fmt.Errorf("could not write L2 transactions. Cause: %w", err)
	}
	return nil
}

// ReadBody retrieves the rollup body corresponding to the hash.
func ReadBody(db ethdb.KeyValueReader, hash gethcommon.Hash) ([]*common.L2Tx, error) {
	data, err := readBodyRLP(db, hash)
	if err != nil {
		return nil, fmt.Errorf("could not read body. Cause: %w", err)
	}
	body := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		return nil, fmt.Errorf("could not decode L2 transactions. Cause: %w", err)
	}
	return *body, nil
}

// Stores an RLP encoded block body into the database.
func writeBodyRLP(db ethdb.KeyValueWriter, hash gethcommon.Hash, rlp rlp.RawValue) error {
	if err := db.Put(rollupBodyKey(hash), rlp); err != nil {
		return fmt.Errorf("could not put rollup body into DB. Cause: %w", err)
	}
	return nil
}

// Retrieves the block body (transactions and uncles) in RLP encoding.
func readBodyRLP(db ethdb.KeyValueReader, hash gethcommon.Hash) (rlp.RawValue, error) {
	data, err := db.Get(rollupBodyKey(hash))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup body from DB. Cause: %w", err)
	}
	return data, nil
}

func ReadRollupsForHeight(db ethdb.Database, number uint64) ([]*core.Rollup, error) {
	hashes := readAllHashes(db, number)
	rollups := make([]*core.Rollup, len(hashes))
	for i, hash := range hashes {
		rollup, err := ReadRollup(db, hash)
		if err != nil {
			return nil, err
		}
		rollups[i] = rollup
	}
	return rollups, nil
}

// Retrieves all the hashes assigned to blocks at a certain heights, both canonical and reorged forks included.
func readAllHashes(db ethdb.Iteratee, number uint64) []gethcommon.Hash {
	prefix := headerKeyPrefix(number)

	hashes := make([]gethcommon.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix)+32 {
			hashes = append(hashes, gethcommon.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

func WriteL2Head(db ethdb.KeyValueWriter, l1Head common.L1RootHash, l2Head common.L2RootHash) error {
	if err := db.Put(headsAfterL1BlockKey(l1Head), l2Head.Bytes()); err != nil {
		return fmt.Errorf("could not put chain heads in DB. Cause: %w", err)
	}
	return nil
}

func ReadL2Head(kv ethdb.KeyValueReader, l1Head common.L1RootHash) (*common.L2RootHash, error) {
	data, err := kv.Get(headsAfterL1BlockKey(l1Head))
	if err != nil {
		return nil, errutil.ErrNotFound
	}
	l2Head := gethcommon.BytesToHash(data)
	return &l2Head, nil
}

func WriteBlockLogs(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, logs []*types.Log) error {
	// Geth serialises its logs in a reduced form to minimise storage space. For now, it is more straightforward for us
	// to serialise all the fields by converting the logs to this type.
	logsForStorage := make([]*logForStorage, len(logs))
	for idx, fullFatLog := range logs {
		logsForStorage[idx] = toLogForStorage(fullFatLog)
	}

	logBytes, err := rlp.EncodeToBytes(logsForStorage)
	if err != nil {
		return fmt.Errorf("could not encode logs. Cause: %w", err)
	}

	if err := db.Put(logsKey(blockHash), logBytes); err != nil {
		return fmt.Errorf("could not put logs in DB. Cause: %w", err)
	}
	return nil
}

func ReadBlockLogs(kv ethdb.KeyValueReader, blockHash gethcommon.Hash) ([]*types.Log, error) {
	data, err := kv.Get(logsKey(blockHash))
	if err != nil {
		return nil, err
	}

	logsForStorage := new([]*logForStorage)
	if err := rlp.Decode(bytes.NewReader(data), logsForStorage); err != nil {
		return nil, fmt.Errorf("could not decode logs. Cause: %w", err)
	}

	logs := make([]*types.Log, len(*logsForStorage))
	for idx, logToStore := range *logsForStorage {
		logs[idx] = logToStore.toLog()
	}

	return logs, nil
}

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func ReadCanonicalHash(db ethdb.Reader, number uint64) (*gethcommon.Hash, error) {
	// Get it by hash from leveldb
	data, err := db.Get(headerHashKey(number))
	if err != nil {
		return nil, err
	}
	hash := gethcommon.BytesToHash(data)
	return &hash, nil
}

// WriteCanonicalHash stores the hash assigned to a canonical rollup number.
func WriteCanonicalHash(db ethdb.KeyValueWriter, l2Head *core.Rollup) error {
	if err := db.Put(headerHashKey(l2Head.NumberU64()), l2Head.Hash().Bytes()); err != nil {
		return fmt.Errorf("failed to store number to hash mapping. Cause: %w", err)
	}
	return nil
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db ethdb.KeyValueWriter, number uint64) error {
	if err := db.Delete(headerHashKey(number)); err != nil {
		return fmt.Errorf("failed to delete number to hash mapping. Cause: %w", err)
	}
	return nil
}
