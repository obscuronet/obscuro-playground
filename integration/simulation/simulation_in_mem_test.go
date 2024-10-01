package simulation

import (
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

// This test creates a network of in memory L1 and L2 nodes, then injects transactions, and finally checks the resulting output blockchain.
// Running it long enough with various parameters will test many corner cases without having to explicitly write individual tests for them.
// The unit of time is the "AvgBlockDuration" - which is the average time between L1 blocks, which are the carriers of rollups.
// Everything else is reported to this value. This number has to be adjusted in conjunction with the number of nodes. If it's too low,
// the CPU usage will be very high during the simulation which might give inconclusive results.
func TestInMemoryMonteCarloSimulation(t *testing.T) {
	setupSimTestLog("in-mem")

	// todo (#718) - try increasing this back to 7 once faster-finality model is optimised
	numberOfNodes := 5
	numberOfSimWallets := 10
	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.TenChainID)

	simParams := params.SimParams{
		NumberOfNodes:              5,
		AvgBlockDuration:           250 * time.Millisecond, // Increased from 180ms
		SimulationTime:             45 * time.Second,       // Increased from 45s
		L1EfficiencyThreshold:      0.2,
		MgmtContractLib:            ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib:           ethereummock.NewERC20ContractLibMock(),
		BlobResolver:               ethereummock.NewBlobResolver(0, ethereummock.SecondsPerSlot),
		Wallets:                    wallets,
		StartPort:                  integration.TestPorts.TestInMemoryMonteCarloSimulationPort,
		IsInMem:                    true,
		L1TenData:                  &params.L1TenData{},
		ReceiptTimeout:             10 * time.Second, // Increased from 5s
		StoppingDelay:              20 * time.Second, // Increased from 15s
		NodeWithInboundP2PDisabled: 2,
		L1BeaconPort:               integration.TestPorts.TestInMemoryMonteCarloSimulationPort + integration.DefaultPrysmGatewayPortOffset,
	}

	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15

	testSimulation(t, network.NewBasicNetworkOfInMemoryNodes(), &simParams)
}
