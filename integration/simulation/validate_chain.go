package simulation

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/ZenBase"

	testcommon "github.com/ten-protocol/go-ten/integration/common"
	"github.com/ten-protocol/go-ten/integration/ethereummock"

	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/ten-protocol/go-ten/integration/simulation/network"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
)

const (
	// The threshold number of transactions below which we consider the simulation to have failed. We generally expect far
	// more than this, but this is a sanity check to ensure the simulation doesn't stop after a single transaction of each
	// type, for example.
	txThreshold = 5
	// The maximum number of blocks an TEN node can fall behind
	maxBlockDelay = 5
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// After a simulation has run, check as much as possible that the outputs of the simulation are expected.
// For example, all injected transactions were processed correctly, the height of the rollup chain is a function of the total
// time of the simulation and the average block duration, that all TEN nodes are roughly in sync, etc
func checkNetworkValidity(t *testing.T, s *Simulation) {
	checkTransactionsInjected(t, s)
	l1MaxHeight := checkEthereumBlockchainValidity(t, s)
	checkTenBlockchainValidity(t, s, l1MaxHeight)
	checkReceivedLogs(t, s)
	checkTenscan(t, s)
	checkZenBaseMinting(t, s)
}

// Ensures that L1 and L2 txs were actually issued.
func checkTransactionsInjected(t *testing.T, s *Simulation) {
	if len(s.TxInjector.TxTracker.L1Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d L1 transactions. At least %d expected", len(s.TxInjector.TxTracker.L1Transactions), txThreshold)
	}
	if len(s.TxInjector.TxTracker.TransferL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d transfer L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.TransferL2Transactions), txThreshold)
	}
	// todo (@stefan) - reenable when old contract deployer phased out.
	/*if len(s.TxInjector.TxTracker.WithdrawalL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d withdrawal L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.WithdrawalL2Transactions), txThreshold)
	}*/
	if len(s.TxInjector.TxTracker.NativeValueTransferL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d transfer L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.NativeValueTransferL2Transactions), txThreshold)
	}
}

