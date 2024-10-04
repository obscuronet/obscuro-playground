package ethereummock

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// findNotIncludedTxs - given a list of transactions, it keeps only the ones that were not included in the block
// todo (#1491) - inefficient
func findNotIncludedTxs(head *types.Block, txs []*types.Transaction, r *blockResolverInMem, db TxDB) []*types.Transaction {
	included := allIncludedTransactions(head, r, db)
	return removeExisting(txs, included)
}

func allIncludedTransactions(b *types.Block, r *blockResolverInMem, db TxDB) map[common.TxHash]*types.Transaction {
	val, found := db.Txs(b)
	if found {
		return val
	}
	if b.NumberU64() == common.L1GenesisHeight {
		return makeMap(b.Transactions())
	}
	newMap := make(map[common.TxHash]*types.Transaction)
	p, err := r.FetchBlock(context.Background(), b.ParentHash())
	if err != nil {
		panic(fmt.Errorf("should not happen. Could not retrieve parent. Cause: %w", err))
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

func removeExisting(base []*types.Transaction, toRemove map[common.TxHash]*types.Transaction) (r []*types.Transaction) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}

func makeMap(txs types.Transactions) map[common.TxHash]*types.Transaction {
	m := make(map[common.TxHash]*types.Transaction)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}
