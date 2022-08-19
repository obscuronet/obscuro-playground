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
    echo "  jampkstring        *Optional* Set the pkstring to deploy JAM contract"
    echo ""
    echo "  ethpkstring        *Optional* Set the pkstring to deploy ETH contract"
    echo ""
    echo "  l2port             *Optional* Set the l2 port. Defaults to 10000"
    echo ""
    echo "  docker_image       *Optional* Sets the docker image to use. Defaults to testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest"
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
l2port=13000
# todo: get rid of these defaults and require them to be passed in, using github secrets for testnet values (requires bridge.go changes)
jampkstring="6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
ethpkstring="4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
jamerc20address="0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
docker_image="testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l2host)                   l2host=${value} ;;
            --l2port)                   l2port=${value} ;;
            --jampkstring)              jampkstring=${value} ;;
            --ethpkstring)              ethpkstring=${value} ;;
            --docker_image)             docker_image=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l2host:-} || -z ${jampkstring:-} || -z ${ethpkstring:-}  ]];
then
    help_and_exit
fi

# deploy contracts to the obscuro network
echo "Deploying JAM ERC20 contract to the obscuro network..."
docker network create --driver bridge node_network || true
docker run --name=jamL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="L2ERC20" \
    --privateKey=${jampkstring}\
    --constructorParams="JAM,JAM,1000000000000000000000000000000"
echo ""

echo "Deploying ETH ERC20 contract to the obscuro network..."
docker run --name=ethL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="L2ERC20" \
    --privateKey=${ethpkstring}\
    --constructorParams="ETH,ETH,1000000000000000000000000000000"
echo ""

echo "Deploying Guessing game contract to the obscuro network..."
docker run --name=guessingGameL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="GUESS" \
    --privateKey=${ethpkstring}\
    --constructorParams="100,${jamerc20address}"
echo ""
