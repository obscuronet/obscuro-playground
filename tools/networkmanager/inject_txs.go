package networkmanager

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"

	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/simulation"
)

func InjectTransactions(cfg Config, args []string, logger gethlog.Logger) {
	ctx := context.Background()
	println("Connecting to L1 node...")
	l1Client, err := ethadapter.NewEthClient(cfg.l1NodeHost, cfg.l1NodeWebsocketPort, cfg.l1RPCTimeout, common.HexToAddress("0x0"), logger)
	if err != nil {
		panic(fmt.Sprintf("could not create L1 client. Cause: %s", err))
	}
	println("Connecting to Obscuro node...")
	l2Client, err := rpc.NewNetworkClient(cfg.obscuroClientAddress)
	obscuroClient := obsclient.NewObsClient(l2Client)
	if err != nil {
		panic(err)
	}

	// We store the block at which we start injecting transactions.
	startBlock := l1Client.FetchHeadBlock()

	simStats := stats.NewStats(0)
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.mgmtContractAddress, logger)
	erc20ContractLib := erc20contractlib.NewERC20ContractLib(&cfg.mgmtContractAddress, &cfg.erc20ContractAddress)
	avgBlockDuration := time.Second

	wallets := createWallets(cfg, l1Client, l2Client, logger)
	walletClients := createWalletRPCClients(wallets, cfg.obscuroClientAddress, logger)

	rpcHandles := &network.RPCHandles{
		EthClients:     []ethadapter.EthClient{l1Client},
		ObscuroClients: []*obsclient.ObsClient{obscuroClient},
		AuthObsClients: walletClients,
	}

	txInjector := simulation.NewTransactionInjector(
		avgBlockDuration,
		simStats,
		rpcHandles,
		wallets,
		&cfg.mgmtContractAddress,
		mgmtContractLib,
		erc20ContractLib,
		parseNumOfTxs(args),
	)

	println("Injecting transactions into network...")
	txInjector.Start()
	time.Sleep(5 * avgBlockDuration) // We sleep to allow transactions to propagate.

	println(fmt.Sprintf("Stopped injecting transactions into network.\n"+
		"Attempted to inject %d L1 transactions, %d L2 transfer transactions, and %d L2 withdrawal transactions.",
		len(txInjector.TxTracker.L1Transactions), len(txInjector.TxTracker.TransferL2Transactions), len(txInjector.TxTracker.WithdrawalL2Transactions),
	))

	checkDepositsSuccessful(txInjector, l1Client, simStats, erc20ContractLib, mgmtContractLib, startBlock)
	checkL2TxsSuccessful(ctx, rpcHandles, txInjector)

	os.Exit(0)
}

// createWalletRPCClients creates map of wallet address to list of wallet clients (of length 1 because we have 1 node)
func createWalletRPCClients(wallets *params.SimWallets, obscuroNodeAddr string, logger gethlog.Logger) map[string][]*obsclient.AuthObsClient {
	clients := make(map[string][]*obsclient.AuthObsClient)

	for _, w := range wallets.SimObsWallets {
		vk, err := rpc.GenerateAndSignViewingKey(w)
		if err != nil {
			panic(err)
		}
		client, err := rpc.NewEncNetworkClient(obscuroNodeAddr, vk, logger)
		if err != nil {
			panic(err)
		}
		authClient := obsclient.NewAuthObsClient(client)

		clients[w.Address().String()] = []*obsclient.AuthObsClient{authClient}
	}
	for _, t := range wallets.Tokens {
		w := t.L2Owner
		vk, err := rpc.GenerateAndSignViewingKey(w)
		if err != nil {
			panic(err)
		}
		client, err := rpc.NewEncNetworkClient(obscuroNodeAddr, vk, logger)
		if err != nil {
			panic(err)
		}
		authClient := obsclient.NewAuthObsClient(client)

		clients[w.Address().String()] = []*obsclient.AuthObsClient{authClient}
	}

	return clients
}

