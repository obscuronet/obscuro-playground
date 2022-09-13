package wallet

import (
	"crypto/ecdsa"
	"math/big"
	"sync/atomic"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet interface {
	// Address returns the pubkey Address of the wallet
	Address() common.Address
	// SignTransaction returns a signed transaction
	SignTransaction(tx types.TxData) (*types.Transaction, error)

	// SetNonce overrides the current nonce
	// The GetTransactionCount is expected to be the next nonce to use in a transaction, not the current account GetTransactionCount
	SetNonce(nonce uint64)
	// GetNonceAndIncrement atomically increments the nonce by one and returns the previous value
	GetNonceAndIncrement() uint64
	GetNonce() uint64

	// PrivateKey returns the wallets private key
	PrivateKey() *ecdsa.PrivateKey
}

type inMemoryWallet struct {
	prvKey     *ecdsa.PrivateKey
	pubKey     *ecdsa.PublicKey
	pubKeyAddr common.Address
	nonce      uint64
	chainID    *big.Int
}

func NewInMemoryWalletFromPK(chainID *big.Int, pk *ecdsa.PrivateKey) Wallet {
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Panic("error casting public key to ECDSA")
	}

	return &inMemoryWallet{
		chainID:    chainID,
		prvKey:     pk,
		pubKey:     publicKeyECDSA,
		pubKeyAddr: crypto.PubkeyToAddress(*publicKeyECDSA),
	}
}

func NewInMemoryWalletFromConfig(config config.HostConfig) Wallet {
	privateKey, err := crypto.HexToECDSA(config.PrivateKeyString)
	if err != nil {
		log.Panic("could not recover private key from hex. Cause: %s", err)
	}
	return NewInMemoryWalletFromPK(big.NewInt(config.L1ChainID), privateKey)
}

// SignTransaction returns a signed transaction
func (m *inMemoryWallet) SignTransaction(tx types.TxData) (*types.Transaction, error) {
	return types.SignNewTx(m.prvKey, types.NewLondonSigner(m.chainID), tx)
}

// Address returns the current wallet address
func (m *inMemoryWallet) Address() common.Address {
	return m.pubKeyAddr
}

func (m *inMemoryWallet) GetNonceAndIncrement() uint64 {
	return atomic.AddUint64(&m.nonce, 1) - 1
}

func (m *inMemoryWallet) GetNonce() uint64 {
	return atomic.LoadUint64(&m.nonce)
}

func (m *inMemoryWallet) SetNonce(nonce uint64) {
	m.nonce = nonce
}

func (m *inMemoryWallet) PrivateKey() *ecdsa.PrivateKey {
	return m.prvKey
}
