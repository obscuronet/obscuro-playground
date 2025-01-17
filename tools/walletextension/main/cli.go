package main

import (
	"flag"
	"fmt"
	"time"

	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	// Flag names, defaults and usages.
	walletExtensionHostName    = "host"
	walletExtensionHostDefault = "127.0.0.1"
	walletExtensionHostUsage   = "The host where the wallet extension should open the port."
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension. Default: 3000."

	walletExtensionPortWSName    = "portWS"
	walletExtensionPortWSDefault = 3001
	walletExtensionPortWSUsage   = "The port on which to serve websocket JSON RPC requests. Default: 3001."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "erpc.sepolia-testnet.ten.xyz"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `erpc.sepolia-testnet.ten.xyz`."

	nodeHTTPPortName    = "nodePortHTTP"
	nodeHTTPPortDefault = 80
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 80."

	nodeWebsocketPortName    = "nodePortWS"
	nodeWebsocketPortDefault = 81
	nodeWebsocketPortUsage   = "The port on which to connect to the Obscuro node via RPC over websockets. Default: 81."

	logPathName    = "logPath"
	logPathDefault = "sys_out"
	logPathUsage   = "The path to use for the wallet extension's log file"

	databasePathName    = "databasePath"
	databasePathDefault = ".obscuro/gateway_database.db"
	databasePathUsage   = "The path for the wallet extension's database file. Default: .obscuro/gateway_database.db"

	verboseFlagName    = "verbose"
	verboseFlagDefault = false
	verboseFlagUsage   = "Flag to enable verbose logging of wallet extension traffic"

	dbTypeFlagName    = "dbType"
	dbTypeFlagDefault = "sqlite"
	dbTypeFlagUsage   = "Defined the db type (sqlite or mariaDB)"

	dbConnectionURLFlagName    = "dbConnectionURL"
	dbConnectionURLFlagDefault = ""
	dbConnectionURLFlagUsage   = "If dbType is set to mariaDB, this must be set. ex: obscurouser:password@tcp(127.0.0.1:3306)/ogdb"

	tenChainIDName      = "tenChainID"
	tenChainIDDefault   = 443
	tenChainIDFlagUsage = "ChainID of TEN network that the gateway is communicating with"

	storeIncomingTxs        = "storeIncomingTxs"
	storeIncomingTxsDefault = true
	storeIncomingTxsUsage   = "Flag to enable storing incoming transactions in the database for debugging purposes. Default: true"

	rateLimitUserComputeTimeName    = "rateLimitUserComputeTime"
	rateLimitUserComputeTimeDefault = 10 * time.Second
	rateLimitUserComputeTimeUsage   = "rateLimitUserComputeTime represents how much compute time is user allowed to used in rateLimitWindow time. If rateLimitUserComputeTime is set to 0, rate limiting is turned off. Default: 10s."

	rateLimitWindowName    = "rateLimitWindow"
	rateLimitWindowDefault = 1 * time.Minute
	rateLimitWindowUsage   = "rateLimitWindow represents time window in which we allow one user to use compute time defined with rateLimitUserComputeTimeMs  Default: 1m"

	rateLimitMaxConcurrentRequestsName    = "maxConcurrentRequestsPerUser"
	rateLimitMaxConcurrentRequestsDefault = 3
	rateLimitMaxConcurrentRequestsUsage   = "Number of concurrent requests allowed per user. Default: 3"

	insideEnclaveFlagName    = "insideEnclave"
	insideEnclaveFlagDefault = false
	insideEnclaveFlagUsage   = "Flag to indicate if the program is running inside an enclave. Default: false"

	keyExchangeURLFlagName    = "keyExchangeURL"
	keyExchangeURLFlagDefault = ""
	keyExchangeURLFlagUsage   = "URL to exchange the key with another enclave. Default: empty"

	enableTLSFlagName    = "enableTLS"
	enableTLSFlagDefault = false
	enableTLSFlagUsage   = "Flag to enable TLS/HTTPS"

	tlsDomainFlagName    = "tlsDomain"
	tlsDomainFlagDefault = ""
	tlsDomainFlagUsage   = "Domain name for TLS certificate"

	encryptingCertificateEnabledFlagName    = "encryptingCertificateEnabled"
	encryptingCertificateEnabledFlagDefault = false
	encryptingCertificateEnabledFlagUsage   = "Flag to enable encrypting certificate functionality. Default: false"

	disableCachingFlagName    = "disableCaching"
	disableCachingFlagDefault = false
	disableCachingFlagUsage   = "Flag to disable response caching in the gateway. Default: false"
)

