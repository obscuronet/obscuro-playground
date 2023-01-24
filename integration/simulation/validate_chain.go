package simulation

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/ethereummock"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/integration/simulation/network"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/rollupextractor"

	"github.com/obscuronet/go-obscuro/go/rpc"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"

	erc20 "github.com/obscuronet/go-obscuro/integration/erc20contract/generated/EthERC20"
)

const (
	// The threshold number of transactions below which we consider the simulation to have failed. We generally expect far
	// more than this, but this is a sanity check to ensure the simulation doesn't stop after a single transaction of each
	// type, for example.
	txThreshold = 5
	// As above, but for the number of logs received via subscriptions.
	logsThreshold = 5
	// The maximum number of blocks an Obscuro node can fall behind
	maxBlockDelay = 5
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// After a simulation has run, check as much as possible that the outputs of the simulation are expected.
// For example, all injected transactions were processed correctly, the height of the rollup chain is a function of the total
// time of the simulation and the average block duration, that all Obscuro nodes are roughly in sync, etc
func checkNetworkValidity(t *testing.T, s *Simulation) {
	checkTransactionsWereInjected(t, s)
	l1MaxHeight := checkEthereumBlockchainValidity(t, s)
	checkObscuroBlockchainValidity(t, s, l1MaxHeight)
	checkReceivedLogs(t, s)
	checkObscuroscan(t, s)
}

// Ensures that the transaction injector actually injected some minimum amount of transactions (to avoid edge-cases
// where zero or very few transactions were injected, and thus the L1 and L2 are trivially aligned).
func checkTransactionsWereInjected(t *testing.T, s *Simulation) {
	if len(s.TxInjector.TxTracker.L1Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d L1 transactions. At least %d expected", len(s.TxInjector.TxTracker.L1Transactions), txThreshold)
	}
	if len(s.TxInjector.TxTracker.TransferL2Transactions) < txThreshold {
		t.Errorf("Simulation only issued %d transfer L2 transactions. At least %d expected", len(s.TxInjector.TxTracker.TransferL2Transactions), txThreshold)
	}
	//TODO: Reenable when old contract deployer phased out.
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

// checkObscuroBlockchainValidity - perform the following checks
// - minimum height - the chain has a minimum number of rollups
// - check height is similar
// - check no duplicate txs
// - check efficiency - no of created blocks/ height
// - check amount in the system
// - check withdrawals/deposits
func checkObscuroBlockchainValidity(t *testing.T, s *Simulation, maxL1Height uint64) {
	// Sanity check number for a minimum height
	minHeight := uint64(float64(s.Params.SimulationTime.Microseconds()) / (2 * float64(s.Params.AvgBlockDuration)))

	// process the blockchain of each node in parallel to minimize the difference between them since they are still running
	heights := make([]uint64, len(s.RPCHandles.ObscuroClients))
	var wg sync.WaitGroup
	for idx := range s.RPCHandles.ObscuroClients {
		wg.Add(1)
		go checkBlockchainOfObscuroNode(t, s.RPCHandles, minHeight, maxL1Height, s, &wg, heights, idx)
	}
	wg.Wait()
	min, max := minMax(heights)
	// This checks that all the nodes are in sync. When a node falls behind with processing blocks it might highlight a problem.
	if max-min > max/10 {
		t.Errorf("There is a problem with the Obscuro chain. Nodes fell out of sync. Max height: %d. Min height: %d", max, min)
	}
}

func checkBlockchainOfEthereumNode(t *testing.T, node ethadapter.EthClient, minHeight uint64, s *Simulation, nodeIdx int) uint64 {
	head, err := node.FetchHeadBlock()
	if err != nil {
		t.Errorf("Node %d: Could not find head block. Cause: %s", nodeIdx, err)
	}
	height := head.NumberU64()

	if height < minHeight {
		t.Errorf("Node %d: There were only %d blocks mined. Expected at least: %d.", nodeIdx, height, minHeight)
	}

	deposits, rollups, _, blockCount, _ := ExtractDataFromEthereumChain(ethereummock.MockGenesisBlock, head, node, s, nodeIdx)
	s.Stats.TotalL1Blocks = uint64(blockCount)

	if len(findHashDups(deposits)) > 0 {
		dups := findHashDups(deposits)
		t.Errorf("Node %d: Found Deposit duplicates: %v", nodeIdx, dups)
	}
	if len(findRollupDups(rollups)) > 0 {
		dups := findRollupDups(rollups)
		t.Errorf("Node %d: Found Rollup duplicates: %v", nodeIdx, dups)
	}

	/* TODO: Reenable while old contract deployer is phased out.
	if s.Stats.TotalDepositedAmount.Cmp(totalDeposited) != 0 {
		t.Errorf("Node %d: Deposit[%d] amounts don't match. Found %d , expected %d", nodeIdx, len(deposits), totalDeposited, s.Stats.TotalDepositedAmount)
	}

	if s.Stats.TotalDepositedAmount.Cmp(gethcommon.Big0) == 0 {
		t.Errorf("Node %d: No deposits", nodeIdx)
	} */

	efficiency := float64(s.Stats.TotalL1Blocks-height) / float64(s.Stats.TotalL1Blocks)
	if efficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: Efficiency in L1 is %f. Expected:%f. Number: %d.", nodeIdx, efficiency, s.Params.L1EfficiencyThreshold, height)
	}

	// compare the number of reorgs for this node against the height
	reorgs := s.Stats.NoL1Reorgs[node.Info().L2ID]
	reorgEfficiency := float64(reorgs) / float64(height)
	if reorgEfficiency > s.Params.L1EfficiencyThreshold {
		t.Errorf("Node %d: The number of reorgs is too high: %d. ", nodeIdx, reorgs)
	}

	// Check that all the rollups are produced by aggregators.
	for _, rollup := range rollups {
		nodeID, err := strconv.ParseInt(rollup.Header.Agg.Hex()[2:], 16, 64)
		if err != nil {
			t.Errorf("Node %d: Could not parse node's integer ID. Cause: %s", nodeIdx, err)
			continue
		}
		if network.GetNodeType(int(nodeID)) != common.Sequencer {
			t.Errorf("Node %d: Found rollup produced by non-sequencer %d", nodeIdx, nodeID)
		}
	}

	checkInjectedTxsInRollups(t, s, nodeIdx, rollups)

	return height
}

// ExtractDataFromEthereumChain returns the deposits, rollups, total amount deposited and length of the blockchain
// between the start block and the end block.
func ExtractDataFromEthereumChain(
	startBlock *types.Block,
	endBlock *types.Block,
	node ethadapter.EthClient,
	s *Simulation,
	nodeIdx int,
) ([]gethcommon.Hash, []*common.ExtRollup, *big.Int, int, uint64) {
	deposits := make([]gethcommon.Hash, 0)
	rollups := make([]*common.ExtRollup, 0)
	totalDeposited := big.NewInt(0)

	blockchain := node.BlocksBetween(startBlock, endBlock)
	succesfulDeposits := uint64(0)
	for _, block := range blockchain {
		for _, tx := range block.Transactions() {
			t := s.Params.ERC20ContractLib.DecodeTx(tx)
			if t == nil {
				t = s.Params.MgmtContractLib.DecodeTx(tx)
			}

			if t == nil {
				continue
			}
			switch l1tx := t.(type) {
			case *ethadapter.L1DepositTx:
				// TODO: Remove this hack once the old integrated bridge is removed.
				if receipt, err := node.TransactionReceipt(tx.Hash()); err == nil && receipt.Status == 1 {
					deposits = append(deposits, tx.Hash())
					totalDeposited.Add(totalDeposited, l1tx.Amount)
					succesfulDeposits++
				}
			case *ethadapter.L1RollupTx:
				r, err := common.DecodeRollup(l1tx.Rollup)
				if err != nil {
					testlog.Logger().Crit("could not decode rollup. ", log.ErrKey, err)
				}
				rollups = append(rollups, r)
				if node.IsBlockAncestor(block, r.Header.L1Proof) {
					// only count the rollup if it is published in the right branch
					// todo - once logic is added to the l1 - this can be made into a check
					s.Stats.NewRollup(nodeIdx)
				}
			}
		}
	}
	return deposits, rollups, totalDeposited, len(blockchain), succesfulDeposits
}

// Check all the injected L2 transactions were included in the rollups.
func checkInjectedTxsInRollups(t *testing.T, s *Simulation, nodeIdx int, rollups []*common.ExtRollup) {
	rollupTxs := make(map[common.TxHash]bool)
	for _, rollup := range rollups {
		for _, batch := range rollup.Batches {
			for _, tx := range batch.TxHashes {
				rollupTxs[tx] = true
			}
		}
	}

	transferTxs, withdrawalTxs, nativeTransferTxs := s.TxInjector.TxTracker.GetL2Transactions()
	notFoundTransfers, notFoundWithdrawals, notFoundNativeTransfers := 0, 0, 0
	for _, tx := range transferTxs {
		if _, found := rollupTxs[tx.Hash()]; !found {
			notFoundTransfers++
		}
	}
	for _, tx := range withdrawalTxs {
		if _, found := rollupTxs[tx.Hash()]; !found {
			notFoundWithdrawals++
		}
	}
	for _, tx := range nativeTransferTxs {
		if _, found := rollupTxs[tx.Hash()]; !found {
			notFoundNativeTransfers++
		}
	}

	if notFoundTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Transfer Txs not found in rollups",
			nodeIdx, notFoundTransfers, len(s.TxInjector.TxTracker.TransferL2Transactions))
	}
	if notFoundWithdrawals > 0 {
		t.Errorf("Node %d: %d out of %d Withdrawal Txs not found in rollups",
			nodeIdx, notFoundWithdrawals, len(s.TxInjector.TxTracker.WithdrawalL2Transactions))
	}
	if notFoundNativeTransfers > 0 {
		t.Errorf("Node %d: %d out of %d Native Transfer Txs not found in rollups",
			nodeIdx, notFoundNativeTransfers, len(s.TxInjector.TxTracker.NativeValueTransferL2Transactions))
	}
}

