package walletextension

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/metrics"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/tools/walletextension/httpapi"

	"github.com/ten-protocol/go-ten/tools/walletextension/rpcapi"

	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/keymanager"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
	"golang.org/x/crypto/acme/autocert"
)

type Container struct {
	stopControl     *stopcontrol.StopControl
	logger          gethlog.Logger
	rpcServer       node.Server
	services        *services.Services
	newHeadsService *subscription.NewHeadsService
}

func NewContainerFromConfig(config wecommon.Config, logger gethlog.Logger) *Container {
	// create the account manager with a single unauthenticated connection
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress

	// get the encryption key (method is determined by the config)
	encryptionKey, err := keymanager.GetEncryptionKey(config, logger)
	if err != nil {
		logger.Crit("unable to get encryption key", log.ErrKey, err)
		os.Exit(1)
	}

	// Create metrics tracker
	var metricsTracker metrics.Metrics
	if config.DBType == "cosmosDB" {
		metricsStorage, err := storage.NewMetricsStorage(config.DBType, config.DBConnectionURL)
		if err != nil {
			logger.Crit("unable to create metrics storage", log.ErrKey, err)
			os.Exit(1)
		}
		metricsTracker = metrics.NewMetricsTracker(metricsStorage)
	} else {
		metricsTracker = metrics.NewNoOpMetricsTracker()
	}

	// start the database with the encryption key
	userStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride, encryptionKey, logger)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
		os.Exit(1)
	}

	// captures version in the env vars
	version := os.Getenv("OBSCURO_GATEWAY_VERSION")
	if version == "" {
		version = "dev"
	}

	stopControl := stopcontrol.New()
	walletExt := services.NewServices(hostRPCBindAddrHTTP, hostRPCBindAddrWS, userStorage, stopControl, version, logger, metricsTracker, &config)
	cfg := &node.RPCConfig{
		EnableHTTP: true,
		HTTPPort:   config.WalletExtensionPortHTTP,
		EnableWs:   true,
		WsPort:     config.WalletExtensionPortWS,
		WsPath:     wecommon.APIVersion1 + "/",
		HTTPPath:   wecommon.APIVersion1 + "/",
		Host:       config.WalletExtensionHost,
	}

	// check if TLS is enabled
	if config.EnableTLS {
		// Create autocert manager for automatic certificate management
		// Generating a certificate consists of the following steps:
		// generating a new private key
		// domain ownership verification (HTTP-01 challenge since certManager.HTTPHandler(nil) is set)
		// Certificate Signing Request (CRS) is generated
		// CRS is sent to CA (Let's Encrypt) via ACME (automated certificate management environment) client
		// CA verifies CRS and issues a certificate
		// Store certificate and private key in certificate storage based on the database type
		certStorage, err := storage.NewCertStorage(config.DBType, config.DBConnectionURL, encryptionKey, config.EncryptingCertificateEnabled, logger)
		if err != nil {
			logger.Crit("unable to create certificate storage", log.ErrKey, err)
			os.Exit(1)
		}

		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.TLSDomain),
			Cache:      certStorage,
		}

		// Create HTTP-01 challenge handler
		httpServer := &http.Server{
			Addr:    ":http", // Port 80
			Handler: certManager.HTTPHandler(nil),
		}
		go httpServer.ListenAndServe() // Start HTTP server for ACME challenges

		tlsConfig := &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}

		// Update RPC server config to use TLS
		cfg.TLSConfig = tlsConfig
	}

	rpcServer := node.NewServer(cfg, logger)

	rpcServer.RegisterRoutes(httpapi.NewHTTPRoutes(walletExt))

	// register all RPC endpoints exposed by a typical Geth node
	rpcServer.RegisterAPIs([]gethrpc.API{
		{
			Namespace: "eth",
			Service:   rpcapi.NewEthereumAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewBlockChainAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewTransactionAPI(walletExt),
		}, {
			Namespace: "txpool",
			Service:   rpcapi.NewTxPoolAPI(walletExt),
		}, {
			Namespace: "debug",
			Service:   rpcapi.NewDebugAPI(walletExt),
		}, {
			Namespace: "sessionkeys",
			Service:   rpcapi.NewSessionKeyAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewFilterAPI(walletExt),
		}, {
			Namespace: "net",
			Service:   rpcapi.NewNetAPI(walletExt),
		}, {
			Namespace: "web3",
			Service:   rpcapi.NewWeb3API(walletExt),
		},
	})

	// Add metrics tracker to stop sequence
	stopControl.OnStop(func() {
		metricsTracker.Stop()
	})

	return &Container{
		stopControl:     stopControl,
		rpcServer:       rpcServer,
		newHeadsService: walletExt.NewHeadsService,
		services:        walletExt,
		logger:          logger,
	}
}

// Start starts the wallet extension container
func (w *Container) Start() error {
	err := w.newHeadsService.Start()
	if err != nil {
		return err
	}

	err = w.rpcServer.Start()
	if err != nil {
		return err
	}
	return nil
}

func (w *Container) Stop() error {
	w.stopControl.Stop()
	_ = w.newHeadsService.Stop()

	if w.rpcServer != nil {
		// rpc server cannot be stopped synchronously as it will kill current request
		go func() {
			// make sure it's not killing the connection before returning the response
			time.Sleep(time.Second) // todo review this sleep
			w.rpcServer.Stop()
		}()
	}

	w.services.Stop()
	return nil
}
