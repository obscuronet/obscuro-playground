package container

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/naoina/toml"

	"github.com/obscuronet/go-obscuro/go/config"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// HostConfigToml is the structure that a host's .toml config is parsed into.
type HostConfigToml struct {
	IsGenesis                 bool
	NodeType                  string
	HasClientRPCHTTP          bool
	ClientRPCPortHTTP         uint
	HasClientRPCWebsockets    bool
	ClientRPCPortWS           uint
	ClientRPCHost             string
	EnclaveRPCAddress         string
	P2PBindAddress            string
	P2PPublicAddress          string
	L1NodeHost                string
	L1NodeWebsocketPort       uint
	EnclaveRPCTimeout         int
	L1RPCTimeout              int
	P2PConnectionTimeout      int
	ManagementContractAddress string
	MessageBusAddress         string
	LogLevel                  int
	LogPath                   string
	PrivateKeyString          string
	L1ChainID                 int64
	ObscuroChainID            int64
	ProfilerEnabled           bool
	L1StartHash               string
	SequencerID               string
	MetricsEnabled            bool
	MetricsHTTPPort           uint
	UseInMemoryDB             bool
	LevelDBPath               string
	DebugNamespaceEnabled     bool
	BatchInterval             string
	RollupInterval            string
}

// ParseConfig returns a config.HostInputConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() (*config.HostInputConfig, error) {
	cfg := config.DefaultHostParsedConfig()
	flagUsageMap := getFlagUsageMap()

	configPath := flag.String(configName, "", flagUsageMap[configName])
	isGenesis := flag.Bool(isGenesisName, cfg.IsGenesis, flagUsageMap[isGenesisName])
	nodeTypeStr := flag.String(nodeTypeName, cfg.NodeType.String(), flagUsageMap[nodeTypeName])
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPName, cfg.ClientRPCPortHTTP, flagUsageMap[clientRPCPortHTTPName])
	clientRPCPortWS := flag.Uint64(clientRPCPortWSName, cfg.ClientRPCPortWS, flagUsageMap[clientRPCPortWSName])
	clientRPCHost := flag.String(clientRPCHostName, cfg.ClientRPCHost, flagUsageMap[clientRPCHostName])
	enclaveRPCAddress := flag.String(enclaveRPCAddressName, cfg.EnclaveRPCAddress, flagUsageMap[enclaveRPCAddressName])
	p2pBindAddress := flag.String(p2pBindAddressName, cfg.P2PBindAddress, flagUsageMap[p2pBindAddressName])
	p2pPublicAddress := flag.String(p2pPublicAddressName, cfg.P2PPublicAddress, flagUsageMap[p2pPublicAddressName])
	l1NodeHost := flag.String(l1NodeHostName, cfg.L1NodeHost, flagUsageMap[l1NodeHostName])
	l1NodePort := flag.Uint64(l1NodePortName, uint64(cfg.L1NodeWebsocketPort), flagUsageMap[l1NodePortName])
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, uint64(cfg.EnclaveRPCTimeout.Seconds()), flagUsageMap[enclaveRPCTimeoutSecsName])
	l1RPCTimeoutSecs := flag.Uint64(l1RPCTimeoutSecsName, uint64(cfg.L1RPCTimeout.Seconds()), flagUsageMap[l1RPCTimeoutSecsName])
	p2pConnectionTimeoutSecs := flag.Uint64(p2pConnectionTimeoutSecsName, uint64(cfg.P2PConnectionTimeout.Seconds()), flagUsageMap[p2pConnectionTimeoutSecsName])
	managementContractAddress := flag.String(managementContractAddrName, cfg.ManagementContractAddress.Hex(), flagUsageMap[managementContractAddrName])
	messageBusContractAddress := flag.String(messageBusContractAddrName, cfg.MessageBusAddress.Hex(), flagUsageMap[messageBusContractAddrName])
	logLevel := flag.Int(logLevelName, cfg.LogLevel, flagUsageMap[logLevelName])
	logPath := flag.String(logPathName, cfg.LogPath, flagUsageMap[logPathName])
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, flagUsageMap[l1ChainIDName])
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, flagUsageMap[obscuroChainIDName])
	privateKeyStr := flag.String(privateKeyName, cfg.PrivateKeyString, flagUsageMap[privateKeyName])
	profilerEnabled := flag.Bool(profilerEnabledName, cfg.ProfilerEnabled, flagUsageMap[profilerEnabledName])
	l1StartHash := flag.String(l1StartHashName, cfg.L1StartHash.Hex(), flagUsageMap[l1StartHashName])
	sequencerID := flag.String(sequencerIDName, cfg.SequencerID.Hex(), flagUsageMap[sequencerIDName])
	metricsEnabled := flag.Bool(metricsEnabledName, cfg.MetricsEnabled, flagUsageMap[metricsEnabledName])
	metricsHTPPPort := flag.Uint(metricsHTTPPortName, cfg.MetricsHTTPPort, flagUsageMap[metricsHTTPPortName])
	useInMemoryDB := flag.Bool(useInMemoryDBName, cfg.UseInMemoryDB, flagUsageMap[useInMemoryDBName])
	levelDBPath := flag.String(levelDBPathName, cfg.LevelDBPath, flagUsageMap[levelDBPathName])
	debugNamespaceEnabled := flag.Bool(debugNamespaceEnabledName, cfg.DebugNamespaceEnabled, flagUsageMap[debugNamespaceEnabledName])
	batchInterval := flag.String(batchIntervalName, cfg.BatchInterval.String(), flagUsageMap[batchIntervalName])
	rollupInterval := flag.String(rollupIntervalName, cfg.RollupInterval.String(), flagUsageMap[rollupIntervalName])

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	nodeType, err := common.ToNodeType(*nodeTypeStr)
	if err != nil {
		return &config.HostInputConfig{}, fmt.Errorf("unrecognised node type '%s'", *nodeTypeStr)
	}

	cfg.IsGenesis = *isGenesis
	cfg.NodeType = nodeType
	cfg.HasClientRPCHTTP = true
	cfg.ClientRPCPortHTTP = *clientRPCPortHTTP
	cfg.HasClientRPCWebsockets = true
	cfg.ClientRPCPortWS = *clientRPCPortWS
	cfg.ClientRPCHost = *clientRPCHost
	cfg.EnclaveRPCAddress = *enclaveRPCAddress
	cfg.P2PBindAddress = *p2pBindAddress
	cfg.P2PPublicAddress = *p2pPublicAddress
	cfg.L1NodeHost = *l1NodeHost
	cfg.L1NodeWebsocketPort = uint(*l1NodePort)
	cfg.EnclaveRPCTimeout = time.Duration(*enclaveRPCTimeoutSecs) * time.Second
	cfg.L1RPCTimeout = time.Duration(*l1RPCTimeoutSecs) * time.Second
	cfg.P2PConnectionTimeout = time.Duration(*p2pConnectionTimeoutSecs) * time.Second
	cfg.ManagementContractAddress = gethcommon.HexToAddress(*managementContractAddress)
	cfg.MessageBusAddress = gethcommon.HexToAddress(*messageBusContractAddress)
	cfg.PrivateKeyString = *privateKeyStr
	cfg.LogLevel = *logLevel
	cfg.LogPath = *logPath
	cfg.L1ChainID = *l1ChainID
	cfg.ObscuroChainID = *obscuroChainID
	cfg.ProfilerEnabled = *profilerEnabled
	cfg.L1StartHash = gethcommon.HexToHash(*l1StartHash)
	cfg.SequencerID = gethcommon.HexToAddress(*sequencerID)
	cfg.MetricsEnabled = *metricsEnabled
	cfg.MetricsHTTPPort = *metricsHTPPPort
	cfg.UseInMemoryDB = *useInMemoryDB
	cfg.LevelDBPath = *levelDBPath
	cfg.DebugNamespaceEnabled = *debugNamespaceEnabled
	cfg.BatchInterval, err = time.ParseDuration(*batchInterval)
	if err != nil {
		return nil, err
	}
	cfg.RollupInterval, err = time.ParseDuration(*rollupInterval)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Parses the config from the .toml file at configPath.
