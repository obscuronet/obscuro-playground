import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script deploys the L2 side of the cross chain messenger.
    It requires the address of the message bus on the L2 side, which is generated by the enclave
    on genesis. 
*/


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts,
        companionNetworks,
    } = hre;
    // Use the contract addresses from the management contract deployment.
    const mgmtContractAddress = process.env.MGMT_CONTRACT_ADDRESS!!

    // Get the prefunded L2 deployer account to use for deploying.
    const {deployer} = await getNamedAccounts();

    console.log(`Script: 001_deploy_cross_chain_messenger.ts - address used: ${deployer}`);

    // TODO: Remove hardcoded L2 message bus address when properly exposed.
    const messageBusAddress = hre.ethers.getAddress("0x526c84529b2b8c11f57d93d3f5537aca3aecef9b");
    // Deploy the L2 Cross chain messenger and use the L2 bus for validation
    const crossChainDeployment = await deployments.deploy('CrossChainMessenger', {
        from: deployer,
        log: true,
        proxy: {
            proxyContract: "OpenZeppelinTransparentProxy",
            execute: {
                init: {
                    methodName: "initialize",
                    args: [ messageBusAddress ]
                }
            }
        }
    });
    // get L1 management contract and write the cross chain messenger address to it
    const mgmtContract = (await hre.ethers.getContractFactory('ManagementContract')).attach(mgmtContractAddress);
    const tx = await mgmtContract.getFunction("SetImportantContractAddress").send("L2CrossChainMessenger", crossChainDeployment.address);
    const receipt = await tx.wait();
    if (receipt!!.status !== 1) {
        console.log("Failed to set L2CrossChainMessenger in management contract");
    }
    console.log(`L2CrossChainMessenger=${crossChainDeployment.address}`);
};

export default func;
func.tags = ['CrossChainMessenger', 'CrossChainMessenger_deploy'];
