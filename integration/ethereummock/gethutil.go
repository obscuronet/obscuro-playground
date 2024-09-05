package ethereummock

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

// Utilities for working with geth structures

// EmptyHash is useful for comparisons to check if hash has been set
var EmptyHash = gethcommon.Hash{}

// LCA - returns the latest common ancestor of the 2 blocks or an error if no common ancestor is found
// it also returns the blocks that became canonical, and the once that are now the fork
func LCA(ctx context.Context, newCanonical *types.Block, oldCanonical *types.Block, resolver *blockResolverInMem) (*common.ChainFork, error) {
	b, cp, ncp, err := internalLCA(ctx, newCanonical, oldCanonical, resolver, []common.L1BlockHash{}, []common.L1BlockHash{})
	return &common.ChainFork{
		NewCanonical:     newCanonical.Header(),
		OldCanonical:     oldCanonical.Header(),
		CommonAncestor:   b.Header(),
		CanonicalPath:    cp,
		NonCanonicalPath: ncp,
	}, err
}

func internalLCA(ctx context.Context, newCanonical *types.Block, oldCanonical *types.Block, resolver *blockResolverInMem, canonicalPath []common.L1BlockHash, nonCanonicalPath []common.L1BlockHash) (*types.Block, []common.L1BlockHash, []common.L1BlockHash, error) {
	if newCanonical.NumberU64() == common.L1GenesisHeight || oldCanonical.NumberU64() == common.L1GenesisHeight {
		return newCanonical, canonicalPath, nonCanonicalPath, nil
	}
	if newCanonical.Hash() == oldCanonical.Hash() {
		// this is where we reach the common ancestor, which we add to the canonical path
		return newCanonical, append(canonicalPath, newCanonical.Hash()), nonCanonicalPath, nil
	}
	if newCanonical.NumberU64() > oldCanonical.NumberU64() {
		p, err := resolver.FetchBlock(ctx, newCanonical.ParentHash())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", newCanonical.ParentHash, err)
		}

		return internalLCA(ctx, p, oldCanonical, resolver, append(canonicalPath, newCanonical.Hash()), nonCanonicalPath)
	}
	if oldCanonical.NumberU64() > newCanonical.NumberU64() {
		p, err := resolver.FetchBlock(ctx, oldCanonical.ParentHash())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", oldCanonical.ParentHash, err)
		}

		return internalLCA(ctx, newCanonical, p, resolver, canonicalPath, append(nonCanonicalPath, oldCanonical.Hash()))
	}
	parentBlockA, err := resolver.FetchBlock(ctx, newCanonical.ParentHash())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", newCanonical.ParentHash, err)
	}
	parentBlockB, err := resolver.FetchBlock(ctx, oldCanonical.ParentHash())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", oldCanonical.ParentHash, err)
	}

	return internalLCA(ctx, parentBlockA, parentBlockB, resolver, append(canonicalPath, newCanonical.Hash()), append(nonCanonicalPath, oldCanonical.Hash()))
}
