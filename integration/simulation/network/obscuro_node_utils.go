package network

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ten-protocol/go-ten/integration/ethereummock"

	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
	"github.com/ten-protocol/go-ten/integration/simulation/p2p"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"golang.org/x/sync/errgroup"
)

const (
	protocolSeparator = "://"
	networkTCP        = "tcp"
)

func startInMemoryTenNodes(params *params.SimParams, l1Clients []ethadapter.EthClient) []rpc.Client {
	// Create the in memory TEN nodes, each connect each to a geth node
	tenNodes := make([]*hostcontainer.HostContainer, params.NumberOfNodes)
	tenHosts := make([]host.Host, params.NumberOfNodes)
	mockP2PNetw := p2p.NewMockP2PNetwork(params.AvgBlockDuration, params.AvgNetworkLatency, params.NodeWithInboundP2PDisabled)
	blobResolver := ethereummock.NewMockBlobResolver()
	for i := 0; i < params.NumberOfNodes; i++ {
		isGenesis := i == 0

		tenNodes[i] = createInMemTenNode(
			int64(i),
			isGenesis,
			GetNodeType(i),
			params.MgmtContractLib,
			params.Wallets.NodeWallets[i],
			l1Clients[i],
			mockP2PNetw.NewNode(i),
			params.L1TenData.MessageBusAddr,
			params.L1TenData.TenStartBlock,
			params.AvgBlockDuration/3,
			true,
			params.AvgBlockDuration,
			blobResolver,
		)
		tenHosts[i] = tenNodes[i].Host()
	}

	// start each TEN node
	for _, m := range tenNodes {
		t := m
		go func() {
			err := t.Start()
			if err != nil {
				panic(err)
			}
		}()
	}

	// Create a handle to each node
	tenClients := make([]rpc.Client, params.NumberOfNodes)
	for i, node := range tenNodes {
		tenClients[i] = p2p.NewInMemTenClient(node)
	}
	time.Sleep(100 * time.Millisecond)

	return tenClients
}

func createAuthClientsPerWallet(clients []rpc.Client, wallets *params.SimWallets) map[string][]*obsclient.AuthObsClient {
	walletClients := make(map[string][]*obsclient.AuthObsClient)
	// loop through all the L2 wallets we're using and round-robin allocate them the rpc clients we have for each host
	for _, w := range append(wallets.SimObsWallets, wallets.L2FaucetWallet) {
		walletClients[w.Address().String()] = CreateAuthClients(clients, w)
	}
	for _, t := range wallets.Tokens {
		w := t.L2Owner
		walletClients[w.Address().String()] = CreateAuthClients(clients, w)
	}
	return walletClients
}

func CreateAuthClients(clients []rpc.Client, wal wallet.Wallet) []*obsclient.AuthObsClient {
	rpcKey, err := rpc.ReadEnclaveKey(clients[0])
	if err != nil {
		return nil
	}
	authClients := make([]*obsclient.AuthObsClient, len(clients))
	for i, client := range clients {
		vk, err := viewingkey.GenerateViewingKeyForWallet(wal)
		if err != nil {
			panic(err)
		}
		// todo - use a child logger
		encClient, err := rpc.NewEncRPCClient(client, vk, rpcKey, testlog.Logger())
		if err != nil {
			panic(err)
		}
		authClients[i] = obsclient.NewAuthObsClient(encClient)
	}
	return authClients
}

// StopTenNodes stops the TEN nodes and their RPC clients.
func StopTenNodes(clients []rpc.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	for _, client := range clients {
		c := client
		eg.Go(func() error {
			err := c.Call(nil, rpc.StopHost)
			if err != nil {
				testlog.Logger().Error("Could not stop TEN node.", log.ErrKey, err)
				return err
			}
			c.Stop()
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		testlog.Logger().Error(fmt.Sprintf("Error waiting for the TEN nodes to stop - %s", err))
	}

	testlog.Logger().Info("TEN nodes stopped")
}

// CheckHostRPCServersStopped checks whether the hosts' RPC server addresses have been freed up.
func CheckHostRPCServersStopped(hostWSURLS []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	for _, hostWSURL := range hostWSURLS {
		url := hostWSURL
		// We cannot stop the RPC server synchronously. This is because the host itself is being stopped by an RPC
		// call, so there is a deadlock. The RPC server is waiting for all connections to close, but a single
		// connection remains open, waiting for the RPC server to close. Instead, we check whether the RPC port
		// becomes free.
		eg.Go(func() error {
			for !isAddressAvailable(url) {
				time.Sleep(100 * time.Millisecond)
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		panic(fmt.Sprintf("Timed out waiting for the TEN host RPC addresses to become available - %s", err))
	}

	testlog.Logger().Info("TEN host RPC addresses freed")
}

func isAddressAvailable(address string) bool {
	// `net.Listen` requires us to strip the protocol, if it exists.
	addressNoProtocol := address
	splitAddress := strings.Split(address, protocolSeparator)
	if len(splitAddress) == 2 {
		addressNoProtocol = splitAddress[1]
	}

	ln, err := net.Listen(networkTCP, addressNoProtocol)
	if ln != nil {
		err = ln.Close()
		if err != nil {
			testlog.Logger().Error(fmt.Sprintf("could not close listener when checking if address %s was available", address))
		}
	}
	if err != nil {
		return false
	}

	return true
}

// GetNodeType returns the type of the node based on its ID.
func GetNodeType(i int) common.NodeType {
	// Only the genesis node is assigned the role of sequencer.
	if i == 0 {
		return common.Sequencer
	}
	return common.Validator
}
