package common

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside an enclave.
type ExtRollup struct {
	Header               *RollupHeader   // the fields required by the management contract
	CalldataRollupHeader []byte          // encrypted header useful for recreating the batches
	BatchPayloads        []byte          // The transactions included in the rollup, in external/encrypted form.
	BlobData             []*kzg4844.Blob // The blob data to be published
	hash                 atomic.Value
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollup) Hash() L2RollupHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RollupHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}