// checkEthereumBlockchainValidity: sanity check of the state of all L1 nodes
// - the chain has a minimum number of blocks
// - the chain height is similar across all ethereum nodes
// - there are no duplicate txs
// - efficiency - number of created blocks/height
// - no reorgs
func checkEthereumBlockchainValidity(t *testing.T, s *Simulation) uint64 {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	heights := make([]uint64, len(s.RPCHandles.EthClients))
	for i, node := range s.RPCHandles.EthClients {
		heights[i] = checkBlockchainOfEthereumNode(t, node, minHeight, s, i)
	}

	min, max := minMax(heights)
	if max-min > max/10 {
		t.Logf("Node Heights: %v", heights)
		t.Errorf("There is a problem with the mock Ethereum chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}

	return max
}

// checkTenBlockchainValidity - perform the following checks
// - minimum height - the chain has a minimum number of rollups
// - check height is similar
// - check no duplicate txs
// - check efficiency - no of created blocks/ height
// - check amount in the system
// - check withdrawals/deposits
func checkTenBlockchainValidity(t *testing.T, s *Simulation, maxL1Height uint64) {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	// process the blockchain of each node in parallel to minimize the difference between them since they are still running
	heights := make([]uint64, len(s.RPCHandles.TenClients))
	var wg sync.WaitGroup
	for idx := range s.RPCHandles.TenClients {
		wg.Add(1)
		go checkBlockchainOfTenNode(t, s.RPCHandles, minHeight, maxL1Height, s, &wg, heights, idx)
	}
	wg.Wait()
	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	// since there is one node that only listens to rollups it will be naturally behind.
	if max-min > max/3 {
		t.Errorf("There is a problem with the TEN chain. Nodes fell out of sync. Max height: %d. Min height: %d -> %+v", max, min, heights)
	}
}

// the cost of an empty rollup - adjust if the management contract changes. This is the rollup overhead.
const emptyRollupGas = 110_000

func checkCollectedL1Fees(_ *testing.T, node ethadapter.EthClient, s *Simulation, nodeIdx int, rollupReceipts types.Receipts) {
	costOfRollupsWithTransactions := big.NewInt(0)
	costOfEmptyRollups := big.NewInt(0)

	if s.Params.IsInMem {
		// not supported for in memory tests
		return
	}

	for _, receipt := range rollupReceipts {
		block, err := node.EthClient().BlockByHash(context.Background(), receipt.BlockHash)
		if err != nil {
			panic(err)
		}

		txCost := big.NewInt(0).Mul(block.BaseFee(), big.NewInt(0).SetUint64(receipt.GasUsed))
		// only calculate the fees collected for non-empty rollups, because the empty ones are subsidized
		if receipt.GasUsed > emptyRollupGas {
			costOfRollupsWithTransactions.Add(costOfRollupsWithTransactions, txCost)
		} else {
			costOfEmptyRollups.Add(costOfEmptyRollups, txCost)
		}
	}

	l2FeesWallet := s.Params.Wallets.L2FeesWallet
	obsClients := network.CreateAuthClients(s.RPCHandles.RPCClients, l2FeesWallet)
	_, err := obsClients[nodeIdx].BalanceAt(context.Background(), nil)
	if err != nil {
		panic(fmt.Errorf("failed getting balance for bridge transfer receiver. Cause: %w", err))
	}

	// if balance of collected fees is less than cost of published rollups fail
	// todo - reenable when gas payments are behaving themselves
	//if feeBalance.Cmp(costOfRollupsWithTransactions) == -1 {
	//	t.Errorf("Node %d: Sequencer has collected insufficient fees. Has: %d, needs: %d", nodeIdx, feeBalance, costOfRollupsWithTransactions)
	//}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethadapter.EthClient, minHeight uint64, s *Simulation, nodeIdx int) uint64 {
	head, err := node.FetchHeadBlock()
	if err != nil {
		t.Errorf("Node %d: Could not find head block. Cause: %s", nodeIdx, err)
	}
	height := head.Number.Uint64()

	if height < minHeight {
		t.Errorf("Node %d: There were only %d blocks mined. Expected at least: %d.", nodeIdx, height, minHeight)
	}

	deposits, rollups, _, blockCount, _, rollupReceipts := ExtractDataFromEthereumChain(ethereummock.MockGenesisBlock.Header(), head, node, s, nodeIdx)
	s.Stats.TotalL1Blocks = uint64(blockCount)

	checkCollectedL1Fees(t, node, s, nodeIdx, rollupReceipts)

	hashDups := findHashDups(deposits)
	if len(hashDups) > 0 {
		t.Errorf("Node %d: Found Deposit duplicates: %v", nodeIdx, hashDups)
	}

	rollupDups := findRollupDups(rollups)
	if len(rollupDups) > 0 {
		// todo @siliev - fix in memory rollups, lack of real client breaks the normal ask smart contract flow.
		if !s.Params.IsInMem {
			t.Errorf("Node %d: Found Rollup duplicates: %v", nodeIdx, rollupDups)
		}
	}

	/* todo (@stefan) - reenable while old contract deployer is phased out.
	if s.Stats.TotalDepositedAmount.Cmp(totalDeposited) != 0 {
		t.Errorf("Node %d: Deposit[%d] amounts don't match. Found %d , expected %d", nodeIdx, len(deposits), totalDeposited, s.Stats.TotalDepositedAmount)
	}

	if s.Stats.TotalDepositedAmount.Cmp(gethcommon.Big0) == 0 {
		t.Errorf("Node %d: No deposits", nodeIdx)
	} */

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().L2ID]
	reorgEfficiency := float64(reorgs) / float64(height)
	if reorgEfficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: The number of reorgs is too high: %d. ", nodeIdx, reorgs)
	}
	if !s.Params.IsInMem {
		checkRollups(t, s, nodeIdx, rollups)
	}

	return height
}

// this function only performs a very brief check.
// the ultimate check that everything works fine is that each node is able to respond to queries
// and has processed all batches correctly.
func checkRollups(t *testing.T, _ *Simulation, nodeIdx int, rollups []*common.ExtRollup) {
	if len(rollups) < 2 {
		t.Errorf("Node %d: Found less than two submitted rollups! Successful simulation should always produce more than 2", nodeIdx)
	}

	sort.Slice(rollups, func(i, j int) bool {
		// Ascending order sort.
		return rollups[i].Header.LastBatchSeqNo < rollups[j].Header.LastBatchSeqNo
	})

	for _, rollup := range rollups {
		// todo (@matt) verify the rollup was produced by a sequencer enclave from the whitelisted ID set

		if len(rollup.BatchPayloads) == 0 {
			t.Errorf("Node %d: No batches in rollup!", nodeIdx)
			continue
		}
	}
}

