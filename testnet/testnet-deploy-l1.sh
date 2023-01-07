#!/usr/bin/env bash

#
# This script deploys contracts to the L1 of testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l1host=gethnetwork --pkstring=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  pkstring           *Required* Set the pkstring to deploy contracts"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
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
l1port=8025
docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
    
# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --pkstring)                 pkstring=${value} ;;
            --docker_image)             docker_image=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l1host:-} || -z ${pkstring:-}  ]];
then
    help_and_exit
fi

network_cfg='{ 
        "layer1" : {
            "url" : '"\"http://${l1host}:${l1port}\""',
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ "deploy_l1/" ],
            "accounts": [ "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" ]
        }
    }'

# deploy contracts to the geth network
echo "Creating docker network..."
docker network create --driver bridge node_network || true

echo "Deploying contracts to Layer 1 using obscuro hardhat container..."
docker run --name=hh-l1-deployer \
    --network=node_network \
    -e NETWORK_JSON="${network_cfg}" \
    -v deploymentsvol:/home/go-obscuro/contracts/deployments \
    "${docker_image}" \
    deploy \
    --network layer1
