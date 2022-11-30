package datagenerator

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common"
)

func RandomRollup() common.ExtRollup {
	extRollup := common.ExtRollup{
		Header: &common.Header{
			ParentHash:  randomHash(),
			Agg:         RandomAddress(),
			L1Proof:     randomHash(),
			Root:        randomHash(),
			Number:      big.NewInt(int64(RandomUInt64())),
			Withdrawals: randomWithdrawals(10),
		},
		TxHashes:        []gethcommon.Hash{randomHash()},
		EncryptedTxBlob: RandomBytes(10),
	}
	return extRollup
}
