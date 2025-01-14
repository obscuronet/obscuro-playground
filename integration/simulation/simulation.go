package simulation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	gethparams "github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/PublicCallbacksTest"
	"github.com/ten-protocol/go-ten/contracts/generated/TransactionPostProcessor"
	"github.com/ten-protocol/go-ten/contracts/generated/ZenBase"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/erc20contract"
	"github.com/ten-protocol/go-ten/integration/ethereummock"
	"github.com/ten-protocol/go-ten/integration/simulation/network"

	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"

	gethcommon "github.com/ethereum/go-ethereum/common"
	testcommon "github.com/ten-protocol/go-ten/integration/common"
)

var initialBalance = common.ValueInWei(big.NewInt(5000))

// Simulation represents all the data required to inject transactions on a network
type Simulation struct {
	RPCHandles       *network.RPCHandles
	AvgBlockDuration uint64
	TxInjector       *TransactionInjector
	SimulationTime   time.Duration
	Stats            *stats.Stats
	Params           *params.SimParams
	ZenBaseAddress   gethcommon.Address
	LogChannels      map[string][]chan types.Log // Maps an owner to the channels on which they receive logs for each client.
	Subscriptions    []ethereum.Subscription     // A slice of all created event subscriptions.
	ctx              context.Context
}

// Start executes the simulation given all the Params. Injects transactions.
func (s *Simulation) Start() {
	testlog.Logger().Info(fmt.Sprintf("Genesis block: b_%s.", ethereummock.MockGenesisBlock.Hash()))
	s.ctx = context.Background() // use injected context for graceful shutdowns

	fmt.Printf("Waiting for TEN genesis on L1\n")
	s.waitForTenGenesisOnL1()

	// Arbitrary sleep to wait for RPC clients to get up and running
	// and for all l2 nodes to receive the genesis l2 batch
	// todo - instead of sleeping, it would be better to poll
	time.Sleep(10 * time.Second)

	cfg, err := s.RPCHandles.TenWalletRndClient(s.Params.Wallets.L2FaucetWallet).GetConfig()
	if err != nil {
		panic(err)
	}
	jsonCfg, err := json.Marshal(cfg)
	if err == nil {
		fmt.Printf("Config: %v\n", string(jsonCfg))
	}

	fmt.Printf("Funding the bridge to TEN\n")
	s.bridgeFundingToTen()

	fmt.Printf("Deploying ZenBase contract\n")
	s.deployTenZen() // Deploy the ZenBase contract

	fmt.Printf("Deploying PublicCallbacksTest contract\n")
	s.deployPublicCallbacksTest()

	fmt.Printf("Creating log subscriptions\n")
	s.trackLogs() // Create log subscriptions, to validate that they're working correctly later.

	fmt.Printf("Prefunding L2 wallets\n")
	s.prefundTenAccounts() // Prefund every L2 wallet

	fmt.Printf("Deploying TEN ERC20 contracts\n")
	s.deployTenERC20s() // Deploy the TEN HOC and POC ERC20 contracts

	fmt.Printf("Prefunding L1 wallets\n")
	s.prefundL1Accounts() // Prefund every L1 wallet

	fmt.Printf("Checking health status\n")
	s.checkHealthStatus() // Checks the nodes health status

	timer := time.Now()
	fmt.Printf("Starting injection\n")
	testlog.Logger().Info("Starting injection")
	go s.TxInjector.Start()

	// Allow for some time after tx injection was stopped so that the network can process all transactions, catch up
	// on missed batches, etc.

	// Wait for the simulation time
	time.Sleep(s.SimulationTime - s.Params.StoppingDelay)
	fmt.Printf("Stopping injection\n")
	testlog.Logger().Info("Stopping injection")

	s.TxInjector.Stop()

	time.Sleep(s.Params.StoppingDelay)

	fmt.Printf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime)
	testlog.Logger().Info(fmt.Sprintf("Ran simulation for %f secs, configured to run for: %s ... \n", time.Since(timer).Seconds(), s.SimulationTime))
}

