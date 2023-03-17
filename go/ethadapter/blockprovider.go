package ethadapter

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/retry"

	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

const (
	waitingForBlockTimeout = 30 * time.Second
)

var (
	one                   = big.NewInt(1)
	backoffRetryIntervals = []time.Duration{10 * time.Millisecond, 100 * time.Millisecond, 1 * time.Second, 10 * time.Second}
)

func NewEthBlockProvider(ethClient EthClient, logger gethlog.Logger) *EthBlockProvider {
	return &EthBlockProvider{
		ethClient: ethClient,
		logger:    logger,
	}
}

// EthBlockProvider streams blocks from the ethereum L1 client in the order expected by the enclave (consecutive, canonical blocks)
type EthBlockProvider struct {
	ethClient EthClient
	logger    gethlog.Logger
}

// StartStreamingFromHash will look up the hash block, find the appropriate height (LCA if there have been forks) and
// then call StartStreamingFromHeight from there
func (e *EthBlockProvider) StartStreamingFromHash(latestHash gethcommon.Hash) (*host.BlockStream, error) {
	ancestorBlk, err := e.latestCanonAncestor(latestHash)
	if err != nil {
		return nil, err
	}
	return e.StartStreamingFromHeight(ancestorBlk.Number())
}

// StartStreamingFromHeight will start streaming from the given height
// returning the fresh channel (the next block will be the requested height) and a cancel function to kill the stream
func (e *EthBlockProvider) StartStreamingFromHeight(height *big.Int) (*host.BlockStream, error) {
	// block heights start at 1
	if height.Cmp(one) < 0 {
		height = one
	}
	ctx, cancel := context.WithCancel(context.Background())
	streamCh := make(chan *types.Block)
	go e.streamBlocks(ctx, height, streamCh)

	return &host.BlockStream{Stream: streamCh, Stop: cancel}, nil
}

// IsLatest returns true iff the block has the same hash as the L1 head block (from the eth client)
func (e *EthBlockProvider) IsLatest(b *types.Block) bool {
	l1Head, err := e.ethClient.FetchHeadBlock()
	if err != nil {
		e.logger.Warn("unable to fetch head eth block - %w", err)
		return false
	}
	isLatest := b.Hash() == l1Head.Hash()
	// this log message is helpful for visibility on how far behind the block feeding is
	e.logger.Info("L1 block provider live-monitoring", "currBlock", b.NumberU64(), "head", l1Head.NumberU64(), "isLatest", isLatest)
	return isLatest
}

// streamBlocks is the main loop. It should be run in a separate go routine. It will stream catch-up blocks from requested height until it
// reaches the latest live block, then it will block until next live block arrives
// It blocks when:
// - publishing a block, it blocks on the outbound channel until the block is consumed
// - awaiting a live block, when consumer is completely up-to-date it waits for a live block to arrive
func (e *EthBlockProvider) streamBlocks(ctx context.Context, fromHeight *big.Int, streamCh chan *types.Block) {
	var latestSent *types.Header // most recently sent block (reset to nil when a `StartStreaming...` method is called)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var block *types.Block
			// this will retry forever with a back-off, if L1 client is unavailable we should still recover eventually
			err := retry.Do(func() error {
				var fetchErr error
				// this will block if we're up-to-date with live blocks
				block, fetchErr = e.fetchNextCanonicalBlock(ctx, fromHeight, latestSent)
				if fetchErr != nil {
					// this shouldn't happen often, it's important that node operator has visibility on it, and that host can
					// eventually recover when L1 client issue is resolved
					e.logger.Warn("unexpected error while preparing block to stream, will retry periodically", log.ErrKey, fetchErr)
					return fetchErr
				}
				return nil
			}, retry.NewBackoffAndRetryForeverStrategy(backoffRetryIntervals, 30*time.Second))
			if err != nil {
				// retry block above should retry forever until success, so if it does error we treat it as unrecoverable
				e.logger.Crit("unable to fetch next canonical block, unexpected failure", log.ErrKey, err)
			}
			e.logger.Trace("blockProvider streaming block", "height", block.Number(), "hash", block.Hash())
			streamCh <- block // we block here until consumer takes it
			// update stream state
			latestSent = block.Header()
		}
	}
}

