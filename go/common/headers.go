package common

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
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

	// The custom Obscuro fields.
	L1Proof                       L1BlockHash                           `json:"l1Proof"` // the L1 block used by the enclave to generate the current batch
	R, S                          *big.Int                              // signature values
	CrossChainMessages            []MessageBus.StructsCrossChainMessage `json:"crossChainMessages"`
	LatestInboundCrossChainHash   common.Hash                           `json:"inboundCrossChainHash"`   // The block hash of the latest block that has been scanned for cross chain messages.
	LatestInboundCrossChainHeight *big.Int                              `json:"inboundCrossChainHeight"` // The block height of the latest block that has been scanned for cross chain messages.
	TransfersTree                 common.Hash                           `json:"transfersTree"`
}

// MarshalJSON custom marshals the BatchHeader into a json
func (b *BatchHeader) MarshalJSON() ([]byte, error) {
	type Alias BatchHeader
	return json.Marshal(struct {
		*Alias
		Hash       common.Hash       `json:"hash"`
		UncleHash  *common.Hash      `json:"sha3Uncles"`
		Coinbase   *common.Address   `json:"miner"`
		Bloom      *types.Bloom      `json:"logsBloom"`
		Difficulty *big.Int          `json:"difficulty"`
		Nonce      *types.BlockNonce `json:"nonce"`

		// BaseFee was added by EIP-1559 and is ignored in legacy headers.
		BaseFee *big.Int `json:"baseFeePerGas"`
	}{
		(*Alias)(b),
		b.Hash(),
		nil,
		&b.Coinbase,
		nil,
		nil,
		nil,
		b.BaseFee,
	})
}

// RollupHeader is a public / plaintext struct that holds common properties of rollups.
// All these fields are processed by the Management contract
type RollupHeader struct {
	L1Proof       L1BlockHash // the L1 block hash used by the enclave to generate the current rollup
	L1ProofNumber *big.Int    // the height of the proof - used by the management contract to check
	Coinbase      common.Address

	CrossChainMessages []MessageBus.StructsCrossChainMessage `json:"crossChainMessages"`

	PayloadHash common.Hash // The hash of the compressed batches. TODO
	R, S        *big.Int    // signature values

	LastBatchSeqNo uint64
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
	cp.R = nil
	cp.S = nil
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
	cp.R = nil
	cp.S = nil
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
