package gethencoding

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

const (
	// The relevant fields in an eth_call request's params.
	callFieldTo                   = "to"
	CallFieldFrom                 = "from"
	callFieldData                 = "data" // this field was renamed in geth CallMsg to 'input' but we need to support both
	callFieldInput                = "input"
	callFieldValue                = "value"
	callFieldGas                  = "gas"
	callFieldNonce                = "nonce"
	callFieldGasPrice             = "gasprice"
	callFieldMaxFeePerGas         = "maxfeepergas"
	callFieldMaxPriorityFeePerGas = "maxpriorityfeepergas"
	callFieldAccessList           = "accesslist"
)

// EncodingService handles conversion to Geth data structures
type EncodingService interface {
	CreateEthHeaderForBatch(ctx context.Context, h *common.BatchHeader) (*types.Header, error)
	CreateEthBlockFromBatch(ctx context.Context, b *core.Batch) (*types.Block, error)
}

type gethEncodingServiceImpl struct {
	storage        storage.Storage
	logger         gethlog.Logger
	cachingService *storage.CacheService
	entropyService *crypto.EvmEntropyService
}

func NewGethEncodingService(storage storage.Storage, cachingService *storage.CacheService, entropyService *crypto.EvmEntropyService, logger gethlog.Logger) EncodingService {
	return &gethEncodingServiceImpl{
		storage:        storage,
		logger:         logger,
		entropyService: entropyService,
		cachingService: cachingService,
	}
}

// ExtractAddress returns a gethcommon.Address given an interface{}, errors if unexpected values are used
func ExtractAddress(param interface{}) (*gethcommon.Address, error) {
	if param == nil {
		return nil, fmt.Errorf("no address specified")
	}

	paramStr, ok := param.(string)
	if !ok {
		return nil, fmt.Errorf("unexpectd address value")
	}

	if len(strings.TrimSpace(paramStr)) == 0 {
		return nil, fmt.Errorf("no address specified")
	}

	addr := gethcommon.HexToAddress(paramStr)
	return &addr, nil
}

// ExtractOptionalBlockNumber defaults nil or empty block number params to latest block number
func ExtractOptionalBlockNumber(params []interface{}, idx int) (*gethrpc.BlockNumberOrHash, error) {
	latest := gethrpc.BlockNumberOrHashWithNumber(gethrpc.LatestBlockNumber)
	if len(params) <= idx {
		return &latest, nil
	}
	if params[idx] == nil {
		return &latest, nil
	}
	if emptyStr, ok := params[idx].(string); ok && len(strings.TrimSpace(emptyStr)) == 0 {
		return &latest, nil
	}

	return ExtractBlockNumber(params[idx])
}

func ExtractBlockNumber(param interface{}) (*gethrpc.BlockNumberOrHash, error) {
	if param == nil {
		latest := gethrpc.BlockNumberOrHashWithNumber(gethrpc.LatestBlockNumber)
		return &latest, nil
	}

	// when the param is a single string we try to convert it to a block number
	blockString, ok := param.(string)
	if ok {
		blockNumber := gethrpc.BlockNumber(0)
		err := blockNumber.UnmarshalJSON([]byte(blockString))
		if err != nil {
			return nil, fmt.Errorf("invalid block number %s - %w", blockString, err)
		}
		return &gethrpc.BlockNumberOrHash{BlockNumber: &blockNumber}, nil
	}

	var blockNo *gethrpc.BlockNumber
	var blockHa *gethcommon.Hash
	var reqCanon bool

	blockAndHash, ok := param.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid block or hash parameter")
	}
	if blockAndHash["blockNumber"] != nil {
		b, ok := blockAndHash["blockNumber"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid blockNumber parameter")
		}
		blockNumber := gethrpc.BlockNumber(0)
		err := blockNumber.UnmarshalJSON([]byte(b))
		if err != nil {
			return nil, fmt.Errorf("invalid block number %s - %w", b, err)
		}
		blockNo = &blockNumber
	}
	if blockAndHash["blockHash"] != nil {
		bh, ok := blockAndHash["blockHash"].(gethcommon.Hash)
		if !ok {
			return nil, fmt.Errorf("invalid blockhash parameter")
		}
		blockHa = &bh
	}
	if blockAndHash["RequireCanonical"] != nil {
		reqCanon, ok = blockAndHash["RequireCanonical"].(bool)
		if !ok {
			return nil, fmt.Errorf("invalid RequireCanonical parameter")
		}
	}

	return &gethrpc.BlockNumberOrHash{
		BlockNumber:      blockNo,
		BlockHash:        blockHa,
		RequireCanonical: reqCanon,
	}, nil
}

