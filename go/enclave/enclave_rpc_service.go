package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/enclave/debugger"
	"github.com/ten-protocol/go-ten/go/enclave/events"
	"github.com/ten-protocol/go-ten/go/enclave/rpc"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/system"
	"github.com/ten-protocol/go-ten/go/responses"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type enclaveRPCService struct {
	rpcEncryptionManager *rpc.EncryptionManager
	registry             components.BatchRegistry
	subscriptionManager  *events.SubscriptionManager
	config               *enclaveconfig.EnclaveConfig
	debugger             *debugger.Debugger
	storage              storage.Storage
	crossChainProcessors *crosschain.Processors
	scb                  system.SystemContractCallbacks
}

func NewEnclaveRPCService(rpcEncryptionManager *rpc.EncryptionManager, registry components.BatchRegistry, subscriptionManager *events.SubscriptionManager, config *enclaveconfig.EnclaveConfig, debugger *debugger.Debugger, storage storage.Storage, crossChainProcessors *crosschain.Processors, scb system.SystemContractCallbacks) common.EnclaveClientRPC {
	return &enclaveRPCService{
		rpcEncryptionManager: rpcEncryptionManager,
		registry:             registry,
		subscriptionManager:  subscriptionManager,
		config:               config,
		debugger:             debugger,
		storage:              storage,
		crossChainProcessors: crossChainProcessors,
		scb:                  scb,
	}
}

func (e *enclaveRPCService) EncryptedRPC(ctx context.Context, encryptedParams common.EncryptedRequest) (*responses.EnclaveResponse, common.SystemError) {
	return rpc.HandleEncryptedRPC(ctx, e.rpcEncryptionManager, encryptedParams)
}

func (e *enclaveRPCService) GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash gethrpc.BlockNumberOrHash) ([]byte, common.SystemError) {
	stateDB, err := e.registry.GetBatchState(ctx, blockNrOrHash)
	if err != nil {
		return nil, responses.ToInternalError(fmt.Errorf("could not create stateDB. Cause: %w", err))
	}
	return stateDB.GetCode(address), nil
}

func (e *enclaveRPCService) Subscribe(ctx context.Context, id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) common.SystemError {
	encodedSubscription, err := e.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	return e.subscriptionManager.AddSubscription(id, encodedSubscription)
}

func (e *enclaveRPCService) Unsubscribe(id gethrpc.ID) common.SystemError {
	e.subscriptionManager.RemoveSubscription(id)
	return nil
}

func (e *enclaveRPCService) DebugTraceTransaction(ctx context.Context, txHash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	// ensure the debug namespace is enabled
	if !e.config.DebugNamespaceEnabled {
		return nil, responses.ToInternalError(fmt.Errorf("debug namespace not enabled"))
	}

	jsonMsg, err := e.debugger.DebugTraceTransaction(ctx, txHash, config)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return nil, responses.ToInternalError(err)
		}
		// TODO *Pedro* MOVE THIS TO Enclave Response
		return json.RawMessage(err.Error()), nil
	}

	return jsonMsg, nil
}

func (e *enclaveRPCService) GetTotalContractCount(ctx context.Context) (*big.Int, common.SystemError) {
	return e.storage.GetContractCount(ctx)
}

func (e *enclaveRPCService) EnclavePublicConfig(context.Context) (*common.EnclavePublicConfig, common.SystemError) {
	address, systemError := e.crossChainProcessors.GetL2MessageBusAddress()
	if systemError != nil {
		return nil, systemError
	}
	analyzerAddress := e.scb.TransactionPostProcessor()
	if analyzerAddress == nil {
		analyzerAddress = &gethcommon.Address{}
	}
	publicCallbacksAddress := e.scb.PublicCallbackHandler()
	if publicCallbacksAddress == nil {
		publicCallbacksAddress = &gethcommon.Address{}
	}

	return &common.EnclavePublicConfig{
		L2MessageBusAddress:             address,
		TransactionPostProcessorAddress: *analyzerAddress,
		PublicSystemContracts: map[string]gethcommon.Address{
			"PublicCallbacks": *publicCallbacksAddress,
		},
	}, nil
}