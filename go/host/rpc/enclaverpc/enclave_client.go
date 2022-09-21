package enclaverpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc/connectivity"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client implements enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type Client struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      config.HostConfig
	nodeShortID uint64
}

// TODO - Avoid panicking and return errors instead where appropriate.

func NewClient(config config.HostConfig) *Client {
	nodeShortID := common.ShortAddress(config.ID)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		common.PanicWithID(nodeShortID, "Failed to connect to enclave RPC service. Cause: %s", err)
	}

	// We wait for the RPC connection to be ready.
	currentTime := time.Now()
	deadline := currentTime.Add(30 * time.Second)
	currentState := connection.GetState()
	for currentState == connectivity.Idle || currentState == connectivity.Connecting || currentState == connectivity.TransientFailure {
		connection.Connect()
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
		currentState = connection.GetState()
	}

	if currentState != connectivity.Ready {
		common.PanicWithID(nodeShortID, "RPC connection failed to establish. Current state is %s", currentState)
	}

	return &Client{
		generated.NewEnclaveProtoClient(connection),
		connection,
		config,
		nodeShortID,
	}
}

func (c *Client) StopClient() error {
	return c.connection.Close()
}

func (c *Client) Status() (common.Status, error) {
	if c.connection.GetState() != connectivity.Ready {
		return common.Unavailable, errors.New("RPC connection is not ready")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.Status(timeoutCtx, &generated.StatusRequest{})
	if err != nil {
		return common.Unavailable, err
	}
	if resp.GetError() != "" {
		return common.Unavailable, errors.New(resp.GetError())
	}
	return common.Status(resp.GetStatus()), nil
}

func (c *Client) Attestation() (*common.AttestationReport, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attestation. Cause: %s", err)
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg), nil
}

func (c *Client) GenerateSecret() (common.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret. Cause: %s", err)
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *Client) ShareSecret(report *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	attestationReportMsg := rpc.ToAttestationReportMsg(report)
	request := generated.FetchSecretRequest{AttestationReportMsg: &attestationReportMsg}
	response, err := c.protoClient.ShareSecret(timeoutCtx, &request)
	if err != nil {
		return nil, err
	}
	if response.GetError() != "" {
		return nil, errors.New(response.GetError())
	}
	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *Client) InitEnclave(secret common.EncryptedSharedEnclaveSecret) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}

func (c *Client) ProduceGenesis(blkHash gethcommon.Hash) (common.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ProduceGenesis(timeoutCtx, &generated.ProduceGenesisRequest{BlockHash: blkHash.Bytes()})
	if err != nil {
		return common.BlockSubmissionResponse{}, fmt.Errorf("could not produce genesis block. Cause: %w", err)
	}

	blockSubmissionResponse, err := rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
	if err != nil {
		return common.BlockSubmissionResponse{}, fmt.Errorf("could not produce block submission response. Cause: %w", err)
	}
	return blockSubmissionResponse, nil
}

func (c *Client) IngestBlocks(blocks []*types.Block) []common.BlockSubmissionResponse {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	encodedBlocks := make([][]byte, 0)
	for _, block := range blocks {
		encodedBlock, err := common.EncodeBlock(block)
		if err != nil {
			common.PanicWithID(c.nodeShortID, "Failed to ingest blocks. Cause: %s", err)
		}
		encodedBlocks = append(encodedBlocks, encodedBlock)
	}
	response, err := c.protoClient.IngestBlocks(timeoutCtx, &generated.IngestBlocksRequest{EncodedBlocks: encodedBlocks})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to ingest blocks. Cause: %s", err)
	}
	responses := response.GetBlockSubmissionResponses()
	result := make([]common.BlockSubmissionResponse, len(responses))
	for i, r := range responses {
		blockSubmissionResponse, err := rpc.FromBlockSubmissionResponseMsg(r)
		if err != nil {
			common.PanicWithID(c.nodeShortID, "Failed to produce block submission response. Cause: %s", err)
		}
		result[i] = blockSubmissionResponse
	}
	return result
}

func (c *Client) Start(block types.Block) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return fmt.Errorf("could not encode block. Cause: %w", err)
	}

	_, err := c.protoClient.Start(timeoutCtx, &generated.StartRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		return fmt.Errorf("could not start enclave. Cause: %w", err)
	}

	return nil
}

