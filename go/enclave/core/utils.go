package core

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

func MakeMap(txs []*common.L2Tx) map[gethcommon.Hash]*common.L2Tx {
	m := make(map[gethcommon.Hash]*common.L2Tx)
	for _, tx := range txs {
		m[tx.Hash()] = tx
	}
	return m
}

func ToMap(txs []*common.L2Tx) map[gethcommon.Hash]gethcommon.Hash {
	m := make(map[gethcommon.Hash]gethcommon.Hash)
	for _, tx := range txs {
		m[tx.Hash()] = tx.Hash()
	}
	return m
}

func PrintTxs(txs []*common.L2Tx) (txsString []string) {
	for _, t := range txs {
		txsString = printTx(t, txsString)
	}
	return txsString
}

func printTx(t *common.L2Tx, txsString []string) []string {
	txsString = append(txsString, fmt.Sprintf("%s,", t.Hash().Hex()))
	return txsString
}

// VerifySignature - Checks that the L2Tx has a valid signature.
func VerifySignature(chainID int64, tx *types.Transaction) error {
	signer := types.NewLondonSigner(big.NewInt(chainID))
	_, err := types.Sender(signer, tx)
	return err
}

// GetAuthenticatedSender - Get sender and tx nonce from transaction
func GetAuthenticatedSender(chainID int64, tx *types.Transaction) (*gethcommon.Address, error) {
	signer := types.NewLondonSigner(big.NewInt(chainID))
	sender, err := types.Sender(signer, tx)
	if err != nil {
		fmt.Println("errorGetAuthenticatedSenderTx: ", tx.Hash())
		return nil, err
	}
	return &sender, nil
}
