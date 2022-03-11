package ethereummock

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// LCA - returns the least nodecommon ancestor of the 2 blocks
func LCA(blockA *types.Block, blockB *types.Block, resolver obscurocommon.BlockResolver) *types.Block {
	if resolver.HeightBlock(blockA) == obscurocommon.L1GenesisHeight || resolver.HeightBlock(blockB) == obscurocommon.L1GenesisHeight {
		return blockA
	}
	if blockA.Hash() == blockB.Hash() {
		return blockA
	}
	if resolver.HeightBlock(blockA) > resolver.HeightBlock(blockB) {
		p, f := resolver.ParentBlock(blockA)
		if !f {
			panic("wtf")
		}
		return LCA(p, blockB, resolver)
	}
	if resolver.HeightBlock(blockB) > resolver.HeightBlock(blockA) {
		p, f := resolver.ParentBlock(blockB)
		if !f {
			panic("wtf")
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, f := resolver.ParentBlock(blockA)
	if !f {
		panic("wtf")
	}
	parentBlockB, f := resolver.ParentBlock(blockB)
	if !f {
		panic("wtf")
	}

	return LCA(parentBlockA, parentBlockB, resolver)
}

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findNotIncludedTxs(head *types.Block, txs []*obscurocommon.L1Tx, r obscurocommon.BlockResolver, db TxDB) []*obscurocommon.L1Tx {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *types.Block, r obscurocommon.BlockResolver, db TxDB) map[obscurocommon.TxHash]*obscurocommon.L1Tx {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if r.HeightBlock(b) == obscurocommon.L1GenesisHeight {
		return makeMap(b.Transactions())
	}
	newMap := make(map[obscurocommon.TxHash]*obscurocommon.L1Tx)
	p, f := r.ParentBlock(b)
	if !f {
		panic("wtf")
	}
	for k, v := range allIncludedTransactions(p, r, db) {
		newMap[k] = v
	}
	for _, tx := range b.Transactions() {
		newMap[tx.Hash()] = tx
	}
	db.AddTxs(b, newMap)
	return newMap
}

func removeExisting(base []*obscurocommon.L1Tx, toRemove map[obscurocommon.TxHash]*obscurocommon.L1Tx) (r []*obscurocommon.L1Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs types.Transactions) map[obscurocommon.TxHash]*obscurocommon.L1Tx {
	m := make(map[obscurocommon.TxHash]*obscurocommon.L1Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func BlocksBetween(blockA *types.Block, blockB *types.Block, resolver obscurocommon.BlockResolver) []*types.Block {
	if blockA.Hash() == blockB.Hash() {
		return []*types.Block{blockA}
	}
	blocks := make([]*types.Block, 0)
	tempBlock := blockB
	var found bool
	for {
		blocks = append(blocks, tempBlock)
		if tempBlock.Hash() == blockA.Hash() {
			break
		}
		tempBlock, found = resolver.ParentBlock(tempBlock)
		if !found {
			panic("should not happen")
		}
	}
	n := len(blocks)
	result := make([]*types.Block, n)
	for i, block := range blocks {
		result[n-i-1] = block
	}
	return result
}
