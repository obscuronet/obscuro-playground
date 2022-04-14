package playground

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"io/ioutil"
	"math/big"
	"os"
	"path"
)

func newSignedTransaction(blockchain *core.BlockChain, key *ecdsa.PrivateKey, data types.TxData) *types.Transaction {
	signer := types.MakeSigner(blockchain.Config(), blockchain.CurrentBlock().Number())
	tx, err := types.SignNewTx(key, signer, data)
	panicIfErr(err)

	return tx
}

func newChildBlock(parentBlock *types.Block, txs []*types.Transaction, receipts types.Receipts) *types.Block {
	gasUsed := uint64(0)
	for _, tx := range txs {
		gasUsed += tx.Gas()
	}

	header := &types.Header{
		ParentHash: parentBlock.Hash(),
		Root:       statedb.IntermediateRoot,
		Number:     big.NewInt(parentBlock.Number().Int64() + 1),
		GasLimit:   parentBlock.GasLimit() * 2, // todo - joel - required to be set this way, but not sure why
		GasUsed:    gasUsed,
		BaseFee:    big.NewInt(1000000000), // todo - joel - required to be set this way, but not sure why
	}
	block := types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil))
	return block
}

func createBlockchain() (*core.BlockChain, ethdb.Database) {
	dataDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		panic(err)
	}

	db := createDB(dataDir)
	cacheConfig := createCacheConfig(dataDir)
	chainConfig := createChainConfig(db)
	engine := createEngine(dataDir)
	vmConfig := createVMConfig()
	shouldPreserve := createShouldPreserve()
	txLookupLimit := uint64(2_350_000) // Default.

	blockchain, err := core.NewBlockChain(db, cacheConfig, chainConfig, engine, vmConfig, shouldPreserve, &txLookupLimit)
	panicIfErr(err)
	return blockchain, db
}

func createDB(dataDir string) ethdb.Database {
	root := path.Join(dataDir, "geth/chaindata")            // Defaults to `geth/chaindata` in the node's data directory.
	cache := 2048                                           // Default.
	handles := 2048                                         // Default.
	freezer := path.Join(dataDir, "geth/chaindata/ancient") // Defaults to `geth/chaindata/ancient` in the node's data directory.
	namespace := ""                                         // Defaults to `eth/db/chaindata`.
	readonly := false                                       // Default.

	db, err := rawdb.NewLevelDBDatabaseWithFreezer(root, cache, handles, freezer, namespace, readonly)
	panicIfErr(err)
	return db
}

func createCacheConfig(dataDir string) *core.CacheConfig {
	return &core.CacheConfig{
		TrieCleanLimit:      614,                                  // Default.
		TrieCleanJournal:    path.Join(dataDir, "geth/triecache"), // Defaults to `geth/triecache` in the node's data directory.
		TrieCleanRejournal:  3600000000000,                        // Default.
		TrieCleanNoPrefetch: false,                                // Default.
		TrieDirtyLimit:      1024,                                 // Default.
		TrieDirtyDisabled:   false,                                // Default.
		TrieTimeLimit:       3600000000000,                        // Default.
		SnapshotLimit:       409,                                  // Default.
		Preimages:           false,                                // Default.
	}
}

func createChainConfig(db ethdb.Database) *params.ChainConfig {
	chainConfig, _, genesisErr := core.SetupGenesisBlockWithOverride(
		db,
		nil, // Default.
		nil, // Default.
		nil, // Default.
	)
	panicIfErr(genesisErr)
	return chainConfig
}

// Recreates the standard path through `eth/ethconfig/config.go/CreateConsensusEngine()`.
func createEngine(dataDir string) consensus.Engine {
	var engine consensus.Engine
	engine = ethash.New(ethash.Config{
		PowMode:          ethash.ModeNormal,                 // Default.
		CacheDir:         path.Join(dataDir, "geth/ethash"), // Defaults to `geth/ethash` in the node's data directory.
		CachesInMem:      2,                                 // Default.
		CachesOnDisk:     3,                                 // Default.
		CachesLockMmap:   false,                             // Default.
		DatasetDir:       "",                                // Defaults to `~/Library/Ethash` in the node's data directory.
		DatasetsInMem:    1,                                 // Default.
		DatasetsOnDisk:   2,                                 // Default.
		DatasetsLockMmap: false,                             // Default.
		NotifyFull:       false,                             // Default.
	}, nil, false) // Defaults.
	engine.(*ethash.Ethash).SetThreads(-1)
	return beacon.New(engine)
}

func createVMConfig() vm.Config {
	return vm.Config{
		EnablePreimageRecording: false, // Default.
	}
}

// We indicate that no blocks are authored by local accounts, and thus all blocks are discarded during reorgs.
func createShouldPreserve() func(header *types.Header) bool {
	return func(header *types.Header) bool {
		return false
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
