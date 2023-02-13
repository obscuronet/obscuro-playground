package container

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/api"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaveclient"
	"github.com/obscuronet/go-obscuro/go/host/rpc/server"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
)

type HostContainer struct {
	host           commonhost.Host
	logger         gethlog.Logger
	metricsService *metrics.Service
	rpcServer      server.Server
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
	h.metricsService.Stop()

	// make sure the rpc server does not request services from a stopped host
	if h.rpcServer != nil {
		h.rpcServer.Stop()
	}

	h.host.Stop()
	return nil
}

func (h *HostContainer) Host() commonhost.Host {
	return h.host
}

func (h *HostContainer) Logger() gethlog.Logger {
	return h.logger
}

// NewHostContainerFromConfig uses config to create all HostContainer dependencies and inject them into a new HostContainer
// (Note: it does not start the HostContainer process, `Start()` must be called on the container)
func NewHostContainerFromConfig(parsedConfig *config.HostInputConfig) *HostContainer {
	cfg := parsedConfig.ToHostConfig()

	logger := log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath, log.NodeIDKey, cfg.ID)
	fmt.Printf("Building host container with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Building host container with config: %+v", cfg))

	// set the Host ID as the Public Key Address
	ethWallet := wallet.NewInMemoryWalletFromConfig(cfg.PrivateKeyString, cfg.L1ChainID, log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath))
	cfg.ID = ethWallet.Address()

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
	enclaveClient := enclaveclient.NewClient(cfg, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	metricsService := metrics.New(cfg.MetricsEnabled, cfg.MetricsHTTPPort, logger)
	aggP2P := p2p.NewSocketP2PLayer(cfg, p2pLogger, metricsService.Registry())
	rpcServer := server.NewServer(cfg, logger)

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
	rpcServer server.Server, // For communication with Obscuro client applications
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
		var apis []rpc.API
		if cfg.NodeType == common.Validator {
			apis = api.ValidatorAPIs(hostContainer, h, logger)
		} else if cfg.NodeType == common.Sequencer {
			apis = api.SequencerAPIs(hostContainer, h, logger)
		}

		rpcServer.RegisterAPIs(apis)
	}

	return hostContainer
}
