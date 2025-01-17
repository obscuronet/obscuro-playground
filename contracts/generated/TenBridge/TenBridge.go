// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TenBridge

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

// TenBridgeMetaData contains all meta data concerning the TenBridge contract.
var TenBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messengerAddress\",\"type\":\"address\"}],\"name\":\"configure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b19169055348015601e575f5ffd5b50611adf8061002c5f395ff3fe60806040526004361061013d575f3560e01c806375b238fc116100bb578063a217fddf11610071578063c4d66de811610057578063c4d66de81461038e578063d547741f146103ad578063e4c3ebc7146103cc575f5ffd5b8063a217fddf14610368578063a381c8e21461037b575f5ffd5b806383bece4d116100a157806383bece4d146102e657806391d148541461030557806393b3744214610349575f5ffd5b806375b238fc1461029457806375cb2672146102c7575f5ffd5b80632f2ff15d11610110578063498d82ab116100f6578063498d82ab146102235780635d872970146102425780635fa7b58414610275575f5ffd5b80632f2ff15d146101e557806336568abe14610204575f5ffd5b806301ffc9a71461014157806316ce8149146101765780631888d71214610197578063248a9ca3146101aa575b5f5ffd5b34801561014c575f5ffd5b5061016061015b3660046112d3565b6103ff565b60405161016d91906112fa565b60405180910390f35b348015610181575f5ffd5b5061019561019036600461132c565b610467565b005b6101956101a536600461132c565b6104c1565b3480156101b5575f5ffd5b506101d86101c436600461135a565b5f9081526002602052604090206001015490565b60405161016d919061137d565b3480156101f0575f5ffd5b506101956101ff36600461138b565b6105c6565b34801561020f575f5ffd5b5061019561021e36600461138b565b6105f0565b34801561022e575f5ffd5b5061019561023d36600461140f565b610641565b34801561024d575f5ffd5b506101d87f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b348015610280575f5ffd5b5061019561028f36600461132c565b610705565b34801561029f575f5ffd5b506101d87fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b3480156102d2575f5ffd5b506101956102e136600461132c565b610759565b3480156102f1575f5ffd5b50610195610300366004611496565b610831565b348015610310575f5ffd5b5061016061031f36600461138b565b5f9182526002602090815260408084206001600160a01b0393909316845291905290205460ff1690565b348015610354575f5ffd5b5061019561036336600461132c565b610948565b348015610373575f5ffd5b506101d85f81565b610195610389366004611496565b61099c565b348015610399575f5ffd5b506101956103a836600461132c565b610a87565b3480156103b8575f5ffd5b506101956103c736600461138b565b610c1d565b3480156103d7575f5ffd5b506101d87fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b5f6001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061046157507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561049181610c41565b506003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b5f34116104e95760405162461bcd60e51b81526004016104e090611510565b60405180910390fd5b5f6040518060400160405280348152602001836001600160a01b03168152506040516020016105189190611548565b60408051601f19818403018152919052600354909150610547906001600160a01b03168260025b5f5f5f610c4e565b6001546040517f346633fb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063346633fb9034906105949086908390600401611556565b5f604051808303818588803b1580156105ab575f5ffd5b505af11580156105bd573d5f5f3e3d5ffd5b50505050505050565b5f828152600260205260409020600101546105e081610c41565b6105ea8383610d62565b50505050565b6001600160a01b0381163314610632576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61063c8282610e0d565b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561066b81610c41565b6106957f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a87610d62565b505f63458ffd6360e01b87878787876040516024016106b895949392919061159c565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526003549091506105bd906001600160a01b031682600161053f565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561072f81610c41565b61063c7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a83610e0d565b610761610e92565b5f805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa1580156107dd573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061080191906115e8565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905550565b6003545f546001600160a01b03918216911633146108615760405162461bcd60e51b81526004016104e09061165f565b806001600160a01b0316610873610efb565b6001600160a01b0316146108995760405162461bcd60e51b81526004016104e0906116c7565b6001600160a01b0384165f9081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16156108e8576108e3848484610f74565b6105ea565b6001600160a01b0384165f9081527f13ad2d85210d477fe1a6e25654c8250308cf29b050a4bf0b039d70467486712c602052604090205460ff1615610930576108e382610f7f565b60405162461bcd60e51b81526004016104e09061172f565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561097281610c41565b61063c7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610d62565b5f82116109bb5760405162461bcd60e51b81526004016104e090611771565b6001600160a01b0383165f9081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16610a115760405162461bcd60e51b81526004016104e090611781565b610a1d83333085610fee565b5f6383bece4d60e01b848484604051602401610a3b93929190611808565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526003549091506105ea906001600160a01b0316825f61053f565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f81158015610ad15750825b90505f8267ffffffffffffffff166001148015610aed5750303b155b905081158015610afb575080155b15610b32576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610b6657845468ff00000000000000001916680100000000000000001785555b610b6f86610759565b610b997fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533610d62565b50610bc47fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105785f610d62565b508315610c1557845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610c0c90600190611852565b60405180910390a15b505050505050565b5f82815260026020526040902060010154610c3781610c41565b6105ea8383610e0d565b610c4b8133611048565b50565b5f6040518060600160405280886001600160a01b0316815260200187815260200185815250604051602001610c8391906118e4565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa91859174010000000000000000000000000000000000000000900463ffffffff16906014610cd783611922565b91906101000a81548163ffffffff021916908363ffffffff1602179055508885886040518663ffffffff1660e01b8152600401610d179493929190611959565b60206040518083038185885af1158015610d33573d5f5f3e3d5ffd5b50505050506040513d601f19601f82011682018060405250810190610d5891906119b9565b5050505050505050565b5f8281526002602090815260408083206001600160a01b038516845290915281205460ff16610e06575f8381526002602090815260408083206001600160a01b03861684529091529020805460ff19166001179055610dbe3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610461565b505f610461565b5f8281526002602090815260408083206001600160a01b038516845290915281205460ff1615610e06575f8381526002602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4506001610461565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610ef9576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f5f5f9054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f4b573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f6f91906115e8565b905090565b61063c8382846110a7565b6040515f906001600160a01b038316908281818181865af19150503d805f8114610fc4576040519150601f19603f3d011682016040523d82523d5f602084013e610fc9565b606091505b5050905080610fea5760405162461bcd60e51b81526004016104e090611a08565b5050565b6105ea84856001600160a01b03166323b872dd86868660405160240161101693929190611a18565b604051602081830303815290604052915060e01b6020820180516001600160e01b0383818316178352505050506110cd565b5f8281526002602090815260408083206001600160a01b038516845290915290205460ff16610fea5780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016104e0929190611556565b61063c83846001600160a01b031663a9059cbb8585604051602401611016929190611556565b5f6110e16001600160a01b0384168361113e565b905080515f141580156111055750808060200190518101906111039190611a53565b155b1561063c57826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016104e09190611a70565b606061114b83835f611152565b9392505050565b60608147101561119057306040517fcd7860590000000000000000000000000000000000000000000000000000000081526004016104e09190611a70565b5f5f856001600160a01b031684866040516111ab9190611a9f565b5f6040518083038185875af1925050503d805f81146111e5576040519150601f19603f3d011682016040523d82523d5f602084013e6111ea565b606091505b50915091506111fa868383611204565b9695505050505050565b6060826112195761121482611270565b61114b565b815115801561123057506001600160a01b0384163b155b1561126957836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016104e09190611a70565b508061114b565b8051156112805780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160e01b031981165b8114610c4b575f5ffd5b8035610461816112b2565b5f602082840312156112e6576112e65f5ffd5b61114b83836112c8565b8015155b82525050565b6020810161046182846112f0565b5f6001600160a01b038216610461565b6112be81611308565b803561046181611318565b5f6020828403121561133f5761133f5f5ffd5b61114b8383611321565b806112be565b803561046181611349565b5f6020828403121561136d5761136d5f5ffd5b61114b838361134f565b806112f4565b602081016104618284611377565b5f5f6040838503121561139f5761139f5f5ffd5b6113a9848461134f565b91506113b88460208501611321565b90509250929050565b5f5f83601f8401126113d4576113d45f5ffd5b50813567ffffffffffffffff8111156113ee576113ee5f5ffd5b602083019150836001820283011115611408576114085f5ffd5b9250929050565b5f5f5f5f5f60608688031215611426576114265f5ffd5b6114308787611321565b9450602086013567ffffffffffffffff81111561144e5761144e5f5ffd5b61145a888289016113c1565b9450945050604086013567ffffffffffffffff81111561147b5761147b5f5ffd5b611487888289016113c1565b92509250509295509295909350565b5f5f5f606084860312156114ab576114ab5f5ffd5b6114b58585611321565b92506114c4856020860161134f565b91506114d38560408601611321565b90509250925092565b600f8152602081017f456d707479207472616e736665722e0000000000000000000000000000000000815290505b60200190565b60208082528101610461816114dc565b6112f481611308565b80516115358382611377565b50602081015161063c6020840182611520565b604081016104618284611529565b604081016115648285611520565b61114b6020830184611377565b82818337505f910152565b818352602083019250611590828483611571565b50601f01601f19160190565b606081016115aa8288611520565b81810360208301526115bd81868861157c565b905081810360408301526115d281848661157c565b979650505050505050565b805161046181611318565b5f602082840312156115fb576115fb5f5ffd5b61114b83836115dd565b60308152602081017f436f6e74726163742063616c6c6572206973206e6f742074686520726567697381527f7465726564206d657373656e6765722100000000000000000000000000000000602082015290505b60400190565b6020808252810161046181611605565b60318152602081017f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2081527f696e636f72726563742073656e6465722100000000000000000000000000000060208201529050611659565b602080825281016104618161166f565b60258152602081017f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206181527f737365742e00000000000000000000000000000000000000000000000000000060208201529050611659565b60208082528101610461816116d7565b601a8152602081017f417474656d7074696e6720656d707479207472616e736665722e0000000000008152905061150a565b602080825281016104618161173f565b6020808252810161046181604e81527f54686973206164647265737320686173206e6f74206265656e20676976656e2060208201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60408201527f742077686974656c69737465642e000000000000000000000000000000000000606082015260800190565b606081016118168286611520565b6118236020830185611377565b6118306040830184611520565b949350505050565b5f67ffffffffffffffff8216610461565b6112f481611838565b602081016104618284611849565b8281835e505f910152565b5f611874825190565b80845260208401935061188b818560208601611860565b601f01601f19169290920192915050565b80515f9060608401906118af8582611520565b50602083015184820360208601526118c7828261186b565b91505060408301516118dc6040860182611377565b509392505050565b6020808252810161114b818461189c565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b63ffffffff165f63fffffffe19820161193d5761193d6118f5565b5060010190565b63ffffffff81166112f4565b60ff81166112f4565b608081016119678287611944565b6119746020830186611944565b8181036040830152611986818561186b565b90506119956060830184611950565b95945050505050565b67ffffffffffffffff81166112be565b80516104618161199e565b5f602082840312156119cc576119cc5f5ffd5b61114b83836119ae565b60148152602081017f4661696c656420746f2073656e642045746865720000000000000000000000008152905061150a565b60208082528101610461816119d6565b60608101611a268286611520565b611a336020830185611520565b6118306040830184611377565b8015156112be565b805161046181611a40565b5f60208284031215611a6657611a665f5ffd5b61114b8383611a48565b602081016104618284611520565b5f611a87825190565b611a95818560208601611860565b9290920192915050565b6104618183611a7e56fea26469706673582212209dcfe92519eb101dc9fd8937e3c4f599ad995e20466aea28eb78608ac2d8d01864736f6c634300081c0033",
}

// TenBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use TenBridgeMetaData.ABI instead.
var TenBridgeABI = TenBridgeMetaData.ABI

// TenBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TenBridgeMetaData.Bin instead.
var TenBridgeBin = TenBridgeMetaData.Bin

// DeployTenBridge deploys a new Ethereum contract, binding an instance of TenBridge to it.
func DeployTenBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TenBridge, error) {
	parsed, err := TenBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TenBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TenBridge{TenBridgeCaller: TenBridgeCaller{contract: contract}, TenBridgeTransactor: TenBridgeTransactor{contract: contract}, TenBridgeFilterer: TenBridgeFilterer{contract: contract}}, nil
}

// TenBridge is an auto generated Go binding around an Ethereum contract.
type TenBridge struct {
	TenBridgeCaller     // Read-only binding to the contract
	TenBridgeTransactor // Write-only binding to the contract
	TenBridgeFilterer   // Log filterer for contract events
}

// TenBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TenBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TenBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TenBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TenBridgeSession struct {
	Contract     *TenBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TenBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TenBridgeCallerSession struct {
	Contract *TenBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TenBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TenBridgeTransactorSession struct {
	Contract     *TenBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TenBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TenBridgeRaw struct {
	Contract *TenBridge // Generic contract binding to access the raw methods on
}

// TenBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TenBridgeCallerRaw struct {
	Contract *TenBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// TenBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TenBridgeTransactorRaw struct {
	Contract *TenBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTenBridge creates a new instance of TenBridge, bound to a specific deployed contract.
func NewTenBridge(address common.Address, backend bind.ContractBackend) (*TenBridge, error) {
	contract, err := bindTenBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TenBridge{TenBridgeCaller: TenBridgeCaller{contract: contract}, TenBridgeTransactor: TenBridgeTransactor{contract: contract}, TenBridgeFilterer: TenBridgeFilterer{contract: contract}}, nil
}

// NewTenBridgeCaller creates a new read-only instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeCaller(address common.Address, caller bind.ContractCaller) (*TenBridgeCaller, error) {
	contract, err := bindTenBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeCaller{contract: contract}, nil
}

// NewTenBridgeTransactor creates a new write-only instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*TenBridgeTransactor, error) {
	contract, err := bindTenBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTransactor{contract: contract}, nil
}