// ExtractDataFromEthereumChain returns the deposits, rollups, total amount deposited and length of the blockchain
// between the start block and the end block.
func ExtractDataFromEthereumChain(startBlock *types.Header, endBlock *types.Header, node ethadapter.EthClient, s *Simulation, nodeIdx int) ([]gethcommon.Hash, []*common.ExtRollup, *big.Int, int, uint64, types.Receipts) {
	deposits := make([]gethcommon.Hash, 0)
	rollups := make([]*common.ExtRollup, 0)
	rollupReceipts := make(types.Receipts, 0)
	totalDeposited := big.NewInt(0)

	blockchain := node.BlocksBetween(startBlock, endBlock)
	successfulDeposits := uint64(0)
	for _, header := range blockchain {
		block, err := node.BlockByHash(header.Hash())
		if err != nil {
			panic(err)
		}
		for _, tx := range block.Transactions() {
			t := s.Params.ERC20ContractLib.DecodeTx(tx)
			if t == nil {
				t = s.Params.MgmtContractLib.DecodeTx(tx)
			}

			if t == nil {
				continue
			}
			receipt, err := node.TransactionReceipt(tx.Hash())

			if err != nil || receipt.Status != types.ReceiptStatusSuccessful {
				continue
			}

			switch l1tx := t.(type) {
			case *common.L1DepositTx:
				// todo (@stefan) - remove this hack once the old integrated bridge is removed.
				deposits = append(deposits, tx.Hash())
				totalDeposited.Add(totalDeposited, l1tx.Amount)
				successfulDeposits++
			case *common.L1RollupHashes:
				r, err := getRollupFromBlobHashes(s.ctx, s.Params.BlobResolver, block, l1tx.BlobHashes)
				if err != nil {
					testlog.Logger().Crit("could not decode rollup. ", log.ErrKey, err)
				}
				rollups = append(rollups, r)
				rollupReceipts = append(rollupReceipts, receipt)
				s.Stats.NewRollup(nodeIdx)
			}
		}
	}
	return deposits, rollups, totalDeposited, len(blockchain), successfulDeposits, rollupReceipts
}

func verifyGasBridgeTransactions(t *testing.T, s *Simulation, nodeIdx int) {
	// takes longer for the funds to be bridged across
	time.Sleep(45 * time.Second)
	mbusABI, _ := abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	gasBridgeRecords := s.TxInjector.TxTracker.GasBridgeTransactions
	for _, record := range gasBridgeRecords {
		inputs, err := mbusABI.Methods["sendValueToL2"].Inputs.Unpack(record.L1BridgeTx.Data()[4:])
		if err != nil {
			panic(err)
		}

		receiver := inputs[0].(gethcommon.Address)
		amount := inputs[1].(*big.Int)

		if receiver != record.ReceiverWallet.Address() {
			panic("Test setup is broken. Receiver in tx should match recorded wallet.")
		}
		obsClients := network.CreateAuthClients(s.RPCHandles.RPCClients, record.ReceiverWallet)
		balance, err := obsClients[nodeIdx].BalanceAt(context.Background(), nil)
		if err != nil {
			panic(fmt.Errorf("failed getting balance for bridge transfer receiver. Cause: %w", err))
		}

		if balance.Cmp(amount) != 0 {
			t.Errorf("Node %d: Balance doesnt match the bridged amount. Have: %d, Want: %d", nodeIdx, balance, amount)
		}
	}
}

