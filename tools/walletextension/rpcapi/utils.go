package rpcapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	tenrpc "github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	notAuthorised = "not authorised"
	serverBusy    = "server busy. please retry later"

	// hardcoding the maximum time for an RPC request
	// this value will be propagated to the node and enclave and all the operations
	maximumRPCCallDuration  = 5 * time.Second
	sendTransactionDuration = 20 * time.Second
)

var rpcNotImplemented = fmt.Errorf("rpc endpoint not implemented")

type AuthExecCfg struct {
	// these 4 fields specify the account(s) that should make the backend call
	account             *gethcommon.Address
	computeFromCallback func(user *common.GWUser) *gethcommon.Address
	tryAll              bool
	tryUntilAuthorised  bool

	adjustArgs func(acct *common.GWAccount) []any
	cacheCfg   *cache.Cfg
	timeout    time.Duration
}

func UnauthenticatedTenRPCCall[R any](ctx context.Context, w *services.Services, cfg *cache.Cfg, method string, args ...any) (*R, error) {
	if ctx == nil {
		return nil, errors.New("invalid call. nil Context")
	}
	audit(w, "RPC start method=%s args=%v", method, args)
	requestStartTime := time.Now()
	cacheArgs := []any{method}
	cacheArgs = append(cacheArgs, args...)

	res, err := cache.WithCache(w.RPCResponsesCache, cfg, generateCacheKey(cacheArgs), func() (*R, error) {
		return services.WithPlainRPCConnection(ctx, w.BackendRPC, func(client *rpc.Client) (*R, error) {
			var resp *R = new(R)
			var err error

			// wrap the context with a timeout to prevent long executions
			timeoutContext, cancelCtx := context.WithTimeout(ctx, maximumRPCCallDuration)
			defer cancelCtx()

			err = client.CallContext(timeoutContext, &resp, method, args...)
			return resp, err
		})
	})

	if err != nil {
		audit(w, "RPC call failed. method=%s args=%v error=%+v time=%d", method, args, err, time.Since(requestStartTime).Milliseconds())
		return nil, err
	}

	audit(w, "RPC call succeeded. method=%s args=%v result=%+v time=%d", method, args, res, time.Since(requestStartTime).Milliseconds())
	return res, err
}

