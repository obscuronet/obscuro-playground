package network

import (
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

type networkInMemGeth struct {
	obscuroClients []rpcclientlib.Client

	// geth
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethadapter.EthClient
	wallets     *params.SimWallets
}

func NewNetworkInMemoryGeth(wallets *params.SimWallets) Network {
	return &networkInMemGeth{
		wallets: wallets,
	}
}

// Create inits and starts the nodes, wires them up, and populates the network objects
func (n *networkInMemGeth) Create(params *params.SimParams, stats *stats.Stats) (*Clients, error) {
	// kickoff the network with the prefunded wallet addresses
	params.MgmtContractAddr, params.BtcErc20Address, params.EthErc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)

	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.BtcErc20Address, params.EthErc20Address)

	// Start the obscuro nodes and return the handles
	var walletClients map[string]rpcclientlib.Client
	n.obscuroClients, walletClients = startInMemoryObscuroNodes(params, stats, n.gethNetwork.GenesisJSON, n.gethClients)

	return &Clients{
		EthClients:     n.gethClients,
		ObscuroClients: n.obscuroClients,
		walletClients:  walletClients,
	}, nil
}

func (n *networkInMemGeth) TearDown() {
	// Stop the obscuro nodes first
	StopObscuroNodes(n.obscuroClients)

	// Stop geth last
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
