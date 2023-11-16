---
---
# Essential Information
Essential information, parameters and configuration settings for the Evan's Cat testnet.

## Limitations on Evan's Cat
These are the key limitations to be aware of when developing for the Evan's Cat testnet:

1. New nodes cannot be added. For now the number of Ten nodes is fixed.
1. Data revelation is not implemented yet.
1. Security is not fully implemented. Some keys are still hardcoded.
1. The decentralised bridge is limited to two hardcoded ERC20 tokens.
1. The Layer 1 is currently a hosted network. In the next iteration we'll connect Ten to an Ethereum testnet.
1. The "Wallet Extension" is not fully polished yet. You can expect a better UX as Ten develops

## Connection to an Ten Node
- **RPC http address:** `erpc.sepolia-testnet.obscu.ro:80`
- **RPC websocket address:** `erpc.sepolia-testnet.obscu.ro:81`

## Rollup Encryption/Decryption Key
The symmetric key used to encrypt and decrypt transaction blobs in rollups on the Ten Testnet:

```
bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c
```

N.B. Decrypting transaction blobs is only possible on testnet, where the rollup encryption key is long-lived and 
well-known. On mainnet, rollups will use rotating keys that are not known to anyone - or anything - other than the 
Ten enclaves.

## ERC20 Contracts
We have a couple of testnet ERC20 tokens (HOC & POC) that are automatically deployed with a static address every time 
testnet is restarted. The bridging mechanism has been setup for these tokens so they can be deposited from and withdrawn
to their respective ERC20 contracts on the L1.

Please contact us on [Discord](https://discord.gg/qA3FBYyZ) if you'd like some tokens to test with these contracts.

Hocus (HOC):

```
0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC
```

Pocus (POC):

```
0x9802F661d17c65527D7ABB59DAAD5439cb125a67
```