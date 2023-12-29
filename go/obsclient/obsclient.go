package obsclient

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

// ObsClient provides access to general Obscuro functionality that doesn't require viewing keys.
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type ObsClient struct {
	rpcClient rpc.Client
}

func Dial(rawurl string) (*ObsClient, error) {
	rc, err := rpc.NewNetworkClient(rawurl)
	if err != nil {
		return nil, err
	}
	return NewObsClient(rc), nil
}

func NewObsClient(c rpc.Client) *ObsClient {
	return &ObsClient{c}
}

func (oc *ObsClient) Close() {
	oc.rpcClient.Stop()
}

// Blockchain Access

// ChainID retrieves the current chain ID for transaction replay protection.
func (oc *ObsClient) ChainID() (*big.Int, error) {
	var result hexutil.Big
	err := oc.rpcClient.Call(&result, rpc.ChainID)
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// BatchNumber returns the height of the head rollup
func (oc *ObsClient) BatchNumber() (uint64, error) {
	var result hexutil.Uint64
	err := oc.rpcClient.Call(&result, rpc.BatchNumber)
	return uint64(result), err
}

// BatchByHash returns the batch with the given hash.
func (oc *ObsClient) BatchByHash(hash gethcommon.Hash) (*common.ExtBatch, error) {
	var batch *common.ExtBatch
	err := oc.rpcClient.Call(&batch, rpc.GetFullBatchByHash, hash)
	if err == nil && batch == nil {
		err = ethereum.NotFound
	}
	return batch, err
}

// BatchHeaderByNumber returns the header of the rollup with the given number
func (oc *ObsClient) BatchHeaderByNumber(number *big.Int) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetBatchByNumber, toBlockNumArg(number), false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// BatchHeaderByHash returns the block header with the given hash.
func (oc *ObsClient) BatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetBatchByHash, hash, false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// Health returns the health of the node.
func (oc *ObsClient) Health() (bool, error) {
	var healthy *hostcommon.HealthCheck
	err := oc.rpcClient.Call(&healthy, rpc.Health)
	if err != nil {
		return false, err
	}
	if !healthy.OverallHealth {
		return false, errors.New(strings.Join(healthy.Errors, ", "))
	}
	return healthy.OverallHealth, nil
}

// GetTotalContractCount returns the total count of created contracts
func (oc *ObsClient) GetTotalContractCount() (int, error) {
	var count int
	err := oc.rpcClient.Call(&count, rpc.GetTotalContractCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalTransactionCount returns the total count of executed transactions
func (oc *ObsClient) GetTotalTransactionCount() (int, error) {
	var count int
	err := oc.rpcClient.Call(&count, rpc.GetTotalTransactionCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetLatestRollupHeader returns the header of the rollup at tip
func (oc *ObsClient) GetLatestRollupHeader() (*common.RollupHeader, error) {
	var header *common.RollupHeader
	err := oc.rpcClient.Call(&header, rpc.GetLatestRollupHeader)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// GetPublicTxListing returns a list of public transactions
func (oc *ObsClient) GetPublicTxListing(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	var result common.TransactionListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetPublicTransactionData, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBatchesListing returns a list of batches
func (oc *ObsClient) GetBatchesListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	var result common.BatchListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetBatchListing, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBlockListing returns a list of block headers
func (oc *ObsClient) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	var result common.BlockListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetBlockListing, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetConfig returns the network config for obscuro
func (oc *ObsClient) GetConfig() (*common.ObscuroNetworkInfo, error) {
	var result common.ObscuroNetworkInfo
	err := oc.rpcClient.Call(&result, rpc.Config)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