func (s *Simulation) Stop() {
	// nothing to do for now
}

func (s *Simulation) waitForTenGenesisOnL1() {
	// grab an L1 client
	client := s.RPCHandles.EthClients[0]

	for {
		// spin through the L1 blocks periodically to see if the genesis rollup has arrived
		head, err := client.FetchHeadBlock()
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			panic(fmt.Errorf("could not fetch head block. Cause: %w", err))
		}
		if err == nil {
			for _, h := range client.BlocksBetween(ethereummock.MockGenesisBlock.Header(), head) {
				b, err := client.BlockByHash(h.Hash())
				if err != nil {
					panic(err)
				}
				for _, tx := range b.Transactions() {
					t := s.Params.MgmtContractLib.DecodeTx(tx)
					if t == nil {
						continue
					}
					if _, ok := t.(*common.L1RollupHashes); ok {
						// exit at the first TEN rollup we see
						return
					}
				}
			}
		}
		time.Sleep(s.Params.AvgBlockDuration)
		testlog.Logger().Trace("Waiting for the TEN genesis rollup...")
	}
}

func (s *Simulation) bridgeFundingToTen() {
	if s.Params.IsInMem {
		return
	}

	testlog.Logger().Info("Funding the bridge to TEN")

	destAddr := s.Params.L1TenData.MessageBusAddr
	value, _ := big.NewInt(0).SetString("7400000000000000000000000000000", 10)

	wallets := []wallet.Wallet{
		s.Params.Wallets.PrefundedEthWallets.Faucet,
		s.Params.Wallets.PrefundedEthWallets.HOC,
		s.Params.Wallets.PrefundedEthWallets.POC,
	}

	receivers := []gethcommon.Address{
		gethcommon.HexToAddress("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77"),
		gethcommon.HexToAddress("0x987E0a0692475bCc5F13D97E700bb43c1913EFfe"),
		gethcommon.HexToAddress("0xDEe530E22045939e6f6a0A593F829e35A140D3F1"),
	}

	busCtr, err := MessageBus.NewMessageBus(destAddr, s.RPCHandles.RndEthClient().EthClient())
	if err != nil {
		panic(err)
	}

	for idx, w := range wallets {
		opts, err := bind.NewKeyedTransactorWithChainID(w.PrivateKey(), w.ChainID())
		if err != nil {
			panic(err)
		}
		opts.Value = value

		_, err = busCtr.SendValueToL2(opts, receivers[idx], value)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(15 * time.Second)
	// todo - fix the wait group, for whatever reason it does not find a receipt...
	/*wg := sync.WaitGroup{}
	for _, tx := range transactions {
		wg.Add(1)
		transaction := tx
		go func() {
			defer wg.Done()
			err := testcommon.AwaitReceiptEth(s.ctx, s.RPCHandles.RndEthClient(), transaction.Hash(), 20*time.Second)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()*/
}

// We subscribe to logs on every client for every wallet.
func (s *Simulation) trackLogs() {
	// In-memory clients cannot handle subscriptions for now.
	if s.Params.IsInMem {
		return
	}

	testlog.Logger().Info("Subscribing to logs")

	for owner, clients := range s.RPCHandles.AuthObsClients {
		// There is a subscription, and corresponding log channel, per owner per client.
		s.LogChannels[owner] = []chan types.Log{}

		for _, client := range clients {
			channel := make(chan types.Log, 1000)

			// To exercise the filtering mechanism, we subscribe for HOC events only, ignoring POC events.
			hocFilter := common.FilterCriteria{
				Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + testcommon.HOCAddr)},
			}
			sub, err := client.SubscribeFilterLogsTEN(context.Background(), hocFilter, channel)
			if err != nil {
				panic(fmt.Errorf("subscription failed. Cause: %w", err))
			}
			s.Subscriptions = append(s.Subscriptions, sub)

			s.LogChannels[owner] = append(s.LogChannels[owner], channel)
		}
	}
}

