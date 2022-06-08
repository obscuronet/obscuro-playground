package host

import (
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	allOrigins = "*"

	apiNamespaceObscuro  = "obscuro"
	apiNamespaceEthereum = "eth"
	apiNamespaceNetwork  = "net"
	apiVersion1          = "1.0"
)

// An implementation of `host.RPCServer` that reuses the Geth `node` package for client communication.
type rpcServerImpl struct {
	node *node.Node
}

func NewRPCServer(config config.HostConfig, host *Node) RPCServer {
	rpcConfig := node.Config{}
	if config.HasClientRPCHTTP {
		rpcConfig.HTTPHost = config.ClientRPCHost
		rpcConfig.HTTPPort = int(config.ClientRPCPortHTTP)
	}
	if config.HasClientRPCWebsockets {
		rpcConfig.WSHost = config.ClientRPCHost
		rpcConfig.WSPort = int(config.ClientRPCPortWS)
		rpcConfig.WSOrigins = []string{allOrigins}
	}

	rpcServerNode, err := node.New(&rpcConfig)
	if err != nil {
		log.Panic("could not create new client server. Cause: %s", err)
	}

	rpcAPIs := []rpc.API{
		{
			Namespace: apiNamespaceObscuro,
			Version:   apiVersion1,
			Service:   NewObscuroAPI(host),
			Public:    true,
		},
		{
			Namespace: apiNamespaceEthereum,
			Version:   apiVersion1,
			Service:   NewEthereumAPI(host),
			Public:    true,
		},
		{
			Namespace: apiNamespaceNetwork,
			Version:   apiVersion1,
			Service:   NewNetworkAPI(),
			Public:    true,
		},
	}
	rpcServerNode.RegisterAPIs(rpcAPIs)

	return rpcServerImpl{node: rpcServerNode}
}

func (s rpcServerImpl) Start() {
	if err := s.node.Start(); err != nil {
		log.Panic("could not start node client server. Cause: %s", err)
	}
}

func (s rpcServerImpl) Stop() error {
	return s.node.Close()
}