// fetchNextCanonicalBlock finds the next block to send downstream to the consumer.
// It looks at:
//   - the latest block that was sent to the consumer (`latestSent`)
//   - the current head of the L1 according to the Eth client (`l1Block`)
//
// If the consumer is up-to-date: this method will wait until a new block arrives from the L1
// If the consumer is behind or there has been a fork: it returns the next canonical block that the consumer needs to see (no waiting)
func (e *EthBlockProvider) fetchNextCanonicalBlock(ctx context.Context, fromHeight *big.Int, latestSent *types.Header) (*types.Block, error) {
	// check for the initial case, when consumer has first requested a stream
	if latestSent == nil {
		// we try to return the canonical block at the height requested
		blk, err := e.ethClient.BlockByNumber(fromHeight)
		if err != nil {
			return nil, fmt.Errorf("could not fetch block at requested height=%d - %w", fromHeight, err)
		}
		return blk, nil
	}

	l1Block, err := e.ethClient.FetchHeadBlock()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head block. Cause: %w", err)
	}
	l1Head := l1Block.Header()

	// check if the block provider is already up-to-date with the latest L1 block (if it is, then wait for a new block)
	if l1Head.Hash() == latestSent.Hash() {
		var err error
		// block provider's current head is up-to-date, we will wait for the next L1 block to arrive
		l1Head, err = e.awaitNewBlock(ctx)
		if err != nil {
			return nil, fmt.Errorf("no new block found from stream - %w", err)
		}
	}

	// most common path should be: new head block that arrived is the next block, and needs to be sent
	if l1Head.ParentHash == latestSent.Hash() {
		blk, err := e.ethClient.BlockByHash(l1Head.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not fetch block with hash=%s - %w", l1Head.Hash(), err)
		}
		return blk, nil
	}

	// and if not then, we walk back up the tree from the current head, to find the most recent **canonical** block we sent
	latestCanon, err := e.latestCanonAncestor(latestSent.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not find ancestor on canonical chain for hash=%s - %w", latestSent.Hash(), err)
	}

	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := e.ethClient.BlockByNumber(increment(latestCanon.Number()))
	if err != nil {
		return nil, fmt.Errorf("could not find block after canon fork branch, height=%s - %w", increment(latestCanon.Number()), err)
	}
	return blk, nil
}

func (e *EthBlockProvider) latestCanonAncestor(blkHash gethcommon.Hash) (*types.Block, error) {
	blk, err := e.ethClient.BlockByHash(blkHash)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block with hash=%s - %w", blkHash, err)
	}
	canonAtSameHeight, err := e.ethClient.BlockByNumber(blk.Number())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block at height=%d - %w", blk.Number(), err)
	}
	if blk.Hash() != canonAtSameHeight.Hash() {
		return e.latestCanonAncestor(blk.ParentHash())
	}
	return blk, nil
}

// awaitNewBlock will block, waiting for the next live L1 block from the L1 client.
// It is used when the block provider is up-to-date and waiting for a new block to forward to its consumer.
func (e *EthBlockProvider) awaitNewBlock(ctx context.Context) (*types.Header, error) {
	liveStream, streamSub := e.ethClient.BlockListener()

	for {
		select {
		case blkHead := <-liveStream:
			e.logger.Trace("received new L1 head", "height", blkHead.Number, "hash", blkHead.Hash())
			streamSub.Unsubscribe()
			return blkHead, nil

		case <-ctx.Done():
			return nil, fmt.Errorf("context closed before block was received")

		case err := <-streamSub.Err():
			if errors.Is(err, ErrSubscriptionNotSupported) {
				return nil, err
			}
			e.logger.Error("L1 block monitoring error", log.ErrKey, err)

			e.logger.Info("Restarting L1 block Monitoring...")
			liveStream, streamSub = e.ethClient.BlockListener()

		case <-time.After(waitingForBlockTimeout):
			return nil, fmt.Errorf("no block received from L1 client stream for over %s", waitingForBlockTimeout)
		}
	}
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
