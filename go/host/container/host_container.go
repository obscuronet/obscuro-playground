package container

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

type HostContainer struct {
	host   commonhost.Host
	logger gethlog.Logger
}

func (h *HostContainer) Start() error {
	fmt.Println("Starting Obscuro host...")
	h.logger.Info("Starting Obscuro host...")
	h.host.Start()
	return nil
}

func (h *HostContainer) Stop() error {
	h.host.Stop()
	return nil
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

	enclaveClient := enclaverpc.NewClient(cfg, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	aggP2P := p2p.NewSocketP2PLayer(cfg, p2pLogger)

	return NewHostContainer(cfg, aggP2P, l1Client, enclaveClient, ethWallet, logger)
}

// NewHostContainer builds a host container with dependency injection rather than from config.
// Useful for testing etc. (want to be able to pass in logger, and also have option to mock out dependencies)
func NewHostContainer(
	cfg *config.HostConfig, // provides various parameters that the host needs to function
	p2p commonhost.P2P, // provides the inbound and outbound p2p communication layer
	l1Client ethadapter.EthClient, // provides inbound and outbound L1 connectivity
	enclaveClient common.Enclave, // provides RPC connection to this host's Enclave
	hostWallet wallet.Wallet, // provides an L1 wallet for the host's transactions
	logger gethlog.Logger, // provides logging with context
) *HostContainer {
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.RollupContractAddress, logger)
	h := host.NewHost(cfg, p2p, l1Client, enclaveClient, hostWallet, mgmtContractLib, logger)
	return &HostContainer{
		host:   h,
		logger: logger,
	}
}
