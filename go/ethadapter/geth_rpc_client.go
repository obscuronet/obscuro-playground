package ethadapter

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/hashicorp/golang-lru/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
	"github.com/ten-protocol/go-ten/go/common/retry"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	connRetryMaxWait        = 10 * time.Minute // after this duration, we will stop retrying to connect and return the failure
	connRetryInterval       = 500 * time.Millisecond
	_maxRetryPriceIncreases = 5
	_retryPriceMultiplier   = 1.2
	_defaultBlockCacheSize  = 51 // enough for 50 request batch size and one for previous block
)

// gethRPCClient implements the EthClient interface and allows connection to a real ethereum node
type gethRPCClient struct {
	client     *ethclient.Client  // the underlying eth rpc client
	l2ID       gethcommon.Address // the address of the Obscuro node this client is dedicated to
	timeout    time.Duration      // the timeout for connecting to, or communicating with, the L1 node
	logger     gethlog.Logger
	rpcURL     string
	blockCache *lru.Cache[gethcommon.Hash, *types.Block]
}

// NewEthClientFromURL instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClientFromURL(rpcURL string, timeout time.Duration, l2ID gethcommon.Address, logger gethlog.Logger) (EthClient, error) {
	client, err := connect(rpcURL, timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the eth node (%s) - %w", rpcURL, err)
	}

	logger.Trace(fmt.Sprintf("Initialized eth node connection - addr: %s", rpcURL))

	// cache recent blocks to avoid re-fetching them (they are often re-used for checking for forks etc.)
	blkCache, err := lru.New[gethcommon.Hash, *types.Block](_defaultBlockCacheSize)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize block cache - %w", err)
	}

	return &gethRPCClient{
		client:     client,
		l2ID:       l2ID,
		timeout:    timeout,
		logger:     logger,
		rpcURL:     rpcURL,
		blockCache: blkCache,
	}, nil
}

// NewEthClient instantiates a new ethadapter.EthClient that connects to an ethereum node
func NewEthClient(ipaddress string, port uint, timeout time.Duration, l2ID gethcommon.Address, logger gethlog.Logger) (EthClient, error) {
	return NewEthClientFromURL(fmt.Sprintf("ws://%s:%d", ipaddress, port), timeout, l2ID, logger)
}

func (e *gethRPCClient) FetchHeadBlock() (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByNumber(ctx, nil)
}

func (e *gethRPCClient) Info() Info {
	return Info{
		L2ID: e.l2ID,
	}
}

func (e *gethRPCClient) BlocksBetween(startingBlock *types.Block, lastBlock *types.Block) []*types.Block {
	var blocksBetween []*types.Block
	var err error

	for currentBlk := lastBlock; currentBlk != nil && !bytes.Equal(currentBlk.Hash().Bytes(), startingBlock.Hash().Bytes()) && !bytes.Equal(currentBlk.ParentHash().Bytes(), gethcommon.HexToHash("").Bytes()); {
		currentBlk, err = e.BlockByHash(currentBlk.ParentHash())
		if err != nil {
			e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", currentBlk.ParentHash().String()), log.ErrKey, err)
		}
		blocksBetween = append(blocksBetween, currentBlk)
	}

	return blocksBetween
}

func (e *gethRPCClient) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) || bytes.Equal(maybeAncestor.Bytes(), (common.L1BlockHash{}).Bytes()) {
		return true
	}

	if block.Number().Int64() == int64(common.L1GenesisHeight) {
		return false
	}

	resolvedBlock, err := e.BlockByHash(maybeAncestor)
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s.", maybeAncestor.String()), log.ErrKey, err)
	}
	if resolvedBlock == nil {
		if resolvedBlock.Number().Int64() >= block.Number().Int64() {
			return false
		}
	}

	p, err := e.BlockByHash(block.ParentHash())
	if err != nil {
		e.logger.Crit(fmt.Sprintf("could not fetch parent block with hash %s", block.ParentHash().String()), log.ErrKey, err)
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

func (e *gethRPCClient) Nonce(account gethcommon.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.PendingNonceAt(ctx, account)
}

func (e *gethRPCClient) BlockListener() (chan *types.Header, ethereum.Subscription) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// we do not buffer here, we expect the consumer to always be ready to receive new blocks and not fall behind
	ch := make(chan *types.Header, 1)
	var sub ethereum.Subscription
	var err error
	err = retry.Do(func() error {
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

func (e *gethRPCClient) BlockByNumber(n *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return e.client.BlockByNumber(ctx, n)
}

func (e *gethRPCClient) BlockByHash(hash gethcommon.Hash) (*types.Block, error) {
	block, found := e.blockCache.Get(hash)
	if found {
		return block, nil
	}

	// not in cache, fetch from RPC
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.blockCache.Add(hash, block)
	return block, nil
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
	e.client.Close()
}

func (e *gethRPCClient) FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error) {
	contract, err := ManagementContract.NewManagementContract(address, e.EthClient())
	if err != nil {
		return nil, err
	}

	return contract.LastBatchSeqNo(&bind.CallOpts{})
}

// PrepareTransactionToSend takes a txData type and overrides the From, Nonce, Gas and Gas Price field with current values
func (e *gethRPCClient) PrepareTransactionToSend(txData types.TxData, from gethcommon.Address, nonce uint64) (types.TxData, error) {
	return e.PrepareTransactionToRetry(txData, from, nonce, 0)
}

// PrepareTransactionToRetry takes a txData type and overrides the From, Nonce, Gas and Gas Price field with current values
// it bumps the price by a multiplier for retries. retryNumber is zero on first attempt (no multiplier on price)
func (e *gethRPCClient) PrepareTransactionToRetry(txData types.TxData, from gethcommon.Address, nonce uint64, retryNumber int) (types.TxData, error) {
	unEstimatedTx := types.NewTx(txData)
	gasPrice, err := e.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// it should never happen but to avoid any risk of repeated price increases we cap the possible retry price bumps to 5
	retryFloat := math.Max(_maxRetryPriceIncreases, float64(retryNumber))
	// we apply a 20% gas price increase for each retry (retrying with similar price gets rejected by mempool)
	multiplier := math.Pow(_retryPriceMultiplier, retryFloat)

	gasPriceFloat := new(big.Float).SetInt(gasPrice)
	retryPriceFloat := big.NewFloat(0).Mul(gasPriceFloat, big.NewFloat(multiplier))
	// prices aren't big enough for float error to matter
	retryPrice, _ := retryPriceFloat.Int(nil)

	gasLimit, err := e.EthClient().EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    unEstimatedTx.To(),
		Value: unEstimatedTx.Value(),
		Data:  unEstimatedTx.Data(),
	})
	if err != nil {
		return nil, err
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
