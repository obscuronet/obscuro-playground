package ethereummock

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/obscuronet/go-obscuro/go/common/async"

	"github.com/google/uuid"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common/gethutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	ethclient_ethereum "github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b common.EncodedL1Block, p common.EncodedL1Block)
	BroadcastTx(tx *types.Transaction)
}

type MiningConfig struct {
	PowTime common.Latency
	LogFile string
}

type TxDB interface {
	Txs(block *types.Block) (map[common.TxHash]*types.Transaction, bool)
	AddTxs(*types.Block, map[common.TxHash]*types.Transaction)
}

type StatsCollector interface {
	// L1Reorg registers when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id gethcommon.Address)
}

type Node struct {
	l2ID     gethcommon.Address // the address of the Obscuro node this client is dedicated to
	cfg      MiningConfig
	Network  L1Network
	mining   bool
	stats    StatsCollector
	Resolver storage.BlockResolver
	db       TxDB
	subs     map[uuid.UUID]*mockSubscription // active subscription for mock blocks
	subMu    sync.Mutex

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *types.Block       // this is where blocks received from peers are dropped
	miningCh    chan *types.Block       // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *types.Block       // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *types.Transaction // where l1 transactions to be published in the next block are added

	// internal
	headInCh         chan bool
	headOutCh        chan *types.Block
	erc20ContractLib erc20contractlib.ERC20ContractLib
	mgmtContractLib  mgmtcontractlib.MgmtContractLib

	logger gethlog.Logger
}

func (m *Node) EstimateGasAndGasPrice(txData types.TxData, _ gethcommon.Address) (types.TxData, error) {
	return txData, nil
}

func (m *Node) SendTransaction(tx *types.Transaction) error {
	m.Network.BroadcastTx(tx)
	return nil
}

func (m *Node) TransactionReceipt(_ gethcommon.Hash) (*types.Receipt, error) {
	// all transactions are immediately processed
	return &types.Receipt{
		Status: types.ReceiptStatusSuccessful,
	}, nil
}

func (m *Node) Nonce(gethcommon.Address) (uint64, error) {
	return 0, nil
}

func (m *Node) getRollupFromBlock(block *types.Block) *common.ExtRollup {
	for _, tx := range block.Transactions() {
		decodedTx := m.mgmtContractLib.DecodeTx(tx)
		if decodedTx == nil {
			continue
		}
		switch l1tx := decodedTx.(type) {
		case *ethadapter.L1RollupTx:
			r, err := common.DecodeRollup(l1tx.Rollup)
			if err == nil {
				return r
			}
		}
	}
	return nil
}

func (m *Node) FetchLastBatchSeqNo(gethcommon.Address) (*big.Int, error) {
	startingBlock, err := m.FetchHeadBlock()
	if err != nil {
		return nil, err
	}

	for currentBlock := startingBlock; currentBlock.NumberU64() != 0; currentBlock, _ = m.BlockByHash(currentBlock.Header().ParentHash) {
		rollup := m.getRollupFromBlock(currentBlock)
		if rollup != nil {
			return big.NewInt(int64(rollup.Header.LastBatchSeqNo)), nil
		}
	}
	return big.NewInt(0), nil
}

// BlockListener provides stream of latest mock head headers as they are created
func (m *Node) BlockListener() (chan *types.Header, ethereum.Subscription) {
	id := uuid.New()
	mockSub := &mockSubscription{
		node:   m,
		id:     id,
		headCh: make(chan *types.Header),
	}
	m.subMu.Lock()
	defer m.subMu.Unlock()
	m.subs[id] = mockSub
	return mockSub.headCh, mockSub
}

func (m *Node) BlockNumber() (uint64, error) {
	blk, err := m.Resolver.FetchHeadBlock()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return 0, ethereum.NotFound
		}
		return 0, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return blk.NumberU64(), nil
}

func (m *Node) BlockByNumber(n *big.Int) (*types.Block, error) {
	if n.Int64() == 0 {
		return MockGenesisBlock, nil
	}
	// TODO this should be a method in the resolver
	blk, err := m.Resolver.FetchHeadBlock()
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, ethereum.NotFound
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	for !bytes.Equal(blk.ParentHash().Bytes(), (common.L1BlockHash{}).Bytes()) {
		if blk.NumberU64() == n.Uint64() {
			return blk, nil
		}

		blk, err = m.Resolver.FetchBlock(blk.ParentHash())
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent for block in chain. Cause: %w", err)
		}
	}
	return nil, ethereum.NotFound
}

func (m *Node) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	blk, err := m.Resolver.FetchBlock(id)
	if err != nil {
		return nil, fmt.Errorf("block could not be retrieved. Cause: %w", err)
	}
	return blk, nil
}

