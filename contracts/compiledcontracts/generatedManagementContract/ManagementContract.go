// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

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
)

// GeneratedManagementContractMetaData contains all meta data concerning the GeneratedManagementContract contract.
var GeneratedManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"AddHostAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506122fa806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063324ff86614610067578063597a972314610085578063981214ba146100a1578063c719bf50146100bd578063e0fd84bd146100d9578063e34fbfc8146100f5575b600080fd5b61006f610111565b60405161007c9190610ebf565b60405180910390f35b61009f600480360381019061009a919061102a565b6101ea565b005b6100bb60048036038101906100b69190611172565b610222565b005b6100d760048036038101906100d29190611271565b6103bf565b005b6100f360048036038101906100ee9190611393565b610451565b005b61010f600480360381019061010a919061142d565b61058e565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156101e1578382906000526020600020018054610154906114a9565b80601f0160208091040260200160405190810160405280929190818152602001828054610180906114a9565b80156101cd5780601f106101a2576101008083540402835291602001916101cd565b820191906000526020600020905b8154815290600101906020018083116101b057829003601f168201915b505050505081526020019060010190610135565b50505050905090565b60038190806001815401808255809150506001900390600052602060002001600090919091909150908161021e9190611686565b5050565b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061027d57600080fd5b60006102ab868685604051602001610297939291906117e7565b6040516020818303038152906040526105e1565b905060006102b9828661061c565b90508673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146102f388610643565b6102fc83610643565b60405160200161030d9291906119b2565b6040516020818303038152906040529061035d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103549190611a41565b60405180910390fd5b506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050505050505050565b600460009054906101000a900460ff16156103d957600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550505050565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166104a757600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826105dc929190611a6e565b505050565b60006105ed8251610806565b826040516020016105ff929190611b8a565b604051602081830303815290604052805190602001209050919050565b600080600061062b8585610966565b91509150610638816109e7565b819250505092915050565b60606000602867ffffffffffffffff81111561066257610661610eff565b5b6040519080825280601f01601f1916602001820160405280156106945781602001600182028036833780820191505090505b50905060005b60148110156107fc5760008160136106b29190611be8565b60086106be9190611c1c565b60026106ca9190611da9565b8573ffffffffffffffffffffffffffffffffffffffff166106eb9190611e23565b60f81b9050600060108260f81c6107029190611e61565b60f81b905060008160f81c60106107199190611e92565b8360f81c6107279190611ecd565b60f81b905061073582610bb3565b858560026107439190611c1c565b8151811061075457610753611f01565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061078c81610bb3565b85600186600261079c9190611c1c565b6107a69190611f30565b815181106107b7576107b6611f01565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050505080806107f490611f86565b91505061069a565b5080915050919050565b60606000820361084d576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610961565b600082905060005b6000821461087f57808061086890611f86565b915050600a826108789190611e23565b9150610855565b60008167ffffffffffffffff81111561089b5761089a610eff565b5b6040519080825280601f01601f1916602001820160405280156108cd5781602001600182028036833780820191505090505b5090505b6000851461095a576001826108e69190611be8565b9150600a856108f59190611fce565b60306109019190611f30565b60f81b81838151811061091757610916611f01565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a856109539190611e23565b94506108d1565b8093505050505b919050565b60008060418351036109a75760008060006020860151925060408601519150606086015160001a905061099b87828585610bf9565b945094505050506109e0565b60408351036109d75760008060208501519150604085015190506109cc868383610d05565b9350935050506109e0565b60006002915091505b9250929050565b600060048111156109fb576109fa611fff565b5b816004811115610a0e57610a0d611fff565b5b0315610bb05760016004811115610a2857610a27611fff565b5b816004811115610a3b57610a3a611fff565b5b03610a7b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a729061207a565b60405180910390fd5b60026004811115610a8f57610a8e611fff565b5b816004811115610aa257610aa1611fff565b5b03610ae2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad9906120e6565b60405180910390fd5b60036004811115610af657610af5611fff565b5b816004811115610b0957610b08611fff565b5b03610b49576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b4090612178565b60405180910390fd5b600480811115610b5c57610b5b611fff565b5b816004811115610b6f57610b6e611fff565b5b03610baf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ba69061220a565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610bde5760308260f81c610bd4919061222a565b60f81b9050610bf4565b60578260f81c610bee919061222a565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610c34576000600391509150610cfc565b601b8560ff1614158015610c4c5750601c8560ff1614155b15610c5e576000600491509150610cfc565b600060018787878760405160008152602001604052604051610c83949392919061227f565b6020604051602081039080840390855afa158015610ca5573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610cf357600060019250925050610cfc565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610d489190611f30565b9050610d5687828885610bf9565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610dca578082015181840152602081019050610daf565b83811115610dd9576000848401525b50505050565b6000601f19601f8301169050919050565b6000610dfb82610d90565b610e058185610d9b565b9350610e15818560208601610dac565b610e1e81610ddf565b840191505092915050565b6000610e358383610df0565b905092915050565b6000602082019050919050565b6000610e5582610d64565b610e5f8185610d6f565b935083602082028501610e7185610d80565b8060005b85811015610ead5784840389528151610e8e8582610e29565b9450610e9983610e3d565b925060208a01995050600181019050610e75565b50829750879550505050505092915050565b60006020820190508181036000830152610ed98184610e4a565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610f3782610ddf565b810181811067ffffffffffffffff82111715610f5657610f55610eff565b5b80604052505050565b6000610f69610ee1565b9050610f758282610f2e565b919050565b600067ffffffffffffffff821115610f9557610f94610eff565b5b610f9e82610ddf565b9050602081019050919050565b82818337600083830152505050565b6000610fcd610fc884610f7a565b610f5f565b905082815260208101848484011115610fe957610fe8610efa565b5b610ff4848285610fab565b509392505050565b600082601f83011261101157611010610ef5565b5b8135611021848260208601610fba565b91505092915050565b6000602082840312156110405761103f610eeb565b5b600082013567ffffffffffffffff81111561105e5761105d610ef0565b5b61106a84828501610ffc565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061109e82611073565b9050919050565b6110ae81611093565b81146110b957600080fd5b50565b6000813590506110cb816110a5565b92915050565b600067ffffffffffffffff8211156110ec576110eb610eff565b5b6110f582610ddf565b9050602081019050919050565b6000611115611110846110d1565b610f5f565b90508281526020810184848401111561113157611130610efa565b5b61113c848285610fab565b509392505050565b600082601f83011261115957611158610ef5565b5b8135611169848260208601611102565b91505092915050565b6000806000806080858703121561118c5761118b610eeb565b5b600061119a878288016110bc565b94505060206111ab878288016110bc565b935050604085013567ffffffffffffffff8111156111cc576111cb610ef0565b5b6111d887828801611144565b925050606085013567ffffffffffffffff8111156111f9576111f8610ef0565b5b61120587828801611144565b91505092959194509250565b600080fd5b600080fd5b60008083601f84011261123157611230610ef5565b5b8235905067ffffffffffffffff81111561124e5761124d611211565b5b60208301915083600182028301111561126a57611269611216565b5b9250929050565b60008060006040848603121561128a57611289610eeb565b5b6000611298868287016110bc565b935050602084013567ffffffffffffffff8111156112b9576112b8610ef0565b5b6112c58682870161121b565b92509250509250925092565b6000819050919050565b6112e4816112d1565b81146112ef57600080fd5b50565b600081359050611301816112db565b92915050565b6000819050919050565b61131a81611307565b811461132557600080fd5b50565b60008135905061133781611311565b92915050565b60008083601f84011261135357611352610ef5565b5b8235905067ffffffffffffffff8111156113705761136f611211565b5b60208301915083600182028301111561138c5761138b611216565b5b9250929050565b60008060008060008060a087890312156113b0576113af610eeb565b5b60006113be89828a016112f2565b96505060206113cf89828a016110bc565b95505060406113e089828a016112f2565b94505060606113f189828a01611328565b935050608087013567ffffffffffffffff81111561141257611411610ef0565b5b61141e89828a0161133d565b92509250509295509295509295565b6000806020838503121561144457611443610eeb565b5b600083013567ffffffffffffffff81111561146257611461610ef0565b5b61146e8582860161133d565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806114c157607f821691505b6020821081036114d4576114d361147a565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261153c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826114ff565b61154686836114ff565b95508019841693508086168417925050509392505050565b6000819050919050565b600061158361157e61157984611307565b61155e565b611307565b9050919050565b6000819050919050565b61159d83611568565b6115b16115a98261158a565b84845461150c565b825550505050565b600090565b6115c66115b9565b6115d1818484611594565b505050565b5b818110156115f5576115ea6000826115be565b6001810190506115d7565b5050565b601f82111561163a5761160b816114da565b611614846114ef565b81016020851015611623578190505b61163761162f856114ef565b8301826115d6565b50505b505050565b600082821c905092915050565b600061165d6000198460080261163f565b1980831691505092915050565b6000611676838361164c565b9150826002028217905092915050565b61168f82610d90565b67ffffffffffffffff8111156116a8576116a7610eff565b5b6116b282546114a9565b6116bd8282856115f9565b600060209050601f8311600181146116f057600084156116de578287015190505b6116e8858261166a565b865550611750565b601f1984166116fe866114da565b60005b8281101561172657848901518255600182019150602085019450602081019050611701565b86831015611743578489015161173f601f89168261164c565b8355505b6001600288020188555050505b505050505050565b60008160601b9050919050565b600061177082611758565b9050919050565b600061178282611765565b9050919050565b61179a61179582611093565b611777565b82525050565b600081519050919050565b600081905092915050565b60006117c1826117a0565b6117cb81856117ab565b93506117db818560208601610dac565b80840191505092915050565b60006117f38286611789565b6014820191506118038285611789565b60148201915061181382846117b6565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611887602d83611820565b91506118928261182b565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b60006118f9602483611820565b91506119048261189d565b602482019050919050565b600061191a82610d90565b6119248185611820565b9350611934818560208601610dac565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b600061199c602483611820565b91506119a782611940565b602482019050919050565b60006119bd8261187a565b91506119c8826118ec565b91506119d4828561190f565b91506119df8261198f565b91506119eb828461190f565b91508190509392505050565b600082825260208201905092915050565b6000611a1382610d90565b611a1d81856119f7565b9350611a2d818560208601610dac565b611a3681610ddf565b840191505092915050565b60006020820190508181036000830152611a5b8184611a08565b905092915050565b600082905092915050565b611a788383611a63565b67ffffffffffffffff811115611a9157611a90610eff565b5b611a9b82546114a9565b611aa68282856115f9565b6000601f831160018114611ad55760008415611ac3578287013590505b611acd858261166a565b865550611b35565b601f198416611ae3866114da565b60005b82811015611b0b57848901358255600182019150602085019450602081019050611ae6565b86831015611b285784890135611b24601f89168261164c565b8355505b6001600288020188555050505b50505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611b74601a83611820565b9150611b7f82611b3e565b601a82019050919050565b6000611b9582611b67565b9150611ba1828561190f565b9150611bad82846117b6565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611bf382611307565b9150611bfe83611307565b925082821015611c1157611c10611bb9565b5b828203905092915050565b6000611c2782611307565b9150611c3283611307565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611c6b57611c6a611bb9565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b6001851115611ccd57808604811115611ca957611ca8611bb9565b5b6001851615611cb85780820291505b8081029050611cc685611c76565b9450611c8d565b94509492505050565b600082611ce65760019050611da2565b81611cf45760009050611da2565b8160018114611d0a5760028114611d1457611d43565b6001915050611da2565b60ff841115611d2657611d25611bb9565b5b8360020a915084821115611d3d57611d3c611bb9565b5b50611da2565b5060208310610133831016604e8410600b8410161715611d785782820a905083811115611d7357611d72611bb9565b5b611da2565b611d858484846001611c83565b92509050818404811115611d9c57611d9b611bb9565b5b81810290505b9392505050565b6000611db482611307565b9150611dbf83611307565b9250611dec7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484611cd6565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611e2e82611307565b9150611e3983611307565b925082611e4957611e48611df4565b5b828204905092915050565b600060ff82169050919050565b6000611e6c82611e54565b9150611e7783611e54565b925082611e8757611e86611df4565b5b828204905092915050565b6000611e9d82611e54565b9150611ea883611e54565b92508160ff0483118215151615611ec257611ec1611bb9565b5b828202905092915050565b6000611ed882611e54565b9150611ee383611e54565b925082821015611ef657611ef5611bb9565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000611f3b82611307565b9150611f4683611307565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611f7b57611f7a611bb9565b5b828201905092915050565b6000611f9182611307565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611fc357611fc2611bb9565b5b600182019050919050565b6000611fd982611307565b9150611fe483611307565b925082611ff457611ff3611df4565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b60006120646018836119f7565b915061206f8261202e565b602082019050919050565b6000602082019050818103600083015261209381612057565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b60006120d0601f836119f7565b91506120db8261209a565b602082019050919050565b600060208201905081810360008301526120ff816120c3565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006121626022836119f7565b915061216d82612106565b604082019050919050565b6000602082019050818103600083015261219181612155565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006121f46022836119f7565b91506121ff82612198565b604082019050919050565b60006020820190508181036000830152612223816121e7565b9050919050565b600061223582611e54565b915061224083611e54565b92508260ff0382111561225657612255611bb9565b5b828201905092915050565b61226a816112d1565b82525050565b61227981611e54565b82525050565b60006080820190506122946000830187612261565b6122a16020830186612270565b6122ae6040830185612261565b6122bb6060830184612261565b9594505050505056fea2646970667358221220a77c3bbcee09aaa636b57ff6ac3cbfbb9fb8840610db2624d9a069c2ccbacc9c64736f6c634300080f0033",
}

// GeneratedManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GeneratedManagementContractMetaData.ABI instead.
var GeneratedManagementContractABI = GeneratedManagementContractMetaData.ABI

// GeneratedManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GeneratedManagementContractMetaData.Bin instead.
var GeneratedManagementContractBin = GeneratedManagementContractMetaData.Bin

// DeployGeneratedManagementContract deploys a new Ethereum contract, binding an instance of GeneratedManagementContract to it.
func DeployGeneratedManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GeneratedManagementContract, error) {
	parsed, err := GeneratedManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GeneratedManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// GeneratedManagementContract is an auto generated Go binding around an Ethereum contract.
type GeneratedManagementContract struct {
	GeneratedManagementContractCaller     // Read-only binding to the contract
	GeneratedManagementContractTransactor // Write-only binding to the contract
	GeneratedManagementContractFilterer   // Log filterer for contract events
}

// GeneratedManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedManagementContractSession struct {
	Contract     *GeneratedManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GeneratedManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedManagementContractCallerSession struct {
	Contract *GeneratedManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// GeneratedManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedManagementContractTransactorSession struct {
	Contract     *GeneratedManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// GeneratedManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedManagementContractRaw struct {
	Contract *GeneratedManagementContract // Generic contract binding to access the raw methods on
}

// GeneratedManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCallerRaw struct {
	Contract *GeneratedManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactorRaw struct {
	Contract *GeneratedManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeneratedManagementContract creates a new instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContract(address common.Address, backend bind.ContractBackend) (*GeneratedManagementContract, error) {
	contract, err := bindGeneratedManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// NewGeneratedManagementContractCaller creates a new read-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractCaller(address common.Address, caller bind.ContractCaller) (*GeneratedManagementContractCaller, error) {
	contract, err := bindGeneratedManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractCaller{contract: contract}, nil
}

// NewGeneratedManagementContractTransactor creates a new write-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedManagementContractTransactor, error) {
	contract, err := bindGeneratedManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractTransactor{contract: contract}, nil
}

// NewGeneratedManagementContractFilterer creates a new log filterer instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedManagementContractFilterer, error) {
	contract, err := bindGeneratedManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractFilterer{contract: contract}, nil
}

// bindGeneratedManagementContract binds a generic wrapper to an already deployed contract.
func bindGeneratedManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transact(opts, method, params...)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCaller) GetHostAddresses(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "GetHostAddresses")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddHostAddress(opts *bind.TransactOpts, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddHostAddress", hostAddress)
}

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddHostAddress(hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddHostAddress(&_GeneratedManagementContract.TransactOpts, hostAddress)
}

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddHostAddress(hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddHostAddress(&_GeneratedManagementContract.TransactOpts, hostAddress)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}