func checkBlockchainOfObscuroNode(t *testing.T, rpcHandles *network.RPCHandles, minObscuroHeight uint64, maxEthereumHeight uint64, s *Simulation, wg *sync.WaitGroup, heights []uint64, nodeIdx int) {
	defer wg.Done()
	obscuroClient := rpcHandles.ObscuroClients[nodeIdx]

	// check that the L1 view is consistent with the L1 network.
	// We cast to int64 to avoid an overflow when l1Height is greater than maxEthereumHeight (due to additional blocks
	// produced since maxEthereumHeight was calculated from querying all L1 nodes - the simulation is still running, so
	// new blocks might have been added in the meantime).
	l1Height, err := rpcHandles.EthClients[nodeIdx].BlockNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve L1 height. Cause: %s", nodeIdx, err)
	}
	if int(maxEthereumHeight)-int(l1Height) > maxBlockDelay {
		t.Errorf("Node %d: Obscuro node fell behind by %d blocks.", nodeIdx, maxEthereumHeight-l1Height)
	}

	// check that the height of the Rollup chain is higher than a minimum expected value.
	headRollupHeader := getHeadRollupHeader(obscuroClient, nodeIdx)
	if headRollupHeader == nil {
		t.Errorf("Node %d: No head rollup recorded. Skipping any further checks for this node.\n", nodeIdx)
		return
	}
	l2Height := headRollupHeader.Number
	if l2Height.Uint64() < minObscuroHeight {
		t.Errorf("Node %d: Node only mined %d rollups. Expected at least: %d.", nodeIdx, l2Height, minObscuroHeight)
	}

	// check that the height from the rollup header is consistent with the height returned by eth_blockNumber.
	l2HeightFromRollupNumber, err := obscuroClient.RollupNumber()
	if err != nil {
		t.Errorf("Node %d: Could not retrieve block number. Cause: %s", nodeIdx, err)
	}
	if l2HeightFromRollupNumber != l2Height.Uint64() {
		t.Errorf("Node %d: Node's head rollup had a height %d, but %s height was %d", nodeIdx, l2Height, rpc.RollupNumber, l2HeightFromRollupNumber)
	}

	totalL2Blocks := s.Stats.NoL2Blocks[nodeIdx]
	// in case the blockchain has advanced above what was collected, there is no longer a point to this check
	if l2Height.Uint64() <= totalL2Blocks {
		efficiencyL2 := float64(totalL2Blocks-l2Height.Uint64()) / float64(totalL2Blocks)
		if efficiencyL2 > s.Params.L2EfficiencyThreshold {
			t.Errorf("Node %d: Efficiency in L2 is %f. Expected:%f", nodeIdx, efficiencyL2, s.Params.L2EfficiencyThreshold)
		}
	}

	// check that the pobi protocol doesn't waste too many blocks.
	// todo - find the block where the genesis was published
	efficiency := float64(l1Height-l2Height.Uint64()) / float64(l1Height)
	if efficiency > s.Params.L2ToL1EfficiencyThreshold {
		t.Errorf("Node %d: L2 to L1 Efficiency is %f. Expected:%f", nodeIdx, efficiency, s.Params.L2ToL1EfficiencyThreshold)
	}

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

	totalSuccessfullyWithdrawn := extractWithdrawals(t, obscuroClient, nodeIdx)

	totalAmountLogged := getLoggedWithdrawals(minObscuroHeight, obscuroClient, headRollupHeader)
	if totalAmountLogged.Cmp(totalSuccessfullyWithdrawn) != 0 {
		t.Errorf("Node %d: Logged withdrawals do not match!", nodeIdx)
	}

	injectorDepositedAmt := big.NewInt(0)
	for _, tx := range s.TxInjector.TxTracker.GetL1Transactions() {
		if depTx, ok := tx.(*ethadapter.L1DepositTx); ok {
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
		client := rpcHandles.ObscuroWalletClient(wallet.Address(), nodeIdx)
		bal := balance(s.ctx, client, wallet.Address(), s.Params.Wallets.Tokens[rollupextractor.HOC].L2ContractAddress)
		total.Add(total, bal)
	}

	/* TODO: Reenable following check once old contract deployer is phased out.
	if total.Cmp(totalAmountInSystem) != 0 {
		t.Errorf("Node %d: The amount of money in accounts does not match the amount deposited. Found %d , expected %d", nodeIdx, total, totalAmountInSystem)
	} */
	// TODO Check that processing transactions in the order specified in the list results in the same balances
	// (execute deposits and transactions and compare to the state in the rollup)

	heights[nodeIdx] = l2Height.Uint64()

	// check that the headers are serialised and deserialised correctly, by recomputing a header's hash
	parentHeader, err := obscuroClient.RollupHeaderByHash(headRollupHeader.ParentHash)
	if err != nil {
		t.Errorf("could not retrieve parent of head rollup")
	}
	if parentHeader.Hash() != headRollupHeader.ParentHash {
		t.Errorf("mismatch in hash of retrieved header")
	}
}