func checkZenBaseMinting(t *testing.T, s *Simulation) {
	// Map to track the number of transactions per sender
	txCountPerSender := make(map[gethcommon.Address]int)

	// Aggregate transaction counts from Transfer and Withdrawal transactions
	for _, tx := range s.TxInjector.TxTracker.TransferL2Transactions {
		sender := getSender(tx)
		txCountPerSender[sender]++
	}

	for _, tx := range s.TxInjector.TxTracker.WithdrawalL2Transactions {
		sender := getSender(tx)
		txCountPerSender[sender]++
	}

	for _, tx := range s.TxInjector.TxTracker.NativeValueTransferL2Transactions {
		sender := getSender(tx)
		txCountPerSender[sender]++
	}

	// Iterate through each sender and verify ZenBase balance
	for sender, expectedMinted := range txCountPerSender {
		senderRpc := s.RPCHandles.TenWalletClient(sender, 1)
		zenBaseContract, err := ZenBase.NewZenBase(s.ZenBaseAddress, senderRpc)
		if err != nil {
			t.Errorf("Sender %s: Failed to create ZenBase contract. Cause: %s", sender.Hex(), err)
		}
		zenBaseBalance, err := zenBaseContract.BalanceOf(&bind.CallOpts{
			From: sender,
		}, sender)
		if err != nil {
			t.Errorf("Sender %s: Failed to get ZenBase balance. Cause: %s", sender.Hex(), err)
		}

		expectedBalance := big.NewInt(int64(expectedMinted)) // Assuming 1 ZenBase per transaction
		if zenBaseBalance.Cmp(expectedBalance) < 0 {
			t.Errorf("Sender %s: Expected ZenBase balance %d, but found %d", sender.Hex(), expectedBalance, zenBaseBalance)
		}
	}

	rpc := s.RPCHandles.TenWalletClient(s.Params.Wallets.L2FaucetWallet.Address(), 1)
	zenBaseContract, err := ZenBase.NewZenBase(s.ZenBaseAddress, rpc)
	if err != nil {
		t.Errorf("Failed to create ZenBase contract. Cause: %s", err)
	}
	totalSupply, err := zenBaseContract.TotalSupply(&bind.CallOpts{
		From: s.Params.Wallets.L2FaucetWallet.Address(),
	})
	if err != nil {
		panic(err)
	}

	if totalSupply.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("ZenBase total supply is 0")
	}
}

