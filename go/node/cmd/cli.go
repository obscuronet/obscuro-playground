package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
)

var (
	startAction      = "start"
	upgradeAction    = "upgrade"
	validNodeActions = []string{startAction, upgradeAction}
)

// NodeConfigCLI represents the configurations passed into the node over CLI
type NodeConfigCLI struct {
	nodeAction              string
	nodeType                string
	isGenesis               bool
	numEnclaves             int
	isSGXEnabled            bool
	enclaveDockerImage      string
	hostDockerImage         string
	l1WebsocketURL          string
	hostP2PPort             int
	hostP2PHost             string
	hostP2PPublicAddr       string
	enclaveHTTPPort         int
	enclaveWSPort           int
	privateKey              string
	hostID                  string
	sequencerP2PAddr        string
	sequencerUpgraderAddr   string
	managementContractAddr  string
	messageBusContractAddr  string
	l1Start                 string
	pccsAddr                string
	edgelessDBImage         string
	hostHTTPPort            int
	hostWSPort              int
	nodeName                string
	isDebugNamespaceEnabled bool
	logLevel                int
	isInboundP2PDisabled    bool
	batchInterval           string // format like 500ms or 2s (any time parsable by time.ParseDuration())
	maxBatchInterval        string // format like 500ms or 2s (any time parsable by time.ParseDuration())
	rollupInterval          string // format like 500ms or 2s (any time parsable by time.ParseDuration())
	l1ChainID               int
	postgresDBHost          string
	l1BeaconUrl             string
	l1BlobArchiveUrl        string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *NodeConfigCLI {
	cfg := &NodeConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	nodeName := flag.String(nodeNameFlag, "obscuronode", flagUsageMap[nodeNameFlag])
	nodeType := flag.String(nodeTypeFlag, "", flagUsageMap[nodeTypeFlag])
	isGenesis := flag.Bool(isGenesisFlag, false, flagUsageMap[isGenesisFlag])
	numEnclaves := flag.Int(numEnclavesFlag, 1, flagUsageMap[numEnclavesFlag])
	isSGXEnabled := flag.Bool(isSGXEnabledFlag, false, flagUsageMap[isSGXEnabledFlag])
	enclaveDockerImage := flag.String(enclaveDockerImageFlag, "", flagUsageMap[enclaveDockerImageFlag])
	hostDockerImage := flag.String(hostDockerImageFlag, "", flagUsageMap[hostDockerImageFlag])
	l1WebsocketURL := flag.String(l1WebsocketURLFlag, "ws://eth2network:9000", flagUsageMap[l1WebsocketURLFlag])
	hostP2PPort := flag.Int(hostP2PPortFlag, 14000, flagUsageMap[hostP2PPortFlag])
	hostP2PHost := flag.String(hostP2PHostFlag, "0.0.0.0", flagUsageMap[hostP2PHostFlag])
	hostP2PPublicAddr := flag.String(hostP2PPublicAddrFlag, "", flagUsageMap[hostP2PPublicAddrFlag])
	hostHTTPPort := flag.Int(hostHTTPPortFlag, 80, flagUsageMap[hostHTTPPortFlag])
	hostWSPort := flag.Int(hostWSPortFlag, 81, flagUsageMap[hostWSPortFlag])
	enclaveHTTPPort := flag.Int(enclaveHTTPPortFlag, 11000, flagUsageMap[enclaveHTTPPortFlag])
	enclaveWSPort := flag.Int(enclaveWSPortFlag, 11001, flagUsageMap[enclaveWSPortFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	hostID := flag.String(hostIDFlag, "", flagUsageMap[hostIDFlag])
	sequencerP2PAddr := flag.String(sequencerP2PAddrFlag, "", flagUsageMap[sequencerP2PAddrFlag])
	managementContractAddr := flag.String(managementContractAddrFlag, "", flagUsageMap[managementContractAddrFlag])
	messageBusContractAddr := flag.String(messageBusContractAddrFlag, "", flagUsageMap[messageBusContractAddrFlag])
	l1Start := flag.String(l1StartBlockFlag, "", flagUsageMap[l1StartBlockFlag])
	pccsAddr := flag.String(pccsAddrFlag, "", flagUsageMap[pccsAddrFlag])
	edgelessDBImage := flag.String(edgelessDBImageFlag, "ghcr.io/edgelesssys/edgelessdb-sgx-4gb:v0.3.2", flagUsageMap[edgelessDBImageFlag])
	isDebugNamespaceEnabled := flag.Bool(isDebugNamespaceEnabledFlag, false, flagUsageMap[isDebugNamespaceEnabledFlag])
	logLevel := flag.Int(logLevelFlag, 3, flagUsageMap[logLevelFlag])
	isInboundP2PDisabled := flag.Bool(isInboundP2PDisabledFlag, false, flagUsageMap[isInboundP2PDisabledFlag])
	batchInterval := flag.String(batchIntervalFlag, "1s", flagUsageMap[batchIntervalFlag])
	maxBatchInterval := flag.String(maxBatchIntervalFlag, "1s", flagUsageMap[maxBatchIntervalFlag])
	rollupInterval := flag.String(rollupIntervalFlag, "3s", flagUsageMap[rollupIntervalFlag])
	l1ChainID := flag.Int(l1ChainIDFlag, 1337, flagUsageMap[l1ChainIDFlag])
	postgresDBHost := flag.String(postgresDBHostFlag, "dd", flagUsageMap[postgresDBHostFlag])
	l1BeaconUrl := flag.String(l1BeaconUrlFlag, "eth2network:126000", flagUsageMap[l1BeaconUrlFlag])
	l1BlobArchiveUrl := flag.String(l1BlobArchiveUrlFlag, "", flagUsageMap[l1BlobArchiveUrlFlag])
	systemContractsUpgrader := flag.String(systemContractsUpgraderFlag, "", flagUsageMap[systemContractsUpgraderFlag])
	flag.Parse()
	cfg.nodeName = *nodeName
	cfg.nodeType = *nodeType
	cfg.isGenesis = *isGenesis
	cfg.numEnclaves = *numEnclaves
	cfg.isSGXEnabled = *isSGXEnabled
	cfg.enclaveDockerImage = *enclaveDockerImage
	cfg.hostDockerImage = *hostDockerImage
	cfg.l1WebsocketURL = *l1WebsocketURL
	cfg.hostP2PPort = *hostP2PPort
	cfg.hostP2PHost = *hostP2PHost
	cfg.hostP2PPublicAddr = *hostP2PPublicAddr
	cfg.enclaveHTTPPort = *enclaveHTTPPort
	cfg.enclaveWSPort = *enclaveWSPort
	cfg.privateKey = *privateKey
	cfg.hostID = *hostID
	cfg.sequencerP2PAddr = *sequencerP2PAddr
	cfg.managementContractAddr = *managementContractAddr
	cfg.messageBusContractAddr = *messageBusContractAddr
	cfg.l1Start = *l1Start
	cfg.pccsAddr = *pccsAddr
	cfg.edgelessDBImage = *edgelessDBImage
	cfg.hostHTTPPort = *hostHTTPPort
	cfg.hostWSPort = *hostWSPort
	cfg.isDebugNamespaceEnabled = *isDebugNamespaceEnabled
	cfg.logLevel = *logLevel
	cfg.isInboundP2PDisabled = *isInboundP2PDisabled
	cfg.batchInterval = *batchInterval
	cfg.maxBatchInterval = *maxBatchInterval
	cfg.rollupInterval = *rollupInterval
	cfg.l1ChainID = *l1ChainID
	cfg.postgresDBHost = *postgresDBHost
	cfg.l1BeaconUrl = *l1BeaconUrl
	cfg.l1BlobArchiveUrl = *l1BlobArchiveUrl
	cfg.sequencerUpgraderAddr = *systemContractsUpgrader

	cfg.nodeAction = flag.Arg(0)
	if !validateNodeAction(cfg.nodeAction) {
		if cfg.nodeAction == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), cfg.nodeAction)
		}
		os.Exit(1)
	}

	return cfg
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}

