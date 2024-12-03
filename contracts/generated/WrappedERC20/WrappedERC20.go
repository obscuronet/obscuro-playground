// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package WrappedERC20

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// WrappedERC20MetaData contains all meta data concerning the WrappedERC20 contract.
var WrappedERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"giver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issueFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f317905534801561003657600080fd5b506040516114dc3803806114dc83398101604081905261005591610276565b81818181600361006583826103ce565b50600461007282826103ce565b50505050506100a77fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775336100af60201b60201c565b50505061048d565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff166101555760008381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561010d3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610159565b5060005b92915050565b634e487b7160e01b600052604160045260246000fd5b601f19601f83011681016001600160401b038111828210171561019a5761019a61015f565b6040525050565b60006101ac60405190565b90506101b88282610175565b919050565b60006001600160401b038211156101d6576101d661015f565b601f19601f83011660200192915050565b60005b838110156102025781810151838201526020016101ea565b50506000910152565b600061021e610219846101bd565b6101a1565b905082815283838301111561023557610235600080fd5b6102438360208301846101e7565b9392505050565b600082601f83011261025e5761025e600080fd5b815161026e84826020860161020b565b949350505050565b6000806040838503121561028c5761028c600080fd5b82516001600160401b038111156102a5576102a5600080fd5b6102b18582860161024a565b602085015190935090506001600160401b038111156102d2576102d2600080fd5b6102de8582860161024a565b9150509250929050565b634e487b7160e01b600052602260045260246000fd5b60028104600182168061031257607f821691505b602082108103610324576103246102e8565b50919050565b60006101596103368381565b90565b6103428361032a565b815460001960089490940293841b1916921b91909117905550565b600061036a818484610339565b505050565b8181101561038a5761038260008261035d565b60010161036f565b5050565b601f82111561036a576000818152602090206020601f850104810160208510156103b55750805b6103c76020601f86010483018261036f565b5050505050565b81516001600160401b038111156103e7576103e761015f565b6103f182546102fe565b6103fc82828561038e565b506020601f82116001811461043157600083156104195750848201515b600019600885021c19811660028502178555506103c7565b600084815260208120601f198516915b828110156104615787850151825560209485019460019092019101610441565b508482101561047e5783870151600019601f87166008021c191681555b50505050600202600101905550565b6110408061049c6000396000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c806336568abe116100cd578063979005ad11610081578063a9059cbb11610066578063a9059cbb146102ce578063d547741f146102e1578063dd62ed3e146102f457600080fd5b8063979005ad146102b3578063a217fddf146102c657600080fd5b806375b238fc116100b257806375b238fc1461024b57806391d148541461027257806395d89b41146102ab57600080fd5b806336568abe1461022557806370a082311461023857600080fd5b80631dd319cb11610124578063248a9ca311610109578063248a9ca3146101e05780632f2ff15d14610203578063313ce5671461021657600080fd5b80631dd319cb146101b857806323b872dd146101cd57600080fd5b806301ffc9a71461015657806306fdde031461017f578063095ea7b31461019457806318160ddd146101a7575b600080fd5b610169610164366004610c49565b610307565b6040516101769190610c72565b60405180910390f35b6101876103a0565b6040516101769190610cd6565b6101696101a2366004610d1d565b610432565b6002545b6040516101769190610d5b565b6101cb6101c6366004610d1d565b61044a565b005b6101696101db366004610d69565b6104b4565b6101ab6101ee366004610db2565b60009081526007602052604090206001015490565b6101cb610211366004610dd1565b6104da565b60126040516101769190610e09565b6101cb610233366004610dd1565b610505565b6101ab610246366004610e17565b610551565b6101ab7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b610169610280366004610dd1565b60009182526007602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6101876105c7565b6101cb6102c1366004610d1d565b6105d6565b6101ab600081565b6101696102dc366004610d1d565b61060a565b6101cb6102ef366004610dd1565b610618565b6101ab610302366004610e36565b61063d565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061039a57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546103af90610e6c565b80601f01602080910402602001604051908101604052809291908181526020018280546103db90610e6c565b80156104285780601f106103fd57610100808354040283529160200191610428565b820191906000526020600020905b81548152906001019060200180831161040b57829003601f168201915b5050505050905090565b6000336104408185856106f8565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561047481610705565b8161047e84610551565b10156104a55760405162461bcd60e51b815260040161049c90610ecc565b60405180910390fd5b6104af8383610712565b505050565b6000336104c285828561074c565b6104cd8585856107b2565b60019150505b9392505050565b6000828152600760205260409020600101546104f581610705565b6104ff8383610811565b50505050565b6001600160a01b0381163314610547576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104af82826108bf565b60006001600160a01b0382163203610581576001600160a01b03821660009081526020819052604090205461039a565b6001600160a01b03821633036105af576001600160a01b03821660009081526020819052604090205461039a565b60405162461bcd60e51b815260040161049c90610f0e565b6060600480546103af90610e6c565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561060081610705565b6104af8383610946565b6000336104408185856107b2565b60008281526007602052604090206001015461063381610705565b6104ff83836108bf565b6000326001600160a01b038416148061065e5750326001600160a01b038316145b15610691576001600160a01b038084166000908152600160209081526040808320938616835292905220545b905061039a565b336001600160a01b03841614806106b05750336001600160a01b038316145b156106e0576001600160a01b0380841660009081526001602090815260408083209386168352929052205461068a565b60405162461bcd60e51b815260040161049c90610f1e565b6104af838383600161097c565b61070f8133610a81565b50565b6001600160a01b03821661073c576000604051634b637e8f60e11b815260040161049c9190610f88565b61074882600083610ae1565b5050565b6000610758848461063d565b905060001981146104ff57818110156107a3578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161049c93929190610f96565b6104ff8484848403600061097c565b6001600160a01b0383166107dc576000604051634b637e8f60e11b815260040161049c9190610f88565b6001600160a01b03821661080657600060405163ec442f0560e01b815260040161049c9190610f88565b6104af838383610ae1565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff166108b75760008381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561086f3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161039a565b50600061039a565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff16156108b75760008381526007602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161039a565b6001600160a01b03821661097057600060405163ec442f0560e01b815260040161049c9190610f88565b61074860008383610ae1565b6001600160a01b0384166109bf5760006040517fe602df0500000000000000000000000000000000000000000000000000000000815260040161049c9190610f88565b6001600160a01b038316610a025760006040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161049c9190610f88565b6001600160a01b03808516600090815260016020908152604080832093871683529290522082905580156104ff57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610a739190610d5b565b60405180910390a350505050565b60008281526007602090815260408083206001600160a01b038516845290915290205460ff166107485780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161049c929190610fc6565b6001600160a01b038316610b0c578060026000828254610b019190610ff7565b90915550610b849050565b6001600160a01b03831660009081526020819052604090205481811015610b65578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161049c93929190610f96565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b038216610ba057600280548290039055610bbf565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610c029190610d5b565b60405180910390a3505050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b811461070f57600080fd5b803561039a81610c0f565b600060208284031215610c5e57610c5e600080fd5b6104d38383610c3e565b8015155b82525050565b6020810161039a8284610c68565b60005b83811015610c9b578181015183820152602001610c83565b50506000910152565b6000610cae825190565b808452602084019350610cc5818560208601610c80565b601f01601f19169290920192915050565b602080825281016104d38184610ca4565b60006001600160a01b03821661039a565b610c3381610ce7565b803561039a81610cf8565b80610c33565b803561039a81610d0c565b60008060408385031215610d3357610d33600080fd5b610d3d8484610d01565b9150610d4c8460208501610d12565b90509250929050565b80610c6c565b6020810161039a8284610d55565b600080600060608486031215610d8157610d81600080fd5b610d8b8585610d01565b9250610d9a8560208601610d01565b9150610da98560408601610d12565b90509250925092565b600060208284031215610dc757610dc7600080fd5b6104d38383610d12565b60008060408385031215610de757610de7600080fd5b610df18484610d12565b9150610d4c8460208501610d01565b60ff8116610c6c565b6020810161039a8284610e00565b600060208284031215610e2c57610e2c600080fd5b6104d38383610d01565b60008060408385031215610e4c57610e4c600080fd5b610df18484610d01565b634e487b7160e01b600052602260045260246000fd5b600281046001821680610e8057607f821691505b602082108103610e9257610e92610e56565b50919050565b60158152602081017f496e73756666696369656e742062616c616e63652e0000000000000000000000815290505b60200190565b6020808252810161039a81610e98565b601f8152602081017f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e63650081529050610ec6565b6020808252810161039a81610edc565b6020808252810161039a81602181527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360208201527f6500000000000000000000000000000000000000000000000000000000000000604082015260600190565b610c6c81610ce7565b6020810161039a8284610f7f565b60608101610fa48286610f7f565b610fb16020830185610d55565b610fbe6040830184610d55565b949350505050565b60408101610fd48285610f7f565b6104d36020830184610d55565b634e487b7160e01b600052601160045260246000fd5b8082018082111561039a5761039a610fe156fea2646970667358221220be3380d12a8a9c04f488bd6d863b29f9c1ea464ebf1b28b5c552d2e3f79e9fda64736f6c634300081c0033",
}

// WrappedERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use WrappedERC20MetaData.ABI instead.
var WrappedERC20ABI = WrappedERC20MetaData.ABI

// WrappedERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WrappedERC20MetaData.Bin instead.
var WrappedERC20Bin = WrappedERC20MetaData.Bin

// DeployWrappedERC20 deploys a new Ethereum contract, binding an instance of WrappedERC20 to it.
func DeployWrappedERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *WrappedERC20, error) {
	parsed, err := WrappedERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WrappedERC20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WrappedERC20{WrappedERC20Caller: WrappedERC20Caller{contract: contract}, WrappedERC20Transactor: WrappedERC20Transactor{contract: contract}, WrappedERC20Filterer: WrappedERC20Filterer{contract: contract}}, nil
}

// WrappedERC20 is an auto generated Go binding around an Ethereum contract.
type WrappedERC20 struct {
	WrappedERC20Caller     // Read-only binding to the contract
	WrappedERC20Transactor // Write-only binding to the contract
	WrappedERC20Filterer   // Log filterer for contract events
}

// WrappedERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type WrappedERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type WrappedERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrappedERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrappedERC20Session struct {
	Contract     *WrappedERC20     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrappedERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrappedERC20CallerSession struct {
	Contract *WrappedERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// WrappedERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrappedERC20TransactorSession struct {
	Contract     *WrappedERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// WrappedERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type WrappedERC20Raw struct {
	Contract *WrappedERC20 // Generic contract binding to access the raw methods on
}

// WrappedERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrappedERC20CallerRaw struct {
	Contract *WrappedERC20Caller // Generic read-only contract binding to access the raw methods on
}

// WrappedERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrappedERC20TransactorRaw struct {
	Contract *WrappedERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewWrappedERC20 creates a new instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20(address common.Address, backend bind.ContractBackend) (*WrappedERC20, error) {
	contract, err := bindWrappedERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20{WrappedERC20Caller: WrappedERC20Caller{contract: contract}, WrappedERC20Transactor: WrappedERC20Transactor{contract: contract}, WrappedERC20Filterer: WrappedERC20Filterer{contract: contract}}, nil
}

// NewWrappedERC20Caller creates a new read-only instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Caller(address common.Address, caller bind.ContractCaller) (*WrappedERC20Caller, error) {
	contract, err := bindWrappedERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Caller{contract: contract}, nil
}

// NewWrappedERC20Transactor creates a new write-only instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*WrappedERC20Transactor, error) {
	contract, err := bindWrappedERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Transactor{contract: contract}, nil
}

// NewWrappedERC20Filterer creates a new log filterer instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*WrappedERC20Filterer, error) {
	contract, err := bindWrappedERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Filterer{contract: contract}, nil
}

// bindWrappedERC20 binds a generic wrapper to an already deployed contract.
func bindWrappedERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WrappedERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedERC20 *WrappedERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedERC20.Contract.WrappedERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedERC20 *WrappedERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedERC20.Contract.WrappedERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedERC20 *WrappedERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedERC20.Contract.WrappedERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedERC20 *WrappedERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedERC20 *WrappedERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedERC20 *WrappedERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedERC20.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) ADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.ADMINROLE(&_WrappedERC20.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) ADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.ADMINROLE(&_WrappedERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.DEFAULTADMINROLE(&_WrappedERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.DEFAULTADMINROLE(&_WrappedERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.Allowance(&_WrappedERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.Allowance(&_WrappedERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.BalanceOf(&_WrappedERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.BalanceOf(&_WrappedERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20Session) Decimals() (uint8, error) {
	return _WrappedERC20.Contract.Decimals(&_WrappedERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20CallerSession) Decimals() (uint8, error) {
	return _WrappedERC20.Contract.Decimals(&_WrappedERC20.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedERC20.Contract.GetRoleAdmin(&_WrappedERC20.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedERC20.Contract.GetRoleAdmin(&_WrappedERC20.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedERC20.Contract.HasRole(&_WrappedERC20.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedERC20.Contract.HasRole(&_WrappedERC20.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20Session) Name() (string, error) {
	return _WrappedERC20.Contract.Name(&_WrappedERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20CallerSession) Name() (string, error) {
	return _WrappedERC20.Contract.Name(&_WrappedERC20.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedERC20.Contract.SupportsInterface(&_WrappedERC20.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedERC20.Contract.SupportsInterface(&_WrappedERC20.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20Session) Symbol() (string, error) {
	return _WrappedERC20.Contract.Symbol(&_WrappedERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20CallerSession) Symbol() (string, error) {
	return _WrappedERC20.Contract.Symbol(&_WrappedERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) TotalSupply() (*big.Int, error) {
	return _WrappedERC20.Contract.TotalSupply(&_WrappedERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _WrappedERC20.Contract.TotalSupply(&_WrappedERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Approve(&_WrappedERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Approve(&_WrappedERC20.TransactOpts, spender, value)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Transactor) BurnFor(opts *bind.TransactOpts, giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "burnFor", giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Session) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.BurnFor(&_WrappedERC20.TransactOpts, giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.BurnFor(&_WrappedERC20.TransactOpts, giver, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.GrantRole(&_WrappedERC20.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.GrantRole(&_WrappedERC20.TransactOpts, role, account)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Transactor) IssueFor(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "issueFor", receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Session) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.IssueFor(&_WrappedERC20.TransactOpts, receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.IssueFor(&_WrappedERC20.TransactOpts, receiver, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RenounceRole(&_WrappedERC20.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RenounceRole(&_WrappedERC20.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RevokeRole(&_WrappedERC20.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RevokeRole(&_WrappedERC20.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Transfer(&_WrappedERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Transfer(&_WrappedERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.TransferFrom(&_WrappedERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.TransferFrom(&_WrappedERC20.TransactOpts, from, to, value)
}

// WrappedERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WrappedERC20 contract.
type WrappedERC20ApprovalIterator struct {
	Event *WrappedERC20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedERC20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20Approval represents a Approval event raised by the WrappedERC20 contract.
type WrappedERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WrappedERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20ApprovalIterator{contract: _WrappedERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WrappedERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20Approval)
				if err := _WrappedERC20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) ParseApproval(log types.Log) (*WrappedERC20Approval, error) {
	event := new(WrappedERC20Approval)
	if err := _WrappedERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the WrappedERC20 contract.
type WrappedERC20RoleAdminChangedIterator struct {
	Event *WrappedERC20RoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedERC20RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedERC20RoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedERC20RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleAdminChanged represents a RoleAdminChanged event raised by the WrappedERC20 contract.
type WrappedERC20RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*WrappedERC20RoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleAdminChangedIterator{contract: _WrappedERC20.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleAdminChanged)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleAdminChanged(log types.Log) (*WrappedERC20RoleAdminChanged, error) {
	event := new(WrappedERC20RoleAdminChanged)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the WrappedERC20 contract.
type WrappedERC20RoleGrantedIterator struct {
	Event *WrappedERC20RoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedERC20RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedERC20RoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedERC20RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleGranted represents a RoleGranted event raised by the WrappedERC20 contract.
type WrappedERC20RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedERC20RoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleGrantedIterator{contract: _WrappedERC20.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleGranted)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleGranted(log types.Log) (*WrappedERC20RoleGranted, error) {
	event := new(WrappedERC20RoleGranted)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the WrappedERC20 contract.
type WrappedERC20RoleRevokedIterator struct {
	Event *WrappedERC20RoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedERC20RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedERC20RoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedERC20RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleRevoked represents a RoleRevoked event raised by the WrappedERC20 contract.
type WrappedERC20RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedERC20RoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleRevokedIterator{contract: _WrappedERC20.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleRevoked)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleRevoked(log types.Log) (*WrappedERC20RoleRevoked, error) {
	event := new(WrappedERC20RoleRevoked)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WrappedERC20 contract.
type WrappedERC20TransferIterator struct {
	Event *WrappedERC20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WrappedERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WrappedERC20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WrappedERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20Transfer represents a Transfer event raised by the WrappedERC20 contract.
type WrappedERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20TransferIterator{contract: _WrappedERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WrappedERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20Transfer)
				if err := _WrappedERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) ParseTransfer(log types.Log) (*WrappedERC20Transfer, error) {
	event := new(WrappedERC20Transfer)
	if err := _WrappedERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