func checkBlockchainOfTenNode(t *testing.T, rpcHandles *network.RPCHandles, minTenHeight uint64, maxEthereumHeight uint64, s *Simulation, wg *sync.WaitGroup, heights []uint64, nodeIdx int) {
	defer wg.Done()
	tenClient := rpcHandles.TenClients[nodeIdx]

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	l1Height, err := rpcHandles.EthClients[nodeIdx].BlockNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve L1 height. Cause: %s", nodeIdx, err)
	}
	if int(maxEthereumHeight)-int(l1Height) > maxBlockDelay {
		t.Errorf("Node %d: TEN node fell behind by %d blocks.", nodeIdx, maxEthereumHeight-l1Height)
	}

	// check that the height of the l2 chain is higher than a minimum expected value.
	headBatchHeader, err := getHeadBatchHeader(tenClient)
	if err != nil {
		t.Error(fmt.Errorf("node %d: %w", nodeIdx, err))
	}

	if headBatchHeader == nil {
		t.Errorf("Node %d: No head rollup recorded. Skipping any further checks for this node.\n", nodeIdx)
		return
	}
	l2Height := headBatchHeader.Number
	if l2Height.Uint64() < minTenHeight {
		t.Errorf("Node %d: Node only mined %d rollups. Expected at least: %d.", nodeIdx, l2Height, minTenHeight)
	}

	// check that the height from the head batch header is consistent with the height returned by eth_blockNumber.
	l2HeightFromBatchNumber, err := tenClient.BatchNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve block number. Cause: %s", nodeIdx, err)
	}
	// due to the difference in calling time, the enclave could produce another batch
	const maxAcceptedDiff = 2
	heightDiff := int(l2HeightFromBatchNumber) - int(l2Height.Uint64())
	if heightDiff > maxAcceptedDiff || heightDiff < -maxAcceptedDiff {
		t.Errorf("Node %d: Node's head batch had a height %d, but %s height was %d", nodeIdx, l2Height, rpc.BatchNumber, l2HeightFromBatchNumber)
	}

	verifyGasBridgeTransactions(t, s, nodeIdx)
	notFoundTransfers, notFoundWithdrawals, notFoundNativeTransfers := FindNotIncludedL2Txs(s.ctx, nodeIdx, rpcHandles, s.TxInjector)
	if notFoundTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Transfer Txs not found in the enclave",
			nodeIdx, notFoundTransfers, len(s.TxInjector.TxTracker.TransferL2Transactions))
	}
	if notFoundWithdrawals > 0 {
		t.Errorf("Node %d: %d out of %d Withdrawal Txs not found in the enclave",
			nodeIdx, notFoundWithdrawals, len(s.TxInjector.TxTracker.WithdrawalL2Transactions))
	}
	if notFoundNativeTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Native Transfer Txs not found in the enclave",
			nodeIdx, notFoundNativeTransfers, len(s.TxInjector.TxTracker.NativeValueTransferL2Transactions))
	}

	checkTransactionReceipts(s.ctx, t, nodeIdx, rpcHandles, s.TxInjector)
	totalSuccessfullyWithdrawn := extractWithdrawals(t, tenClient, nodeIdx)

	// todo: @siliev check that the total amount deposited is the same as the total amount withdraw
	// the previous check is no longer possible due to header removal of xchain messages; and we cannot
	// take the preimages of the crossChainTree as we dont really know them so they must be saved somewhere
	// We can rely on the e2e tests for now.

	injectorDepositedAmt := big.NewInt(0)
	for _, tx := range s.TxInjector.TxTracker.GetL1Transactions() {
		if depTx, ok := tx.(*common.L1DepositTx); ok {
			injectorDepositedAmt.Add(injectorDepositedAmt, depTx.Amount)
		}
	}

	// expected condition : some Txs (stats) did not make it to the blockchain
	// best condition : all Txs (stats) were issue and consumed in the blockchain
	// can't happen : sum of headers withdraws greater than issued Txs (stats)
	if totalSuccessfullyWithdrawn.Cmp(s.Stats.TotalWithdrawalRequestedAmount) > 0 {
		t.Errorf("Node %d: The amount withdrawn %d exceeds the actual amount requested %d", nodeIdx, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// As withdrawals are highly random, here we check that some at least are successful.
	if totalSuccessfullyWithdrawn.Cmp(gethcommon.Big0) < 0 {
		t.Errorf("Node %d: The amount withdrawn %d is far smaller than the amount requested %d", nodeIdx, totalSuccessfullyWithdrawn, s.Stats.TotalWithdrawalRequestedAmount)
	}

	// check that the sum of all balances matches the total amount of money that must be in the system
	// totalAmountInSystem := big.NewInt(0).Sub(s.Stats.TotalDepositedAmount, totalSuccessfullyWithdrawn)
	total := big.NewInt(0)
	for _, wallet := range s.Params.Wallets.SimObsWallets {
		client := rpcHandles.TenWalletClient(wallet.Address(), nodeIdx)
		bal := balance(s.ctx, client, wallet.Address(), s.Params.Wallets.Tokens[testcommon.HOC].L2ContractAddress, nodeIdx)
		total.Add(total, bal)
	}

	/* todo (@stefan) - reenable following check once old contract deployer is phased out.
	if total.Cmp(totalAmountInSystem) != 0 {
		t.Errorf("Node %d: The amount of money in accounts does not match the amount deposited. Found %d , expected %d", nodeIdx, total, totalAmountInSystem)
	} */
	// todo (@stefan) - check that processing transactions in the order specified in the list results in the same balances
	// (execute deposits and transactions and compare to the state in the rollup)

	heights[nodeIdx] = l2Height.Uint64()

	if headBatchHeader.SequencerOrderNo.Uint64() == common.L2GenesisSeqNo {
		return
	}
	// check that the headers are serialised and deserialised correctly, by recomputing a header's hash
	parentHeader, err := tenClient.GetBatchHeaderByHash(headBatchHeader.ParentHash)
	if err != nil {
		t.Errorf("could not retrieve parent of head batch")
		return
	}
	if parentHeader.Hash() != headBatchHeader.ParentHash {
		t.Errorf("mismatch in hash of retrieved header. Parent: %+v\nCurrent: %+v", parentHeader, headBatchHeader)
	}
}

// FindNotIncludedL2Txs returns the number of transfers and withdrawals that were injected but are not present in the L2 blockchain.
func FindNotIncludedL2Txs(ctx context.Context, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) (int, int, int) {
	transfers, withdrawals, nativeTransfers := txInjector.TxTracker.GetL2Transactions()

	notFoundTransfers := 0
	for _, tx := range transfers {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.TenWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundTransfers++
		}
	}

	notFoundWithdrawals := 0
	for _, tx := range withdrawals {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.TenWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundWithdrawals++
		}
	}

	notFoundNativeTransfers := 0
	for _, tx := range nativeTransfers {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.TenWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundNativeTransfers++
		}
	}

	return notFoundTransfers, notFoundWithdrawals, notFoundNativeTransfers
}

func getSender(tx *common.L2Tx) gethcommon.Address {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		panic(fmt.Errorf("couldn't find sender to verify transaction - %w", err))
	}
	return from
}