// Prefunds the L2 wallets with `allocObsWallets` each.
func (s *Simulation) prefundTenAccounts() {
	testlog.Logger().Info("Prefunding L2 wallets")

	faucetWallet := s.Params.Wallets.L2FaucetWallet
	faucetClient := s.RPCHandles.TenWalletClient(faucetWallet.Address(), 0) // get sequencer, else errors on submission get swallowed
	// in memory test needs this to allow head batch to be set
	time.Sleep(5 * time.Second)
	nonce := NextNonce(s.ctx, s.RPCHandles, faucetWallet)

	// Give 1000 ether per account - ether is 1e18 so best convert it by code
	// as a lot of the hardcodes were giving way too little and choking the gas payments
	allocObsWallets := big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(gethparams.Ether))
	testcommon.PrefundWallets(s.ctx, faucetWallet, faucetClient, nonce, s.Params.Wallets.AllObsWallets(), allocObsWallets, s.Params.ReceiptTimeout)
}

func (s *Simulation) deployPublicCallbacksTest() {
	testlog.Logger().Info("Deploying PublicCallbacksTest contract")

	auth, err := bind.NewKeyedTransactorWithChainID(s.Params.Wallets.L2FaucetWallet.PrivateKey(), s.Params.Wallets.L2FaucetWallet.ChainID())
	if err != nil {
		panic(fmt.Errorf("failed to create transactor in order to bootstrap sim test: %w", err))
	}
	rpcClient := s.RPCHandles.TenWalletClient(s.Params.Wallets.L2FaucetWallet.Address(), 1)
	var cfg *common.TenNetworkInfo
	for cfg == nil || cfg.TransactionPostProcessorAddress.Cmp(gethcommon.Address{}) == 0 {
		cfg, err = rpcClient.GetConfig()
		if err != nil {
			s.TxInjector.logger.Info("failed to get config", log.ErrKey, err)
		}
		time.Sleep(2 * time.Second)
	}

	publicCallbacksAddress := cfg.PublicSystemContracts["PublicCallbacks"]
	if publicCallbacksAddress.Cmp(gethcommon.Address{}) == 0 {
		panic(fmt.Errorf("public callbacks address is not set"))
	}

	auth.Nonce = big.NewInt(0).SetUint64(NextNonce(s.ctx, s.RPCHandles, s.Params.Wallets.L2FaucetWallet))
	auth.GasPrice = big.NewInt(0).SetUint64(gethparams.InitialBaseFee)
	auth.Context = s.ctx
	auth.Value = big.NewInt(0).Mul(big.NewInt(1), big.NewInt(gethparams.Ether))

	_, tx, instance, err := PublicCallbacksTest.DeployPublicCallbacksTest(auth, rpcClient, publicCallbacksAddress)
	if err != nil {
		panic(fmt.Errorf("failed to deploy public callbacks test contract: %w", err))
	}

	receipt, err := bind.WaitMined(s.ctx, rpcClient, tx)
	if err != nil || receipt.Status != types.ReceiptStatusSuccessful {
		panic(fmt.Errorf("failed to deploy public callbacks test contract"))
	}

	success, err := instance.IsLastCallSuccess(&bind.CallOpts{Context: s.ctx, From: s.Params.Wallets.L2FaucetWallet.Address()})
	if err != nil {
		panic(fmt.Errorf("failed to check if last call was successful: %w", err))
	}
	if !success {
		panic(fmt.Errorf("last call was not successful"))
	}
}

