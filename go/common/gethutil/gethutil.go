package gethutil

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

// Utilities for working with geth structures

type BlockResolver interface {
	FetchBlock(ctx context.Context, blockHash common.L1BlockHash) (*types.Header, error)
}

// EmptyHash is useful for comparisons to check if hash has been set
var EmptyHash = gethcommon.Hash{}

// LCA - returns the latest common ancestor of the 2 blocks or an error if no common ancestor is found
// it also returns the blocks that became canonical, and the once that are now the fork
func LCA(ctx context.Context, newCanonical *types.Header, oldCanonical *types.Header, resolver BlockResolver) (*common.ChainFork, error) {
	b, cp, ncp, err := internalLCA(ctx, newCanonical, oldCanonical, resolver, []common.L1BlockHash{}, []common.L1BlockHash{})
	return &common.ChainFork{
		NewCanonical:     newCanonical,
		OldCanonical:     oldCanonical,
		CommonAncestor:   b,
		CanonicalPath:    cp,
		NonCanonicalPath: ncp,
	}, err
}

func internalLCA(ctx context.Context, newCanonical *types.Header, oldCanonical *types.Header, resolver BlockResolver, canonicalPath []common.L1BlockHash, nonCanonicalPath []common.L1BlockHash) (*types.Header, []common.L1BlockHash, []common.L1BlockHash, error) {
	if newCanonical.Number.Uint64() == common.L1GenesisHeight || oldCanonical.Number.Uint64() == common.L1GenesisHeight {
		return oldCanonical, canonicalPath, nonCanonicalPath, nil
	}
	if newCanonical.Hash() == oldCanonical.Hash() {
		// this is where we reach the common ancestor, which we add to the canonical path
		return newCanonical, append(canonicalPath, newCanonical.Hash()), nonCanonicalPath, nil
	}
	if newCanonical.Number.Uint64() > oldCanonical.Number.Uint64() {
		p, err := resolver.FetchBlock(ctx, newCanonical.ParentHash)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", newCanonical.ParentHash, err)
		}

		return internalLCA(ctx, p, oldCanonical, resolver, append(canonicalPath, newCanonical.Hash()), nonCanonicalPath)
	}
	if oldCanonical.Number.Uint64() > newCanonical.Number.Uint64() {
		p, err := resolver.FetchBlock(ctx, oldCanonical.ParentHash)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", oldCanonical.ParentHash, err)
		}

		return internalLCA(ctx, newCanonical, p, resolver, canonicalPath, append(nonCanonicalPath, oldCanonical.Hash()))
	}
	parentBlockA, err := resolver.FetchBlock(ctx, newCanonical.ParentHash)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", newCanonical.ParentHash, err)
	}
	parentBlockB, err := resolver.FetchBlock(ctx, oldCanonical.ParentHash)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not retrieve parent block %s. Cause: %w", oldCanonical.ParentHash, err)
	}

	return internalLCA(ctx, parentBlockA, parentBlockB, resolver, append(canonicalPath, newCanonical.Hash()), append(nonCanonicalPath, oldCanonical.Hash()))
}
