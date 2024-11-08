package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"

	subscriptioncommon "github.com/ten-protocol/go-ten/go/common/subscription"

	tencommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/ratelimiter"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// Services handles the various business logic for the api endpoints
type Services struct {
	HostAddrHTTP      string // The HTTP address on which the TEN host can be reached
	HostAddrWS        string // The WS address on which the TEN host can be reached
	Storage           storage.UserStorage
	logger            gethlog.Logger
	stopControl       *stopcontrol.StopControl
	version           string
	RPCResponsesCache cache.Cache
	BackendRPC        *BackendRPC
	RateLimiter       *ratelimiter.RateLimiter
	Config            *common.Config
	NewHeadsService   *subscriptioncommon.NewHeadsService
}

type NewHeadNotifier interface {
	onNewHead(header *tencommon.BatchHeader)
}

// number of rpc responses to cache
const rpcResponseCacheSize = 1_000_000

func NewServices(hostAddrHTTP string, hostAddrWS string, storage storage.UserStorage, stopControl *stopcontrol.StopControl, version string, logger gethlog.Logger, config *common.Config) *Services {
	newGatewayCache, err := cache.NewCache(rpcResponseCacheSize, logger)
	if err != nil {
		logger.Error(fmt.Errorf("could not create cache. Cause: %w", err).Error())
		panic(err)
	}

	rateLimiter := ratelimiter.NewRateLimiter(config.RateLimitUserComputeTime, config.RateLimitWindow, uint32(config.RateLimitMaxConcurrentRequests), logger)

	services := Services{
		HostAddrHTTP:      hostAddrHTTP,
		HostAddrWS:        hostAddrWS,
		Storage:           storage,
		logger:            logger,
		stopControl:       stopControl,
		version:           version,
		RPCResponsesCache: newGatewayCache,
		BackendRPC:        NewBackendRPC(hostAddrHTTP, hostAddrWS, logger),
		RateLimiter:       rateLimiter,
		Config:            config,
	}

	services.NewHeadsService = subscriptioncommon.NewNewHeadsService(
		func() (chan *tencommon.BatchHeader, <-chan error, error) {
			logger.Info("Connecting to new heads service...")
			// clear the cache to avoid returning stale data during reconnecting.
			services.RPCResponsesCache.EvictShortLiving()
			ch := make(chan *tencommon.BatchHeader)
			errCh, err := subscribeToNewHeadsWithRetry(ch, services, logger)
			logger.Info("Connected to new heads service.", log.ErrKey, err)
			return ch, errCh, err
		},
		true,
		logger,
		func(newHead *tencommon.BatchHeader) error {
			services.RPCResponsesCache.EvictShortLiving()
			return nil
		})

	return &services
}

func subscribeToNewHeadsWithRetry(ch chan *tencommon.BatchHeader, services Services, logger gethlog.Logger) (<-chan error, error) {
	var sub *gethrpc.ClientSubscription
	err := retry.Do(
		func() error {
			connectionObj, err := services.BackendRPC.PlainConnectWs(context.Background())
			if err != nil {
				return fmt.Errorf("cannot fetch rpc connection to backend node %w", err)
			}
			sub, err = connectionObj.Subscribe(context.Background(), rpc.SubscribeNamespace, ch, rpc.SubscriptionTypeNewHeads)
			if err != nil {
				logger.Info("could not subscribe for new head blocks", log.ErrKey, err)
				_ = services.BackendRPC.ReturnConnWS(connectionObj)
			}
			return err
		},
		retry.NewTimeoutStrategy(20*time.Minute, 1*time.Second),
	)
	if err != nil {
		logger.Error("could not subscribe for new head blocks.", log.ErrKey, err)
		return nil, fmt.Errorf("cannot subscribe to new heads to the backend %w", err)
	}

	return sub.Err(), nil
}

// IsStopping returns whether the WE is stopping
func (w *Services) IsStopping() bool {
	return w.stopControl.IsStopping()
}

// Logger returns the WE set logger
func (w *Services) Logger() gethlog.Logger {
	return w.logger
}

// GenerateAndStoreNewUser generates new key-pair and userID, stores it in the database and returns hex encoded userID and error
func (w *Services) GenerateAndStoreNewUser() ([]byte, error) {
	audit(w, "Generating and storing new user")
	requestStartTime := time.Now()
	// generate new key-pair
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		w.Logger().Error(fmt.Sprintf("could not generate new keypair: %s", err))
		return nil, err
	}

	// create UserID and store it in the database with the private key
	userID := viewingkey.CalculateUserID(common.PrivateKeyToCompressedPubKey(viewingPrivateKeyEcies))
	err = w.Storage.AddUser(userID, crypto.FromECDSA(viewingPrivateKeyEcies.ExportECDSA()))
	if err != nil {
		w.Logger().Error(fmt.Sprintf("failed to save user to the database: %s", err))
		return nil, err
	}

	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	audit(w, "Storing new userID: %s, duration: %d ", hexutils.BytesToHex(userID), duration.Milliseconds())
	return userID, nil
}

