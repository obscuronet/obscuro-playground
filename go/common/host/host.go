package host

import (
	"context"

	"github.com/ten-protocol/go-ten/go/common"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/host/storage"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	Config() *hostconfig.HostConfig
	EnclaveClient() common.Enclave
	Storage() storage.Storage
	// Start initializes the main loop of the host.
	Start() error
	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to the other hosts on the network.
	SubmitAndBroadcastTx(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.RawTx, error)
	// SubscribeLogs feeds logs matching the encrypted log subscription to the matchedLogs channel.
	SubscribeLogs(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogs chan []byte) error
	// UnsubscribeLogs terminates a log subscription between the host and the enclave.
	UnsubscribeLogs(id rpc.ID)
	// Stop gracefully stops the host execution.
	Stop() error

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck(context.Context) (*HealthCheck, error)

	// TenConfig returns the info of the Obscuro network
	TenConfig() (*common.TenNetworkInfo, error)

	// NewHeadsChan returns live batch headers
	// Note - do not use directly. This is meant only for the NewHeadsManager, which multiplexes the headers
	NewHeadsChan() chan *common.BatchHeader
}

type BatchMsg struct {
	Batches []*common.ExtBatch // The batches being sent.
	IsLive  bool               // true if these batches are being sent as new, false if in response to a p2p request
}

type P2PHostService interface {
	Service
	P2P
}

type L1RepoService interface {
	Service
	L1DataService
}
