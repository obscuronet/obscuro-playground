#!/usr/bin/env bash

#
# This script starts a local testnet using the recommended / default arguments
#

echo Building the required docker images
./testnet-local-build_images.sh

echo Starting up the L1 network
./testnet-local-eth2network.sh --pkaddresses=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944,0x0654D8B60033144D567f25bF41baC1FB0D60F23B

echo Deploying the l1 contracts
./testnet-deploy-contracts.sh --l1host=-eth2network -pkstring=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb

echo Starting up the Obscuro node
./start-obscuro-node.sh --sgx_enabled=false --l1host=eth2network --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --hocerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --pocerc20addr=0x51D43a3Ca257584E770B6188232b199E76B022A2 --is_genesis=true --node_type=sequencer

echo Deploying the L2 contracts
./testnet-deploy-l2-contracts.sh --l2host=testnet-host-1 --l1host=eth2network

echo Starting obscuroscan
./start-obscuroscan.sh --rpcServerAddress=http://testnet-host-1:13000 --receivingPort=8098

