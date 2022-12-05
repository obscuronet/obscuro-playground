package contractdeployer

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"

	"github.com/obscuronet/go-obscuro/tools/contractdeployer"

	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const (
	contractDeployerPrivateKeyHex = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
	latestBlock                   = "latest"
	emptyCode                     = "0x"
	erc20ParamOne                 = "Hocus"
	erc20ParamTwo                 = "Hoc"
	erc20ParamThree               = "1000000000000000000"
	testLogs                      = "../.build/noderunner/"
)

var (
	config = &contractdeployer.Config{
		NodeHost:          network.Localhost,
		NodePort:          integration.StartPortContractDeployerTest + network.DefaultHostRPCWSOffset,
		IsL1Deployment:    false,
		PrivateKey:        contractDeployerPrivateKeyHex,
		ChainID:           big.NewInt(integration.ObscuroChainID),
		ContractName:      contractdeployer.Layer2Erc20Contract,
		ConstructorParams: []string{erc20ParamOne, erc20ParamTwo, erc20ParamThree},
	}
	nodeAddress = fmt.Sprintf("ws://%s:%d", config.NodeHost, config.NodePort)
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "noderunner",
		TestSubtype: "test",
	})
}

func TestCanDeployLayer2ERC20Contract(t *testing.T) {
	createObscuroNetwork(t)
	// This sleep is required to ensure the initial rollup exists, and thus contract deployer can check its balance.
	time.Sleep(2 * time.Second)
	contractAddr, err := contractdeployer.Deploy(config, testlog.Logger())
	if err != nil {
		panic(err)
	}

	contractDeployerWallet := getWallet(contractDeployerPrivateKeyHex)
	contractDeployerClient := getClient(contractDeployerWallet)

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
	createObscuroNetwork(t)

	faucetWallet := getWallet(rollupchain.FaucetPrivateKeyHex)
	faucetClient := getClient(faucetWallet)

	contractDeployerWallet := getWallet(contractDeployerPrivateKeyHex)
	// We send more than enough to the contract deployer, to make sure prefunding won't be needed.
	excessivePrealloc := big.NewInt(contractdeployer.Prealloc * 3)
	testcommon.PrefundWallets(context.Background(), faucetWallet, obsclient.NewAuthObsClient(faucetClient), 0, []wallet.Wallet{contractDeployerWallet}, excessivePrealloc)

	// We check the faucet's balance before and after the deployment. Since the contract deployer has already been sent
	// sufficient funds, the faucet should have been to dispense any more, leaving its balance unchanged.
	var faucetInitialBalance string
	err := faucetClient.Call(&faucetInitialBalance, rpc.GetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	_, err = contractdeployer.Deploy(config, testlog.Logger())
	if err != nil {
		panic(err)
	}

	var faucetBalanceAfterDeploy string
	// We create a new faucet client because deploying the contract will have overwritten the faucet's viewing key on the node.
	faucetClient = getClient(faucetWallet)
	err = faucetClient.Call(&faucetBalanceAfterDeploy, rpc.GetBalance, faucetWallet.Address().Hex(), latestBlock)
	if err != nil {
		panic(err)
	}

	if faucetInitialBalance != faucetBalanceAfterDeploy {
		t.Fatal("contract deployment allocated extra funds to contract deployer, despite sufficient funds")
	}
}

func getWallet(privateKeyHex string) wallet.Wallet {
	faucetPrivKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic("could not initialise faucet private key")
	}
	faucetWallet := wallet.NewInMemoryWalletFromPK(config.ChainID, faucetPrivKey, testlog.Logger())
	return faucetWallet
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T) {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        integration.StartPortContractDeployerTest,
	}
	simStats := stats.NewStats(simParams.NumberOfNodes)
	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, simStats)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}

// Returns a viewing-key client with a registered viewing key.
func getClient(wallet wallet.Wallet) *rpc.EncRPCClient {
	viewingKey, err := rpc.GenerateAndSignViewingKey(wallet)
	if err != nil {
		panic(err)
	}
	client, err := rpc.NewEncNetworkClient(nodeAddress, viewingKey, testlog.Logger())
	if err != nil {
		panic(err)
	}
	return client
}
