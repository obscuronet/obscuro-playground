package rpcapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	rpc2 "github.com/ten-protocol/go-ten/go/rpc"

	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/privacy"
	"github.com/ten-protocol/go-ten/go/common/subscription"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type BlockChainAPI struct {
	we               *services.Services
	storageWhitelist *privacy.Whitelist
}

func NewBlockChainAPI(we *services.Services) *BlockChainAPI {
	whitelist := privacy.NewWhitelist()
	return &BlockChainAPI{
		we:               we,
		storageWhitelist: whitelist,
	}
}

func (api *BlockChainAPI) ChainId() *hexutil.Big { //nolint:stylecheck
	chainID, _ := UnauthenticatedTenRPCCall[hexutil.Big](context.Background(), api.we, &cache.Cfg{Type: cache.LongLiving}, "ten_chainId")
	return chainID
}

func (api *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	nr, err := UnauthenticatedTenRPCCall[hexutil.Uint64](context.Background(), api.we, &cache.Cfg{Type: cache.LatestBatch}, rpc2.BatchNumber)
	if err != nil {
		return hexutil.Uint64(0)
	}
	return *nr
}

func (api *BlockChainAPI) GetBalance(ctx context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return nil, err
	}

	// Add detailed logging
	api.we.Logger().Info("GetBalance request received",
		"address", address.Hex(),
		"blockNrOrHash", blockNrOrHash,
		"user", user)

	balance, err := ExecAuthRPC[hexutil.Big](
		ctx,
		api.we,
		&AuthExecCfg{
			cacheCfg: &cache.Cfg{
				DynamicType: func() cache.Strategy {
					return cacheBlockNumberOrHash(blockNrOrHash)
				},
			},
			account:            &address,
			tryUntilAuthorised: true, // the user can request the balance of a contract account
		},
		tenrpc.ERPCGetBalance,
		address,
		blockNrOrHash,
	)

	// Log the result
	if err != nil {
		api.we.Logger().Error("GetBalance request failed",
			"address", address.Hex(),
			"blockNrOrHash", blockNrOrHash,
			"error", err)
	} else {
		api.we.Logger().Info("GetBalance request completed",
			"address", address.Hex(),
			"blockNrOrHash", blockNrOrHash,
			"balance", balance)
	}

	return balance, err
}

