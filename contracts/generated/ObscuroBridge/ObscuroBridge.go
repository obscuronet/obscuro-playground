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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b191690553480156200002157600080fd5b5060405162001bbb38038062001bbb8339810160408190526200004491620001fc565b600080546001600160a01b0319166001600160a01b038316908117909155604080516350d113fd60e11b8152905183929163a1a227fa916004808301926020929190829003018186803b1580156200009b57600080fd5b505afa158015620000b0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000d69190620001fc565b600180546001600160a01b0319166001600160a01b039290921691909117905550620001237fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217753362000157565b620001507fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad2110578600062000157565b506200022e565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff16620001f85760008281526002602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620001b73390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b6000602082840312156200020f57600080fd5b81516001600160a01b03811681146200022757600080fd5b9392505050565b61197d806200023e6000396000f3fe60806040526004361061010e5760003560e01c80635fa7b584116100a557806393b3744211610074578063a381c8e211610059578063a381c8e21461033e578063d547741f1461035e578063e4c3ebc71461037e57600080fd5b806393b3744214610309578063a217fddf1461032957600080fd5b80635fa7b5841461024f57806375b238fc1461026f57806383bece4d146102a357806391d14854146102c357600080fd5b80632f2ff15d116100e15780632f2ff15d146101bb57806336568abe146101db578063498d82ab146101fb5780635d8729701461021b57600080fd5b806301ffc9a71461011357806316ce8149146101485780631888d7121461016a578063248a9ca31461017d575b600080fd5b34801561011f57600080fd5b5061013361012e36600461147f565b6103b2565b60405190151581526020015b60405180910390f35b34801561015457600080fd5b506101686101633660046114c1565b61041b565b005b6101686101783660046114c1565b610481565b34801561018957600080fd5b506101ad6101983660046114de565b60009081526002602052604090206001015490565b60405190815260200161013f565b3480156101c757600080fd5b506101686101d63660046114f7565b610532565b3480156101e757600080fd5b506101686101f63660046114f7565b61055d565b34801561020757600080fd5b50610168610216366004611570565b6105e5565b34801561022757600080fd5b506101ad7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b34801561025b57600080fd5b5061016861026a3660046114c1565b6106b3565b34801561027b57600080fd5b506101ad7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b3480156102af57600080fd5b506101686102be3660046115f3565b610708565b3480156102cf57600080fd5b506101336102de3660046114f7565b60009182526002602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561031557600080fd5b506101686103243660046114c1565b61092b565b34801561033557600080fd5b506101ad600081565b34801561034a57600080fd5b506101686103593660046115f3565b610980565b34801561036a57600080fd5b506101686103793660046114f7565b610b2d565b34801561038a57600080fd5b506101ad7fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b60006001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061041557507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756104468133610b53565b50600380547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b600034116104d65760405162461bcd60e51b815260206004820152600f60248201527f456d707479207472616e736665722e000000000000000000000000000000000060448201526064015b60405180910390fd5b604080518082018252348082526001600160a01b03848116602093840190815284519384019290925290518116828401528251808303840181526060909201909252600354909161052e91168260005b600080610bd3565b5050565b60008281526002602052604090206001015461054e8133610b53565b6105588383610cee565b505050565b6001600160a01b03811633146105db5760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c66000000000000000000000000000000000060648201526084016104cd565b61052e8282610d90565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106108133610b53565b61063a7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a87610cee565b600063458ffd6360e01b878787878760405160240161065d95949392919061165e565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526003549091506106aa906001600160a01b0316826002610526565b50505050505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106de8133610b53565b61052e7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a83610d90565b6003546000546001600160a01b039182169116331461078f5760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e676572210000000000000000000000000000000060648201526084016104cd565b806001600160a01b03166107a1610e13565b6001600160a01b03161461081d5760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e6465722100000000000000000000000000000060648201526084016104cd565b6001600160a01b03841660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff161561086d57610868848484610e9f565b610925565b6001600160a01b03841660009081527f13ad2d85210d477fe1a6e25654c8250308cf29b050a4bf0b039d70467486712c602052604090205460ff16156108b7576108688383610eaa565b60405162461bcd60e51b815260206004820152602560248201527f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206160448201527f737365742e00000000000000000000000000000000000000000000000000000060648201526084016104cd565b50505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756109568133610b53565b61052e7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610cee565b600082116109d05760405162461bcd60e51b815260206004820152601a60248201527f417474656d7074696e6720656d707479207472616e736665722e00000000000060448201526064016104cd565b6001600160a01b03831660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16610aa35760405162461bcd60e51b815260206004820152604e60248201527f54686973206164647265737320686173206e6f74206265656e20676976656e2060448201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60648201527f742077686974656c69737465642e000000000000000000000000000000000000608482015260a4016104cd565b610aaf83333085610f4d565b604080516001600160a01b038581166024830152604482018590528381166064808401919091528351808403909101815260849092019092526020810180516001600160e01b03167f83bece4d0000000000000000000000000000000000000000000000000000000017905260035490916109259116826001610526565b600082815260026020526040902060010154610b498133610b53565b6105588383610d90565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff1661052e57610b91816001600160a01b03166014610fd1565b610b9c836020610fd1565b604051602001610bad9291906116cc565b60408051601f198184030181529082905262461bcd60e51b82526104cd91600401611779565b60006040518060600160405280876001600160a01b0316815260200186815260200184815250604051602001610c09919061178c565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa917401000000000000000000000000000000000000000090910463ffffffff16906014610c5c836117e7565b91906101000a81548163ffffffff021916908363ffffffff1602179055508684866040518563ffffffff1660e01b8152600401610c9c949392919061180b565b602060405180830381600087803b158015610cb657600080fd5b505af1158015610cca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106aa9190611848565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff1661052e5760008281526002602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610d4c3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff161561052e5760008281526002602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b60008060009054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b815260040160206040518083038186803b158015610e6257600080fd5b505afa158015610e76573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e9a9190611872565b905090565b6105588382846111b9565b6000816001600160a01b03168360405160006040518083038185875af1925050503d8060008114610ef7576040519150601f19603f3d011682016040523d82523d6000602084013e610efc565b606091505b50509050806105585760405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e6420457468657200000000000000000000000060448201526064016104cd565b6040516001600160a01b03808516602483015283166044820152606481018290526109259085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152611202565b60606000610fe083600261188f565b610feb9060026118ae565b67ffffffffffffffff811115611003576110036118c6565b6040519080825280601f01601f19166020018201604052801561102d576020820181803683370190505b5090507f300000000000000000000000000000000000000000000000000000000000000081600081518110611064576110646118dc565b60200101906001600160f81b031916908160001a9053507f7800000000000000000000000000000000000000000000000000000000000000816001815181106110af576110af6118dc565b60200101906001600160f81b031916908160001a90535060006110d384600261188f565b6110de9060016118ae565b90505b6001811115611163577f303132333435363738396162636465660000000000000000000000000000000085600f166010811061111f5761111f6118dc565b1a60f81b828281518110611135576111356118dc565b60200101906001600160f81b031916908160001a90535060049490941c9361115c816118f2565b90506110e1565b5083156111b25760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016104cd565b9392505050565b6040516001600160a01b0383166024820152604481018290526105589084907fa9059cbb0000000000000000000000000000000000000000000000000000000090606401610f9a565b6000611257826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166112e79092919063ffffffff16565b80519091501561055857808060200190518101906112759190611909565b6105585760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016104cd565b60606112f684846000856112fe565b949350505050565b6060824710156113765760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016104cd565b6001600160a01b0385163b6113cd5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016104cd565b600080866001600160a01b031685876040516113e9919061192b565b60006040518083038185875af1925050503d8060008114611426576040519150601f19603f3d011682016040523d82523d6000602084013e61142b565b606091505b509150915061143b828286611446565b979650505050505050565b606083156114555750816111b2565b8251156114655782518084602001fd5b8160405162461bcd60e51b81526004016104cd9190611779565b60006020828403121561149157600080fd5b81356001600160e01b0319811681146111b257600080fd5b6001600160a01b03811681146114be57600080fd5b50565b6000602082840312156114d357600080fd5b81356111b2816114a9565b6000602082840312156114f057600080fd5b5035919050565b6000806040838503121561150a57600080fd5b82359150602083013561151c816114a9565b809150509250929050565b60008083601f84011261153957600080fd5b50813567ffffffffffffffff81111561155157600080fd5b60208301915083602082850101111561156957600080fd5b9250929050565b60008060008060006060868803121561158857600080fd5b8535611593816114a9565b9450602086013567ffffffffffffffff808211156115b057600080fd5b6115bc89838a01611527565b909650945060408801359150808211156115d557600080fd5b506115e288828901611527565b969995985093965092949392505050565b60008060006060848603121561160857600080fd5b8335611613816114a9565b925060208401359150604084013561162a816114a9565b809150509250925092565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0386168152606060208201526000611681606083018688611635565b8281036040840152611694818587611635565b98975050505050505050565b60005b838110156116bb5781810151838201526020016116a3565b838111156109255750506000910152565b7f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008152600083516117048160178501602088016116a0565b7f206973206d697373696e6720726f6c652000000000000000000000000000000060179184019182015283516117418160288401602088016116a0565b01602801949350505050565b600081518084526117658160208601602086016116a0565b601f01601f19169290920160200192915050565b6020815260006111b2602083018461174d565b602081526001600160a01b03825116602082015260006020830151606060408401526117bb608084018261174d565b9050604084015160608401528091505092915050565b634e487b7160e01b600052601160045260246000fd5b600063ffffffff80831681811415611801576118016117d1565b6001019392505050565b600063ffffffff808716835280861660208401525060806040830152611834608083018561174d565b905060ff8316606083015295945050505050565b60006020828403121561185a57600080fd5b815167ffffffffffffffff811681146111b257600080fd5b60006020828403121561188457600080fd5b81516111b2816114a9565b60008160001904831182151516156118a9576118a96117d1565b500290565b600082198211156118c1576118c16117d1565b500190565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b600081611901576119016117d1565b506000190190565b60006020828403121561191b57600080fd5b815180151581146111b257600080fd5b6000825161193d8184602087016116a0565b919091019291505056fea2646970667358221220519ff8f34202b5c162e1f65b3da902137fc038957a2647d8b23853505847c0de64736f6c63430008090033",
}

// ObscuroBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroBridgeMetaData.ABI instead.
var ObscuroBridgeABI = ObscuroBridgeMetaData.ABI

// ObscuroBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroBridgeMetaData.Bin instead.
var ObscuroBridgeBin = ObscuroBridgeMetaData.Bin

// DeployObscuroBridge deploys a new Ethereum contract, binding an instance of ObscuroBridge to it.
func DeployObscuroBridge(auth *bind.TransactOpts, backend bind.ContractBackend, messenger common.Address) (common.Address, *types.Transaction, *ObscuroBridge, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroBridgeBin), backend, messenger)
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
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, account)
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
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendERC20(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
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
