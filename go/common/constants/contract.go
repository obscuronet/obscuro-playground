package constants

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"
)

// TODO move this out of the constants package
func Bytecode(seqAddress common.Address) ([]byte, error) {
	parsed, err := ManagementContract.ManagementContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("", seqAddress)
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(ManagementContract.ManagementContractMetaData.Bin)
	return append(bytecode, input...), nil
}
