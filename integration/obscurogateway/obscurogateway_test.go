package faucet

import (
	"context"
	"encoding/json"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/obsclient"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/ethereummock"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
	"github.com/obscuronet/go-obscuro/tools/walletextension/container"
	"github.com/stretchr/testify/assert"
)

func init() { //nolint:gochecknoinits
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "obscurogateway",
		TestSubtype: "test",
		LogLevel:    log.LvlInfo,
	})
}

const (
	testLogs = "../.build/obscurogateway/"
)

func TestObscuroGateway(t *testing.T) {
	//t.Skip("skipping until Tudor's DB changes simplify the enclave logic")
	startPort := integration.StartPortObscuroGatewayUnitTest
	wallets := createObscuroNetwork(t, startPort)

	obscuroGatewayConf := config.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: startPort + integration.DefaultObscuroGatewayHTTPPortOffset,
		WalletExtensionPortWS:   startPort + integration.DefaultObscuroGatewayWSPortOffset,
		NodeRPCHTTPAddress:      fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCHTTPOffset),
		NodeRPCWebsocketAddress: fmt.Sprintf("127.0.0.1:%d", startPort+integration.DefaultHostRPCWSOffset),
		LogPath:                 "sys_out",
		VerboseFlag:             false,
	}

	obscuroGwContainer := container.NewWalletExtensionContainerFromConfig(obscuroGatewayConf, testlog.Logger())
	go func() {
		err := obscuroGwContainer.Start()
		if err != nil {
			fmt.Printf("error stopping WE - %s", err)
		}
	}()

	// wait for the msg bus contract to be deployed
	time.Sleep(5 * time.Second)

	// make sure the server is ready to receive requests
	// TODO Implement health endpoint
	serverAddress := fmt.Sprintf("http://%s:%d", obscuroGatewayConf.WalletExtensionHost, obscuroGatewayConf.WalletExtensionPortHTTP)

	w := wallets.L2FaucetWallet

	vk, err := viewingkey.GenerateViewingKeyForWallet(w)
	assert.Nil(t, err)
	client, err := rpc.NewEncNetworkClient(fmt.Sprintf("ws://%s", obscuroGatewayConf.NodeRPCWebsocketAddress), vk, testlog.Logger())
	assert.Nil(t, err)
	authClient := obsclient.NewAuthObsClient(client)

	balance, err := authClient.BalanceAt(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotEqual(t, big.NewInt(0), balance)

	txHash := transferRandomAddr(t, authClient, w)

	// Issue tests
	jsonRPCRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getStorageAt",
		"params": []interface{}{
			w.Address().String(),
			nil,
			nil,
		},
		"id": "1",
	}
	bodyReq, err := json.Marshal(jsonRPCRequest)
	assert.NoError(t, err)

	statusCode, body, err := httputil.PostDataJSON(serverAddress, bodyReq)
	assert.NoError(t, err)
	require.Equal(t, 200, statusCode)

	type request struct {
		Result []common.PublicTxData `json:"result"`
	}

	pubDataReq := request{}
	err = json.Unmarshal(body, &pubDataReq)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(pubDataReq.Result))
	assert.Equal(t, w.Address(), pubDataReq.Result[0].SenderAddress)
	assert.Equal(t, txHash.Hex(), pubDataReq.Result[0].TransactionHash.Hex())

	// Gracefully shutdown
	err = obscuroGwContainer.Stop()
	assert.NoError(t, err)
}

func transferRandomAddr(t *testing.T, authClient *obsclient.AuthObsClient, w wallet.Wallet) common.TxHash {
	ctx := context.Background()
	toAddr := datagenerator.RandomAddress()
	nonce, err := authClient.NonceAt(ctx, nil)
	assert.Nil(t, err)

	w.SetNonce(nonce)
	estimatedTx := authClient.EstimateGasAndGasPrice(&types.LegacyTx{
		Nonce:    w.GetNonceAndIncrement(),
		To:       &toAddr,
		Value:    big.NewInt(100),
		Gas:      uint64(1_000_000),
		GasPrice: gethcommon.Big1,
	})
	assert.Nil(t, err)

	fmt.Println("Transferring from:", w.Address(), " to:", toAddr)

	signedTx, err := w.SignTransaction(estimatedTx)
	assert.Nil(t, err)

	err = authClient.SendTransaction(ctx, signedTx)
	assert.Nil(t, err)

	fmt.Printf("Created Tx: %s \n", signedTx.Hash().Hex())
	fmt.Printf("Checking for tx receipt for %s \n", signedTx.Hash())
	var receipt *types.Receipt
	for start := time.Now(); time.Since(start) < time.Minute; time.Sleep(time.Second) {
		receipt, err = authClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			break
		}
		//
		// Currently when a receipt is not available the obscuro node is returning nil instead of err ethereum.NotFound
		// once that's fixed this commented block should be removed
		//if !errors.Is(err, ethereum.NotFound) {
		//	t.Fatal(err)
		//}
		if receipt != nil && receipt.Status == 1 {
			break
		}
		fmt.Printf("no tx receipt after %s - %s\n", time.Since(start), err)
	}

	if receipt == nil {
		t.Fatalf("Did not mine the transaction after %s seconds  - receipt: %+v", 30*time.Second, receipt)
	}
	if receipt.Status == 0 {
		t.Fatalf("Tx Failed")
	}
	fmt.Println("Successfully minted the transaction - ", signedTx.Hash())
	return signedTx.Hash()
}

// Creates a single-node Obscuro network for testing.
func createObscuroNetwork(t *testing.T, startPort int) *params.SimWallets {
	// Create the Obscuro network.
	numberOfNodes := 1
	wallets := params.NewSimWallets(1, numberOfNodes, integration.EthereumChainID, integration.ObscuroChainID)
	simParams := params.SimParams{
		NumberOfNodes:    numberOfNodes,
		AvgBlockDuration: 1 * time.Second,
		MgmtContractLib:  ethereummock.NewMgmtContractLibMock(),
		ERC20ContractLib: ethereummock.NewERC20ContractLibMock(),
		Wallets:          wallets,
		StartPort:        startPort,
	}

	obscuroNetwork := network.NewNetworkOfSocketNodes(wallets)
	t.Cleanup(obscuroNetwork.TearDown)
	_, err := obscuroNetwork.Create(&simParams, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to create test Obscuro network. Cause: %s", err))
	}
	return wallets
}
