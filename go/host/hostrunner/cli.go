package hostrunner

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/naoina/toml"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/common"
)

// HostConfigToml is the structure that a host's .toml config is parsed into.
type HostConfigToml struct {
	ID                     string
	IsGenesis              bool
	GossipRoundDuration    int
	HasClientRPCHTTP       bool
	ClientRPCPortHTTP      uint
	HasClientRPCWebsockets bool
	ClientRPCPortWS        uint
	ClientRPCHost          string
	ClientRPCTimeout       int
	EnclaveRPCAddress      string
	EnclaveRPCTimeout      int
	P2PBindAddress         string
	P2PPublicAddress       string
	L1NodeHost             string
	L1NodeWebsocketPort    uint
	L1ConnectionTimeout    int
	RollupContractAddress  string
	LogLevel               string
	LogPath                string
	PrivateKeyString       string
	L1ChainID              int64
	ObscuroChainID         int64
	ProfilerEnabled        bool
}

// ParseConfig returns a config.HostConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() config.HostConfig {
	cfg := config.DefaultHostConfig()

	configPath := flag.String(configName, "", configUsage)
	nodeID := flag.String(nodeIDName, cfg.ID.Hex(), nodeIDUsage)
	isGenesis := flag.Bool(isGenesisName, cfg.IsGenesis, isGenesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(cfg.GossipRoundDuration), gossipRoundNanosUsage)
	clientRPCPortHTTP := flag.Uint64(clientRPCPortHTTPName, cfg.ClientRPCPortHTTP, clientRPCPortHTTPUsage)
	clientRPCPortWS := flag.Uint64(clientRPCPortWSName, cfg.ClientRPCPortWS, clientRPCPortWSUsage)
	clientRPCHost := flag.String(clientRPCHostName, cfg.ClientRPCHost, clientRPCHostUsage)
	clientRPCTimeoutSecs := flag.Uint64(clientRPCTimeoutSecsName, uint64(cfg.ClientRPCTimeout.Seconds()), clientRPCTimeoutSecsUsage)
	enclaveRPCAddress := flag.String(enclaveRPCAddressName, cfg.EnclaveRPCAddress, enclaveRPCAddressUsage)
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, uint64(cfg.EnclaveRPCTimeout.Seconds()), enclaveRPCTimeoutSecsUsage)
	p2pBindAddress := flag.String(p2pBindAddressName, cfg.P2PBindAddress, p2pBindAddressUsage)
	p2pPublicAddress := flag.String(p2pPublicAddressName, cfg.P2PPublicAddress, p2pPublicAddressUsage)
	l1NodeHost := flag.String(l1NodeHostName, cfg.L1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(cfg.L1NodeWebsocketPort), l1NodePortUsage)
	l1ConnectionTimeoutSecs := flag.Uint64(l1ConnectionTimeoutSecsName, uint64(cfg.L1ConnectionTimeout.Seconds()), l1ConnectionTimeoutSecsUsage)
	rollupContractAddress := flag.String(rollupContractAddrName, cfg.RollupContractAddress.Hex(), rollupContractAddrUsage)
	logLevel := flag.String(logLevelName, cfg.LogLevel, logLevelUsage)
	logPath := flag.String(logPathName, cfg.LogPath, logPathUsage)
	l1ChainID := flag.Int64(l1ChainIDName, cfg.L1ChainID, l1ChainIDUsage)
	obscuroChainID := flag.Int64(obscuroChainIDName, cfg.ObscuroChainID, obscuroChainIDUsage)
	privateKeyStr := flag.String(privateKeyName, cfg.PrivateKeyString, privateKeyUsage)
	profilerEnabled := flag.Bool(profilerEnabledName, cfg.ProfilerEnabled, profilerEnabledUsage)

	flag.Parse()

	if *configPath != "" {
		return fileBasedConfig(*configPath)
	}

	cfg.ID = common.HexToAddress(*nodeID)
	cfg.IsGenesis = *isGenesis
	cfg.GossipRoundDuration = time.Duration(*gossipRoundNanos)
	cfg.HasClientRPCHTTP = true
	cfg.ClientRPCPortHTTP = *clientRPCPortHTTP
	cfg.HasClientRPCWebsockets = true
	cfg.ClientRPCPortWS = *clientRPCPortWS
	cfg.ClientRPCHost = *clientRPCHost
	cfg.ClientRPCTimeout = time.Duration(*enclaveRPCTimeoutSecs) * time.Second
	cfg.EnclaveRPCAddress = *enclaveRPCAddress
	cfg.EnclaveRPCTimeout = time.Duration(*clientRPCTimeoutSecs) * time.Second
	cfg.P2PBindAddress = *p2pBindAddress
	cfg.P2PPublicAddress = *p2pPublicAddress
	cfg.L1NodeHost = *l1NodeHost
	cfg.L1NodeWebsocketPort = uint(*l1NodePort)
	cfg.L1ConnectionTimeout = time.Duration(*l1ConnectionTimeoutSecs) * time.Second
	cfg.RollupContractAddress = common.HexToAddress(*rollupContractAddress)
	cfg.PrivateKeyString = *privateKeyStr
	cfg.LogLevel = *logLevel
	cfg.LogPath = *logPath
	cfg.L1ChainID = *l1ChainID
	cfg.ObscuroChainID = *obscuroChainID
	cfg.ProfilerEnabled = *profilerEnabled

	return cfg
}

// Parses the config from the .toml file at configPath.
func fileBasedConfig(configPath string) config.HostConfig {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	var tomlConfig HostConfigToml
	err = toml.Unmarshal(bytes, &tomlConfig)
	if err != nil {
		panic(fmt.Sprintf("could not read config file at %s. Cause: %s", configPath, err))
	}

	return config.HostConfig{
		ID:                     common.HexToAddress(tomlConfig.ID),
		IsGenesis:              tomlConfig.IsGenesis,
		GossipRoundDuration:    time.Duration(tomlConfig.GossipRoundDuration),
		HasClientRPCHTTP:       tomlConfig.HasClientRPCHTTP,
		ClientRPCPortHTTP:      uint64(tomlConfig.ClientRPCPortHTTP),
		HasClientRPCWebsockets: tomlConfig.HasClientRPCWebsockets,
		ClientRPCPortWS:        uint64(tomlConfig.ClientRPCPortWS),
		ClientRPCHost:          tomlConfig.ClientRPCHost,
		ClientRPCTimeout:       time.Duration(tomlConfig.ClientRPCTimeout) * time.Second,
		EnclaveRPCAddress:      tomlConfig.EnclaveRPCAddress,
		EnclaveRPCTimeout:      time.Duration(tomlConfig.EnclaveRPCTimeout) * time.Second,
		P2PBindAddress:         tomlConfig.P2PBindAddress,
		P2PPublicAddress:       tomlConfig.P2PPublicAddress,
		L1NodeHost:             tomlConfig.L1NodeHost,
		L1NodeWebsocketPort:    tomlConfig.L1NodeWebsocketPort,
		L1ConnectionTimeout:    time.Duration(tomlConfig.L1ConnectionTimeout) * time.Second,
		RollupContractAddress:  common.HexToAddress(tomlConfig.RollupContractAddress),
		LogLevel:               tomlConfig.LogLevel,
		LogPath:                tomlConfig.LogPath,
		PrivateKeyString:       tomlConfig.PrivateKeyString,
		L1ChainID:              tomlConfig.L1ChainID,
		ObscuroChainID:         tomlConfig.ObscuroChainID,
		ProfilerEnabled:        tomlConfig.ProfilerEnabled,
	}
}
