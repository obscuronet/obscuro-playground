package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

var (
	sharedSecret = []byte("SharedSecret")

	attestationKeyPrefix           = []byte("oAK")  // attestationKeyPrefix + address -> key
	syntheticTransactionsKeyPrefix = []byte("oSTX") // attestationKeyPrefix + address -> key

	rollupHeaderPrefix       = []byte("oh")  // rollupHeaderPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix         = []byte("on")  // rollupHeaderPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	rollupBodyPrefix         = []byte("ob")  // rollupBodyPrefix + num (uint64 big endian) + hash -> rollup body
	rollupHeaderNumberPrefix = []byte("oH")  // rollupHeaderNumberPrefix + hash -> num (uint64 big endian)
	headsAfterL1BlockPrefix  = []byte("och") // headsAfterL1BlockPrefix + hash -> num (uint64 big endian)
	logsPrefix               = []byte("olg") // logsPrefix + hash -> block logs
	rollupReceiptsPrefix     = []byte("or")  // rollupReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	contractReceiptPrefix    = []byte("ocr") // contractReceiptPrefix + address -> tx hash
	txLookupPrefix           = []byte("ol")  // txLookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix          = []byte("oB")  // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits
)

// encodeRollupNumber encodes a rollup number as big endian uint64
func encodeRollupNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = rollupHeaderPrefix + hash
func headerKey(hash common.Hash) []byte {
	return append(rollupHeaderPrefix, hash.Bytes()...)
}

// headerKeyPrefix = headerPrefix + num (uint64 big endian)
func headerKeyPrefix(number uint64) []byte {
	return append(rollupHeaderPrefix, encodeRollupNumber(number)...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(rollupHeaderNumberPrefix, hash.Bytes()...)
}

// rollupBodyKey = rollupBodyPrefix + hash
func rollupBodyKey(hash common.Hash) []byte {
	return append(rollupBodyPrefix, hash.Bytes()...)
}

// headsAfterL1BlockKey = headsAfterL1BlockPrefix + hash
func headsAfterL1BlockKey(hash common.Hash) []byte {
	return append(headsAfterL1BlockPrefix, hash.Bytes()...)
}

// logsKey = logsPrefix + hash
func logsKey(hash common.Hash) []byte {
	return append(logsPrefix, hash.Bytes()...)
}

// rollupReceiptsKey = rollupReceiptsPrefix + hash
func rollupReceiptsKey(hash common.Hash) []byte {
	return append(rollupReceiptsPrefix, hash.Bytes()...)
}

func contractReceiptKey(contractAddress common.Address) []byte {
	return append(contractReceiptPrefix, contractAddress.Bytes()...)
}

// txLookupKey = txLookupPrefix + hash
func txLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func bloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(rollupHeaderPrefix, encodeRollupNumber(number)...), headerHashSuffix...)
}

func attestationPkKey(aggregator common.Address) []byte {
	return append(attestationKeyPrefix, aggregator.Bytes()...)
}

func crossChainMessagesKey(blockHash common.Hash) []byte {
	return append(syntheticTransactionsKeyPrefix, blockHash.Bytes()...)
}
