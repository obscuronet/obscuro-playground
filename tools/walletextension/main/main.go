package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/tools/walletextension"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	// @fixme -
	// this is a temporary fix as out forked version of log.go does not map with gethlog.Level<Level>
	// and should be fixed as part of logging refactoring in the future
	legacyLevelDebug = 4
	legacyLevelError = 1
)

func main() {
	config := parseCLIArgs()
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")

	// Setup logging first
	if config.LogPath != log.SysOut {
		_, err := os.Create(config.LogPath)
		if err != nil {
			panic(fmt.Sprintf("could not create log file. Cause: %s", err))
		}
	}
	logLvl := legacyLevelError
	if config.VerboseFlag {
		logLvl = legacyLevelDebug
	}
	logger := log.New(log.WalletExtCmp, logLvl, config.LogPath)

	logger.Info("Welcome to the TEN gateway")
	logger.Info("Starting with following config", "config", string(jsonConfig))

	// Start the wallet extension right away
	walletExtContainer := walletextension.NewContainerFromConfig(config, logger)
	err := walletExtContainer.Start()
	if err != nil {
		logger.Error("Failed to start wallet extension", "error", err)
		os.Exit(1)
	}

	walletExtensionAddr := fmt.Sprintf("%s:%d", common.Localhost, config.WalletExtensionPortHTTP)
	logger.Info("TEN gateway started", "url", fmt.Sprintf("http://%s/v1/network-config", walletExtensionAddr))

	select {}
}
