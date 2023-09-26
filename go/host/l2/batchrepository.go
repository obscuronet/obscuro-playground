package l2

import (
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

const (
	// if request asks for batches from seq no. X we don't want to return potentially thousands of batches, so we limit
	// the number of batches we return with this cap
	// (recipient will request the next ones as required, and they should be catching up from roll-ups first)
	_maxBatchesInP2PResponse      = 50
	_timeoutWaitingForP2PResponse = 30 * time.Second
)

// This private interface enforces the services that the guardian depends on
type batchRepoServiceLocator interface {
	P2P() host.P2P
	Enclaves() host.EnclaveService
}

// Repository is responsible for storing and retrieving batches from the database
// If it can't find a batch it will request it from peers. It also subscribes for batch requests from peers and responds to them.
type Repository struct {
	subscribers []host.L2BatchHandler

	sl          batchRepoServiceLocator
	db          *db.DB
	isSequencer bool

	// high watermark for batch sequence numbers seen so far. If we can't find batch for seq no < this, then we should ask peers for missing batches
	latestBatchSeqNo *big.Int
	latestSeqNoMutex sync.Mutex

	// The repository requests batches from peers asynchronously, we don't want to repeatedly spam out requests if we
	// haven't received a response yet, but we also don't want to wait forever if there's no response.
	// So we keep track of the last request time and what was requested, using a mutex to avoid concurrent access errors on them
	p2pReqMutex          sync.Mutex
	p2pInFlightRequested *big.Int
	p2pInFlightReqTime   *time.Time

	running atomic.Bool
	logger  gethlog.Logger
}

func NewBatchRepository(cfg *config.HostConfig, hostService batchRepoServiceLocator, database *db.DB, logger gethlog.Logger) *Repository {
	return &Repository{
		sl:               hostService,
		db:               database,
		isSequencer:      cfg.NodeType == common.Sequencer,
		latestBatchSeqNo: big.NewInt(0),
		running:          atomic.Bool{},
		logger:           logger,
	}
}

func (r *Repository) Start() error {
	r.running.Store(true)

	// register ourselves for new batches from p2p
	r.sl.P2P().SubscribeForBatches(r)
	r.sl.P2P().SubscribeForBatchRequests(r)

	return nil
}

func (r *Repository) Stop() error {
	r.running.Store(false)
	return nil
}

func (r *Repository) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on last received batch or something
	errMsg := ""
	if !r.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

// HandleBatches receives new batches from the p2p network, it also handles batches that are requested from peers
// If the batch is the new head of the L2 then it notifies subscribers to this service that a new batch has arrived
func (r *Repository) HandleBatches(batches []*common.ExtBatch, isLive bool) {
	// if these batches resolve the in-flight request we made then clear the in-flight request (see type def for details)
	r.p2pReqMutex.Lock()
	if !isLive && len(batches) > 0 && r.p2pInFlightRequested != nil && batches[0].Header.SequencerOrderNo.Cmp(r.p2pInFlightRequested) == 0 {
		// the first bach in the response is the one we requested, so clear the in-flight request
		r.p2pInFlightRequested = nil
		r.p2pInFlightReqTime = nil
	}
	r.p2pReqMutex.Unlock()

	// try to add all the batches to the db, and notify subscribers if they are new and live
	for _, batch := range batches {
		err := r.AddBatch(batch)
		if err != nil {
			if !errors.Is(err, errutil.ErrAlreadyExists) {
				r.logger.Warn("unable to add p2p batch to L2 batch repository", log.ErrKey, err)
			}
			// we've already seen this batch or failed to store it for another reason - do not notify subscribers
			return
		}
		if isLive {
			// notify subscribers if the batch is new
			for _, subscriber := range r.subscribers {
				go subscriber.HandleBatch(batch)
			}
		}
	}
}

// HandleBatchRequest handles a request for a batch from a peer, sending batches to the requester asynchronously
// todo (#1625) - only allow requests for batches since last rollup, to avoid DoS attacks.
func (r *Repository) HandleBatchRequest(requesterID string, fromSeqNo *big.Int) {
	batches := make([]*common.ExtBatch, 0)
	nextSeqNum := fromSeqNo
	for len(batches) <= _maxBatchesInP2PResponse {
		batch, err := r.db.GetBatchBySequenceNumber(nextSeqNum)
		if err != nil {
			if !errors.Is(err, errutil.ErrNotFound) {
				r.logger.Warn("unexpected error fetching batches for peer req", log.BatchSeqNoKey, nextSeqNum, log.ErrKey, err)
			}
			break // once one batch lookup fails we don't expect to find any of them
		}
		batches = append(batches, batch)
		nextSeqNum = nextSeqNum.Add(nextSeqNum, big.NewInt(1))
	}
	if len(batches) == 0 {
		return // nothing to send
	}

	err := r.sl.P2P().RespondToBatchRequest(requesterID, batches)
	if err != nil {
		r.logger.Warn("unable to send batches to peer", "peer", requesterID, log.ErrKey, err)
	}
}

// Subscribe registers a handler to be notified of new head batches as they arrive
func (r *Repository) Subscribe(subscriber host.L2BatchHandler) {
	r.subscribers = append(r.subscribers, subscriber)
}

func (r *Repository) FetchBatchBySeqNo(seqNo *big.Int) (*common.ExtBatch, error) {
	b, err := r.db.GetBatchBySequenceNumber(seqNo)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) && seqNo.Cmp(r.latestBatchSeqNo) < 0 {
			if r.isSequencer {
				// sequencer does not request batches from peers, it checks if its enclave has the batch
				return r.fetchBatchFallbackToEnclave(seqNo)
			}
			// we haven't seen this batch before, but it is older than the latest batch we have seen so far
			// Request missing batches from peers (the batches from any response will be added asynchronously, so
			// we will return the not found error and hopefully future attempts will succeed)
			go r.requestMissingBatchesFromPeers(seqNo)
		}
		return nil, err
	}
	return b, nil
}

