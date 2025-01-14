package ethadapter

import (
	"context"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ethereum/go-ethereum"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient defines the interface for RPC communications with the ethereum nodes
// todo (#1617) - some of these methods are composed calls that should be decoupled in the future (ie: BlocksBetween or IsBlockAncestor)
type EthClient interface {
	BlockNumber() (uint64, error)                                                 // retrieves the number of the head block
	HeaderByHash(id gethcommon.Hash) (*types.Header, error)                       // retrieves a block header given a hash
	BlockByHash(id gethcommon.Hash) (*types.Block, error)                         // retrieves a block given a hash
	HeaderByNumber(n *big.Int) (*types.Header, error)                             // retrieves a block given a number - returns head block if n is nil
	SendTransaction(signedTx *types.Transaction) error                            // issues an ethereum transaction (expects signed tx)
	TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error)              // fetches the ethereum transaction receipt
	TransactionByHash(hash gethcommon.Hash) (*types.Transaction, bool, error)     // fetches the ethereum tx
	Nonce(address gethcommon.Address) (uint64, error)                             // fetches the account nonce to use in the next transaction
	BalanceAt(account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) // fetches the balance of the account
	GetLogs(q ethereum.FilterQuery) ([]types.Log, error)                          // fetches the logs for a given query

	Info() Info                                                                 // retrieves the node Info
	FetchHeadBlock() (*types.Header, error)                                     // retrieves the block at head height
	BlocksBetween(block *types.Header, head *types.Header) []*types.Header      // returns the blocks between two blocks
	IsBlockAncestor(block *types.Header, maybeAncestor common.L1BlockHash) bool // returns if the node considers a block the ancestor
	BlockListener() (chan *types.Header, ethereum.Subscription)                 // subscribes to new blocks and returns a listener with the blocks heads and the subscription handler

	CallContract(msg ethereum.CallMsg) ([]byte, error) // Runs the provided call message on the latest block.

	// PrepareTransactionToSend updates the tx with from address, current nonce and current estimates for the gas and the gas price
	PrepareTransactionToSend(ctx context.Context, txData types.TxData, from gethcommon.Address) (types.TxData, error)
	PrepareTransactionToRetry(ctx context.Context, txData types.TxData, from gethcommon.Address, nonce uint64, retries int) (types.TxData, error)

	FetchLastBatchSeqNo(address gethcommon.Address) (*big.Int, error)

	Stop() // tries to cleanly stop the client and release any resources

	EthClient() *ethclient.Client // returns the underlying eth client
	ReconnectIfClosed() error     // closes and creates a new connection
	Alive() bool                  // returns whether the connection is live or not
}

// Info forces the RPC EthClient to return the data in the same format (independently of its implementation)
type Info struct {
	L2ID gethcommon.Address // the address of the Obscuro node this client is dedicated to
}
