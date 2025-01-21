# Two-Phase Game Move Execution

Games running on platforms with Smart Transparency and native RNG are vulnerable to a couple of attacks:

### 1. The Proxy execute-revert attack

The attacker deploys a contract that acts as a proxy between the user and the game contract.
The proxy calls the `move` function of the game. Then it queries the result by invoking a view function on the game. 
If it doesn't like the result it just reverts.

This attack can be made smarter by the proxy trying moves in a loop until it wins the prize.


### 2. Gas estimate side-channel attack

The attacker calls the free gas estimate for different moves. `Gas estimate` does not normally return the result and does not modify the state.
But, if the game has a slight difference between the execution paths of winning and losing, then the attacker could infer from the gas used whether the simulated move is a winner.


## Proposed solution - Two-Phase execution

This solution combines a pattern that game developers must implement with some primitives offered by the platform.

### Game pattern - Decouple move recording from move execution

When the user interacts with a game, they just commit to a move which the game stores in an internal queue. No processing should happen in this first phase.
The execution will be performed in phase 2 - as a separate transaction.

Note: The "execution" transaction can't be triggered by the user because then they can perform the proxy attack by executing both phases inside the proxy.

To implement this pattern securely, there must be a trusted third party that calls "execute" on behalf of users. 
This solution is not ideal because it affects the UX - the user will have to wait at least 2 blocks for a response, and it introduces a third party. 
There is also the issue of the execution cost. The third party has to be somehow paid for this service.


### Platform support - execution callback

TEN already has a system contract that is executed at the end of every block. The execution of the sys contract is triggered by the platform as a "synthetic" transaction.

If games had a way to register themselves with this system contract, then it could call the execution phase automatically at the end of each batch.

This solves the UX problem - because the user would get the response at the same time, and it also removes the need for a trusted third party.

#### Implementation 

Ideally, the registration of the callback would be a "smart contract"-only action.  
Basically, the system contract that is called automatically will have a function where it allows anyone to register a callback.


#### Paying for the execution

This is the trickiest component of the solution.

On a high level, the most likely requirement is that the Game would ask the users to pay an amount that covers both the recording of the move and the execution of an average move.
As part of the "phase 1" transaction, the money for the move recording goes to the platform, while the estimated money for "phase 2" goes to an account owned by the game.

Assuming we implement the callback as described in the previous section, we need a way from solidity to calculate the gas cost of executing each move. 
This appears to be possible using a workaround with the `gasLeft()` function.

##### Possible Implementation of execution payment

1. The game developer has to estimate how much an average move will cost to execute and add some margin.
2. During the "phase 1" transaction, the game dev will transfer some `value` from the caller to an address controlled by the sys contract.
   This address was allocated to this game during the callback registration (similar to how you have to reference a payment when you move money to some broker).
   Note: To the user, it will be equivalent to gas spent for the transaction.
4. At the end of the batch, the sys contract will invoke the "phase 2" function of the game - one move at a time.
   It will meter the cost of each move, and check whether there is enough left in the prepaid account to continue.


