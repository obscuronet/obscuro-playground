package gethencoding

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	// The relevant fields in an eth_call request's params.
	callFieldTo                   = "to"
	CallFieldFrom                 = "from"
	callFieldData                 = "input"
	callFieldValue                = "value"
	callFieldGas                  = "gas"
	callFieldNonce                = "nonce"
	callFieldGasPrice             = "gasprice"
	callFieldMaxFeePerGas         = "maxfeepergas"
	callFieldMaxPriorityFeePerGas = "maxpriorityfeepergas"
)

// EncodingService handles conversion to Geth data structures
type EncodingService interface {
	CreateEthHeaderForBatch(h *common.BatchHeader) (*types.Header, error)
	CreateEthBlockFromBatch(b *core.Batch) (*types.Block, error)
}

type gethEncodingServiceImpl struct {
	// conversion is expensive. Cache the converted headers. The key is the hash of the batch.
	gethHeaderCache *cache.Cache[*types.Header]

	storage storage.Storage
	logger  gethlog.Logger
}

func NewGethEncodingService(storage storage.Storage, logger gethlog.Logger) EncodingService {
	// todo (tudor) figure out the best values
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,    // number of keys to track frequency of.
		MaxCost:     1 << 28, // maximum cost of cache (256MB).
		BufferItems: 64,      // number of keys per Get buffer. Todo - what is this
	})
	if err != nil {
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	return &gethEncodingServiceImpl{
		gethHeaderCache: cache.New[*types.Header](ristrettoStore),
		storage:         storage,
		logger:          logger,
	}
}

// ExtractEthCallMapString extracts the eth_call gethapi.TransactionArgs from an interface{}
// it ensures that :
// - All types are string
// - All keys are lowercase
// - There is only one key per value
// - From field is set by default
func ExtractEthCallMapString(paramBytes interface{}) (map[string]string, error) {
	// geth lowercase the field name and uses the last seen value
	var valString string
	var ok bool
	callMsg := map[string]string{
		// From field is set by default
		"from": gethcommon.HexToAddress("0x0").Hex(),
	}
	for field, val := range paramBytes.(map[string]interface{}) {
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
			callMsg[callFieldTo] = valString
		case CallFieldFrom:
			callMsg[CallFieldFrom] = valString
		case callFieldData:
			callMsg[callFieldData] = valString
		case callFieldValue:
			callMsg[callFieldValue] = valString
		case callFieldGas:
			callMsg[callFieldGas] = valString
		case callFieldMaxFeePerGas:
			callMsg[callFieldMaxFeePerGas] = valString
		case callFieldMaxPriorityFeePerGas:
			callMsg[callFieldMaxPriorityFeePerGas] = valString
		default:
			callMsg[field] = valString
		}
	}

	return callMsg, nil
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

	addr := gethcommon.HexToAddress(param.(string))
	return &addr, nil
}

// ExtractOptionalBlockNumber defaults nil or empty block number params to latest block number
func ExtractOptionalBlockNumber(params []interface{}, idx int) (*gethrpc.BlockNumber, error) {
	if len(params) <= idx {
		return ExtractBlockNumber("latest")
	}
	if params[idx] == nil {
		return ExtractBlockNumber("latest")
	}
	if emptyStr, ok := params[idx].(string); ok && len(strings.TrimSpace(emptyStr)) == 0 {
		return ExtractBlockNumber("latest")
	}

	return ExtractBlockNumber(params[idx])
}

// ExtractBlockNumber returns a gethrpc.BlockNumber given an interface{}, errors if unexpected values are used
func ExtractBlockNumber(param interface{}) (*gethrpc.BlockNumber, error) {
	if param == nil {
		return nil, errutil.ErrNotFound
	}

	blockNumber := gethrpc.BlockNumber(0)
	err := blockNumber.UnmarshalJSON([]byte(param.(string)))
	if err != nil {
		return nil, fmt.Errorf("could not parse requested rollup number %s - %w", param.(string), err)
	}

	return &blockNumber, err
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

	for field, val := range param.(map[string]interface{}) {
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
		case callFieldData:
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
func (enc *gethEncodingServiceImpl) CreateEthHeaderForBatch(h *common.BatchHeader) (*types.Header, error) {
	// wrap in a caching layer
	return common.GetCachedValue(enc.gethHeaderCache, enc.logger, h.Hash(), func(a any) (*types.Header, error) {
		// deterministically calculate the private randomness that will be exposed to the EVM
		secret, err := enc.storage.FetchSecret()
		if err != nil {
			enc.logger.Crit("Could not fetch shared secret. Exiting.", log.ErrKey, err)
		}
		perBatchRandomness := crypto.CalculateRootBatchEntropy(secret[:], h.Number)

		// calculate the converted hash of the parent, for a correct converted chain
		// default to the genesis
		convertedParentHash := common.GethGenesisParentHash

		if h.SequencerOrderNo.Uint64() > common.L2GenesisSeqNo {
			convertedParentHash, err = enc.storage.FetchConvertedHash(h.ParentHash)
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
			UncleHash:       gethcommon.Hash{},
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

func (enc *gethEncodingServiceImpl) CreateEthBlockFromBatch(b *core.Batch) (*types.Block, error) {
	blockHeader, err := enc.CreateEthHeaderForBatch(b.Header)
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

// DecodeParamBytes decodes the parameters byte array into a slice of interfaces
// Helps each calling method to manage the positional data
func DecodeParamBytes(paramBytes []byte) ([]interface{}, error) {
	var paramList []interface{}

	if err := json.Unmarshal(paramBytes, &paramList); err != nil {
		return nil, fmt.Errorf("unable to unmarshal params - %w", err)
	}
	return paramList, nil
}

// ExtractViewingKey returns the viewingkey pubkey and the signature from the request
func ExtractViewingKey(vkBytesIntf interface{}) ([]byte, []byte, error) {
	vkBytesList, ok := vkBytesIntf.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("unable to cast the vk to []interface")
	}

	if len(vkBytesList) != 2 {
		return nil, nil, fmt.Errorf("wrong size of viewing key params")
	}

	vkPubkeyHexBytes, err := hexutil.Decode(vkBytesList[0].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode data in vk pub key - %w", err)
	}

	accountSignatureHexBytes, err := hexutil.Decode(vkBytesList[1].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode data in vk signature - %w", err)
	}

	return vkPubkeyHexBytes, accountSignatureHexBytes, nil
}

func ExtractPrivateCustomQuery(_ interface{}, query interface{}) (*common.PrivateCustomQueryListTransactions, error) {
	// Convert the map to a JSON string
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	var result common.PrivateCustomQueryListTransactions
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
