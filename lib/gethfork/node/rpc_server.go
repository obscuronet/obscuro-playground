package node

import (
	"net/http"

	"github.com/ten-protocol/go-ten/tools/walletextension/httpapi"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

const (
	allOrigins = "*"
)

type RPCConfig struct {
	Host       string
	EnableHTTP bool
	HTTPPort   int
	EnableWs   bool
	WsPort     int
	WsPath     string
	HTTPPath   string

	// ExposedURLParamNames - url prams that are available in the services
	ExposedURLParamNames []string
}

// Server manages the lifeycle of an RPC Server
type Server interface {
	Start() error
	Stop()
	RegisterAPIs(apis []rpc.API)
	RegisterRoutes(routes []httpapi.Route)
}

// An implementation of `host.Server` that reuses the Geth `node` package for client communication.
type serverImpl struct {
	node   *Node
	logger gethlog.Logger
}

func NewServer(config *RPCConfig, logger gethlog.Logger) Server {
	rpcConfig := Config{
		Logger:               logger,
		ExposedURLParamNames: config.ExposedURLParamNames,
	}
	if config.EnableHTTP {
		rpcConfig.HTTPHost = config.Host
		rpcConfig.HTTPPort = config.HTTPPort
		// todo (@pedro) - review if this poses a security issue
		rpcConfig.HTTPVirtualHosts = []string{allOrigins}
		rpcConfig.HTTPPathPrefix = config.HTTPPath
	}
	if config.EnableWs {
		rpcConfig.WSHost = config.Host
		rpcConfig.WSPort = config.WsPort
		// todo (@pedro) - review if this poses a security issue
		rpcConfig.WSOrigins = []string{allOrigins}
		rpcConfig.WSPathPrefix = config.WsPath
	}

	rpcServerNode, err := New(&rpcConfig)
	if err != nil {
		logger.Crit("could not create new client server.", log.ErrKey, err)
	}

	return &serverImpl{node: rpcServerNode, logger: logger}
}

func (s *serverImpl) RegisterAPIs(apis []rpc.API) {
	s.node.RegisterAPIs(apis)
}

func (s *serverImpl) RegisterRoutes(routes []httpapi.Route) {
	for _, route := range routes {
		s.node.RegisterHandler(route.Name, route.Name, http.HandlerFunc(route.Func))
	}
}

func (s *serverImpl) Start() error {
	return s.node.Start()
}

func (s *serverImpl) Stop() {
	err := s.node.Close()
	if err != nil {
		s.logger.Crit("could not stop node client server.", log.ErrKey, err)
	}
}
