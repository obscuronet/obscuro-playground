import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const UUPSTestV1 = await hre.ethers.getContractFactory("UUPSTestV1");
    const uupsTestv1 = await hre.upgrades.deployProxy(UUPSTestV1, [], {
      initializer: 'initialize',
      kind: 'uups',
      txOverrides: {
      }
     })
  
    await uupsTestv1.waitForDeployment();
    const address = await uupsTestv1.getAddress()
    console.log("Proxy deployed at", address);
  
    const UUPSTestV2 = await hre.ethers.getContractFactory("UUPSTestV2");
    const uuupsTestV2 = await hre.upgrades.upgradeProxy(address, UUPSTestV2);
    console.log("UUPS Test V2 deployed to:", await uuupsTestV2.getAddress());
  
};


export default func;
func.tags = ['GasDebug'];
