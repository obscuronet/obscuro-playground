package db

import (
	"crypto/ecdsa"
	"io"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash.
	FetchBlock(blockHash common.L1BlockHash) (*types.Block, error)
	// FetchHeadBlock - returns the head of the current chain.
	FetchHeadBlock() (*types.Block, error)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	// todo (low priority) - this is super confusing, analyze the usage
	IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool
}

type BatchResolver interface {
	// FetchBatch returns the batch with the given hash.
	FetchBatch(hash common.L2BatchHash) (*core.Batch, error)
	// FetchBatchByHeight returns the batch on the canonical chain with the given height.
	FetchBatchByHeight(height uint64) (*core.Batch, error)
	// FetchHeadBatch returns the current head batch of the canonical chain.
	FetchHeadBatch() (*core.Batch, error)
	// StoreBatch stores a batch.
	StoreBatch(batch *core.Batch, receipts []*types.Receipt) error
}

type RollupResolver interface {
	// StoreRollup stores a rollup.
	StoreRollup(rollup *core.Rollup) error
}

type HeadsAfterL1BlockStorage interface {
	// FetchHeadBatchForBlock returns the hash of the head batch at a given L1 block.
	FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error)
	// FetchHeadRollupForBlock returns the hash of the latest (i.e. highest-numbered) rollup in the given L1 block, or
	// nil if the block contains no rollups.
	FetchHeadRollupForBlock(blockHash *common.L1BlockHash) (*core.Rollup, error)
	// UpdateL1Head updates the L1 head.
	UpdateL1Head(l1Head common.L1BlockHash) error
	// UpdateHeadBatch updates the canonical L2 head batch for a given L1 block.
	UpdateHeadBatch(l1Head common.L1BlockHash, l2Head *core.Batch, receipts []*types.Receipt) error
	// SetHeadBatchPointer updates the canonical L2 head batch for a given L1 block.
	SetHeadBatchPointer(l2Head *core.Batch) error
	// UpdateHeadRollup just updates the canonical L2 head batch, leaving data untouched (used to rewind after L1 fork or data corruption)
	UpdateHeadRollup(l1Head *common.L1BlockHash, l2Head *common.L2BatchHash) error
	// CreateStateDB creates a database that can be used to execute transactions
	CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error)
	// EmptyStateDB creates the original empty StateDB
	EmptyStateDB() (*state.StateDB, error)
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret.
	FetchSecret() (*crypto.SharedEnclaveSecret, error)
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret crypto.SharedEnclaveSecret) error
}

type TransactionStorage interface {
	// GetTransaction - returns the positional metadata of the tx by hash
	GetTransaction(txHash common.L2TxHash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error)
	// GetTransactionReceipt - returns the receipt of a tx by tx hash
	GetTransactionReceipt(txHash common.L2TxHash) (*types.Receipt, error)
	// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
	GetReceiptsByHash(hash common.L2BatchHash) (types.Receipts, error)
	// GetSender returns the sender of the tx by hash
	GetSender(txHash common.L2TxHash) (gethcommon.Address, error)
	// GetContractCreationTx returns the hash of the tx that created a contract
	GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error)
}

type AttestationStorage interface {
	// FetchAttestedKey returns the public key of an attested aggregator
	FetchAttestedKey(aggregator gethcommon.Address) (*ecdsa.PublicKey, error)
	// StoreAttestedKey - store the public key of an attested aggregator
	StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error
}

type CrossChainMessagesStorage interface {
	StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error
	GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	BatchResolver
	RollupResolver
	SharedSecretStorage
	HeadsAfterL1BlockStorage
	TransactionStorage
	AttestationStorage
	CrossChainMessagesStorage
	io.Closer

	// HealthCheck returns whether the storage is deemed healthy or not
	HealthCheck() (bool, error)
	// FilterLogs - applies the properties the relevancy checks for the requestingAccount to all the stored log events
	// nil values will be ignored. Make sure to set all fields to the right values before calling this function
	// the blockHash should always be nil.
	FilterLogs(requestingAccount *gethcommon.Address, fromBlock, toBlock *big.Int, blockHash *common.L2BatchHash, addresses []gethcommon.Address, topics [][]gethcommon.Hash) ([]*types.Log, error)
}
