package main

import (
	"fmt"
	"os"

	l1cp "github.com/ten-protocol/go-ten/testnet/launcher/l1challengeperiod"
)

func main() {
	cliConfig := ParseConfigCLI()

	l1challengeperiod, err := l1cp.NewSetChallengePeriod(
		l1cp.NewChallengePeriodConfig(
			l1cp.WithL1HTTPURL(cliConfig.l1HTTPURL),
			l1cp.WithPrivateKey(cliConfig.privateKey),
			l1cp.WithDockerImage(cliConfig.dockerImage),
			l1cp.WithMgmtContractAddress(cliConfig.mgmtContractAddress),
			l1cp.WithChallengePeriod(cliConfig.challengePeriod),
		),
	)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1challengeperiod.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1challengeperiod.WaitForFinish()
	if err != nil {
		fmt.Println("unexpected error waiting for set challenge period script to finish - %w", err)
		os.Exit(1)
	}
	fmt.Println("L1 challenge period was successfully set...")
	os.Exit(0)
}
