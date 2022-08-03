package simulation

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"
)

const azureTestEnv = "AZURE_TEST_ENABLED"

// TODO: we really need tests to demonstrate the unhappy-cases in the attestation scenario:
//	 - if someone puts a dodgy public key on a request with a genuine attestation report they shouldn't get secret
//	 - if owner doesn't match - they shouldn't get secret

// Todo: replace with the IPs of the VMs you are testing, see the azuredeployer README for more info.
var vmIPs = []string{"20.90.164.68"}

// This test creates a network of L2 nodes consisting of just the Azure nodes configured above.
//
// It then injects transactions, and finally checks the resulting output blockchain
// The L2 nodes communicate with each other via sockets, and with their enclave servers via RPC.
// All nodes and enclaves live in the same process, and the Ethereum nodes are mocked out.
func TestAzureEnclaveNodesMonteCarloSimulation(t *testing.T) {
	if os.Getenv(azureTestEnv) == "" {
		t.Skipf("set the variable to run this test: `%s=true`", azureTestEnv)
	}
	setupSimTestLog("azure-enclave")

	numberOfNodes := 1
	numberOfSimWallets := 5

	wallets := params.NewSimWallets(numberOfSimWallets, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)

	simParams := params.SimParams{
		NumberOfNodes:             numberOfNodes,
		AvgBlockDuration:          time.Second,
		SimulationTime:            60 * time.Second,
		L1EfficiencyThreshold:     0.2,
		L2EfficiencyThreshold:     0.3,
		L2ToL1EfficiencyThreshold: 0.4,
		Wallets:                   wallets,
		StartPort:                 integration.StartPortSimulationAzureEnclave,
	}
	simParams.AvgNetworkLatency = simParams.AvgBlockDuration / 15
	simParams.AvgGossipPeriod = simParams.AvgBlockDuration / 3

	if len(vmIPs) > simParams.NumberOfNodes {
		panic(fmt.Sprintf("have %d VMs but only %d nodes", len(vmIPs), simParams.NumberOfNodes))
	}

	// define the network to use
	netw := network.NewNetworkWithAzureEnclaves(vmIPs, wallets)

	testSimulation(t, netw, &simParams)
}
