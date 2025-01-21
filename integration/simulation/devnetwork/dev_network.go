package devnetwork

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/tools/walletextension"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/simulation/network"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

const (
	// these ports were picked arbitrarily, if we want plan to use these tests on CI we need to use ports in the constants.go file
	_gwHTTPPort = 11180
	_gwWSPort   = 11181
)

var _defaultFaucetAmount = big.NewInt(750_000_000_000_000)

// InMemDevNetwork is a local dev network (L1 and L2) - the TEN nodes are in-memory in a single go process, the L1 nodes are a docker geth network
//
// It can play the role of node operators and network admins to reproduce complex scenarios around nodes joining/leaving/failing.
//
// It also implements networktest.NetworkConnector to allow us to run the same NetworkTests against it that we can run against Testnets.
type InMemDevNetwork struct {
	logger gethlog.Logger

	// todo (@matt) - replace this with a struct for accs/contracts that are controlled by network admins
	// 	(don't pollute with "user" wallets, they will be controlled by the individual network test runners)
	networkWallets *params.SimWallets

	l1Network L1Network

	// When TEN network has been initialised on the L1 network, this will be populated
	// - if reconnecting to an existing network it needs to be populated when initialising this object
	// - if it is nil when `Start()` is called then TEN contracts will be deployed on the L1
	l1SetupData *params.L1TenData

	tenConfig           *TenConfig
	tenSequencer        *InMemNodeOperator
	tenValidators       []*InMemNodeOperator
	tenGatewayContainer *walletextension.Container

	faucet     userwallet.User
	faucetLock sync.Mutex
}

func (s *InMemDevNetwork) GetGatewayURL() (string, error) {
	if !s.tenConfig.TenGatewayEnabled {
		return "", fmt.Errorf("TEN gateway not enabled")
	}
	return fmt.Sprintf("http://localhost:%d", _gwHTTPPort), nil
}

func (s *InMemDevNetwork) GetGatewayWSURL() (string, error) {
	if !s.tenConfig.TenGatewayEnabled {
		return "", fmt.Errorf("TEN gateway not enabled")
	}
	return fmt.Sprintf("ws://localhost:%d", _gwWSPort), nil
}

func (s *InMemDevNetwork) GetMCOwnerWallet() (wallet.Wallet, error) {
	return s.networkWallets.MCOwnerWallet, nil
}

func (s *InMemDevNetwork) ChainID() int64 {
	return integration.TenChainID
}

func (s *InMemDevNetwork) FaucetWallet() wallet.Wallet {
	return s.networkWallets.L2FaucetWallet
}

func (s *InMemDevNetwork) AllocateFaucetFunds(ctx context.Context, account gethcommon.Address) error {
	// ensure only one test account is getting faucet funds at a time, faucet client isn't thread-safe
	s.faucetLock.Lock()
	defer s.faucetLock.Unlock()

	txHash, err := s.faucet.SendFunds(ctx, account, _defaultFaucetAmount)
	if err != nil {
		return err
	}

	receipt, err := s.faucet.AwaitReceipt(ctx, txHash)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("faucet transaction receipt status not successful - %v", receipt.Status)
	}
	return nil
}

func (s *InMemDevNetwork) SequencerRPCAddress() string {
	seq := s.GetSequencerNode()
	return seq.HostRPCWSAddress()
}

func (s *InMemDevNetwork) ValidatorRPCAddress(idx int) string {
	val := s.GetValidatorNode(idx)
	return val.HostRPCWSAddress()
}

// GetL1Client returns the first client we have for our local L1 network
// todo (@matt) - this allows tests some basic L1 verification but in future this will need support more manipulation of L1 nodes,
//
//	(to allow us to simulate various scenarios where L1 is unavailable, under attack, etc.)
func (s *InMemDevNetwork) GetL1Client() (ethadapter.EthClient, error) {
	return s.l1Network.GetClient(0), nil
}

func (s *InMemDevNetwork) GetSequencerNode() networktest.NodeOperator {
	return s.tenSequencer
}

func (s *InMemDevNetwork) GetValidatorNode(i int) networktest.NodeOperator {
	return s.tenValidators[i]
}

func (s *InMemDevNetwork) NumValidators() int {
	return len(s.tenValidators)
}

