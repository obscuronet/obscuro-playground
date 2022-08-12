package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common/log"
)

const (
	logPattern    = ">   Agg%d: %s"
	txExecPattern = "!Tx %s: %s"
)

// LogWithID logs a message at INFO level with the aggregator's identity prepended.
func LogWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Info(logPattern, nodeID, formattedMsg)
}

// WarnWithID logs a message at WARN level with the aggregator's identity prepended.
func WarnWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Warn(logPattern, nodeID, formattedMsg)
}

// TraceWithID logs a message at TRACE level with the aggregator's identity prepended.
func TraceWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Trace(logPattern, nodeID, formattedMsg)
}

// ErrorWithID logs a message at ERROR level with the aggregator's identity prepended.
func ErrorWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Error(logPattern, nodeID, formattedMsg)
}

// PanicWithID logs a message at PANIC level with the aggregator's identity prepended.
func PanicWithID(nodeID uint64, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Panic(logPattern, nodeID, formattedMsg)
}

// LogTXExecution - logs at INFO level with the tx hash prepended
func LogTXExecution(txHash common.Hash, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Info(txExecPattern, txHash.Hex(), formattedMsg)
}

// TraceTXExecution - logs at Trace level with the tx hash prepended
func TraceTXExecution(txHash common.Hash, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Trace(txExecPattern, txHash.Hex(), formattedMsg)
}

// ErrorTXExecution - logs at Error level with the tx hash prepended
func ErrorTXExecution(txHash common.Hash, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	log.Error(txExecPattern, txHash.Hex(), formattedMsg)
}
