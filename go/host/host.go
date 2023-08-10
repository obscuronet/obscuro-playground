package host

import (
	"encoding/json"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"os"

	"github.com/kamilsk/breaker"

	"github.com/obscuronet/go-obscuro/go/host/l2"

	"github.com/obscuronet/go-obscuro/go/host/enclave"
	"github.com/obscuronet/go-obscuro/go/host/l1"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/events"
	"github.com/obscuronet/go-obscuro/go/responses"
)

// Implementation of host.Host.
type host struct {
	hostcommon.APIDBRepository  // Access to the hosts public available db
	hostcommon.APIEnclaveClient // Access to the Enclave operations

	config  *config.HostConfig
	shortID uint64

	// ignore incoming requests
	stopControl *stopcontrol.StopControl

	logger gethlog.Logger

	interrupter breaker.Interface
	services    *Services
	subsService *events.LogEventManager
}

func NewHost(
	config *config.HostConfig,
	database *db.DB,
	l1Client ethadapter.EthClient,
	l1Publisher *l1.Publisher,
	p2p P2PHostService,
	enclaveClient common.Enclave,
	logger gethlog.Logger,
) hostcommon.Host {
	// set up host metadata
	hostIdentity := hostcommon.NewIdentity(config)

	// set up control mechanisms
	stopControl := stopcontrol.New()
	interrupter := breaker.Multiplex(
		breaker.BreakBySignal(
			os.Kill,
			os.Interrupt,
		),
	)

	// create the Host controlled Services
	l1RepoService := l1.NewL1Repository(l1Client, logger)
	l2RepoService := l2.NewBatchRepository(config, p2p, database, logger)
	logEventService := events.NewLogEventManager(enclaveClient, logger)
	// todo review why there are 2 inits for the guardian ( guardian component + service )
	// both exist at the same level and both have start/stop but the service should be the only available interface to the host
	enclaveGuardianLib := enclave.NewGuardian(
		config,
		hostIdentity,
		p2p,
		l1Publisher,
		l1RepoService,
		l2RepoService,
		logEventService,
		enclaveClient,
		database,
		interrupter,
		logger,
	)
	enclaveService := enclave.NewService(hostIdentity, p2p, enclaveGuardianLib, logger)

	// log startup data for sanity control / tests
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Host service created with following config:", log.CfgKey, string(jsonConfig))

	return &host{
		APIDBRepository:  database,
		APIEnclaveClient: enclaveClient,

		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),
		logger:  logger,

		stopControl: stopControl,
		interrupter: interrupter,
		subsService: logEventService,
		services: &Services{
			P2P:            p2p,
			L1Repo:         l1RepoService,
			L2Repo:         l2RepoService,
			EnclaveService: enclaveService,
		},
	}
}

// Start validates the host config and starts the Host in a go routine - immediately returns after
func (h *host) Start() error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Start with the host stopping"))
	}

	h.validateConfig()

	// start all registered services
	for i, service := range h.services.All() {
		err := service.Start()
		if err != nil {
			return fmt.Errorf("could not start service=%d: %w", i, err)
		}
	}

	tomlConfig, err := toml.Marshal(h.config)
	if err != nil {
		return fmt.Errorf("could not print host config - %w", err)
	}
	h.logger.Info("Host started with following config", log.CfgKey, string(tomlConfig))

	return nil
}

func (h *host) Config() *config.HostConfig {
	return h.config
}

func (h *host) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitAndBroadcastTx with the host stopping"))
	}
	return h.services.EnclaveService.SubmitAndBroadcastTx(encryptedParams)
}

func (h *host) Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Subscribe with the host stopping"))
	}
	return h.subsService.Subscribe(id, encryptedLogSubscription, matchedLogsCh)
}

func (h *host) Unsubscribe(id rpc.ID) {
	if h.stopControl.IsStopping() {
		h.logger.Error("requested Subscribe with the host stopping")
	}
	h.subsService.Unsubscribe(id)
}

func (h *host) Stop() error {
	// block all incoming requests
	h.stopControl.Stop()

	h.logger.Info("Host received a stop command. Attempting shutdown...")
	h.interrupter.Close()

	// stop all registered services
	for i, service := range h.services.All() {
		if err := service.Stop(); err != nil {
			h.logger.Error("failed to stop service", "service", i, log.ErrKey, err)
		}
	}

	h.logger.Info("Host shut down complete.")
	return nil
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (*hostcommon.HealthCheck, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested HealthCheck with the host stopping"))
	}

	healthErrors := make([]string, 0)

	// loop through all registered services and collect their health statuses
	for i, service := range h.services.All() {
		status := service.HealthStatus()
		if !status.OK() {
			healthErrors = append(healthErrors, fmt.Sprintf("[%d] not healthy - %s", i, status.Message()))
		}
	}

	return &hostcommon.HealthCheck{
		OverallHealth: len(healthErrors) == 0,
		Errors:        healthErrors,
	}, nil
}

// Checks the host config is valid.
func (h *host) validateConfig() {
	if h.config.IsGenesis && h.config.NodeType != common.Sequencer {
		h.logger.Crit("genesis node must be the sequencer")
	}
	if !h.config.IsGenesis && h.config.NodeType == common.Sequencer {
		h.logger.Crit("only the genesis node can be a sequencer")
	}

	if h.config.P2PPublicAddress == "" {
		h.logger.Crit("the host must specify a public P2P address")
	}
}
