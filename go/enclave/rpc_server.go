package enclave

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/config"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"google.golang.org/grpc"
)

// Receives RPC calls to the enclave process and relays them to the enclave.Enclave.
type server struct {
	generated.UnimplementedEnclaveProtoServer
	enclave     common.Enclave
	rpcServer   *grpc.Server
	nodeShortID uint64
}

// StartServer starts a server on the given port on a separate thread. It creates an enclave.Enclave for the provided nodeID,
// and uses it to respond to incoming RPC messages from the host.
func StartServer(
	enclaveConfig config.EnclaveConfig,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	collector StatsCollector,
) (func(), error) {
	lis, err := net.Listen("tcp", enclaveConfig.Address)
	if err != nil {
		return nil, fmt.Errorf("enclave RPC server could not listen on port: %w", err)
	}

	enclaveServer := server{
		enclave:     NewEnclave(enclaveConfig, mgmtContractLib, erc20ContractLib, collector),
		rpcServer:   grpc.NewServer(),
		nodeShortID: common.ShortAddress(enclaveConfig.HostID),
	}
	generated.RegisterEnclaveProtoServer(enclaveServer.rpcServer, &enclaveServer)

	go func(lis net.Listener) {
		common.LogWithID(enclaveServer.nodeShortID, "Enclave server listening on address %s.", enclaveConfig.Address)
		err = enclaveServer.rpcServer.Serve(lis)
		if err != nil {
			common.LogWithID(enclaveServer.nodeShortID, "enclave RPC server could not serve: %s", err)
		}
	}(lis)

	closeHandle := func() {
		go enclaveServer.Stop(context.Background(), nil) //nolint:errcheck
	}

	return closeHandle, nil
}

// Status returns the current status of the server as an enum value (see common.Status for details)
func (s *server) Status(context.Context, *generated.StatusRequest) (*generated.StatusResponse, error) {
	errStr := ""
	status, err := s.enclave.Status()
	if err != nil {
		errStr = err.Error()
	}
	return &generated.StatusResponse{Status: int32(status), Error: errStr}, nil
}

func (s *server) Attestation(context.Context, *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation, err := s.enclave.Attestation()
	if err != nil {
		return nil, err
	}
	msg := rpc.ToAttestationReportMsg(attestation)
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *server) GenerateSecret(context.Context, *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret, err := s.enclave.GenerateSecret()
	if err != nil {
		return nil, err
	}
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) ShareSecret(_ context.Context, request *generated.FetchSecretRequest) (*generated.ShareSecretResponse, error) {
	attestationReport := rpc.FromAttestationReportMsg(request.AttestationReportMsg)
	secret, err := s.enclave.ShareSecret(attestationReport)
	return &generated.ShareSecretResponse{EncryptedSharedEnclaveSecret: secret}, err
}

func (s *server) InitEnclave(_ context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	errStr := ""
	if err := s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret); err != nil {
		errStr = err.Error()
	}
	return &generated.InitEnclaveResponse{Error: errStr}, nil
}

func (s *server) ProduceGenesis(_ context.Context, request *generated.ProduceGenesisRequest) (*generated.ProduceGenesisResponse, error) {
	genesisRollup, err := s.enclave.ProduceGenesis(gethcommon.BytesToHash(request.GetBlockHash()))
	if err != nil {
		return nil, err
	}

	blockSubmissionResponse, err := rpc.ToBlockSubmissionResponseMsg(genesisRollup)
	if err != nil {
		return nil, err
	}

	return &generated.ProduceGenesisResponse{BlockSubmissionResponse: &blockSubmissionResponse}, nil
}

func (s *server) Start(_ context.Context, request *generated.StartRequest) (*generated.StartResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	err := s.enclave.Start(bl)
	if err != nil {
		return nil, err
	}
	return &generated.StartResponse{}, nil
}

func (s *server) SubmitBlock(_ context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	blockSubmissionResponse, err := s.enclave.SubmitBlock(bl, request.IsLatest)
	if err != nil {
		return nil, err
	}

	msg, err := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
}

func (s *server) SubmitRollup(_ context.Context, request *generated.SubmitRollupRequest) (*generated.SubmitRollupResponse, error) {
	extRollup := rpc.FromExtRollupMsg(request.ExtRollup)
	err := s.enclave.SubmitRollup(extRollup)
	if err != nil {
		return nil, err
	}
	return &generated.SubmitRollupResponse{}, nil
}

func (s *server) SubmitTx(_ context.Context, request *generated.SubmitTxRequest) (*generated.SubmitTxResponse, error) {
	encryptedHash, err := s.enclave.SubmitTx(request.EncryptedTx)
	return &generated.SubmitTxResponse{EncryptedHash: encryptedHash}, err
}

