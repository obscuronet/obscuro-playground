package obsclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
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

// RollupNumber returns the height of the head rollup
func (oc *ObsClient) RollupNumber() (uint64, error) {
	var result hexutil.Uint64
	err := oc.rpcClient.Call(&result, rpc.RollupNumber)
	return uint64(result), err
}

// RollupHeaderByNumber returns the header of the rollup with the given number
func (oc *ObsClient) RollupHeaderByNumber(number *big.Int) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetRollupByNumber, toBlockNumArg(number), false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// RollupHeaderByHash returns the block header with the given hash.
func (oc *ObsClient) RollupHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetRollupByHash, hash, false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// Health returns the health of the node.
func (oc *ObsClient) Health() (bool, error) {
	var healthy *hostcommon.HealthCheck
	err := oc.rpcClient.Call(&healthy, rpc.Health)
	return healthy.OverallHealth, err
}
