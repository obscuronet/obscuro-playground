package rpcapi

//goland:noinspection ALL
import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type BlockChainAPI struct {
	we *Services
}

func NewBlockChainAPI(we *Services) *BlockChainAPI {
	return &BlockChainAPI{we}
}

func (api *BlockChainAPI) ChainId() *hexutil.Big { //nolint:stylecheck
	// chainid, _ := UnauthenticatedTenRPCCall[hexutil.Big](nil, api.we, &CacheCfg{TTL: longCacheTTL}, "eth_chainId")
	// return chainid
	chainID := big.NewInt(int64(api.we.Config.TenChainID))
	return (*hexutil.Big)(chainID)
}

func (api *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	nr, err := UnauthenticatedTenRPCCall[hexutil.Uint64](nil, api.we, &CacheCfg{TTL: shortCacheTTL}, "eth_blockNumber")
	if err != nil {
		return hexutil.Uint64(0)
	}
	return *nr
}

func (api *BlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	// todo - how do you handle getBalance for contracts
	return ExecAuthRPC[hexutil.Big](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
						return shortCacheTTL
					}
					return longCacheTTL
				},
			},
			tryUntilAuthorised: true, // the user can request the balance of a contract account
		},
		"eth_getBalance",
		address,
		blockNrOrHash,
	)
}

// Result structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

/*
	func (s *BlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
		// not implemented
		return nil, nil
	}
*/
func (api *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{TTLCallback: func() time.Duration {
		if number > 0 {
			return longCacheTTL
		}
		return shortCacheTTL
	}}, "eth_getHeaderByNumber", number)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	resp, _ := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{TTL: longCacheTTL}, "eth_getHeaderByHash", hash)
	if resp == nil {
		return nil
	}
	return *resp
}

func (api *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](
		ctx,
		api.we,
		&CacheCfg{
			TTLCallback: func() time.Duration {
				if number > 0 {
					return longCacheTTL
				}
				return shortCacheTTL
			},
		}, "eth_getBlockByNumber", number, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{TTL: longCacheTTL}, "eth_getBlockByHash", hash, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
						return shortCacheTTL
					}
					return longCacheTTL
				},
			},
			account: &address,
		},
		"eth_getCode",
		address,
		blockNrOrHash,
	)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// GetStorageAt is repurposed to return the userID
	if address.Hex() == wecommon.GetStorageAtUserIDRequestMethodName {
		userID, err := extractUserID(ctx, api.we)
		if err != nil {
			return nil, err
		}

		_, err = getUser(userID, api.we.Storage)
		if err != nil {
			return nil, err
		}
		return userID, nil
	}

	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, api.we, &ExecCfg{account: &address}, "eth_getStorageAt", address, hexKey, blockNrOrHash)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

/*
	func (s *BlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
		// not implemented
		return nil, nil
	}
*/
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}
type (
	StateOverride  map[common.Address]OverrideAccount
	BlockOverrides struct {
		Number     *hexutil.Big
		Difficulty *hexutil.Big
		Time       *hexutil.Uint64
		GasLimit   *hexutil.Uint64
		Coinbase   *common.Address
		Random     *common.Hash
		BaseFee    *hexutil.Big
	}
)

func (api *BlockChainAPI) Call(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, api.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			TTLCallback: func() time.Duration {
				if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
					return shortCacheTTL
				}
				return longCacheTTL
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *GWAccount) []any {
			// set the from
			args.From = acct.address
			return []any{args, blockNrOrHash, overrides, blockOverrides}
		},
	}, "eth_call", args, blockNrOrHash, overrides, blockOverrides)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) EstimateGas(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	resp, err := ExecAuthRPC[hexutil.Uint64](ctx, api.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			TTLCallback: func() time.Duration {
				if blockNrOrHash != nil && blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
					return shortCacheTTL
				}
				return longCacheTTL
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		// is this a security risk?
		useDefaultUser: true,
	}, "eth_estimateGas", args, blockNrOrHash, overrides)
	if resp == nil {
		return 0, err
	}
	return *resp, err
}

/*
type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	// not implemented
	return nil, nil
}
*/
