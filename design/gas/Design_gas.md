# Obscuro Gas Design

## Overview

The revenue generation mechanism for Obscuro hinges on the collection of gas fees. Notably, as Obscuro operates as a layer 2 protocol, the cost of operations surpasses that of layer 1, since it has to cover layer 1 gas costs in addition to its operational expenses. Unlike layer 1 miners, who only shoulder operational costs, there's no expenditure prerequisite for block publication on Obscuro except for the static cost of stake.

The gas mechanics of a protocol normally have the following functions:
1) Revenue for node operators.
2) Congestion management and prevention of denial of service.
3) Price discovery mechanism for transactions.

## Arbitrum Gas Mechanics

Arbitrum does not describe how their gas mechanics work and how would price change with network congestion. It is implied in [this](https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9) article that there is indeed a mechanism that increases the L2 cost, but their public documentation seems to just call the gas price unit `ArbGas` without elaborating how its derived. It's possible that it is actually static and what is implied in the article is that the L1 price influenced the change in the L2 gas cost. 

## Optimism Gas Mechanics

Optimism is EIP-1559 compliant and has a [public dashboard](https://public-grafana.optimism.io/d/9hkhMxn7z/public-dashboard?orgId=1&refresh=5m) which showcases gas prices for both L1 and L2. During time of writing, the gas dashboard clearly showed different spikes between L1 prices and L2 prices. They do not elaborate however when is the base fee modified, presumably its at the time of a rollup as they have no concept of batches in their current live version. They mention that the parameters of EIP-1559 are different, presumably the gas fee increase/decrease params. Optimism bedrock should however introduce blocks and those parameters should change. 

### Initial Bootstrapping of Network Gas

Similar to other layer 2 protocols, Obscuro seeks to leverage Ethereum for gas usage. To facilitate this, Ethereum must first be bridged over to Obscuro. Although the bridge typically operates with a relayer, relaying a message incurs a gas cost. Hence, to bootstrap the entire layer 2 protocol, including relayers and bridges, cross-chain messages concerning Ethereum depositing into an account will be automatically executed on layer 2 at no cost and without the need for relaying. This ensures no additional transactions that necessitate publishing to layer 1.
    

### Types of Transactions

Ethereum has developed to accommodate several [transaction types](https://docs.infura.io/networks/ethereum/concepts/transaction-types):
1) Regular transactions - they work based on a gas auction.
2) Access list transactions - same as previous ones, but they also announce what storage access they are supposed to perform when executed.
3) EIP-1559 - the new transactions that burn the gas fee and work with tips in order to achieve a fair system.

Usually layer 2's support EIP-1559 so there is a pending question of what we would need to support. Normally all support EIP-1559, but not all support gas auction (or mention if they do). As most wallets have moved over to EIP-1559 it is preferable to start with support for it and add the support for the gas auction later on after everything else is finished.
Note that Binance for example still uses Type 1 transactions and has not migrated to EIP-1559. 

## Requirements

1. Ethereum must be used for gas operations and bridged from the layer 1 protocol.
2. Bridged gas should automatically deposit into the respective addresses.
3. Gas-consuming transactions should include a layer 1 cost component for calldata publishing and a layer 2 component that covers execution costs and prevents endless transaction execution.
4. L2 gas price should be configurable. Furthermore it should support being free whilst still not allowing non halting execution.
5. The gas mechanics need to prevent denial of service attacks.
6. EIP-1559 transactions should have their base fee + tip deposited to the sequencer without any ETH being burned.

## Gas Cost Display

Metamask can disassemble the cost components for layer 1 and layer 2 for specific transactions on Optimism, such as token balance approvals. However, it remains unclear how they achieve this as it's neither fully documented nor supports all transactions.


## Design for Gas Deposit

Given that Ethereum as a currency operates on the layer 1 protocol, it needs to be bridged to Obscuro for usage. This transition mechanism will be facilitated by the standard Obscuro bridge. The deposited value will produce a cross-chain message indicating the Ethereum recipient through the bridge smart contract.

```EthereumDeposited(address receiver, uint256 amount)```

The enclave will automatically pick up this message, along with other cross-chain messages, but it will be treated specially. Upon detection, the enclave will increase the balance of the receiver account by the amount specified in the message. This automatic relaying is needed in order to avoid the chicken and egg problem one would have initially where no relayer has any gas to relay the gas deposit messages. It would also ensure relayers who run out of gas can easily recharge without relying on other relayers if we were to use a different bootstrap mechanism. Note that those deposits will not call the fallback function of the address if it is a smart contract and any deposits to such contracts would later on need a message relayed to verify the message as they would normally.

## Gas Estimation

When an RPC call for gas estimation is made to Obscuro, the returned estimate should include the expected layer 1 calldata cost alongside the usual execution-based gas estimation. This can be achieved by taking an estimation for layer 1 minFee per gas, determining the calldata cost of the transaction, and adding this to the standard gas estimate.

`Gas estimate = L1_gas_fee * calldata_gas + estimated_l2_gas`

## Gas Expenditure

Before executing a transaction at Obscuro, we must automatically deduct the layer 1 costs for the transaction from the sender's balance. This can be done directly through the stateDB before processing. Ideally, the remaining balance will be refunded, simulating a scenario where the gas limit has been spent. However, potential overruns of execution costs could result in transactions exceeding their set gas limit.

The best approach would be to both decrease the balance and reduce the transaction's gas limit, assuming that layer 1 prices will be factored into this gas limit, given that the limit is set based on the eth_estimateGas call.


## Monitoring and Adjusting the L1 Fee

In order to pay for rollups without losing money, Obscuro needs to collect sufficient gas fees for the L1 cost including any differences caused due to time. EIP-1559 made the gas fee predictable. With each block a miner might choose to increase or reduce it by a % of the previous block. This means we can implement a prediction mechanism that prices the worst case scenario by increasing the current L1 gas cost with the required percentage based on how many blocks will the L1 produce until we hit our rollup duration limit.

It's important to note that Optimism has chosen not to track L1 prices precisely and impose a limit on how much transaction fees can be affected. They initially had a 10% overcharge on the l1 gas cost to cover for potential spikes, but have now migrated to a model they don't really explain in their docs.

This leaves us with two options: 
1) Attempt to precisely estimate gas fees and refund anything overcharged when the rollup is published
2) Do a simple overcharge and follow what Optimism have done as this model has proven quite successful for them revenue wise.


## Sequencer Compensation

Upon the production of a batch, all transactions enclosed within it would have made payments for gas. The aggregate sum, encompassing both layer 1 and layer 2 costs, will subsequently be debited to the address denoted as coinbase within the batch. It would be more practical for this address to be configurable. The reason being that the sequencer's address is derived from a key that must be safeguarded, and implementing any logic for fund transfer from it would necessitate code development.

This entire logic can happen in the `evm_facade.go`.

## Batch gas limit

In order to prevent denial of service attacks we need to have a batch gas limit. This would lend itself well to Obscuro being EIP-1559 complaint making us able to dynamically price L2 fees based on congestion. For the initial version of the gas mechanics we can use a fixed L2 gas price and later on migrate to EIP-1559 compliancy. Taking cue from Optimism's approach to modifying the parameters, we should divide the parameters for Obscuro proportionally to L1 block creation time vs Obscuro L2 batch creation time. This should yield approximately identical behaviour.