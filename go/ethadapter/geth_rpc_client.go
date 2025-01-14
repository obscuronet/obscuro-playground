package ethadapter

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/TwiN/gocache/v2"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
)

const (
	connRetryMaxWait        = 10 * time.Minute // after this duration, we will stop retrying to connect and return the failure
	connRetryInterval       = 500 * time.Millisecond
	_maxRetryPriceIncreases = 5
	_retryPriceMultiplier   = 2  // geth now wants 100% increase or for you to wait
	_defaultBlockCacheSize  = 51 // enough for 50 request batch size and one for previous block

)

// geth enforces a 1 gwei minimum for blob tx fee
var minBlobTxFee = big.NewInt(params.GWei)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client     *ethclient.Client // the underlying eth rpc client
	timeout    time.Duration     // the timeout for connecting to, or communicating with, the L1 node
	logger     gethlog.Logger
	rpcURL     string
	blockCache *gocache.Cache
}

// NewEthClientFromURL instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClientFromURL(rpcURL string, timeout time.Duration, logger gethlog.Logger) (EthClient, error) {
	client, err := connect(rpcURL, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcURL, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s", rpcURL))

	return &gethRPCClient{
		client:     client,
		timeout:    timeout,
		logger:     logger,
		rpcURL:     rpcURL,
		blockCache: newFifoCache(_defaultBlockCacheSize, 5*time.Minute),
	}, nil
}

func newFifoCache(nrElem int, ttl time.Duration) *gocache.Cache {
	cache := gocache.NewCache().WithMaxSize(nrElem).WithEvictionPolicy(gocache.FirstInFirstOut).WithDefaultTTL(ttl)
	err := cache.StartJanitor()
	if err != nil {
		panic("failed to start cache.")
	}
	return cache
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(ipaddress string, port uint, timeout time.Duration, logger gethlog.Logger) (EthClient, error) {
	return NewEthClientFromURL(fmt.Sprintf("ws://%s:%d", ipaddress, port), timeout, logger)
}

func (e *gethRPCClient) FetchHeadBlock() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.HeaderByNumber(ctx, nil)
}

func (e *gethRPCClient) Info() Info {
	return Info{}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Header, lastBlock *types.Header) []*types.Header {
	var blocksBetween []*types.Header
	var err error

	for currentBlk := lastBlock; currentBlk != nil && !bytes.Equal(currentBlk.Hash().Bytes(), startingBlock.Hash().Bytes()) && !bytes.Equal(currentBlk.ParentHash.Bytes(), gethcommon.HexToHash("").Bytes()); {
		c := currentBlk.ParentHash
		currentBlk, err = e.HeaderByHash(currentBlk.ParentHash)
		if err != nil {
			e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", c.String()), log.ErrKey, err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Header, maybeAncestor common.L1BlockHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) || bytes.Equal(maybeAncestor.Bytes(), (common.L1BlockHash{}).Bytes()) {
		return true
	}

	if block.Number.Int64() == int64(common.L1GenesisHeight) {
		return false
	}

	resolvedAncestorBlock, err := e.HeaderByHash(maybeAncestor)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", maybeAncestor.String()), log.ErrKey, err)
	}
	if resolvedAncestorBlock.Number.Int64() >= block.Number.Int64() {
		return false
	}

	p, err := e.HeaderByHash(block.ParentHash)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s", block.ParentHash.String()), log.ErrKey, err)
	}
	if p == nil {
		return false
	}

	return e.IsBlockAncestor(p, maybeAncestor)
}

func (e *gethRPCClient) SendTransaction(signedTx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.SendTransaction(ctx, signedTx)
}

func (e *gethRPCClient) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.TransactionReceipt(ctx, hash)
}

func (e *gethRPCClient) TransactionByHash(hash gethcommon.Hash) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.TransactionByHash(ctx, hash)
}

func (e *gethRPCClient) Nonce(account gethcommon.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.PendingNonceAt(ctx, account)
}

func (e *gethRPCClient) BlockListener() (chan *types.Header, ethereum.Subscription) {
	// we do not buffer here, we expect the consumer to always be ready to receive new blocks and not fall behind
	ch := make(chan *types.Header)
	var sub ethereum.Subscription
	var err error
	err = retry.Do(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		sub, err = e.client.SubscribeNewHead(ctx, ch)
		if err != nil {
			e.logger.Warn("could not subscribe for new head blocks", log.ErrKey, err)
		}
		return err
	}, retry.NewTimeoutStrategy(connRetryMaxWait, connRetryInterval))
	if err != nil {
		// todo (#1638) - handle this scenario better. Health monitor to report L1 unavailable to node operator, be able to recover without restarting host.
		// couldn't connect after timeout period, will not continue
		e.logger.Crit("could not subscribe for new head blocks.", log.ErrKey, err)
	}

	return ch, sub
}