func ExecAuthRPC[R any](ctx context.Context, w *services.Services, cfg *AuthExecCfg, method string, args ...any) (*R, error) {
	audit(w, "RPC start method=%s args=%v", method, args)
	requestStartTime := time.Now()
	user, err := extractUserForRequest(ctx, w)
	if err != nil {
		return nil, err
	}

	w.MetricsTracker.RecordUserActivity(hexutils.BytesToHex(user.ID))

	rateLimitAllowed, requestUUID := w.RateLimiter.Allow(gethcommon.Address(user.ID))
	defer w.RateLimiter.SetRequestEnd(gethcommon.Address(user.ID), requestUUID)
	if !rateLimitAllowed {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	cacheArgs := []any{user.ID, method}
	cacheArgs = append(cacheArgs, args...)

	res, err := cache.WithCache(w.RPCResponsesCache, cfg.cacheCfg, generateCacheKey(cacheArgs), func() (*R, error) {
		// determine candidate "from"
		candidateAccts, err := getCandidateAccounts(user, w, cfg)
		if err != nil {
			return nil, err
		}
		if len(candidateAccts) == 0 {
			return nil, fmt.Errorf("illegal access")
		}

		var rpcErr error
		for i := range candidateAccts {
			acct := candidateAccts[i]
			result, err := services.WithEncRPCConnection(ctx, w.BackendRPC, acct, func(rpcClient *tenrpc.EncRPCClient) (*R, error) {
				var result *R
				adjustedArgs := args
				if cfg.adjustArgs != nil {
					adjustedArgs = cfg.adjustArgs(acct)
				}

				// wrap the context with a timeout to prevent long executions
				deadline := cfg.timeout
				// if not set, use default
				if deadline == 0 {
					deadline = maximumRPCCallDuration
				}
				timeoutContext, cancelCtx := context.WithTimeout(ctx, deadline)
				defer cancelCtx()

				err := rpcClient.CallContext(timeoutContext, &result, method, adjustedArgs...)
				// return a friendly error to the user
				if err != nil && errors.Is(err, context.DeadlineExceeded) {
					return nil, fmt.Errorf(serverBusy)
				}
				return result, err
			})
			if err != nil {
				// for calls where we know the expected error we can return early
				if cfg.tryUntilAuthorised && err.Error() != notAuthorised {
					return nil, err
				}
				rpcErr = err
				continue
			}
			return result, nil
		}
		return nil, rpcErr
	})
	audit(w, "RPC call. uid=%s, method=%s args=%v result=%s error=%s time=%d", hexutils.BytesToHex(user.ID), method, args, SafeGenericToString(res), err, time.Since(requestStartTime).Milliseconds())
	return res, err
}

func getCandidateAccounts(user *common.GWUser, we *services.Services, cfg *AuthExecCfg) ([]*common.GWAccount, error) {
	candidateAccts := make([]*common.GWAccount, 0)
	// for users with multiple accounts try to determine a candidate account based on the available information
	switch {
	case cfg.account != nil:
		acc := user.AllAccounts()[*cfg.account]
		if acc != nil {
			candidateAccts = append(candidateAccts, acc)
			return candidateAccts, nil
		}

	case cfg.computeFromCallback != nil:
		suggestedAddress := cfg.computeFromCallback(user)
		if suggestedAddress != nil {
			acc := user.AllAccounts()[*suggestedAddress]
			if acc != nil {
				candidateAccts = append(candidateAccts, acc)
				return candidateAccts, nil
			} else {
				// this can only happen when the "from" is not one of the registered accounts.
				return nil, fmt.Errorf("account: %s not registered to current user. Please register first", suggestedAddress.Hex())
			}
		}
	}

	if cfg.tryAll || cfg.tryUntilAuthorised {
		for _, acc := range user.AllAccounts() {
			candidateAccts = append(candidateAccts, acc)
		}
	}

	return candidateAccts, nil
}

func extractUserID(ctx context.Context, _ *services.Services) ([]byte, error) {
	token, ok := ctx.Value(rpc.GWTokenKey{}).(string)
	if !ok {
		return nil, fmt.Errorf("invalid authentication token: %s", ctx.Value(rpc.GWTokenKey{}))
	}
	userID := gethcommon.FromHex(token)
	if len(userID) != viewingkey.UserIDLength {
		return nil, fmt.Errorf("invalid authentication token: %s", token)
	}
	return userID, nil
}

func extractUserForRequest(ctx context.Context, w *services.Services) (*common.GWUser, error) {
	userID, err := extractUserID(ctx, w)
	if err != nil {
		return nil, err
	}
	user, err := w.Storage.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	return user, nil
}

// generateCacheKey generates a cache key for the given method, encryptionToken and parameters
// encryptionToken is used to generate a unique cache key for each user and empty string should be used for public data
func generateCacheKey(params []any) []byte {
	// Serialize parameters
	rawKey, err := json.Marshal(params)
	if err != nil {
		return nil
	}

	// Optional: Apply hashing
	hasher := sha256.New()
	hasher.Write(rawKey)

	return hasher.Sum(nil)
}

func audit(services *services.Services, msg string, params ...any) {
	if services.Config.VerboseFlag {
		services.Logger().Info(fmt.Sprintf(msg, params...))
	}
}

func cacheBlockNumberOrHash(blockNrOrHash rpc.BlockNumberOrHash) cache.Strategy {
	if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
		return cache.LatestBatch
	}
	return cache.LongLiving
}

func cacheBlockNumber(lastBlock rpc.BlockNumber) cache.Strategy {
	if lastBlock > 0 {
		return cache.LongLiving
	}
	return cache.LatestBatch
}

func SafeGenericToString[R any](r *R) string {
	if r == nil {
		return "nil"
	}

	v := reflect.ValueOf(r).Elem()
	t := v.Type()

	switch v.Kind() {
	case reflect.Struct:
		return structToString(v, t)
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func structToString(v reflect.Value, t reflect.Type) string {
	var parts []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name

		if !fieldType.IsExported() {
			parts = append(parts, fmt.Sprintf("%s: <unexported>", fieldName))
			continue
		}

		fieldStr := fmt.Sprintf("%s: ", fieldName)

		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				fieldStr += "nil"
			} else {
				fieldStr += fmt.Sprintf("%v", field.Elem().Interface())
			}
		case reflect.Slice, reflect.Array:
			if field.Len() > 10 {
				fieldStr += fmt.Sprintf("%v (length: %d)", field.Slice(0, 10).Interface(), field.Len())
			} else {
				fieldStr += fmt.Sprintf("%v", field.Interface())
			}
		case reflect.Struct:
			fieldStr += "{...}" // Avoid recursive calls for nested structs
		default:
			fieldStr += fmt.Sprintf("%v", field.Interface())
		}

		parts = append(parts, fieldStr)
	}

	return fmt.Sprintf("%s{%s}", t.Name(), strings.Join(parts, ", "))
}
