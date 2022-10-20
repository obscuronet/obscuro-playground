package rawdb

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/status-im/keycard-go/hexutils"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

func ReadRollup(db ethdb.KeyValueReader, hash gethcommon.Hash) *core.Rollup {
	height := ReadHeaderNumber(db, hash)
	if height == nil {
		return nil
	}
	return &core.Rollup{
		Header:       ReadHeader(db, hash, *height),
		Transactions: ReadBody(db, hash, *height),
	}
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash gethcommon.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

func WriteRollup(db ethdb.KeyValueWriter, rollup *core.Rollup) {
	WriteHeader(db, rollup.Header)
	WriteBody(db, rollup.Hash(), rollup.Header.Number.Uint64(), rollup.Transactions)
}

// WriteHeader stores a rollup header into the database and also stores the hash-
// to-number mapping.
func WriteHeader(db ethdb.KeyValueWriter, header *common.Header) {
	var (
		hash   = header.Hash()
		number = header.Number.Uint64()
	)
	// Write the hash -> number mapping
	WriteHeaderNumber(db, hash, number)

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Panic("could not encode rollup header. Cause: %s", err)
	}
	key := headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		log.Panic("could not put header in DB. Cause: %s", err)
	}
}

// WriteHeaderNumber stores the hash->number mapping.
func WriteHeaderNumber(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) {
	key := headerNumberKey(hash)
	enc := encodeRollupNumber(number)
	if err := db.Put(key, enc); err != nil {
		log.Panic("could not put header number in DB. Cause: %s", err)
	}
}

// ReadHeader retrieves the rollup header corresponding to the hash.
func ReadHeader(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64) *common.Header {
	data := ReadHeaderRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	header := new(common.Header)
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Panic("could not decode rollup header. Cause: %s", err)
	}
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func ReadHeaderRLP(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(headerKey(number, hash))
	if err != nil {
		log.Panic("could not retrieve block header. Cause: %s", err)
	}
	return data
}

func WriteBody(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64, body []*common.L2Tx) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		log.Panic("could not encode L2 transactions. Cause: %s", err)
	}
	WriteBodyRLP(db, hash, number, data)
}

// ReadBody retrieves the rollup body corresponding to the hash.
func ReadBody(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64) []*common.L2Tx {
	data := ReadBodyRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new([]*common.L2Tx)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		log.Panic("could not decode L2 transactions. Cause: %s", err)
	}
	return *body
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64, rlp rlp.RawValue) {
	if err := db.Put(rollupBodyKey(number, hash), rlp); err != nil {
		log.Panic("could not put rollup body into DB. Cause: %s", err)
	}
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyRLP(db ethdb.KeyValueReader, hash gethcommon.Hash, number uint64) rlp.RawValue {
	data, err := db.Get(rollupBodyKey(number, hash))
	if err != nil {
		log.Panic("could not retrieve rollup body :r_%d from DB. Cause: %s. Key: %s", common.ShortHash(hash), err, hexutils.BytesToHex(rollupBodyKey(number, hash)))
	}
	return data
}

func ReadRollupsForHeight(db ethdb.Database, number uint64) []*core.Rollup {
	hashes := ReadAllHashes(db, number)
	rollups := make([]*core.Rollup, len(hashes))
	for i, hash := range hashes {
		rollups[i] = ReadRollup(db, hash)
	}
	return rollups
}

// ReadAllHashes retrieves all the hashes assigned to blocks at a certain heights,
// both canonical and reorged forks included.
func ReadAllHashes(db ethdb.Iteratee, number uint64) []gethcommon.Hash {
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

func WriteBlockState(db ethdb.KeyValueWriter, bs *core.BlockState) {
	blockStateBytes, err := rlp.EncodeToBytes(bs)
	if err != nil {
		log.Panic("could not encode block state. Cause: %s", err)
	}
	if err := db.Put(blockStateKey(bs.Block), blockStateBytes); err != nil {
		log.Panic("could not put block state in DB. Cause: %s", err)
	}
}

func ReadBlockState(kv ethdb.KeyValueReader, hash gethcommon.Hash) *core.BlockState {
	data, _ := kv.Get(blockStateKey(hash))
	if data == nil {
		return nil
	}
	bs := new(core.BlockState)
	if err := rlp.Decode(bytes.NewReader(data), bs); err != nil {
		log.Panic("could not decode block state. Cause: %s", err)
	}
	return bs
}

func WriteBlockLogs(db ethdb.KeyValueWriter, blockHash gethcommon.Hash, logs []*types.Log) {
	// Geth serialises its logs in a reduced form to minimise storage space. For now, it is more straightforward for us
	// to serialise all the fields by converting the logs to this type.
	logsForStorage := []*logForStorage{}
	for _, fullFatLog := range logs {
		logsForStorage = append(logsForStorage, toLogForStorage(fullFatLog))
	}

	logBytes, err := rlp.EncodeToBytes(logsForStorage)
	if err != nil {
		log.Panic("could not encode logs. Cause: %s", err)
	}

	if err := db.Put(logsKey(blockHash), logBytes); err != nil {
		log.Panic("could not put logs in DB. Cause: %s", err)
	}
}

func ReadBlockLogs(kv ethdb.KeyValueReader, blockHash gethcommon.Hash) []*types.Log {
	data, _ := kv.Get(logsKey(blockHash))
	if data == nil {
		return nil
	}

	logsForStorage := new([]*logForStorage)
	if err := rlp.Decode(bytes.NewReader(data), logsForStorage); err != nil {
		log.Panic("could not decode logs. Cause: %s", err)
	}

	logs := make([]*types.Log, len(*logsForStorage))
	for idx, logToStore := range *logsForStorage {
		logs[idx] = logToStore.toLog()
	}

	return logs
}

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func ReadCanonicalHash(db ethdb.Reader, number uint64) gethcommon.Hash {
	// Get it by hash from leveldb
	data, _ := db.Get(headerHashKey(number))
	return gethcommon.BytesToHash(data)
}

// WriteCanonicalHash stores the hash assigned to a canonical block number.
func WriteCanonicalHash(db ethdb.KeyValueWriter, hash gethcommon.Hash, number uint64) {
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.Panic("Failed to store number to hash mapping. Cause: %s", err)
	}
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db ethdb.KeyValueWriter, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.Panic("Failed to delete number to hash mapping. Cause: %s", err)
	}
}

// ReadHeadRollupHash retrieves the hash of the current canonical head block.
func ReadHeadRollupHash(db ethdb.KeyValueReader) gethcommon.Hash {
	data, _ := db.Get(headRollupKey)
	if len(data) == 0 {
		return gethcommon.Hash{}
	}
	return gethcommon.BytesToHash(data)
}

// WriteHeadRollupHash stores the head block's hash.
func WriteHeadRollupHash(db ethdb.KeyValueWriter, hash gethcommon.Hash) {
	if err := db.Put(headRollupKey, hash.Bytes()); err != nil {
		log.Panic("Failed to store last block's hash. Cause: %s", err)
	}
}

// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
func ReadHeadHeaderHash(db ethdb.KeyValueReader) gethcommon.Hash {
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return gethcommon.Hash{}
	}
	return gethcommon.BytesToHash(data)
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
func WriteHeadHeaderHash(db ethdb.KeyValueWriter, hash gethcommon.Hash) {
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		log.Panic("Failed to store last header's hash. Cause: %s", err)
	}
}
