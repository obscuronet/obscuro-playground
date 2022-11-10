package ethereummock

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// MockEthNetwork - models a full network including artificial random latencies
// This is the gateway through which the mock L1 nodes communicate with each other
type MockEthNetwork struct {
	CurrentNode *Node

	AllNodes []*Node

	// config
	avgLatency       time.Duration
	avgBlockDuration time.Duration

	Stats *stats.Stats
}

// NewMockEthNetwork returns an instance of a configured L1 Network (no nodes)
func NewMockEthNetwork(avgBlockDuration time.Duration, avgLatency time.Duration, stats *stats.Stats) *MockEthNetwork {
	return &MockEthNetwork{
		Stats:            stats,
		avgLatency:       avgLatency,
		avgBlockDuration: avgBlockDuration,
	}
}

// BroadcastBlock broadcast a block to the l1 nodes
func (n *MockEthNetwork) BroadcastBlock(b common.EncodedBlock, p common.EncodedBlock) {
	bl, _ := b.DecodeBlock()
	for _, m := range n.AllNodes {
		if m.Info().L2ID != n.CurrentNode.Info().L2ID {
			t := m
			common.Schedule(n.delay(), func() { t.P2PReceiveBlock(b, p) })
		} else {
			m.logger.Info(printBlock(bl, *m))
		}
	}

	n.Stats.NewBlock(bl)
}

// BroadcastTx Broadcasts the L1 tx containing the rollup to the L1 network
func (n *MockEthNetwork) BroadcastTx(tx *types.Transaction) {
	for _, m := range n.AllNodes {
		if m.Info().L2ID != n.CurrentNode.Info().L2ID {
			t := m
			// the time to broadcast a tx is half that of a L1 block, because it is smaller.
			// todo - find a better way to express this
			d := n.delay() / 2
			common.Schedule(d, func() { t.P2PGossipTx(tx) })
		}
	}
}

// delay returns an expected delay on the l1 network
func (n *MockEthNetwork) delay() time.Duration {
	return testcommon.RndBtwTime(n.avgLatency/10, 2*n.avgLatency)
}

func printBlock(b *types.Block, m Node) string {
	// This is just for printing
	var txs []string
	for _, tx := range b.Transactions() {
		t := m.erc20ContractLib.DecodeTx(tx)
		if t == nil {
			t = m.mgmtContractLib.DecodeTx(tx)
		}

		if t == nil {
			continue
		}

		switch l1Tx := t.(type) {
		case *ethadapter.L1RollupTx:
			r, err := common.DecodeRollup(l1Tx.Rollup)
			if err != nil {
				testlog.Logger().Crit("failed to decode rollup")
			}
			txs = append(txs, fmt.Sprintf("r_%d(nonce=%d)", common.ShortHash(r.Hash()), tx.Nonce()))

		case *ethadapter.L1DepositTx:
			var to uint64
			if l1Tx.To != nil {
				to = common.ShortAddress(*l1Tx.To)
			}
			txs = append(txs, fmt.Sprintf("deposit(%d=%d)", to, l1Tx.Amount))
		}
	}
	p, f := m.Resolver.ParentBlock(b)
	if !f {
		testlog.Logger().Crit("Should not happen. Parent not found")
	}

	return fmt.Sprintf(" create b_%d(Height=%d, RollupNonce=%d)[parent=b_%d]. Txs: %v",
		common.ShortHash(b.Hash()), b.NumberU64(), common.ShortNonce(b.Header().Nonce), common.ShortHash(p.Hash()), txs)
}
