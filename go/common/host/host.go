package host

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// Host is the half of the Obscuro node that lives outside the enclave.
type Host interface {
	Config() *config.HostConfig
	DB() *sql.DB
	EnclaveClient() common.Enclave

	// Start initializes the main loop of the host.
	Start() error
	// SubmitAndBroadcastTx submits an encrypted transaction to the enclave, and broadcasts it to the other hosts on the network.
	SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error)
	// Subscribe feeds logs matching the encrypted log subscription to the matchedLogs channel.
	Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogs chan []byte) error
	// Unsubscribe terminates a log subscription between the host and the enclave.
	Unsubscribe(id rpc.ID)
	// Stop gracefully stops the host execution.
	Stop() error

	// HealthCheck returns the health status of the host + enclave + db
	HealthCheck() (*HealthCheck, error)

	// ObscuroConfig returns the info of the Obscuro network
	ObscuroConfig() (*common.ObscuroNetworkInfo, error)
}

type BlockStream struct {
	Stream <-chan *types.Block // the channel which will receive the consecutive, canonical blocks
	Stop   func()              // function to permanently stop the stream and clean up any associated processes/resources
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
	L1BlockRepository
}
