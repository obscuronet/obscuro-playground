// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

// The ERC20 token bridge interface.
// Calling functions on it will result in assets being bridged over to the other layer automatically.
interface IBridge {
    enum Topics {
        TRANSFER,
        MANAGEMENT,
        VALUE
    }

    // Sends the native currency to the other layer. On Layer 1 the native currency is ETH, while on Layer 2 it is OBX.
    // When it arrives on the other side it will be wrapped as a token.
    // receiver - the L2 address that will receive the assets on the other network.
    function sendNative(address receiver) external payable;

    // Sends ERC20 assets over to the other network. The user must grant allowance to the bridge
    // before calling this function for more or equal to the amount being bridged over. 
    // This can be done using IERC20(asset).approve(bridge, amount); 
    // asset - the address of the smart contract of the ERC20 token.
    // amount - the number of tokens being transferred.
    // receiver - the L2 address receiving the assets.
    function sendERC20(
        address asset,
        uint256 amount,
        address receiver
    ) external;

    // This function is called to retrieve assets that have been sent on the other layer.
    // In the basic implementation it is only callable from the CrossChainMessenger when a message is
    // being relayed.
    function receiveAssets(
        address asset,
        uint256 amount,
        address receiver
    ) external;

    struct ValueTransfer {
        uint256 amount;
        address recipient;
    }
}