// AddBatch allows the host to add a batch to the repository, this is used:
// - when the node is a sequencer to store newly produced batches (the only way the sequencer host receives batches)
// - when the node is a validator to store batches read from roll-ups
// If the repository already has the batch it returns an AlreadyExists error which is typically ignored.
func (r *Repository) AddBatch(batch *common.ExtBatch) error {
	r.logger.Debug("Saving batch", log.BatchSeqNoKey, batch.Header.SequencerOrderNo, log.BatchHashKey, batch.Hash())
	err := r.db.AddBatch(batch)
	if err != nil {
		return err
	}
	// atomically compare and swap latest batch sequence number if successfully added batch is newer
	r.latestSeqNoMutex.Lock()
	defer r.latestSeqNoMutex.Unlock()
	if batch.Header.SequencerOrderNo.Cmp(r.latestBatchSeqNo) > 0 {
		r.latestBatchSeqNo = batch.Header.SequencerOrderNo
	}
	return nil
}

func (r *Repository) fetchBatchFallbackToEnclave(seqNo *big.Int) (*common.ExtBatch, error) {
	b, err := r.sl.Enclaves().LookupBatchBySeqNo(seqNo)
	if err != nil {
		return nil, err
	}

	// asynchronously add that batch to the repo, so we have it for the next request
	go func() {
		err := r.AddBatch(b)
		if err != nil {
			r.logger.Info("unable to add batch that was returned from the enclave", log.ErrKey, err)
		}
	}()

	return b, nil
}

// RequestMissingBatches requests batches from peers from the specified sequence number.
// It is an asynchronous request and the repository does not expect to be notified of the result.
func (r *Repository) requestMissingBatchesFromPeers(fromSeqNo *big.Int) {
	r.p2pReqMutex.Lock()
	defer r.p2pReqMutex.Unlock()
	if r.p2pInFlightReqTime != nil && time.Since(*r.p2pInFlightReqTime) < _timeoutWaitingForP2PResponse {
		// don't send request if we have sent one too recently
		r.logger.Trace("not requesting missing batches from sequencer - too soon since last request", "fromSeqNo", fromSeqNo, "lastReq", r.p2pInFlightReqTime)
		return
	}

	r.logger.Debug("requesting missing batches from sequencer", "fromSeqNo", fromSeqNo)
	err := r.sl.P2P().RequestBatchesFromSequencer(fromSeqNo)
	if err != nil {
		r.logger.Warn("unable to request missing batches from sequencer", "fromSeqNo", fromSeqNo, log.ErrKey, err)
		return
	}
	now := time.Now()
	r.p2pInFlightReqTime = &now
}
