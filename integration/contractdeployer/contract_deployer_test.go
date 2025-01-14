package contractdeployer

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"
	contractdeployer "github.com/ten-protocol/go-ten/tools/hardhatdeployer"

	testcommon "github.com/ten-protocol/go-ten/integration/common"
)

const (
	contractDeployerPrivateKeyHex = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8" // Used only in tests.
	latestBlock                   = "latest"
	emptyCode                     = "0x"
	erc20ParamOne                 = "Hocus"
	erc20ParamTwo                 = "Hoc"
	erc20ParamThree               = "1000000000000000000"
	testLogs                      = "../.build/noderunner/"
	receiptTimeout                = 30 * time.Second // The time to wait for a receipt for a transaction.
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "contractdeployer",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

func TestCanDeployLayer2ERC20Contract(t *testing.T) {
	startPort := integration.TestPorts.TestCanDeployLayer2ERC20ContractPort
	hostWSPort := startPort + integration.DefaultHostRPCWSOffset
	creatTenNetwork(t, startPort)
	// This sleep is required to ensure the initial rollup exists, and thus contract deployer can check its balance.
	time.Sleep(2 * time.Second)

	config := &contractdeployer.Config{
		NodeHost:          network.Localhost,
		NodePort:          uint(hostWSPort),
		IsL1Deployment:    false,
		PrivateKey:        contractDeployerPrivateKeyHex,
		ChainID:           big.NewInt(integration.TenChainID),
		ContractName:      contractdeployer.Layer2Erc20Contract,
		ConstructorParams: []string{erc20ParamOne, erc20ParamTwo, erc20ParamThree},
	}

	contractAddr, err := contractdeployer.Deploy(config, testlog.Logger())
	if err != nil {
		panic(err)
	}

	contractDeployerWallet := wallet.NewInMemoryWalletFromConfig(contractDeployerPrivateKeyHex, integration.TenChainID, testlog.Logger())
	contractDeployerClient := getClient(hostWSPort, contractDeployerWallet)

	var deployedCode string
	err = contractDeployerClient.Call(&deployedCode, rpc.GetCode, contractAddr, latestBlock)
	if err != nil {
		panic(err)
	}

	if deployedCode == emptyCode {
		t.Fatal("contract was deployed but could not get code")
	}
}

func TestFaucetSendsFundsOnlyIfNeeded(t *testing.T) {
	startPort := integration.TestPorts.TestFaucetSendsFundsOnlyIfNeededPort
	hostWSPort := startPort + integration.DefaultHostRPCWSOffset
	creatTenNetwork(t, startPort)

	faucetWallet := wallet.NewInMemoryWalletFromConfig(testcommon.TestnetPrefundedPK, integration.TenChainID, testlog.Logger())
	faucetClient := getClient(hostWSPort, faucetWallet)

	contractDeployerWallet := wallet.NewInMemoryWalletFromConfig(contractDeployerPrivateKeyHex, integration.TenChainID, testlog.Logger())
	// We send more than enough to the contract deployer, to make sure prefunding won't be needed.
	excessivePrealloc := big.NewInt(contractdeployer.Prealloc * 3)
	testcommon.PrefundWallets(context.Background(), faucetWallet, obsclient.NewAuthObsClient(faucetClient), 0, []wallet.Wallet{contractDeployerWallet}, excessivePrealloc, receiptTimeout)

	// We check the faucet's balance before and after the deployment. Since the contract deployer has already been sent
	// sufficient funds, the faucet should have been to dispense any more, leaving its balance unchanged.
	var faucetInitialBalance string
	err := faucetClient.Call(&faucetInitialBalance, tenrpc.ERPCGetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	config := &contractdeployer.Config{
		NodeHost:          network.Localhost,
		NodePort:          uint(startPort + integration.DefaultHostRPCWSOffset),
		IsL1Deployment:    false,
		PrivateKey:        contractDeployerPrivateKeyHex,
		ChainID:           big.NewInt(integration.TenChainID),
		ContractName:      contractdeployer.Layer2Erc20Contract,
		ConstructorParams: []string{erc20ParamOne, erc20ParamTwo, erc20ParamThree},
	}

	_, err = contractdeployer.Deploy(config, testlog.Logger())
	if err != nil {
		panic(err)
	}

	var faucetBalanceAfterDeploy string
	// We create a new faucet client because deploying the contract will have overwritten the faucet's viewing key on the node.
	faucetClient = getClient(hostWSPort, faucetWallet)
	err = faucetClient.Call(&faucetBalanceAfterDeploy, tenrpc.ERPCGetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	if faucetInitialBalance != faucetBalanceAfterDeploy {
		t.Fatal("contract deployment allocated extra funds to contract deployer, despite sufficient funds")
	}
}

// Creates a single-node Obscuro network for testing.
func creatTenNetwork(t *testing.T, startPort int) {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.TenChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 2 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
		WithPrefunding:   true,
	}
	simStats := stats.NewStats(simParams.NumberOfNodes)
	tenNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(tenNetwork.TearDown)
	_, err := tenNetwork.Create(&simParams, simStats)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}

// Returns a viewing-key client with a registered viewing key.
func getClient(hostWSPort int, wallet wallet.Wallet) *rpc.EncRPCClient {
	viewingKey, err := viewingkey.GenerateViewingKeyForWallet(wallet)
	if err != nil {
		panic(err)
	}
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s:%d", network.Localhost, hostWSPort), viewingKey, testlog.Logger())
	if err != nil {
		panic(err)
	}
	return client
}
