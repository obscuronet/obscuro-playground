package network

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"

	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkOfSocketNodes struct {
	l2Clients        []rpc.Client
	hostRPCAddresses []string
	enclaveAddresses []string

	// geth
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethadapter.EthClient
	wallets     *params.SimWallets
}

func NewNetworkOfSocketNodes(wallets *params.SimWallets) Network {
	return &networkOfSocketNodes{
		wallets: wallets,
	}
}

func (n *networkOfSocketNodes) Create(simParams *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	// kickoff the network with the prefunded wallet addresses
	simParams.L1SetupData, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		simParams.StartPort,
		simParams.NumberOfNodes,
		int(simParams.AvgBlockDuration.Seconds()),
	)

	simParams.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(simParams.L1SetupData.MgmtContractAddress, testlog.Logger())
	simParams.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(simParams.L1SetupData.MgmtContractAddress,
		simParams.L1SetupData.ObxErc20Address, simParams.L1SetupData.EthErc20Address)

	// Start the enclaves
	startRemoteEnclaveServers(simParams)

	n.enclaveAddresses = make([]string, simParams.NumberOfNodes)
	for i := 0; i < simParams.NumberOfNodes; i++ {
		n.enclaveAddresses[i] = fmt.Sprintf("%s:%d", Localhost, simParams.StartPort+DefaultEnclaveOffset+i)
	}

	l2Clients, hostRPCAddresses := startStandaloneObscuroNodes(simParams, n.gethClients, n.enclaveAddresses)
	n.l2Clients = l2Clients
	n.hostRPCAddresses = hostRPCAddresses

	obscuroClients := make([]*obsclient.ObsClient, simParams.NumberOfNodes)
	for idx, l2Client := range n.l2Clients {
		obscuroClients[idx] = obsclient.NewObsClient(l2Client)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, simParams.Wallets)

	return &RPCHandles{
		EthClients:     n.gethClients,
		ObscuroClients: obscuroClients,
		RPCClients:     n.l2Clients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *networkOfSocketNodes) TearDown() {
	// Stop the Obscuro nodes first (each host will attempt to shut down its enclave as part of shutdown).
	StopObscuroNodes(n.l2Clients)
	StopGethNetwork(n.gethClients, n.gethNetwork)
	CheckHostRPCServersStopped(n.hostRPCAddresses)
}
