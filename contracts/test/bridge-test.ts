import { expect } from "chai";
import hre, { ethers } from "hardhat";
import { time } from "@nomicfoundation/hardhat-network-helpers";
import { bridge } from "../typechain-types/src";
import { MessageBus, ObscuroBridge, WrappedERC20__factory } from "../typechain-types";
import { EthereumBridge } from "../typechain-types/src/bridge/L2/EthereumBridge";
import { CrossChainMessenger } from "../typechain-types/src/messaging/messenger";
import { Contract } from "hardhat/internal/hardhat-network/stack-traces/model";


import type {
  BaseContract,
  ContractTransaction, ContractTransactionResponse
} from 'ethers';
import { WrappedERC20 } from "../typechain-types/src/common";

describe("Bridge", function () {

  let busL1: MessageBus
  let busL2: MessageBus

  let messengerL1: CrossChainMessenger
  let messengerL2: CrossChainMessenger

  let bridgeL1 : ObscuroBridge
  let bridgeL2 : EthereumBridge

  let erc20address : any

  this.beforeEach(async function(){
    const MessageBus = await hre.ethers.getContractFactory("MessageBus");
    const Messenger = await hre.ethers.getContractFactory("CrossChainMessenger");
    const L1Bridge = await hre.ethers.getContractFactory("ObscuroBridge");
    const L2Bridge = await hre.ethers.getContractFactory("EthereumBridge");

    const [owner] = await ethers.getSigners();

    const ERC20 = await hre.ethers.getContractFactory("ConstantSupplyERC20", owner);

    console.log(`Deploying erc20`);
    try {
      const erc20 = await ERC20.deploy("XXX", "XXX", 100000);
      erc20address = await erc20.getAddress();
    } catch(err) {
      console.error(err);
    }


    console.log(`Deployed erc20`);

    busL1 = await MessageBus.deploy();
    busL2 = await MessageBus.deploy();

    messengerL1 = await Messenger.deploy();
    await messengerL1.initialize(busL1.getAddress());
    messengerL2 = await Messenger.deploy();
    await messengerL2.initialize(busL2.getAddress())

    bridgeL1 = await L1Bridge.deploy();
    bridgeL1.initialize(messengerL1.getAddress());
    bridgeL2 = await L2Bridge.deploy();
    bridgeL2.initialize(messengerL2.getAddress(), bridgeL1.getAddress());

    const tx = await bridgeL1.setRemoteBridge(bridgeL2.getAddress());
    await tx.wait();
  });

  it ("Contracts exists", async function() {
    // This test feels redundant as beforeEach would fail ... but I don't trust javascript.
    expect(await busL1.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
    expect(await busL2.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
    expect(await messengerL1.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
    expect(await messengerL2.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
    expect(await bridgeL1.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
    expect(await bridgeL2.getAddress()).to.not.hexEqual(ethers.ZeroAddress);
  });

  async function submitMessagesFromTx(tx: ContractTransactionResponse) {

      const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";
      const topic = ethers.id(eventSignature)
      let eventIface = new ethers.Interface([ `event LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)`]);

      const receipt = (await  tx.wait())!!;

      const events = receipt.logs.filter((x)=> { 
        return x.topics.find((t)=> t == topic) != undefined;
      }) || [];

      if (events.length == 0) {
        return null
      }
     
      const promises = events.map(async (event) => {
          const decodedEvent = eventIface.parseLog({
            topics: event!.topics!,
            data: event!.data
          })!!;
              
          const xchainMessage = {
            sender: decodedEvent.args[0],
            sequence: decodedEvent.args[1],
            nonce: decodedEvent.args[2],
            topic: decodedEvent.args[3],
            payload: decodedEvent.args[4],
            consistencyLevel: decodedEvent.args[5]
          };

          // If the event was emitted from L1 then we want to submit on L2, otherwise reverse.
          // same for messenger.
          let bus : MessageBus = event.address == await busL1.getAddress() ? busL2 : busL1;
          let messenger : CrossChainMessenger = event.address == await busL1.getAddress() ? messengerL2 : messengerL1;
          await (await bus.storeCrossChainMessage(xchainMessage, 1)).wait();
           
          return { 
              msg: xchainMessage,
              messenger : messenger,
          };
      });
      const bindings = await Promise.all(promises);
      
      // This allows to relay selectively or not to in order to enable test cases.
      return { 
        relayAll: async ()=> {
          const receipts = bindings.map(async (x)=>{
            const tx = await x.messenger.relayMessage(x.msg)
            return await tx.wait();
          })
          return await Promise.all(receipts);
        },
        bindings : bindings
      };
  }

  it("Bridge owned wrapped token should be inaccessible externally", async function () {
      const wrappedERC20 : WrappedERC20__factory = await hre.ethers.getContractFactory("WrappedERC20");
      const [owner] = await ethers.getSigners();

      const whitelistTx = bridgeL1.whitelistToken(erc20address, "o.ZZZ", "o.ZZZ");
          
      await expect(whitelistTx, "Transaction whitelisting the erc20 token failed!").to.not.be.reverted;
      let messages = await submitMessagesFromTx(await whitelistTx);
      await expect(messages, "Missing message to create wrapped tokens on L2 bridge.").not.null;
      await messages!.relayAll();

      const localERC = await bridgeL2.remoteToLocalToken(erc20address);
      const l2Erc20 : WrappedERC20 = wrappedERC20.attach(localERC);

      await expect(l2Erc20.issueFor(owner.address, 5_000_000)).reverted
  });

  it("Bridge relaying published message from different sender should fail", async function () {
      const whitelistTx = await bridgeL1.whitelistToken(erc20address, "o.ZZZ", "o.ZZZ");
        
      await expect(whitelistTx, "Transaction whitelisting the erc20 token failed!").to.not.be.reverted;
      let messages = await submitMessagesFromTx(whitelistTx);
      expect(messages, "Missing message to create wrapped tokens on L2 bridge.").not.null;
      await messages!.relayAll();

      const [owner] = await ethers.getSigners();
      await expect(bridgeL2.receiveAssets(erc20address, 500, owner.address), "Only messenger should be able to call receiveAssets")
        .revertedWith("Contract caller is not the registered messenger!");

      const encodedData = bridgeL2.interface.encodeFunctionData("receiveAssets", [erc20address, 500, owner.address]);

      const encodedCalldata = await messengerL2.encodeCall(await bridgeL2.getAddress(), encodedData);

      const tx = busL1.publishMessage(0, 0, encodedCalldata, 0);
      expect(tx, "Anyone should be able to publish a message!");

      messages = await submitMessagesFromTx(await tx);
      expect(messages, "publishing a message should create a cross chain event").not.null;
      const publishedFakeMessage = messages!.bindings[0].msg

      expect(messengerL2.relayMessage(publishedFakeMessage))
        .revertedWithCustomError
  });
  
  it("Bridge relay unpublished message should fail", async function () {
      const whitelistTx = bridgeL1.whitelistToken(erc20address, "o.ZZZ", "o.ZZZ");
          
      await expect(whitelistTx, "Transaction whitelisting the erc20 token failed!").to.not.be.reverted;
      let messages = await submitMessagesFromTx(await whitelistTx);
      await expect(messages, "Missing message to create wrapped tokens on L2 bridge.").not.null;
      await messages!.relayAll();

      const [owner] = await ethers.getSigners();

      const xCrossChainCallData = ethers.AbiCoder.defaultAbiCoder().encode(
          ['address', 'bytes', 'uint256'],
          [await bridgeL2.getAddress(), "0x00", 0]
      );

      const unpublishedFakeMessage = {
        sender: owner.address,
        sequence: 0,
        nonce: 0,
        topic: 0,
        payload: xCrossChainCallData,
        consistencyLevel: 0,
      };

      await expect(messengerL2.relayMessage(unpublishedFakeMessage), "Attempting to relay fake message should revert")
        .revertedWith("Message not found or finalized.");
  });

  /* //TODO: @siliev - reenable; new ethers is fucking it up, but e2e tests pass so they cover this
  it("Bridge mock environment full test.", async function () {
      const [owner] = await ethers.getSigners();

      const wrappedERC20 = await hre.ethers.getContractFactory("WrappedERC20");
      const l1Erc20 : WrappedERC20 = await wrappedERC20.deploy("ZZZ", "ZZZ");
      const whitelistTx = bridgeL1.whitelistToken(await l1Erc20.getAddress(), "o.ZZZ", "o.ZZZ");
      
      
      await expect(whitelistTx, "Transaction whitelisting the erc20 token failed!").to.not.be.reverted;
      let messages = await submitMessagesFromTx(await whitelistTx);
      await expect(messages, "Missing message to create wrapped tokens on L2 bridge.").not.null;
      await messages!.relayAll();

      await expect(await bridgeL2.wrappedTokens(erc20address), "L2 bridge should return zero for non whitelisted contracts.")
        .to.hexEqual(ethers.ZeroAddress);
        
      const localErc = await bridgeL2.remoteToLocalToken(await l1Erc20.getAddress());
      const l2Erc20 : WrappedERC20 = wrappedERC20.attach(localErc);

      const l2Erc20ForOwner = l2Erc20.connect(owner);
      const l1Erc20ForOwner = l1Erc20.connect(owner);


      expect(await bridgeL2.wrappedTokens(await l2Erc20.getAddress()), "L2 bridge should not return zero for whitelisted contract.")
        .to.not.hexEqual(ethers.ZeroAddress);

      await expect(l1Erc20.issueFor(owner.address, 10_000_000), "Failed to mint L1 token").not.reverted;
      await expect(l1Erc20.approve(bridgeL1.getAddress(), 9_000_000), "Failed to increase allowance!").not.reverted;

      await expect(bridgeL1.sendERC20(l1Erc20.getAddress(), 10_000_000, owner.address), "Sending more than allowed should revert").reverted;

      const sendAssetsTx = bridgeL1.sendERC20(l1Erc20.getAddress(), 9_000_000, owner.address);      
      await expect(sendAssetsTx, "Sending as much as allowed should not revert").not.reverted;

      await expect(await l1Erc20ForOwner.balanceOf(owner.address, { from: owner.address }), "Remaining L1 balance should be initial minus bridged amount!")
        .to.equal(10_000_000 - 9_000_000);

      messages = await submitMessagesFromTx(await sendAssetsTx);
      await expect(messages, "Sending assets to L2 resulted in no messages!").not.null;


      await expect(await l2Erc20ForOwner.balanceOf(owner.address, { from: owner.address }), "There should be no balance before relaying stored messages!").to.equal(0);

      await messages!.relayAll();

      await expect(await l2Erc20ForOwner.balanceOf(owner.address, { from: owner.address }), "Relay should have granted balance").to.equal(9_000_000);

      await expect(l2Erc20.approve(bridgeL2.address, 8_000_000), "L2 allowance increase should not revert.").not.reverted;

      const bridgeBackTx = bridgeL2.sendERC20(l2Erc20.getAddress(), 8_000_000, owner.address);
      await expect(bridgeBackTx, "Sending assets back to L1 should not revert").not.reverted;
    
      messages = await submitMessagesFromTx(await bridgeBackTx);
      await expect(messages, "Sending assets back to L1 should produce cross chain messages").not.null;
      await messages!.relayAll();

      await expect(await l2Erc20ForOwner.balanceOf(owner.address, { from: owner.address }), "Remaining L2 balance should be reduced!").to.equal(1_000_000);
      await expect(await l1Erc20ForOwner.balanceOf(owner.address, { from: owner.address }), "New L1 balance should match leftover + bridged amount")
        .to.equal(1_000_000 + 8_000_000);
  }); */

  it("Whitelisting tokens works and relaying creates L2 contracts.", async function () {

      const whitelistTx = bridgeL1.whitelistToken(erc20address, "XXX", "XXX");

      await expect(whitelistTx)
        .to.emit(busL1, "LogMessagePublished");


      const eventSignature = "LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)";

      const topic = ethers.id(eventSignature)
      const event = (await (await whitelistTx)!!.wait())!!.logs?.find((x: any)=> { 
          return x.topics.find((t: string)=> t == topic) != undefined;
      });

      await expect(event).to.not.be.undefined;

      let eventIface = new ethers.Interface([ `event LogMessagePublished(address,uint64,uint32,uint32,bytes,uint8)`]);

      const decodedEvent = eventIface.parseLog({
        topics: event!.topics!.map((v)=>v),
        data: event!.data
      })!!;

      const xchainMessage = {
        sender: decodedEvent.args[0],
        sequence: decodedEvent.args[1],
        nonce: decodedEvent.args[2],
        topic: decodedEvent.args[3],
        payload: decodedEvent.args[4],
        consistencyLevel: decodedEvent.args[5]
      };

      const storeMessage = busL2.storeCrossChainMessage(xchainMessage, 1);
      await expect(storeMessage).to.not.be.reverted;

      const tx = messengerL2.relayMessage(xchainMessage);
      await expect(tx).to.not.be.reverted;

      const localErc = await bridgeL2.remoteToLocalToken(erc20address);

      //bridge L1 sent cross chain message for erc20address when we whitelisted it.
      await expect(await bridgeL2.wrappedTokens(localErc))
        .to.not.hexEqual(ethers.ZeroAddress);

      //random address should not work.
      await expect(await bridgeL2.wrappedTokens(ethers.getAddress("0x8ba1f109551bd432803012645ac136ddd64dba72")))
        .to.hexEqual(ethers.ZeroAddress);

      await expect(messengerL2.relayMessage({
        sender: decodedEvent.args[0],
        sequence: decodedEvent.args[1],
        nonce: 1,
        topic: decodedEvent.args[3],
        payload: decodedEvent.args[4],
        consistencyLevel: decodedEvent.args[5]
      })).to.be.revertedWith("Message not found or finalized.");
  });

  it("MessageBus retrieveAllFunds method should allow owner to extract all native funds from the message bus", async function() {
    const [owner] = await ethers.getSigners();
    const amount = ethers.parseEther("0.01");

    const tx = await busL1.sendValueToL2(owner.address, amount, {
      value: amount
    });
    await expect(tx).to.not.be.reverted;

    // check that the funds were received
    await expect(await ethers.provider.getBalance(busL1.getAddress())).to.equal(amount);

    // retrieve all native funds from the message bus contract on the L1
    const retrieveAllFundsTx = busL1.retrieveAllFunds(owner.address);
    await expect(retrieveAllFundsTx).to.not.be.reverted;

    // check that the funds were drained
    await expect(await ethers.provider.getBalance(busL1.getAddress())).to.equal(0);
  });
});
