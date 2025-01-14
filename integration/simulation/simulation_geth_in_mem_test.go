package simulation

import (
	"os"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration"

	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

const (
	gethTestEnv = "GETH_TEST_ENABLED"
)

// TestGethSimulation runs the simulation against a private geth network using Clique (PoA)
func TestGethSimulation(t *testing.T) {
	if os.Getenv(gethTestEnv) == "" {
		t.Skipf("set the variable to run this test: `%s=true`", gethTestEnv)
	}
	setupSimTestLog("geth-in-mem")

	numberOfNodes := 5
	numberOfSimWallets := 5

	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.TenChainID)

	simParams := &params.SimParams{
		NumberOfNodes:         numberOfNodes,
		AvgBlockDuration:      2 * time.Second,
		SimulationTime:        35 * time.Second,
		L1EfficiencyThreshold: 0.2,
		Wallets:               wallets,
		StartPort:             integration.TestPorts.TestGethSimulationPort,
		IsInMem:               true,
		ReceiptTimeout:        30 * time.Second,
		StoppingDelay:         10 * time.Second,
		L1BeaconPort:          integration.TestPorts.TestGethSimulationPort + integration.DefaultPrysmGatewayPortOffset,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15

	testSimulation(t, network.NewNetworkInMemoryGeth(wallets), simParams)
}
