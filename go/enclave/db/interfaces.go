package db

import (
	"crypto/ecdsa"

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
	FetchBlock(blockHash common.L1RootHash) (*types.Block, error)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block)
	// ParentBlock returns the L1 Block's parent.
	ParentBlock(block *types.Block) (*types.Block, error)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	// Todo - this is super confusing, analyze the usage
	IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool
	// FetchHeadBlock - returns the head of the current chain.
	FetchHeadBlock() (*types.Block, error)
	// FetchLogs returns the block's logs.
	FetchLogs(blockHash common.L1RootHash) ([]*types.Log, error)
}

type RollupResolver interface {
	// FetchRollup returns the rollup with the given hash.
	FetchRollup(hash common.L2RootHash) (*core.Rollup, error)
	// FetchRollupByHeight returns the rollup with the given height.
	FetchRollupByHeight(height uint64) (*core.Rollup, error)
	// ParentRollup returns the rollup's parent rollup.
	ParentRollup(rollup *core.Rollup) (*core.Rollup, error)
	// StoreGenesisRollupHash stores the hash of the genesis rollup.
	StoreGenesisRollupHash(rollupHash common.L2RootHash) error
	// StoreRollup stores a rollup.
	StoreRollup(rollup *core.Rollup, receipts []*types.Receipt) error
	// FetchGenesisRollup returns the rollup genesis.
	FetchGenesisRollup() (*core.Rollup, error)
	// FetchHeadRollup returns the current head rollup
	FetchHeadRollup() (*core.Rollup, error)
	// Proof - returns the block used as proof for the rollup
	Proof(rollup *core.Rollup) (*types.Block, error)
}

type HeadsAfterL1BlockStorage interface {
	// FetchL2Head returns the hash of the head rollup at a given L1 block.
	FetchL2Head(blockHash common.L1RootHash) (*common.L2RootHash, error)
	// UpdateL1Head updates the L1 head.
	UpdateL1Head(l1Head common.L1RootHash) error
	// UpdateL2Head updates the canonical L2 head for a given L1 block.
	UpdateL2Head(l1Head common.L1RootHash, l2Head *core.Rollup, receipts []*types.Receipt) error
	// CreateStateDB creates a database that can be used to execute transactions
	CreateStateDB(hash common.L2RootHash) (*state.StateDB, error)
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
	GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error)
	// GetTransactionReceipt - returns the receipt of a tx by tx hash
	GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error)
	// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
	GetReceiptsByHash(hash gethcommon.Hash) (types.Receipts, error)
	// GetSender returns the sender of the tx by hash
	GetSender(txHash gethcommon.Hash) (gethcommon.Address, error)
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
	StoreL1Messages(blockHash common.L1RootHash, messages common.CrossChainMessages) error
	GetL1Messages(blockHash common.L1RootHash) (common.CrossChainMessages, error)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	RollupResolver
	SharedSecretStorage
	HeadsAfterL1BlockStorage
	TransactionStorage
	AttestationStorage
	CrossChainMessagesStorage

	// HealthCheck returns whether the storage is deemed healthy or not
	HealthCheck() (bool, error)
}
