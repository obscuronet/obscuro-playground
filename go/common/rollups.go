package common

import (
	"sync/atomic"
)

// ExtRollup is an encrypted form of rollup used when passing the rollup around outside an enclave.
type ExtRollup struct {
	Header *RollupHeader
	// TODO - #718 - Consider compressing these batches before submitting to the L1.
	Batches []*ExtBatch // The batches included in the rollup, in external/encrypted form.
	hash    atomic.Value
}

// Hash returns the keccak256 hash of the rollup's header.
// The hash is computed on the first call and cached thereafter.
func (r *ExtRollup) Hash() L2RootHash {
	if hash := r.hash.Load(); hash != nil {
		return hash.(L2RootHash)
	}
	v := r.Header.Hash()
	r.hash.Store(v)
	return v
}

// ExtRollupFromExtBatches converts a set of ExtBatch into an ExtRollup. The head batch is the one stored in the
// ExtRollup header's `HeadBatchHash` field.
func ExtRollupFromExtBatches(headBatch *ExtBatch, additionalBatches []*ExtBatch) *ExtRollup {
	return &ExtRollup{
		Header:  headBatch.Header.ToRollupHeader(),
		Batches: append(additionalBatches, headBatch),
	}
}
