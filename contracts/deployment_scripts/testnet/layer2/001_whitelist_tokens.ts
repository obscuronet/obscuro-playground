import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { Receipt } from 'hardhat-deploy/dist/types';

/* This script whitelists the L1 tokens through the bridge and makes sure their wrapped
   versions are created on the L2.
   This is the new version of the HOC and POC on the L2. 
*/


async function sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(resolve, ms);
    });
}

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l1Network = hre.companionNetworks.layer1;
    const l2Network = hre; 

    const l1Accounts = await l1Network.getNamedAccounts();
    const l2Accounts = await l2Network.getNamedAccounts();

    // Get the HOC POC layer 1 deployments.
    const HOCDeployment = await l1Network.deployments.get("HOCERC20");
    const POCDeployment = await l1Network.deployments.get("POCERC20");

    console.log(` Using deployers for bridge interaction L1 address=${l1Accounts.deployer} L2 Address=${l2Accounts.deployer}`);

    // Request the message bus address from the config endpoint
    const networkConfig: any = await hre.network.provider.request({ method: 'net_config' });
    if (!networkConfig || !networkConfig.L2MessageBusAddress) {
        throw new Error("Failed to retrieve L2MessageBusAddress from network config");
    }
    const l2messageBusAddress = networkConfig.L2MessageBusAddress;
    console.log(`Loaded message bus address = ${l2messageBusAddress}`);

    // Tell the bridge to whitelist the address of HOC token. This generates a cross chain message.
    let hocResult = await l1Network.deployments.execute("TenBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", HOCDeployment.address, "HOC", "HOC");

    if (hocResult.status != 1) {
        console.error("Unable to whitelist HOC token!");
        throw Error("Unable to whitelist HOC token!");
    }

    // Tell the bridge to whitelist POC. This also generates a cross chain message.
    const pocResult = await l1Network.deployments.execute("TenBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "whitelistToken", POCDeployment.address, "POC", "POC");
    
    if (pocResult.status != 1) {
        console.error("Unable to whitelist POC token!");
        throw Error("Unable to whitelist POC token!");
    }

    const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
    // Get the hash id of the event signature
    const topic = hre.ethers.id(eventSignature)

    // Get the interface for the event in order to convert it to cross chain message.
    let eventIface = new hre.ethers.Interface([ `event LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)`]);

    // This function converts the logs from transaction receipts into cross chain messages
    function getXChainMessages(result: Receipt) {
        
        // Get events with matching topic as the event id of LogMessagePublished
        const events = result.logs?.filter((x)=> { 
            return x.topics.find((t: string)=> t == topic) != undefined;
        });

        const messages = events!.map((event)=> {
            // Parse the rlp encoded log into event.
            const decodedEvent = eventIface.parseLog({
                topics: event!.topics!,
                data: event!.data
            })!!;
        
            //Construct the cross chain message.
            const xchainMessage = {
                sender: decodedEvent.args[0],
                sequence: decodedEvent.args[1],
                nonce: decodedEvent.args[2],
                topic: decodedEvent.args[3],
                payload: decodedEvent.args[4],
                consistencyLevel: decodedEvent.args[5]
            };

            return xchainMessage;
        })

        return messages;
    }

    let messages = getXChainMessages(hocResult);
    messages = messages.concat(getXChainMessages(pocResult));

    console.log("Attempting to verify cross chain message transfer.");

    // Poll message submission 
    await new Promise(async (resolve, fail)=> { 
        setTimeout(fail, 30_000)
        const messageBusContract = (await hre.ethers.getContractAt('MessageBus', l2messageBusAddress));
        const gasLimit = await messageBusContract.getFunction('verifyMessageFinalized').estimateGas(messages[1], {
            maxFeePerGas: 1000000001,
        })
        try {
            while (await messageBusContract.getFunction('verifyMessageFinalized').staticCall(messages[1], {
                maxFeePerGas: 1000000001,
                gasLimit: gasLimit + (gasLimit/BigInt(2)),
                from: l2Accounts.deployer
            }) != true) {
                console.log(`Messages not stored on L2 yet, retrying...`);
                await sleep(1_000);
            } 
        }catch (err) {
            console.log(err)
            fail(err)
        }

        resolve(true);
    });
    
    // Perform message relay. This will forward the whitelist command to the L2 subordinate bridge.
    // Get the balance of l2Accounts.deployer using provider
    const provider = l2Network.ethers.provider;
    const balance = await provider.getBalance(l2Accounts.deployer);
    console.log(`Balance of l2Accounts.deployer: ${balance}`);

    console.log(`Relaying messages using account ${l2Accounts.deployer}`);
    const relayMsg = async (msg: any) => {
        return l2Network.deployments.execute("CrossChainMessenger", {
            from: l2Accounts.deployer, 
            log: true,
            gasLimit: 5_000_000
        }, "relayMessage", msg);
    };

    console.log(`Relaying message - 0`);
    let results = [await relayMsg(messages[0])];

    console.log(`Relaying message - 1`);
    results = results.concat(await relayMsg(messages[1]))

    results.forEach(res=>{
        if (res.status != 1) {
            throw Error("Unable to relay messages...");
        } 
    });
};

export default func;
func.tags = ['Whitelist', 'Whitelist_deploy'];
func.dependencies = ['EthereumBridge'];
