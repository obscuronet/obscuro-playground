package simulation

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	testcommon "github.com/obscuronet/go-obscuro/integration/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
)

const (
	testLogs = "../.build/simulations/"
)

func setupSimTestLog(simType string) {
	testlog.Setup(&testlog.Cfg{
		LogDir:      testLogs,
		TestType:    "sim-log",
		TestSubtype: simType,
		LogLevel:    log.LvlTrace,
	})
}

func minMax(arr []uint64) (min uint64, max uint64) {
	min = ^uint64(0)
	for _, no := range arr {
		if no < min {
			min = no
		}
		if no > max {
			max = no
		}
	}
	return
}

// Uses the client to retrieve the current rollup head.
func getHeadRollupHeader(client *obsclient.ObsClient, nodeIdx int) *common.BatchHeader {
	headRollupHeight, err := client.RollupNumber()
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed attempt to retrieve head rollup height. Cause: %w", err))
	}

	headRollupHeader, err := client.RollupHeaderByNumber(big.NewInt(int64(headRollupHeight)))
	if err != nil {
		// panic(fmt.Errorf("node %d: simulation failed due to failed attempt to retrieve rollup with height %d. Cause: %w", nodeIdx, headRollupHeight, err))
		return nil
	}
	return headRollupHeader
}

// Uses the client to retrieve the balance of the wallet with the given address.
func balance(ctx context.Context, client *obsclient.AuthObsClient, address gethcommon.Address, l2ContractAddress *gethcommon.Address) *big.Int {
	balanceData := erc20contractlib.CreateBalanceOfData(address)

	callMsg := ethereum.CallMsg{
		From: address,
		To:   l2ContractAddress,
		Data: balanceData,
	}

	response, err := client.CallContract(ctx, callMsg, nil)
	if err != nil {
		panic(fmt.Errorf("simulation failed due to failed RPC call. Cause: %w", err))
	}
	b := new(big.Int)
	// remove the "0x" prefix (we already confirmed it is present), convert the remaining hex value (base 16) to a balance number
	b.SetString(string(response)[2:], 16)
	return b
}

// FindHashDups - returns a map of all hashes that appear multiple times, and how many times
func findHashDups(list []gethcommon.Hash) map[gethcommon.Hash]int {
	elementCount := make(map[gethcommon.Hash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item]
		if exist {
			elementCount[item]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item] = 1 // else start counting from 1
		}
	}
	dups := make(map[gethcommon.Hash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: %s\n", u)
		}
	}
	return dups
}

// FindRollupDups - returns a map of all L2 root hashes that appear multiple times, and how many times
func findRollupDups(list []*common.ExtRollup) map[common.L2BatchHash]int {
	elementCount := make(map[common.L2BatchHash]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := elementCount[item.Hash()]
		if exist {
			elementCount[item.Hash()]++ // increase counter by 1 if already in the map
		} else {
			elementCount[item.Hash()] = 1 // else start counting from 1
		}
	}
	dups := make(map[common.L2BatchHash]int)
	for u, i := range elementCount {
		if i > 1 {
			dups[u] = i
			fmt.Printf("Dup: r_%d\n", common.ShortHash(u))
		}
	}
	return dups
}

func sleepRndBtw(min time.Duration, max time.Duration) {
	time.Sleep(testcommon.RndBtwTime(min, max))
}
