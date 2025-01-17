// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package GasConsumerBalance

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

// GasConsumerBalanceMetaData contains all meta data concerning the GasConsumerBalance contract.
var GasConsumerBalanceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"destroy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_balance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"resetOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b505f80546001600160a01b0319163317905561022f8061002d5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c806373cc802a1461004e57806383197ef0146100a45780638da5cb5b146100ac578063c1cfb99a146100a2575b5f5ffd5b6100a261005c36600461017d565b5f80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b005b6100a26100e1565b5f546100cb9073ffffffffffffffffffffffffffffffffffffffff1681565b6040516100d891906101b0565b60405180910390f35b5f5473ffffffffffffffffffffffffffffffffffffffff16331461013a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610131906101be565b60405180910390fd5b30ff5b5f73ffffffffffffffffffffffffffffffffffffffff82165b92915050565b6101658161013d565b811461016f575f5ffd5b50565b80356101568161015c565b5f60208284031215610190576101905f5ffd5b61019a8383610172565b9392505050565b6101aa8161013d565b82525050565b6020810161015682846101a1565b6020808252810161015681601581527f596f7520617265206e6f7420746865206f776e6572000000000000000000000060208201526040019056fea26469706673582212202a35857581c80dfb42043b9bc221c33e7afeeda7ce917e39c6f919b2b08c310964736f6c634300081c0033",
}

// GasConsumerBalanceABI is the input ABI used to generate the binding from.
// Deprecated: Use GasConsumerBalanceMetaData.ABI instead.
var GasConsumerBalanceABI = GasConsumerBalanceMetaData.ABI

// GasConsumerBalanceBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GasConsumerBalanceMetaData.Bin instead.
var GasConsumerBalanceBin = GasConsumerBalanceMetaData.Bin

// DeployGasConsumerBalance deploys a new Ethereum contract, binding an instance of GasConsumerBalance to it.
func DeployGasConsumerBalance(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GasConsumerBalance, error) {
	parsed, err := GasConsumerBalanceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GasConsumerBalanceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GasConsumerBalance{GasConsumerBalanceCaller: GasConsumerBalanceCaller{contract: contract}, GasConsumerBalanceTransactor: GasConsumerBalanceTransactor{contract: contract}, GasConsumerBalanceFilterer: GasConsumerBalanceFilterer{contract: contract}}, nil
}

// GasConsumerBalance is an auto generated Go binding around an Ethereum contract.
type GasConsumerBalance struct {
	GasConsumerBalanceCaller     // Read-only binding to the contract
	GasConsumerBalanceTransactor // Write-only binding to the contract
	GasConsumerBalanceFilterer   // Log filterer for contract events
}

// GasConsumerBalanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type GasConsumerBalanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasConsumerBalanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GasConsumerBalanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasConsumerBalanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GasConsumerBalanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasConsumerBalanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GasConsumerBalanceSession struct {
	Contract     *GasConsumerBalance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// GasConsumerBalanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GasConsumerBalanceCallerSession struct {
	Contract *GasConsumerBalanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// GasConsumerBalanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GasConsumerBalanceTransactorSession struct {
	Contract     *GasConsumerBalanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// GasConsumerBalanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type GasConsumerBalanceRaw struct {
	Contract *GasConsumerBalance // Generic contract binding to access the raw methods on
}

// GasConsumerBalanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GasConsumerBalanceCallerRaw struct {
	Contract *GasConsumerBalanceCaller // Generic read-only contract binding to access the raw methods on
}

// GasConsumerBalanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GasConsumerBalanceTransactorRaw struct {
	Contract *GasConsumerBalanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGasConsumerBalance creates a new instance of GasConsumerBalance, bound to a specific deployed contract.
func NewGasConsumerBalance(address common.Address, backend bind.ContractBackend) (*GasConsumerBalance, error) {
	contract, err := bindGasConsumerBalance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GasConsumerBalance{GasConsumerBalanceCaller: GasConsumerBalanceCaller{contract: contract}, GasConsumerBalanceTransactor: GasConsumerBalanceTransactor{contract: contract}, GasConsumerBalanceFilterer: GasConsumerBalanceFilterer{contract: contract}}, nil
}

