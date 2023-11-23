package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ten-protocol/go-ten/integration/eth2network"
)

// Spins up a new eth 2 network.
func main() {
	config := parseCLIArgs()

	fmt.Printf("Starting eth2network with params: %+v\n", config)

	binariesPath, err := eth2network.EnsureBinariesExist()
	if err != nil {
		panic(err)
	}

	if config.onlyDownload {
		os.Exit(0)
	}

	eth2Network := eth2network.NewEth2Network(
		binariesPath,
		config.logToFile,
		config.gethHTTPStartPort,
		config.gethWSStartPort,
		config.gethAuthRPCStartPort,
		config.gethNetworkStartPort,
		config.prysmBeaconRPCStartPort,
		config.prysmBeaconRPCStartPort,
		config.chainID,
		config.numNodes,
		config.blockTimeSecs,
		config.slotsPerEpoch,
		config.secondsPerSlot,
		config.prefundedAddrs,
		5*time.Minute,
	)

	err = eth2Network.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("eth 2 network started..")

	handleInterrupt(eth2Network)
}

// Shuts down the Geth network when an interrupt is received.
func handleInterrupt(network eth2network.Eth2Network) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	err := network.Stop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("eth2 network stopping...")
	os.Exit(0)
}
