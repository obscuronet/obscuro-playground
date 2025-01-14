package faucet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/datagenerator"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/tools/faucet/container"
	"github.com/ten-protocol/go-ten/tools/faucet/faucet"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "faucet",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	contractDeployerPrivateKeyHex = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
	testLogs                      = "../.build/faucet/"
)

func TestFaucet(t *testing.T) {
	t.Skip("Skipping because it is too flaky")

	startPort := integration.TestPorts.TestFaucetPort
	createObscuroNetwork(t, startPort)
	// This sleep is required to ensure the initial rollup exists, and thus contract deployer can check its balance.
	time.Sleep(2 * time.Second)

	faucetConfig := &faucet.Config{
		Host:              "localhost",
		HTTPPort:          startPort + integration.DefaultHostRPCHTTPOffset,
		PK:                "0x" + contractDeployerPrivateKeyHex,
		JWTSecret:         "This_is_secret",
		ChainID:           big.NewInt(integration.TenChainID),
		ServerPort:        integration.TestPorts.TestFaucetHTTPPort,
		DefaultFundAmount: new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)),
	}
	faucetContainer, err := container.NewFaucetContainerFromConfig(faucetConfig)
	assert.NoError(t, err)

	err = faucetContainer.Start()
	defer func(faucetContainer *container.FaucetContainer) {
		err := faucetContainer.Stop()
		if err != nil {
			fmt.Printf("Could not stop faucet %s", err.Error())
		}
	}(faucetContainer)
	assert.NoError(t, err)

	initialFaucetBal, err := getFaucetBalance(faucetConfig.ServerPort)
	require.NoError(t, err)
	require.NotZero(t, initialFaucetBal)

	rndWallet := datagenerator.RandomWallet(integration.TenChainID)
	err = fundWallet(faucetConfig.ServerPort, rndWallet)
	require.NoError(t, err)

	obsClient, err := obsclient.DialWithAuth(fmt.Sprintf("http://%s:%d", network.Localhost, startPort+integration.DefaultHostRPCHTTPOffset), rndWallet, testlog.Logger())
	require.NoError(t, err)

	currentBalance, err := obsClient.BalanceAt(context.Background(), nil)
	require.NoError(t, err)

	if currentBalance.Cmp(big.NewInt(0)) <= 0 {
		t.Fatalf("Unexpected balance, got: %d, expected > 0", currentBalance.Int64())
	}

	endFaucetBal, err := getFaucetBalance(faucetConfig.ServerPort)
	require.NoError(t, err)
	assert.NotZero(t, endFaucetBal)
	// faucet balance should have decreased
	assert.Less(t, endFaucetBal.Cmp(initialFaucetBal), 0)
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) {
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

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
}

func fundWallet(port int, w wallet.Wallet) error {
	url := fmt.Sprintf("http://localhost:%d/auth/fund/eth", port)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"address":"%s"}`, w.Address()))

	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.xDOI1Cc30Zuj7VYKiRTqB2VntEKpZ5SkJW1heSsvzFw")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func getFaucetBalance(port int) (*big.Int, error) {
	url := fmt.Sprintf("http://localhost:%d/balance", port)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp struct {
		Balance string `json:"balance"`
	}
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	bal, success := new(big.Int).SetString(resp.Balance, 10)
	if !success {
		return nil, fmt.Errorf("failed to parse balance - %s", resp.Balance)
	}

	return bal, nil
}