// NewGasConsumerBalanceCaller creates a new read-only instance of GasConsumerBalance, bound to a specific deployed contract.
func NewGasConsumerBalanceCaller(address common.Address, caller bind.ContractCaller) (*GasConsumerBalanceCaller, error) {
	contract, err := bindGasConsumerBalance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GasConsumerBalanceCaller{contract: contract}, nil
}

// NewGasConsumerBalanceTransactor creates a new write-only instance of GasConsumerBalance, bound to a specific deployed contract.
func NewGasConsumerBalanceTransactor(address common.Address, transactor bind.ContractTransactor) (*GasConsumerBalanceTransactor, error) {
	contract, err := bindGasConsumerBalance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GasConsumerBalanceTransactor{contract: contract}, nil
}

// NewGasConsumerBalanceFilterer creates a new log filterer instance of GasConsumerBalance, bound to a specific deployed contract.
func NewGasConsumerBalanceFilterer(address common.Address, filterer bind.ContractFilterer) (*GasConsumerBalanceFilterer, error) {
	contract, err := bindGasConsumerBalance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GasConsumerBalanceFilterer{contract: contract}, nil
}

// bindGasConsumerBalance binds a generic wrapper to an already deployed contract.
func bindGasConsumerBalance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GasConsumerBalanceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasConsumerBalance *GasConsumerBalanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasConsumerBalance.Contract.GasConsumerBalanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasConsumerBalance *GasConsumerBalanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.GasConsumerBalanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasConsumerBalance *GasConsumerBalanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.GasConsumerBalanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasConsumerBalance *GasConsumerBalanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasConsumerBalance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasConsumerBalance *GasConsumerBalanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasConsumerBalance *GasConsumerBalanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GasConsumerBalance *GasConsumerBalanceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GasConsumerBalance.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GasConsumerBalance *GasConsumerBalanceSession) Owner() (common.Address, error) {
	return _GasConsumerBalance.Contract.Owner(&_GasConsumerBalance.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GasConsumerBalance *GasConsumerBalanceCallerSession) Owner() (common.Address, error) {
	return _GasConsumerBalance.Contract.Owner(&_GasConsumerBalance.CallOpts)
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactor) Destroy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasConsumerBalance.contract.Transact(opts, "destroy")
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_GasConsumerBalance *GasConsumerBalanceSession) Destroy() (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.Destroy(&_GasConsumerBalance.TransactOpts)
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactorSession) Destroy() (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.Destroy(&_GasConsumerBalance.TransactOpts)
}

// GetBalance is a paid mutator transaction binding the contract method 0xc1cfb99a.
//
// Solidity: function get_balance() returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactor) GetBalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasConsumerBalance.contract.Transact(opts, "get_balance")
}

// GetBalance is a paid mutator transaction binding the contract method 0xc1cfb99a.
//
// Solidity: function get_balance() returns()
func (_GasConsumerBalance *GasConsumerBalanceSession) GetBalance() (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.GetBalance(&_GasConsumerBalance.TransactOpts)
}

// GetBalance is a paid mutator transaction binding the contract method 0xc1cfb99a.
//
// Solidity: function get_balance() returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactorSession) GetBalance() (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.GetBalance(&_GasConsumerBalance.TransactOpts)
}

// ResetOwner is a paid mutator transaction binding the contract method 0x73cc802a.
//
// Solidity: function resetOwner(address _owner) returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactor) ResetOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _GasConsumerBalance.contract.Transact(opts, "resetOwner", _owner)
}

// ResetOwner is a paid mutator transaction binding the contract method 0x73cc802a.
//
// Solidity: function resetOwner(address _owner) returns()
func (_GasConsumerBalance *GasConsumerBalanceSession) ResetOwner(_owner common.Address) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.ResetOwner(&_GasConsumerBalance.TransactOpts, _owner)
}

// ResetOwner is a paid mutator transaction binding the contract method 0x73cc802a.
//
// Solidity: function resetOwner(address _owner) returns()
func (_GasConsumerBalance *GasConsumerBalanceTransactorSession) ResetOwner(_owner common.Address) (*types.Transaction, error) {
	return _GasConsumerBalance.Contract.ResetOwner(&_GasConsumerBalance.TransactOpts, _owner)
}