func getLoggedWithdrawals(minObscuroHeight uint64, obscuroClient *obsclient.ObsClient, currentHeader *common.BatchHeader) *big.Int {
	totalAmountLogged := big.NewInt(0)
	for i := minObscuroHeight; i < currentHeader.Number.Uint64(); i++ {
		header, err := obscuroClient.RollupHeaderByNumber(big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}

		for _, msg := range header.CrossChainMessages {
			contractAbi, err := abi.JSON(strings.NewReader(erc20.EthERC20MetaData.ABI))
			if err != nil {
				panic(err)
			}

			transfer := map[string]interface{}{}
			err = contractAbi.Methods["transferFrom"].Inputs.UnpackIntoMap(transfer, msg.Payload) // can't figure out how to unpack it without cheating, geth is kinda clunky
			if err != nil {
				panic(err)
			}

			amount := transfer["amount"].(*big.Int)
			totalAmountLogged = totalAmountLogged.Add(totalAmountLogged, amount)
		}
	}
	return totalAmountLogged
}

// FindNotIncludedL2Txs returns the number of transfers and withdrawals that were injected but are not present in the L2 blockchain.
func FindNotIncludedL2Txs(ctx context.Context, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) (int, int, int) {
	transfers, withdrawals, nativeTransfers := txInjector.TxTracker.GetL2Transactions()
	notFoundTransfers := 0
	for _, tx := range transfers {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundTransfers++
		}
	}

	notFoundWithdrawals := 0
	for _, tx := range withdrawals {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundWithdrawals++
		}
	}

	notFoundNativeTransfers := 0
	for _, tx := range nativeTransfers {
		sender := getSender(tx)
		// because of viewing key encryption we need to get the RPC client for this specific node for the wallet that sent the transaction
		l2tx, _, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionByHash(ctx, tx.Hash())
		if err != nil || l2tx == nil {
			notFoundNativeTransfers++
		}
	}

	return notFoundTransfers, notFoundWithdrawals, notFoundNativeTransfers
}

