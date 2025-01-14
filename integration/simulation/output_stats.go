package simulation

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/integration/common/testlog"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ten-protocol/go-ten/go/common"
)

// OutputStats decouples the processing of data and the collection of statistics
// there's a bit more to do around this, this serves as a first iteration
type OutputStats struct {
	simulation *Simulation

	l2RollupCountInHeaders    int // Number of rollups counted while node rollup header traversing
	l2RollupCountInL1Blocks   int // Number of rollups counted while traversing the node block header and searching the txs
	l2RollupTxCountInL1Blocks int // Number of rollup Txs counted while traversing the node block header
	l1Height                  int // Last known l1 block height
	l2Height                  int // Last known l2 block height

	canonicalERC20DepositCount int // Number of erc20 deposits on the canonical chain
}

// NewOutputStats processes the simulation and retrieves the output statistics
func NewOutputStats(simulation *Simulation) *OutputStats {
	outputStats := &OutputStats{
		simulation: simulation,
	}

	outputStats.countBlockChain()
	outputStats.populateHeights()

	return outputStats
}

func (o *OutputStats) populateHeights() {
	l1Height, err := o.simulation.RPCHandles.EthClients[0].BlockNumber()
	if err != nil {
		panic(fmt.Errorf("simulation failed because could not read L1 height. Cause: %w", err))
	}
	o.l1Height = int(l1Height)

	tenClient := o.simulation.RPCHandles.TenClients[0]
	hRollup, err := getHeadBatchHeader(tenClient)
	if err != nil {
		panic(err)
	}

	o.l2Height = int(hRollup.Number.Uint64())
}

func (o *OutputStats) countBlockChain() {
	l1Node := o.simulation.RPCHandles.EthClients[0]
	tenClient := o.simulation.RPCHandles.TenClients[0]

	// iterate the Node Headers and get the rollups
	header, err := getHeadBatchHeader(tenClient)
	if err != nil {
		panic(err)
	}

	for {
		if header != nil && !bytes.Equal(header.Hash().Bytes(), (common.L1BlockHash{}).Bytes()) {
			break
		}

		o.l2RollupCountInHeaders++

		header, err = tenClient.GetBatchHeaderByHash(header.ParentHash)
		if err != nil {
			testlog.Logger().Crit("could not retrieve rollup by hash.", log.ErrKey, err)
		}
	}

	// iterate the L1 Blocks and get the rollups
	for block, _ := l1Node.FetchHeadBlock(); block != nil && !bytes.Equal(block.Hash().Bytes(), (common.L1BlockHash{}).Bytes()); block, _ = l1Node.HeaderByHash(block.ParentHash) {
		b, err := l1Node.BlockByHash(block.Hash())
		if err != nil {
			panic(err)
		}
		o.incrementStats(b, l1Node)
	}
}

func (o *OutputStats) incrementStats(block *types.Block, _ ethadapter.EthClient) {
	for _, tx := range block.Transactions() {
		t := o.simulation.Params.MgmtContractLib.DecodeTx(tx)
		if t == nil {
			t = o.simulation.Params.ERC20ContractLib.DecodeTx(tx)
		}

		if t == nil {
			continue
		}

		switch l1Tx := t.(type) {
		case *common.L1RollupTx:
			_, err := common.DecodeRollup(l1Tx.Rollup)
			if err != nil {
				testlog.Logger().Crit("could not decode rollup.", log.ErrKey, err)
			}
			//if l1Node.IsBlockAncestor(block, r.Header.L1Proof) {
			//	o.l2RollupCountInL1Blocks++
			//	for _, batch := range r.BatchPayloads {
			//		o.l2RollupTxCountInL1Blocks += len(batch.TxHashes)
			//	}
			//}

		case *common.L1DepositTx:
			o.canonicalERC20DepositCount++
		}
	}
}

func (o *OutputStats) String() string {
	return fmt.Sprintf("\n"+
		"nrMiners: %d\n"+
		"l1Height: %d\n"+
		"l2Height: %d\n"+
		"totalL1Blocks: %d\n"+
		"totalL2Blocks: %v\n"+
		"l2RollupCountInHeaders: %d\n"+
		"l2RollupCountInL1Blocks: %d\n"+
		"l2RollupTxCountInL1Blocks: %d\n"+
		"maxRollupsPerBlock: %d \n"+
		"nrEmptyBlocks: %d\n"+
		"noL1Reorgs: %+v\n"+
		"noL2Recalcs: %+v\n"+
		"totalDepositedAmount: %d\n"+
		"totalWithdrawnAmount: %d\n"+
		"rollupWithMoreRecentProof: %d\n"+
		"nrTransferTransactions: %d\n"+
		"nrNativeTransferTransactions: %d\n"+
		"nrBlockParsedERC20Deposits: %d\n"+
		"gasBridgeCount: %d\n",
		o.simulation.Stats.NrMiners,
		o.l1Height,
		o.l2Height,
		o.simulation.Stats.TotalL1Blocks,
		o.simulation.Stats.NoL2Blocks,
		o.l2RollupCountInHeaders,
		o.l2RollupCountInL1Blocks,
		o.l2RollupTxCountInL1Blocks,
		o.simulation.Stats.MaxRollupsPerBlock,
		o.simulation.Stats.NrEmptyBlocks,
		o.simulation.Stats.NoL1Reorgs,
		o.simulation.Stats.NoL2Recalcs,
		o.simulation.Stats.TotalDepositedAmount,
		o.simulation.Stats.TotalWithdrawalRequestedAmount,
		o.simulation.Stats.RollupWithMoreRecentProofCount,
		o.simulation.Stats.NrTransferTransactions,
		o.simulation.Stats.NrNativeTransferTransactions,
		o.canonicalERC20DepositCount,
		len(o.simulation.TxInjector.TxTracker.GasBridgeTransactions),
	)
}
