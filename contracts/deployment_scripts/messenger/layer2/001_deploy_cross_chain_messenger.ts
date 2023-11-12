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
        getNamedAccounts
    } = hre;

    // Get the prefunded L2 deployer account to use for deploying.
    const {deployer} = await getNamedAccounts();

    console.log(`Deployer acc ${deployer}`);

    // TODO: Remove hardcoded L2 message bus address when properly exposed.
    const busAddress = hre.ethers.utils.getAddress("0x526c84529b2b8c11f57d93d3f5537aca3aecef9b")

    console.log(`Beginning deploy of cross chain messenger`);

    // Deploy the L2 Cross chain messenger and use the L2 bus for validation
    await deployments.deploy('CrossChainMessengerL2', {
        from: deployer,
        args: [ ],
        log: true,
        contract: 'CrossChainMessenger'
    });
};

export default func;
func.tags = ['CrossChainMessengerL2', 'CrossChainMessenger_deployL2'];
func.dependencies = ['CrossChainMessenger']; //TODO: Remove HPERC20, this is only to have matching addresses.
