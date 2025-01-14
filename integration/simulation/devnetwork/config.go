package devnetwork

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/wallet"
	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

// L1Config tells network admin how to setup the L1 network
type L1Config struct {
	PortStart          int
	WebsocketPortStart int
	NumNodes           int
	AvgBlockDuration   time.Duration
}

type TenConfigOption func(*TenConfig) // option pattern - typically used as overrides to DefaultTenConfig

// TenConfig describes the L2 network configuration we want to spin up
type TenConfig struct {
	PortStart          int
	InitNumValidators  int
	BatchInterval      time.Duration
	RollupInterval     time.Duration
	CrossChainInterval time.Duration
	NumNodes           int
	TenGatewayEnabled  bool
	NumSeqEnclaves     int
	L1BeaconPort       int
	L1BlockTime        time.Duration
	DeployerPK         string
}

func DefaultTenConfig() *TenConfig {
	return &TenConfig{
		PortStart:          integration.TestPorts.NetworkTestsPort,
		NumNodes:           4,
		InitNumValidators:  3,
		BatchInterval:      1 * time.Second,
		RollupInterval:     10 * time.Second,
		CrossChainInterval: 11 * time.Second,
		TenGatewayEnabled:  false,
		NumSeqEnclaves:     1, // increase for HA simulation
		L1BeaconPort:       integration.TestPorts.NetworkTestsPort + integration.DefaultPrysmGatewayPortOffset,
		DeployerPK:         "",
	}
}

func LocalDevNetwork(tenConfigOpts ...TenConfigOption) *InMemDevNetwork {
	tenConfig := DefaultTenConfig()
	for _, opt := range tenConfigOpts {
		opt(tenConfig)
	}

	// 1 wallet per node
	nodeOpL1Wallets := params.NewSimWallets(0, tenConfig.NumNodes, integration.EthereumChainID, integration.TenChainID)

	if tenConfig.DeployerPK != "" {
		privKey, err := crypto.HexToECDSA(tenConfig.DeployerPK)
		if err != nil {
			panic("could not initialise deployer private key")
		}
		nodeOpL1Wallets.MCOwnerWallet = wallet.NewInMemoryWalletFromPK(big.NewInt(integration.EthereumChainID), privKey, testlog.Logger())
	}

	l1Config := &L1Config{
		PortStart:        integration.TestPorts.NetworkTestsPort,
		NumNodes:         tenConfig.NumNodes, // we'll have 1 L1 node per L2 node
		AvgBlockDuration: 2 * time.Second,
	}
	l1Network := NewGethNetwork(nodeOpL1Wallets, l1Config)

	return NewInMemDevNetwork(tenConfig, l1Network, nodeOpL1Wallets)
}

// NewInMemDevNetwork provides an off-the-shelf default config for a sim network
// tenConfig - the requirements of the L2 network we are spinning up
// l1Network - the L1 network we are running the L2 network on
// nodeOpL1Wallets - the funded wallets for the node operators on the L1 network (expecting 1 per node)
func NewInMemDevNetwork(tenConfig *TenConfig, l1Network L1Network, nodeOpL1Wallets *params.SimWallets) *InMemDevNetwork {
	// update tenConfig references to be consistent with the L1 setup
	tenConfig.L1BlockTime = l1Network.GetBlockTime()

	return &InMemDevNetwork{
		networkWallets: nodeOpL1Wallets,
		l1Network:      l1Network,
		tenConfig:      tenConfig,
		faucetLock:     sync.Mutex{},
	}
}

// LiveL1DevNetwork provides a local obscuro network running on a live L1
// Caller should provide a wallet per node and ideally an RPC URL per node (may not be necessary but can avoid conflicts, e.g. Infura seems to require an API key per connection)
func LiveL1DevNetwork(seqWallet wallet.Wallet, validatorWallets []wallet.Wallet, rpcURLs []string) *InMemDevNetwork {
	// setup the host and deployer wallets to be the prefunded wallets

	// create the L2 faucet wallet
	l2FaucetPrivKey, err := crypto.HexToECDSA(testcommon.TestnetPrefundedPK)
	if err != nil {
		panic("could not initialise L2 faucet private key")
	}
	l2FaucetWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.TenChainID), l2FaucetPrivKey, testlog.Logger())
	networkWallets := &params.SimWallets{
		MCOwnerWallet:  seqWallet,
		NodeWallets:    append([]wallet.Wallet{seqWallet}, validatorWallets...),
		L2FaucetWallet: l2FaucetWallet,
		Tokens:         map[testcommon.ERC20]*params.SimToken{},
	}

	l1Network := &liveL1Network{
		deployWallet:     seqWallet, // use the same wallet for deploying the contracts
		seqWallet:        seqWallet,
		validatorWallets: validatorWallets,
		rpcURLs:          rpcURLs,
	}

	return &InMemDevNetwork{
		logger:         testlog.Logger(),
		networkWallets: networkWallets,
		l1Network:      l1Network,
		tenConfig:      DefaultTenConfig(),
	}
}

func WithGateway() TenConfigOption {
	return func(tc *TenConfig) {
		tc.TenGatewayEnabled = true
	}
}

func WithHASequencer() TenConfigOption {
	return func(tc *TenConfig) {
		tc.NumSeqEnclaves = 2
	}
}

// WithPredictableDeployer - will use a known private key for the deployer instead of generating a random one.
// Useful for being able to run administrative functions from hardhat when debugging on the deployed network.
func WithPredictableDeployer() TenConfigOption {
	return func(tc *TenConfig) {
		tc.DeployerPK = "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
	}
}
