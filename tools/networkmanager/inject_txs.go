package networkmanager

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"

	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
)

func InjectTransactions(cfg Config) {
	hostConfig := config.HostConfig{
		L1NodeHost:          cfg.l1NodeHost,
		L1NodeWebsocketPort: cfg.l1NodeWebsocketPort,
		L1ConnectionTimeout: cfg.l1ConnectionTimeout,
	}
	l1Client, err := ethclient.NewEthClient(hostConfig)
	if err != nil {
		panic(fmt.Sprintf("could not create L1 client. Cause: %s", err))
	}
	l2Client := obscuroclient.NewClient(cfg.obscuroClientAddress)

	txInjector := simulation.NewTransactionInjector(
		1*time.Second,
		stats.NewStats(1),
		[]ethclient.EthClient{l1Client},
		createWallets(cfg, l1Client, l2Client),
		&cfg.mgmtContractAddress,
		[]obscuroclient.Client{l2Client},
		mgmtcontractlib.NewMgmtContractLib(&cfg.mgmtContractAddress),
		erc20contractlib.NewERC20ContractLib(&cfg.mgmtContractAddress, &cfg.erc20ContractAddress),
	)

	// We listen for interrupts, to log statistics before exiting.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		for range interruptChan {
			println(fmt.Sprintf(
				"Stopped injecting transactions into network\nInjected %d L1 transactions, %d L2 transfer transactions, and %d L2 withdrawal transactions.",
				len(txInjector.Counter.L1Transactions), len(txInjector.Counter.TransferL2Transactions), len(txInjector.Counter.WithdrawalL2Transactions),
			))
			os.Exit(0)
		}
	}()

	println("Injecting transactions into network...")
	txInjector.Start()
}

func createWallets(nmConfig Config, l1Client ethclient.EthClient, l2Client obscuroclient.Client) *params.SimWallets {
	wallets := params.NewSimWallets(len(nmConfig.privateKeys), 0, nmConfig.l1ChainID, nmConfig.obscuroChainID)

	// We override the autogenerated Ethereum wallets with ones using the provided private keys.
	wallets.SimEthWallets = make([]wallet.Wallet, len(nmConfig.privateKeys))
	for idx, privateKeyString := range nmConfig.privateKeys {
		privateKey, err := crypto.HexToECDSA(privateKeyString)
		if err != nil {
			panic(fmt.Errorf("could not recover private key from hex. Cause: %w", err))
		}
		l1Wallet := wallet.NewInMemoryWalletFromPK(big.NewInt(nmConfig.l1ChainID), privateKey)
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
		err := l2Client.Call(&nonce, obscuroclient.RPCNonce, l2Wallet.Address())
		if err != nil {
			panic(fmt.Errorf("could not set L2 wallet nonce. Cause: %w", err))
		}
		l2Wallet.SetNonce(nonce)
	}

	return wallets
}