func (e *gethRPCClient) BlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockNumber(ctx)
}

func (e *gethRPCClient) HeaderByNumber(n *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	b, err := e.client.BlockByNumber(ctx, n)
	if err != nil {
		return nil, err
	}
	return b.Header(), nil
}

func (e *gethRPCClient) HeaderByHash(hash gethcommon.Hash) (*types.Header, error) {
	cachedBlock, found := e.blockCache.Get(hash.Hex())
	if found {
		h, ok := cachedBlock.(types.Header)
		if !ok {
			return nil, fmt.Errorf("should not happen. could not cast cached block to header")
		}
		return &h, nil
	}

	// not in cache, fetch from RPC
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.blockCache.Set(hash.Hex(), *block.Header())
	return block.Header(), nil
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByHash(ctx, hash)
}

func (e *gethRPCClient) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.CallContract(ctx, msg, nil)
}

func (e *gethRPCClient) EthClient() *ethclient.Client {
	return e.client
}

func (e *gethRPCClient) BalanceAt(address gethcommon.Address, blockNum *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BalanceAt(ctx, address, blockNum)
}

func (e *gethRPCClient) GetLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.FilterLogs(ctx, q)
}

func (e *gethRPCClient) Stop() {
	e.blockCache.StopJanitor()
	e.client.Close()
}

func (e *gethRPCClient) FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error) {
	contract, err := ManagementContract.NewManagementContract(address, e.EthClient())
	if err != nil {
		return nil, err
	}

	return contract.LastBatchSeqNo(&bind.CallOpts{})
}

// PrepareTransactionToSend takes a txData type and overrides the From, Gas and Gas Price field with current values
func (e *gethRPCClient) PrepareTransactionToSend(ctx context.Context, txData types.TxData, from gethcommon.Address) (types.TxData, error) {
	nonce, err := e.EthClient().PendingNonceAt(ctx, from)
	if err != nil {
		return nil, fmt.Errorf("could not get nonce - %w", err)
	}
	return e.PrepareTransactionToRetry(ctx, txData, from, nonce, 0)
}

// PrepareTransactionToRetry takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func (e *gethRPCClient) PrepareTransactionToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
	switch tx := txData.(type) {
	case *types.LegacyTx:
		return e.prepareLegacyTxToRetry(ctx, tx, from, nonce, retryNumber)
	case *types.BlobTx:
		return e.prepareBlobTxToRetry(ctx, tx, from, nonce, retryNumber)
	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", tx)
	}
}

func (e *gethRPCClient) prepareLegacyTxToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
	unEstimatedTx := types.NewTx(txData)
	gasPrice, err := e.EthClient().SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not suggest gas price - %w", err)
	}

	// it should never happen but to avoid any risk of repeated price increases we cap the possible retry price bumps to 5
	retryFloat := math.Min(_maxRetryPriceIncreases, float64(retryNumber))
	// we apply a 20% gas price increase for each retry (retrying with similar price gets rejected by mempool)
	// Retry '0' is the first attempt, gives multiplier of 1.0
	multiplier := math.Pow(_retryPriceMultiplier, retryFloat)

	gasPriceFloat := new(big.Float).SetInt(gasPrice)
	retryPriceFloat := big.NewFloat(0).Mul(gasPriceFloat, big.NewFloat(multiplier))
	// prices aren't big enough for float error to matter
	retryPrice, _ := retryPriceFloat.Int(nil)

	gasLimit, err := e.EthClient().EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    unEstimatedTx.To(),
		Value: unEstimatedTx.Value(),
		Data:  unEstimatedTx.Data(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not estimate gas - %w", err)
	}

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: retryPrice,
		Gas:      gasLimit,
		To:       unEstimatedTx.To(),
		Value:    unEstimatedTx.Value(),
		Data:     unEstimatedTx.Data(),
	}, nil
}

