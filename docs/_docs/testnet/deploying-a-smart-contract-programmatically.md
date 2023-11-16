---
---
# Deploying a Smart Contract to Ten Testnet Programmatically
The steps below demonstrate how to programmatically create a new contract on to Ten Testnet and interact with it via 
call functions. The example uses [Python](https://www.python.org/) and [web3.py](https://web3py.readthedocs.io/en/stable/) 
as a reference but the principles of usage will be the same in any web3 language implementation. 

A full working example can be seen in [deploying-a-smart-contract-programmatically.py](https://github.com/ten-protocol/go-ten/blob/main/docs/_docs/testnet/deploying-a-smart-contract-programmatically.py).
Usage of the example requires Python > 3.9.13, solc 0.8.15 and the web3, requests, and json modules. It is assumed solc 
has been installed using homebrew and resides in `/opt/homebrew/bin/solc` and that the wallet extension is running on 
the local host with default values `WHOST=127.0.0.1` and `WPORT=3000`.

A walk through and explanation of the steps performed is given below;

## Connect to the network and create a local private key
The [wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension/) acts as an HTTP server to mediate RPC requests. In the 
below a connection is made on the wallet extension host and port, a private key is locally created and the associated 
account stored for later usage.

```python
    w3 = Web3(Web3.HTTPProvider('http://%s:%d' % (WHOST, WPORT)))
    private_key = secrets.token_hex(32)
    account = w3.eth.account.privateKeyToAccount(private_key)
    logging.info('Using account with address %s' % account.address)
```

## Request ETH from the faucet server for native ETH
An account needs gas to perform transactions on Ten, where gas is paid in native ETH. Requests of native ETH can be 
made through a POST to the faucet server where the address is supplied in the data payload.
```python
    headers = {'Content-Type': 'application/json'}
    data = {"address": account.address}
    requests.post(FAUCET_URL, data=json.dumps(data), headers=headers)
```

## Generate a viewing key, sign and post back to the wallet extension
The enclave encodes all communication to the wallet extension using viewing keys. HTTP endpoints exist in the wallet 
extension to facilitate requesting a viewing key, and to sign and return it to the enclave. 
```python 
    headers = {'Accept': 'application/json', 'Content-Type': 'application/json'}
    
    data = {"address": account.address}
    response = requests.post('http://%s:%d/generateviewingkey/' % (WHOST, WPORT), data=json.dumps(data), headers=headers)
    signed_msg = w3.eth.account.sign_message(encode_defunct(text='vk' + response.text), private_key=private_key)
    data = {"signature": signed_msg.signature.hex(), "address": account.address}
    response = requests.post('http://%s:%d/submitviewingkey/' % (WHOST, WPORT), data=json.dumps(data), headers=headers)
```

## Compile the contract and build the local deployment transaction
A contract can be compiled using solc and a transaction created for deploying the contract. Construction of the transaction 
requires `gasPrice` and `gas` to be explicitly defined (the need to perform this will be removed in a later 
release). An arbitrary `gasPrice` should be given e.g. the current price on the Ropsten test network. 
```python 
    compiled_sol = compile_source(guesser, output_values=['abi', 'bin'], solc_binary='/opt/homebrew/bin/solc')
    contract_id, contract_interface = compiled_sol.popitem()
    bytecode = contract_interface['bin']
    abi = contract_interface['abi']
    contract = w3.eth.contract(abi=abi, bytecode=bytecode)
    build_tx = contract.constructor(random.randrange(0, 100)).buildTransaction(
        {
            'from': account.address,
            'nonce': w3.eth.getTransactionCount(account.address),
            'gasPrice': 1499934385,
            'gas': 720000,
            'chainId': 443
        }
    )
```

## Sign the transaction and send to the network 
Using the account the transaction can be signed and submitted to the Ten Testnet. 
```python
    signed_tx = account.signTransaction(build_tx)
    tx_hash = None
    try:
        tx_hash = w3.eth.sendRawTransaction(signed_tx.rawTransaction)
    except Exception as e:
        logging.error('Error sending raw transaction %s' % e)
        return
```

## Wait for the transaction receipt 
Once submitted the transaction receipt can be obtained in order to get the deployed contract address. 
```python
    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    if tx_receipt.status == 0:
        logging.error('Transaction receipt has failed status ... aborting')
        return
    else:
        logging.info('Received transaction receipt')
```

## Create the contract using the abi and contract address
Once the transaction receipt is received function calls can be made against the contract. 
```python
    contract = w3.eth.contract(address=tx_receipt.contractAddress, abi=abi)
    contract.functions.guess(guess).call()
```
