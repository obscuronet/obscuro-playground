package ethereummock

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/common/async"

	"github.com/google/uuid"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	ethclient_ethereum "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

type (
	Latency func() time.Duration
)

type L1Network interface {
	// BroadcastBlock - send the block and the parent to make sure there are no gaps
	BroadcastBlock(b EncodedL1Block, p EncodedL1Block)
	BroadcastTx(tx *types.Transaction)
}

type MiningConfig struct {
	PowTime      Latency
	LogFile      string
	L1BeaconPort int
}

type TxDB interface {
	Txs(block *types.Block) (map[common.TxHash]*types.Transaction, bool)
	AddTxs(*types.Block, map[common.TxHash]*types.Transaction)
}

type StatsCollector interface {
	// L1Reorg registers when a miner has to process a reorg (a winning block from a fork)
	L1Reorg(id gethcommon.Address)
}

type BlockWithBlobs struct {
	Block *types.Block
	Blobs []*kzg4844.Blob
}

type Node struct {
	l2ID          gethcommon.Address // the address of the Obscuro node this client is dedicated to
	cfg           MiningConfig
	Network       L1Network
	mining        bool
	stats         StatsCollector
	BlockResolver *blockResolverInMem
	BlobResolver  l1.BlobResolver
	db            TxDB
	subs          map[uuid.UUID]*mockSubscription // active subscription for mock blocks
	subMu         sync.Mutex

	// Channels
	exitCh       chan bool // the Node stops
	exitMiningCh chan bool // the mining loop is notified to stop
	interrupt    *int32

	p2pCh       chan *types.Block       // this is where blocks received from peers are dropped
	miningCh    chan *types.Block       // this is where blocks created by the mining setup of the current node are dropped
	canonicalCh chan *types.Header      // this is where the main processing routine drops blocks that are canonical
	mempoolCh   chan *types.Transaction // where l1 transactions to be published in the next block are added

	// internal
	erc20ContractLib erc20contractlib.ERC20ContractLib
	mgmtContractLib  mgmtcontractlib.MgmtContractLib

	logger gethlog.Logger

	// this mock state is to simulate the permissioning of the sequencer enclave, the L1 now 'knows the seq enclave ID'
	tenSeqEnclaveID common.EnclaveID
}

func (m *Node) PrepareTransactionToSend(_ context.Context, txData types.TxData, _ gethcommon.Address) (types.TxData, error) {
	switch tx := txData.(type) {
	case *types.LegacyTx:
		return createLegacyTx(txData)
	case *types.BlobTx:
		return createBlobTx(txData)
	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", tx)
	}
}

func (m *Node) PrepareTransactionToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, _ uint64, _ int) (types.TxData, error) {
	return m.PrepareTransactionToSend(ctx, txData, from)
}

func createLegacyTx(txData types.TxData) (types.TxData, error) {
	tx := types.NewTx(txData)
	return &types.LegacyTx{
		Nonce:    123,
		GasPrice: tx.GasPrice(),
		Gas:      tx.Gas(),
		To:       tx.To(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil
}

func createBlobTx(txData types.TxData) (types.TxData, error) {
	tx := types.NewTx(txData)
	return &types.BlobTx{
		To:         *tx.To(),
		Data:       tx.Data(),
		BlobHashes: tx.BlobHashes(),
		Sidecar:    tx.BlobTxSidecar(),
	}, nil
}

func (m *Node) SendTransaction(tx *types.Transaction) error {
	m.Network.BroadcastTx(tx)
	return nil
}

func (m *Node) TransactionReceipt(_ gethcommon.Hash) (*types.Receipt, error) {
	// all transactions are immediately processed
	return &types.Receipt{
		BlockNumber: big.NewInt(1),
		Status:      types.ReceiptStatusSuccessful,
	}, nil
}

func (m *Node) TransactionByHash(hash gethcommon.Hash) (*types.Transaction, bool, error) {
	// First check mempool/pending transactions
	select {
	case tx := <-m.mempoolCh:
		if tx.Hash() == hash {
			// Put the tx back in the channel
			m.mempoolCh <- tx
			return tx, false, nil
		}
		// Put the tx back in the channel
		m.mempoolCh <- tx
	default:
		// Don't block if mempool channel is empty
	}

	// Then check if the transaction exists in any block
	blk, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		return nil, false, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}

	// Traverse the chain looking for the transaction
	for !bytes.Equal(blk.ParentHash().Bytes(), (common.L1BlockHash{}).Bytes()) {
		for _, tx := range blk.Transactions() {
			if tx.Hash() == hash {
				return tx, true, nil
			}
		}

		blk, err = m.BlockResolver.FetchFullBlock(context.Background(), blk.ParentHash())
		if err != nil {
			return nil, false, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
		}
	}

	// Check genesis block
	for _, tx := range MockGenesisBlock.Transactions() {
		if tx.Hash() == hash {
			return tx, true, nil
		}
	}

	// If we get here, the transaction wasn't found
	return nil, false, ethereum.NotFound
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
		case *common.L1RollupHashes:
			ctx := context.TODO()
			blobs, _ := m.BlobResolver.FetchBlobs(ctx, block.Header(), l1tx.BlobHashes)
			r, err := ethadapter.ReconstructRollup(blobs)
			if err != nil {
				m.logger.Error("could not recreate rollup from blobs. Cause: %w", err)
				return nil
			}
			return r
		}
	}
	return nil
}