func (s *server) ExecuteOffChainTransaction(_ context.Context, request *generated.OffChainRequest) (*generated.OffChainResponse, error) {
	result, err := s.enclave.ExecuteOffChainTransaction(request.EncryptedParams)
	var errorJSON []byte
	var merr error
	if err != nil {
		// serialise the error object returned by the evm into a json
		errorJSON, merr = json.Marshal(err)
		if merr != nil {
			return nil, merr
		}
		return &generated.OffChainResponse{Error: errorJSON}, nil
	}
	return &generated.OffChainResponse{Result: result}, nil
}

func (s *server) GetTransactionCount(_ context.Context, request *generated.GetTransactionCountRequest) (*generated.GetTransactionCountResponse, error) {
	result, err := s.enclave.GetTransactionCount(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionCountResponse{Result: result}, nil
}

func (s *server) RoundWinner(_ context.Context, request *generated.RoundWinnerRequest) (*generated.RoundWinnerResponse, error) {
	extRollup, winner, err := s.enclave.RoundWinner(gethcommon.BytesToHash(request.Parent))
	if err != nil {
		return nil, err
	}
	extRollupMsg := rpc.ToExtRollupMsg(&extRollup)
	return &generated.RoundWinnerResponse{Winner: winner, ExtRollup: &extRollupMsg}, nil
}

func (s *server) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	defer s.rpcServer.GracefulStop()
	err := s.enclave.Stop()
	return &generated.StopResponse{}, err
}

func (s *server) GetTransaction(_ context.Context, request *generated.GetTransactionRequest) (*generated.GetTransactionResponse, error) {
	encryptedTx, err := s.enclave.GetTransaction(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionResponse{EncryptedTx: encryptedTx}, nil
}

func (s *server) GetTransactionReceipt(_ context.Context, request *generated.GetTransactionReceiptRequest) (*generated.GetTransactionReceiptResponse, error) {
	encryptedTxReceipt, err := s.enclave.GetTransactionReceipt(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionReceiptResponse{EncryptedTxReceipt: encryptedTxReceipt}, nil
}

func (s *server) GetRollup(_ context.Context, request *generated.GetRollupRequest) (*generated.GetRollupResponse, error) {
	extRollup, err := s.enclave.GetRollup(gethcommon.BytesToHash(request.RollupHash))
	if err != nil {
		return nil, err
	}

	extRollupMsg := rpc.ToExtRollupMsg(extRollup)
	return &generated.GetRollupResponse{ExtRollup: &extRollupMsg}, nil
}

func (s *server) AddViewingKey(_ context.Context, request *generated.AddViewingKeyRequest) (*generated.AddViewingKeyResponse, error) {
	err := s.enclave.AddViewingKey(request.ViewingKey, request.Signature)
	if err != nil {
		return nil, err
	}
	return &generated.AddViewingKeyResponse{}, nil
}

func (s *server) GetBalance(_ context.Context, request *generated.GetBalanceRequest) (*generated.GetBalanceResponse, error) {
	encryptedBalance, err := s.enclave.GetBalance(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetBalanceResponse{EncryptedBalance: encryptedBalance}, nil
}

func (s *server) GetCode(_ context.Context, request *generated.GetCodeRequest) (*generated.GetCodeResponse, error) {
	address := gethcommon.BytesToAddress(request.Address)
	rollupHash := gethcommon.BytesToHash(request.RollupHash)

	code, err := s.enclave.GetCode(address, &rollupHash)
	if err != nil {
		return nil, err
	}
	return &generated.GetCodeResponse{Code: code}, nil
}

func (s *server) StoreAttestation(_ context.Context, req *generated.StoreAttestationRequest) (*generated.StoreAttestationResponse, error) {
	err := s.enclave.StoreAttestation(rpc.FromAttestationReportMsg(req.AttestationReportMsg))
	resp := ""
	if err != nil {
		resp = err.Error()
	}
	return &generated.StoreAttestationResponse{Error: resp}, nil
}

func (s *server) Subscribe(_ context.Context, req *generated.SubscribeRequest) (*generated.SubscribeResponse, error) {
	err := s.enclave.Subscribe(gethrpc.ID(req.Id), req.EncryptedSubscription)
	return &generated.SubscribeResponse{}, err
}

func (s *server) Unsubscribe(_ context.Context, req *generated.UnsubscribeRequest) (*generated.UnsubscribeResponse, error) {
	err := s.enclave.Unsubscribe(gethrpc.ID(req.Id))
	return &generated.UnsubscribeResponse{}, err
}

func (s *server) EstimateGas(_ context.Context, req *generated.EstimateGasRequest) (*generated.EstimateGasResponse, error) {
	encryptedBalance, err := s.enclave.EstimateGas(req.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.EstimateGasResponse{EncryptedResponse: encryptedBalance}, nil
}

func (s *server) decodeBlock(encodedBlock []byte) types.Block {
	block := types.Block{}
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		common.LogWithID(s.nodeShortID, "failed to decode block sent to enclave: %v", err)
	}
	return block
}
