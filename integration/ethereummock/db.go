package ethereummock

import (
	"bytes"
	"context"
	"math/big"
	"sync"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

// HeightCommittedBlocks is the number of blocks deep a transaction must be to be considered safe from reorganisations.
const HeightCommittedBlocks = 15

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[common.L1BlockHash]*types.Block
	m          sync.RWMutex
}

func (n *blockResolverInMem) IsBlockCanonical(ctx context.Context, blockHash common.L1BlockHash) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (n *blockResolverInMem) FetchCanonicaBlockByHeight(_ context.Context, _ *big.Int) (*types.Block, error) {
	panic("implement me")
}

func (n *blockResolverInMem) Proof(_ context.Context, _ *core.Rollup) (*types.Block, error) {
	panic("implement me")
}

func NewResolver() *blockResolverInMem {
	return &blockResolverInMem{
		blockCache: map[common.L1BlockHash]*types.Block{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) StoreBlock(_ context.Context, block *types.Block, _ *common.ChainFork) error {
	n.m.Lock()
	defer n.m.Unlock()
	n.blockCache[block.Hash()] = block
	return nil
}

func (n *blockResolverInMem) FetchFullBlock(_ context.Context, hash common.L1BlockHash) (*types.Block, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	if !f {
		return nil, errutil.ErrNotFound
	}
	return block, nil
}

func (n *blockResolverInMem) FetchBlock(_ context.Context, hash common.L1BlockHash) (*types.Header, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	if !f {
		return nil, errutil.ErrNotFound
	}
	return block.Header(), nil
}

func (n *blockResolverInMem) FetchHeadBlock(_ context.Context) (*types.Block, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	var max *types.Block
	for k := range n.blockCache {
		bh := n.blockCache[k]
		if max == nil || max.NumberU64() < bh.NumberU64() {
			max = bh
		}
	}
	if max == nil {
		return nil, errutil.ErrNotFound
	}
	return max, nil
}

func (n *blockResolverInMem) ParentBlock(ctx context.Context, b *types.Header) (*types.Block, error) {
	return n.FetchFullBlock(ctx, b.ParentHash)
}

func (n *blockResolverInMem) IsAncestor(ctx context.Context, block *types.Header, maybeAncestor *types.Header) bool {
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.Number.Uint64() >= block.Number.Uint64() {
		return false
	}

	p, err := n.ParentBlock(ctx, block)
	if err != nil {
		return false
	}

	return n.IsAncestor(ctx, p.Header(), maybeAncestor)
}

func (n *blockResolverInMem) IsBlockAncestor(ctx context.Context, block *types.Header, maybeAncestor common.L1BlockHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) {
		return true
	}

	if bytes.Equal(maybeAncestor.Bytes(), MockGenesisBlock.Hash().Bytes()) {
		return true
	}

	if block.Number.Uint64() == common.L1GenesisHeight {
		return false
	}

	resolvedBlock, err := n.FetchFullBlock(ctx, maybeAncestor)
	if err == nil {
		if resolvedBlock.NumberU64() >= block.Number.Uint64() {
			return false
		}
	}

	p, err := n.ParentBlock(ctx, block)
	if err != nil {
		// todo (@tudor) - if error is not `errutil.ErrNotFound`, throw
		return false
	}

	return n.IsBlockAncestor(ctx, p.Header(), maybeAncestor)
}

// The cache of included transactions
type txDBInMem struct {
	transactionsPerBlockCache map[common.L1BlockHash]map[common.TxHash]*types.Transaction
	rpbcM                     *sync.RWMutex
}

func NewTxDB() TxDB {
	return &txDBInMem{
		transactionsPerBlockCache: make(map[common.L1BlockHash]map[common.TxHash]*types.Transaction),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDBInMem) Txs(b *types.Block) (map[common.TxHash]*types.Transaction, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()

	return val, found
}

func (n *txDBInMem) AddTxs(b *types.Block, newMap map[common.TxHash]*types.Transaction) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func (m *Node) removeCommittedTransactions(
	ctx context.Context,
	cb *types.Block,
	mempool []*types.Transaction,
	resolver *blockResolverInMem,
	db TxDB,
) []*types.Transaction {
	if cb.NumberU64() <= HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == HeightCommittedBlocks {
			break
		}

		p, err := resolver.FetchFullBlock(ctx, b.ParentHash())
		if err != nil {
			m.logger.Crit("Could not retrieve parent block.", log.ErrKey, err)
		}

		b = p
		i++
	}

	val, _ := db.Txs(b)
	//if !found {
	//	panic("should not fail here")
	//}

	return removeExisting(mempool, val)
}