func (m *Node) FetchLastBatchSeqNo(gethcommon.Address) (*big.Int, error) {
	startingBlock, err := m.FetchHeadBlock()
	if err != nil {
		return nil, err
	}

	for currentBlock := startingBlock; currentBlock.Number.Uint64() != 0; {
		cb, err := m.BlockByHash(currentBlock.ParentHash)
		if err != nil {
			m.logger.Error("Error fetching block by hash", "error", err)
			break
		}
		rollup := m.getRollupFromBlock(cb)
		if rollup != nil {
			return big.NewInt(int64(rollup.Header.LastBatchSeqNo)), nil
		}
		currentBlock = cb.Header()
	}
	// the first batch is number 1
	return big.NewInt(int64(common.L2GenesisSeqNo)), nil
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
	blk, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return 0, ethereum.NotFound
		}
		return 0, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return blk.NumberU64(), nil
}

func (m *Node) HeaderByNumber(n *big.Int) (*types.Header, error) {
	if n == nil {
		return m.FetchHeadBlock()
	}
	if n.Int64() == 0 {
		return MockGenesisBlock.Header(), nil
	}
	blk, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, ethereum.NotFound
		}
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	for !bytes.Equal(blk.ParentHash().Bytes(), (common.L1BlockHash{}).Bytes()) {
		if blk.NumberU64() == n.Uint64() {
			return blk.Header(), nil
		}

		blk, err = m.BlockResolver.FetchFullBlock(context.Background(), blk.ParentHash())
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent for block in chain. Cause: %w", err)
		}
	}
	return nil, ethereum.NotFound
}

func (m *Node) HeaderByHash(id gethcommon.Hash) (*types.Header, error) {
	blk, err := m.BlockResolver.FetchFullBlock(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("block could not be retrieved. Cause: %w", err)
	}
	return blk.Header(), nil
}

func (m *Node) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	blk, err := m.BlockResolver.FetchFullBlock(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("block could not be retrieved. Cause: %w", err)
	}
	return blk, nil
}

func (m *Node) FetchHeadBlock() (*types.Header, error) {
	block, err := m.BlockResolver.FetchHeadBlock(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	return block.Header(), nil
}

func (m *Node) Info() ethadapter.Info {
	return ethadapter.Info{
		L2ID: m.l2ID,
	}
}

func (m *Node) IsBlockAncestor(block *types.Header, proof common.L1BlockHash) bool {
	return m.BlockResolver.IsBlockAncestor(context.Background(), block, proof)
}

func (m *Node) BalanceAt(gethcommon.Address, *big.Int) (*big.Int, error) {
	panic("not implemented")
}

// GetLogs is a mock method - we create logs with topics matching the real contract events
func (m *Node) GetLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	logs := make([]types.Log, 0)
	if fq.BlockHash == nil {
		return logs, nil
	}

	blk, err := m.BlockByHash(*fq.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve block. Cause: %w", err)
	}

	for _, tx := range blk.Transactions() {
		if tx.To() == nil {
			continue
		}

		// map transaction types to their corresponding event topics
		var topic gethcommon.Hash
		var data []byte
		switch tx.To().Hex() {
		case rollupTxAddr.Hex():
			topic = crosschain.RollupAddedID
		case messageBusAddr.Hex():
			topic = crosschain.CrossChainEventID
		case depositTxAddr.Hex():
			topic = crosschain.ValueTransferEventID
		case storeSecretTxAddr.Hex():
			topic = crosschain.NetworkSecretRespondedID
		case requestSecretTxAddr.Hex():
			topic = crosschain.NetworkSecretRequestedID
		case initializeSecretTxAddr.Hex():
			topic = crosschain.SequencerEnclaveGrantedEventID
		case grantSeqTxAddr.Hex():
			topic = crosschain.SequencerEnclaveGrantedEventID
			// enclave ID address, padded out to 32 bytes to match standard eth fields
			data = make([]byte, 32)
			copy(data[12:], m.tenSeqEnclaveID[:])
		default:
			continue
		}

		dummyLog := types.Log{
			Address:     *tx.To(),
			BlockHash:   blk.Hash(),
			TxHash:      tx.Hash(),
			Topics:      []gethcommon.Hash{topic},
			BlockNumber: blk.NumberU64(),
			Index:       uint(len(logs)),
			Data:        data,
		}
		logs = append(logs, dummyLog)
	}
	return logs, nil
}

