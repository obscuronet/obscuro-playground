package container

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/metrics"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/host"
	"github.com/ten-protocol/go-ten/go/host/p2p"
	"github.com/ten-protocol/go-ten/go/host/rpc/clientapi"
	"github.com/ten-protocol/go-ten/go/host/rpc/enclaverpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethlog "github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

const (
	APIVersion1         = "1.0"
	APINamespaceObscuro = "obscuro"
	APINamespaceEth     = "eth"
	APINamespaceScan    = "scan"
	APINamespaceNetwork = "net"
	APINamespaceTest    = "test"
	APINamespaceDebug   = "debug"
)

type HostContainer struct {
	host           hostcommon.Host
	logger         gethlog.Logger
	metricsService *metrics.Service
	rpcServer      node.Server
}

func (h *HostContainer) Start() error {
	h.metricsService.Start()

	// make sure the rpc server has a host to render requests
	err := h.host.Start()
	if err != nil {
		return err
	}
	h.logger.Info("Started Obscuro host...")
	fmt.Println("Started Obscuro host...")

	if h.rpcServer != nil {
		err = h.rpcServer.Start()
		if err != nil {
			return err
		}

		h.logger.Info("Started Obscuro host RPC Server...")
		fmt.Println("Started Obscuro host RPC Server...")
	}

	return nil
}

func (h *HostContainer) Stop() error {
	// host will not respond to further external requests
	err := h.host.Stop()
	if err != nil {
		return err
	}

	h.metricsService.Stop()

	if h.rpcServer != nil {
		// rpc server cannot be stopped synchronously as it will kill current request
		go func() {
			// make sure it's not killing the connection before returning the response
			time.Sleep(time.Second) // todo review this sleep
			h.rpcServer.Stop()
		}()
	}

	return nil
}

func (h *HostContainer) Host() hostcommon.Host {
	return h.host
}

// NewHostContainerFromConfig uses config to create all HostContainer dependencies and inject them into a new HostContainer
// (Note: it does not start the HostContainer process, `Start()` must be called on the container)
func NewHostContainerFromConfig(parsedConfig *config.HostInputConfig, logger gethlog.Logger) *HostContainer {
	cfg := parsedConfig.ToHostConfig()

	addr, err := wallet.RetrieveAddress(parsedConfig.PrivateKeyString)
	if err != nil {
		panic("unable to retrieve the Node ID")
	}
	cfg.ID = *addr

	// create the logger if not set - used when the testlogger is injected
	if logger == nil {
		logger = log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath, log.NodeIDKey, cfg.ID)
	}

	fmt.Printf("Building host container with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Building host container with config: %+v", cfg))

	ethWallet := wallet.NewInMemoryWalletFromConfig(cfg.PrivateKeyString, cfg.L1ChainID, log.New("wallet", cfg.LogLevel, cfg.LogPath))

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethadapter.NewEthClientFromURL(cfg.L1WebsocketURL, cfg.L1RPCTimeout, cfg.ID, logger)
	if err != nil {
		logger.Crit("could not create Ethereum client.", log.ErrKey, err)
	}

	// update the wallet nonce
	nonce, err := l1Client.Nonce(ethWallet.Address())
	if err != nil {
		logger.Crit("could not retrieve Ethereum account nonce.", log.ErrKey, err)
	}
	ethWallet.SetNonce(nonce)

	// set the Host ID as the Public Key Address
	cfg.ID = ethWallet.Address()

	fmt.Println("Connecting to the enclave...")
	services := host.NewServicesRegistry(logger)
	enclaveClients := make([]common.Enclave, len(cfg.EnclaveRPCAddresses))
	for i, addr := range cfg.EnclaveRPCAddresses {
		enclaveClients[i] = enclaverpc.NewClient(addr, cfg.EnclaveRPCTimeout, logger)
	}
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	metricsService := metrics.New(cfg.MetricsEnabled, cfg.MetricsHTTPPort, logger)

	aggP2P := p2p.NewSocketP2PLayer(cfg, services, p2pLogger, metricsService.Registry())

	rpcServer := node.NewServer(&node.RPCConfig{
		EnableHTTP: cfg.HasClientRPCHTTP,
		HTTPPort:   int(cfg.ClientRPCPortHTTP),
		EnableWs:   cfg.HasClientRPCWebsockets,
		WsPort:     int(cfg.ClientRPCPortWS),
		Host:       cfg.ClientRPCHost,
	}, logger)

	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.ManagementContractAddress, logger)
	obscuroRelevantContracts := []gethcommon.Address{cfg.ManagementContractAddress, cfg.MessageBusAddress}
	l1Repo := l1.NewL1Repository(l1Client, obscuroRelevantContracts, logger)

	return NewHostContainer(cfg, services, aggP2P, l1Client, l1Repo, enclaveClients, mgmtContractLib, ethWallet, rpcServer, logger, metricsService)
}

// NewHostContainer builds a host container with dependency injection rather than from config.
// Useful for testing etc. (want to be able to pass in logger, and also have option to mock out dependencies)
func NewHostContainer(cfg *config.HostConfig, services *host.ServicesRegistry, p2p hostcommon.P2PHostService, l1Client ethadapter.EthClient, l1Repo hostcommon.L1RepoService, enclaveClients []common.Enclave, contractLib mgmtcontractlib.MgmtContractLib, hostWallet wallet.Wallet, rpcServer node.Server, logger gethlog.Logger, metricsService *metrics.Service) *HostContainer {
	h := host.NewHost(cfg, services, p2p, l1Client, l1Repo, enclaveClients, hostWallet, contractLib, logger, metricsService.Registry())

	hostContainer := &HostContainer{
		host:           h,
		logger:         logger,
		rpcServer:      rpcServer,
		metricsService: metricsService,
	}

	if cfg.HasClientRPCHTTP || cfg.HasClientRPCWebsockets {
		filterAPI := clientapi.NewFilterAPI(h, logger)
		rpcServer.RegisterAPIs([]rpc.API{
			{
				Namespace: APINamespaceObscuro,
				Service:   clientapi.NewObscuroAPI(h),
			},
			{
				Namespace: APINamespaceEth,
				Service:   clientapi.NewEthereumAPI(h, logger),
			},
			{
				Namespace: APINamespaceScan,
				Version:   APIVersion1,
				Service:   clientapi.NewScanAPI(h, logger),
				Public:    true,
			},
			{
				Namespace: APINamespaceNetwork,
				Service:   clientapi.NewNetworkAPI(h),
			},
			{
				Namespace: APINamespaceTest,
				Service:   clientapi.NewTestAPI(hostContainer),
			},
			{
				Namespace: APINamespaceEth,
				Service:   filterAPI,
			},
		})

		if cfg.DebugNamespaceEnabled {
			rpcServer.RegisterAPIs([]rpc.API{
				{
					Namespace: APINamespaceDebug,
					Service:   clientapi.NewNetworkDebug(h),
				},
			})
		}
		services.RegisterService(hostcommon.FilterAPIServiceName, filterAPI.NewHeadsService)
	}
	return hostContainer
}
