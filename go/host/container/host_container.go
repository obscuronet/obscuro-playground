package container

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	APIVersion1             = "1.0"
	APINamespaceObscuro     = "obscuro"
	APINamespaceEth         = "eth"
	APINamespaceObscuroScan = "obscuroscan"
	APINamespaceNetwork     = "net"
	APINamespaceTest        = "test"
	APINamespaceDebug       = "debug"
)

type HostContainer struct {
	host           commonhost.Host
	logger         gethlog.Logger
	metricsService *metrics.Service
	rpcServer      clientrpc.Server
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

func (h *HostContainer) Host() commonhost.Host {
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
		logger = log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath)
	}
	logger = logger.New(log.NodeIDKey, cfg.ID, log.CmpKey, log.HostCmp)

	fmt.Printf("Building host container with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Building host container with config: %+v", cfg))

	ethWallet := wallet.NewInMemoryWalletFromConfig(cfg.PrivateKeyString, cfg.L1ChainID, log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath))

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethadapter.NewEthClient(cfg.L1NodeHost, cfg.L1NodeWebsocketPort, cfg.L1RPCTimeout, cfg.ID, logger)
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
	enclaveClient := enclaverpc.NewClient(cfg, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	metricsService := metrics.New(cfg.MetricsEnabled, cfg.MetricsHTTPPort, logger)
	aggP2P := p2p.NewSocketP2PLayer(cfg, p2pLogger, metricsService.Registry())
	rpcServer := clientrpc.NewServer(cfg, logger)

	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.ManagementContractAddress, logger)

	return NewHostContainer(cfg, aggP2P, l1Client, enclaveClient, mgmtContractLib, ethWallet, rpcServer, logger, metricsService)
}

// NewHostContainer builds a host container with dependency injection rather than from config.
// Useful for testing etc. (want to be able to pass in logger, and also have option to mock out dependencies)
func NewHostContainer(
	cfg *config.HostConfig, // provides various parameters that the host needs to function
	p2p commonhost.P2P, // provides the inbound and outbound p2p communication layer
	l1Client ethadapter.EthClient, // provides inbound and outbound L1 connectivity
	enclaveClient common.Enclave, // provides RPC connection to this host's Enclave
	contractLib mgmtcontractlib.MgmtContractLib, // provides the management contract lib injection
	hostWallet wallet.Wallet, // provides an L1 wallet for the host's transactions
	rpcServer clientrpc.Server, // For communication with Obscuro client applications
	logger gethlog.Logger, // provides logging with context
	metricsService *metrics.Service, // provides the metrics service for other packages to use
) *HostContainer {
	h := host.NewHost(cfg, p2p, l1Client, enclaveClient, hostWallet, contractLib, logger, metricsService.Registry())

	hostContainer := &HostContainer{
		host:           h,
		logger:         logger,
		rpcServer:      rpcServer,
		metricsService: metricsService,
	}

	if cfg.HasClientRPCHTTP || cfg.HasClientRPCWebsockets {
		rpcServer.RegisterAPIs([]rpc.API{
			{
				Namespace: APINamespaceObscuro,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewEthereumAPI(h, logger),
				Public:    true,
			},
			{
				Namespace: APINamespaceObscuroScan,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroScanAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceNetwork,
				Version:   APIVersion1,
				Service:   clientapi.NewNetworkAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceTest,
				Version:   APIVersion1,
				Service:   clientapi.NewTestAPI(hostContainer),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewFilterAPI(h, logger),
				Public:    true,
			},
		})

		if cfg.DebugNamespaceEnabled {
			rpcServer.RegisterAPIs([]rpc.API{
				{
					Namespace: APINamespaceDebug,
					Version:   APIVersion1,
					Service:   clientapi.NewNetworkDebug(h),
					Public:    true,
				},
			})
		}
	}

	return hostContainer
}