// ExtractEthCall extracts the eth_call gethapi.TransactionArgs from an interface{}
func ExtractEthCall(param interface{}) (*gethapi.TransactionArgs, error) {
	// geth lowercases the field name and uses the last seen value
	var valString string
	var to, from *gethcommon.Address
	var data *hexutil.Bytes
	var value, gasPrice, maxFeePerGas, maxPriorityFeePerGas *hexutil.Big
	var ok bool
	zeroUint := hexutil.Uint64(0)
	nonce := &zeroUint
	// if gas is not set it should be null
	gas := (*hexutil.Uint64)(nil)

	ethCallMap, ok := param.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid eth call parameter")
	}

	for field, val := range ethCallMap {
		if val == nil {
			continue
		}
		valString, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type supplied in `%s` field", field)
		}
		if len(strings.TrimSpace(valString)) == 0 {
			continue
		}
		switch strings.ToLower(field) {
		case callFieldTo:
			toVal := gethcommon.HexToAddress(valString)
			to = &toVal
		case CallFieldFrom:
			fromVal := gethcommon.HexToAddress(valString)
			from = &fromVal
		case callFieldData, callFieldInput:
			dataVal, err := hexutil.Decode(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode data in CallMsg - %w", err)
			}
			data = (*hexutil.Bytes)(&dataVal)
		case callFieldValue:
			valueVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			value = (*hexutil.Big)(valueVal)
		case callFieldNonce:
			nonceVal, err := hexutil.DecodeUint64(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			nonce = (*hexutil.Uint64)(&nonceVal)
		case callFieldGas:
			gasVal, err := hexutil.DecodeUint64(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			gas = (*hexutil.Uint64)(&gasVal)

		case callFieldGasPrice:
			gasPriceVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			gasPrice = (*hexutil.Big)(gasPriceVal)

		case callFieldMaxFeePerGas:
			maxFeePerGasVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			maxFeePerGas = (*hexutil.Big)(maxFeePerGasVal)

		case callFieldMaxPriorityFeePerGas:
			maxPriorityFeePerGasVal, err := hexutil.DecodeBig(valString)
			if err != nil {
				return nil, fmt.Errorf("could not decode value in CallMsg - %w", err)
			}
			maxPriorityFeePerGas = (*hexutil.Big)(maxPriorityFeePerGasVal)
		case callFieldAccessList:
			// ignore access list for now
		}
	}

	// convert the params[0] into an ethereum.CallMsg
	callMsg := &gethapi.TransactionArgs{
		From:                 from,
		To:                   to,
		Gas:                  gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		Value:                value,
		Data:                 data,
		Nonce:                nonce,
		AccessList:           nil,
	}

	return callMsg, nil
}

