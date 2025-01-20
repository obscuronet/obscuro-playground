package common

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// Used to hash headers.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// BatchHeader is a public / plaintext struct that holds common properties of batches.
// Making changes to this struct will require GRPC + GRPC Converters regen
type BatchHeader struct {
	ParentHash       L2BatchHash    `json:"parentHash"`
	Root             StateRoot      `json:"stateRoot"`
	TxHash           common.Hash    `json:"transactionsRoot"`
	ReceiptHash      common.Hash    `json:"receiptsRoot"`
	Number           *big.Int       `json:"number"`           // height of the batch
	SequencerOrderNo *big.Int       `json:"sequencerOrderNo"` // multiple batches can be created with the same height in case of L1 reorgs. The sequencer is responsible for including all of them in the rollups.
	GasLimit         uint64         `json:"gasLimit"`
	GasUsed          uint64         `json:"gasUsed"`
	Time             uint64         `json:"timestamp"`
	Extra            []byte         `json:"extraData"`
	BaseFee          *big.Int       `json:"baseFee"`
	Coinbase         common.Address `json:"coinbase"`

	// The custom TEN fields.
	L1Proof        L1BlockHash              `json:"l1Proof"` // the L1 block used by the enclave to generate the current batch
	Signature      []byte                   `json:"signature"`
	CrossChainRoot common.Hash              `json:"crossChainTreeHash"` // This is the root hash of a merkle tree, built from all the cross chain messages and transfers that need to go on MainNet.
	CrossChainTree SerializedCrossChainTree `json:"crossChainTree"`     // Those are the leafs of the merkle tree hashed for privacy. Necessary for clients to be able to build proofs as they have no access to all transactions in a batch or their receipts.
}

// TODO - use exposed headers once #3987 is completed.
func ConvertBatchHeaderToHeader(batchHeader *BatchHeader) *types.Header {
	var header types.Header

	// Convert ParentHash if necessary
	header.ParentHash = batchHeader.ParentHash

	// Convert Root if necessary
	header.Root = batchHeader.Root

	// Directly assign fields that have the same type
	header.TxHash = batchHeader.TxHash
	header.ReceiptHash = batchHeader.ReceiptHash
	header.Number = batchHeader.Number
	header.GasLimit = batchHeader.GasLimit
	header.GasUsed = batchHeader.GasUsed
	header.Time = batchHeader.Time
	header.Extra = batchHeader.Extra
	header.BaseFee = batchHeader.BaseFee
	header.Coinbase = batchHeader.Coinbase

	// Set default values for fields not present in BatchHeader
	header.UncleHash = common.Hash{}  // Default to zero hash
	header.Bloom = types.Bloom{}      // Default to zero bloom filter
	header.Difficulty = big.NewInt(0) // Default difficulty
	header.MixDigest = common.Hash{}  // Default mix digest
	header.Nonce = types.BlockNonce{} // Default nonce

	// Set optional fields to nil or default values
	header.WithdrawalsHash = nil
	header.BlobGasUsed = nil
	header.ExcessBlobGas = nil
	header.ParentBeaconRoot = nil

	return &header
}

// IsGenesis indicates whether the batch is the genesis batch.
// todo (#718) - Change this to a check against a hardcoded genesis hash.
func (b *BatchHeader) IsGenesis() bool {
	return b.Number.Cmp(big.NewInt(int64(L2GenesisHeight))) == 0
}

type batchHeaderEncoding struct {
	Hash             common.Hash     `json:"hash"`
	ParentHash       L2BatchHash     `json:"parentHash"`
	Root             common.Hash     `json:"stateRoot"`
	TxHash           common.Hash     `json:"transactionsRoot"`
	ReceiptHash      common.Hash     `json:"receiptsRoot"`
	Number           *hexutil.Big    `json:"number"`
	SequencerOrderNo *hexutil.Big    `json:"sequencerOrderNo"`
	GasLimit         hexutil.Uint64  `json:"gasLimit"`
	GasUsed          hexutil.Uint64  `json:"gasUsed"`
	Time             hexutil.Uint64  `json:"timestamp"`
	Extra            []byte          `json:"extraData"`
	BaseFee          *hexutil.Big    `json:"baseFeePerGas"`
	Coinbase         *common.Address `json:"miner"`

	// The custom Obscuro fields.
	L1Proof            L1BlockHash              `json:"l1Proof"` // the L1 block used by the enclave to generate the current batch
	Signature          []byte                   `json:"signature"`
	CrossChainRootHash common.Hash              `json:"crossChainTreeHash"`
	CrossChainTree     SerializedCrossChainTree `json:"crossChainTree"`
}