func (m *Node) Start() {
	if m.mining {
		// This starts the mining
		go m.startMining()
	}

	err := m.BlockResolver.StoreBlock(context.Background(), MockGenesisBlock, nil)
	if err != nil {
		m.logger.Crit("Failed to store block")
	}
	head := m.setHead(MockGenesisBlock.Header())

	for {
		select {
		case p2pb := <-m.p2pCh: // Received from peers
			_, err := m.BlockResolver.FetchFullBlock(context.Background(), p2pb.Hash())
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
			if bytes.Equal(head.Hash().Bytes(), mb.Hash().Bytes()) { // Only broadcast if it's the new head
				p, err := m.BlockResolver.FetchFullBlock(context.Background(), mb.ParentHash())
				if err != nil {
					panic(fmt.Errorf("could not retrieve parent. Cause: %w", err))
				}
				encodedBlock, err := EncodeBlock(mb)
				if err != nil {
					panic(fmt.Errorf("could not encode block. Cause: %w", err))
				}
				encodedParentBlock, err := EncodeBlock(p)
				if err != nil {
					panic(fmt.Errorf("could not encode parent block. Cause: %w", err))
				}
				m.Network.BroadcastBlock(encodedBlock, encodedParentBlock)
			}
		case <-m.exitCh:
			return
		}
	}
}

func (m *Node) processBlock(b *types.Block, head *types.Header) *types.Header {
	err := m.BlockResolver.StoreBlock(context.Background(), b, nil)
	if err != nil {
		m.logger.Crit("Failed to store block. Cause: %w", err)
	}

	_, err = m.BlockResolver.FetchFullBlock(context.Background(), b.Header().ParentHash)
	// only proceed if the parent is available
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			m.logger.Info(fmt.Sprintf("Parent block not found=b_%s", b.Header().ParentHash))
			return head
		}
		m.logger.Crit("Could not fetch block parent. Cause: %w", err)
	}

	// Ignore superseded blocks
	if b.NumberU64() <= head.Number.Uint64() {
		return head
	}

	// Check for Reorgs
	if !m.BlockResolver.IsAncestor(context.Background(), b.Header(), head) {
		m.stats.L1Reorg(m.l2ID)
		fork, err := gethutil.LCA(context.Background(), head, b.Header(), m.BlockResolver)
		if err != nil {
			m.logger.Error("Should not happen.", log.ErrKey, err)
			return head
		}
		m.logger.Info(
			fmt.Sprintf("L1Reorg new=b_%s(%d), old=b_%s(%d), fork=b_%s(%d)", b.Hash(), b.NumberU64(), head.Hash(), head.Number.Uint64(), fork.CommonAncestor.Hash(), fork.CommonAncestor.Number.Uint64()))
		return m.setFork(m.BlocksBetween(fork.CommonAncestor, b.Header()))
	}
	if b.NumberU64() > (head.Number.Uint64() + 1) {
		m.logger.Error("Should not happen. Blocks are skewed")
	}

	return m.setHead(b.Header())
}