// NewTenBridgeFilterer creates a new log filterer instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*TenBridgeFilterer, error) {
	contract, err := bindTenBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TenBridgeFilterer{contract: contract}, nil
}

// bindTenBridge binds a generic wrapper to an already deployed contract.
func bindTenBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TenBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridge *TenBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridge.Contract.TenBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridge *TenBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridge.Contract.TenBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridge *TenBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridge.Contract.TenBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridge *TenBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridge *TenBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridge *TenBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridge.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) ADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.ADMINROLE(&_TenBridge.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) ADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.ADMINROLE(&_TenBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.DEFAULTADMINROLE(&_TenBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.DEFAULTADMINROLE(&_TenBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) ERC20TOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "ERC20_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.ERC20TOKENROLE(&_TenBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.ERC20TOKENROLE(&_TenBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) NATIVETOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "NATIVE_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.NATIVETOKENROLE(&_TenBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.NATIVETOKENROLE(&_TenBridge.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridge.Contract.GetRoleAdmin(&_TenBridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridge.Contract.GetRoleAdmin(&_TenBridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridge.Contract.HasRole(&_TenBridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridge.Contract.HasRole(&_TenBridge.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridge.Contract.SupportsInterface(&_TenBridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridge.Contract.SupportsInterface(&_TenBridge.CallOpts, interfaceId)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeTransactor) Configure(opts *bind.TransactOpts, messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "configure", messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Configure(&_TenBridge.TransactOpts, messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeTransactorSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Configure(&_TenBridge.TransactOpts, messengerAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.GrantRole(&_TenBridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.GrantRole(&_TenBridge.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_TenBridge *TenBridgeTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "initialize", messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_TenBridge *TenBridgeSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Initialize(&_TenBridge.TransactOpts, messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_TenBridge *TenBridgeTransactorSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Initialize(&_TenBridge.TransactOpts, messenger)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeTransactor) PromoteToAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "promoteToAdmin", newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PromoteToAdmin(&_TenBridge.TransactOpts, newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeTransactorSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PromoteToAdmin(&_TenBridge.TransactOpts, newAdmin)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.ReceiveAssets(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.ReceiveAssets(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_TenBridge *TenBridgeTransactor) RemoveToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "removeToken", asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_TenBridge *TenBridgeSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RemoveToken(&_TenBridge.TransactOpts, asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_TenBridge *TenBridgeTransactorSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RemoveToken(&_TenBridge.TransactOpts, asset)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RenounceRole(&_TenBridge.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RenounceRole(&_TenBridge.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RevokeRole(&_TenBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RevokeRole(&_TenBridge.TransactOpts, role, account)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendERC20(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendERC20(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendNative(&_TenBridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendNative(&_TenBridge.TransactOpts, receiver)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeTransactor) SetRemoteBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "setRemoteBridge", bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SetRemoteBridge(&_TenBridge.TransactOpts, bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeTransactorSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SetRemoteBridge(&_TenBridge.TransactOpts, bridge)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeTransactor) WhitelistToken(opts *bind.TransactOpts, asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "whitelistToken", asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.Contract.WhitelistToken(&_TenBridge.TransactOpts, asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeTransactorSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.Contract.WhitelistToken(&_TenBridge.TransactOpts, asset, name, symbol)
}

// TenBridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TenBridge contract.
type TenBridgeInitializedIterator struct {
	Event *TenBridgeInitialized // Event containing the contract specifics and raw log

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
func (it *TenBridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeInitialized)
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
		it.Event = new(TenBridgeInitialized)
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
func (it *TenBridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeInitialized represents a Initialized event raised by the TenBridge contract.
type TenBridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridge *TenBridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*TenBridgeInitializedIterator, error) {

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TenBridgeInitializedIterator{contract: _TenBridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridge *TenBridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TenBridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeInitialized)
				if err := _TenBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseInitialized(log types.Log) (*TenBridgeInitialized, error) {
	event := new(TenBridgeInitialized)
	if err := _TenBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TenBridge contract.
type TenBridgeRoleAdminChangedIterator struct {
	Event *TenBridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleAdminChanged)
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
		it.Event = new(TenBridgeRoleAdminChanged)
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
func (it *TenBridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleAdminChanged represents a RoleAdminChanged event raised by the TenBridge contract.
type TenBridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridge *TenBridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TenBridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleAdminChangedIterator{contract: _TenBridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridge *TenBridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleAdminChanged)
				if err := _TenBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleAdminChanged(log types.Log) (*TenBridgeRoleAdminChanged, error) {
	event := new(TenBridgeRoleAdminChanged)
	if err := _TenBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TenBridge contract.
type TenBridgeRoleGrantedIterator struct {
	Event *TenBridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleGranted)
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
		it.Event = new(TenBridgeRoleGranted)
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
func (it *TenBridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleGranted represents a RoleGranted event raised by the TenBridge contract.
type TenBridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleGrantedIterator{contract: _TenBridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleGranted)
				if err := _TenBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleGranted(log types.Log) (*TenBridgeRoleGranted, error) {
	event := new(TenBridgeRoleGranted)
	if err := _TenBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TenBridge contract.
type TenBridgeRoleRevokedIterator struct {
	Event *TenBridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleRevoked)
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
		it.Event = new(TenBridgeRoleRevoked)
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
func (it *TenBridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleRevoked represents a RoleRevoked event raised by the TenBridge contract.
type TenBridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleRevokedIterator{contract: _TenBridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleRevoked)
				if err := _TenBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleRevoked(log types.Log) (*TenBridgeRoleRevoked, error) {
	event := new(TenBridgeRoleRevoked)
	if err := _TenBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
