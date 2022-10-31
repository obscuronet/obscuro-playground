package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/obscuronet/go-obscuro/integration/gethnetwork"
)

// Spins up a new Geth network.
func main() {
	config := parseCLIArgs()

	gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
	if err != nil {
		panic(err)
	}
	gethNetwork := gethnetwork.NewGethNetwork(config.startPort, config.websocketStartPort, gethBinaryPath, config.numNodes, config.blockTimeSecs, config.prefundedAddrs, config.logPath, config.logLevel)
	fmt.Println("Geth network started.")

	handleInterrupt(gethNetwork)
}

// Shuts down the Geth network when an interrupt is received.
func handleInterrupt(gethNetwork *gethnetwork.GethNetwork) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	gethNetwork.StopNodes()
	fmt.Println("Geth network stopping...")
	os.Exit(1)
}