// Notifies the Miner to start mining on the new block and the aggregator to produce rollups
func (m *Node) setHead(b *types.Header) *types.Header {
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

func (m *Node) setFork(blocks []*types.Header) *types.Header {
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
func (m *Node) P2PReceiveBlock(b EncodedL1Block, p EncodedL1Block) {
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

		case canonicalBlockHeader := <-m.canonicalCh:
			// A new canonical block was found. Start a new round based on that block.
			canonicalBlock, err := m.BlockResolver.FetchFullBlock(context.Background(), canonicalBlockHeader.Hash())
			if err != nil {
				panic(fmt.Errorf("could not fetch block. Cause: %w", err))
			}

			// remove transactions that are already considered committed
			mempool = m.removeCommittedTransactions(context.Background(), canonicalBlock, mempool, m.BlockResolver, m.db)

			// notify the existing mining go routine to stop mining
			atomic.StoreInt32(interrupt, 1)
			c := int32(0)
			interrupt = &c

			// Generate a random number, and wait for that number of ms. Equivalent to PoW
			// Include all rollups received during this period.
			blockTime := uint64(time.Now().Unix())
			async.Schedule(m.cfg.PowTime(), func() {
				toInclude := findNotIncludedTxs(canonicalBlock, mempool, m.BlockResolver, m.db)
				// todo - iterate through the rollup transactions and include only the ones with the proof on the canonical chain
				if atomic.LoadInt32(m.interrupt) == 1 {
					return
				}
				b := NewBlock(canonicalBlock, m.l2ID, toInclude, blockTime)
				// there is a race condition if we process this at the same time as the blocks, so it has to be placed here
				err := m.ProcessBlobs(b)
				if err != nil {
					m.logger.Crit("Failed to store blobs. Cause: %w", err)
				}
				m.miningCh <- NewBlock(canonicalBlock, m.l2ID, toInclude, blockTime)
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

func (m *Node) BlocksBetween(blockA *types.Header, blockB *types.Header) []*types.Header {
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return []*types.Header{blockB}
	}
	blocks := make([]*types.Header, 0)
	tempBlock := blockB
	for {
		blocks = append(blocks, tempBlock)
		if bytes.Equal(tempBlock.Hash().Bytes(), blockA.Hash().Bytes()) {
			break
		}
		tb, err := m.BlockResolver.FetchFullBlock(context.Background(), tempBlock.ParentHash)
		if err != nil {
			panic(fmt.Errorf("could not retrieve parent block. Cause: %w", err))
		}
		tempBlock = tb.Header()
	}
	n := len(blocks)
	result := make([]*types.Header, n)
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

func (m *Node) ReconnectIfClosed() error {
	return nil
}

func (m *Node) Alive() bool {
	return true
}

func (m *Node) ProcessBlobs(b *types.Block) error {
	var blobs []*kzg4844.Blob
	for _, tx := range b.Transactions() {
		if tx.BlobHashes() != nil {
			for i := range tx.BlobTxSidecar().Blobs {
				blobPtr := &tx.BlobTxSidecar().Blobs[i]
				blobs = append(blobs, blobPtr)
			}
		}
	}

	blobPointers := make([]*kzg4844.Blob, len(blobs))
	copy(blobPointers, blobs)

	if len(blobs) > 0 {
		// dummy slot to simplify the in memory blob storage
		err := m.BlobResolver.StoreBlobs(0, blobs)
		if err != nil {
			return fmt.Errorf("could not store blobs. Cause: %w", err)
		}
	}
	return nil
}

func (m *Node) PromoteEnclave(id common.EnclaveID) {
	m.tenSeqEnclaveID = id
}

func NewMiner(
	id gethcommon.Address,
	cfg MiningConfig,
	network L1Network,
	statsCollector StatsCollector,
	blobResolver l1.BlobResolver,
	logger gethlog.Logger,
) *Node {
	return &Node{
		l2ID:             id,
		mining:           true,
		cfg:              cfg,
		stats:            statsCollector,
		BlockResolver:    NewResolver(),
		BlobResolver:     blobResolver,
		db:               NewTxDB(),
		Network:          network,
		exitCh:           make(chan bool),
		exitMiningCh:     make(chan bool),
		interrupt:        new(int32),
		p2pCh:            make(chan *types.Block),
		miningCh:         make(chan *types.Block),
		canonicalCh:      make(chan *types.Header),
		mempoolCh:        make(chan *types.Transaction),
		erc20ContractLib: NewERC20ContractLibMock(),
		mgmtContractLib:  NewMgmtContractLibMock(),
		logger:           logger,
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

func (sub *mockSubscription) publish(b *types.Header) {
	if sub.headCh != nil {
		sub.headCh <- b
	}
}

func (sub *mockSubscription) publishAll(blocks []*types.Header) {
	for _, b := range blocks {
		sub.publish(b)
	}
}
