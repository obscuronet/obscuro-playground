// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObscuroBridge

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

// ObscuroBridgeMetaData contains all meta data concerning the ObscuroBridge contract.
var ObscuroBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messengerAddress\",\"type\":\"address\"}],\"name\":\"configure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b19169055348015601f57600080fd5b50611b5d8061002f6000396000f3fe6080604052600436106101445760003560e01c806375b238fc116100c0578063a217fddf11610074578063c4d66de811610059578063c4d66de8146103a7578063d547741f146103c7578063e4c3ebc7146103e757600080fd5b8063a217fddf1461037f578063a381c8e21461039457600080fd5b806383bece4d116100a557806383bece4d146102f957806391d148541461031957806393b374421461035f57600080fd5b806375b238fc146102a557806375cb2672146102d957600080fd5b80632f2ff15d11610117578063498d82ab116100fc578063498d82ab146102315780635d872970146102515780635fa7b5841461028557600080fd5b80632f2ff15d146101f157806336568abe1461021157600080fd5b806301ffc9a71461014957806316ce81491461017f5780631888d712146101a1578063248a9ca3146101b4575b600080fd5b34801561015557600080fd5b50610169610164366004611314565b61041b565b604051610176919061133d565b60405180910390f35b34801561018b57600080fd5b5061019f61019a366004611370565b610484565b005b61019f6101af366004611370565b6104de565b3480156101c057600080fd5b506101e46101cf3660046113a0565b60009081526002602052604090206001015490565b60405161017691906113c5565b3480156101fd57600080fd5b5061019f61020c3660046113d3565b6105e9565b34801561021d57600080fd5b5061019f61022c3660046113d3565b610614565b34801561023d57600080fd5b5061019f61024c36600461145d565b610665565b34801561025d57600080fd5b506101e47f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b34801561029157600080fd5b5061019f6102a0366004611370565b61072a565b3480156102b157600080fd5b506101e47fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b3480156102e557600080fd5b5061019f6102f4366004611370565b61077e565b34801561030557600080fd5b5061019f6103143660046114ea565b610859565b34801561032557600080fd5b506101696103343660046113d3565b60009182526002602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561036b57600080fd5b5061019f61037a366004611370565b610973565b34801561038b57600080fd5b506101e4600081565b61019f6103a23660046114ea565b6109c7565b3480156103b357600080fd5b5061019f6103c2366004611370565b610ab6565b3480156103d357600080fd5b5061019f6103e23660046113d3565b610c4f565b3480156103f357600080fd5b506101e47fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b60006001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061047e57507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756104ae81610c74565b506003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b600034116105075760405162461bcd60e51b81526004016104fe90611567565b60405180910390fd5b60006040518060400160405280348152602001836001600160a01b0316815250604051602001610537919061159f565b60408051601f19818403018152919052600354909150610566906001600160a01b03168260025b600080610c81565b6001546040517f346633fb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063346633fb9034906105b390869083906004016115ad565b6000604051808303818588803b1580156105cc57600080fd5b505af11580156105e0573d6000803e3d6000fd5b50505050505050565b60008281526002602052604090206001015461060481610c74565b61060e8383610d8e565b50505050565b6001600160a01b0381163314610656576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106608282610e3c565b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561068f81610c74565b6106b97f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a87610d8e565b50600063458ffd6360e01b87878787876040516024016106dd9594939291906115f4565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526003549091506105e0906001600160a01b031682600161055e565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561075481610c74565b6106607f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a83610e3c565b610786610ec3565b6000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa158015610805573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108299190611640565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905550565b6003546000546001600160a01b039182169116331461088a5760405162461bcd60e51b81526004016104fe906116b9565b806001600160a01b031661089c610f2c565b6001600160a01b0316146108c25760405162461bcd60e51b81526004016104fe90611721565b6001600160a01b03841660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16156109125761090d848484610fa9565b61060e565b6001600160a01b03841660009081527f13ad2d85210d477fe1a6e25654c8250308cf29b050a4bf0b039d70467486712c602052604090205460ff161561095b5761090d82610fb4565b60405162461bcd60e51b81526004016104fe90611789565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561099d81610c74565b6106607fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610d8e565b600082116109e75760405162461bcd60e51b81526004016104fe906117cb565b6001600160a01b03831660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16610a3e5760405162461bcd60e51b81526004016104fe906117db565b610a4a83333085611026565b60006383bece4d60e01b848484604051602401610a6993929190611862565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915260035490915061060e906001600160a01b031682600061055e565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff16600081158015610b015750825b905060008267ffffffffffffffff166001148015610b1e5750303b155b905081158015610b2c575080155b15610b63576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610b9757845468ff00000000000000001916680100000000000000001785555b610ba08661077e565b610bca7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533610d8e565b50610bf67fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105786000610d8e565b508315610c4757845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610c3e906001906118ad565b60405180910390a15b505050505050565b600082815260026020526040902060010154610c6a81610c74565b61060e8383610e3c565b610c7e8133611080565b50565b60006040518060600160405280876001600160a01b0316815260200186815260200184815250604051602001610cb7919061195a565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa91349174010000000000000000000000000000000000000000900463ffffffff16906014610d0b8361199a565b91906101000a81548163ffffffff021916908363ffffffff1602179055508785876040518663ffffffff1660e01b8152600401610d4b94939291906119d2565b60206040518083038185885af1158015610d69573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906105e09190611a32565b60008281526002602090815260408083206001600160a01b038516845290915281205460ff16610e345760008381526002602090815260408083206001600160a01b03861684529091529020805460ff19166001179055610dec3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161047e565b50600061047e565b60008281526002602090815260408083206001600160a01b038516845290915281205460ff1615610e345760008381526002602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161047e565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610f2a576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008060009054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f80573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fa49190611640565b905090565b6106608382846110e0565b6040516000906001600160a01b038316908281818181865af19150503d8060008114610ffc576040519150601f19603f3d011682016040523d82523d6000602084013e611001565b606091505b50509050806110225760405162461bcd60e51b81526004016104fe90611a83565b5050565b61060e84856001600160a01b03166323b872dd86868660405160240161104e93929190611a93565b604051602081830303815290604052915060e01b6020820180516001600160e01b038381831617835250505050611106565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff166110225780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016104fe9291906115ad565b61066083846001600160a01b031663a9059cbb858560405160240161104e9291906115ad565b600061111b6001600160a01b03841683611179565b9050805160001415801561114057508080602001905181019061113e9190611ace565b155b1561066057826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016104fe9190611aed565b60606111878383600061118e565b9392505050565b6060814710156111cc57306040517fcd7860590000000000000000000000000000000000000000000000000000000081526004016104fe9190611aed565b600080856001600160a01b031684866040516111e89190611b1d565b60006040518083038185875af1925050503d8060008114611225576040519150601f19603f3d011682016040523d82523d6000602084013e61122a565b606091505b509150915061123a868383611244565b9695505050505050565b60608261125957611254826112b0565b611187565b815115801561127057506001600160a01b0384163b155b156112a957836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016104fe9190611aed565b5080611187565b8051156112c05780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160e01b031981165b8114610c7e57600080fd5b803561047e816112f2565b60006020828403121561132957611329600080fd5b6111878383611309565b8015155b82525050565b6020810161047e8284611333565b60006001600160a01b03821661047e565b6112fe8161134b565b803561047e8161135c565b60006020828403121561138557611385600080fd5b6111878383611365565b806112fe565b803561047e8161138f565b6000602082840312156113b5576113b5600080fd5b6111878383611395565b80611337565b6020810161047e82846113bf565b600080604083850312156113e9576113e9600080fd5b6113f38484611395565b91506114028460208501611365565b90509250929050565b60008083601f84011261142057611420600080fd5b50813567ffffffffffffffff81111561143b5761143b600080fd5b60208301915083600182028301111561145657611456600080fd5b9250929050565b60008060008060006060868803121561147857611478600080fd5b6114828787611365565b9450602086013567ffffffffffffffff8111156114a1576114a1600080fd5b6114ad8882890161140b565b9450945050604086013567ffffffffffffffff8111156114cf576114cf600080fd5b6114db8882890161140b565b92509250509295509295909350565b60008060006060848603121561150257611502600080fd5b61150c8585611365565b925061151b8560208601611395565b915061152a8560408601611365565b90509250925092565b600f8152602081017f456d707479207472616e736665722e0000000000000000000000000000000000815290505b60200190565b6020808252810161047e81611533565b6113378161134b565b805161158c83826113bf565b5060208101516106606020840182611577565b6040810161047e8284611580565b604081016115bb8285611577565b61118760208301846113bf565b82818337506000910152565b8183526020830192506115e88284836115c8565b50601f01601f19160190565b606081016116028288611577565b81810360208301526116158186886115d4565b9050818103604083015261162a8184866115d4565b979650505050505050565b805161047e8161135c565b60006020828403121561165557611655600080fd5b6111878383611635565b60308152602081017f436f6e74726163742063616c6c6572206973206e6f742074686520726567697381527f7465726564206d657373656e6765722100000000000000000000000000000000602082015290505b60400190565b6020808252810161047e8161165f565b60318152602081017f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2081527f696e636f72726563742073656e64657221000000000000000000000000000000602082015290506116b3565b6020808252810161047e816116c9565b60258152602081017f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206181527f737365742e000000000000000000000000000000000000000000000000000000602082015290506116b3565b6020808252810161047e81611731565b601a8152602081017f417474656d7074696e6720656d707479207472616e736665722e00000000000081529050611561565b6020808252810161047e81611799565b6020808252810161047e81604e81527f54686973206164647265737320686173206e6f74206265656e20676976656e2060208201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60408201527f742077686974656c69737465642e000000000000000000000000000000000000606082015260800190565b606081016118708286611577565b61187d60208301856113bf565b61188a6040830184611577565b949350505050565b600067ffffffffffffffff821661047e565b61133781611892565b6020810161047e82846118a4565b60005b838110156118d65781810151838201526020016118be565b50506000910152565b60006118e9825190565b8084526020840193506119008185602086016118bb565b601f01601f19169290920192915050565b805160009060608401906119258582611577565b506020830151848203602086015261193d82826118df565b915050604083015161195260408601826113bf565b509392505050565b602080825281016111878184611911565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b63ffffffff16600063fffffffe1982016119b6576119b661196b565b5060010190565b63ffffffff8116611337565b60ff8116611337565b608081016119e082876119bd565b6119ed60208301866119bd565b81810360408301526119ff81856118df565b9050611a0e60608301846119c9565b95945050505050565b67ffffffffffffffff81166112fe565b805161047e81611a17565b600060208284031215611a4757611a47600080fd5b6111878383611a27565b60148152602081017f4661696c656420746f2073656e6420457468657200000000000000000000000081529050611561565b6020808252810161047e81611a51565b60608101611aa18286611577565b611aae6020830185611577565b61188a60408301846113bf565b8015156112fe565b805161047e81611abb565b600060208284031215611ae357611ae3600080fd5b6111878383611ac3565b6020810161047e8284611577565b6000611b05825190565b611b138185602086016118bb565b9290920192915050565b61047e8183611afb56fea264697066735822122032433aabb942ea26f27da50ceccad24f1c40e371b93725b5e6d103d9184127dc64736f6c634300081c0033",
}

// ObscuroBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroBridgeMetaData.ABI instead.
var ObscuroBridgeABI = ObscuroBridgeMetaData.ABI

// ObscuroBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroBridgeMetaData.Bin instead.
var ObscuroBridgeBin = ObscuroBridgeMetaData.Bin

// DeployObscuroBridge deploys a new Ethereum contract, binding an instance of ObscuroBridge to it.
func DeployObscuroBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ObscuroBridge, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// ObscuroBridge is an auto generated Go binding around an Ethereum contract.
type ObscuroBridge struct {
	ObscuroBridgeCaller     // Read-only binding to the contract
	ObscuroBridgeTransactor // Write-only binding to the contract
	ObscuroBridgeFilterer   // Log filterer for contract events
}

// ObscuroBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ObscuroBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObscuroBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObscuroBridgeSession struct {
	Contract     *ObscuroBridge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObscuroBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObscuroBridgeCallerSession struct {
	Contract *ObscuroBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ObscuroBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObscuroBridgeTransactorSession struct {
	Contract     *ObscuroBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ObscuroBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ObscuroBridgeRaw struct {
	Contract *ObscuroBridge // Generic contract binding to access the raw methods on
}

// ObscuroBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObscuroBridgeCallerRaw struct {
	Contract *ObscuroBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ObscuroBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactorRaw struct {
	Contract *ObscuroBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewObscuroBridge creates a new instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridge(address common.Address, backend bind.ContractBackend) (*ObscuroBridge, error) {
	contract, err := bindObscuroBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// NewObscuroBridgeCaller creates a new read-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeCaller(address common.Address, caller bind.ContractCaller) (*ObscuroBridgeCaller, error) {
	contract, err := bindObscuroBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeCaller{contract: contract}, nil
}

// NewObscuroBridgeTransactor creates a new write-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ObscuroBridgeTransactor, error) {
	contract, err := bindObscuroBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeTransactor{contract: contract}, nil
}

// NewObscuroBridgeFilterer creates a new log filterer instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ObscuroBridgeFilterer, error) {
	contract, err := bindObscuroBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeFilterer{contract: contract}, nil
}

// bindObscuroBridge binds a generic wrapper to an already deployed contract.
func bindObscuroBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.ObscuroBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) ADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ADMINROLE(&_ObscuroBridge.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) ADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ADMINROLE(&_ObscuroBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.DEFAULTADMINROLE(&_ObscuroBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.DEFAULTADMINROLE(&_ObscuroBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) ERC20TOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "ERC20_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) ERC20TOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ERC20TOKENROLE(&_ObscuroBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) ERC20TOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ERC20TOKENROLE(&_ObscuroBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) NATIVETOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "NATIVE_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) NATIVETOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.NATIVETOKENROLE(&_ObscuroBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) NATIVETOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.NATIVETOKENROLE(&_ObscuroBridge.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroBridge.Contract.GetRoleAdmin(&_ObscuroBridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroBridge.Contract.GetRoleAdmin(&_ObscuroBridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroBridge.Contract.HasRole(&_ObscuroBridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroBridge.Contract.HasRole(&_ObscuroBridge.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroBridge.Contract.SupportsInterface(&_ObscuroBridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroBridge.Contract.SupportsInterface(&_ObscuroBridge.CallOpts, interfaceId)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) Configure(opts *bind.TransactOpts, messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "configure", messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Configure(&_ObscuroBridge.TransactOpts, messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Configure(&_ObscuroBridge.TransactOpts, messengerAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.GrantRole(&_ObscuroBridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.GrantRole(&_ObscuroBridge.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "initialize", messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Initialize(&_ObscuroBridge.TransactOpts, messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Initialize(&_ObscuroBridge.TransactOpts, messenger)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) PromoteToAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "promoteToAdmin", newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.PromoteToAdmin(&_ObscuroBridge.TransactOpts, newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.PromoteToAdmin(&_ObscuroBridge.TransactOpts, newAdmin)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RemoveToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "removeToken", asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RemoveToken(&_ObscuroBridge.TransactOpts, asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RemoveToken(&_ObscuroBridge.TransactOpts, asset)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RevokeRole(&_ObscuroBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RevokeRole(&_ObscuroBridge.TransactOpts, role, account)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendERC20(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendERC20(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, receiver)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SetRemoteBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "setRemoteBridge", bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SetRemoteBridge(&_ObscuroBridge.TransactOpts, bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SetRemoteBridge(&_ObscuroBridge.TransactOpts, bridge)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) WhitelistToken(opts *bind.TransactOpts, asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "whitelistToken", asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.WhitelistToken(&_ObscuroBridge.TransactOpts, asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.WhitelistToken(&_ObscuroBridge.TransactOpts, asset, name, symbol)
}

// ObscuroBridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ObscuroBridge contract.
type ObscuroBridgeInitializedIterator struct {
	Event *ObscuroBridgeInitialized // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeInitialized)
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
		it.Event = new(ObscuroBridgeInitialized)
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
func (it *ObscuroBridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeInitialized represents a Initialized event raised by the ObscuroBridge contract.
type ObscuroBridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*ObscuroBridgeInitializedIterator, error) {

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeInitializedIterator{contract: _ObscuroBridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeInitialized)
				if err := _ObscuroBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseInitialized(log types.Log) (*ObscuroBridgeInitialized, error) {
	event := new(ObscuroBridgeInitialized)
	if err := _ObscuroBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleAdminChangedIterator struct {
	Event *ObscuroBridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleAdminChanged)
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
		it.Event = new(ObscuroBridgeRoleAdminChanged)
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
func (it *ObscuroBridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleAdminChanged represents a RoleAdminChanged event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ObscuroBridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleAdminChangedIterator{contract: _ObscuroBridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleAdminChanged)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleAdminChanged(log types.Log) (*ObscuroBridgeRoleAdminChanged, error) {
	event := new(ObscuroBridgeRoleAdminChanged)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleGrantedIterator struct {
	Event *ObscuroBridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleGranted)
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
		it.Event = new(ObscuroBridgeRoleGranted)
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
func (it *ObscuroBridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleGranted represents a RoleGranted event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroBridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleGrantedIterator{contract: _ObscuroBridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleGranted)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleGranted(log types.Log) (*ObscuroBridgeRoleGranted, error) {
	event := new(ObscuroBridgeRoleGranted)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleRevokedIterator struct {
	Event *ObscuroBridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleRevoked)
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
		it.Event = new(ObscuroBridgeRoleRevoked)
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
func (it *ObscuroBridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleRevoked represents a RoleRevoked event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroBridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleRevokedIterator{contract: _ObscuroBridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleRevoked)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleRevoked(log types.Log) (*ObscuroBridgeRoleRevoked, error) {
	event := new(ObscuroBridgeRoleRevoked)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
