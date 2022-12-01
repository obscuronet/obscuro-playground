package p2p

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/obscuronet/go-obscuro/go/common/host"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/obscuronet/go-obscuro/go/common"
)

// MockP2P - models a full network of in memory nodes including artificial random latencies
// Implements the P2p interface
// Will be plugged into each node
type MockP2P struct {
	CurrentNode host.Host
	Nodes       []host.Host

	avgLatency       time.Duration
	avgBlockDuration time.Duration

	listenerInterrupt *int32
}

// NewMockP2P returns an instance of a configured L2 Network (no nodes)
func NewMockP2P(avgBlockDuration time.Duration, avgLatency time.Duration) *MockP2P {
	i := int32(0)
	return &MockP2P{
		avgLatency:        avgLatency,
		avgBlockDuration:  avgBlockDuration,
		listenerInterrupt: &i,
	}
}

func (netw *MockP2P) StartListening(host.Host) {
	// nothing to do here, since communication is direct through the in memory objects
}

func (netw *MockP2P) StopListening() error {
	atomic.StoreInt32(netw.listenerInterrupt, 1)
	return nil
}

func (netw *MockP2P) UpdatePeerList([]string) {
	// Do nothing.
}

func (netw *MockP2P) BroadcastTx(tx common.EncryptedTx) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	for _, node := range netw.Nodes {
		if node.Config().ID.Hex() != netw.CurrentNode.Config().ID.Hex() {
			tempNode := node
			common.Schedule(netw.delay()/2, func() { tempNode.ReceiveTx(tx) })
		}
	}

	return nil
}

func (netw *MockP2P) BroadcastBatch(batch *common.ExtBatch) error {
	if atomic.LoadInt32(netw.listenerInterrupt) == 1 {
		return nil
	}

	encodedBatches, err := rlp.EncodeToBytes([]*common.ExtBatch{batch})
	if err != nil {
		return fmt.Errorf("could not encode batch using RLP. Cause: %w", err)
	}

	for _, node := range netw.Nodes {
		if node.Config().ID.Hex() != netw.CurrentNode.Config().ID.Hex() {
			tempNode := node
			common.Schedule(netw.delay()/2, func() { tempNode.ReceiveBatches(encodedBatches) })
		}
	}

	return nil
}

func (netw *MockP2P) RequestBatches(_ *common.BatchRequest) error {
	panic(errutil.ErrNoImpl)
}

func (netw *MockP2P) SendBatches(_ []*common.ExtBatch, _ string) error {
	panic(errutil.ErrNoImpl)
}

// delay returns an expected delay on the l2
func (netw *MockP2P) delay() time.Duration {
	return testcommon.RndBtwTime(netw.avgLatency/10, 2*netw.avgLatency)
}