func (m *Node) FetchHeadBlock() (*types.Block, error) {
	block, err := m.Resolver.FetchHeadBlock()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return block, nil
}

func (m *Node) Info() ethadapter.Info {
	return ethadapter.Info{
		L2ID: m.l2ID,
	}
}

func (m *Node) IsBlockAncestor(block *types.Block, proof common.L1BlockHash) bool {
	return m.Resolver.IsBlockAncestor(block, proof)
}

func (m *Node) BalanceAt(gethcommon.Address, *big.Int) (*big.Int, error) {
	panic("not implemented")
}

// Start runs an infinite loop that listens to the two block producing channels and processes them.
func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	err := m.Resolver.StoreBlock(MockGenesisBlock, nil, nil)
	if err != nil {
		m.logger.Crit("Failed to store block")
	}
	head := m.setHead(MockGenesisBlock)

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, err := m.Resolver.FetchBlock(p2pb.Hash())
			// only process blocks if they haven't been processed before
			if err != nil {
				if errors.Is(err, errutil.ErrNotFound) {
					head = m.processBlock(p2pb, head)
				} else {
					panic(fmt.Errorf("could not retrieve parent block. Cause: %w", err))
				}
			}

		case mb := <-m.miningCh: // Received from the local mining
			head = m.processBlock(mb, head)
			if bytes.Equal(head.Hash().Bytes(), mb.Hash().Bytes()) { // Ignore the locally produced block if someone else found one already
				p, err := m.Resolver.FetchBlock(mb.ParentHash())
				if err != nil {
					panic(fmt.Errorf("could not retrieve parent. Cause: %w", err))
				}
				encodedBlock, err := common.EncodeBlock(mb)
				if err != nil {
					panic(fmt.Errorf("could not encode block. Cause: %w", err))
				}
				encodedParentBlock, err := common.EncodeBlock(p)
				if err != nil {
					panic(fmt.Errorf("could not encode parent block. Cause: %w", err))
				}
				m.Network.BroadcastBlock(encodedBlock, encodedParentBlock)
			}
		case <-m.headInCh:
			m.headOutCh <- head
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *types.Block, head *types.Block) *types.Block {
	err := m.Resolver.StoreBlock(b, nil, nil)
	if err != nil {
		m.logger.Crit("Failed to store block. Cause: %w", err)
	}
	_, err = m.Resolver.FetchBlock(b.Header().ParentHash)
	// only proceed if the parent is available
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			m.logger.Info(fmt.Sprintf("Parent block not found=b_%d", common.ShortHash(b.Header().ParentHash)))
			return head
		}
		m.logger.Crit("Could not fetch block parent. Cause: %w", err)
	}

	// Ignore superseded blocks
	if b.NumberU64() <= head.NumberU64() {
		return head
	}

	// Check for Reorgs
	if !m.Resolver.IsAncestor(b, head) {
		m.stats.L1Reorg(m.l2ID)
		fork, _, _, err := gethutil.LCA(head, b, m.Resolver)
		if err != nil {
			panic(err)
		}
		m.logger.Info(
			fmt.Sprintf("L1Reorg new=b_%d(%d), old=b_%d(%d), fork=b_%d(%d)", common.ShortHash(b.Hash()), b.NumberU64(), common.ShortHash(head.Hash()), head.NumberU64(), common.ShortHash(fork.Hash()), fork.NumberU64()))
		return m.setFork(m.BlocksBetween(fork, b))
	}

	if b.NumberU64() > (head.NumberU64() + 1) {
		m.logger.Crit("Should not happen")
	}

	return m.setHead(b)
}

// Notifies the Miner to start mining on the new block and the aggregator to produce rollups
func (m *Node) setHead(b *types.Block) *types.Block {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return b
	}

	// notify the client subscriptions
	m.subMu.Lock()
	for _, s := range m.subs {
		sub := s
		go sub.publish(b)
	}
	m.subMu.Unlock()
	m.canonicalCh <- b

	return b
}

func (m *Node) setFork(blocks []*types.Block) *types.Block {
	head := blocks[len(blocks)-1]
	if atomic.LoadInt32(m.interrupt) == 1 {
		return head
	}

	// notify the client subs
	m.subMu.Lock()
	for _, s := range m.subs {
		sub := s
		go sub.publishAll(blocks)
	}
	m.subMu.Unlock()

	m.canonicalCh <- head

	return head
}

// P2PReceiveBlock is called by counterparties when there is a block to broadcast
// All it does is drop the blocks in a channel for processing.
func (m *Node) P2PReceiveBlock(b common.EncodedL1Block, p common.EncodedL1Block) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}
	decodedBlock, err := b.DecodeBlock()
	if err != nil {
		panic(fmt.Errorf("could not decode block. Cause: %w", err))
	}
	decodedParentBlock, err := p.DecodeBlock()
	if err != nil {
		panic(fmt.Errorf("could not decode parent block. Cause: %w", err))
	}
	m.p2pCh <- decodedParentBlock
	m.p2pCh <- decodedBlock
}

