// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacks

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

// PublicCallbacksMetaData contains all meta data concerning the PublicCallbacks contract.
var PublicCallbacksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBefore\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfter\",\"type\":\"uint256\"}],\"name\":\"CallbackExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeNextCallbacks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"reattemptCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callback\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5060156019565b60c9565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560685760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c65780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b610fa7806100d65f395ff3fe608060405260043610610057575f3560e01c806382fbdc9c1161004157806382fbdc9c146100a9578063929d34e9146100c9578063a67e1760146100e8575f5ffd5b8062e0d3b51461005b5780638129fc1c14610093575b5f5ffd5b348015610066575f5ffd5b5061007a610075366004610a43565b6100fc565b60405161008a9493929190610ac1565b60405180910390f35b34801561009e575f5ffd5b506100a76101b4565b005b6100bc6100b7366004610b54565b6102f3565b60405161008a9190610b99565b3480156100d4575f5ffd5b506100a76100e3366004610a43565b61035a565b3480156100f3575f5ffd5b506100a76104f4565b5f60208190529081526040902080546001820180546001600160a01b03909216929161012790610bbb565b80601f016020809104026020016040519081016040528092919081815260200182805461015390610bbb565b801561019e5780601f106101755761010080835404028352916020019161019e565b820191905f5260205f20905b81548152906001019060200180831161018157829003601f168201915b5050505050908060020154908060030154905084565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156101fe5750825b90505f8267ffffffffffffffff16600114801561021a5750303b155b905081158015610228575080155b1561025f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561029357845468ff00000000000000001916680100000000000000001785555b5f600181905560025583156102ec57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906102e390600190610c0a565b60405180910390a15b5050505050565b5f5f341161031c5760405162461bcd60e51b815260040161031390610c4c565b60405180910390fd5b61520861032834610545565b116103455760405162461bcd60e51b815260040161031390610c5c565b61035133848434610550565b90505b92915050565b5f8181526020818152604080832081516080810190925280546001600160a01b03168252600181018054929391929184019161039590610bbb565b80601f01602080910402602001604051908101604052809291908181526020018280546103c190610bbb565b801561040c5780601f106103e35761010080835404028352916020019161040c565b820191905f5260205f20905b8154815290600101906020018083116103ef57829003601f168201915b505050505081526020016002820154815260200160038201548152505090505f815f01516001600160a01b0316826020015160405161044b9190610cde565b5f604051808303815f865af19150503d805f8114610484576040519150601f19603f3d011682016040523d82523d5f602084013e610489565b606091505b50509050806104aa5760405162461bcd60e51b815260040161031390610d1a565b5f838152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff19168155906104e060018301826109f2565b505f60028201819055600390910155505050565b5f610500600130610d3e565b9050336001600160a01b0382161461052a5760405162461bcd60e51b815260040161031390610d93565b600254600154146105425761053d610637565b61052a565b50565b5f6103544883610db7565b5f60015490506040518060800160405280866001600160a01b0316815260200185858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201829052509385525050506020820185905248604090920191909152600180548291826105c883610dca565b9091555081526020808201929092526040015f208251815473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039091161781559082015160018201906106189082610e8b565b5060408201516002820155606090910151600390910155949350505050565b6002546001540361064457565b5f5f61064e610754565b915091505f826060015190505f81846040015161066b9190610db7565b90505f5a90505f855f01516001600160a01b03168387602001516040516106929190610cde565b5f604051808303815f8787f1925050503d805f81146106cc576040519150601f19603f3d011682016040523d82523d5f602084013e6106d1565b606091505b505090505f5a90505f6106e48285610f47565b90505f8186111561070757866106fa8388610f47565b6107049190610f5a565b90505b5f818a604001516107189190610f47565b8a51909150851561072b5761072b610866565b6107336108b0565b61073e83828c6108c6565b610747826109a0565b5050505050505050505050565b61078560405180608001604052805f6001600160a01b03168152602001606081526020015f81526020015f81525090565b6002545f8181526020818152604080832081516080810190925280546001600160a01b03168252600181018054949591949193859290840191906107c890610bbb565b80601f01602080910402602001604051908101604052809291908181526020018280546107f490610bbb565b801561083f5780601f106108165761010080835404028352916020019161083f565b820191905f5260205f20905b81548152906001019060200180831161082257829003601f168201915b50505050508152602001600282015481526020016003820154815250509150915091509091565b6002545f908152602081905260408120805473ffffffffffffffffffffffffffffffffffffffff191681559061089f60018301826109f2565b505f60028201819055600390910155565b60028054905f6108bf83610dca565b9190505550565b5f816040516024016108d89190610b99565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f5ea3955800000000000000000000000000000000000000000000000000000000179052519091505f906001600160a01b0385169061afc8908790610952908690610cde565b5f60405180830381858888f193505050503d805f811461098d576040519150601f19603f3d011682016040523d82523d5f602084013e610992565b606091505b50509050806102ec576102ec855b805f036109aa5750565b604051419082905f81818185875af1925050503d805f81146109e7576040519150601f19603f3d011682016040523d82523d5f602084013e6109ec565b606091505b50505050565b5080546109fe90610bbb565b5f825580601f10610a0d575050565b601f0160209004905f5260205f209081019061054291905b80821115610a38575f8155600101610a25565b5090565b8035610354565b5f60208284031215610a5657610a565f5ffd5b6103518383610a3c565b5f6001600160a01b038216610354565b610a7981610a60565b82525050565b8281835e505f910152565b5f610a93825190565b808452602084019350610aaa818560208601610a7f565b601f01601f19169290920192915050565b80610a79565b60808101610acf8287610a70565b8181036020830152610ae18186610a8a565b9050610af06040830185610abb565b610afd6060830184610abb565b95945050505050565b5f5f83601f840112610b1957610b195f5ffd5b50813567ffffffffffffffff811115610b3357610b335f5ffd5b602083019150836001820283011115610b4d57610b4d5f5ffd5b9250929050565b5f5f60208385031215610b6857610b685f5ffd5b823567ffffffffffffffff811115610b8157610b815f5ffd5b610b8d85828601610b06565b92509250509250929050565b602081016103548284610abb565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610bcf57607f821691505b602082108103610be157610be1610ba7565b50919050565b5f61035482610bf4565b90565b67ffffffffffffffff1690565b610a7981610be7565b602081016103548284610c01565b600d8152602081017f4e6f2076616c75652073656e7400000000000000000000000000000000000000815290505b60200190565b6020808252810161035481610c18565b6020808252810161035481602481527f47617320746f6f206c6f7720636f6d706172656420746f20636f7374206f662060208201527f63616c6c00000000000000000000000000000000000000000000000000000000604082015260600190565b5f610cc6825190565b610cd4818560208601610a7f565b9290920192915050565b6103548183610cbd565b60198152602081017f43616c6c6261636b20657865637574696f6e206661696c65640000000000000081529050610c46565b6020808252810161035481610ce8565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b0391821691908116908282039081111561035457610354610d2a565b60088152602081017f4e6f742073656c6600000000000000000000000000000000000000000000000081529050610c46565b6020808252810161035481610d61565b634e487b7160e01b5f52601260045260245ffd5b5f82610dc557610dc5610da3565b500490565b5f60018201610ddb57610ddb610d2a565b5060010190565b634e487b7160e01b5f52604160045260245ffd5b5f610354610bf18381565b610e0a83610df6565b81545f1960089490940293841b1916921b91909117905550565b5f610e30818484610e01565b505050565b81811015610e4f57610e475f82610e24565b600101610e35565b5050565b601f821115610e30575f818152602090206020601f85010481016020851015610e795750805b6102ec6020601f860104830182610e35565b815167ffffffffffffffff811115610ea557610ea5610de2565b610eaf8254610bbb565b610eba828285610e53565b506020601f821160018114610eed575f8315610ed65750848201515b5f19600885021c19811660028502178555506102ec565b5f84815260208120601f198516915b82811015610f1c5787850151825560209485019460019092019101610efc565b5084821015610f3857838701515f19601f87166008021c191681555b50505050600202600101905550565b8181038181111561035457610354610d2a565b818102811582820484141761035457610354610d2a56fea2646970667358221220d077d37ad6740210bd32dc45c1bb1fda2c3418880ff1ee5e28a8f5f61443e06964736f6c634300081c0033",
}

// PublicCallbacksABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksMetaData.ABI instead.
var PublicCallbacksABI = PublicCallbacksMetaData.ABI

// PublicCallbacksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksMetaData.Bin instead.
var PublicCallbacksBin = PublicCallbacksMetaData.Bin

// DeployPublicCallbacks deploys a new Ethereum contract, binding an instance of PublicCallbacks to it.
func DeployPublicCallbacks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PublicCallbacks, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// PublicCallbacks is an auto generated Go binding around an Ethereum contract.
type PublicCallbacks struct {
	PublicCallbacksCaller     // Read-only binding to the contract
	PublicCallbacksTransactor // Write-only binding to the contract
	PublicCallbacksFilterer   // Log filterer for contract events
}

// PublicCallbacksCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksSession struct {
	Contract     *PublicCallbacks  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicCallbacksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksCallerSession struct {
	Contract *PublicCallbacksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// PublicCallbacksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTransactorSession struct {
	Contract     *PublicCallbacksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PublicCallbacksRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksRaw struct {
	Contract *PublicCallbacks // Generic contract binding to access the raw methods on
}

// PublicCallbacksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksCallerRaw struct {
	Contract *PublicCallbacksCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactorRaw struct {
	Contract *PublicCallbacksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacks creates a new instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacks(address common.Address, backend bind.ContractBackend) (*PublicCallbacks, error) {
	contract, err := bindPublicCallbacks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// NewPublicCallbacksCaller creates a new read-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksCaller, error) {
	contract, err := bindPublicCallbacks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCaller{contract: contract}, nil
}

// NewPublicCallbacksTransactor creates a new write-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTransactor, error) {
	contract, err := bindPublicCallbacks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTransactor{contract: contract}, nil
}

// NewPublicCallbacksFilterer creates a new log filterer instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksFilterer, error) {
	contract, err := bindPublicCallbacks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksFilterer{contract: contract}, nil
}

