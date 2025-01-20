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
	Bin: "0x6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f3179055348015610035575f5ffd5b5060405161144738038061144783398101604081905261005491610251565b818181816003610064838261039e565b506004610071828261039e565b50505050506100a67fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775336100ae60201b60201c565b505050610459565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff16610152575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561010a3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610155565b505f5b92915050565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b03811182821017156101945761019461015b565b6040525050565b5f6101a560405190565b90506101b1828261016f565b919050565b5f6001600160401b038211156101ce576101ce61015b565b601f19601f83011660200192915050565b8281835e505f910152565b5f6101fc6101f7846101b6565b61019b565b9050828152838383011115610212576102125f5ffd5b6102208360208301846101df565b9392505050565b5f82601f830112610239576102395f5ffd5b81516102498482602086016101ea565b949350505050565b5f5f60408385031215610265576102655f5ffd5b82516001600160401b0381111561027d5761027d5f5ffd5b61028985828601610227565b602085015190935090506001600160401b038111156102a9576102a95f5ffd5b6102b585828601610227565b9150509250929050565b634e487b7160e01b5f52602260045260245ffd5b6002810460018216806102e757607f821691505b6020821081036102f9576102f96102bf565b50919050565b5f61015561030a8381565b90565b610316836102ff565b81545f1960089490940293841b1916921b91909117905550565b5f61033c81848461030d565b505050565b8181101561035b576103535f82610330565b600101610341565b5050565b601f82111561033c575f818152602090206020601f850104810160208510156103855750805b6103976020601f860104830182610341565b5050505050565b81516001600160401b038111156103b7576103b761015b565b6103c182546102d3565b6103cc82828561035f565b506020601f8211600181146103ff575f83156103e85750848201515b5f19600885021c1981166002850217855550610397565b5f84815260208120601f198516915b8281101561042e578785015182556020948501946001909201910161040e565b508482101561044a57838701515f19601f87166008021c191681555b50505050600202600101905550565b610fe1806104665f395ff3fe608060405234801561000f575f5ffd5b5060043610610149575f3560e01c806336568abe116100c7578063979005ad1161007d578063a9059cbb11610063578063a9059cbb146102c2578063d547741f146102d5578063dd62ed3e146102e8575f5ffd5b8063979005ad146102a8578063a217fddf146102bb575f5ffd5b806375b238fc116100ad57806375b238fc1461024157806391d148541461026857806395d89b41146102a0575f5ffd5b806336568abe1461021b57806370a082311461022e575f5ffd5b80631dd319cb1161011c578063248a9ca311610102578063248a9ca3146101d75780632f2ff15d146101f9578063313ce5671461020c575f5ffd5b80631dd319cb146101af57806323b872dd146101c4575f5ffd5b806301ffc9a71461014d57806306fdde0314610176578063095ea7b31461018b57806318160ddd1461019e575b5f5ffd5b61016061015b366004610c18565b6102fb565b60405161016d9190610c3f565b60405180910390f35b61017e610393565b60405161016d9190610c89565b610160610199366004610ccf565b610423565b6002545b60405161016d9190610d0b565b6101c26101bd366004610ccf565b61043a565b005b6101606101d2366004610d19565b6104a4565b6101a26101e5366004610d5f565b5f9081526007602052604090206001015490565b6101c2610207366004610d7c565b6104c9565b601260405161016d9190610db2565b6101c2610229366004610d7c565b6104f3565b6101a261023c366004610dc0565b61053f565b6101a27fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b610160610276366004610d7c565b5f9182526007602090815260408084206001600160a01b0393909316845291905290205460ff1690565b61017e6105b2565b6101c26102b6366004610ccf565b6105c1565b6101a25f81565b6101606102d0366004610ccf565b6105f5565b6101c26102e3366004610d7c565b610602565b6101a26102f6366004610ddd565b610626565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061038d57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546103a290610e0f565b80601f01602080910402602001604051908101604052809291908181526020018280546103ce90610e0f565b80156104195780601f106103f057610100808354040283529160200191610419565b820191905f5260205f20905b8154815290600101906020018083116103fc57829003601f168201915b5050505050905090565b5f336104308185856106de565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610464816106eb565b8161046e8461053f565b10156104955760405162461bcd60e51b815260040161048c90610e6f565b60405180910390fd5b61049f83836106f8565b505050565b5f336104b1858285610730565b6104bc858585610793565b60019150505b9392505050565b5f828152600760205260409020600101546104e3816106eb565b6104ed83836107f0565b50505050565b6001600160a01b0381163314610535576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61049f828261089b565b5f6001600160a01b038216320361056d576001600160a01b0382165f9081526020819052604090205461038d565b6001600160a01b038216330361059a576001600160a01b0382165f9081526020819052604090205461038d565b60405162461bcd60e51b815260040161048c90610eb1565b6060600480546103a290610e0f565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756105eb816106eb565b61049f8383610920565b5f33610430818585610793565b5f8281526007602052604090206001015461061c816106eb565b6104ed838361089b565b5f326001600160a01b03841614806106465750326001600160a01b038316145b15610678576001600160a01b038084165f908152600160209081526040808320938616835292905220545b905061038d565b336001600160a01b03841614806106975750336001600160a01b038316145b156106c6576001600160a01b038084165f90815260016020908152604080832093861683529290522054610671565b60405162461bcd60e51b815260040161048c90610ec1565b61049f8383836001610954565b6106f58133610a56565b50565b6001600160a01b038216610721575f604051634b637e8f60e11b815260040161048c9190610f2b565b61072c825f83610ab5565b5050565b5f61073b8484610626565b90505f1981146104ed5781811015610785578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161048c93929190610f39565b6104ed84848484035f610954565b6001600160a01b0383166107bc575f604051634b637e8f60e11b815260040161048c9190610f2b565b6001600160a01b0382166107e5575f60405163ec442f0560e01b815260040161048c9190610f2b565b61049f838383610ab5565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff16610894575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561084c3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161038d565b505f61038d565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff1615610894575f8381526007602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161038d565b6001600160a01b038216610949575f60405163ec442f0560e01b815260040161048c9190610f2b565b61072c5f8383610ab5565b6001600160a01b038416610996575f6040517fe602df0500000000000000000000000000000000000000000000000000000000815260040161048c9190610f2b565b6001600160a01b0383166109d8575f6040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161048c9190610f2b565b6001600160a01b038085165f90815260016020908152604080832093871683529290522082905580156104ed57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610a489190610d0b565b60405180910390a350505050565b5f8281526007602090815260408083206001600160a01b038516845290915290205460ff1661072c5780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161048c929190610f69565b6001600160a01b038316610adf578060025f828254610ad49190610f98565b90915550610b559050565b6001600160a01b0383165f9081526020819052604090205481811015610b37578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161048c93929190610f39565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b038216610b7157600280548290039055610b8f565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610bd29190610d0b565b60405180910390a3505050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b81146106f5575f5ffd5b803561038d81610bdf565b5f60208284031215610c2b57610c2b5f5ffd5b6104c28383610c0d565b8015155b82525050565b6020810161038d8284610c35565b8281835e505f910152565b5f610c61825190565b808452602084019350610c78818560208601610c4d565b601f01601f19169290920192915050565b602080825281016104c28184610c58565b5f6001600160a01b03821661038d565b610c0381610c9a565b803561038d81610caa565b80610c03565b803561038d81610cbe565b5f5f60408385031215610ce357610ce35f5ffd5b610ced8484610cb3565b9150610cfc8460208501610cc4565b90509250929050565b80610c39565b6020810161038d8284610d05565b5f5f5f60608486031215610d2e57610d2e5f5ffd5b610d388585610cb3565b9250610d478560208601610cb3565b9150610d568560408601610cc4565b90509250925092565b5f60208284031215610d7257610d725f5ffd5b6104c28383610cc4565b5f5f60408385031215610d9057610d905f5ffd5b610d9a8484610cc4565b9150610cfc8460208501610cb3565b60ff8116610c39565b6020810161038d8284610da9565b5f60208284031215610dd357610dd35f5ffd5b6104c28383610cb3565b5f5f60408385031215610df157610df15f5ffd5b610d9a8484610cb3565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610e2357607f821691505b602082108103610e3557610e35610dfb565b50919050565b60158152602081017f496e73756666696369656e742062616c616e63652e0000000000000000000000815290505b60200190565b6020808252810161038d81610e3b565b601f8152602081017f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e63650081529050610e69565b6020808252810161038d81610e7f565b6020808252810161038d81602181527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360208201527f6500000000000000000000000000000000000000000000000000000000000000604082015260600190565b610c3981610c9a565b6020810161038d8284610f22565b60608101610f478286610f22565b610f546020830185610d05565b610f616040830184610d05565b949350505050565b60408101610f778285610f22565b6104c26020830184610d05565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561038d5761038d610f8456fea2646970667358221220ba6b4996116e560a54ab485db6a674eeef7255da269323f7a3bdf4b4aa8995c964736f6c634300081c0033",
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