// MarshalJSON custom marshals the BatchHeader into a json
func (b *BatchHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(batchHeaderEncoding{
		b.Hash(),
		b.ParentHash,
		b.Root,
		b.TxHash,
		b.ReceiptHash,
		(*hexutil.Big)(b.Number),
		(*hexutil.Big)(b.SequencerOrderNo),
		hexutil.Uint64(b.GasLimit),
		hexutil.Uint64(b.GasUsed),
		hexutil.Uint64(b.Time),
		b.Extra,
		(*hexutil.Big)(b.BaseFee),
		&b.Coinbase,
		b.L1Proof,
		b.Signature,
		b.CrossChainRoot,
		b.CrossChainTree,
	})
}

func (b *BatchHeader) UnmarshalJSON(data []byte) error {
	dec := new(batchHeaderEncoding)
	err := json.Unmarshal(data, dec)
	if err != nil {
		return err
	}

	b.ParentHash = dec.ParentHash
	b.Root = dec.Root
	b.TxHash = dec.TxHash
	b.ReceiptHash = dec.ReceiptHash
	b.Number = (*big.Int)(dec.Number)
	b.SequencerOrderNo = (*big.Int)(dec.SequencerOrderNo)
	b.GasLimit = uint64(dec.GasLimit)
	b.GasUsed = uint64(dec.GasUsed)
	b.Time = uint64(dec.Time)
	b.Extra = dec.Extra
	b.BaseFee = (*big.Int)(dec.BaseFee)
	b.Coinbase = *dec.Coinbase
	b.L1Proof = dec.L1Proof
	b.Signature = dec.Signature
	b.CrossChainRoot = dec.CrossChainRootHash
	b.CrossChainTree = dec.CrossChainTree
	return nil
}

// RollupHeader is a public / plaintext struct that holds common properties of rollups.
// All these fields are processed by the Management contract
type RollupHeader struct {
	BlockHash      L1BlockHash // the l1 block that the sequencer considers canonical at the time when this rollup is created
	PayloadHash    common.Hash // The hash of the compressed batches. TODO
	Signature      []byte      // The signature of the sequencer enclave on the payload hash
	LastBatchSeqNo uint64

	// New fields needed for blob-based rollups
	BlobHash    common.Hash // Hash of the blob data, computed from the blob
	MessageRoot common.Hash // Root hash of cross-chain messages (can be 0 for now)
	BlockNumber uint64      // L1 block number for binding
}

// CalldataRollupHeader contains all information necessary to reconstruct the batches included in the rollup.
// This data structure is serialised, compressed, and encrypted, before being serialised again in the rollup.
type CalldataRollupHeader struct {
	FirstBatchSequence    *big.Int
	FirstCanonBatchHeight *big.Int
	FirstCanonParentHash  L2BatchHash

	Coinbase common.Address
	BaseFee  *big.Int
	GasLimit uint64

	StartTime       uint64
	BatchTimeDeltas [][]byte // todo - minimize assuming a default of 1 sec and then store only exceptions

	L1HeightDeltas [][]byte // delta of the block height. Stored as a byte array because rlp can't encode negative numbers

	// these fields are for debugging the compression. Uncomment if there are issues
	// BatchHashes  []L2BatchHash
	// BatchHeaders []*BatchHeader

	ReOrgs [][]byte `rlp:"optional"` // sparse list of reorged headers - non null only for reorgs.
}

// PublicRollupMetadata contains internal rollup data that can be requested from the enclave.
type PublicRollupMetadata struct {
	FirstBatchSequence *big.Int
	StartTime          uint64
}

// MarshalJSON custom marshals the RollupHeader into a json
func (r *RollupHeader) MarshalJSON() ([]byte, error) {
	type Alias RollupHeader
	return json.Marshal(struct {
		*Alias
		Hash common.Hash `json:"hash"`
	}{
		(*Alias)(r),
		r.Hash(),
	})
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (b *BatchHeader) Hash() L2BatchHash {
	cp := *b
	cp.Signature = nil
	hash, err := rlpHash(cp)
	if err != nil {
		panic("err hashing batch header")
	}
	return hash
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding excluding the signature.
func (r *RollupHeader) Hash() L2RollupHash {
	cp := *r
	cp.Signature = nil
	hash, err := rlpHash(cp)
	if err != nil {
		panic("err hashing rollup header")
	}
	return hash
}

// Encodes value, hashes the encoded bytes and returns the hash.
func rlpHash(value interface{}) (common.Hash, error) {
	var hash common.Hash

	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()

	err := rlp.Encode(sha, value)
	if err != nil {
		return hash, fmt.Errorf("unable to encode Value. %w", err)
	}

	_, err = sha.Read(hash[:])
	if err != nil {
		return hash, fmt.Errorf("unable to read encoded value. %w", err)
	}

	return hash, nil
}
