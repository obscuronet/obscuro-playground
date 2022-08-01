package network

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"

	"github.com/obscuronet/go-obscuro/go/rpcclientlib"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkOfSocketNodes struct {
	obscuroClients   []rpcclientlib.Client
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

func (n *networkOfSocketNodes) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	// kickoff the network with the prefunded wallet addresses
	params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)

	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address)

	// Start the enclaves
	startRemoteEnclaveServers(0, params, stats)

	n.enclaveAddresses = make([]string, params.NumberOfNodes)
	for i := 0; i < params.NumberOfNodes; i++ {
		n.enclaveAddresses[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
	}

	obscuroClients, walletClients := startStandaloneObscuroNodes(params, stats, n.gethClients, n.enclaveAddresses)
	n.obscuroClients = obscuroClients

	return &RPCHandles{
		EthClients:                    n.gethClients,
		ObscuroClients:                n.obscuroClients,
		VirtualWalletExtensionClients: walletClients,
	}, nil
}

func (n *networkOfSocketNodes) TearDown() {
	// Stop the Obscuro nodes first (each host will attempt to shut down its enclave as part of shutdown).
	StopObscuroNodes(n.obscuroClients)
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