func NodeCLIConfigToTenConfig(cliCfg *NodeConfigCLI) *config.TenConfig {
	nodeType, err := common.ToNodeType(cliCfg.nodeType)
	if err != nil {
		fmt.Printf("Error converting node type: %v\n", err)
		os.Exit(1)
	}

	// load default Ten config before we apply the CLI overrides
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Printf("Error loading default Ten config: %v\n", err)
		os.Exit(1)
	}
	enclaveAddresses := make([]string, cliCfg.numEnclaves)
	for i := 0; i < cliCfg.numEnclaves; i++ {
		enclaveAddresses[i] = fmt.Sprintf("%s-enclave-%d:%d",
			cliCfg.nodeName, i, cliCfg.enclaveWSPort)
	}

	tenCfg.Network.L1.ChainID = int64(cliCfg.l1ChainID)
	tenCfg.Network.L1.L1Contracts.ManagementContract = gethcommon.HexToAddress(cliCfg.managementContractAddr)
	tenCfg.Network.L1.L1Contracts.MessageBusContract = gethcommon.HexToAddress(cliCfg.messageBusContractAddr)
	tenCfg.Network.L1.StartHash = gethcommon.HexToHash(cliCfg.l1Start)
	tenCfg.Network.Batch.Interval, err = time.ParseDuration(cliCfg.batchInterval)
	if err != nil {
		fmt.Printf("Error parsing batch interval '%s': %v\n", cliCfg.batchInterval, err)
		os.Exit(1)
	}
	tenCfg.Network.Batch.MaxInterval, err = time.ParseDuration(cliCfg.maxBatchInterval)
	if err != nil {
		fmt.Printf("Error parsing max batch interval '%s': %v\n", cliCfg.maxBatchInterval, err)
		os.Exit(1)
	}
	tenCfg.Network.Rollup.Interval, err = time.ParseDuration(cliCfg.rollupInterval)
	if err != nil {
		fmt.Printf("Error parsing rollup interval '%s': %v\n", cliCfg.rollupInterval, err)
		os.Exit(1)
	}
	tenCfg.Network.Sequencer.P2PAddress = cliCfg.sequencerP2PAddr
	tenCfg.Network.Sequencer.SystemContractsUpgrader = gethcommon.HexToAddress(cliCfg.sequencerUpgraderAddr)

	tenCfg.Node.ID = cliCfg.hostID
	tenCfg.Node.Name = cliCfg.nodeName
	tenCfg.Node.NodeType = nodeType
	tenCfg.Node.IsGenesis = cliCfg.isGenesis
	tenCfg.Node.HostAddress = cliCfg.hostP2PPublicAddr
	tenCfg.Node.PrivateKeyString = cliCfg.privateKey

	tenCfg.Host.DB.UseInMemory = false // these nodes always use a persistent DB
	tenCfg.Host.DB.PostgresHost = cliCfg.postgresDBHost
	tenCfg.Host.Debug.EnableDebugNamespace = cliCfg.isDebugNamespaceEnabled
	tenCfg.Host.Enclave.RPCAddresses = enclaveAddresses
	tenCfg.Host.L1.WebsocketURL = cliCfg.l1WebsocketURL
	tenCfg.Host.L1.L1BeaconUrl = cliCfg.l1BeaconUrl
	tenCfg.Host.L1.L1BlobArchiveUrl = cliCfg.l1BlobArchiveUrl
	tenCfg.Host.P2P.BindAddress = fmt.Sprintf("%s:%d", cliCfg.hostP2PHost, cliCfg.hostP2PPort)
	tenCfg.Host.P2P.IsDisabled = cliCfg.isInboundP2PDisabled
	tenCfg.Host.RPC.HTTPPort = uint64(cliCfg.hostHTTPPort)
	tenCfg.Host.RPC.WSPort = uint64(cliCfg.hostWSPort)
	tenCfg.Host.Log.Level = cliCfg.logLevel

	tenCfg.Enclave.DB.UseInMemory = false                                     // these nodes always use a persistent DB
	tenCfg.Enclave.DB.EdgelessDBHost = cliCfg.nodeName + "-edgelessdb-" + "0" // will be dynamically set for HA
	tenCfg.Enclave.Debug.EnableDebugNamespace = cliCfg.isDebugNamespaceEnabled
	tenCfg.Enclave.EnableAttestation = cliCfg.isSGXEnabled
	tenCfg.Enclave.RPC.BindAddress = fmt.Sprintf("0.0.0.0:%d", cliCfg.enclaveWSPort)
	tenCfg.Enclave.Log.Level = cliCfg.logLevel

	// the sequencer does not store the executed transactions
	// todo - once we replace this launcher we'll configure this flag explicitly via an environment variable
	if nodeType == common.Sequencer {
		tenCfg.Enclave.StoreExecutedTransactions = false
	}

	return tenCfg
}
