package l1

import (
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/host"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
)

var (
	// todo (@matt) make this configurable?
	_timeoutNoBlocks = 30 * time.Second
	one              = big.NewInt(1)
	ErrNoNextBlock   = errors.New("no next block")
)

// Repository is a host service for subscribing to new blocks and looking up L1 data
type Repository struct {
	subscribers []host.L1BlockHandler
	// this eth client should only be used by the repository, the repository may "reconnect" it at any time and don't want to interfere with other processes
	ethClient ethadapter.EthClient
	logger    gethlog.Logger

	running atomic.Bool
	head    gethcommon.Hash
}

func NewL1Repository(ethClient ethadapter.EthClient, logger gethlog.Logger) *Repository {
	return &Repository{
		ethClient: ethClient,
		running:   atomic.Bool{},
		logger:    logger,
	}
}

func (r *Repository) Start() error {
	r.running.Store(true)

	// Repository constantly streams new blocks and forwards them to subscribers
	go r.streamLiveBlocks()
	return nil
}

func (r *Repository) Stop() error {
	r.running.Store(false)
	return nil
}

func (r *Repository) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on last received block or something
	errMsg := ""
	if !r.running.Load() {
		errMsg = "not running"
	}
	return &host.L1RepositoryStatus{ErrMsg: errMsg}
}

// Subscribe will register a new block handler to receive new blocks as they arrive
func (r *Repository) Subscribe(handler host.L1BlockHandler) {
	r.subscribers = append(r.subscribers, handler)
}

// todo (@matt) add unsubscribe (+ mutex on subscribers list), will be needed if enclaves can come and go

func (r *Repository) FetchNextBlock(prevBlock gethcommon.Hash) (*types.Block, bool, error) {
	if prevBlock == r.head {
		// prevBlock is the latest known head
		return nil, false, ErrNoNextBlock
	}
	// the latestCanonAncestor will usually return the prevBlock itself but this step is necessary to walk back if there was a fork
	lca, err := r.latestCanonAncestor(prevBlock)
	if err != nil {
		return nil, false, err
	}
	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := r.ethClient.BlockByNumber(increment(lca.Number()))
	if err != nil {
		return nil, false, fmt.Errorf("could not find block after latest canon ancestor, height=%s - %w", increment(lca.Number()), err)
	}

	return blk, blk.Hash() == r.head, nil
}

func (r *Repository) latestCanonAncestor(blkHash gethcommon.Hash) (*types.Block, error) {
	blk, err := r.ethClient.BlockByHash(blkHash)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block with hash=%s - %w", blkHash, err)
	}
	canonAtSameHeight, err := r.ethClient.BlockByNumber(blk.Number())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block at height=%d - %w", blk.Number(), err)
	}
	if blk.Hash() != canonAtSameHeight.Hash() {
		return r.latestCanonAncestor(blk.ParentHash())
	}
	return blk, nil
}

func (r *Repository) FetchReceipts(block *common.L1Block) types.Receipts {
	receipts := make(types.Receipts, 0)

	for _, transaction := range block.Transactions() {
		receipt, err := r.ethClient.TransactionReceipt(transaction.Hash())

		if err != nil || receipt == nil {
			r.logger.Error("Problem with retrieving the receipt on the host!", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
			continue
		}

		r.logger.Trace("Adding receipt", "status", receipt.Status, log.TxKey, transaction.Hash(),
			log.BlockHashKey, block.Hash(), log.CmpKey, log.CrossChainCmp)

		receipts = append(receipts, receipt)
	}

	return receipts
}

// stream blocks from L1 as they arrive and forward them to subscribers, no guarantee of perfect ordering or that there won't be gaps.
// If streaming is interrupted it will carry on from latest, it won't try to replay missed blocks.
func (r *Repository) streamLiveBlocks() {
	liveStream, streamSub := r.resetLiveStream()
	for r.running.Load() {
		select {
		case header := <-liveStream:
			r.head = header.Hash()
			for _, handler := range r.subscribers {
				block, err := r.ethClient.BlockByHash(header.Hash())
				if err != nil {
					r.logger.Error("error fetching new block", log.ErrKey, err)
				}
				go handler.HandleBlock(block)
			}
		case <-time.After(_timeoutNoBlocks):
			r.logger.Warn("no new blocks received since timeout", "timeout", _timeoutNoBlocks)
			// reset stream to ensure it has not died
			liveStream, streamSub = r.resetLiveStream()
		}
	}

	streamSub.Unsubscribe()
}

func (r *Repository) resetLiveStream() (chan *types.Header, ethereum.Subscription) {
	err := retry.Do(func() error {
		if !r.running.Load() {
			// break out of the loop if repository has stopped
			return retry.FailFast(errors.New("repository is stopped"))
		}
		err := r.ethClient.Reconnect()
		if err != nil {
			r.logger.Warn("failed to reconnect to L1", log.ErrKey, err)
			return err
		}
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{100 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))
	if err != nil {
		// this should only happen if repository has been stopped, because we retry reconnecting forever
		r.logger.Warn("unable to reconnect to L1", log.ErrKey, err)
		return nil, nil
	}
	return r.ethClient.BlockListener()
}

func (r *Repository) FetchBlockByHeight(height int) (*types.Block, error) {
	return r.ethClient.BlockByNumber(big.NewInt(int64(height)))
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