func getSender(tx *common.L2Tx) gethcommon.Address {
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		panic(fmt.Errorf("couldn't find sender to verify transaction - %w", err))
	}
	return msg.From()
}

// Checks that there is a receipt available for each L2 transaction.
func checkTransactionReceipts(ctx context.Context, t *testing.T, nodeIdx int, rpcHandles *network.RPCHandles, txInjector *TransactionInjector) {
	l2Txs := append(txInjector.TxTracker.TransferL2Transactions, txInjector.TxTracker.WithdrawalL2Transactions...)

	for _, tx := range l2Txs {
		sender := getSender(tx)

		// We check that there is a receipt available for each transaction
		receipt, err := rpcHandles.ObscuroWalletClient(sender, nodeIdx).TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			t.Errorf("node %d: could not retrieve receipt for transaction %s. Cause: %s", nodeIdx, tx.Hash().Hex(), err)
			continue
		}

		// We check that the receipt's logs are relevant to the sender.
		for _, retrievedLog := range receipt.Logs {
			assertRelevantLogsOnly(t, sender.Hex(), *retrievedLog)
		}

		if receipt.Status == types.ReceiptStatusFailed {
			testlog.Logger().Info("Transaction receipt had failed status.", log.TxKey, tx.Hash().Hex())
		}
	}
}

