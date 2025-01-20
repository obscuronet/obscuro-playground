// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacksTest

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

// PublicCallbacksTestMetaData contains all meta data concerning the PublicCallbacksTest contract.
var PublicCallbacksTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_callbacks\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"callbackRefundees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"contractIPublicCallbacks\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleAllCallbacksRan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedGas\",\"type\":\"uint256\"}],\"name\":\"handleCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"handleCallbackFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"handleRefund\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isLastCallSuccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingRefunds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060408190525f805461ffff60a01b1916815560035561087b3881900390819083398101604081905261003291610304565b5f80546001600160b01b0319166001600160a01b03831617905561005461005a565b506103d6565b5f4861006760033461033c565b610071919061033c565b90505f63a072d7b060e01b8260405160240161008d919061034f565b60408051601f19818403018152918152602080830180516001600160e01b039081166001600160e01b031990961695909517905281516004808252602480830185528284018051881663a4c016fb60e01b179052845191825281019093529082018051909416629e79db60e81b179093525f805492945090916001600160a01b03166382fbdc9c61011f60033461033c565b866040518363ffffffff1660e01b815260040161013c9190610397565b60206040518083038185885af1158015610158573d5f5f3e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019061017d91906103b9565b5f81815260016020526040812080546001600160a01b03191633179055549091506001600160a01b03166382fbdc9c6101b760033461033c565b856040518363ffffffff1660e01b81526004016101d49190610397565b60206040518083038185885af11580156101f0573d5f5f3e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019061021591906103b9565b5f81815260016020526040812080546001600160a01b03191633179055549091506001600160a01b03166382fbdc9c61024f60033461033c565b846040518363ffffffff1660e01b815260040161026c9190610397565b60206040518083038185885af1158015610288573d5f5f3e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906102ad91906103b9565b5f90815260016020526040902080546001600160a01b031916331790555050505050565b5f6001600160a01b0382165b92915050565b6102ec816102d1565b81146102f6575f5ffd5b50565b80516102dd816102e3565b5f60208284031215610317576103175f5ffd5b61032183836102f9565b9392505050565b634e487b7160e01b5f52601260045260245ffd5b5f8261034a5761034a610328565b500490565b818152602081016102dd565b8281835e505f910152565b5f61036f825190565b80845260208401935061038681856020860161035b565b601f01601f19169290920192915050565b602080825281016103218184610366565b806102ec565b80516102dd816103a8565b5f602082840312156103cc576103cc5f5ffd5b61032183836103ae565b610498806103e35f395ff3fe608060405260043610610078575f3560e01c8063a072d7b01161004c578063a072d7b014610144578063a4c016fb14610163578063b613b11414610177578063ee1d5872146101af575f5ffd5b8062b127831461007c5780635ea39558146100b05780638103ab13146100c55780639e79db0014610106575b5f5ffd5b348015610087575f5ffd5b505f5461009a906001600160a01b031681565b6040516100a791906102ee565b60405180910390f35b6100c36100be366004610316565b6101d0565b005b3480156100d0575f5ffd5b506100f96100df366004610316565b60016020525f90815260409020546001600160a01b031681565b6040516100a79190610353565b348015610111575f5ffd5b506100c35f80547fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff16600160a81b179055565b34801561014f575f5ffd5b506100c361015e366004610316565b61021f565b34801561016e575f5ffd5b506100c361024c565b348015610182575f5ffd5b506101a2610191366004610375565b60026020525f908152604090205481565b6040516100a79190610398565b3480156101ba575f5ffd5b506101c3610299565b6040516100a791906103ae565b5f818152600160209081526040808320546001600160a01b031683526002909152812080543492906102039084906103e9565b909155505060038054905f610217836103fc565b919050555050565b5f5a905061022f61083483610414565b8110610248575f805460ff60a01b1916600160a01b1790555b5050565b5f805460ff60a01b1916600160a01b1790556040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029090610427565b60405180910390fd5b5f8054600160a81b900460ff1680156102b457506003546003145b905090565b5f6001600160a01b0382165b92915050565b5f6102c5826102b9565b5f6102c5826102cb565b6102e8816102d5565b82525050565b602081016102c582846102df565b805b8114610308575f5ffd5b50565b80356102c5816102fc565b5f60208284031215610329576103295f5ffd5b610333838361030b565b9392505050565b5f6001600160a01b0382166102c5565b6102e88161033a565b602081016102c5828461034a565b6102fe8161033a565b80356102c581610361565b5f60208284031215610388576103885f5ffd5b610333838361036a565b806102e8565b602081016102c58284610392565b8015156102e8565b602081016102c582846103a6565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b808201808211156102c5576102c56103bc565b5f6001820161040d5761040d6103bc565b5060010190565b818103818111156102c5576102c56103bc565b602080825281016102c581601681527f5468697320697320612074657374206661696c7572650000000000000000000060208201526040019056fea264697066735822122014fd76ae51330580df6594926ee6d8a78b07172786d1f84457d5fb9570ce540564736f6c634300081c0033",
}

// PublicCallbacksTestABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksTestMetaData.ABI instead.
var PublicCallbacksTestABI = PublicCallbacksTestMetaData.ABI

// PublicCallbacksTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksTestMetaData.Bin instead.
var PublicCallbacksTestBin = PublicCallbacksTestMetaData.Bin

// DeployPublicCallbacksTest deploys a new Ethereum contract, binding an instance of PublicCallbacksTest to it.
func DeployPublicCallbacksTest(auth *bind.TransactOpts, backend bind.ContractBackend, _callbacks common.Address) (common.Address, *types.Transaction, *PublicCallbacksTest, error) {
	parsed, err := PublicCallbacksTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksTestBin), backend, _callbacks)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacksTest{PublicCallbacksTestCaller: PublicCallbacksTestCaller{contract: contract}, PublicCallbacksTestTransactor: PublicCallbacksTestTransactor{contract: contract}, PublicCallbacksTestFilterer: PublicCallbacksTestFilterer{contract: contract}}, nil
}

// PublicCallbacksTest is an auto generated Go binding around an Ethereum contract.
type PublicCallbacksTest struct {
	PublicCallbacksTestCaller     // Read-only binding to the contract
	PublicCallbacksTestTransactor // Write-only binding to the contract
	PublicCallbacksTestFilterer   // Log filterer for contract events
}

// PublicCallbacksTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksTestSession struct {
	Contract     *PublicCallbacksTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PublicCallbacksTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksTestCallerSession struct {
	Contract *PublicCallbacksTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// PublicCallbacksTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTestTransactorSession struct {
	Contract     *PublicCallbacksTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// PublicCallbacksTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksTestRaw struct {
	Contract *PublicCallbacksTest // Generic contract binding to access the raw methods on
}

// PublicCallbacksTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksTestCallerRaw struct {
	Contract *PublicCallbacksTestCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTestTransactorRaw struct {
	Contract *PublicCallbacksTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacksTest creates a new instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTest(address common.Address, backend bind.ContractBackend) (*PublicCallbacksTest, error) {
	contract, err := bindPublicCallbacksTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTest{PublicCallbacksTestCaller: PublicCallbacksTestCaller{contract: contract}, PublicCallbacksTestTransactor: PublicCallbacksTestTransactor{contract: contract}, PublicCallbacksTestFilterer: PublicCallbacksTestFilterer{contract: contract}}, nil
}

// NewPublicCallbacksTestCaller creates a new read-only instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksTestCaller, error) {
	contract, err := bindPublicCallbacksTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestCaller{contract: contract}, nil
}

// NewPublicCallbacksTestTransactor creates a new write-only instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTestTransactor, error) {
	contract, err := bindPublicCallbacksTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestTransactor{contract: contract}, nil
}

// NewPublicCallbacksTestFilterer creates a new log filterer instance of PublicCallbacksTest, bound to a specific deployed contract.
func NewPublicCallbacksTestFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksTestFilterer, error) {
	contract, err := bindPublicCallbacksTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTestFilterer{contract: contract}, nil
}