// Checks that there is a receipt available for each L2 transaction.
func checkTransactionReceipts(ctx context.Context, t *testing.T, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) {
	l2Txs := append(txInjector.TxTracker.TransferL2Transactions, txInjector.TxTracker.WithdrawalL2Transactions...)
	nrSuccessful := 0
	for _, tx := range l2Txs {
		sender := getSender(tx)

		// We check that there is a receipt available for each transaction
		receipt, err := rpcHandles.TenWalletClient(sender, nodeIdx).TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			t.Errorf("node %d: could not retrieve receipt for transaction %s. Cause: %s", nodeIdx, tx.Hash().Hex(), err)
			continue
		}

		// We check that the receipt's logs are relevant to the sender.
		for _, retrievedLog := range receipt.Logs {
			assertRelevantLogsOnly(t, sender.Hex(), *retrievedLog)
		}

		if receipt.Status == types.ReceiptStatusFailed {
			testlog.Logger().Info("Transaction receipt had failed status.", log.TxKey, tx.Hash())
		} else {
			nrSuccessful++
		}
	}

	if nrSuccessful < len(l2Txs)/2 {
		t.Errorf("node %d: More than half the transactions failed. Successful number: %d", nodeIdx, nrSuccessful)
	}

	rpc := rpcHandles.TenWalletClient(txInjector.rndObsWallet().Address(), nodeIdx)
	cfg, err := rpc.GetConfig()
	if err != nil {
		panic(err)
	}

	msgBusAddr := cfg.L2MessageBusAddress

	for _, tx := range txInjector.TxTracker.WithdrawalL2Transactions {
		sender := getSender(tx)
		receipt, err := rpcHandles.TenWalletClient(sender, nodeIdx).TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			continue
		}

		if receipt.Status == types.ReceiptStatusFailed {
			testlog.Logger().Error("[CrossChain] failed withdrawal")
			continue
		}

		if len(receipt.Logs) == 0 {
			testlog.Logger().Error("[CrossChain] no logs in withdrawal?!")
			continue
		}

		if receipt.Logs[0].Address != msgBusAddr {
			testlog.Logger().Error("[CrossChain] wtf")
			continue
		}

		abi, _ := MessageBus.MessageBusMetaData.GetAbi()
		if receipt.Logs[0].Topics[0] != abi.Events["ValueTransfer"].ID {
			testlog.Logger().Error("[CrossChain] wtf")
			continue
		}

		testlog.Logger().Info("[CrossChain] successful withdrawal")
	}
}

func extractWithdrawals(t *testing.T, tenClient *obsclient.ObsClient, nodeIdx int) (totalSuccessfullyWithdrawn *big.Int) {
	totalSuccessfullyWithdrawn = big.NewInt(0)
	header, err := getHeadBatchHeader(tenClient)
	if err != nil {
		t.Error(fmt.Errorf("node %d: %w", nodeIdx, err))
	}
	if header == nil {
		panic(fmt.Sprintf("Node %d: The current head should not be nil", nodeIdx))
	}

	// sum all the withdrawals by traversing the node headers from Head to Genesis
	for {
		if header == nil {
			t.Errorf(fmt.Sprintf("Node %d: Reached a missing rollup", nodeIdx))
			return
		}
		if header.Number.Uint64() == common.L1GenesisHeight {
			return
		}

		// note this retrieves batches currently.
		newHeader, err := tenClient.GetBatchHeaderByHash(header.ParentHash)
		if err != nil {
			t.Errorf(fmt.Sprintf("Node %d: Could not retrieve batch header %s. Cause: %s", nodeIdx, header.ParentHash, err))
			return
		}

		header = newHeader
	}
}

// Terminates all subscriptions and validates the received events.
func checkReceivedLogs(t *testing.T, s *Simulation) {
	logsFromSnapshots := 0
	// rough estimation. In total there should be 2 relevant events for each successful transfer.
	// We assume 66% will pass
	nrLogs := len(s.TxInjector.TxTracker.TransferL2Transactions) * 2
	nrLogs = nrLogs * 2 / 3
	for _, clients := range s.RPCHandles.AuthObsClients {
		for _, client := range clients {
			logsFromSnapshots += checkSnapshotLogs(t, client)
		}
	}
	if logsFromSnapshots < nrLogs {
		t.Errorf("only received %d logs from snapshots, expected at least %d", logsFromSnapshots, nrLogs)
	}

	// In-memory clients cannot handle subscriptions for now.
	if !s.Params.IsInMem {
		for _, sub := range s.Subscriptions {
			sub.Unsubscribe()
		}

		logsFromSubscriptions := 0
		for owner, channels := range s.LogChannels {
			for _, channel := range channels {
				logsFromSubscriptions += checkSubscribedLogs(t, owner, channel)
			}
		}
		if logsFromSubscriptions < nrLogs {
			t.Errorf("only received %d logs from subscriptions, expected at least %d", logsFromSubscriptions, nrLogs)
		}
	}
}