func (s *Simulation) deployTenZen() {
	testlog.Logger().Info("Deploying ZenBase contract")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		auth, err := bind.NewKeyedTransactorWithChainID(s.Params.Wallets.L2FaucetWallet.PrivateKey(), s.Params.Wallets.L2FaucetWallet.ChainID())
		if err != nil {
			panic(fmt.Errorf("failed to create transactor in order to bootstrap sim test: %w", err))
		}

		// Node one, because random client might yield the no p2p node, which breaks the timings
		rpcClient := s.RPCHandles.TenWalletClient(s.Params.Wallets.L2FaucetWallet.Address(), 1)
		var cfg *common.TenNetworkInfo
		for cfg == nil || cfg.TransactionPostProcessorAddress.Cmp(gethcommon.Address{}) == 0 {
			cfg, err = rpcClient.GetConfig()
			if err != nil {
				s.TxInjector.logger.Info("failed to get config", log.ErrKey, err)
			}
			time.Sleep(2 * time.Second)
		}

		// Wait for balance with retry
		err = retry.Do(func() error {
			balance, err := rpcClient.BalanceAt(context.Background(), nil)
			if err != nil {
				return fmt.Errorf("failed to get balance: %w", err)
			}
			if balance.Cmp(big.NewInt(0)) <= 0 {
				return fmt.Errorf("waiting for positive balance")
			}
			return nil
		}, retry.NewTimeoutStrategy(1*time.Minute, 2*time.Second))
		if err != nil {
			panic(fmt.Errorf("failed to get positive balance after timeout: %w", err))
		}

		owner := s.Params.Wallets.L2FaucetWallet
		ownerRpc := s.RPCHandles.TenWalletClient(owner.Address(), 1)
		auth.GasPrice = big.NewInt(0).SetUint64(gethparams.InitialBaseFee)
		auth.Context = context.Background()
		auth.Nonce = big.NewInt(0).SetUint64(NextNonce(s.ctx, s.RPCHandles, owner))

		zenBaseAddress, signedTx, _, err := ZenBase.DeployZenBase(auth, ownerRpc, cfg.TransactionPostProcessorAddress) //, "ZenBase", "ZEN")
		if err != nil {
			panic(fmt.Errorf("failed to deploy zen base contract: %w", err))
		}
		if receipt, err := bind.WaitMined(s.ctx, s.RPCHandles.TenWalletRndClient(owner), signedTx); err != nil || receipt.Status != types.ReceiptStatusSuccessful {
			panic(fmt.Errorf("failed to deploy zen base contract"))
		}
		s.ZenBaseAddress = zenBaseAddress

		transactionPostProcessor, err := TransactionPostProcessor.NewTransactionPostProcessor(cfg.TransactionPostProcessorAddress, ownerRpc)
		if err != nil {
			panic(fmt.Errorf("failed to deploy transactions analyzer contract: %w", err))
		}

		auth.Nonce = big.NewInt(0).SetUint64(NextNonce(s.ctx, s.RPCHandles, owner))
		tx, err := transactionPostProcessor.AddOnBlockEndCallback(auth, zenBaseAddress)
		if err != nil {
			panic(fmt.Errorf("failed to add on block end callback: %w", err))
		}
		receipt, err := bind.WaitMined(s.ctx, ownerRpc, tx)
		if err != nil || receipt.Status != types.ReceiptStatusSuccessful {
			panic(fmt.Errorf("failed to add on block end callback"))
		}
	}()
	wg.Wait()
}

