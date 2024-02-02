package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func ExtractGetBalanceRequestreqParams(reqParams []any, rpc *EncryptionManager) (*UserRPCRequest1[hexutil.Big], error) {
	// Parameters are [Address, BlockNumber]
	if len(reqParams) != 2 {
		return nil, fmt.Errorf("unexpected number of parameters")
	}
	requestedAddress, err := gethencoding.ExtractAddress(reqParams[0])
	if err != nil {
		return nil, fmt.Errorf("unable to extract requested address - %w", err)
	}

	blockNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
	if err != nil {
		return nil, fmt.Errorf("unable to extract requested block number - %w", err)
	}

	encryptAddress, balance, err := rpc.chain.GetBalance(*requestedAddress, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get balance - %w", err)
	}

	return &UserRPCRequest1[hexutil.Big]{encryptAddress, balance}, nil
}

func ExecuteGetBalance(decodedParams *UserRPCRequest1[hexutil.Big], _ *EncryptionManager) (*UserResponse[hexutil.Big], error) {
	return &UserResponse[hexutil.Big]{decodedParams.Param1, nil}, nil
}