// Checks that a subscription has received the expected logs.
func checkSubscribedLogs(t *testing.T, owner string, channel chan types.Log) int {
	var logs []*types.Log

	for {
		if len(channel) == 0 {
			break
		}
		log := <-channel
		logs = append(logs, &log)
	}

	assertLogsValid(t, owner, logs)
	return len(logs)
}

func checkSnapshotLogs(t *testing.T, client *obsclient.AuthObsClient) int {
	// To exercise the filtering mechanism, we get a snapshot for HOC events only, ignoring POC events.
	hocFilter := common.FilterCriteria{
		Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + testcommon.HOCAddr)},
	}
	logs, err := client.GetLogs(context.Background(), hocFilter)
	if err != nil {
		t.Errorf("could not retrieve logs for client. Cause: %s", err)
	}

	assertLogsValid(t, client.Address().Hex(), logs)

	return len(logs)
}

// Asserts that the logs meet various criteria.
func assertLogsValid(t *testing.T, owner string, logs []*types.Log) {
	for _, receivedLog := range logs {
		assertRelevantLogsOnly(t, owner, *receivedLog)

		logAddrHex := receivedLog.Address.Hex()
		if logAddrHex != "0x"+testcommon.HOCAddr {
			t.Errorf("due to filter, expected logs from the HOC contract only, but got a log from %s", logAddrHex)
		}
	}

	assertNoDupeLogs(t, logs)
}

// Asserts that the log is relevant to the recipient (either a lifecycle event or a relevant user event).
func assertRelevantLogsOnly(t *testing.T, owner string, receivedLog types.Log) {
	// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can (potentially)
	// be user addresses. We filter these out. In theory, we should also check whether the topics exist in the state DB
	// and are contract addresses, but we cannot do this as part of chain validation.
	var userAddrs []string
	for idx, topic := range receivedLog.Topics {
		// The first topic is always the hash of the event.
		if idx == 0 {
			continue
		}

		// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can
		// (potentially) be user addresses.
		topicHex := topic.Hex()
		if topicHex[2:len(zeroBytesHex)+2] == zeroBytesHex {
			userAddrs = append(userAddrs, gethcommon.HexToAddress(topicHex).Hex())
		}
	}

	// If there are no potential user addresses, this is a lifecycle event, and is therefore relevant to everyone.
	if len(userAddrs) == 0 {
		return
	}

	for _, addr := range userAddrs {
		if addr == owner {
			return
		}
	}

	// If we've fallen through to here, it means the log was not relevant.
	t.Errorf("received log that was not relevant (neither a lifecycle event nor relevant to the client's account)")
}

// Asserts that there are no duplicate logs in the provided list.
func assertNoDupeLogs(t *testing.T, logs []*types.Log) {
	logCount := make(map[string]int)

	for _, item := range logs {
		logJSON, err := item.MarshalJSON()
		if err != nil {
			t.Errorf("could not marshal log to JSON to check for duplicate logs")
			continue
		}

		// check if the item/element exist in the duplicate_frequency map
		_, exist := logCount[string(logJSON)]
		if exist {
			logCount[string(logJSON)]++ // increase counter by 1 if already in the map
		} else {
			logCount[string(logJSON)] = 1 // else start counting from 1
		}
	}

	for logJSON, count := range logCount {
		if count > 1 {
			t.Errorf("received duplicate log with body %s", logJSON)
		}
	}
}

// Checks that the various APIs powering Tenscan are working correctly.
func checkTenscan(t *testing.T, s *Simulation) {
	for idx, client := range s.RPCHandles.RPCClients {
		checkTotalTransactions(t, client, idx)
		checkForLatestBatches(t, client, idx)
		checkForLatestRollups(t, client, idx)

		txHashes := getLatestTransactions(t, client, idx)
		for _, txHash := range txHashes {
			checkBatchFromTxs(t, client, txHash, idx)
		}
	}
}