// This deploys an ERC20 contract on Ten, which is used for token arithmetic.
func (s *Simulation) deployTenERC20s() {
	testlog.Logger().Info("Deploying TEN ERC20 contracts")
	tokens := []testcommon.ERC20{testcommon.HOC, testcommon.POC}

	wg := sync.WaitGroup{}
	for _, token := range tokens {
		wg.Add(1)
		go func(token testcommon.ERC20) {
			defer wg.Done()
			owner := s.Params.Wallets.Tokens[token].L2Owner

			cfg, err := s.RPCHandles.TenWalletRndClient(owner).GetConfig()
			if err != nil {
				panic(err)
			}
			contractBytes := erc20contract.L2BytecodeWithDefaultSupply(string(token), cfg.L2MessageBusAddress)

			fmt.Printf("Deploy contract from: %s\n", owner.Address().Hex())
			deployContractTxData := types.DynamicFeeTx{
				Nonce:     NextNonce(s.ctx, s.RPCHandles, owner),
				Gas:       5_000_000,
				GasFeeCap: gethcommon.Big1, // This field is used to derive the gas price for dynamic fee transactions.
				Data:      contractBytes,
				GasTipCap: gethcommon.Big1,
			}

			deployContractTx := s.RPCHandles.TenWalletRndClient(owner).EstimateGasAndGasPrice(&deployContractTxData)
			signedTx, err := owner.SignTransaction(deployContractTx)
			if err != nil {
				panic(err)
			}

			rpc := s.RPCHandles.TenWalletClient(owner.Address(), 1)
			err = rpc.SendTransaction(s.ctx, signedTx)
			if err != nil {
				panic(fmt.Sprintf("ERC20 deployment transaction unsuccessful. Cause: %s", err))
			}

			err = testcommon.AwaitReceipt(s.ctx, rpc, signedTx.Hash(), s.Params.ReceiptTimeout)
			if err != nil {
				panic(fmt.Sprintf("ERC20 deployment transaction unsuccessful. Cause: %s", err))
			}
		}(token)
	}
	wg.Wait()
}

// Sends an amount from the faucet to each L1 account, to pay for transactions.
func (s *Simulation) prefundL1Accounts() {
	testlog.Logger().Info("Prefunding L1 wallets")

	for _, w := range s.Params.Wallets.SimEthWallets {
		ethClient := s.RPCHandles.RndEthClient()
		receiver := w.Address()
		tokenOwner := s.Params.Wallets.Tokens[testcommon.HOC].L1Owner
		ownerAddr := tokenOwner.Address()
		txData := &common.L1DepositTx{
			Amount:        initialBalance,
			To:            &receiver,
			TokenContract: s.Params.Wallets.Tokens[testcommon.HOC].L1ContractAddress,
			Sender:        &ownerAddr,
		}
		tx := s.Params.ERC20ContractLib.CreateDepositTx(txData)
		estimatedTx, err := ethClient.PrepareTransactionToSend(s.ctx, tx, tokenOwner.Address())
		if err != nil {
			// ignore txs that are not able to be estimated/execute
			testlog.Logger().Error("unable to estimate tx", log.ErrKey, err)
			tokenOwner.SetNonce(tokenOwner.GetNonce() - 1)
			continue
		}
		signedTx, err := tokenOwner.SignTransaction(estimatedTx)
		if err != nil {
			panic(err)
		}
		err = s.RPCHandles.RndEthClient().SendTransaction(signedTx)
		if err != nil {
			panic(err)
		}

		go s.TxInjector.TxTracker.trackL1Tx(txData)
	}
}

func (s *Simulation) checkHealthStatus() {
	testlog.Logger().Info("Checking health status")

	for _, client := range s.RPCHandles.TenClients {
		err := retry.Do(func() error {
			healthy, err := client.Health()
			if !healthy.OverallHealth || err != nil {
				return fmt.Errorf("client is not healthy: %w", err)
			}
			return nil
		}, retry.NewTimeoutStrategy(30*time.Second, 100*time.Millisecond))
		if err != nil {
			panic(err)
		}
	}
}

func NextNonce(ctx context.Context, clients *network.RPCHandles, w wallet.Wallet) uint64 {
	counter := 0

	// only returns the nonce when the previous transaction was recorded
	for {
		remoteNonce, err := clients.TenWalletRndClient(w).NonceAt(ctx, nil)
		if err != nil {
			panic(err)
		}
		localNonce := w.GetNonce()
		if remoteNonce == localNonce {
			return w.GetNonceAndIncrement()
		}
		if remoteNonce > localNonce {
			panic("remote nonce exceeds local nonce")
		}

		counter++
		if counter > nonceTimeoutMillis {
			panic(fmt.Sprintf("transaction injector could not retrieve nonce after thirty seconds for address %s. "+
				"Local nonce was %d, remote nonce was %d", w.Address().Hex(), localNonce, remoteNonce))
		}
		time.Sleep(time.Millisecond)
	}
}
