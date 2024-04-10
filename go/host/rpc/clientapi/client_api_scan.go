package clientapi

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// ScanAPI implements metric specific RPC endpoints
type ScanAPI struct {
	host   host.Host
	logger log.Logger
}

func NewScanAPI(host host.Host, logger log.Logger) *ScanAPI {
	return &ScanAPI{
		host:   host,
		logger: logger,
	}
}

// GetTotalContractCount returns the number of recorded contracts on the network.
func (s *ScanAPI) GetTotalContractCount() (*big.Int, error) {
	return s.host.EnclaveClient().GetTotalContractCount()
}

// GetTotalTxCount returns the number of recorded transactions on the network.
func (s *ScanAPI) GetTotalTxCount() (*big.Int, error) {
	return s.host.Storage().FetchTotalTxCount()
}

// GetBatchListing returns a paginated list of batches
func (s *ScanAPI) GetBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return s.host.Storage().FetchBatchListing(pagination)
}

// GetBatchListingDeprecated returns the deprecated version of batch listing
func (s *ScanAPI) GetBatchListingDeprecated(pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	return s.host.Storage().FetchBatchListingDeprecated(pagination)
}

// GetPublicBatchByHash returns the public batch
func (s *ScanAPI) GetPublicBatchByHash(hash common.L2BatchHash) (*common.PublicBatch, error) {
	return s.host.Storage().FetchPublicBatchByHash(hash)
}

// GetBatch returns the `ExtBatch` with the given hash
func (s *ScanAPI) GetBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	return s.host.Storage().FetchBatch(batchHash)
}

// GetBatchByTx returns the `ExtBatch` with the given tx hash
func (s *ScanAPI) GetBatchByTx(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	return s.host.Storage().FetchBatchByTx(txHash)
}

// GetLatestBatch returns the head `BatchHeader`
func (s *ScanAPI) GetLatestBatch() (*common.BatchHeader, error) {
	return s.host.Storage().FetchLatestBatch()
}

// GetBatchByHeight returns the `BatchHeader` with the given height
func (s *ScanAPI) GetBatchByHeight(height *big.Int) (*common.BatchHeader, error) {
	return s.host.Storage().FetchBatchHeaderByHeight(height)
}

// GetRollupListing returns a paginated list of Rollups
func (s *ScanAPI) GetRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	return s.host.Storage().FetchRollupListing(pagination)
}

// GetLatestRollupHeader returns the head `RollupHeader`
func (s *ScanAPI) GetLatestRollupHeader() (*common.RollupHeader, error) {
	return s.host.Storage().FetchLatestRollupHeader()
}

// GetPublicTransactionData returns a paginated list of transaction data
func (s *ScanAPI) GetPublicTransactionData(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return s.host.EnclaveClient().GetPublicTransactionData(pagination)
}

// GetBlockListing returns a paginated list of blocks that include rollups
func (s *ScanAPI) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	return s.host.Storage().FetchBlockListing(pagination)
}
