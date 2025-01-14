package ethereummock

import (
	"context"
	"fmt"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/common/async"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/integration/simulation/stats"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ten-protocol/go-ten/go/common"

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
func (n *MockEthNetwork) BroadcastBlock(b EncodedL1Block, p EncodedL1Block) {
	bl, _ := b.DecodeBlock()
	for _, m := range n.AllNodes {
		if m.Info().L2ID != n.CurrentNode.Info().L2ID {
			t := m
			async.Schedule(n.delay(), func() { t.P2PReceiveBlock(b, p) })
		} else {
			m.logger.Info(printBlock(bl, m))
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
			async.Schedule(d, func() { t.P2PGossipTx(tx) })
		}
	}
}

// delay returns an expected delay on the l1 network
func (n *MockEthNetwork) delay() time.Duration {
	return testcommon.RndBtwTime(n.avgLatency/10, 2*n.avgLatency)
}

func printBlock(b *types.Block, m *Node) string {
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
		case *common.L1RollupTx:
			r, err := common.DecodeRollup(l1Tx.Rollup)
			if err != nil {
				testlog.Logger().Crit("failed to decode rollup")
			}
			txs = append(txs, fmt.Sprintf("r_%s(nonce=%d)", r.Hash(), tx.Nonce()))

		case *common.L1DepositTx:
			var to gethcommon.Address
			if l1Tx.To != nil {
				to = *l1Tx.To
			}
			txs = append(txs, fmt.Sprintf("deposit(%s=%d)", to, l1Tx.Amount))
		}
	}
	p, err := m.BlockResolver.FetchFullBlock(context.Background(), b.ParentHash())
	if err != nil {
		testlog.Logger().Crit("Should not happen. Could not retrieve parent", log.ErrKey, err)
	}

	return fmt.Sprintf(" create b_%s(Height=%d, RollupNonce=%s)[parent=b_%s]. Txs: %v",
		b.Hash(), b.NumberU64(), b.Header().Nonce, p.Hash(), txs)
}
