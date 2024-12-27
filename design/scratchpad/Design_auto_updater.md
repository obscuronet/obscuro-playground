# Ten Auto Updater - Scheduled callbacks

## Introduction

In Ethereum it is well known that transactions cannot originate from contracts. There is no way to automatically run some update. Someone has to init the transaction. If you have a contract on a layer 2 that needs rebalancing for example someone has to initiate it. This results in having to publish the transaction on the L1 and pay for it.


This feature suggestion would allow to circumvent the L1 publishing by creating a definition of "L2 state derived transactions" which are akin to cross chain synthetic transactions.


### Cross Chain Synthetic Transactions Primer

The synthetic transactions we currently have are not published to the L1. They are not included in the batch albeit mutating its state. This is not an issue, because the transactions are deterministically derived from the L1 state or more specifically, the L1 block that the batch points to. This means that when a validator is recomputing the batch, the data for rebuilding the transactions, in the same order is available implicitly. If we were to publish them to the L1 it would only be redundant and increase costs.


## L2 State Derived transactions / Scheduled callbacks

L2 state derived transaction is a special synthetic transaction, also started by the sequencer that would call a specific contract. This contract would have to register itself as updateable and prepay L2 gas costs to the `block.coinbase`. Using this prepayed funding, the sequencer would create the implicit functions calling the entry point function with no calldata and no value. It is possible to extend to support implicit calldata sent through the prepayment, but its best to keep it simple for the initial design.

The updatable contracts would need to export a callable method that determines if an update transaction should happen. Let's say this function signature will be `shouldAutoUpdate() returns (bool)`. 
Whenever a batch is produced, after all the transactions have been applied the sequencer would go through the registered contracts and `derive` the synthetic auto updates that need to happen. Then they would be applied as any other transaction would, where the gas limit and cost would be put as the prepaid amount. Anything unused will be refunded back to the contract.

#### Computation implications

The callable `shouldAutoUpdate` function would need to be subsidized by the sequencer and thus would have hard cap gas limit. This is similar to how optimism grants free gas for some special auto calls on cross chain messages. Alternatively we can use a scheduling mechanic where the contract instructs the sequencer when a transaction should happen, similar to Arbitrum's retryable transactions. When something is scheduled, it would be prepaid and the executed scheduled call should reschedule a new one to keep the automation going.

You can imagine the scheduling approach akin to doing a recursive javascript `setTimeout`. `block.coinbase` is payed and a system function is called on a predeployed contract. Contracts can use this to conviniently have the first call to them in batch trigger a schedule, paid for by the contract. Gas costs with schedules can be problematic however if they are allowed to happen at arbitrary point in time, as contracts will be unable to predict costs and might fail the schedule (and reschedule) transaction due to being out of gas. This would effectively terminate their loop if they operate based on such. Given that excess will be refunded it shouldn't be a huge problem if contracts oversupply the schedule payment. It's the approach they seem to have with chainlink callbacks.

#### Use cases

This feature can be good for a plethora of use cases that already rely on some sort of authorized caller to update the contract periodically. It would make those more decentralised, as the rules for when it should be called would be in the smart contract and early coming transactions would not be derived by the validator. 
For example you can have option protocols that auto excercise and settle on expiry, instead of the current approach that requires people to manually execute the settlement. 
It's also possible to have automatic liquidations for lending protocols. 
Another obvious application would be automatic arbitrage.
Gas intensive dApps could also greatly benefit - imagine a uniswap proxy router contract that users just queue swaps in and they all get batch executed automatically at the end of the batch. If there is a heavily traded pair with a lot of concurrent traffic to it, this would reduce the expensive storage read/writes and might even match trades up front to reduce LP fees.
Any use case that requires a bot to trigger on chain stuff would be a good fit as this feature essentially removes the need to run a bot, unless the bot does a ton of off chain computing.

It's important to note that most dApps aren't structured to be automated in such a way - there is no way to iterate their debt for example, as its stored in maps and people search for liquidation opportunities off chain. Regardless of this, it's reasonable to expect the feature will be utilized. Contracts currently are mostly structured around rebalancing on user inputs directly, along with exposing some public functions that can do another rebalance based on calling other contracts and so on. This is highly inefficient as most of those contracts would be far better off using an automated update system that is paid for by the users who would split the cost instead of paying all of it. While its easy to argue that L2 execution costs are so low right now that this is irrelevant, lowering the barrier even more would enable more use cases and make current ones more efficient. 