func (c *Client) SubmitBlock(block types.Block) (common.BlockSubmissionResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return common.BlockSubmissionResponse{}, fmt.Errorf("could not encode block. Cause: %w", err)
	}

	response, err := c.protoClient.SubmitBlock(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes()})
	if err != nil {
		return common.BlockSubmissionResponse{}, fmt.Errorf("could not submit block. Cause: %w", err)
	}

	blockSubmissionResponse, err := rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
	if err != nil {
		return common.BlockSubmissionResponse{}, fmt.Errorf("could not produce block submission response. Cause: %w", err)
	}
	return blockSubmissionResponse, nil
}

func (c *Client) SubmitRollup(rollup common.ExtRollup) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	extRollupMsg := rpc.ToExtRollupMsg(&rollup)
	_, err := c.protoClient.SubmitRollup(timeoutCtx, &generated.SubmitRollupRequest{ExtRollup: &extRollupMsg})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Could not submit rollup. Cause: %s", err)
	}
}

func (c *Client) SubmitTx(tx common.EncryptedTx) (common.EncryptedResponseSendRawTx, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return nil, err
	}
	return response.EncryptedHash, err
}

func (c *Client) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ExecuteOffChainTransaction(timeoutCtx, &generated.OffChainRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result, nil
}

func (c *Client) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (common.EncryptedResponseGetTxCount, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result, nil
}

func (c *Client) RoundWinner(parent common.L2RootHash) (common.ExtRollup, bool, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.RoundWinner(timeoutCtx, &generated.RoundWinnerRequest{Parent: parent.Bytes()})
	if err != nil {
		common.PanicWithID(c.nodeShortID, "Failed to determine round winner. Cause: %s", err)
	}

	if response.Winner {
		return rpc.FromExtRollupMsg(response.ExtRollup), true, nil
	}
	return common.ExtRollup{}, false, nil
}

func (c *Client) Stop() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return fmt.Errorf("could not stop enclave: %w", err)
	}
	return nil
}

func (c *Client) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (common.EncryptedResponseGetTxByHash, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedTx, nil
}

func (c *Client) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (common.EncryptedResponseGetTxReceipt, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, err
	}
	return response.EncryptedTxReceipt, nil
}

func (c *Client) GetRollup(rollupHash common.L2RootHash) (*common.ExtRollup, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetRollup(timeoutCtx, &generated.GetRollupRequest{RollupHash: rollupHash.Bytes()})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rollup with hash %s. Cause: %w", rollupHash.Hex(), err)
	}

	extRollup := rpc.FromExtRollupMsg(response.ExtRollup)
	return &extRollup, nil
}

func (c *Client) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	_, err := c.protoClient.AddViewingKey(timeoutCtx, &generated.AddViewingKeyRequest{
		ViewingKey: viewingKeyBytes,
		Signature:  signature,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedBalance, nil
}

func (c *Client) GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.GetCode(timeoutCtx, &generated.GetCodeRequest{
		Address:    address.Bytes(),
		RollupHash: rollupHash.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	return resp.Code, nil
}

func (c *Client) StoreAttestation(report *common.AttestationReport) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	msg := rpc.ToAttestationReportMsg(report)
	resp, err := c.protoClient.StoreAttestation(timeoutCtx, &generated.StoreAttestationRequest{
		AttestationReportMsg: &msg,
	})
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}
	return nil
}

func (c *Client) Subscribe(id uuid.UUID, encryptedParams common.EncryptedParamsLogSubscription) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	idBinary, err := id.MarshalBinary()
	if err != nil {
		return fmt.Errorf("could not marshall subscription ID to binary. Cause: %w", err)
	}

	_, err = c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    idBinary,
		EncryptedSubscription: encryptedParams,
	})
	return err
}

func (c *Client) Unsubscribe(id uuid.UUID) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	idBinary, err := id.MarshalBinary()
	if err != nil {
		return fmt.Errorf("could not marshall subscription ID to binary. Cause: %w", err)
	}

	_, err = c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: idBinary,
	})
	return err
}

func (c *Client) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (common.EncryptedResponseEstimateGas, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	resp, err := c.protoClient.EstimateGas(timeoutCtx, &generated.EstimateGasRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, err
	}
	return resp.EncryptedResponse, nil
}
