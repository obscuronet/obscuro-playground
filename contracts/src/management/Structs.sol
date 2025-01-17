// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import * as MessageBusStructs from "../messaging/Structs.sol";

interface Structs {
    struct MetaRollup{
        bytes32 Hash;
        bytes Signature;
        uint256 LastSequenceNumber;
        bytes32 BlockHash;
        bytes32 MessageRoot;
        uint256 BlockNumber;
        bytes32 BlobHash;
    }

    struct RollupStorage {
        mapping(bytes32=>MetaRollup) byHash;
        mapping(uint256=>bytes32) byOrder;
        mapping(uint256=>bytes32) toUniqueForkID;
        uint256 nextFreeSequenceNumber;
    }

    struct HeaderCrossChainData {
        MessageBusStructs.Structs.CrossChainMessage[] messages;
    }
}