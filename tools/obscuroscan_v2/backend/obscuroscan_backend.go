package backend

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/obsclient"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type Backend struct {
	obsClient *obsclient.ObsClient
}

func NewBackend(obsClient *obsclient.ObsClient) *Backend {
	return &Backend{
		obsClient: obsClient,
	}
}

func (b *Backend) GetLatestBatch() (*common.BatchHeader, error) {
	return b.obsClient.BatchHeaderByNumber(nil)
}

func (b *Backend) GetLatestRollup() (*common.RollupHeader, error) {
	// return b.obsClient.L1RollupHeaderByNumber(nil)
	return &common.RollupHeader{}, nil
}

func (b *Backend) GetNodeCount() (int, error) {
	// return b.obsClient.ActiveNodeCount()
	return 0, nil
}

func (b *Backend) GetTotalContractCount() (int, error) {
	return b.obsClient.GetTotalContractCount()
}

func (b *Backend) GetTotalTransactionCount() (int, error) {
	return b.obsClient.GetTotalTransactionCount()
}

func (b *Backend) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return b.obsClient.GetLatestRollupHeader()
}

func (b *Backend) GetBatch(hash gethcommon.Hash) (*common.BatchHeader, error) {
	return b.obsClient.BatchHeaderByHash(hash)
}