func (s *InMemDevNetwork) Start() {
	if s.logger == nil {
		s.logger = testlog.Logger()
	}
	fmt.Println("Starting L1 network")
	s.l1Network.Prepare()
	if s.l1SetupData == nil {
		// this is a new network, deploy the contracts to the L1
		fmt.Println("Deploying TEN contracts to L1")
		s.deployTenNetworkContracts()
		fmt.Printf("L1 Port - %d\n", integration.TestPorts.NetworkTestsPort)
	}
	fmt.Println("Starting TEN nodes")
	s.startNodes()

	// sleep to allow the nodes to start
	time.Sleep(30 * time.Second) // it takes a while for the secret exchanges etc. to sort themselves out

	seqClient, err := obsclient.Dial(s.SequencerRPCAddress())
	if err != nil {
		panic(err)
	}
	h, _ := seqClient.Health()
	if len(h.Enclaves) == 0 {
		panic("no enclaves available to promote on sequencer")
	}
	for _, e := range h.Enclaves {
		err = network.PermissionTenSequencerEnclave(s.networkWallets.MCOwnerWallet, s.l1Network.GetClient(0), s.l1SetupData.MgmtContractAddress, e.EnclaveID)
		if err != nil {
			panic(err)
		}
	}

	if s.tenConfig.TenGatewayEnabled {
		s.startTenGateway()
	}
}

func (s *InMemDevNetwork) GetGatewayClient() (ethadapter.EthClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemDevNetwork) startNodes() {
	if s.tenSequencer == nil {
		// initialise node operators (sequencer is node 0, validators are 1..N)
		s.tenSequencer = NewInMemNodeOperator(0, s.tenConfig, common.Sequencer, s.l1SetupData, s.l1Network.GetClient(0), s.networkWallets.NodeWallets[0], s.logger)
		for i := 1; i <= s.tenConfig.InitNumValidators; i++ {
			l1Client := s.l1Network.GetClient(i % s.l1Network.NumNodes())
			s.tenValidators = append(s.tenValidators, NewInMemNodeOperator(i, s.tenConfig, common.Validator, s.l1SetupData, l1Client, s.networkWallets.NodeWallets[i], s.logger))
		}
	}

	go func() {
		err := s.tenSequencer.Start()
		if err != nil {
			panic(err)
		}
	}()
	for _, v := range s.tenValidators {
		go func(v networktest.NodeOperator) {
			err := v.Start()
			if err != nil {
				panic(err)
			}
		}(v)
	}
	s.faucet = userwallet.NewUserWallet(s.networkWallets.L2FaucetWallet, s.SequencerRPCAddress(), s.logger)
}

func (s *InMemDevNetwork) startTenGateway() {
	validator := s.GetValidatorNode(0)
	validatorHTTP := validator.HostRPCHTTPAddress()
	// remove http:// prefix for the gateway config
	validatorHTTP = validatorHTTP[len("http://"):]
	validatorWS := validator.HostRPCWSAddress()
	// remove ws:// prefix for the gateway config
	validatorWS = validatorWS[len("ws://"):]
	cfg := wecommon.Config{
		WalletExtensionHost:     "127.0.0.1",
		WalletExtensionPortHTTP: _gwHTTPPort,
		WalletExtensionPortWS:   _gwWSPort,
		NodeRPCHTTPAddress:      validatorHTTP,
		NodeRPCWebsocketAddress: validatorWS,
		LogPath:                 "sys_out",
		VerboseFlag:             false,
		DBType:                  "sqlite",
		TenChainID:              integration.TenChainID,
	}
	tenGWContainer := walletextension.NewContainerFromConfig(cfg, s.logger)
	go func() {
		fmt.Println("Starting TEN Gateway, HTTP Port:", _gwHTTPPort, "WS Port:", _gwWSPort)
		err := tenGWContainer.Start()
		if err != nil {
			s.logger.Error("failed to start TEN gateway", "err", err)
			panic(err)
		}
		s.tenGatewayContainer = tenGWContainer
	}()
}

func (s *InMemDevNetwork) CleanUp() {
	for _, v := range s.tenValidators {
		go func(v networktest.NodeOperator) {
			err := v.Stop()
			if err != nil {
				fmt.Println("failed to stop validator", err.Error())
			}
		}(v)
	}
	go func() {
		err := s.tenSequencer.Stop()
		if err != nil {
			fmt.Println("failed to stop sequencer", err.Error())
		}
	}()
	go s.l1Network.CleanUp()
	if s.tenGatewayContainer != nil {
		go func() {
			err := s.tenGatewayContainer.Stop()
			if err != nil {
				fmt.Println("failed to stop TEN gateway", err.Error())
			}
		}()
	}

	s.logger.Info("Waiting for servers to stop.")
	time.Sleep(3 * time.Second)
}

func (s *InMemDevNetwork) deployTenNetworkContracts() {
	client := s.l1Network.GetClient(0)
	// note: we don't currently deploy ERC20s here, don't want to waste gas on sepolia
	l1SetupData, err := network.DeployTenNetworkContracts(client, s.networkWallets, false)
	if err != nil {
		panic(err)
	}
	s.l1SetupData = l1SetupData
}
