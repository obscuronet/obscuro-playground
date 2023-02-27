package main

import (
	"fmt"
	"os"

	l2cd "github.com/obscuronet/go-obscuro/testnet/launcher/l2contractdeployer"
)

func main() {
	cliConfig := ParseConfigCLI()

	l2ContractDeployer, err := l2cd.NewDockerContractDeployer(
		l2cd.NewContractDeployerConfig(
			l2cd.WithL1Host(cliConfig.l1Host),                                    // "eth2network"
			l2cd.WithL1Port(cliConfig.l1HTTPPort),                                // 8025
			l2cd.WithL2Host(cliConfig.l2Host),                                    // "host"
			l2cd.WithL2WSPort(cliConfig.l2WSPort),                                // 13001
			l2cd.WithL1PrivateKey(cliConfig.privateKey),                          // "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
			l2cd.WithMessageBusContractAddress(cliConfig.messageBusContractAddr), // "0xFD03804faCA2538F4633B3EBdfEfc38adafa259B"
			l2cd.WithL2PrivateKey(cliConfig.l2PrivateKey),                        // "8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b"
			l2cd.WithHocPKString(cliConfig.l2HOCPrivateKey),                      // "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"),
			l2cd.WithPocPKString(cliConfig.l2POCPrivateKey),                      // "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"),
			l2cd.WithDockerImage(cliConfig.dockerImage),
		),
	)
	if err != nil {
		fmt.Println("unable to configure the l2 contract deployer - ", err)
		os.Exit(1)
	}

	err = l2ContractDeployer.Start()
	if err != nil {
		fmt.Println("unable to start the l2 contract deployer - ", err)
		os.Exit(1)
	}

	err = l2ContractDeployer.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for l2 contract deployer to finish - ", err)
		os.Exit(1)
	}
	fmt.Println("L2 Contracts were successfully deployed...")
	os.Exit(0)
}