And lastly, there is a very exotic use case for this feature that isn't really possible anywhere - security monitoring. A contract can register an automatic callback that checks some constraints and if something is wrong pause everything automatically. 

## Tech Requirements

We'd need a predeployed system smart contract that allows for registering contracts:

```solidity
interface TenAutomationRegistry {
    function registerCallback(bytes memory callInstructions) external payable returns(uint256)
    function uniqueCaller(address targetContract, bytes memory callInstructions) pure returns (address)
    function cancel(uint256 id)
}
```

dApps will use this system contract to register callbacks. Whenever `registerCallback` is called a record will be put to call `msg.sender` with `callInstructions`.
It is payable in this design, because we might need to put some arbitrary cost to prevent DOS attacks.
The `uniqueCaller` getter would return the unique address who will be the `tx.origin` for a specific contract's callback. This will allow the dApps to limit who can call the special automatic functions if they wish to do so. Futhermore it will be a layer of security as the caller address will be isolated per contract's callback, in order to prevent potential security exploits, albeit I'm not able to come up with an example.


To implement the whole feature we'd need to extend the logic in the `BatchExecutor` component and simply add another layer of transactions at the end to be executed by the `evm_facade`. Those transactions will be executed with gas priced same as normal user transactions. 

**Cancellation** - There must be a way to cancel scheduled callbacks. Assuming that the implementation of the feature is inspired by javascript `setTimeout()` callbacks, we can have the registration return a `uniqueID` which can be used as a parameter to `cancel(id)`. Cancelling should refund collateral and prefunding, but this would require to pay to something else instead of `block.coinbase`. This could be the predeployed contract.

The need to cancel might come from a contract being updated and signatures not matching, contract self destructing or some other imposed constraint that would mean a scheduled callback is a bad idea.


### Performance considerations

If instead of scheduling, we use a callable `shouldAutoUpdate`, having too many registered contracts can start slowing down batch production. This is the only notable detail as we would have to subsidize those, but putting a cost on registering should still be somewhat of a deterrent, given that those calls will use limited gas cap.

As for potential slowdown of batch production due to the expensive computations when doing an automatic update - As this will have to be paid for by the contract/user/whatever to `block.coinbase` beforehand, I don't see any difference between `L2 derived transactions` and `L2 user transactions`. If a user transaction takes a while to finish on the EVM, it would pay for it. Same for the derived transactions.

There shouldn't be any noticeable performance degradation even if this feature is heavily utilised. When looking at the performance of the sequencer previously the bottleneck was never in the EVM processing. Even a ton of storage mutating transactions were taking sub milliseconds. If contracts are engineered smartly, with dirty flags for example, the automatic updates will be blazing fast as they would be hitting storage slots that are already in the in memory state tree. The fact that the auto updates would be derived from the L2 state means that most of its use cases would be triggered by mutations to warm storage slots.


## Futher improvements

This feature is very symbiotic with some newer L2 concepts like cross domain state reads. We can have scheduled callbacks query a system contract that provides them with data from the L1 - transactions for a block, state of a contract, balances and etc. This will enable implicit L1 to L2 state propogation without having to manually do cross chain message submission and consumption transactions on the L1 and L2 respectively. Potentially we can have automatic oracles, synthetics and other cool dApps spin out from this. Maybe some cool bridge. 

Another avenue for imporving this would be to introduce complex scheduling mechanics - for example when registering, a contract would pass data for a static call to another contract. The system contract would evaluate the static call and put listeners on all the storage slots that were accessed through the **public** method provided when registering. When one of them gets updated, the schedule would be executed at the end of a batch. It would give 3rd party contracts the ability to put dirty flags on publicly exposed storage slots. I imagine this would be really great for automated arbitrage. "Please notify me when this pool in balancer changes, but this one in Uniswap does not."; Perhaps the complex scheduling can also include constraints - "Notify me if delta of this storage slot is bigger than uint256(50)".