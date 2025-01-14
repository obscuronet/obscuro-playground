package subscription

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ethereum/go-ethereum/core/types"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// NewHeadsService multiplexes new batch header messages from an input channel into multiple subscribers
// also handles unsubscribe
// Note: this is a service which must be Started and Stopped
type NewHeadsService struct {
	connectFunc        func() (chan *common.BatchHeader, <-chan error, error)
	convertToEthHeader bool
	notifiersMutex     *sync.RWMutex
	newHeadNotifiers   map[rpc.ID]*rpc.Notifier
	onMessage          func(*common.BatchHeader) error
	stopped            *atomic.Bool
	logger             gethlog.Logger
}

// connect - function that returns the input channel
func NewNewHeadsService(connect func() (chan *common.BatchHeader, <-chan error, error), convertToEthHeader bool, logger gethlog.Logger, onMessage func(*common.BatchHeader) error) *NewHeadsService {
	return &NewHeadsService{
		connectFunc:        connect,
		convertToEthHeader: convertToEthHeader,
		onMessage:          onMessage,
		logger:             logger,
		stopped:            &atomic.Bool{},
		newHeadNotifiers:   make(map[rpc.ID]*rpc.Notifier),
		notifiersMutex:     &sync.RWMutex{},
	}
}

func (nhs *NewHeadsService) Start() error {
	nhs.reconnect()
	return nil
}

func (nhs *NewHeadsService) reconnect() {
	// reconnect to the backend and restart the listening
	newCh, errCh, err := nhs.connectFunc()
	if err != nil {
		nhs.logger.Crit("could not connect to new heads: ", log.ErrKey, err)
	}
	nhs._subscribe(newCh, errCh)
}

func (nhs *NewHeadsService) _subscribe(inputCh chan *common.BatchHeader, errChan <-chan error) {
	backedUnsub := &atomic.Bool{}
	go HandleUnsubscribeErrChan([]<-chan error{errChan}, func() {
		backedUnsub.Store(true)
	})
	go ForwardFromChannels(
		[]chan *common.BatchHeader{inputCh},
		func(head *common.BatchHeader) error {
			return nhs.onNewBatch(head)
		},
		func() {
			nhs.logger.Info("Disconnected from new head subscription. Reconnecting...")
			nhs.reconnect()
		},
		backedUnsub,
		nhs.stopped,
		2*time.Minute, // todo - create constant
		nhs.logger,
	)
}

func (nhs *NewHeadsService) onNewBatch(head *common.BatchHeader) error {
	if nhs.onMessage != nil {
		err := nhs.onMessage(head)
		if err != nil {
			nhs.logger.Info("failed invoking onMessage callback.", log.ErrKey, err)
		}
	}

	var msg any = head
	if nhs.convertToEthHeader {
		msg = ConvertBatchHeader(head)
	}

	nhs.notifiersMutex.Lock()
	defer nhs.notifiersMutex.Unlock()

	// for each new head, notify all registered subscriptions
	for id, notifier := range nhs.newHeadNotifiers {
		if nhs.stopped.Load() {
			return nil
		}
		err := notifier.Notify(id, msg)
		if err != nil {
			// on error, remove the notification
			nhs.logger.Info("failed to notify newHead subscription", log.ErrKey, err, log.SubIDKey, id)
			delete(nhs.newHeadNotifiers, id)
		}
	}
	return nil
}

func (nhs *NewHeadsService) RegisterNotifier(notifier *rpc.Notifier, subscription *rpc.Subscription) {
	nhs.notifiersMutex.Lock()
	defer nhs.notifiersMutex.Unlock()
	nhs.newHeadNotifiers[subscription.ID] = notifier

	go HandleUnsubscribe(subscription, func() {
		nhs.notifiersMutex.Lock()
		defer nhs.notifiersMutex.Unlock()
		delete(nhs.newHeadNotifiers, subscription.ID)
	})
}

func (nhs *NewHeadsService) Stop() error {
	nhs.stopped.Store(true)
	return nil
}

func (nhs *NewHeadsService) HealthStatus(context.Context) host.HealthStatus {
	return &host.BasicErrHealthStatus{}
}

func ConvertBatchHeader(head *common.BatchHeader) *types.Header {
	return &types.Header{
		ParentHash:  head.ParentHash,
		UncleHash:   gethutil.EmptyHash,
		Coinbase:    head.Coinbase,
		Root:        head.Root,
		TxHash:      head.TxHash,
		ReceiptHash: head.ReceiptHash,
		Bloom:       types.Bloom{},
		Difficulty:  big.NewInt(0),
		Number:      head.Number,
		GasLimit:    head.GasLimit,
		GasUsed:     head.GasUsed,
		Time:        head.Time,
		Extra:       make([]byte, 0),
		BaseFee:     head.BaseFee,
	}
}
