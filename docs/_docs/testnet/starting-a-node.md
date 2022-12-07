---
---
# Starting a node
How to start a node in the Evan's Cat testnet.

## Requirements
- SGX enabled VM
- Docker

## Steps
#### - Create an SGX enabled VM
Recommended Standard DC4s v2 (4 vcpus, 16 GiB memory) in Azure.

#### - Install Docker

```
sudo apt-get update \
    && sudo apt-get install -y jq \
    && curl -fsSL https://get.docker.com -o get-docker.sh && sh ./get-docker.sh
```

#### - Download Obscuro Start script


Make sure to use the latest `<version>` at https://github.com/obscuronet/go-obscuro/tags.

```
  wget https://raw.githubusercontent.com/obscuronet/go-obscuro/<version>/testnet/start-obscuro-node.sh && \
  wget https://raw.githubusercontent.com/obscuronet/go-obscuro/<version>/testnet/docker-compose.yml
```

#### - Start Obscuro Node

To start the obscuro node some information is required to populate the starting script.

- (p2p_public_address) The external facing address of the network. Where outside peers will connect to. Must be open to outside connections.`curl https://ipinfo.io/ip` provides the external IP.
- (pkstring) Private Key to issue transactions into the Layer 1
- (host_id) Public Key derived from the Private Key

```
chmod +x start-obscuro-node.sh \
&& sudo ./start-obscuro-node.sh  \
  --is_genesis=false \
  --node_type="validator" \
  --sgx_enabled="true" \
  --host_id=PublicKeyAddress \
  --l1host="testnet-gethnetwork.uksouth.azurecontainer.io" \
  --mgmtcontractaddr="0xF886d9e52d38f3C7BEd96B1F45e366C459886694" \
  --hocerc20addr="0xc559903C00ed55d43021cf4ea49f5924BF8b5A4B" \
  --pocerc20addr="0xB46213b1755545261Ce32e8b46B300fB01663889" \
  --pkstring="PrivateKeyString" \
  --p2p_public_address="HOST:10000"