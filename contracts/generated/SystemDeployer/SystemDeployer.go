// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package SystemDeployer

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

// SystemDeployerMetaData contains all meta data concerning the SystemDeployer contract.
var SystemDeployerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"SystemContractDeployed\",\"type\":\"event\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620024a9380380620024a98339810160408190526200003491620001bf565b6200003f8162000046565b5062000307565b6000604051620000569062000169565b604051809103906000f08015801562000073573d6000803e3d6000fd5b509050600063485cc95560e01b833360405160240162000095929190620001fd565b604051602081830303815290604052906001600160e01b0319166020820180516001600160e01b03838183161783525050505090506000620000df8385846200012060201b60201c565b90507fbd64e14789a915ea657e42f2dbf0b973227708fa64b58766287637985d1ade698160405162000112919062000223565b60405180910390a150505050565b600080848484604051620001349062000177565b6200014293929190620002cb565b604051809103906000f0801580156200015f573d6000803e3d6000fd5b5095945050505050565b610e28806200035583390190565b61132c806200117d83390190565b60006001600160a01b0382165b92915050565b620001a38162000185565b8114620001af57600080fd5b50565b8051620001928162000198565b600060208284031215620001d657620001d6600080fd5b6000620001e48484620001b2565b949350505050565b620001f78162000185565b82525050565b604081016200020d8285620001ec565b6200021c6020830184620001ec565b9392505050565b604080825281016200025f81601481527f5472616e73616374696f6e73416e616c797a6572000000000000000000000000602082015260400190565b9050620001926020830184620001ec565b60005b838110156200028d57818101518382015260200162000273565b50506000910152565b6000620002a1825190565b808452602084019350620002ba81856020860162000270565b601f01601f19169290920192915050565b60608101620002db8286620001ec565b620002ea6020830185620001ec565b8181036040830152620002fe818462000296565b95945050505050565b603f80620003166000396000f3fe6080604052600080fdfea2646970667358221220ff9df1e4c29f3a37150345a049596501d6997b1caf7da1bb6cb030b1eeed4a7f64736f6c63430008140033608060405234801561001057600080fd5b50610e08806100206000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c80635f03a66111610081578063d547741f1161005b578063d547741f146101fa578063dfc6cc361461020d578063ee546fd81461022057600080fd5b80635f03a6611461019457806391d14854146101bb578063a217fddf146101f257600080fd5b806336568abe116100b257806336568abe14610147578063485cc9551461015a578063508a50f41461016d57600080fd5b806301ffc9a7146100d9578063248a9ca3146101025780632f2ff15d14610132575b600080fd5b6100ec6100e73660046108a6565b610297565b6040516100f991906108d9565b60405180910390f35b6101256101103660046108f8565b60009081526020819052604090206001015490565b6040516100f9919061091f565b610145610140366004610952565b610330565b005b610145610155366004610952565b61035b565b61014561016836600461098f565b6103ac565b6101257ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59881565b6101257f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db81565b6100ec6101c9366004610952565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b610125600081565b610145610208366004610952565b610548565b61014561021b366004610a03565b61056d565b61014561022e366004610a4b565b6001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061032a57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b60008281526020819052604090206001015461034b816106d0565b61035583836106dd565b50505050565b6001600160a01b038116331461039d576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103a78282610787565b505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156103f75750825b905060008267ffffffffffffffff1660011480156104145750303b155b905081158015610422575080155b15610459576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561048d57845468ff00000000000000001916680100000000000000001785555b6104986000886106dd565b506104c37ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598886106dd565b506104ee7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db876106dd565b50831561053f57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061053690600190610a87565b60405180910390a15b50505050505050565b600082815260208190526040902060010154610563816106d0565b6103558383610787565b7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db610597816106d0565b60008290036105db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d290610a95565b60405180910390fd5b6040517f3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece419061060b90849061091f565b60405180910390a160005b6001548110156103555760006001828154811061063557610635610ad0565b6000918252602090912001546040517fd43827cb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169150819063d43827cb9061068c9088908890600401610d76565b600060405180830381600087803b1580156106a657600080fd5b505af11580156106ba573d6000803e3d6000fd5b5050505050806106c990610d9e565b9050610616565b6106da813361080a565b50565b6000828152602081815260408083206001600160a01b038516845290915281205460ff1661077f576000838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556107373390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161032a565b50600061032a565b6000828152602081815260408083206001600160a01b038516845290915281205460ff161561077f576000838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161032a565b6000828152602081815260408083206001600160a01b038516845290915290205460ff166108685780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016105d2929190610db7565b5050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b81146106da57600080fd5b803561032a8161086c565b6000602082840312156108bb576108bb600080fd5b60006108c7848461089b565b949350505050565b8015155b82525050565b6020810161032a82846108cf565b80610890565b803561032a816108e7565b60006020828403121561090d5761090d600080fd5b60006108c784846108ed565b806108d3565b6020810161032a8284610919565b60006001600160a01b03821661032a565b6108908161092d565b803561032a8161093e565b6000806040838503121561096857610968600080fd5b600061097485856108ed565b925050602061098585828601610947565b9150509250929050565b600080604083850312156109a5576109a5600080fd5b60006109748585610947565b60008083601f8401126109c6576109c6600080fd5b50813567ffffffffffffffff8111156109e1576109e1600080fd5b6020830191508360208202830111156109fc576109fc600080fd5b9250929050565b60008060208385031215610a1957610a19600080fd5b823567ffffffffffffffff811115610a3357610a33600080fd5b610a3f858286016109b1565b92509250509250929050565b600060208284031215610a6057610a60600080fd5b60006108c78484610947565b600067ffffffffffffffff821661032a565b6108d381610a6c565b6020810161032a8284610a7e565b6020808252810161032a81601a81527f4e6f207472616e73616374696f6e7320746f20636f6e76657274000000000000602082015260400190565b634e487b7160e01b600052603260045260246000fd5b60ff8116610890565b803561032a81610ae6565b50600061032a6020830183610aef565b60ff81166108d3565b50600061032a60208301836108ed565b50600061032a6020830183610947565b6108d38161092d565b6000808335601e1936859003018112610b5757610b57600080fd5b830160208101925035905067ffffffffffffffff811115610b7a57610b7a600080fd5b368190038213156109fc576109fc600080fd5b82818337506000910152565b818352602083019250610bad828483610b8d565b50601f01601f19160190565b801515610890565b803561032a81610bb9565b50600061032a6020830183610bc1565b60006101208301610bed8380610afa565b610bf78582610b0a565b50610c056020840184610b13565b610c126020860182610919565b50610c206040840184610b13565b610c2d6040860182610919565b50610c3b6060840184610b13565b610c486060860182610919565b50610c566080840184610b23565b610c636080860182610b33565b50610c7160a0840184610b13565b610c7e60a0860182610919565b50610c8c60c0840184610b3c565b85830360c0870152610c9f838284610b99565b92505050610cb060e0840184610b23565b610cbd60e0860182610b33565b50610ccc610100840184610bcc565b610cda6101008601826108cf565b509392505050565b6000610cee8383610bdc565b9392505050565b6000823561011e1936849003018112610d1057610d10600080fd5b90910192915050565b818352602083019250600083602084028101838060005b87811015610d69578484038952610d478284610cf5565b610d518582610ce2565b94505060208201602099909901989150600101610d30565b5091979650505050505050565b602080825281016108c7818486610d19565b634e487b7160e01b600052601160045260246000fd5b600060018201610db057610db0610d88565b5060010190565b60408101610dc58285610b33565b610cee602083018461091956fea2646970667358221220735d7f7f76e88e4e593c9d0b792373fd020ef0f687d35555e63a4c6ea0dfccee64736f6c6343000814003360a06040526040516200132c3803806200132c8339810160408190526200002691620004c5565b828162000034828262000098565b505081604051620000459062000351565b6200005191906200054c565b604051809103906000f0801580156200006e573d6000803e3d6000fd5b506001600160a01b03166080526200008f6200008960805190565b620000fe565b505050620005ac565b620000a38262000167565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115620000f057620000eb8282620001e4565b505050565b620000fa62000263565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620001406000805160206200130c833981519152546001600160a01b031690565b82604051620001519291906200055c565b60405180910390a1620001648162000285565b50565b806001600160a01b03163b600003620001a05780604051634c9c8ce360e01b81526004016200019791906200054c565b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080846001600160a01b031684604051620002039190620005a0565b600060405180830381855af49150503d806000811462000240576040519150601f19603f3d011682016040523d82523d6000602084013e62000245565b606091505b50909250905062000258858383620002c9565b925050505b92915050565b3415620002835760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b038116620002b2576000604051633173bdd160e11b81526004016200019791906200054c565b806000805160206200130c833981519152620001c3565b606082620002e257620002dc8262000327565b62000320565b8151158015620002fa57506001600160a01b0384163b155b156200031d5783604051639996b31560e01b81526004016200019791906200054c565b50805b9392505050565b805115620003385780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b6106ff8062000c0d83390190565b60006001600160a01b0382166200025d565b6200037c816200035f565b81146200016457600080fd5b80516200025d8162000371565b634e487b7160e01b600052604160045260246000fd5b601f19601f83011681016001600160401b0381118282101715620003d357620003d362000395565b6040525050565b6000620003e660405190565b9050620003f48282620003ab565b919050565b60006001600160401b0382111562000415576200041562000395565b601f19601f83011660200192915050565b60005b838110156200044357818101518382015260200162000429565b50506000910152565b6000620004636200045d84620003f9565b620003da565b905082815260208101848484011115620004805762000480600080fd5b6200048d84828562000426565b509392505050565b600082601f830112620004ab57620004ab600080fd5b8151620004bd8482602086016200044c565b949350505050565b600080600060608486031215620004df57620004df600080fd5b6000620004ed868662000388565b9350506020620005008682870162000388565b604086015190935090506001600160401b03811115620005235762000523600080fd5b620005318682870162000495565b9150509250925092565b62000546816200035f565b82525050565b602081016200025d82846200053b565b604081016200056c82856200053b565b6200032060208301846200053b565b600062000586825190565b6200059681856020860162000426565b9290920192915050565b6200025d81836200057b565b608051610646620005c76000396000601001526106466000f3fe608060405261000c61000e565b005b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633036100c5576000357fffffffff00000000000000000000000000000000000000000000000000000000167f4f1ef28600000000000000000000000000000000000000000000000000000000146100bb576040517fd2b576ec00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6100c36100cd565b565b6100c36100fc565b6000806100dd36600481846103cf565b8101906100ea919061054b565b915091506100f8828261010c565b5050565b6100c3610107610167565b61019f565b610115826101c3565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561015f5761015a828261026b565b505050565b6100f86102e3565b600061019a7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156101be573d6000f35b3d6000fd5b806001600160a01b03163b60000361021257806040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815260040161020991906105b2565b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516102889190610606565b600060405180830381855af49150503d80600081146102c3576040519150601f19603f3d011682016040523d82523d6000602084013e6102c8565b606091505b50915091506102d885838361031b565b925050505b92915050565b34156100c3576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060826103305761032b8261038a565b610383565b815115801561034757506001600160a01b0384163b155b1561038057836040517f9996b31500000000000000000000000000000000000000000000000000000000815260040161020991906105b2565b50805b9392505050565b80511561039a5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50565b600080858511156103e2576103e2600080fd5b838611156103f2576103f2600080fd5b5050820193919092039150565b60006001600160a01b0382166102dd565b610419816103ff565b81146103cc57600080fd5b80356102dd81610410565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b601f19601f830116810181811067ffffffffffffffff821117156104845761048461042f565b6040525050565b600061049660405190565b90506104a2828261045e565b919050565b600067ffffffffffffffff8211156104c1576104c161042f565b601f19601f83011660200192915050565b82818337506000910152565b60006104f16104ec846104a7565b61048b565b90508281526020810184848401111561050c5761050c600080fd5b6105178482856104d2565b509392505050565b600082601f83011261053357610533600080fd5b81356105438482602086016104de565b949350505050565b6000806040838503121561056157610561600080fd5b600061056d8585610424565b925050602083013567ffffffffffffffff81111561058d5761058d600080fd5b6105998582860161051f565b9150509250929050565b6105ac816103ff565b82525050565b602081016102dd82846105a3565b60005b838110156105db5781810151838201526020016105c3565b50506000910152565b60006105ee825190565b6105fc8185602086016105c0565b9290920192915050565b6102dd81836105e456fea2646970667358221220eecd99cd63d826407d5639a8db1821a6b0aadb0feacc1d2f9d4541d896c80bd464736f6c63430008140033608060405234801561001057600080fd5b506040516106ff3803806106ff83398101604081905261002f916100f8565b806001600160a01b038116610063576000604051631e4fbdf760e01b815260040161005a9190610130565b60405180910390fd5b61006c81610073565b505061013e565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006001600160a01b0382165b92915050565b6100df816100c3565b81146100ea57600080fd5b50565b80516100d0816100d6565b60006020828403121561010d5761010d600080fd5b600061011984846100ed565b949350505050565b61012a816100c3565b82525050565b602081016100d08284610121565b6105b28061014d6000396000f3fe60806040526004361061005a5760003560e01c80639623609d116100435780639623609d146100a5578063ad3cb1cc146100b8578063f2fde38b1461010e57600080fd5b8063715018a61461005f5780638da5cb5b14610076575b600080fd5b34801561006b57600080fd5b5061007461012e565b005b34801561008257600080fd5b506000546001600160a01b031660405161009c91906102fa565b60405180910390f35b6100746100b3366004610462565b610142565b3480156100c457600080fd5b506101016040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161009c9190610523565b34801561011a57600080fd5b5061007461012936600461053b565b6101ca565b61013661022a565b6101406000610270565b565b61014a61022a565b6040517f4f1ef2860000000000000000000000000000000000000000000000000000000081526001600160a01b03841690634f1ef286903490610193908690869060040161055c565b6000604051808303818588803b1580156101ac57600080fd5b505af11580156101c0573d6000803e3d6000fd5b5050505050505050565b6101d261022a565b6001600160a01b03811661021e5760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161021591906102fa565b60405180910390fd5b61022781610270565b50565b6000546001600160a01b0316331461014057336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161021591906102fa565b600080546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b60006001600160a01b0382165b92915050565b6102f4816102d8565b82525050565b602081016102e582846102eb565b60006102e5826102d8565b61031c81610308565b811461022757600080fd5b80356102e581610313565b61031c816102d8565b80356102e581610332565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b601f19601f830116810181811067ffffffffffffffff8211171561039b5761039b610346565b6040525050565b60006103ad60405190565b90506103b98282610375565b919050565b600067ffffffffffffffff8211156103d8576103d8610346565b601f19601f83011660200192915050565b82818337506000910152565b6000610408610403846103be565b6103a2565b90508281526020810184848401111561042357610423600080fd5b61042e8482856103e9565b509392505050565b600082601f83011261044a5761044a600080fd5b813561045a8482602086016103f5565b949350505050565b60008060006060848603121561047a5761047a600080fd5b60006104868686610327565b93505060206104978682870161033b565b925050604084013567ffffffffffffffff8111156104b7576104b7600080fd5b6104c386828701610436565b9150509250925092565b60005b838110156104e85781810151838201526020016104d0565b50506000910152565b60006104fb825190565b8084526020840193506105128185602086016104cd565b601f01601f19169290920192915050565b6020808252810161053481846104f1565b9392505050565b60006020828403121561055057610550600080fd5b600061045a848461033b565b6040810161056a82856102eb565b818103602083015261045a81846104f156fea2646970667358221220b0f8a95a6e2425eadd967ffc0cf44240f936f4de811b02bbb90ac6935cf0ce6a64736f6c63430008140033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103",
}

