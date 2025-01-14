package log

import (
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/lib/gethfork/debug"
)

// These are the keys of the log entries
const (
	ErrKey           = "err"
	CtrErrKey        = "ctr_err"
	SubIDKey         = "subscription_id"
	CfgKey           = "cfg"
	TxKey            = "tx"
	DurationKey      = "duration"
	DurationMilliKey = "durationMilli"
	BundleHashKey    = "bundle"
	BatchHashKey     = "batch"
	BatchHeightKey   = "batch_height"
	BatchSeqNoKey    = "batch_seq_num"
	RollupHashKey    = "rollup"
	CmpKey           = "cmp"
	NodeIDKey        = "node_id"
	EnclaveIDKey     = "enclave_id"
	NetworkIDKey     = "network_id"
	BlockHeightKey   = "block_height"
	BlockHashKey     = "block_hash"
	PackageKey       = "package"
)

// Logging is grouped by the component where it was initialised
const (
	EnclaveCmp      = "enclave"
	HostCmp         = "host"
	HostRPCCmp      = "host_rpc"
	TxInjectCmp     = "tx_inject"
	P2PCmp          = "p2p"
	RPCClientCmp    = "rpc_client"
	DeployerCmp     = "deployer"
	NetwMngCmp      = "network_manager"
	WalletExtCmp    = "wallet_extension"
	TestGethNetwCmp = "test_geth_network"
	EthereumL1Cmp   = "l1_host"
	TenscanCmp      = "tenscan"
	CrossChainCmp   = "cross_chain"
)

// SysOut - Used when the logger has to write to Sys.out
const (
	SysOut = "sys_out"
)

// New - helper function used to create a top level logger for a component.
// Note: this expects legacy geth log levels, you will get unexpected behaviour if you use gethlog.<LEVEL> directly.
func New(component string, level int, out string, ctx ...interface{}) gethlog.Logger {
	logFile := ""
	if out != SysOut {
		logFile = out
	}
	verbosity := gethlog.FromLegacyLevel(level)

	err := debug.Setup("terminal", logFile, false, 0, 0, 0, false, false, verbosity, "")
	if err != nil {
		panic(err.Error())
	}

	context := append(ctx, CmpKey, component)
	l := gethlog.New(context...)

	return l
}
