// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Structs

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

// StructsTransaction is an auto generated low-level Go binding around an user-defined struct.
type StructsTransaction struct {
	TxType               uint8
	ChainId              *big.Int
	Nonce                *big.Int
	GasPrice             *big.Int
	GasLimit             *big.Int
	To                   common.Address
	Value                *big.Int
	Data                 []byte
	V                    uint8
	R                    [32]byte
	S                    [32]byte
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
	AccessList           []common.Address
}

// StructsMetaData contains all meta data concerning the Structs contract.
var StructsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"accessList\",\"type\":\"address[]\"}],\"internalType\":\"structStructs.Transaction\",\"name\":\"txData\",\"type\":\"tuple\"}],\"name\":\"recoverSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523461001e5760405161069e61002482393081505061069e90f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c638314c06f0361003657610072565b90816101c09103126100365790565b600080fd5b9060208282031261003657813567ffffffffffffffff8111610036576100619201610027565b90565b6001600160a01b031690565b565b6100a261008861008336600461003b565b6102b6565b604051918291826001600160a01b03909116815260200190565b0390f35b60ff81165b0361003657565b35610061816100a6565b806100ab565b35610061816100bc565b6001600160a01b0381166100ab565b35610061816100cc565b903590601e193682900301821215610036570180359067ffffffffffffffff8211610036576020019136829003831361003657565b903590601e193682900301821215610036570180359067ffffffffffffffff82116100365760200191602082023603831361003657565b90826000939282370152565b919061017b81610174816101859560209181520190565b8095610151565b601f01601f191690565b0190565b90503590610070826100cc565b50610061906020810190610189565b818352602090920191906000825b8282106101c1575050505090565b909192936101f06101e96001926101d88886610196565b6001600160a01b0316815260200190565b9560200190565b939201906101b3565b999b9c9a9895939196949290966101608b019760008c0161021b9160ff169052565b60208b015260408a0152606089015260808801526001600160a01b031660a087015260c086015284820360e08601526102539261015d565b936101008301610261919052565b610120820152808303906101400152610061926101a5565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff8211176102b157604052565b610279565b6102bf816100b2565b6102cb602083016100c2565b6102d7604084016100c2565b906102e4606085016100c2565b926102f1608086016100c2565b6102fd60a087016100db565b9061030a60c088016100c2565b61031760e08901896100e5565b6103246101608b016100c2565b916103326101808c016100c2565b936103416101a08d018d61011a565b97909661034d60405190565b9c8d9c60208e019c61035f9c8e6101f9565b90810382520361036f908261028f565b805190602001206103a9907f19457468657265756d205369676e6564204d6573736167653a0a333200000000600052601c52603c60002090565b906103b761010082016100b2565b906103c561012082016100c2565b90610140016103d3906100c2565b6100619384936103e39391610458565b90929192610550565b6100616100616100619290565b610061906103ec565b6100646100616100619290565b61006190610402565b61044861007094610441606094989795610437608086019a6000870152565b60ff166020850152565b6040830152565b0152565b6040513d6000823e3d90fd5b9091610463846103f9565b61049361048f7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a06103ec565b9190565b1161050d57906104b5602094600094936104ac60405190565b94859485610418565b838052039060015afa15610508576000516000916104d28361040f565b6001600160a01b0381166001600160a01b038416146104fb57506104f5836103ec565b91929190565b91506104f56001936103ec565b61044c565b50505061051a600061040f565b9160039190565b634e487b7160e01b600052602160045260246000fd5b6004111561054157565b610521565b9061007082610537565b61055a6000610546565b61056382610546565b0361056c575050565b6105766001610546565b61057f82610546565b036105b3576040517ff645eedf000000000000000000000000000000000000000000000000000000008152600490fd5b0390fd5b6105bd6002610546565b6105c682610546565b0361060d576105af6105d7836103f9565b6040519182917ffce698f70000000000000000000000000000000000000000000000000000000083526004830190815260200190565b61062061061a6003610546565b91610546565b146106285750565b6105af9061063560405190565b9182917fd78bce0c000000000000000000000000000000000000000000000000000000008352600483019081526020019056fea2646970667358221220c65c0636aae937e4abd84684e03270a4926f575a45974a78aecc6bdde8f5898e64736f6c63430008140033",
}

