package simulation

import (
	"bytes"
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/common"
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
	obscuroClient := o.simulation.RPCHandles.ObscuroClients[0]
	o.l1Height = int(getCurrentBlockHeadHeight(obscuroClient))
	o.l2Height = int(getCurrentRollupHead(obscuroClient).Number.Uint64())
}

func (o *OutputStats) countBlockChain() {
	l1Node := o.simulation.RPCHandles.EthClients[0]
	l2Client := o.simulation.RPCHandles.ObscuroClients[0]

	// iterate the Node Headers and get the rollups
	for header := getCurrentRollupHead(l2Client); header != nil && !bytes.Equal(header.Hash().Bytes(), common.GenesisHash.Bytes()); header = getRollupHeader(l2Client, header.ParentHash) {
		o.l2RollupCountInHeaders++
	}

	// iterate the L1 Blocks and get the rollups
	for headBlock := l1Node.FetchHeadBlock(); headBlock != nil && !bytes.Equal(headBlock.Hash().Bytes(), common.GenesisHash.Bytes()); headBlock, _ = l1Node.BlockByHash(headBlock.ParentHash()) {
		for _, tx := range headBlock.Transactions() {
			t := o.simulation.Params.MgmtContractLib.DecodeTx(tx)
			if t == nil {
				t = o.simulation.Params.ERC20ContractLib.DecodeTx(tx)
			}

			if t == nil {
				continue
			}

			switch l1Tx := t.(type) {
			case *ethadapter.L1RollupTx:
				r, err := common.DecodeRollup(l1Tx.Rollup)
				if err != nil {
					testlog.Logger().Crit("could not decode rollup.", log.ErrKey, err)
				}
				if l1Node.IsBlockAncestor(headBlock, r.Header.L1Proof) {
					o.l2RollupCountInL1Blocks++
					o.l2RollupTxCountInL1Blocks += len(r.EncryptedTxBlob)
				}

			case *ethadapter.L1DepositTx:
				o.canonicalERC20DepositCount++
			}
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
		"nrBlockParsedERC20Deposits: %d\n",
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
		o.canonicalERC20DepositCount,
	)
}
