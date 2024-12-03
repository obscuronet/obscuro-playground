// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./Structs.sol";

// This represents the smart contract API that allows dApps and accounts to send and verify received messages 
// between layer 1 and layer 2.
interface IMessageBus {

    // The event emitted whenever a message is published. 
    // The enclave listens for it on the deployed message bus addresses.
    event LogMessagePublished
    (
        address sender, 
        uint64 sequence, 
        uint32 nonce, 
        uint32 topic, 
        bytes payload, 
        uint8 consistencyLevel
    );

    event ValueTransfer
    (
        address indexed sender,
        address indexed receiver,
        uint256 amount,
        uint64 sequence
    );

    // This method is called from contracts to publish messages to the other linked message bus.
    // nonce - This is provided and serves as deduplication nonce. It can also be used to group a batch of messages together.
    // topic - This is the topic for which the payload is published. 
    // payload - This is the actual message.
    // consistencyLevel - this is how many block confirmations to wait before publishing the message. 
    // Notice that consistencyLevel == 0 is still secure, but might make your protocol result more prone to reorganizations.
    // returns sequence - this is the unique id of the published message for the address calling the function. It can be used
    // to determine the order of incoming messages on the other side and if something is missing.
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) external returns (uint64 sequence);

    function sendValueToL2(
        address receiver,
        uint256 amount
    ) external payable;

    function receiveValueFromL2(
        address receiver,
        uint256 amount
    ) external;

    // This function verifies that a cross chain message provided by the caller has indeed been submitted from the other network
    // and returns true only if the challenge period for the message has passed.
    function verifyMessageFinalized(Structs.CrossChainMessage calldata crossChainMessage) external view returns (bool);
    
    // Returns the time when a message is final (when the rollup challenge period has passed). If the message was never submitted the call will revert.
    function getMessageTimeOfFinality(Structs.CrossChainMessage calldata crossChainMessage) external view returns (uint256);

    // This is the smart contract function which is used to store messages sent from the other linked layer. 
    // The function will be called by the ManagementContract on L1 and the enclave on L2. 
    // It should be access controlled and called according to the consistencyLevel and Obscuro platform rules.
    function storeCrossChainMessage(Structs.CrossChainMessage calldata crossChainMessage, uint256 finalAfterTimestamp) external;

    // This is a testnet function which allows the bridge owner to retrieve all funds from the message bus.
    function retrieveAllFunds(address receiver) external;
}