// StructsABI is the input ABI used to generate the binding from.
// Deprecated: Use StructsMetaData.ABI instead.
var StructsABI = StructsMetaData.ABI

// StructsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StructsMetaData.Bin instead.
var StructsBin = StructsMetaData.Bin

// DeployStructs deploys a new Ethereum contract, binding an instance of Structs to it.
func DeployStructs(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Structs, error) {
	parsed, err := StructsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StructsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Structs{StructsCaller: StructsCaller{contract: contract}, StructsTransactor: StructsTransactor{contract: contract}, StructsFilterer: StructsFilterer{contract: contract}}, nil
}

// Structs is an auto generated Go binding around an Ethereum contract.
type Structs struct {
	StructsCaller     // Read-only binding to the contract
	StructsTransactor // Write-only binding to the contract
	StructsFilterer   // Log filterer for contract events
}

// StructsCaller is an auto generated read-only Go binding around an Ethereum contract.
type StructsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StructsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StructsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StructsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StructsSession struct {
	Contract     *Structs          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StructsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StructsCallerSession struct {
	Contract *StructsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StructsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StructsTransactorSession struct {
	Contract     *StructsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StructsRaw is an auto generated low-level Go binding around an Ethereum contract.
type StructsRaw struct {
	Contract *Structs // Generic contract binding to access the raw methods on
}

// StructsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StructsCallerRaw struct {
	Contract *StructsCaller // Generic read-only contract binding to access the raw methods on
}

// StructsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StructsTransactorRaw struct {
	Contract *StructsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStructs creates a new instance of Structs, bound to a specific deployed contract.
func NewStructs(address common.Address, backend bind.ContractBackend) (*Structs, error) {
	contract, err := bindStructs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Structs{StructsCaller: StructsCaller{contract: contract}, StructsTransactor: StructsTransactor{contract: contract}, StructsFilterer: StructsFilterer{contract: contract}}, nil
}

// NewStructsCaller creates a new read-only instance of Structs, bound to a specific deployed contract.
func NewStructsCaller(address common.Address, caller bind.ContractCaller) (*StructsCaller, error) {
	contract, err := bindStructs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StructsCaller{contract: contract}, nil
}

// NewStructsTransactor creates a new write-only instance of Structs, bound to a specific deployed contract.
func NewStructsTransactor(address common.Address, transactor bind.ContractTransactor) (*StructsTransactor, error) {
	contract, err := bindStructs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StructsTransactor{contract: contract}, nil
}

// NewStructsFilterer creates a new log filterer instance of Structs, bound to a specific deployed contract.
func NewStructsFilterer(address common.Address, filterer bind.ContractFilterer) (*StructsFilterer, error) {
	contract, err := bindStructs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StructsFilterer{contract: contract}, nil
}

// bindStructs binds a generic wrapper to an already deployed contract.
func bindStructs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StructsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Structs *StructsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Structs.Contract.StructsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Structs *StructsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Structs.Contract.StructsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Structs *StructsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Structs.Contract.StructsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Structs *StructsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Structs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Structs *StructsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Structs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Structs *StructsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Structs.Contract.contract.Transact(opts, method, params...)
}

// RecoverSender is a free data retrieval call binding the contract method 0x23bc02dd.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[]) txData) pure returns(address sender)
func (_Structs *StructsCaller) RecoverSender(opts *bind.CallOpts, txData StructsTransaction) (common.Address, error) {
	var out []interface{}
	err := _Structs.contract.Call(opts, &out, "recoverSender", txData)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSender is a free data retrieval call binding the contract method 0x23bc02dd.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[]) txData) pure returns(address sender)
func (_Structs *StructsSession) RecoverSender(txData StructsTransaction) (common.Address, error) {
	return _Structs.Contract.RecoverSender(&_Structs.CallOpts, txData)
}

// RecoverSender is a free data retrieval call binding the contract method 0x23bc02dd.
//
// Solidity: function recoverSender((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[]) txData) pure returns(address sender)
func (_Structs *StructsCallerSession) RecoverSender(txData StructsTransaction) (common.Address, error) {
	return _Structs.Contract.RecoverSender(&_Structs.CallOpts, txData)
}
