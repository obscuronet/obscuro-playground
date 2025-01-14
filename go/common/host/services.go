package host

import (
	"context"
	"math/big"

	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
)

// service names - these are the keys used to register known services with the host
const (
	P2PName                    = "p2p"
	L1DataServiceName          = "l1-data-service"
	L1PublisherName            = "l1-publisher"
	L2BatchRepositoryName      = "l2-batch-repo"
	EnclaveServiceName         = "enclaves"
	LogSubscriptionServiceName = "log-subs"
	FilterAPIServiceName       = "filter-api"
	CrossChainServiceName      = "cross-chain"
)

// The host has a number of services that encapsulate the various responsibilities of the host.
// This file contains service-level abstractions and utilities, as well as all the interfaces for these services, code
// should depend on these interfaces rather than the concrete implementations.

// Service interface allows the host to manage all services in a generic way
// Note: Services may depend on other services but they shouldn't use them during construction, only when 'Start()' is called.
// They should be resilient to services availability, because the construction ordering is not guaranteed.
type Service interface {
	Start() error
	Stop() error
	HealthStatus(context.Context) HealthStatus
}

// P2P provides an interface for the host to interact with the P2P network
type P2P interface {
	// BroadcastBatches sends live batch(es) to every other node on the network
	BroadcastBatches(batches []*common.ExtBatch) error
	// SendTxToSequencer sends the encrypted transaction to the sequencer.
	SendTxToSequencer(tx common.EncryptedTx) error

	// RequestBatchesFromSequencer asynchronously requests batches from the sequencer, from the given sequence number
	RequestBatchesFromSequencer(fromSeqNo *big.Int) error
	// RespondToBatchRequest sends the requested batches to the requesting peer
	RespondToBatchRequest(requestID string, batches []*common.ExtBatch) error

	// SubscribeForBatches will register a handler to receive new batches from peers, returns unsubscribe func
	SubscribeForBatches(handler P2PBatchHandler) func()
	// SubscribeForTx will register a handler to receive new transactions from peers, returns unsubscribe func
	SubscribeForTx(handler P2PTxHandler) func()
	// SubscribeForBatchRequests will register a handler to receive new batch requests from peers, returns unsubscribe func
	// todo (@matt) feels a bit weird to have this in this interface since it relates to serving data rather than receiving
	SubscribeForBatchRequests(handler P2PBatchRequestHandler) func()
}

// P2PBatchHandler is an interface for receiving new batches from the P2P network as they arrive
type P2PBatchHandler interface {
	// HandleBatches will be called in a new goroutine for batches that arrive
	HandleBatches(batch []*common.ExtBatch, isLive bool)
}

// P2PTxHandler is an interface for receiving new transactions from the P2P network as they arrive
type P2PTxHandler interface {
	// HandleTransaction will be called in a new goroutine for each new tx as it arrives
	HandleTransaction(tx common.EncryptedTx)
}

type P2PBatchRequestHandler interface {
	// HandleBatchRequest will be called in a new goroutine for each new batch request as it arrives
	HandleBatchRequest(requestID string, fromSeqNo *big.Int)
}

// L1DataService provides an interface for the host to request L1 block data (live-streaming and historical)
type L1DataService interface {
	// Subscribe will register a block handler to receive new blocks as they arrive, returns unsubscribe func
	Subscribe(handler L1BlockHandler) func()

	FetchBlockByHeight(height *big.Int) (*types.Header, error)
	// FetchNextBlock returns the next canonical block after a given block hash
	// It returns the new block, a bool which is true if the block is the current L1 head and a bool if the block is on a different fork to prevBlock
	FetchNextBlock(prevBlock gethcommon.Hash) (*types.Header, bool, error)
	// GetTenRelevantTransactions returns the events and transactions relevant to Ten
	GetTenRelevantTransactions(block *types.Header) (*common.ProcessedL1Data, error)
}

// L1BlockHandler is an interface for receiving new blocks from the repository as they arrive
type L1BlockHandler interface {
	// HandleBlock will be called in a new goroutine for each new block as it arrives
	HandleBlock(block *types.Header)
}