// Checks that the node has stored sufficient transactions.
func checkTotalTransactions(t *testing.T, client rpc.Client, nodeIdx int) {
	var totalTxs *big.Int
	err := client.Call(&totalTxs, rpc.GetTotalTxCount)
	if err != nil {
		t.Errorf("node %d: could not retrieve total transactions. Cause: %s", nodeIdx, err)
	}
	if totalTxs.Int64() < txThreshold {
		t.Errorf("node %d: expected at least %d transactions, but only received %d", nodeIdx, txThreshold, totalTxs)
	}
}

// Checks that we can retrieve the latest batches
func checkForLatestBatches(t *testing.T, client rpc.Client, nodeIdx int) {
	var latestBatches common.BatchListingResponseDeprecated
	pagination := common.QueryPagination{Offset: uint64(0), Size: uint(20)}
	err := client.Call(&latestBatches, rpc.GetBatchListing, &pagination)
	if err != nil {
		t.Errorf("node %d: could not retrieve latest batches. Cause: %s", nodeIdx, err)
	}
	// the batch listing function returns the last received batches , but it might receive them in a random order
	if len(latestBatches.BatchesData) < 5 {
		t.Errorf("node %d: expected at least %d batches, but only received %d", nodeIdx, 5, len(latestBatches.BatchesData))
	}
}

// Checks that we can retrieve the latest rollups
func checkForLatestRollups(t *testing.T, client rpc.Client, nodeIdx int) {
	var latestRollups common.RollupListingResponse
	pagination := common.QueryPagination{Offset: uint64(0), Size: uint(5)}
	err := client.Call(&latestRollups, rpc.GetRollupListing, &pagination)
	if err != nil {
		t.Errorf("node %d: could not retrieve latest transactions. Cause: %s", nodeIdx, err)
	}
	if len(latestRollups.RollupsData) != 5 {
		t.Errorf("node %d: expected at least %d transactions, but only received %d", nodeIdx, 5, len(latestRollups.RollupsData))
	}
}

func getLatestTransactions(t *testing.T, client rpc.Client, nodeIdx int) []gethcommon.Hash {
	var transactionResponse common.TransactionListingResponse
	var txHashes []gethcommon.Hash
	pagination := common.QueryPagination{Offset: uint64(0), Size: uint(5)}
	err := client.Call(&transactionResponse, rpc.GetPublicTransactionData, &pagination)
	if err != nil {
		t.Errorf("node %d: could not retrieve latest transactions. Cause: %s", nodeIdx, err)
	}

	for _, transaction := range transactionResponse.TransactionsData {
		txHashes = append(txHashes, transaction.TransactionHash)
	}

	return txHashes
}

// Retrieves the batch using the transaction hash, and validates it.
func checkBatchFromTxs(t *testing.T, client rpc.Client, txHash gethcommon.Hash, nodeIdx int) {
	var batchByTx *common.ExtBatch
	err := client.Call(&batchByTx, rpc.GetBatchByTx, txHash)
	if err != nil {
		t.Errorf("node %d: could not retrieve batch for transaction. Cause: %s", nodeIdx, err)
		return
	}

	var containsTx bool
	for _, txHashFromBatch := range batchByTx.TxHashes {
		if txHashFromBatch == txHash {
			containsTx = true
		}
	}
	if !containsTx {
		t.Errorf("node %d: retrieved batch by transaction, but transaction was missing from batch", nodeIdx)
	}

	var batchByHash *common.ExtBatch
	err = client.Call(&batchByHash, rpc.GetBatch, batchByTx.Header.Hash())
	if err != nil {
		t.Errorf("node %d: could not retrieve batch by hash. Cause: %s", nodeIdx, err)
		return
	}
	if batchByHash.Header.Hash() != batchByTx.Header.Hash() {
		t.Errorf("node %d: retrieved batch by hash, but hash was incorrect", nodeIdx)
	}
}

func getRollupFromBlobHashes(ctx context.Context, blobResolver l1.BlobResolver, block *types.Block, blobHashes []gethcommon.Hash) (*common.ExtRollup, error) {
	blobs, err := blobResolver.FetchBlobs(ctx, block.Header(), blobHashes)
	if err != nil {
		return nil, fmt.Errorf("could not fetch blobs from hashes during chain validation. Cause: %w", err)
	}
	data, err := ethadapter.DecodeBlobs(blobs)
	if err != nil {
		return nil, fmt.Errorf("error decoding rollup blob. Cause: %w", err)
	}

	var rollup common.ExtRollup
	if err := rlp.DecodeBytes(data, &rollup); err != nil {
		return nil, fmt.Errorf("could not decode rollup. Cause: %w", err)
	}
	return &rollup, nil
}
