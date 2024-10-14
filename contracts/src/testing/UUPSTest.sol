// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract UUPSTestV1 is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    string public message;

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        message = "This is version 1";
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    function setMessage(string memory newMessage) public {
        message = newMessage;
    }

    function getVersion() public pure returns (uint256) {
        return 1;
    }
}

contract UUPSTestV2 is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    string public message;

    function initialize() public reinitializer(2) {
        message = "This is version 2";
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    function setMessage(string memory newMessage) public {
        message = newMessage;
    }

    function getVersion() public pure returns (uint256) {
        return 2;
    }
}