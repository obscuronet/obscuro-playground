package simulation

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration"

	"github.com/ten-protocol/go-ten/integration/simulation/params"

	"github.com/ten-protocol/go-ten/integration/simulation/network"
)

// This test creates a network of L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process. The L1 network is a private geth network using Clique (PoA).
func TestFullNetworkMonteCarloSimulation(t *testing.T) {
	setupSimTestLog("full-network")

	numberOfNodes := 5
	numberOfSimWallets := 5

	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := &params.SimParams{
		NumberOfNodes:              numberOfNodes,
		AvgBlockDuration:           1 * time.Second,
		SimulationTime:             75 * time.Second,
		L1EfficiencyThreshold:      0.2,
		Wallets:                    wallets,
		StartPort:                  integration.StartPortSimulationFullNetwork,
		ReceiptTimeout:             20 * time.Second,
		StoppingDelay:              15 * time.Second,
		NodeWithInboundP2PDisabled: 2,
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15

	testSimulation(t, network.NewNetworkOfSocketNodes(wallets), simParams)
}