// SystemDeployerABI is the input ABI used to generate the binding from.
// Deprecated: Use SystemDeployerMetaData.ABI instead.
var SystemDeployerABI = SystemDeployerMetaData.ABI

// SystemDeployerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SystemDeployerMetaData.Bin instead.
var SystemDeployerBin = SystemDeployerMetaData.Bin

// DeploySystemDeployer deploys a new Ethereum contract, binding an instance of SystemDeployer to it.
func DeploySystemDeployer(auth *bind.TransactOpts, backend bind.ContractBackend, eoaAdmin common.Address) (common.Address, *types.Transaction, *SystemDeployer, error) {
	parsed, err := SystemDeployerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SystemDeployerBin), backend, eoaAdmin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SystemDeployer{SystemDeployerCaller: SystemDeployerCaller{contract: contract}, SystemDeployerTransactor: SystemDeployerTransactor{contract: contract}, SystemDeployerFilterer: SystemDeployerFilterer{contract: contract}}, nil
}

// SystemDeployer is an auto generated Go binding around an Ethereum contract.
type SystemDeployer struct {
	SystemDeployerCaller     // Read-only binding to the contract
	SystemDeployerTransactor // Write-only binding to the contract
	SystemDeployerFilterer   // Log filterer for contract events
}

// SystemDeployerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SystemDeployerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemDeployerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SystemDeployerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemDeployerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SystemDeployerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemDeployerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SystemDeployerSession struct {
	Contract     *SystemDeployer   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SystemDeployerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SystemDeployerCallerSession struct {
	Contract *SystemDeployerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SystemDeployerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SystemDeployerTransactorSession struct {
	Contract     *SystemDeployerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SystemDeployerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SystemDeployerRaw struct {
	Contract *SystemDeployer // Generic contract binding to access the raw methods on
}

// SystemDeployerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SystemDeployerCallerRaw struct {
	Contract *SystemDeployerCaller // Generic read-only contract binding to access the raw methods on
}

// SystemDeployerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SystemDeployerTransactorRaw struct {
	Contract *SystemDeployerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSystemDeployer creates a new instance of SystemDeployer, bound to a specific deployed contract.
func NewSystemDeployer(address common.Address, backend bind.ContractBackend) (*SystemDeployer, error) {
	contract, err := bindSystemDeployer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SystemDeployer{SystemDeployerCaller: SystemDeployerCaller{contract: contract}, SystemDeployerTransactor: SystemDeployerTransactor{contract: contract}, SystemDeployerFilterer: SystemDeployerFilterer{contract: contract}}, nil
}

// NewSystemDeployerCaller creates a new read-only instance of SystemDeployer, bound to a specific deployed contract.
func NewSystemDeployerCaller(address common.Address, caller bind.ContractCaller) (*SystemDeployerCaller, error) {
	contract, err := bindSystemDeployer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SystemDeployerCaller{contract: contract}, nil
}

// NewSystemDeployerTransactor creates a new write-only instance of SystemDeployer, bound to a specific deployed contract.
func NewSystemDeployerTransactor(address common.Address, transactor bind.ContractTransactor) (*SystemDeployerTransactor, error) {
	contract, err := bindSystemDeployer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SystemDeployerTransactor{contract: contract}, nil
}

// NewSystemDeployerFilterer creates a new log filterer instance of SystemDeployer, bound to a specific deployed contract.
func NewSystemDeployerFilterer(address common.Address, filterer bind.ContractFilterer) (*SystemDeployerFilterer, error) {
	contract, err := bindSystemDeployer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SystemDeployerFilterer{contract: contract}, nil
}

// bindSystemDeployer binds a generic wrapper to an already deployed contract.
func bindSystemDeployer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SystemDeployerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SystemDeployer *SystemDeployerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SystemDeployer.Contract.SystemDeployerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SystemDeployer *SystemDeployerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SystemDeployer.Contract.SystemDeployerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SystemDeployer *SystemDeployerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SystemDeployer.Contract.SystemDeployerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SystemDeployer *SystemDeployerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SystemDeployer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SystemDeployer *SystemDeployerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SystemDeployer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SystemDeployer *SystemDeployerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SystemDeployer.Contract.contract.Transact(opts, method, params...)
}

