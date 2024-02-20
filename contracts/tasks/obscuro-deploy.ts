import { task } from "hardhat/config";
import 'hardhat/types/config';

import {on, exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";

task("obscuro:deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    
    if (!rpcURL) {
        console.log(`obscuro:deploy requires "obscuroEncRpcUrl" to be set as part of the selected network's config.`)
        return;
    } 
    
    // Trigger shutdown on CTRL + C
    process.on('SIGINT', ()=>exit(1));
    
    // Execute the deploy task provided by the HH deploy plugin.
    await hre.run('deploy');
});

// This extends the hardhat config object for networks to have a key for
// the encrypted rpc endpoint of obscuro nodes.
declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}