package rollupchain

import (
	"bytes"

	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// Returns the new head rollup, and a boolean indicating whether this is a new rollup or the existing head rollup.
// todo - add statistics to determine why there are conflicts.
func selectNextRollup(parentRollup *core.Rollup, rollups []*core.Rollup, blockResolver db.BlockResolver) (*core.Rollup, bool) {
	var nextRollup *core.Rollup

	// We iterate over the proposed rollups to select the best next rollup.
	for _, rollup := range rollups {
		// We ignore rollups from L2 forks, or that are older than the parent rollup.
		isFromFork := !bytes.Equal(rollup.Header.ParentHash.Bytes(), parentRollup.Hash().Bytes())
		isOlderThanParent := rollup.Header.Number.Int64() <= parentRollup.Header.Number.Int64()
		if isFromFork || isOlderThanParent {
			continue
		}

		// If this is the first rollup to pass the checks above, or it is newer than the existing candidate, we make it
		// the candidate next rollup.
		if nextRollup == nil || blockResolver.ProofHeight(rollup) > blockResolver.ProofHeight(nextRollup) {
			nextRollup = rollup
		}
	}

	if nextRollup == nil {
		// We remain with the existing head rollup.
		return parentRollup, false
	}
	return nextRollup, true
}
