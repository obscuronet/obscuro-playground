package clientapi

import (
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// ObscuroScanAPI implements ObscuroScan-specific JSON RPC operations.
type ObscuroScanAPI struct {
	host host.Host
}

func NewObscuroScanAPI(host host.Host) *ObscuroScanAPI {
	return &ObscuroScanAPI{
		host: host,
	}
}

// GetBlockHeaderByHash returns the header for the block with the given number.
func (api *ObscuroScanAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader, found := api.host.DB().GetBlockHeader(blockHash)
	if !found {
		return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
	}
	return blockHeader, nil
}

// GetHeadRollupHeader returns the current head rollup's header.
// TODO - #718 - Switch to reading batch header.
func (api *ObscuroScanAPI) GetHeadRollupHeader() *common.Header {
	headerWithHashes, found := api.host.DB().GetHeadRollupHeader()
	if !found {
		return nil
	}
	return headerWithHashes.Header
}

// GetRollup returns the rollup with the given hash.
func (api *ObscuroScanAPI) GetRollup(hash gethcommon.Hash) (*common.ExtRollup, error) {
	return api.host.EnclaveClient().GetRollup(hash)
}

// GetRollupForTx returns the rollup containing a given transaction hash. Required for ObscuroScan.
func (api *ObscuroScanAPI) GetRollupForTx(txHash gethcommon.Hash) (*common.ExtRollup, error) {
	rollupNumber, found := api.host.DB().GetRollupNumber(txHash)
	if !found {
		return nil, fmt.Errorf("no rollup containing a transaction with hash %s is stored", txHash)
	}

	rollupHash, found := api.host.DB().GetRollupHash(rollupNumber)
	if !found {
		return nil, fmt.Errorf("no rollup with number %d is stored", rollupNumber.Int64())
	}

	rollup, err := api.host.EnclaveClient().GetRollup(*rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with hash %s. Cause: %w", rollupNumber, err)
	}

	return rollup, nil
}

// GetLatestTransactions returns the hashes of the latest `num` transactions, or as many as possible if less than `num` transactions exist.
// TODO - Consider introducing paging or similar to prevent a huge number of transactions blowing the node's memory limit.
// TODO - #718 - Switch to retrieving transactions from latest batch.
func (api *ObscuroScanAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	currentRollupHeaderWithHashes, found := api.host.DB().GetHeadRollupHeader()
	if !found {
		return nil, nil
	}
	nextRollupHash := currentRollupHeaderWithHashes.Header.Hash()

	// We walk the chain until we've collected sufficient transactions.
	var txHashes []gethcommon.Hash
	for {
		rollupHeaderWithHashes, found := api.host.DB().GetRollupHeader(nextRollupHash)
		if !found {
			return nil, fmt.Errorf("could not retrieve rollup for hash %s", nextRollupHash)
		}

		for _, txHash := range rollupHeaderWithHashes.TxHashes {
			txHashes = append(txHashes, txHash)
			if len(txHashes) >= num {
				return txHashes, nil
			}
		}

		// If we have reached the top of the chain (i.e. the current rollup's number is one), we stop walking.
		if rollupHeaderWithHashes.Header.Number.Cmp(big.NewInt(0)) == 0 {
			break
		}
		nextRollupHash = rollupHeaderWithHashes.Header.ParentHash
	}

	return txHashes, nil
}

// GetTotalTransactions returns the number of recorded transactions on the network.
func (api *ObscuroScanAPI) GetTotalTransactions() *big.Int {
	totalTransactions := api.host.DB().GetTotalTransactions()
	return totalTransactions
}

// Attestation returns the node's attestation details.
func (api *ObscuroScanAPI) Attestation() (*common.AttestationReport, error) {
	return api.host.EnclaveClient().Attestation()
}
