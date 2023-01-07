#!/usr/bin/env bash

#
# This script deploys contracts to testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l2host=testnet-host-1"
    echo ""
    echo "  l2host             *Required* Set the l2 host address"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  hocpkstring        *Optional* Set the pkstring to deploy HOC contract"
    echo ""
    echo "  pocpkstring        *Optional* Set the pkstring to deploy POC contract"
    echo ""
    echo "  l2port             *Optional* Set the l2 port. Defaults to 13001"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 8025"
    echo ""
    echo "  docker_image       *Optional* Sets the docker image to use. Defaults to testnetobscuronet.azurecr.io/obscuronet/contractdeployer:latest"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"

# Define defaults
l2port=13001
l1port=8025
# todo: get rid of these defaults and require them to be passed in, using github secrets for testnet values (requires bridge.go changes)
hocpkstring="6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
pocpkstring="4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
hocerc20address="0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l2host)                   l2host=${value} ;;
            --l2port)                   l2port=${value} ;;
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --hocpkstring)              hocpkstring=${value} ;;
            --pocpkstring)              pocpkstring=${value} ;;
            --docker_image)             docker_image=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l2host:-} || -z ${hocpkstring:-} || -z ${pocpkstring:-} ||  -z ${l1host:-} ]];
then
    help_and_exit
fi

network_cfg='{ 
        "layer2" : {
            "url" : "http://127.0.0.1:3000",
            "obscuroEncRpcUrl" : '"\"http://${l2host}:${l2port}\""',
            "live" : false,
            "saveDeployments" : true,
            "companionNetworks" : {
                "layer1": "layer1"
            },
            "deploy": [ "deploy_l2/" ],
            "accounts": [ "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b" ]
        },
        "layer1" : {
            "url" : '"\"http://${l1host}:${l1port}\""',
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ "deploy_l1/" ],
            "accounts": [ "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" ]
        }
    }'

# deploy contracts to the obscuro network
echo "Creating network for container..."
docker network create --driver bridge node_network || true

echo "Deploying contracts to layer 2 using hardhat container..."
docker run --name=hh-l2-deployer \
    --network=node_network \
    -e NETWORK_JSON="${network_cfg}" \
    -v deploymentsvol:/home/go-obscuro/contracts/deployments \
    "${docker_image}" \
    obscuro:deploy \
    --network layer2 