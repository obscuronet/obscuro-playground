package ethadapter

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	// excessBlobBits represents the number of bits in the last byte of a BLS scalar
	// that are not fully utilized due to the BLS modulus not being a power of 2.
	// This allows for more efficient packing of data into blobs by utilizing these
	// otherwise unused bits.
	excessBlobBits = 6
	MaxBlobBytes   = 32 * 4096
)

// MakeSidecar builds & returns the BlobTxSidecar and corresponding blob hashes from the raw blob
// data.
func MakeSidecar(blobs []*kzg4844.Blob) (*types.BlobTxSidecar, []gethcommon.Hash, error) {
	sidecar := &types.BlobTxSidecar{}
	var blobHashes []gethcommon.Hash
	for i, blob := range blobs {
		sidecar.Blobs = append(sidecar.Blobs, *blob)
		commitment, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Commitments = append(sidecar.Commitments, commitment)
		proof, err := kzg4844.ComputeBlobProof(blob, commitment)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", i, err)
		}
		sidecar.Proofs = append(sidecar.Proofs, proof)
		blobHashes = append(blobHashes, KZGToVersionedHash(commitment))
	}
	return sidecar, blobHashes, nil
}

// EncodeBlobs converts bytes into blobs used for KZG commitment EIP-4844
// transactions on Ethereum.
func EncodeBlobs(data []byte) ([]*kzg4844.Blob, error) {
	data, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}

	if len(data) >= MaxBlobBytes {
		return nil, fmt.Errorf("data too large to encode in blobs")
	}

	var blobs []*kzg4844.Blob
	for len(data) > 0 {
		var b kzg4844.Blob
		// fill blob with full bytes first
		data = fillBlobBytes(b[:], data)
		// fill the remaining bits
		data, err = fillBlobBits(b[:], data)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, &b)
	}
	return blobs, nil
}

func fillBlobBytes(blob []byte, data []byte) []byte {
	for fieldElement := 0; fieldElement < params.BlobTxFieldElementsPerBlob; fieldElement++ {
		startIdx := fieldElement*32 + 1
		copy(blob[startIdx:startIdx+31], data)
		if len(data) <= 31 {
			return nil
		}
		data = data[31:]
	}
	return data
}

func fillBlobBits(blob []byte, data []byte) ([]byte, error) {
	var acc uint16 // accumulator for bits
	accBits := 0   // number of bits currently in the accumulator

	for fieldElement := 0; fieldElement < params.BlobTxFieldElementsPerBlob; fieldElement++ {
		// if we need more bits and have more data, add a byte to the accumulator
		if accBits < excessBlobBits && len(data) > 0 {
			acc |= uint16(data[0]) << accBits
			accBits += 8
			data = data[1:]
		}

		// fill the excess bits of the current field element
		blob[fieldElement*32] = uint8(acc & ((1 << excessBlobBits) - 1))
		accBits -= excessBlobBits

		if accBits < 0 {
			break
		}

		// shift the used bits out of the accumulator
		acc >>= excessBlobBits
	}

	if accBits > 0 {
		return nil, fmt.Errorf("unexpected %v spare accBits remaining", accBits)
	}

	return data, nil
}

// DecodeBlobs decodes blobs into the data encoded in them accounting for excess blob bits
func DecodeBlobs(blobs []*kzg4844.Blob) ([]byte, error) {
	var rlpData []byte
	for _, blob := range blobs {
		for fieldIndex := 0; fieldIndex < params.BlobTxFieldElementsPerBlob; fieldIndex++ {
			rlpData = append(rlpData, blob[fieldIndex*32+1:(fieldIndex+1)*32]...)
		}
		var acc uint16
		cumulativeBits := 0
		for fieldIndex := 0; fieldIndex < params.BlobTxFieldElementsPerBlob; fieldIndex++ {
			acc |= uint16(blob[fieldIndex*32]) << cumulativeBits
			cumulativeBits += excessBlobBits
			if cumulativeBits >= 8 {
				rlpData = append(rlpData, uint8(acc))
				acc >>= 8
				cumulativeBits -= 8
			}
		}
		if cumulativeBits != 0 {
			return nil, fmt.Errorf("somehow ended up with %v spare cumulative bits", cumulativeBits)
		}
	}
	var outputData []byte
	err := rlp.Decode(bytes.NewReader(rlpData), &outputData)
	return outputData, err
}

// ReconstructRollup decodes and returns the ExtRollup in the blob
func ReconstructRollup(blobs []*kzg4844.Blob) (*common.ExtRollup, error) {
	data, err := DecodeBlobs(blobs)
	if err != nil {
		return nil, fmt.Errorf("could not decode rollup blob. Cause: %w ", err)
	}
	var rollup common.ExtRollup
	if err := rlp.DecodeBytes(data, &rollup); err != nil {
		return nil, fmt.Errorf("could not decode rollup bytes. Cause: %w", err)
	}

	return &rollup, nil
}

// ComputeBlobHash computes the versioned hash of a blob that matches the blobhash opcode
func ComputeBlobHash(blob *kzg4844.Blob) (gethcommon.Hash, error) {
	commitment, err := kzg4844.BlobToCommitment(blob)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	return KZGToVersionedHash(commitment), nil
}
