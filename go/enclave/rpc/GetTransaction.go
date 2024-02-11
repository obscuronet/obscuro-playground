package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

func GetTransactionValidate(reqParams []any, builder *CallBuilder[gethcommon.Hash, RpcTransaction], _ *EncryptionManager) error {
	// Parameters are [Hash]
	if len(reqParams) != 1 {
		builder.Err = fmt.Errorf("wrong parameters")
		return nil
	}
	txHashStr, ok := reqParams[0].(string)
	if !ok {
		builder.Err = fmt.Errorf("unexpected tx hash parameter")
		return nil
	}
	txHash := gethcommon.HexToHash(txHashStr)
	builder.Param = &txHash
	return nil
}

func GetTransactionExecute(builder *CallBuilder[gethcommon.Hash, RpcTransaction], rpc *EncryptionManager) error {
	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := rpc.storage.GetTransaction(*builder.Param)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		return err
	}

	sender, err := core.GetTxSigner(tx)
	if err != nil {
		return fmt.Errorf("could not recover the tx %s sender. Cause: %w", tx.Hash(), err)
	}

	// authorise - only the signer can request the transaction
	if sender.Hex() != builder.VK.AccountAddress.Hex() {
		builder.Status = NotAuthorised
		// builder.ReturnValue= []byte{}
		return nil
	}

	// Unlike in the Geth impl, we hardcode the use of a London signer.
	// todo (#1553) - once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`
	signer := types.NewLondonSigner(tx.ChainId())
	rpcTx := newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer)
	builder.ReturnValue = rpcTx
	return nil
}

// Lifted from Geth's internal `ethapi` package.
type RpcTransaction struct { //nolint
	BlockHash        *gethcommon.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big        `json:"blockNumber"`
	From             gethcommon.Address  `json:"from"`
	Gas              hexutil.Uint64      `json:"gas"`
	GasPrice         *hexutil.Big        `json:"gasPrice"`
	GasFeeCap        *hexutil.Big        `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big        `json:"maxPriorityFeePerGas,omitempty"`
	Hash             gethcommon.Hash     `json:"hash"`
	Input            hexutil.Bytes       `json:"input"`
	Nonce            hexutil.Uint64      `json:"nonce"`
	To               *gethcommon.Address `json:"to"`
	TransactionIndex *hexutil.Uint64     `json:"transactionIndex"`
	Value            *hexutil.Big        `json:"value"`
	Type             hexutil.Uint64      `json:"type"`
	Accesses         *types.AccessList   `json:"accessList,omitempty"`
	ChainID          *hexutil.Big        `json:"chainId,omitempty"`
	V                *hexutil.Big        `json:"v"`
	R                *hexutil.Big        `json:"r"`
	S                *hexutil.Big        `json:"s"`
}

// Lifted from Geth's internal `ethapi` package.
func newRPCTransaction(tx *types.Transaction, blockHash gethcommon.Hash, blockNumber uint64, index uint64, baseFee *big.Int, signer types.Signer) *RpcTransaction {
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RpcTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (gethcommon.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case types.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (gethcommon.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}
	return result
}
