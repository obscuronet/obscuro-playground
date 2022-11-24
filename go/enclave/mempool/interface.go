package mempool

import (
	"github.com/obscuronet/go-obscuro/go/common"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type Manager interface {
	// FetchMempoolTxs returns all transactions in the mempool
	FetchMempoolTxs() []*common.L2Tx
	// AddMempoolTx adds an transaction to the mempool
	AddMempoolTx(tx *common.L2Tx) error
	// RemoveMempoolTxs removes transactions that are considered immune to re-orgs
	RemoveMempoolTxs(r *obscurocore.Rollup, resolver db.RollupResolver) error
	// CurrentTxs Returns the transactions that should be included in the current rollup
	CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) ([]*common.L2Tx, error)
}
