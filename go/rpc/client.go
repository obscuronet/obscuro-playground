package rpc

import (
	"context"
	"errors"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

const (
	BatchNumber           = "eth_blockNumber"
	Call                  = "eth_call"
	ChainID               = "eth_chainId"
	GetBalance            = "eth_getBalance"
	GetBatchByHash        = "eth_getBlockByHash"
	GetBatchByNumber      = "eth_getBlockByNumber"
	GetCode               = "eth_getCode"
	GetTransactionByHash  = "eth_getTransactionByHash"
	GetTransactionCount   = "eth_getTransactionCount"
	GetTransactionReceipt = "eth_getTransactionReceipt"
	SendRawTransaction    = "eth_sendRawTransaction"
	EstimateGas           = "eth_estimateGas"
	GetLogs               = "eth_getLogs"
	GetStorageAt          = "eth_getStorageAt"
	GasPrice              = "eth_gasPrice"

	Health = "obscuro_health"
	Config = "obscuro_config"

	StopHost             = "test_stopHost"
	Subscribe            = "eth_subscribe"
	Unsubscribe          = "eth_unsubscribe"
	SubscribeNamespace   = "eth"
	SubscriptionTypeLogs = "logs"

	GetBatchForTx             = "scan_getBatchForTx"
	GetLatestRollupHeader     = "scan_getLatestRollupHeader"
	GetTotalTransactionCount  = "scan_getTotalTransactionCount"
	GetTotalContractCount     = "scan_getTotalContractCount"
	GetPublicTransactionData  = "scan_getPublicTransactionData"
	GetBatchListing           = "scan_getBatchListing"
	GetBatchListingDeprecated = "scan_getBatchListingDeprecated"
	GetBlockListing           = "scan_getBlockListing"
	GetRollupListing          = "scan_getRollupListing"
	GetFullBatchByHash        = "scan_getFullBatchByHash"
	GetLatestBatch            = "scan_getLatestBatch"
	GetPublicBatchByHash      = "scan_getPublicBatchByHash"
	GetBatchByHeight          = "scan_getBatchByHeight"
)

var ErrNilResponse = errors.New("nil response received from Obscuro node")

// Client is used by client applications to interact with the Obscuro node
type Client interface {
	// Call executes the named method via RPC. (Returns `ErrNilResponse` on nil response from Node, this is used as "not found" for some method calls)
	Call(result interface{}, method string, args ...interface{}) error
	// CallContext If the context is canceled before the call has successfully returned, CallContext returns immediately.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	// Subscribe creates a subscription to the Obscuro host.
	Subscribe(ctx context.Context, result interface{}, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
	// Stop closes the client.
	Stop()
}
