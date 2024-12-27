# Ten Protocol - Async Cross Chain Messages

This document serves to capture some rationale behind preparing the messaging infrastructure for asynchronous out of order publishing of cross chain messages. This serves to allow relayers to fast publish messages before the batches in the DA layer for applications that can work with non finalized (pending) messages. 


## Requirements
1 - The DA layer must be unblocked from reorgs (it currently is, as l1 proof is not checked)
2 - There must be a path forward for fast outgoing xchain messages. This means being able to publish messages
before the relevant DA is posted on the L1 and have them in a `pending` state.
3 - Relayers should be able to publish messages themselves by going through the sequencer and getting a signed bundle.

## Problems

1 - How do we determine the challenge period for bundles?
2 - If bundles are published in reverse order, how do we enforce in order challenge period in order to prevent out of order submission to contracts?
3 - What abstractions do we need to protect smart contracts from low level exploits?
4 - How to differentiate between rollup signatures and cross chain bundle signatures and prevent "malleability"?
5 - How to expose the ability to selectively get cross chain messages and have them signed.


## Proposal

In order to enable relaying of messages and have them out of order, we can migrate to a merkle tree structure.
The management contract will now accept the `rootHash` of the xchain messages for a specific `L1 proof`.
When this `rootHash` is published early by a relayer, the management contract will refund the gas cost of this publishing for the concrete call into it.

The relayer will build transactions that contain `inclusion proofs` pointing to a `rootHash` and messages and feed them into target contracts. Those contracts will call into the message bus to verify the inclusion proofs and proceed with their logic. This way relayers will pay the entirety of the gas cost for external dapps.

Furthermore we can put rules in place for message ordering based on the `emitter` (whomever called the message bus to send the message) and enforce requirements to provide all messages in order. This will limit the potential for attacks where messages are relayed out of order and interfere with administrative functions of a contract while still allow for asynchronous relaying across different `emitters`.

Having a `rootHash` publishing scheme along with per contract ordering rule means that inclusion proofs will also become valid in order of message emission. One wouldn't be able to relay messages in reverse with different inclusion proofs pointing to different `rootHashes` as the verification will fail when trying to supply a message with a nonce in the future. 

For this to work we will need a new endpoint `ten_getMessageBundle` that will take arguments in the form of block hashes which represents a range over which the enclave will collect and bundle messages in order to build the `rootHash` of the merkle tree. Furthermore, we can add the query filter that is used in normal querying for emitted events in order to make the process selective.

`ten_getMessageBundle` will return the messages, inclusion proofs for them and a signature over the `rootHash`. Furthermore there will be metadata (also signed) representing the l1 block for which the `rootHash` is valid to publish and the batch hashes representing the range from which the `rootHash` is built (and sequence numbers?). 

The messages will be in a `pending` state until the management contract receives a DA rollup which ends with a higher `sequenceNo` than the batch from which the message was emitted. After this the challenge period begins. After the challenge the messages are considered `finalized`. 

## Challenges 

Validators will monitor submissions of cross chain messages. When they acquire the relevant batch data for which a `rootHash` was generated, they will verify that iterating over it and building a merkle tree produces the same `rootHash`. If this is not true, they will begin a challenge. The challenge process might need to wait for the DA rollup encompassing all relevant messages to be published. As there are no fraud proofs currently it is unclear how the batch part of the challenge will be processed, but the message bundle challenge will require to prove that how the message where a difference is establish is generated and if it matches the sequencer or the validator's view. 

## Additional details

1 - The management contract will only accept `rootHashes` that end in higher `sequenceNo` than what was last published. This will prevent relayers from republishing hashes encompassing the same messages with slight deviations out of order.
2 - There is not much point of storing received messages once proven by their `inclusionProof`; Contracts asking for validity can store them as consumed and process directly as callbacks because the relayer is paying.
3 - Rollups and message bundles will be signed going through some formatting that makes it impossible to present one as the other. For example we do not want to present a signature over a rollup hash as a signature over a merkle tree root. It shouldn't really be possible to exploit this anyway, but there might be some extreme edge cases which are hard to predict.