func parseCLIArgs() wecommon.Config {
	walletExtensionHost := flag.String(walletExtensionHostName, walletExtensionHostDefault, walletExtensionHostUsage)
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	walletExtensionPortWS := flag.Int(walletExtensionPortWSName, walletExtensionPortWSDefault, walletExtensionPortWSUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	nodeWebsocketPort := flag.Int(nodeWebsocketPortName, nodeWebsocketPortDefault, nodeWebsocketPortUsage)
	logPath := flag.String(logPathName, logPathDefault, logPathUsage)
	databasePath := flag.String(databasePathName, databasePathDefault, databasePathUsage)
	verboseFlag := flag.Bool(verboseFlagName, verboseFlagDefault, verboseFlagUsage)
	dbType := flag.String(dbTypeFlagName, dbTypeFlagDefault, dbTypeFlagUsage)
	dbConnectionURL := flag.String(dbConnectionURLFlagName, dbConnectionURLFlagDefault, dbConnectionURLFlagUsage)
	tenChainID := flag.Int(tenChainIDName, tenChainIDDefault, tenChainIDFlagUsage)
	storeIncomingTransactions := flag.Bool(storeIncomingTxs, storeIncomingTxsDefault, storeIncomingTxsUsage)
	rateLimitUserComputeTime := flag.Duration(rateLimitUserComputeTimeName, rateLimitUserComputeTimeDefault, rateLimitUserComputeTimeUsage)
	rateLimitWindow := flag.Duration(rateLimitWindowName, rateLimitWindowDefault, rateLimitWindowUsage)
	rateLimitMaxConcurrentRequests := flag.Int(rateLimitMaxConcurrentRequestsName, rateLimitMaxConcurrentRequestsDefault, rateLimitMaxConcurrentRequestsUsage)
	insideEnclaveFlag := flag.Bool(insideEnclaveFlagName, insideEnclaveFlagDefault, insideEnclaveFlagUsage)
	keyExchangeURL := flag.String(keyExchangeURLFlagName, keyExchangeURLFlagDefault, keyExchangeURLFlagUsage)
	enableTLSFlag := flag.Bool(enableTLSFlagName, enableTLSFlagDefault, enableTLSFlagUsage)
	tlsDomainFlag := flag.String(tlsDomainFlagName, tlsDomainFlagDefault, tlsDomainFlagUsage)
	encryptingCertificateEnabled := flag.Bool(encryptingCertificateEnabledFlagName, encryptingCertificateEnabledFlagDefault, encryptingCertificateEnabledFlagUsage)
	disableCaching := flag.Bool(disableCachingFlagName, disableCachingFlagDefault, disableCachingFlagUsage)
	flag.Parse()

	return wecommon.Config{
		WalletExtensionHost:            *walletExtensionHost,
		WalletExtensionPortHTTP:        *walletExtensionPort,
		WalletExtensionPortWS:          *walletExtensionPortWS,
		NodeRPCHTTPAddress:             fmt.Sprintf("%s:%d", *nodeHost, *nodeHTTPPort),
		NodeRPCWebsocketAddress:        fmt.Sprintf("%s:%d", *nodeHost, *nodeWebsocketPort),
		LogPath:                        *logPath,
		DBPathOverride:                 *databasePath,
		VerboseFlag:                    *verboseFlag,
		DBType:                         *dbType,
		DBConnectionURL:                *dbConnectionURL,
		TenChainID:                     *tenChainID,
		StoreIncomingTxs:               *storeIncomingTransactions,
		RateLimitUserComputeTime:       *rateLimitUserComputeTime,
		RateLimitWindow:                *rateLimitWindow,
		RateLimitMaxConcurrentRequests: *rateLimitMaxConcurrentRequests,
		InsideEnclave:                  *insideEnclaveFlag,
		KeyExchangeURL:                 *keyExchangeURL,
		EnableTLS:                      *enableTLSFlag,
		TLSDomain:                      *tlsDomainFlag,
		EncryptingCertificateEnabled:   *encryptingCertificateEnabled,
		DisableCaching:                 *disableCaching,
	}
}