// Result structs for GetProof
type AccountResult struct {
	Address      gethcommon.Address `json:"address"`
	AccountProof []string           `json:"accountProof"`
	Balance      *hexutil.Big       `json:"balance"`
	CodeHash     gethcommon.Hash    `json:"codeHash"`
	Nonce        hexutil.Uint64     `json:"nonce"`
	StorageHash  gethcommon.Hash    `json:"storageHash"`
	StorageProof []StorageResult    `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

func (s *BlockChainAPI) GetProof(ctx context.Context, address gethcommon.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	return nil, rpcNotImplemented
}

func (api *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &cache.Cfg{DynamicType: func() cache.Strategy {
		return cacheBlockNumber(number)
	}}, "ten_getHeaderByNumber", number)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash gethcommon.Hash) map[string]interface{} {
	resp, _ := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, "ten_getHeaderByHash", hash)
	if resp == nil {
		return nil
	}
	return *resp
}

func (api *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[common.BatchHeader](
		ctx,
		api.we,
		&cache.Cfg{
			DynamicType: func() cache.Strategy {
				return cacheBlockNumber(number)
			},
		}, rpc2.GetBatchByNumber, number, fullTx)
	if resp == nil {
		return nil, err
	}

	// convert to geth header and marshall
	header := subscription.ConvertBatchHeader(resp)
	fields := RPCMarshalHeader(header)

	// dummy fields
	fields["size"] = hexutil.Uint64(0)
	fields["transactions"] = []any{}

	addExtraTenFields(fields, resp)
	return fields, err
}

func (api *BlockChainAPI) GetBlockByHash(ctx context.Context, hash gethcommon.Hash, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[common.BatchHeader](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, rpc2.GetBatchByHash, hash, fullTx)
	if resp == nil {
		return nil, err
	}

	// convert to geth header and marshall
	header := subscription.ConvertBatchHeader(resp)
	fields := RPCMarshalHeader(header)

	// dummy fields
	fields["size"] = hexutil.Uint64(0)
	fields["transactions"] = []any{}

	addExtraTenFields(fields, resp)
	return fields, err
}

func (api *BlockChainAPI) GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// todo - must be authenticated
	resp, err := UnauthenticatedTenRPCCall[hexutil.Bytes](
		ctx,
		api.we,
		&cache.Cfg{
			DynamicType: func() cache.Strategy {
				return cacheBlockNumberOrHash(blockNrOrHash)
			},
		},
		"ten_getCode",
		address,
		blockNrOrHash,
	)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

// GetStorageAt is not compatible with ETH RPC tooling. TEN network does not getStorageAt because it would
// violate the privacy guarantees of the network.
//
// However, we can repurpose this method to be able to route TEN-specific requests through from an ETH RPC client.
// We call these requests Custom Queries.
//
// This method signature matches eth_getStorageAt, but we use the address field to specify the custom query method,
// the hex-encoded position field to specify the parameters json, and nil for the block number.
//
// In the future, we can support both CustomQueries and some debug version of eth_getStorageAt if needed.
func (api *BlockChainAPI) GetStorageAt(ctx context.Context, address gethcommon.Address, params string, _ rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return nil, err
	}

	switch address.Hex() {
	case common.ListPrivateTransactionsCQMethod:
		// sensitive CustomQuery methods use the convention of having "address" at the top level of the params json
		userAddr, err := extractCustomQueryAddress(params)
		if err != nil {
			return nil, fmt.Errorf("unable to extract address from custom query params: %w", err)
		}
		resp, err := ExecAuthRPC[any](ctx, api.we, &AuthExecCfg{account: userAddr}, "scan_getPersonalTransactions", params)
		if err != nil {
			return nil, fmt.Errorf("unable to execute custom query: %w", err)
		}
		// turn resp object into hexutil.Bytes
		serialised, err := json.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal response object: %w", err)
		}
		return serialised, nil
	case common.CreateSessionKeyCQMethod:
		sk, err := api.we.SKManager.CreateSessionKey(user)
		if err != nil {
			return nil, fmt.Errorf("unable to create session key: %w", err)
		}
		return sk.Account.Address.Bytes(), nil
	case common.ActivateSessionKeyCQMethod:
		res, err := api.we.SKManager.ActivateSessionKey(user)
		return []byte{boolToByte(res)}, err
	case common.DeactivateSessionKeyCQMethod:
		res, err := api.we.SKManager.DeactivateSessionKey(user)
		return []byte{boolToByte(res)}, err
	case common.DeleteSessionKeyCQMethod:
		res, err := api.we.SKManager.DeleteSessionKey(user)
		return []byte{boolToByte(res)}, err
	default: // address was not a recognised custom query method address
		resp, err := ExecAuthRPC[any](ctx, api.we, &AuthExecCfg{tryUntilAuthorised: true}, tenrpc.ERPCGetStorageAt, address, params, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to execute eth_getStorageAt: %w", err)
		}
		if resp == nil {
			return nil, nil
		}

		respHex, ok := (*resp).(string)
		if !ok {
			return nil, fmt.Errorf("unable to decode response")
		}
		// turn resp object into hexutil.Bytes
		return hexutil.MustDecode(respHex), nil
	}
}

func boolToByte(res bool) byte {
	if res {
		return 1
	}
	return 0
}

func (s *BlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	return nil, rpcNotImplemented
}

type OverrideAccount struct {
	Nonce     *hexutil.Uint64                      `json:"nonce"`
	Code      *hexutil.Bytes                       `json:"code"`
	Balance   **hexutil.Big                        `json:"balance"`
	State     *map[gethcommon.Hash]gethcommon.Hash `json:"state"`
	StateDiff *map[gethcommon.Hash]gethcommon.Hash `json:"stateDiff"`
}
type (
	StateOverride  map[gethcommon.Address]OverrideAccount
	BlockOverrides struct {
		Number     *hexutil.Big
		Difficulty *hexutil.Big
		Time       *hexutil.Uint64
		GasLimit   *hexutil.Uint64
		Coinbase   *gethcommon.Address
		Random     *gethcommon.Hash
		BaseFee    *hexutil.Big
	}
)

func (api *BlockChainAPI) Call(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, api.we, &AuthExecCfg{
		cacheCfg: &cache.Cfg{
			DynamicType: func() cache.Strategy {
				return cacheBlockNumberOrHash(blockNrOrHash)
			},
		},
		computeFromCallback: func(user *wecommon.GWUser) *gethcommon.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *wecommon.GWAccount) []any {
			argsClone := populateFrom(acct, args)
			return []any{argsClone, blockNrOrHash, overrides, blockOverrides}
		},
		tryAll: true,
	}, tenrpc.ERPCCall, args, blockNrOrHash, overrides, blockOverrides)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) EstimateGas(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	resp, err := ExecAuthRPC[hexutil.Uint64](ctx, api.we, &AuthExecCfg{
		cacheCfg: &cache.Cfg{
			DynamicType: func() cache.Strategy {
				if blockNrOrHash != nil {
					return cacheBlockNumberOrHash(*blockNrOrHash)
				}
				return cache.LatestBatch
			},
		},
		computeFromCallback: func(user *wecommon.GWUser) *gethcommon.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *wecommon.GWAccount) []any {
			argsClone := populateFrom(acct, args)
			return []any{argsClone, blockNrOrHash, overrides}
		},
		// is this a security risk?
		tryAll: true,
	}, tenrpc.ERPCEstimateGas, args, blockNrOrHash, overrides)
	if resp == nil {
		return 0, err
	}
	return *resp, err
}

func populateFrom(acct *wecommon.GWAccount, args gethapi.TransactionArgs) gethapi.TransactionArgs {
	// clone the args
	argsClone := cloneArgs(args)
	// set the from
	if args.From == nil || args.From.Hex() == (gethcommon.Address{}).Hex() {
		argsClone.From = acct.Address
	}
	return argsClone
}

func cloneArgs(args gethapi.TransactionArgs) gethapi.TransactionArgs {
	serialised, _ := json.Marshal(args)
	var argsClone gethapi.TransactionArgs
	err := json.Unmarshal(serialised, &argsClone)
	if err != nil {
		return gethapi.TransactionArgs{}
	}
	return argsClone
}

type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	return nil, rpcNotImplemented
}

func extractCustomQueryAddress(params any) (*gethcommon.Address, error) {
	// sensitive CustomQuery methods use the convention of having "address" at the top level of the params json
	// we don't care about the params struct overall, just want to extract the address string field
	paramsStr, ok := params.(string)
	if !ok {
		return nil, fmt.Errorf("params must be a json string")
	}
	var paramsJSON map[string]json.RawMessage
	err := json.Unmarshal([]byte(paramsStr), &paramsJSON)
	if err != nil {
		// try to base64 decode the params string and then unmarshal before giving up
		bytesStr, err64 := base64.StdEncoding.DecodeString(paramsStr)
		if err64 != nil {
			// was not base64 encoded, give up
			return nil, fmt.Errorf("unable to unmarshal params string: %w", err)
		}
		// was base64 encoded, try to unmarshal
		err = json.Unmarshal(bytesStr, &paramsJSON)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal params string: %w", err)
		}
	}
	// Extract the RawMessage for the key "address"
	addressRaw, ok := paramsJSON["address"]
	if !ok {
		return nil, fmt.Errorf("params must contain an 'address' field")
	}

	// Unmarshal the RawMessage to a string
	var addressStr string
	err = json.Unmarshal(addressRaw, &addressStr)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal address field to string: %w", err)
	}
	address := gethcommon.HexToAddress(addressStr)
	return &address, nil
}

// RPCMarshalHeader converts the given header to the RPC output .
// duplicated from go-ethereum
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}
	if head.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = head.WithdrawalsHash
	}
	if head.BlobGasUsed != nil {
		result["blobGasUsed"] = hexutil.Uint64(*head.BlobGasUsed)
	}
	if head.ExcessBlobGas != nil {
		result["excessBlobGas"] = hexutil.Uint64(*head.ExcessBlobGas)
	}
	if head.ParentBeaconRoot != nil {
		result["parentBeaconBlockRoot"] = head.ParentBeaconRoot
	}
	return result
}

func addExtraTenFields(fields map[string]interface{}, header *common.BatchHeader) {
	fields["l1Proof"] = header.L1Proof
	fields["signature"] = header.Signature
	fields["crossChainTreeHash"] = header.CrossChainRoot
	fields["crossChainTree"] = header.CrossChainTree
}