// CreateEthHeaderForBatch - the EVM requires an Ethereum header.
// We convert the Batch headers to Ethereum headers to be able to use the Geth EVM.
// Special care must be taken to maintain a valid chain of these converted headers.
func (enc *gethEncodingServiceImpl) CreateEthHeaderForBatch(ctx context.Context, h *common.BatchHeader) (*types.Header, error) {
	// wrap in a caching layer
	return enc.cachingService.ReadConvertedHeader(ctx, h.Hash(), func() (*types.Header, error) {
		// deterministically calculate the private randomness that will be exposed to the EVM
		perBatchRandomness := enc.entropyService.BatchEntropy(h)

		// calculate the converted hash of the parent, for a correct converted chain
		// default to the genesis
		convertedParentHash := common.GethGenesisParentHash
		var err error
		if h.SequencerOrderNo.Uint64() > common.L2GenesisSeqNo {
			convertedParentHash, err = enc.storage.FetchConvertedHash(ctx, h.ParentHash)
			if err != nil {
				enc.logger.Error("Cannot find the converted value for the parent of", log.BatchSeqNoKey, h.SequencerOrderNo)
				return nil, err
			}
		}

		baseFee := uint64(0)
		if h.BaseFee != nil {
			baseFee = h.BaseFee.Uint64()
		}

		gethHeader := types.Header{
			ParentHash:      convertedParentHash,
			UncleHash:       gethutil.EmptyHash,
			Root:            h.Root,
			TxHash:          h.TxHash,
			ReceiptHash:     h.ReceiptHash,
			Difficulty:      big.NewInt(0),
			Number:          h.Number,
			GasLimit:        h.GasLimit,
			GasUsed:         h.GasUsed,
			BaseFee:         big.NewInt(0).SetUint64(baseFee),
			Coinbase:        h.Coinbase,
			Time:            h.Time,
			MixDigest:       perBatchRandomness,
			Nonce:           types.BlockNonce{},
			Extra:           h.SequencerOrderNo.Bytes(),
			WithdrawalsHash: nil,
			BlobGasUsed:     nil,
			ExcessBlobGas:   nil,
			Bloom:           types.Bloom{},
		}
		enc.cachingService.CacheConvertedHash(ctx, h.Hash(), gethHeader.Hash())
		return &gethHeader, nil
	})
}

// The Geth "Block" type doesn't expose the header directly.
// This type is required for adjusting the header.
type localBlock struct {
	header       *types.Header
	uncles       []*types.Header
	transactions types.Transactions
	withdrawals  types.Withdrawals

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

func (enc *gethEncodingServiceImpl) CreateEthBlockFromBatch(ctx context.Context, b *core.Batch) (*types.Block, error) {
	blockHeader, err := enc.CreateEthHeaderForBatch(ctx, b.Header)
	if err != nil {
		return nil, fmt.Errorf("unable to create eth block from batch - %w", err)
	}

	// adjust the header of the returned block to make sure the hashes align
	lb := localBlock{
		header:       blockHeader,
		uncles:       nil,
		transactions: b.Transactions,
		withdrawals:  nil,
	}
	// cast the correct local structure to the standard geth block.
	return (*types.Block)(unsafe.Pointer(&lb)), nil
}

// ExtractPrivateTransactionsQuery is designed to support a wide range of custom TEN queries.
// The first parameter here is the method name, which is used to determine the query type.
// The second parameter is the query parameters.
func ExtractPrivateTransactionsQuery(queryParams any) (*common.ListPrivateTransactionsQueryParams, error) {
	// we expect second param to be a json string
	queryParamsStr, ok := queryParams.(string)
	if !ok {
		return nil, fmt.Errorf("expected queryParams as string but was type %T", queryParams)
	}

	var privateQueryParams common.ListPrivateTransactionsQueryParams
	err := json.Unmarshal([]byte(queryParamsStr), &privateQueryParams)
	if err != nil {
		// if it fails, check if the string was base64 encoded
		bytesStr, err64 := base64.StdEncoding.DecodeString(queryParamsStr)
		if err64 != nil {
			// was not base64 encoded, give up
			return nil, fmt.Errorf("unable to unmarshal params string: %w", err)
		}
		// was base64 encoded, try to unmarshal
		err = json.Unmarshal(bytesStr, &privateQueryParams)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal params string: %w", err)
		}
	}

	return &privateQueryParams, nil
}