func extractWithdrawals(t *testing.T, obscuroClient *obsclient.ObsClient, nodeIdx int) (totalSuccessfullyWithdrawn *big.Int) {
	totalSuccessfullyWithdrawn = big.NewInt(0)
	header := getHeadRollupHeader(obscuroClient, nodeIdx)
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

		newHeader, err := obscuroClient.RollupHeaderByHash(header.ParentHash)
		if err != nil {
			t.Errorf(fmt.Sprintf("Node %d: Could not retrieve rollup header %s. Cause: %s", nodeIdx, header.ParentHash, err))
			return
		}

		header = newHeader
	}
}

// Terminates all subscriptions and validates the received events.
func checkReceivedLogs(t *testing.T, s *Simulation) {
	logsFromSnapshots := 0
	for _, clients := range s.RPCHandles.AuthObsClients {
		for _, client := range clients {
			logsFromSnapshots += checkSnapshotLogs(t, client)
		}
	}
	if logsFromSnapshots < logsThreshold {
		t.Errorf("only received %d logs from snapshots, expected at least %d", logsFromSnapshots, logsThreshold)
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
		if logsFromSubscriptions < logsThreshold {
			t.Errorf("only received %d logs from subscriptions, expected at least %d", logsFromSubscriptions, logsThreshold)
		}
	}
}

// Checks that a subscription has received the expected logs.
func checkSubscribedLogs(t *testing.T, owner string, channel chan common.IDAndLog) int {
	var logs []*types.Log

out:
	for {
		select {
		case idAndLog := <-channel:
			logs = append(logs, idAndLog.Log)

		// The logs will have built up on the channel throughout the simulation, so they should arrive immediately.
		// However, if we use a `default` case, only the first one arrives. Some minimal wait is required.
		case <-time.After(time.Millisecond):
			break out
		}
	}

	assertLogsValid(t, owner, logs)
	return len(logs)
}

func checkSnapshotLogs(t *testing.T, client *obsclient.AuthObsClient) int {
	// To exercise the filtering mechanism, we get a snapshot for HOC events only, ignoring POC events.
	hocFilter := common.FilterCriteriaJSON{
		Addresses: []gethcommon.Address{gethcommon.HexToAddress("0x" + rollupextractor.HOCAddr)},
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
		if logAddrHex != "0x"+rollupextractor.HOCAddr {
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

// Checks that the various APIs powering Obscuroscan are working correctly.
func checkObscuroscan(t *testing.T, s *Simulation) {
	for idx, client := range s.RPCHandles.RPCClients {
		checkTotalTransactions(t, client, idx)
		latestTxHashes := checkLatestTxs(t, client, idx)
		for _, txHash := range latestTxHashes {
			checkBatchFromTxs(t, client, txHash, idx)
		}
	}
}

// Checks that the node has stored sufficient transactions.
func checkTotalTransactions(t *testing.T, client rpc.Client, nodeIdx int) {
	var totalTxs *big.Int
	err := client.Call(&totalTxs, rpc.GetTotalTxs)
	if err != nil {
		t.Errorf("node %d: could not retrieve total transactions. Cause: %s", nodeIdx, err)
	}
	if totalTxs.Int64() < txThreshold {
		t.Errorf("node %d: expected at least %d transactions, but only received %d", nodeIdx, txThreshold, totalTxs)
	}
}

// Checks that we can retrieve the latest transactions for the node.
func checkLatestTxs(t *testing.T, client rpc.Client, nodeIdx int) []gethcommon.Hash {
	var latestTxHashes []gethcommon.Hash
	err := client.Call(&latestTxHashes, rpc.GetLatestTxs, txThreshold)
	if err != nil {
		t.Errorf("node %d: could not retrieve latest transactions. Cause: %s", nodeIdx, err)
	}
	if len(latestTxHashes) != txThreshold {
		t.Errorf("node %d: expected at least %d transactions, but only received %d", nodeIdx, txThreshold, len(latestTxHashes))
	}
	return latestTxHashes
}

// Retrieves the batch using the transaction hash, and validates it.
func checkBatchFromTxs(t *testing.T, client rpc.Client, txHash gethcommon.Hash, nodeIdx int) {
	var batchByTx *common.ExtBatch
	err := client.Call(&batchByTx, rpc.GetBatchForTx, txHash)
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
