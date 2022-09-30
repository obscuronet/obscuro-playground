package evm

import (
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// ExecuteTransactions
// header - the header of the rollup where this transaction will be included
// fromTxIndex - for the receipts and events, the evm needs to know for each transaction the order in which it was executed in the block.
func ExecuteTransactions(txs []*common.L2Tx, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig, fromTxIndex int) map[common.TxHash]interface{} {
	chain, vmCfg, gp := initParams(storage, false)
	zero := uint64(0)
	usedGas := &zero
	result := map[common.TxHash]interface{}{}

	ethHeader := convertToEthHeader(header, secret(storage))

	for i, t := range txs {
		r, err := executeTransaction(s, chainConfig, chain, gp, ethHeader, t, usedGas, vmCfg, fromTxIndex+i)
		if err != nil {
			result[t.Hash()] = err
			common.ErrorTXExecution(t.Hash(), "Error: %s", err)
			continue
		}
		result[t.Hash()] = r
		logReceipt(r)
	}
	s.Finalise(true)
	return result
}

func executeTransaction(s *state.StateDB, cc *params.ChainConfig, chain *ObscuroChainContext, gp *gethcore.GasPool, header *types.Header, t *common.L2Tx, usedGas *uint64, vmCfg vm.Config, tCount int) (*types.Receipt, error) {
	s.Prepare(t.Hash(), tCount)
	snap := s.Snapshot()

	before := header.MixDigest
	// calculate a random value per transaction
	header.MixDigest = gethcommon.BytesToHash(crypto.PerTransactionRnd(before.Bytes(), tCount))
	// todo - Author?
	receipt, err := gethcore.ApplyTransaction(cc, chain, nil, gp, s, header, t, usedGas, vmCfg)
	header.MixDigest = before
	if err != nil {
		s.RevertToSnapshot(snap)
		return nil, err
	}

	return receipt, nil
}

func logReceipt(r *types.Receipt) {
	receiptJSON, err := r.MarshalJSON()
	if err != nil {
		if r.Status == types.ReceiptStatusFailed {
			common.ErrorTXExecution(r.TxHash, "Unsuccessful (status != 1) (but could not print receipt as JSON)")
		} else {
			common.TraceTXExecution(r.TxHash, "Successfully executed (but could not print receipt as JSON)")
		}
	}

	if r.Status == types.ReceiptStatusFailed {
		common.ErrorTXExecution(r.TxHash, "Unsuccessful (status != 1). Receipt: %s", string(receiptJSON))
	} else {
		common.TraceTXExecution(r.TxHash, "Successfully executed. Receipt: %s", string(receiptJSON))
	}
}

// ExecuteOffChainCall - executes the "data" command against the "to" smart contract
func ExecuteOffChainCall(from gethcommon.Address, to *gethcommon.Address, data []byte, s *state.StateDB, header *common.Header, storage db.Storage, chainConfig *params.ChainConfig) (*gethcore.ExecutionResult, error) {
	chain, vmCfg, gp := initParams(storage, true)

	blockContext := gethcore.NewEVMBlockContext(convertToEthHeader(header, secret(storage)), chain, &header.Agg)
	// todo use ToMessage
	// 100_000_000_000 is just a huge number gasLimit for making sure the local tx doesn't fail with lack of gas
	msg := types.NewMessage(from, to, 0, gethcommon.Big0, 100_000_000_000, gethcommon.Big0, gethcommon.Big0, gethcommon.Big0, data, nil, true)

	// sets Tx.origin
	txContext := gethcore.NewEVMTxContext(msg)
	vmenv := vm.NewEVM(blockContext, txContext, s, chainConfig, vmCfg)

	result, err := gethcore.ApplyMessage(vmenv, msg, gp)
	if err != nil {
		// this error is ignored by geth. logging just in case
		log.Error("Error applying msg: %s", err)
	}
	err = s.Error()
	if err != nil {
		return nil, err
	}

	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return nil, newRevertError(result)
	}
	return result, nil
}

func initParams(storage db.Storage, noBaseFee bool) (*ObscuroChainContext, vm.Config, *gethcore.GasPool) {
	chain := &ObscuroChainContext{storage: storage}
	tracer := logger.NewStructLogger(&logger.Config{Debug: true})
	vmCfg := vm.Config{
		NoBaseFee: noBaseFee,
		Debug:     false,
		Tracer:    tracer,
	}
	gp := gethcore.GasPool(math.MaxUint64)
	return chain, vmCfg, &gp
}

// Todo - this is currently just returning the shared secret
// it should not use it directly, but derive some entropy from it
func secret(storage db.Storage) []byte {
	secret := storage.FetchSecret()
	return secret[:]
}

// copy pasted from geth
func newRevertError(result *gethcore.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}