// AddAddressToUser checks if a message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func (w *Services) AddAddressToUser(userID []byte, address string, signature []byte, signatureType viewingkey.SignatureType) error {
	audit(w, "Adding address to user: %s, address: %s", hexutils.BytesToHex(userID), address)
	requestStartTime := time.Now()
	addressFromMessage := gethcommon.HexToAddress(address)
	// check if a message was signed by the correct address and if the signature is valid
	recoveredAddress, err := viewingkey.CheckSignature(userID, signature, int64(w.Config.TenChainID), signatureType)
	if err != nil {
		return fmt.Errorf("signature is not valid: %w", err)
	}

	if recoveredAddress.Hex() != addressFromMessage.Hex() {
		return fmt.Errorf("invalid request. Signature doesn't match address")
	}

	// register the account for that viewing key
	err = w.Storage.AddAccount(userID, addressFromMessage.Bytes(), signature, signatureType)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error while storing account (%s) for user (%s): %w", addressFromMessage.Hex(), userID, err).Error())
		return err
	}

	audit(w, "Storing new address for user: %s, address: %s, duration: %d ", hexutils.BytesToHex(userID), address, time.Since(requestStartTime).Milliseconds())
	return nil
}

// UserHasAccount checks if provided account exist in the database for given userID
func (w *Services) UserHasAccount(userID []byte, address string) (bool, error) {
	audit(w, "Checking if user has account: %s, address: %s", hexutils.BytesToHex(userID), address)
	addressBytes, err := hex.DecodeString(address[2:]) // remove 0x prefix from address
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", address[2:], err).Error())
		return false, err
	}

	// todo - this can be optimised and done in the database if we will have users with large number of accounts
	// get all the accounts for the selected user
	user, err := w.Storage.GetUser(userID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting accounts for user (%s), %w", userID, err).Error())
		return false, err
	}
	accounts := user.Accounts

	// check if any of the account matches given account
	found := false
	for _, account := range accounts {
		if bytes.Equal(account.Address.Bytes(), addressBytes) {
			found = true
		}
	}
	return found, nil
}

func (w *Services) UserExists(userID []byte) bool {
	audit(w, "Checking if user exists: %s", userID)
	// Check if user exists and don't log error if user doesn't exist, because we expect this to happen in case of
	// user revoking encryption token or using different testnet.
	// todo add a counter here in the future
	user, err := w.Storage.GetUser(userID)
	if err != nil {
		return false
	}

	return len(user.UserKey) > 0
}

func (w *Services) Version() string {
	return w.version
}

func (w *Services) GetTenNodeHealthStatus() (bool, error) {
	audit(w, "Getting TEN node health status")
	res, err := WithPlainRPCConnection[bool](context.Background(), w.BackendRPC, func(client *gethrpc.Client) (*bool, error) {
		res, err := obsclient.NewObsClient(client).Health()
		return &res, err
	})
	return *res, err
}

func (w *Services) GetTenNetworkConfig() (tencommon.TenNetworkInfo, error) {
	audit(w, "Getting TEN network config")
	res, err := WithPlainRPCConnection[tencommon.TenNetworkInfo](context.Background(), w.BackendRPC, func(client *gethrpc.Client) (*tencommon.TenNetworkInfo, error) {
		return obsclient.NewObsClient(client).GetConfig()
	})
	if err != nil {
		return tencommon.TenNetworkInfo{}, err
	}
	return *res, err
}

func (w *Services) GenerateUserMessageToSign(encryptionToken []byte, formatsSlice []string) (string, error) {
	audit(w, "Generating user message to sign")
	// Check if the formats are valid
	for _, format := range formatsSlice {
		if _, exists := viewingkey.SignatureTypeMap[format]; !exists {
			return "", fmt.Errorf("invalid format: %s", format)
		}
	}

	messageFormat := viewingkey.GetBestFormat(formatsSlice)
	message, err := viewingkey.GenerateMessage(encryptionToken, int64(w.Config.TenChainID), viewingkey.PersonalSignVersion, messageFormat)
	if err != nil {
		return "", fmt.Errorf("error generating message: %w", err)
	}
	return string(message), nil
}

func (w *Services) Stop() {
	w.BackendRPC.Stop()
}