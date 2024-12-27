# Introduction

This design will formulate a high level protocol to recover the enclave state if all enclaves have disappeared and the secret is lost. The basic premise is that encryption keys will be generated per rollup. Those encryption keys will be derived from the result of an ASIC resistant hash function. The input to the hash function will be a random number in a constrained range, part of the full uint256 range. The constrained range will be randomly selected for each key rotation and the size
of the range will represent the `difficulty` to bruteforce. This would allow for selecting a `difficulty` that is too challenging to decrypt the data before the protocol reveals it anyway, but at the same time easy enough that it is feasible given computational resources and some time.

The intent is that if everything catastrophically fails, it would still be possible to decrypt all the data, boot up a newly promoted sequencer and continue to run from the latest checkpoint (represented by a rollup on the l1). As this is a last resort solution and ideally would never need to be used, the time to wait for a bruteforce should be ok. Note that this creates a trade off between privacy and speed of recovery based on the `difficulty`. 

This design proposal also overlaps in spirit with the escape hatch - when encryption is bruteforced and state is recovered, escaping the network would become feasible again, from the most accurate state possible.

## Primer on RandomX (ASIC resistant hash function)

RandomX is the hashing function used within Monero to make mining viable only on CPUs. The way it works is that we initialize the hasher with a random key. Based on this random key, the hasher builds a random source code for its VM. Then each input to be hashed goes through this VM with unique source code and ends up as a hash result. It's always guaranteed to get a hashing result. RandomX has a ton of upgrades over the years to make it super resistant, but at its core the logical premise is simple - if the code executed changed randomly each time, then ASICs are pointless as the hardware most optimal for arbitrary code execution is ... a CPU!


# Algorithm

The sequencer wants to publish a rollup. It picks a random big number `rollup_encryption_seed`. This number is put in the public rollup header that goes on chain.
To derive the encryption key, the rollup initializes an instance of `RandomX` hasher. The `rollup_encryption_seed` is fed into the hasher as the key that would generate the random hashing bytecode. 
Once the hasher is initialized, the sequencer needs to derive the `constrained input range` for the hasher. This would be done by using `rollup_encryption_seed` as the input for a `keccak256` hash. The resulting hash can be converted to a `uint256` number. 
The converted number is taken as the central point of the `constrained range`. We derieve the range to be `[hashResult-difficulty:hashResult+difficulty]`.
Now the sequencer picks another random number in the `constrained range`. This number is hashed with `RandomX`. The result is the `encryption_secret`.

Now the sequencer can use the `encryption_secret` to derieve an AES256 key and encrypt the rollup. Then its published as it would normally; The header contains the `rollup_encryption_seed` and a keccak256 hash of the `encryption_secret`.


## Bruteforcing

To bruteforce, we have to take the `rollup_encryption_seed` from the header published on chain and derive the `constrained_range` and `RandomX`. Then we have to go through each possible input in the `constrained_range` and hash it with `RandomX`; The result is then hashed through keccak256 and compared the published hash of the `encryption_secret` in the rollup header. This continues until we score a match. The preimage of the keccak256 we have is the `encryption_secret` for the rollup.

Bruteforcing can be done in parallel. With a special orchestrator it can be also done in parallel across many machines. There are up to date online benchmarks for all CPUs and their hash rate for `RandomX` with the monero default params, which we can reuse. Based on those benchmarks we can come up with the initial `difficulty` in terms of expected CPU duration to bruteforce and then convert the duration into a dollar cost based on cloud pricing for similair VMs. The cost should be something we can afford if we ever wanted to bruteforce, but it can also be higher and community subsidized. 

## Security Considerations

As we pick random number within the constrained range, it's possible the number is close to the start where a bruteforcing algo would start. This means that it could successfully decrypt the data faster than the revalation period, but as the range would be huge the odds of having the number in a convinient spot for bruteforcing would be small. As anybody attempting to bruteforce the number would have no idea where the number is, with a proper difficulty and random number distribution, a very huge majority of the time the bruteforcing would be completely pointless as the data would be revealed before its successful. This would in turn logically prevent people from even trying in the first place as it would just burn money aimlessly. 

If someone wanted to skip deriving through `RandomX`, they would have to pick the correct uint256 number in the full range from which the AES key is derived, which is way harder.

The `constrained range` and `RandomX` should be different based on the `rollup_encryption_seed` each time as otherwise one could precompute AES keys in the constrained range and store them. Having the hashing algo change each time, along with the range moving to random spots means that its impossible to prepare any speed ups. 

**The problem** - A state actor (or anybody rich and motivated enough) with virtually unlimited resources could of course bruteforce it in parallel in absolutely no time. A state actor also has resources for SGX sidechannel attacks and potentially the influence to get an attestation for a non SGX execution, so this doesn't really change the dynamic. **Note** that the state actor would need to bruteforce the key each time a new rollup is published as it would be different. If the revelation period is low and the dollar cost of bruteforcing high, this would make it extremely expensive to continiously monitor the network using the bruteforce approach. It's far more likely they will attack SGX.

# Incentivising bruteforcing

We can upgrade the protocol further by having a cut of the sequencer profits go into a bruteforcing pot. A reward from this pot can be claimed by the first address who proves to the contract they know the preimage of the `encryption_seed` hash. The reward will be decaying with the age of the rollup relative. We can make it so it fully decays when the revalation period passes, but we can also make it so it extends further than the revalation period. If we want it to extend further, of course the protocol should take into account how to not reveal the preimage when revealing the encryption key. 

**Why we want a decaying reward?** - It allows us to monitor the `difficulty` in real time. The biggest reward can be claimed by pushing out the preimage for the latest rollup. The timeframe to do so would be the lowest. If someone does it then we know that bruteforcing that `difficulty` takes the elapsed time between publishing the rollup and the bruteforce being announced. And potentially that it was economic to do it. Based on this we can increase the `difficulty` by whatever amount determined suitable. If we allow for rewards past the revalation period, then we can still enjoy the same monitoring, albeit people would probably do it at a loss due to the decaying reward.


We'd expect the pot to be claimed extremely rarely and people who do it to usually do so at a loss. But if computation power increased suddenly or if someone figured out how to do it quick and profitable, the pot would be claimed often and give us real time information if its time to reconsider the approach.

**Using the pot to fund recovery** - If we ever end up needing to recover state by bruteforcing, we can use the pot to incentivise 3rd parties to do it for us. We can add money to it if they find it unattractive or we can use it to pay for renting VMs ourselves and doing it. I imagine people would pool up and get the pot, rather than us having to do it as cloud services are expensive and there is a range where the incentive would carry a decent profit for the bruteforcers while end up being cheaper for us.

**Pot donations** - The pot can made to accept donations from other addresses. This can be useful in the event of a catastrophic failure as it allows interested parties to speed up the recovery process by providing incentive. If the `difficulty` is high, and the TVL locked is huge, its reasonable to expect donations for orderly recovery. 

We can also build a bruteforcing pool ourselves if the need ever arises and push it to the community and parties who are interested in gaining parts of the pot.