// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransparentUpgradeableProxy

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

// TransparentUpgradeableProxyMetaData contains all meta data concerning the TransparentUpgradeableProxy contract.
var TransparentUpgradeableProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProxyDeniedAdminAccess\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60a06040526040516111d43803806111d48339810160408190526100229161043b565b828161002e8282610086565b50508160405161003d90610311565b61004791906104ac565b604051809103905ff080158015610060573d5f5f3e3d5ffd5b506001600160a01b031660805261007e61007960805190565b6100e4565b505050610500565b61008f82610146565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156100d8576100d382826101bf565b505050565b6100e0610234565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6101235f5160206111b45f395f51905f52546001600160a01b031690565b826040516101329291906104ba565b60405180910390a161014381610255565b50565b806001600160a01b03163b5f0361017b5780604051634c9c8ce360e01b815260040161017291906104ac565b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b60605f5f846001600160a01b0316846040516101db91906104f6565b5f60405180830381855af49150503d805f8114610213576040519150601f19603f3d011682016040523d82523d5f602084013e610218565b606091505b509092509050610229858383610292565b925050505b92915050565b34156102535760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b03811661027e575f604051633173bdd160e11b815260040161017291906104ac565b805f5160206111b45f395f51905f5261019e565b6060826102a7576102a2826102e8565b6102e1565b81511580156102be57506001600160a01b0384163b155b156102de5783604051639996b31560e01b815260040161017291906104ac565b50805b9392505050565b8051156102f85780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b6106a680610b0e83390190565b5f6001600160a01b03821661022e565b6103378161031e565b8114610143575f5ffd5b805161022e8161032e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b03811182821017156103855761038561034c565b6040525050565b5f61039660405190565b90506103a28282610360565b919050565b5f6001600160401b038211156103bf576103bf61034c565b601f19601f83011660200192915050565b8281835e505f910152565b5f6103ed6103e8846103a7565b61038c565b9050828152838383011115610403576104035f5ffd5b6102e18360208301846103d0565b5f82601f830112610423576104235f5ffd5b81516104338482602086016103db565b949350505050565b5f5f5f60608486031215610450576104505f5ffd5b61045a8585610341565b92506104698560208601610341565b60408501519092506001600160401b03811115610487576104875f5ffd5b61049386828701610411565b9150509250925092565b6104a68161031e565b82525050565b6020810161022e828461049d565b604081016104c8828561049d565b6102e1602083018461049d565b5f6104de825190565b6104ec8185602086016103d0565b9290920192915050565b61022e81836104d5565b6080516105f76105175f395f601001526105f75ff3fe608060405261000c61000e565b005b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633036100c4575f357fffffffff00000000000000000000000000000000000000000000000000000000167f4f1ef28600000000000000000000000000000000000000000000000000000000146100ba576040517fd2b576ec00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6100c26100cc565b565b6100c26100fa565b5f806100db36600481846103c0565b8101906100e8919061051c565b915091506100f6828261010a565b5050565b6100c2610105610164565b61019b565b610113826101b9565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a280511561015c576101578282610260565b505050565b6100f66102d4565b5f6101967f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b365f5f375f5f365f845af43d5f5f3e8080156101b5573d5ff35b3d5ffd5b806001600160a01b03163b5f0361020757806040517f4c9c8ce30000000000000000000000000000000000000000000000000000000081526004016101fe919061057d565b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60605f5f846001600160a01b03168460405161027c91906105b7565b5f60405180830381855af49150503d805f81146102b4576040519150601f19603f3d011682016040523d82523d5f602084013e6102b9565b606091505b50915091506102c985838361030c565b925050505b92915050565b34156100c2576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060826103215761031c8261037b565b610374565b815115801561033857506001600160a01b0384163b155b1561037157836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016101fe919061057d565b50805b9392505050565b80511561038b5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b5f5f858511156103d1576103d15f5ffd5b838611156103e0576103e05f5ffd5b5050820193919092039150565b5f6001600160a01b0382166102ce565b610406816103ed565b81146103bd575f5ffd5b80356102ce816103fd565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff8211171561046e5761046e61041b565b6040525050565b5f61047f60405190565b905061048b8282610448565b919050565b5f67ffffffffffffffff8211156104a9576104a961041b565b601f19601f83011660200192915050565b82818337505f910152565b5f6104d76104d284610490565b610475565b90508281528383830111156104ed576104ed5f5ffd5b6103748360208301846104ba565b5f82601f83011261050d5761050d5f5ffd5b610374838335602085016104c5565b5f5f60408385031215610530576105305f5ffd5b61053a8484610410565b9150602083013567ffffffffffffffff811115610558576105585f5ffd5b610564858286016104fb565b9150509250929050565b610577816103ed565b82525050565b602081016102ce828461056e565b8281835e505f910152565b5f61059f825190565b6105ad81856020860161058b565b9290920192915050565b6102ce818361059656fea2646970667358221220ef8defc17bb0b61d9786ae4bf6a8bae264cd148835fbd62548b5d583020e977364736f6c634300081c0033608060405234801561000f575f5ffd5b506040516106a63803806106a683398101604081905261002e916100f3565b806001600160a01b038116610061575f604051631e4fbdf760e01b81526004016100589190610126565b60405180910390fd5b61006a81610071565b5050610134565b5f80546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b5f6001600160a01b0382165b92915050565b6100db816100c0565b81146100e5575f5ffd5b50565b80516100cc816100d2565b5f60208284031215610106576101065f5ffd5b61011083836100e8565b9392505050565b610120816100c0565b82525050565b602081016100cc8284610117565b610565806101415f395ff3fe608060405260043610610058575f3560e01c80639623609d116100415780639623609d1461009f578063ad3cb1cc146100b2578063f2fde38b14610107575f5ffd5b8063715018a61461005c5780638da5cb5b14610072575b5f5ffd5b348015610067575f5ffd5b50610070610126565b005b34801561007d575f5ffd5b505f546001600160a01b031660405161009691906102e9565b60405180910390f35b6100706100ad36600461043b565b610139565b3480156100bd575f5ffd5b506100fa6040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161009691906104d9565b348015610112575f5ffd5b506100706101213660046104ea565b6101bd565b61012e61021c565b6101375f610261565b565b61014161021c565b6040517f4f1ef2860000000000000000000000000000000000000000000000000000000081526001600160a01b03841690634f1ef28690349061018a9086908690600401610507565b5f604051808303818588803b1580156101a1575f5ffd5b505af11580156101b3573d5f5f3e3d5ffd5b5050505050505050565b6101c561021c565b6001600160a01b038116610210575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161020791906102e9565b60405180910390fd5b61021981610261565b50565b5f546001600160a01b0316331461013757336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161020791906102e9565b5f80546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b5f6001600160a01b0382165b92915050565b6102e3816102c8565b82525050565b602081016102d482846102da565b5f6102d4826102c8565b61030a816102f7565b8114610219575f5ffd5b80356102d481610301565b61030a816102c8565b80356102d48161031f565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff8211171561038657610386610333565b6040525050565b5f61039760405190565b90506103a38282610360565b919050565b5f67ffffffffffffffff8211156103c1576103c1610333565b601f19601f83011660200192915050565b82818337505f910152565b5f6103ef6103ea846103a8565b61038d565b9050828152838383011115610405576104055f5ffd5b6104138360208301846103d2565b9392505050565b5f82601f83011261042c5761042c5f5ffd5b610413838335602085016103dd565b5f5f5f60608486031215610450576104505f5ffd5b61045a8585610314565b92506104698560208601610328565b9150604084013567ffffffffffffffff811115610487576104875f5ffd5b6104938682870161041a565b9150509250925092565b8281835e505f910152565b5f6104b1825190565b8084526020840193506104c881856020860161049d565b601f01601f19169290920192915050565b6020808252810161041381846104a8565b5f602082840312156104fd576104fd5f5ffd5b6104138383610328565b6040810161051582856102da565b818103602083015261052781846104a8565b94935050505056fea2646970667358221220825046736e9ba7da84b283bda1535f423544d9a64889ef23b00a31fa58c37ed964736f6c634300081c0033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
}

// TransparentUpgradeableProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use TransparentUpgradeableProxyMetaData.ABI instead.
var TransparentUpgradeableProxyABI = TransparentUpgradeableProxyMetaData.ABI

// TransparentUpgradeableProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransparentUpgradeableProxyMetaData.Bin instead.
var TransparentUpgradeableProxyBin = TransparentUpgradeableProxyMetaData.Bin

// DeployTransparentUpgradeableProxy deploys a new Ethereum contract, binding an instance of TransparentUpgradeableProxy to it.
func DeployTransparentUpgradeableProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, initialOwner common.Address, _data []byte) (common.Address, *types.Transaction, *TransparentUpgradeableProxy, error) {
	parsed, err := TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransparentUpgradeableProxyBin), backend, _logic, initialOwner, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransparentUpgradeableProxy{TransparentUpgradeableProxyCaller: TransparentUpgradeableProxyCaller{contract: contract}, TransparentUpgradeableProxyTransactor: TransparentUpgradeableProxyTransactor{contract: contract}, TransparentUpgradeableProxyFilterer: TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

// TransparentUpgradeableProxy is an auto generated Go binding around an Ethereum contract.
type TransparentUpgradeableProxy struct {
	TransparentUpgradeableProxyCaller     // Read-only binding to the contract
	TransparentUpgradeableProxyTransactor // Write-only binding to the contract
	TransparentUpgradeableProxyFilterer   // Log filterer for contract events
}

// TransparentUpgradeableProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransparentUpgradeableProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransparentUpgradeableProxySession struct {
	Contract     *TransparentUpgradeableProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// TransparentUpgradeableProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransparentUpgradeableProxyCallerSession struct {
	Contract *TransparentUpgradeableProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// TransparentUpgradeableProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransparentUpgradeableProxyTransactorSession struct {
	Contract     *TransparentUpgradeableProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// TransparentUpgradeableProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransparentUpgradeableProxyRaw struct {
	Contract *TransparentUpgradeableProxy // Generic contract binding to access the raw methods on
}

// TransparentUpgradeableProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyCallerRaw struct {
	Contract *TransparentUpgradeableProxyCaller // Generic read-only contract binding to access the raw methods on
}

// TransparentUpgradeableProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyTransactorRaw struct {
	Contract *TransparentUpgradeableProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransparentUpgradeableProxy creates a new instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxy(address common.Address, backend bind.ContractBackend) (*TransparentUpgradeableProxy, error) {
	contract, err := bindTransparentUpgradeableProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxy{TransparentUpgradeableProxyCaller: TransparentUpgradeableProxyCaller{contract: contract}, TransparentUpgradeableProxyTransactor: TransparentUpgradeableProxyTransactor{contract: contract}, TransparentUpgradeableProxyFilterer: TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

// NewTransparentUpgradeableProxyCaller creates a new read-only instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyCaller(address common.Address, caller bind.ContractCaller) (*TransparentUpgradeableProxyCaller, error) {
	contract, err := bindTransparentUpgradeableProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyCaller{contract: contract}, nil
}

// NewTransparentUpgradeableProxyTransactor creates a new write-only instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*TransparentUpgradeableProxyTransactor, error) {
	contract, err := bindTransparentUpgradeableProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyTransactor{contract: contract}, nil
}

// NewTransparentUpgradeableProxyFilterer creates a new log filterer instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*TransparentUpgradeableProxyFilterer, error) {
	contract, err := bindTransparentUpgradeableProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyFilterer{contract: contract}, nil
}

// bindTransparentUpgradeableProxy binds a generic wrapper to an already deployed contract.
func bindTransparentUpgradeableProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransparentUpgradeableProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.Fallback(&_TransparentUpgradeableProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.Fallback(&_TransparentUpgradeableProxy.TransactOpts, calldata)
}

// TransparentUpgradeableProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyAdminChangedIterator struct {
	Event *TransparentUpgradeableProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransparentUpgradeableProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransparentUpgradeableProxyAdminChanged)
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
		it.Event = new(TransparentUpgradeableProxyAdminChanged)
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
func (it *TransparentUpgradeableProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransparentUpgradeableProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransparentUpgradeableProxyAdminChanged represents a AdminChanged event raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TransparentUpgradeableProxyAdminChangedIterator, error) {

	logs, sub, err := _TransparentUpgradeableProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyAdminChangedIterator{contract: _TransparentUpgradeableProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TransparentUpgradeableProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TransparentUpgradeableProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransparentUpgradeableProxyAdminChanged)
				if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) ParseAdminChanged(log types.Log) (*TransparentUpgradeableProxyAdminChanged, error) {
	event := new(TransparentUpgradeableProxyAdminChanged)
	if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransparentUpgradeableProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyUpgradedIterator struct {
	Event *TransparentUpgradeableProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *TransparentUpgradeableProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransparentUpgradeableProxyUpgraded)
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
		it.Event = new(TransparentUpgradeableProxyUpgraded)
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
func (it *TransparentUpgradeableProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransparentUpgradeableProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransparentUpgradeableProxyUpgraded represents a Upgraded event raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TransparentUpgradeableProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransparentUpgradeableProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyUpgradedIterator{contract: _TransparentUpgradeableProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TransparentUpgradeableProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransparentUpgradeableProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransparentUpgradeableProxyUpgraded)
				if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) ParseUpgraded(log types.Log) (*TransparentUpgradeableProxyUpgraded, error) {
	event := new(TransparentUpgradeableProxyUpgraded)
	if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
