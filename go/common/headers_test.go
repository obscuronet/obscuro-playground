package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBatchHeader_MarshalJSON(t *testing.T) {
	batchHeader := &BatchHeader{
		ParentHash:                    randomHash(),
		Root:                          randomHash(),
		TxHash:                        randomHash(),
		ReceiptHash:                   randomHash(),
		Number:                        gethcommon.Big0,
		SequencerOrderNo:              gethcommon.Big1,
		GasLimit:                      100,
		GasUsed:                       200,
		Time:                          300,
		Extra:                         []byte("123"),
		BaseFee:                       gethcommon.Big2,
		L1Proof:                       randomHash(),
		R:                             gethcommon.Big3,
		S:                             gethcommon.Big3,
		CrossChainMessages:            nil,
		LatestInboundCrossChainHash:   gethcommon.Hash{},
		LatestInboundCrossChainHeight: nil,
	}

	jsonMarshalled, err := json.Marshal(batchHeader)
	require.NoError(t, err)

	fmt.Println(string(jsonMarshalled))
	fmt.Println("derp")

	batchUnmarshalled := BatchHeader{}
	err = json.Unmarshal(jsonMarshalled, &batchUnmarshalled)
	require.NoError(t, err)

	require.Equal(t, batchHeader.ParentHash, batchUnmarshalled.ParentHash)
	require.Equal(t, batchHeader.Root, batchUnmarshalled.Root)
	require.Equal(t, batchHeader.TxHash, batchUnmarshalled.TxHash)
	require.Equal(t, batchHeader.ReceiptHash, batchUnmarshalled.ReceiptHash)
	require.Equal(t, batchHeader.Number, batchUnmarshalled.Number)
	require.Equal(t, batchHeader.SequencerOrderNo, batchUnmarshalled.SequencerOrderNo)
	require.Equal(t, batchHeader.GasLimit, batchUnmarshalled.GasLimit)
	require.Equal(t, batchHeader.GasUsed, batchUnmarshalled.GasUsed)
	require.Equal(t, batchHeader.Time, batchUnmarshalled.Time)
	require.Equal(t, batchHeader.Extra, batchUnmarshalled.Extra)
	require.Equal(t, batchHeader.BaseFee, batchUnmarshalled.BaseFee)
	require.Equal(t, batchHeader.L1Proof, batchUnmarshalled.L1Proof)
	require.Equal(t, batchHeader.R, batchUnmarshalled.R)
	require.Equal(t, batchHeader.S, batchUnmarshalled.S)
	require.Equal(t, batchHeader.CrossChainMessages, batchUnmarshalled.CrossChainMessages)
	require.Equal(t, batchHeader.LatestInboundCrossChainHash, batchUnmarshalled.LatestInboundCrossChainHash)
	require.Equal(t, batchHeader.LatestInboundCrossChainHeight, batchUnmarshalled.LatestInboundCrossChainHeight)
	require.Equal(t, batchHeader.Hash(), batchUnmarshalled.Hash())
}

func randomHash() gethcommon.Hash {
	byteArr := make([]byte, 32)
	if _, err := rand.Read(byteArr); err != nil {
		panic(err)
	}

	return gethcommon.BytesToHash(byteArr)
}
