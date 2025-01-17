// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC1967Proxy

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

// ERC1967ProxyMetaData contains all meta data concerning the ERC1967Proxy contract.
var ERC1967ProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60806040526040516104923803806104928339810160408190526100229161033c565b61002c8282610033565b50506103d6565b61003c82610091565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115610085576100808282610107565b505050565b61008d61017c565b5050565b806001600160a01b03163b5f036100c65780604051634c9c8ce360e01b81526004016100bd919061039d565b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b60605f5f846001600160a01b03168460405161012391906103cc565b5f60405180830381855af49150503d805f811461015b576040519150601f19603f3d011682016040523d82523d5f602084013e610160565b606091505b50909250905061017185838361019d565b925050505b92915050565b341561019b5760405163b398979f60e01b815260040160405180910390fd5b565b6060826101b2576101ad826101f3565b6101ec565b81511580156101c957506001600160a01b0384163b155b156101e95783604051639996b31560e01b81526004016100bd919061039d565b50805b9392505050565b8051156102035780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b50565b5f6001600160a01b038216610176565b6102388161021f565b811461021c575f5ffd5b80516101768161022f565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b03811182821017156102865761028661024d565b6040525050565b5f61029760405190565b90506102a38282610261565b919050565b5f6001600160401b038211156102c0576102c061024d565b601f19601f83011660200192915050565b8281835e505f910152565b5f6102ee6102e9846102a8565b61028d565b9050828152838383011115610304576103045f5ffd5b6101ec8360208301846102d1565b5f82601f830112610324576103245f5ffd5b81516103348482602086016102dc565b949350505050565b5f5f60408385031215610350576103505f5ffd5b61035a8484610242565b60208401519092506001600160401b03811115610378576103785f5ffd5b61038485828601610312565b9150509250929050565b6103978161021f565b82525050565b60208101610176828461038e565b5f6103b4825190565b6103c28185602086016102d1565b9290920192915050565b61017681836103ab565b60b0806103e25f395ff3fe6080604052600a600c565b005b60186014601a565b605d565b565b5f60587f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b905090565b365f5f375f5f365f845af43d5f5f3e8080156076573d5ff35b3d5ffdfea2646970667358221220009f1261a1ec8b4d87ff54c25d816495d207b04d42d2b605d3be8e6a2406578264736f6c634300081c0033",
}

// ERC1967ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967ProxyMetaData.ABI instead.
var ERC1967ProxyABI = ERC1967ProxyMetaData.ABI

// ERC1967ProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC1967ProxyMetaData.Bin instead.
var ERC1967ProxyBin = ERC1967ProxyMetaData.Bin

// DeployERC1967Proxy deploys a new Ethereum contract, binding an instance of ERC1967Proxy to it.
func DeployERC1967Proxy(auth *bind.TransactOpts, backend bind.ContractBackend, implementation common.Address, _data []byte) (common.Address, *types.Transaction, *ERC1967Proxy, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC1967ProxyBin), backend, implementation, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// ERC1967Proxy is an auto generated Go binding around an Ethereum contract.
type ERC1967Proxy struct {
	ERC1967ProxyCaller     // Read-only binding to the contract
	ERC1967ProxyTransactor // Write-only binding to the contract
	ERC1967ProxyFilterer   // Log filterer for contract events
}

// ERC1967ProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC1967ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC1967ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC1967ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC1967ProxySession struct {
	Contract     *ERC1967Proxy     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967ProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC1967ProxyCallerSession struct {
	Contract *ERC1967ProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC1967ProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC1967ProxyTransactorSession struct {
	Contract     *ERC1967ProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC1967ProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC1967ProxyRaw struct {
	Contract *ERC1967Proxy // Generic contract binding to access the raw methods on
}

// ERC1967ProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC1967ProxyCallerRaw struct {
	Contract *ERC1967ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967ProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC1967ProxyTransactorRaw struct {
	Contract *ERC1967ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Proxy creates a new instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967Proxy(address common.Address, backend bind.ContractBackend) (*ERC1967Proxy, error) {
	contract, err := bindERC1967Proxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// NewERC1967ProxyCaller creates a new read-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyCaller(address common.Address, caller bind.ContractCaller) (*ERC1967ProxyCaller, error) {
	contract, err := bindERC1967Proxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyCaller{contract: contract}, nil
}

// NewERC1967ProxyTransactor creates a new write-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967ProxyTransactor, error) {
	contract, err := bindERC1967Proxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyTransactor{contract: contract}, nil
}

// NewERC1967ProxyFilterer creates a new log filterer instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967ProxyFilterer, error) {
	contract, err := bindERC1967Proxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyFilterer{contract: contract}, nil
}

// bindERC1967Proxy binds a generic wrapper to an already deployed contract.
func bindERC1967Proxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.ERC1967ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// ERC1967ProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgradedIterator struct {
	Event *ERC1967ProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967ProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyUpgraded)
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
		it.Event = new(ERC1967ProxyUpgraded)
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
func (it *ERC1967ProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyUpgraded represents a Upgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967ProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyUpgradedIterator{contract: _ERC1967Proxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyUpgraded)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseUpgraded(log types.Log) (*ERC1967ProxyUpgraded, error) {
	event := new(ERC1967ProxyUpgraded)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