// L1Publisher provides an interface for the host to interact with Ten data (management contract etc.) on L1
type L1Publisher interface {
	// InitializeSecret will send a management contract transaction to initialize the network with the generated secret
	InitializeSecret(attestation *common.AttestationReport, encSecret common.EncryptedSharedEnclaveSecret) error
	// RequestSecret will send a management contract transaction to request a secret from the enclave, returning the L1 head at time of sending
	RequestSecret(report *common.AttestationReport) (gethcommon.Hash, error)
	// FindSecretResponseTx will return the secret response tx from an L1 block
	FindSecretResponseTx(responseTxs []*common.L1TxData) []*common.L1RespondSecretTx
	// PublishRollup will create and publish a rollup tx to the management contract - fire and forget we don't wait for receipt
	// todo (#1624) - With a single sequencer, it is problematic if rollup publication fails; handle this case better
	PublishRollup(producedRollup *common.ExtRollup)
	// PublishSecretResponse will create and publish a secret response tx to the management contract - fire and forget we don't wait for receipt
	PublishSecretResponse(secretResponse *common.ProducedSecretResponse) error

	// PublishCrossChainBundle will create and publish a cross-chain bundle tx to the management contract
	PublishCrossChainBundle(*common.ExtCrossChainBundle, *big.Int, gethcommon.Hash) error

	FetchLatestSeqNo() (*big.Int, error)

	// GetImportantContracts returns a (cached) record of addresses of the important network contracts
	GetImportantContracts() map[string]gethcommon.Address
	// ResyncImportantContracts will fetch the latest important contracts from the management contract, update the cache
	ResyncImportantContracts() error

	// GetBundleRangeFromManagementContract returns the range of batches for which to build a bundle
	GetBundleRangeFromManagementContract(lastRollupNumber *big.Int, lastRollupUID gethcommon.Hash) (*gethcommon.Hash, *big.Int, *big.Int, error)
}

// L2BatchRepository provides an interface for the host to request L2 batch data (live-streaming and historical)
type L2BatchRepository interface {
	// SubscribeNewBatches will register a handler to receive new batches from the publisher as they arrive at the host
	SubscribeNewBatches(handler L2BatchHandler) func()

	// SubscribeValidatedBatches will register a handler to receive batches that have been validated by the enclave
	SubscribeValidatedBatches(handler L2BatchHandler) func()

	FetchBatchBySeqNo(background context.Context, seqNo *big.Int) (*common.ExtBatch, error)

	// AddBatch is used to notify the repository of a new batch, e.g. from the enclave when seq produces one or a rollup is consumed
	// Note: it is fine to add batches that the repo already has, it will just ignore them
	AddBatch(batch *common.ExtBatch) error

	// NotifyNewValidatedHead - called after an enclave validates a batch, to update the repo's validated head and notify subscribers
	NotifyNewValidatedHead(batch *common.ExtBatch)
}

// L2BatchHandler is an interface for receiving new batches from the publisher as they arrive
type L2BatchHandler interface {
	// HandleBatch will be called in a new goroutine for each new batch as it arrives
	HandleBatch(batch *common.ExtBatch)
}

// EnclaveService provides access to the host enclave(s), it monitors and manages the states of the enclaves, so it can
// help with failover, gate-keeping (throttling/load-balancing) and circuit-breaking when the enclave is not available
type EnclaveService interface {
	// LookupBatchBySeqNo is used to fetch batch data from the enclave - it is only used as a fallback for the sequencer
	// host if it's missing a batch (other host services should use L2Repo to fetch batch data)
	LookupBatchBySeqNo(ctx context.Context, seqNo *big.Int) (*common.ExtBatch, error)

	// GetEnclaveClient returns an enclave client // todo (@matt) we probably don't want to expose this
	GetEnclaveClient() common.Enclave

	// GetEnclaveClients returns a list of all enclave clients
	GetEnclaveClients() []common.Enclave

	// EvictEnclave will remove the enclave from the list of enclaves, it is used when an enclave is unhealthy
	// - the enclave guardians are responsible for calling this method when they detect an enclave is unhealthy to notify
	//	 the service that it should failover if possible
	NotifyUnavailable(enclaveID *common.EnclaveID)

	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to other hosts on the network (in particular, to the sequencer)
	SubmitAndBroadcastTx(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.RawTx, error)

	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription) error
	Unsubscribe(id rpc.ID) error
}

// LogSubscriptionManager provides an interface for the host to manage log subscriptions
type LogSubscriptionManager interface {
	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error
	Unsubscribe(id rpc.ID)
	SendLogsToSubscribers(result *common.EncryptedSubscriptionLogs)
}