// startMining - listens on the canonicalCh and schedule a go routine that produces a block after a PowTime and drop it
// on the miningCh channel
func (m *Node) startMining() {
	m.logger.Info(" starting miner...")
	// stores all transactions seen from the beginning of time.
	mempool := make([]*types.Transaction, 0)
	z := int32(0)
	interrupt := &z

	for {
		select {
		case <-m.exitMiningCh:
			return
		case tx := <-m.mempoolCh:
			mempool = append(mempool, tx)

		case canonicalBlock := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.

			// remove transactions that are already considered committed
			mempool = m.removeCommittedTransactions(canonicalBlock, mempool, m.Resolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			async.Schedule(m.cfg.PowTime(), func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.Resolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}

				m.miningCh <- NewBlock(canonicalBlock, m.l2ID, toInclude)
			})
		}
	}
}

// P2PGossipTx receive rollups to publish from the linked aggregators
func (m *Node) P2PGossipTx(tx *types.Transaction) {
	if atomic.LoadInt32(m.interrupt) == 1 {
		return
	}

	m.mempoolCh <- tx
}

func (m *Node) BroadcastTx(tx types.TxData) {
	m.Network.BroadcastTx(types.NewTx(tx))
}

func (m *Node) Stop() {
	// block all requests
	atomic.StoreInt32(m.interrupt, 1)
	time.Sleep(time.Millisecond * 100)

	m.exitMiningCh <- true
	m.exitCh <- true
}

func (m *Node) BlocksBetween(blockA *types.Block, blockB *types.Block) []*types.Block {
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return []*types.Block{blockA}
	}
	blocks := make([]*types.Block, 0)
	tempBlock := blockB
	var err error
	for {
		blocks = append(blocks, tempBlock)
		if bytes.Equal(tempBlock.Hash().Bytes(), blockA.Hash().Bytes()) {
			break
		}
		tempBlock, err = m.Resolver.FetchBlock(tempBlock.ParentHash())
		if err != nil {
			panic(fmt.Errorf("could not retrieve parent block. Cause: %w", err))
		}
	}
	n := len(blocks)
	result := make([]*types.Block, n)
	for i, block := range blocks {
		result[n-i-1] = block
	}
	return result
}

func (m *Node) CallContract(ethereum.CallMsg) ([]byte, error) {
	return nil, nil
}

func (m *Node) EthClient() *ethclient_ethereum.Client {
	return nil
}

func (m *Node) RemoveSubscription(id uuid.UUID) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	delete(m.subs, id)
}

func (m *Node) Reconnect() error {
	return nil
}

func (m *Node) Alive() bool {
	return true
}

func NewMiner(
	id gethcommon.Address,
	cfg MiningConfig,
	network L1Network,
	statsCollector StatsCollector,
) *Node {
	return &Node{
		l2ID:             id,
		mining:           true,
		cfg:              cfg,
		stats:            statsCollector,
		Resolver:         NewResolver(),
		db:               NewTxDB(),
		Network:          network,
		exitCh:           make(chan bool),
		exitMiningCh:     make(chan bool),
		interrupt:        new(int32),
		p2pCh:            make(chan *types.Block),
		miningCh:         make(chan *types.Block),
		canonicalCh:      make(chan *types.Block),
		mempoolCh:        make(chan *types.Transaction),
		headInCh:         make(chan bool),
		headOutCh:        make(chan *types.Block),
		erc20ContractLib: NewERC20ContractLibMock(),
		mgmtContractLib:  NewMgmtContractLibMock(),
		logger:           log.New(log.EthereumL1Cmp, int(gethlog.LvlInfo), cfg.LogFile, log.NodeIDKey, id),
		subs:             map[uuid.UUID]*mockSubscription{},
		subMu:            sync.Mutex{},
	}
}

// implements the ethereum.Subscription
type mockSubscription struct {
	id     uuid.UUID
	headCh chan *types.Header
	node   *Node // we hold a reference to the node to unsubscribe ourselves - not ideal but this is just a mock
}

func (sub *mockSubscription) Err() <-chan error {
	return make(chan error)
}

func (sub *mockSubscription) Unsubscribe() {
	sub.node.RemoveSubscription(sub.id)
}

func (sub *mockSubscription) publish(b *types.Block) {
	sub.headCh <- b.Header()
}

func (sub *mockSubscription) publishAll(blocks []*types.Block) {
	for _, b := range blocks {
		sub.publish(b)
	}
}