// bindPublicCallbacks binds a generic wrapper to an already deployed contract.
func bindPublicCallbacks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.PublicCallbacksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transact(opts, method, params...)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksCaller) Callbacks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbacks", arg0)

	outstruct := new(struct {
		Target  common.Address
		Data    []byte
		Value   *big.Int
		BaseFee *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Data = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Value = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.BaseFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksSession) Callbacks(arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 ) view returns(address target, bytes data, uint256 value, uint256 baseFee)
func (_PublicCallbacks *PublicCallbacksCallerSession) Callbacks(arg0 *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, arg0)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ExecuteNextCallbacks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "executeNextCallbacks")
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ReattemptCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "reattemptCallback", callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactor) Register(opts *bind.TransactOpts, callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "register", callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactorSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// PublicCallbacksCallbackExecutedIterator is returned from FilterCallbackExecuted and is used to iterate over the raw logs and unpacked data for CallbackExecuted events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecutedIterator struct {
	Event *PublicCallbacksCallbackExecuted // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackExecuted)
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
		it.Event = new(PublicCallbacksCallbackExecuted)
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
func (it *PublicCallbacksCallbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackExecuted represents a CallbackExecuted event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecuted struct {
	CallbackId *big.Int
	GasBefore  *big.Int
	GasAfter   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackExecuted is a free log retrieval operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackExecuted(opts *bind.FilterOpts) (*PublicCallbacksCallbackExecutedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackExecutedIterator{contract: _PublicCallbacks.contract, event: "CallbackExecuted", logs: logs, sub: sub}, nil
}

// WatchCallbackExecuted is a free log subscription operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackExecuted(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackExecuted) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackExecuted)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
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

// ParseCallbackExecuted is a log parse operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackExecuted(log types.Log) (*PublicCallbacksCallbackExecuted, error) {
	event := new(PublicCallbacksCallbackExecuted)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicCallbacksInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PublicCallbacks contract.
type PublicCallbacksInitializedIterator struct {
	Event *PublicCallbacksInitialized // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksInitialized)
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
		it.Event = new(PublicCallbacksInitialized)
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
func (it *PublicCallbacksInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksInitialized represents a Initialized event raised by the PublicCallbacks contract.
type PublicCallbacksInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterInitialized(opts *bind.FilterOpts) (*PublicCallbacksInitializedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksInitializedIterator{contract: _PublicCallbacks.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PublicCallbacksInitialized) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksInitialized)
				if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseInitialized(log types.Log) (*PublicCallbacksInitialized, error) {
	event := new(PublicCallbacksInitialized)
	if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