func fileBasedConfig(configPath string) (*config.HostInputConfig, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var tomlConfig HostConfigToml
	err = toml.Unmarshal(bytes, &tomlConfig)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	nodeType, err := common.ToNodeType(tomlConfig.NodeType)
	if err != nil {
		return &config.HostInputConfig{}, fmt.Errorf("unrecognised node type '%s'", tomlConfig.NodeType)
	}

	batchInterval, rollupInterval := 1*time.Second, 5*time.Second
	if interval, err := time.ParseDuration(tomlConfig.BatchInterval); err == nil {
		batchInterval = interval
	}
	if interval, err := time.ParseDuration(tomlConfig.RollupInterval); err == nil {
		rollupInterval = interval
	}

	return &config.HostInputConfig{
		IsGenesis:                 tomlConfig.IsGenesis,
		NodeType:                  nodeType,
		HasClientRPCHTTP:          tomlConfig.HasClientRPCHTTP,
		ClientRPCPortHTTP:         uint64(tomlConfig.ClientRPCPortHTTP),
		HasClientRPCWebsockets:    tomlConfig.HasClientRPCWebsockets,
		ClientRPCPortWS:           uint64(tomlConfig.ClientRPCPortWS),
		ClientRPCHost:             tomlConfig.ClientRPCHost,
		EnclaveRPCAddress:         tomlConfig.EnclaveRPCAddress,
		P2PBindAddress:            tomlConfig.P2PBindAddress,
		P2PPublicAddress:          tomlConfig.P2PPublicAddress,
		L1NodeHost:                tomlConfig.L1NodeHost,
		L1NodeWebsocketPort:       tomlConfig.L1NodeWebsocketPort,
		EnclaveRPCTimeout:         time.Duration(tomlConfig.EnclaveRPCTimeout) * time.Second,
		L1RPCTimeout:              time.Duration(tomlConfig.L1RPCTimeout) * time.Second,
		P2PConnectionTimeout:      time.Duration(tomlConfig.P2PConnectionTimeout) * time.Second,
		ManagementContractAddress: gethcommon.HexToAddress(tomlConfig.ManagementContractAddress),
		MessageBusAddress:         gethcommon.HexToAddress(tomlConfig.MessageBusAddress),
		LogLevel:                  tomlConfig.LogLevel,
		LogPath:                   tomlConfig.LogPath,
		PrivateKeyString:          tomlConfig.PrivateKeyString,
		L1ChainID:                 tomlConfig.L1ChainID,
		ObscuroChainID:            tomlConfig.ObscuroChainID,
		ProfilerEnabled:           tomlConfig.ProfilerEnabled,
		L1StartHash:               gethcommon.HexToHash(tomlConfig.L1StartHash),
		SequencerID:               gethcommon.HexToAddress(tomlConfig.SequencerID),
		MetricsEnabled:            tomlConfig.MetricsEnabled,
		MetricsHTTPPort:           tomlConfig.MetricsHTTPPort,
		UseInMemoryDB:             tomlConfig.UseInMemoryDB,
		LevelDBPath:               tomlConfig.LevelDBPath,
		BatchInterval:             batchInterval,
		RollupInterval:            rollupInterval,
	}, nil
}