func createWallets(nmConfig Config, l1Client ethadapter.EthClient, l2Client rpc.Client, logger gethlog.Logger) *params.SimWallets {
	wallets := params.NewSimWallets(len(nmConfig.privateKeys), 0, nmConfig.l1ChainID, nmConfig.obscuroChainID)

	// We override the autogenerated Ethereum wallets with ones using the provided private keys.
	wallets.SimEthWallets = make([]wallet.Wallet, len(nmConfig.privateKeys))
	for idx, privateKeyString := range nmConfig.privateKeys {
		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			panic(fmt.Errorf("could not recover private key from hex. Cause: %w", err))
		}
		l1Wallet := wallet.NewInMemoryWalletFromPK(big.NewInt(nmConfig.l1ChainID), privateKey, logger)
		wallets.SimEthWallets[idx] = l1Wallet
	}

	// We update the L1 and L2 wallet nonces.
	for _, l1Wallet := range wallets.AllEthWallets() {
		nonce, err := l1Client.Nonce(l1Wallet.Address())
		if err != nil {
			panic(fmt.Errorf("could not set L1 wallet nonce. Cause: %w", err))
		}
		l1Wallet.SetNonce(nonce)
	}
	for _, l2Wallet := range wallets.AllObsWallets() {
		var nonce uint64
		err := l2Client.Call(&nonce, rpc.GetTransactionCount, l2Wallet.Address())
		if err != nil {
			panic(fmt.Errorf("could not set L2 wallet nonce. Cause: %w", err))
		}
		l2Wallet.SetNonce(nonce)
	}

	// We set the ERC20 contract for the tokens.
	for _, token := range wallets.Tokens {
		token.L1ContractAddress = &nmConfig.erc20ContractAddress
	}

	return wallets
}

// Extracts the number of transactions to inject from the command-line arguments.
func parseNumOfTxs(args []string) int {
	if len(args) != 1 {
		panic(fmt.Errorf("expected one argument to %s command, got %d", injectTxsName, len(args)))
	}
	numOfTxs, err := strconv.Atoi(args[0])
	if err != nil {
		panic(fmt.Errorf("could not parse number of transactions to inject. Cause: %w", err))
	}
	return numOfTxs
}

func checkDepositsSuccessful(txInjector *simulation.TransactionInjector, l1Client ethadapter.EthClient, stats *stats.Stats, erc20ContractLib erc20contractlib.ERC20ContractLib, mgmtContractLib mgmtcontractlib.MgmtContractLib, startBlock *types.Block) {
	currentBlock := l1Client.FetchHeadBlock()
	dummySim := simulation.Simulation{
		Stats: stats,
		Params: &params.SimParams{
			ERC20ContractLib: erc20ContractLib,
			MgmtContractLib:  mgmtContractLib,
		},
	}
	deposits, _, _, _ := simulation.ExtractDataFromEthereumChain(startBlock, currentBlock, l1Client, &dummySim, 0) //nolint:dogsled

	if len(deposits) != len(txInjector.TxTracker.L1Transactions) {
		println(fmt.Sprintf("Injected %d deposits into the L1 but %d were missing.",
			len(deposits), len(txInjector.TxTracker.L1Transactions)))
	} else {
		println(fmt.Sprintf("Successfully injected %d deposits into the L1.", len(deposits)))
	}
}

func checkL2TxsSuccessful(ctx context.Context, rpcHandles *network.RPCHandles, txInjector *simulation.TransactionInjector) {
	injectedTransfers := len(txInjector.TxTracker.TransferL2Transactions)
	injectedWithdrawals := len(txInjector.TxTracker.WithdrawalL2Transactions)
	notFoundTransfers, notFoundWithdrawals := simulation.FindNotIncludedL2Txs(ctx, 0, rpcHandles, txInjector)

	if notFoundTransfers != 0 {
		println(fmt.Sprintf("Injected %d transfers into the L2 but %d were missing.", injectedTransfers, notFoundTransfers))
	} else {
		println(fmt.Sprintf("Successfully injected %d transfers into the L1.", injectedTransfers))
	}

	if notFoundWithdrawals != 0 {
		println(fmt.Sprintf("Injected %d withdrawals into the L2 but %d were missing.", injectedWithdrawals, notFoundWithdrawals))
	} else {
		println(fmt.Sprintf("Successfully injected %d withdrawals into the L1.", injectedWithdrawals))
	}
}
