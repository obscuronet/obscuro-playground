// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConstantSupplyERC20

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

// ConstantSupplyERC20MetaData contains all meta data concerning the ConstantSupplyERC20 contract.
var ConstantSupplyERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051610db7380380610db783398101604081905261002e916102b5565b8282600361003c838261040e565b506004610049828261040e565b50505061005c338261006460201b60201c565b505050610559565b6001600160a01b038216610096575f60405163ec442f0560e01b815260040161008d91906104e8565b60405180910390fd5b6100a15f83836100a5565b5050565b6001600160a01b0383166100cf578060025f8282546100c4919061050a565b9091555061012c9050565b6001600160a01b0383165f908152602081905260409020548181101561010e5783818360405163391434e360e21b815260040161008d93929190610523565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b03821661014857600280548290039055610166565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516101a9919061054b565b60405180910390a3505050565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b03811182821017156101ef576101ef6101b6565b6040525050565b5f61020060405190565b905061020c82826101ca565b919050565b5f6001600160401b03821115610229576102296101b6565b601f19601f83011660200192915050565b8281835e505f910152565b5f61025761025284610211565b6101f6565b905082815283838301111561026d5761026d5f5ffd5b61027b83602083018461023a565b9392505050565b5f82601f830112610294576102945f5ffd5b81516102a4848260208601610245565b949350505050565b80515b92915050565b5f5f5f606084860312156102ca576102ca5f5ffd5b83516001600160401b038111156102e2576102e25f5ffd5b6102ee86828701610282565b602086015190945090506001600160401b0381111561030e5761030e5f5ffd5b61031a86828701610282565b92505061032a85604086016102ac565b90509250925092565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061035b57607f821691505b60208210810361036d5761036d610333565b50919050565b5f6102af61037e8381565b90565b61038a83610373565b81545f1960089490940293841b1916921b91909117905550565b5f6103b0818484610381565b505050565b818110156100a1576103c75f826103a4565b6001016103b5565b601f8211156103b0575f818152602090206020601f850104810160208510156103f55750805b6104076020601f8601048301826103b5565b5050505050565b81516001600160401b03811115610427576104276101b6565b6104318254610347565b61043c8282856103cf565b506020601f82116001811461046f575f83156104585750848201515b5f19600885021c1981166002850217855550610407565b5f84815260208120601f198516915b8281101561049e578785015182556020948501946001909201910161047e565b50848210156104ba57838701515f19601f87166008021c191681555b50505050600202600101905550565b5f6001600160a01b0382166102af565b6104e2816104c9565b82525050565b602081016102af82846104d9565b634e487b7160e01b5f52601160045260245ffd5b808201808211156102af576102af6104f6565b806104e2565b6060810161053182866104d9565b61053e602083018561051d565b6102a4604083018461051d565b602081016102af828461051d565b610851806105665f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c8063313ce5671161007257806395d89b411161005857806395d89b411461013c578063a9059cbb14610144578063dd62ed3e14610157575f5ffd5b8063313ce5671461010557806370a0823114610114575f5ffd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e157806323b872dd146100f2575b5f5ffd5b6100ab61018f565b6040516100b8919061060a565b60405180910390f35b6100d46100cf366004610664565b61021f565b6040516100b891906106a4565b6002545b6040516100b891906106b8565b6100d46101003660046106c6565b610238565b60126040516100b89190610715565b6100e5610122366004610723565b6001600160a01b03165f9081526020819052604090205490565b6100ab61025b565b6100d4610152366004610664565b61026a565b6100e5610165366004610740565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b60606003805461019e90610781565b80601f01602080910402602001604051908101604052809291908181526020018280546101ca90610781565b80156102155780601f106101ec57610100808354040283529160200191610215565b820191905f5260205f20905b8154815290600101906020018083116101f857829003601f168201915b5050505050905090565b5f3361022c818585610277565b60019150505b92915050565b5f33610245858285610289565b610250858585610313565b506001949350505050565b60606004805461019e90610781565b5f3361022c818585610313565b61028483838360016103a2565b505050565b6001600160a01b038381165f908152600160209081526040808320938616835292905220545f19811461030d57818110156102ff578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016102f6939291906107b6565b60405180910390fd5b61030d84848484035f6103a2565b50505050565b6001600160a01b038316610355575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016102f691906107e6565b6001600160a01b038216610397575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016102f691906107e6565b6102848383836104a4565b6001600160a01b0384166103e4575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016102f691906107e6565b6001600160a01b038316610426575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016102f691906107e6565b6001600160a01b038085165f908152600160209081526040808320938716835292905220829055801561030d57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161049691906106b8565b60405180910390a350505050565b6001600160a01b0383166104ce578060025f8282546104c39190610808565b909155506105449050565b6001600160a01b0383165f9081526020819052604090205481811015610526578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016102f6939291906107b6565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b0382166105605760028054829003905561057e565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105c191906106b8565b60405180910390a3505050565b8281835e505f910152565b5f6105e2825190565b8084526020840193506105f98185602086016105ce565b601f01601f19169290920192915050565b6020808252810161061b81846105d9565b9392505050565b5f6001600160a01b038216610232565b61063b81610622565b8114610645575f5ffd5b50565b803561023281610632565b8061063b565b803561023281610653565b5f5f60408385031215610678576106785f5ffd5b6106828484610648565b91506106918460208501610659565b90509250929050565b8015155b82525050565b60208101610232828461069a565b8061069e565b6020810161023282846106b2565b5f5f5f606084860312156106db576106db5f5ffd5b6106e58585610648565b92506106f48560208601610648565b91506107038560408601610659565b90509250925092565b60ff811661069e565b60208101610232828461070c565b5f60208284031215610736576107365f5ffd5b61061b8383610648565b5f5f60408385031215610754576107545f5ffd5b61075e8484610648565b91506106918460208501610648565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061079557607f821691505b6020821081036107a7576107a761076d565b50919050565b61069e81610622565b606081016107c482866107ad565b6107d160208301856106b2565b6107de60408301846106b2565b949350505050565b6020810161023282846107ad565b634e487b7160e01b5f52601160045260245ffd5b80820180821115610232576102326107f456fea26469706673582212202625812a0e573f7616e7b6dec452cbafd36cca84132156ab03c28e0e9e707efb64736f6c634300081c0033",
}

// ConstantSupplyERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ConstantSupplyERC20MetaData.ABI instead.
var ConstantSupplyERC20ABI = ConstantSupplyERC20MetaData.ABI

// ConstantSupplyERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConstantSupplyERC20MetaData.Bin instead.
var ConstantSupplyERC20Bin = ConstantSupplyERC20MetaData.Bin

// DeployConstantSupplyERC20 deploys a new Ethereum contract, binding an instance of ConstantSupplyERC20 to it.
func DeployConstantSupplyERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int) (common.Address, *types.Transaction, *ConstantSupplyERC20, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConstantSupplyERC20Bin), backend, name, symbol, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// ConstantSupplyERC20 is an auto generated Go binding around an Ethereum contract.
type ConstantSupplyERC20 struct {
	ConstantSupplyERC20Caller     // Read-only binding to the contract
	ConstantSupplyERC20Transactor // Write-only binding to the contract
	ConstantSupplyERC20Filterer   // Log filterer for contract events
}

// ConstantSupplyERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstantSupplyERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstantSupplyERC20Session struct {
	Contract     *ConstantSupplyERC20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstantSupplyERC20CallerSession struct {
	Contract *ConstantSupplyERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ConstantSupplyERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstantSupplyERC20TransactorSession struct {
	Contract     *ConstantSupplyERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ConstantSupplyERC20Raw struct {
	Contract *ConstantSupplyERC20 // Generic contract binding to access the raw methods on
}

// ConstantSupplyERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20CallerRaw struct {
	Contract *ConstantSupplyERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ConstantSupplyERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20TransactorRaw struct {
	Contract *ConstantSupplyERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewConstantSupplyERC20 creates a new instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20(address common.Address, backend bind.ContractBackend) (*ConstantSupplyERC20, error) {
	contract, err := bindConstantSupplyERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// NewConstantSupplyERC20Caller creates a new read-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Caller(address common.Address, caller bind.ContractCaller) (*ConstantSupplyERC20Caller, error) {
	contract, err := bindConstantSupplyERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Caller{contract: contract}, nil
}

// NewConstantSupplyERC20Transactor creates a new write-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ConstantSupplyERC20Transactor, error) {
	contract, err := bindConstantSupplyERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Transactor{contract: contract}, nil
}

// NewConstantSupplyERC20Filterer creates a new log filterer instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ConstantSupplyERC20Filterer, error) {
	contract, err := bindConstantSupplyERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Filterer{contract: contract}, nil
}

// bindConstantSupplyERC20 binds a generic wrapper to an already deployed contract.
func bindConstantSupplyERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// ConstantSupplyERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20ApprovalIterator struct {
	Event *ConstantSupplyERC20Approval // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Approval)
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
		it.Event = new(ConstantSupplyERC20Approval)
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
func (it *ConstantSupplyERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Approval represents a Approval event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ConstantSupplyERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20ApprovalIterator{contract: _ConstantSupplyERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Approval)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseApproval(log types.Log) (*ConstantSupplyERC20Approval, error) {
	event := new(ConstantSupplyERC20Approval)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConstantSupplyERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20TransferIterator struct {
	Event *ConstantSupplyERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Transfer)
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
		it.Event = new(ConstantSupplyERC20Transfer)
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
func (it *ConstantSupplyERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Transfer represents a Transfer event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConstantSupplyERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20TransferIterator{contract: _ConstantSupplyERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Transfer)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseTransfer(log types.Log) (*ConstantSupplyERC20Transfer, error) {
	event := new(ConstantSupplyERC20Transfer)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
