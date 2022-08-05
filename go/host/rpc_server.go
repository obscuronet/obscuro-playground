package host

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
)

const (
	allOrigins = "*"
)

// An implementation of `host.RPCServer` that reuses the Geth `node` package for client communication.
type rpcServerImpl struct {
	node *node.Node
}

func NewRPCServer(config config.HostConfig, rpcAPIs []rpc.API) RPCServer {
	rpcConfig := node.Config{}
	if config.HasClientRPCHTTP {
		rpcConfig.HTTPHost = config.ClientRPCHost
		rpcConfig.HTTPPort = int(config.ClientRPCPortHTTP)
		// TODO review if this poses a security issue
		rpcConfig.HTTPVirtualHosts = []string{allOrigins}
	}
	if config.HasClientRPCWebsockets {
		rpcConfig.WSHost = config.ClientRPCHost
		rpcConfig.WSPort = int(config.ClientRPCPortWS)
		// TODO review if this poses a security issue
		rpcConfig.WSOrigins = []string{allOrigins}
	}

	rpcServerNode, err := node.New(&rpcConfig)
	if err != nil {
		log.Panic("could not create new client server. Cause: %s", err)
	}
	rpcServerNode.RegisterAPIs(rpcAPIs)

	return &rpcServerImpl{node: rpcServerNode}
}

func (s *rpcServerImpl) Start() {
	if err := s.node.Start(); err != nil {
		log.Panic("could not start node client server. Cause: %s", err)
	}
}

func (s *rpcServerImpl) Stop() error {
	return s.node.Close()
}