// PrepareBlobTransactionToRetry takes a txData type and overrides the From, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func (e *gethRPCClient) prepareBlobTxToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
	unEstimatedTx := types.NewTx(txData)
	to := unEstimatedTx.To()
	value := unEstimatedTx.Value()
	data := unEstimatedTx.Data()
	gasPrice, err := e.EthClient().SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not suggest gas price - %w", err)
	}

	// TODO move all this into common function
	// it should never happen but to avoid any risk of repeated price increases we cap the possible retry price bumps to 5
	retryFloat := math.Min(_maxRetryPriceIncreases, float64(retryNumber))
	// we apply a 20% gas price increase for each retry (retrying with similar price gets rejected by mempool)
	// Retry '0' is the first attempt, gives multiplier of 1.0
	multiplier := math.Pow(_retryPriceMultiplier, retryFloat)

	gasPriceFloat := new(big.Float).SetInt(gasPrice)
	retryPriceFloat := big.NewFloat(0).Mul(gasPriceFloat, big.NewFloat(multiplier))
	// prices aren't big enough for float error to matter
	retryPrice, _ := retryPriceFloat.Int(nil)

	gasLimit, err := e.EthClient().EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return nil, fmt.Errorf("could not estimate gas - %w", err)
	}

	// The base blob fee is calculated from the parent header
	head, err := e.EthClient().HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the suggested base fee: %w", err)
	} else if head.BaseFee == nil {
		return nil, fmt.Errorf("txmgr does not support pre-london blocks that do not have a base fee")
	}
	var blobBaseFee *big.Int
	if head.ExcessBlobGas != nil {
		blobBaseFee = eip4844.CalcBlobFee(*head.ExcessBlobGas)
	}
	blobFeeCap := calcBlobFeeCap(blobBaseFee, retryNumber)

	baseFee := head.BaseFee
	gasFeeCap := new(big.Int).Mul(baseFee, big.NewInt(2))
	gasFeeCap.Add(gasFeeCap, retryPrice)

	return &types.BlobTx{
		Nonce:      nonce,
		GasTipCap:  uint256.MustFromBig(retryPrice), // maxPriorityFeePerGas
		GasFeeCap:  uint256.MustFromBig(gasFeeCap),  // maxFeePerGas = (baseFee * 2) + maxPriorityFeePerGas
		Gas:        gasLimit,
		To:         *unEstimatedTx.To(),
		Value:      uint256.MustFromBig(value),
		Data:       unEstimatedTx.Data(),
		BlobFeeCap: uint256.MustFromBig(blobFeeCap),
		BlobHashes: unEstimatedTx.BlobHashes(),
		Sidecar:    unEstimatedTx.BlobTxSidecar(),
	}, nil
}

// ReconnectIfClosed closes the existing client connection and creates a new connection to the same address:port
func (e *gethRPCClient) ReconnectIfClosed() error {
	if e.Alive() {
		// connection is not closed
		return nil
	}
	e.client.Close()

	client, err := connect(e.rpcURL, e.timeout)
	if err != nil {
		return fmt.Errorf("unable to connect to the eth node (%s) - %w", e.rpcURL, err)
	}
	e.client = client
	return nil
}

// Alive tests the client
func (e *gethRPCClient) Alive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	_, err := e.client.BlockNumber(ctx)
	if err != nil {
		e.logger.Error("Unable to fetch BlockNumber rpc endpoint - client connection is in error state")
		return false
	}
	return err == nil
}

func connect(rpcURL string, connectionTimeout time.Duration) (*ethclient.Client, error) {
	var err error
	var c *ethclient.Client
	for start := time.Now(); time.Since(start) < connectionTimeout; time.Sleep(time.Second) {
		c, err = ethclient.Dial(rpcURL)
		if err == nil {
			break
		}
	}

	return c, err
}

// calcBlobFeeCap computes a suggested blob fee cap that is twice the current header's blob base fee
// value, with a minimum value of minBlobTxFee. It also doubles the blob fee cap based on the retry number.
func calcBlobFeeCap(blobBaseFee *big.Int, retryNumber int) *big.Int {
	// Base calculation: twice the current blob base fee
	blobFeeCap := new(big.Int).Mul(blobBaseFee, big.NewInt(2))

	// Ensure the blob fee cap is at least the minimum value
	if blobFeeCap.Cmp(minBlobTxFee) < 0 {
		blobFeeCap.Set(minBlobTxFee)
	}

	// Double the blob fee cap for each retry attempt
	if retryNumber > 0 {
		multiplier := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(retryNumber)), nil)
		blobFeeCap.Mul(blobFeeCap, multiplier)
	}

	return blobFeeCap
}
