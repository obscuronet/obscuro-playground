// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./ICrossChainMessenger.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";


// TODO: We need to upgrade the open zeppelin version to the 4.x that adds cross chain enabled
abstract contract CrossChainEnabledObscuro is Initializable {
    ICrossChainMessenger messenger;
    IMessageBus messageBus;
    uint32 nonce = 0;

    // The messenger contract passed will be the authority that we trust to tell us
    // who has sent the cross chain message and that the message is indeed cross chain.
    function configure(address messengerAddress) public onlyInitializing {
        messenger = ICrossChainMessenger(messengerAddress);
        messageBus = IMessageBus(messenger.messageBus());
    }

    // Returns if the message is considered to be a cross chain one.
    function _isCrossChain() internal view returns (bool) {
        return msg.sender == address(messenger);
    }

    function _messageBus() internal view returns (IMessageBus) {
        return messageBus;
    }

    // Returns the address of the sender of the current cross chain message.
    // address 0x0 is considered null/no sender.
    function _crossChainSender() internal view returns (address) {
        return messenger.crossChainSender();
    }

    // Ensures that the message is coming from another chain and for this contract.
    // Combined usage of _isCrossChain and _crossChainSender prevents attacks where
    // a message is relayed from the correct sender to a different contract that
    // maliciously calls into this contract.
    modifier onlyCrossChainSender(address sender) {
        require(
            _isCrossChain(),
            "Contract caller is not the registered messenger!"
        );
        require(
            _crossChainSender() == sender,
            "Cross chain message coming from incorrect sender!"
        );
        _;
    }

    // This function should be called to queue messages that are encodings of function calls
    // using the abi.encodeWithSelector(bytes4, arg) pattern.
    // target - the contract on which the call will be executed.
    // message - the result of abi.encodeWithSelector(bytes4, arg)
    // topic - TODO: determine if should be removed.
    // gas - when doing target.call{} this gas will be the limit set for the call. *Currently Unused*
    // consinstencyLevel - Block confirmations before finalizing message. Obscuro allows that even 0 is secure,
    // but it might not be suitable for all protocols.
    function queueMessage(
        address target,
        bytes memory message,
        uint32 topic,
        uint256 gas,
        uint8 consistencyLevel
    ) internal {
        bytes memory payload = abi.encode(
            ICrossChainMessenger.CrossChainCall(target, message, gas)
        );
        messageBus.publishMessage(nonce++, topic, payload, consistencyLevel);
    }
}