// SystemDeployerSystemContractDeployedIterator is returned from FilterSystemContractDeployed and is used to iterate over the raw logs and unpacked data for SystemContractDeployed events raised by the SystemDeployer contract.
type SystemDeployerSystemContractDeployedIterator struct {
	Event *SystemDeployerSystemContractDeployed // Event containing the contract specifics and raw log

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
func (it *SystemDeployerSystemContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SystemDeployerSystemContractDeployed)
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
		it.Event = new(SystemDeployerSystemContractDeployed)
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
func (it *SystemDeployerSystemContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SystemDeployerSystemContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SystemDeployerSystemContractDeployed represents a SystemContractDeployed event raised by the SystemDeployer contract.
type SystemDeployerSystemContractDeployed struct {
	Name            string
	ContractAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterSystemContractDeployed is a free log retrieval operation binding the contract event 0xbd64e14789a915ea657e42f2dbf0b973227708fa64b58766287637985d1ade69.
//
// Solidity: event SystemContractDeployed(string name, address contractAddress)
func (_SystemDeployer *SystemDeployerFilterer) FilterSystemContractDeployed(opts *bind.FilterOpts) (*SystemDeployerSystemContractDeployedIterator, error) {

	logs, sub, err := _SystemDeployer.contract.FilterLogs(opts, "SystemContractDeployed")
	if err != nil {
		return nil, err
	}
	return &SystemDeployerSystemContractDeployedIterator{contract: _SystemDeployer.contract, event: "SystemContractDeployed", logs: logs, sub: sub}, nil
}

// WatchSystemContractDeployed is a free log subscription operation binding the contract event 0xbd64e14789a915ea657e42f2dbf0b973227708fa64b58766287637985d1ade69.
//
// Solidity: event SystemContractDeployed(string name, address contractAddress)
func (_SystemDeployer *SystemDeployerFilterer) WatchSystemContractDeployed(opts *bind.WatchOpts, sink chan<- *SystemDeployerSystemContractDeployed) (event.Subscription, error) {

	logs, sub, err := _SystemDeployer.contract.WatchLogs(opts, "SystemContractDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SystemDeployerSystemContractDeployed)
				if err := _SystemDeployer.contract.UnpackLog(event, "SystemContractDeployed", log); err != nil {
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

// ParseSystemContractDeployed is a log parse operation binding the contract event 0xbd64e14789a915ea657e42f2dbf0b973227708fa64b58766287637985d1ade69.
//
// Solidity: event SystemContractDeployed(string name, address contractAddress)
func (_SystemDeployer *SystemDeployerFilterer) ParseSystemContractDeployed(log types.Log) (*SystemDeployerSystemContractDeployed, error) {
	event := new(SystemDeployerSystemContractDeployed)
	if err := _SystemDeployer.contract.UnpackLog(event, "SystemContractDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
