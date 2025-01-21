# Session keys on Ten

EIP-4337 (Account Abstraction) has popularised the "Session keys" concept to minimize clicks.

We want to offer a similar user-friendly experience to Ten dApp devs.

To test the concept, we will convert the Battleships game to a no-click UX.


## The Ten Gateway as a Session Key manager

*Reminder: A session key (SK) is a key that is not managed in a wallet and can be used to sign behind-the-scenes operations without user action.*

Without smart contract wallets, a SK must have a balance to pay for gas. The browser can manage SKs on behalf of the user, but it can lose the gas if the browser crashes.


The Ten Gateway (TG) already manages VKs on behalf of users, so adding SKs is relatively straightforward.

The advantage of the TG managing SKs is that it can return the funds to the EOA anytime.

Another advantage is that they can be reused between sessions to avoid unnecessary transactions.


### The flow

#### Create
The logic of the game UI (javascript) will call an `/session_key/create` endpoint on the TG, which will return the account address corresponding to the SK. 
The game UI will create a value transfer transaction from the main account to this address with enough money to prepay for the moves. The user must sign it with their wallet.

The game will call `/session_key/activate/${sk_account}`. From then on, all transactions originating from that userID will be considered as part of the session. 

After the activation, the game will create unsigned move transactions and submit them to the TG. 
The TG will sign these transactions with the active SK and submit them to the network.

#### Reuse SK with a balance

If the user already has a session key, the game can retrieve it with `/session_key/list`. This will return a list of `$address, $amount`.
The game has the option to top up the value of an account. If there is enough balance, activate it.

#### Ending the game

When the game ends, it will create "cleanup" transactions that move the values accumulated by the Sk to the main account, and submit them to the TG. 

*Note that the TG will sign these with the SK.*

To finish the session, the game will call: `/session_key/deactivate` which returns `address, amount`.

#### Deleting a SK

SKs need to be stored on the TG to allow the user to query the data of previous games in the future.
If that data is not relevant, it is recommended for performance reasons to delete the SK with: ``/session_key/delete/${sk_account}``, after making sure it has no balance.

#### Exporting the SKs

`/session_key/export/${sk_account}` - returns the private key

This can be used for advanced use cases, or as a backup after deletion. 


### Implementation on the Ten Gateway

1. Store SKs in the database and manage them via the `session_key` endpoints: create/delete, activate/deactivate, list
2. On create: generate a k-v pair and sign the SK account (from the new k-p) with the vk of the user. Note: maybe create a new type of "signature" for sks.
3. When the `activate` endpoint is called, we need to set a flag for that user. And add extra logic in the `sendTransaction(Raw)` endpoint to sign transactions with the SK.
4. We should limit the number of SKs per user to 1 or 2.
5. When the user queries data, also include the SKs as candidates. Todo: think about when the data needs concatenation 

### Security consideration

The UserId or the RPC URL should not be exposed to the browser.
Find out why we need it currently.

If we need it for tenscan functionality, we need to find alternatives.
