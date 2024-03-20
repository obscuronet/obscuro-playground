package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/syserr"
)

func TenCallValidate(reqParams []any, builder *CallBuilder[CallParamsWithBlock, string], _ *EncryptionManager) error {
	// Parameters are [TransactionArgs, BlockNumber, 2 more which we don't support yet]
	if len(reqParams) < 2 && len(reqParams) > 4 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	apiArgs, err := gethencoding.ExtractEthCall(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return nil
	}

	if apiArgs.From == nil {
		builder.Err = fmt.Errorf("no from address provided")
		return nil
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}

	builder.From = apiArgs.From
	// todo - support BlockNumberOrHash
	builder.Param = &CallParamsWithBlock{apiArgs, blkNumber.BlockNumber}

	return nil
}

func TenCallExecute(builder *CallBuilder[CallParamsWithBlock, string], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	apiArgs := builder.Param.callParams
	blkNumber := builder.Param.block
	execResult, err := rpc.chain.ObsCall(apiArgs, blkNumber)
	if err != nil {
		rpc.logger.Debug("Failed eth_call.", log.ErrKey, err)

		// return system errors to the host
		if errors.Is(err, syserr.InternalError{}) {
			return err
		}

		// extract the EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		builder.Err = err
		return nil
	}

	var encodedResult string
	if len(execResult.ReturnData) != 0 {
		encodedResult = hexutil.Encode(execResult.ReturnData)
		builder.ReturnValue = &encodedResult
	} else {
		builder.ReturnValue = nil
	}
	return nil
}

func serializeEVMError(err error) ([]byte, error) {
	var errReturn interface{}

	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *errutil.EVMSerialisableError
	if ok := errors.As(err, &e); ok {
		errReturn = e
	} else {
		// it's a generic error, serialise it
		errReturn = &errutil.EVMSerialisableError{Err: err.Error()}
	}

	// serialise the error object returned by the evm into a json
	errSerializedBytes, marshallErr := json.Marshal(errReturn)
	if marshallErr != nil {
		return nil, marshallErr
	}
	return errSerializedBytes, nil
}
