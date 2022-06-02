package datagenerator

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/obscuronode/config"

	"github.com/obscuronet/go-obscuro/go/obscuronode/wallet"
)

// RandomWallet returns a wallet with a random private key
func RandomWallet(chainID int64) wallet.Wallet {
	pk, err := randomHex(32)
	if err != nil {
		panic(err) // this should never panic - world should stop if it does
	}
	walletConfig := config.HostConfig{
		PrivateKeyString: pk,
		ChainID:          *big.NewInt(chainID),
	}
	return wallet.NewInMemoryWalletFromString(walletConfig)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
