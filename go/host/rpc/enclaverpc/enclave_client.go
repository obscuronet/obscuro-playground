package enclaverpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/common/rpc"
	"github.com/ten-protocol/go-ten/go/common/rpc/generated"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/responses"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// Client implements enclave.Enclave and should be used by the host when communicating with the enclave via RPC.
type Client struct {
	protoClient generated.EnclaveProtoClient
	connection  *grpc.ClientConn
	config      *config.HostConfig
	logger      gethlog.Logger
}

func NewClient(config *config.HostConfig, logger gethlog.Logger) common.Enclave {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	connection, err := grpc.Dial(config.EnclaveRPCAddress, opts...)
	if err != nil {
		logger.Crit("Failed to connect to enclave RPC service.", log.ErrKey, err)
	}
	connection.Connect()
	// perform an initial sleep because that Connect() method is not blocking and the retry immediately checks the status
	time.Sleep(500 * time.Millisecond)

	// We wait for the RPC connection to be ready.
	err = retry.Do(func() error {
		currState := connection.GetState()
		if currState != connectivity.Ready {
			logger.Info("retrying connection until enclave is available", "status", currState.String(), "rpcAddr", config.EnclaveRPCAddress)
			connection.Connect()
			return fmt.Errorf("connection is not ready, status=%s", currState)
		}
		// connection is ready, break out of the loop
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{500 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))
	if err != nil {
		// this should not happen as we retry forever...
		logger.Crit("failed to connect to enclave", log.ErrKey, err)
	}

	return &Client{
		protoClient: generated.NewEnclaveProtoClient(connection),
		connection:  connection,
		config:      config,
		logger:      logger,
	}
}

func (c *Client) StopClient() common.SystemError {
	c.logger.Info("Closing rpc server connection.")
	return c.connection.Close()
}

func (c *Client) Status() (common.Status, common.SystemError) {
	if c.connection.GetState() != connectivity.Ready {
		return common.Status{StatusCode: common.Unavailable}, syserr.NewInternalError(fmt.Errorf("RPC connection is not ready"))
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Status(timeoutCtx, &generated.StatusRequest{})
	if err != nil {
		return common.Status{StatusCode: common.Unavailable}, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return common.Status{StatusCode: common.Unavailable}, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return common.Status{
		StatusCode: common.StatusCode(response.StatusCode),
		L1Head:     gethcommon.BytesToHash(response.L1Head),
		L2Head:     big.NewInt(0).SetBytes(response.L2Head),
	}, nil
}

func (c *Client) Attestation() (*common.AttestationReport, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Attestation(timeoutCtx, &generated.AttestationRequest{})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return rpc.FromAttestationReportMsg(response.AttestationReportMsg), nil
}

func (c *Client) GenerateSecret() (common.EncryptedSharedEnclaveSecret, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GenerateSecret(timeoutCtx, &generated.GenerateSecretRequest{})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return response.EncryptedSharedEnclaveSecret, nil
}

func (c *Client) InitEnclave(secret common.EncryptedSharedEnclaveSecret) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.InitEnclave(timeoutCtx, &generated.InitEnclaveRequest{EncryptedSharedEnclaveSecret: secret})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return nil
}

func (c *Client) EnclaveID() (common.EnclaveID, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.EnclaveID(timeoutCtx, &generated.EnclaveIDRequest{})
	if err != nil {
		return common.EnclaveID{}, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return common.EnclaveID{}, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return common.EnclaveID(response.EnclaveID), nil
}

func (c *Client) SubmitL1Block(block types.Block, receipts types.Receipts, isLatest bool) (*common.BlockSubmissionResponse, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	var buffer bytes.Buffer
	if err := block.EncodeRLP(&buffer); err != nil {
		return nil, fmt.Errorf("could not encode block. Cause: %w", err)
	}

	serialized, err := rlp.EncodeToBytes(receipts)
	if err != nil {
		return nil, fmt.Errorf("could not encode receipts. Cause: %w", err)
	}

	response, err := c.protoClient.SubmitL1Block(timeoutCtx, &generated.SubmitBlockRequest{EncodedBlock: buffer.Bytes(), EncodedReceipts: serialized, IsLatest: isLatest})
	if err != nil {
		return nil, fmt.Errorf("could not submit block. Cause: %w", err)
	}

	blockSubmissionResponse, err := rpc.FromBlockSubmissionResponseMsg(response.BlockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return blockSubmissionResponse, nil
}

func (c *Client) SubmitTx(tx common.EncryptedTx) (*responses.RawTx, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.SubmitTx(timeoutCtx, &generated.SubmitTxRequest{EncryptedTx: tx})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) SubmitBatch(batch *common.ExtBatch) common.SystemError {
	defer core.LogMethodDuration(c.logger, measure.NewStopwatch(), "SubmitBatch rpc call")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg := rpc.ToExtBatchMsg(batch)

	response, err := c.protoClient.SubmitBatch(timeoutCtx, &generated.SubmitBatchRequest{Batch: &batchMsg})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return nil
}

func (c *Client) ObsCall(encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.ObsCall(timeoutCtx, &generated.ObsCallRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionCount(timeoutCtx, &generated.GetTransactionCountRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) Stop() common.SystemError {
	c.logger.Info("Shutting down enclave client.")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Stop(timeoutCtx, &generated.StopRequest{})
	if err != nil {
		return syserr.NewRPCError(fmt.Errorf("could not stop enclave: %w", err))
	}
	if response != nil && response.SystemError != nil {
		return syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return nil
}

func (c *Client) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (*responses.TxByHash, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransaction(timeoutCtx, &generated.GetTransactionRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (*responses.TxReceipt, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTransactionReceipt(timeoutCtx, &generated.GetTransactionReceiptRequest{EncryptedParams: encryptedParams})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (*responses.Balance, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetBalance(timeoutCtx, &generated.GetBalanceRequest{
		EncryptedParams: encryptedParams,
	})
	response.String()
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetCode(address gethcommon.Address, batchHash *gethcommon.Hash) ([]byte, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetCode(timeoutCtx, &generated.GetCodeRequest{
		Address:    address.Bytes(),
		RollupHash: batchHash.Bytes(),
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return response.Code, nil
}

func (c *Client) Subscribe(id gethrpc.ID, encryptedParams common.EncryptedParamsLogSubscription) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Subscribe(timeoutCtx, &generated.SubscribeRequest{
		Id:                    []byte(id),
		EncryptedSubscription: encryptedParams,
	})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return nil
}

func (c *Client) Unsubscribe(id gethrpc.ID) common.SystemError {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.Unsubscribe(timeoutCtx, &generated.UnsubscribeRequest{
		Id: []byte(id),
	})
	if err != nil {
		return syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return nil
}

func (c *Client) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (*responses.Gas, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.EstimateGas(timeoutCtx, &generated.EstimateGasRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetLogs(timeoutCtx, &generated.GetLogsRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) HealthCheck() (bool, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.HealthCheck(timeoutCtx, &generated.EmptyArgs{})
	if err != nil {
		return false, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return false, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return response.Status, nil
}

func (c *Client) CreateBatch(skipIfEmpty bool) common.SystemError {
	defer core.LogMethodDuration(c.logger, measure.NewStopwatch(), "CreateBatch rpc call")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.CreateBatch(timeoutCtx, &generated.CreateBatchRequest{SkipIfEmpty: skipIfEmpty})
	if err != nil {
		return syserr.NewInternalError(err)
	}
	if response != nil && response.Error != "" {
		return syserr.NewInternalError(fmt.Errorf("%s", response.Error))
	}
	return err
}

func (c *Client) CreateRollup(fromSeqNo uint64) (*common.ExtRollup, common.SystemError) {
	defer core.LogMethodDuration(c.logger, measure.NewStopwatch(), "CreateRollup rpc call")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.CreateRollup(timeoutCtx, &generated.CreateRollupRequest{
		FromSequenceNumber: &fromSeqNo,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return rpc.FromExtRollupMsg(response.Msg), nil
}

func (c *Client) DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	confBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := c.protoClient.DebugTraceTransaction(timeoutCtx, &generated.DebugTraceTransactionRequest{
		TxHash: hash.Bytes(),
		Config: confBytes,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return json.RawMessage(response.Msg), nil
}

func (c *Client) GetBatch(hash common.L2BatchHash) (*common.ExtBatch, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg, err := c.protoClient.GetBatch(timeoutCtx, &generated.GetBatchRequest{KnownHead: hash.Bytes()})
	if err != nil {
		return nil, fmt.Errorf("rpc GetBatch failed. Cause: %w", err)
	}

	return common.DecodeExtBatch(batchMsg.Batch)
}

func (c *Client) GetBatchBySeqNo(seqNo uint64) (*common.ExtBatch, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	batchMsg, err := c.protoClient.GetBatchBySeqNo(timeoutCtx, &generated.GetBatchBySeqNoRequest{SeqNo: seqNo})
	if err != nil {
		return nil, fmt.Errorf("rpc GetBatchBySeqNo failed. Cause: %w", err)
	}

	return common.DecodeExtBatch(batchMsg.Batch)
}

func (c *Client) GetRollupData(hash common.L2RollupHash) (*common.PublicRollupMetadata, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetRollupData(timeoutCtx, &generated.GetRollupDataRequest{Hash: hash.Bytes()})
	if err != nil {
		return nil, fmt.Errorf("rpc GetRollupData failed. Cause: %w", err)
	}
	return rpc.FromRollupDataMsg(response.Msg)
}

func (c *Client) StreamL2Updates() (chan common.StreamL2UpdatesResponse, func()) {
	// channel size is 10 to allow for some buffering but caller is expected to read immediately to avoid blocking
	batchChan := make(chan common.StreamL2UpdatesResponse, 10)
	cancelCtx, cancel := context.WithCancel(context.Background())

	stream, err := c.protoClient.StreamL2Updates(cancelCtx, &generated.StreamL2UpdatesRequest{})
	if err != nil {
		c.logger.Error("Error opening batch stream.", log.ErrKey, err)
		cancel()
		close(batchChan)
		// return closed channel and no-op cancel func
		return batchChan, func() {}
	}

	stopIt := func() {
		c.logger.Info("Closing batch stream.")
		if err := stream.CloseSend(); err != nil {
			c.logger.Error("Client is unable to close batch stream", log.ErrKey, err)
		}

		cancel()
		close(batchChan)
	}

	go func() {
		defer stopIt()
		for {
			batchMsg, err := stream.Recv()
			if err != nil {
				c.logger.Error("Error receiving batch from stream.", log.ErrKey, err)
				break
			}

			var decoded common.StreamL2UpdatesResponse
			if err := json.Unmarshal(batchMsg.Batch, &decoded); err != nil {
				c.logger.Error("Error unmarshalling batch from stream.", log.ErrKey, err)
				break
			}

			batchChan <- decoded
		}
	}()

	return batchChan, cancel
}

func (c *Client) DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.DebugEventLogRelevancy(timeoutCtx, &generated.DebugEventLogRelevancyRequest{
		TxHash: hash.Bytes(),
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return json.RawMessage(response.Msg), nil
}

func (c *Client) GetTotalContractCount() (*big.Int, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetTotalContractCount(timeoutCtx, &generated.GetTotalContractCountRequest{})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return big.NewInt(response.Count), nil
}

func (c *Client) GetCustomQuery(encryptedParams common.EncryptedParamsGetStorageAt) (*responses.Receipts, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetReceiptsByAddress(timeoutCtx, &generated.GetReceiptsByAddressRequest{
		EncryptedParams: encryptedParams,
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return responses.ToEnclaveResponse(response.EncodedEnclaveResponse), nil
}

func (c *Client) GetPublicTransactionData(pagination *common.QueryPagination) (*common.TransactionListingResponse, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.GetPublicTransactionData(timeoutCtx, &generated.GetPublicTransactionDataRequest{
		Pagination: &generated.Pagination{
			Offset: int32(pagination.Offset),
			Size:   int32(pagination.Size),
		},
	})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	var result common.TransactionListingResponse
	err = json.Unmarshal(response.PublicTransactionData, &result)
	if err != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}

	return &result, nil
}

func (c *Client) EnclavePublicConfig() (*common.EnclavePublicConfig, common.SystemError) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.config.EnclaveRPCTimeout)
	defer cancel()

	response, err := c.protoClient.EnclavePublicConfig(timeoutCtx, &generated.EnclavePublicConfigRequest{})
	if err != nil {
		return nil, syserr.NewRPCError(err)
	}
	if response != nil && response.SystemError != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("%s", response.SystemError.ErrorString))
	}
	return &common.EnclavePublicConfig{
		L2MessageBusAddress: gethcommon.BytesToAddress(response.L2MessageBusAddress),
	}, nil
}