// bindPublicCallbacksTest binds a generic wrapper to an already deployed contract.
func bindPublicCallbacksTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacksTest *PublicCallbacksTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.PublicCallbacksTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacksTest *PublicCallbacksTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacksTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacksTest *PublicCallbacksTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacksTest *PublicCallbacksTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.contract.Transact(opts, method, params...)
}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) CallbackRefundees(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "callbackRefundees", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestSession) CallbackRefundees(arg0 *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, arg0)
}

// CallbackRefundees is a free data retrieval call binding the contract method 0x8103ab13.
//
// Solidity: function callbackRefundees(uint256 ) view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) CallbackRefundees(arg0 *big.Int) (common.Address, error) {
	return _PublicCallbacksTest.Contract.CallbackRefundees(&_PublicCallbacksTest.CallOpts, arg0)
}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) Callbacks(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "callbacks")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestSession) Callbacks() (common.Address, error) {
	return _PublicCallbacksTest.Contract.Callbacks(&_PublicCallbacksTest.CallOpts)
}

// Callbacks is a free data retrieval call binding the contract method 0x00b12783.
//
// Solidity: function callbacks() view returns(address)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) Callbacks() (common.Address, error) {
	return _PublicCallbacksTest.Contract.Callbacks(&_PublicCallbacksTest.CallOpts)
}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) IsLastCallSuccess(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "isLastCallSuccess")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestSession) IsLastCallSuccess() (bool, error) {
	return _PublicCallbacksTest.Contract.IsLastCallSuccess(&_PublicCallbacksTest.CallOpts)
}

// IsLastCallSuccess is a free data retrieval call binding the contract method 0xee1d5872.
//
// Solidity: function isLastCallSuccess() view returns(bool)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) IsLastCallSuccess() (bool, error) {
	return _PublicCallbacksTest.Contract.IsLastCallSuccess(&_PublicCallbacksTest.CallOpts)
}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestCaller) PendingRefunds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicCallbacksTest.contract.Call(opts, &out, "pendingRefunds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestSession) PendingRefunds(arg0 common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, arg0)
}

// PendingRefunds is a free data retrieval call binding the contract method 0xb613b114.
//
// Solidity: function pendingRefunds(address ) view returns(uint256)
func (_PublicCallbacksTest *PublicCallbacksTestCallerSession) PendingRefunds(arg0 common.Address) (*big.Int, error) {
	return _PublicCallbacksTest.Contract.PendingRefunds(&_PublicCallbacksTest.CallOpts, arg0)
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleAllCallbacksRan(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleAllCallbacksRan")
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleAllCallbacksRan() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleAllCallbacksRan(&_PublicCallbacksTest.TransactOpts)
}

// HandleAllCallbacksRan is a paid mutator transaction binding the contract method 0x9e79db00.
//
// Solidity: function handleAllCallbacksRan() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleAllCallbacksRan() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleAllCallbacksRan(&_PublicCallbacksTest.TransactOpts)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleCallback(opts *bind.TransactOpts, expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleCallback", expectedGas)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleCallback(expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallback(&_PublicCallbacksTest.TransactOpts, expectedGas)
}

// HandleCallback is a paid mutator transaction binding the contract method 0xa072d7b0.
//
// Solidity: function handleCallback(uint256 expectedGas) returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleCallback(expectedGas *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallback(&_PublicCallbacksTest.TransactOpts, expectedGas)
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleCallbackFail(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleCallbackFail")
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleCallbackFail() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallbackFail(&_PublicCallbacksTest.TransactOpts)
}

// HandleCallbackFail is a paid mutator transaction binding the contract method 0xa4c016fb.
//
// Solidity: function handleCallbackFail() returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleCallbackFail() (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleCallbackFail(&_PublicCallbacksTest.TransactOpts)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactor) HandleRefund(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.contract.Transact(opts, "handleRefund", callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}

// HandleRefund is a paid mutator transaction binding the contract method 0x5ea39558.
//
// Solidity: function handleRefund(uint256 callbackId) payable returns()
func (_PublicCallbacksTest *PublicCallbacksTestTransactorSession) HandleRefund(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacksTest.Contract.HandleRefund(&_PublicCallbacksTest.TransactOpts, callbackId)
}
