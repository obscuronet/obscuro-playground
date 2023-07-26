package gethutil

import (
	"bytes"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// Utilities for working with geth structures

// EmptyHash is useful for comparisons to check if hash has been set
var EmptyHash = gethcommon.Hash{}

// LCA - returns the latest common ancestor of the 2 blocks or an error if no common ancestor is found
// it also returns the blocks that became canonincal, and the once that are now the fork
func LCA(newCanonical *types.Block, oldCanonical *types.Block, resolver storage.BlockResolver) (*types.Block, []common.L1BlockHash, []common.L1BlockHash, error) {
	b, cp, ncp, err := internalLCA(newCanonical, oldCanonical, resolver, []common.L1BlockHash{}, []common.L1BlockHash{oldCanonical.Hash()})
	// remove the common ancestor
	if len(cp) > 0 {
		cp = cp[0 : len(cp)-1]
	}
	if len(ncp) > 0 {
		ncp = ncp[0 : len(ncp)-1]
	}
	return b, cp, ncp, err
}

func internalLCA(newCanonical *types.Block, oldCanonical *types.Block, resolver storage.BlockResolver, canonicalPath []common.L1BlockHash, nonCanonicalPath []common.L1BlockHash) (*types.Block, []common.L1BlockHash, []common.L1BlockHash, error) {
	if newCanonical.NumberU64() == common.L1GenesisHeight || oldCanonical.NumberU64() == common.L1GenesisHeight {
		return newCanonical, canonicalPath, nonCanonicalPath, nil
	}
	if bytes.Equal(newCanonical.Hash().Bytes(), oldCanonical.Hash().Bytes()) {
		return newCanonical, canonicalPath, nonCanonicalPath, nil
	}
	if newCanonical.NumberU64() > oldCanonical.NumberU64() {
		p, err := resolver.FetchBlock(newCanonical.ParentHash())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
		}

		return internalLCA(p, oldCanonical, resolver, append(canonicalPath, p.Hash()), nonCanonicalPath)
	}
	if oldCanonical.NumberU64() > newCanonical.NumberU64() {
		p, err := resolver.FetchBlock(oldCanonical.ParentHash())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
		}

		return internalLCA(newCanonical, p, resolver, canonicalPath, append(nonCanonicalPath, p.Hash()))
	}
	parentBlockA, err := resolver.FetchBlock(newCanonical.ParentHash())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
	}
	parentBlockB, err := resolver.FetchBlock(oldCanonical.ParentHash())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
	}

	return internalLCA(parentBlockA, parentBlockB, resolver, append(canonicalPath, parentBlockA.Hash()), append(nonCanonicalPath, parentBlockB.Hash()))
}
