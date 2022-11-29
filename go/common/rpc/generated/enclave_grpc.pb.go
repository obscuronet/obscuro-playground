// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: enclave.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EnclaveProtoClient is the client API for EnclaveProto service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnclaveProtoClient interface {
	// Status is used to check whether the server is ready for requests.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation(ctx context.Context, in *AttestationRequest, opts ...grpc.CallOption) (*AttestationResponse, error)
	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret(ctx context.Context, in *GenerateSecretRequest, opts ...grpc.CallOption) (*GenerateSecretResponse, error)
	// Init - initialise an enclave with a seed received by another enclave
	InitEnclave(ctx context.Context, in *InitEnclaveRequest, opts ...grpc.CallOption) (*InitEnclaveResponse, error)
	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis(ctx context.Context, in *ProduceGenesisRequest, opts ...grpc.CallOption) (*ProduceGenesisResponse, error)
	// Start - start speculative execution
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	// SubmitL1Block - Used for the host to submit blocks to the enclave, these may be:
	//
	//	a. historic block - if the enclave is behind and in the process of catching up with the L1 state
	//	b. the latest block published by the L1, to which the enclave should respond with a rollup
	//
	// It is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving ancestors of it, will result in it being ignored
	SubmitL1Block(ctx context.Context, in *SubmitBlockRequest, opts ...grpc.CallOption) (*SubmitBlockResponse, error)
	// ProduceRollup creates a new rollup.
	ProduceRollup(ctx context.Context, in *ProduceRollupRequest, opts ...grpc.CallOption) (*ProduceRollupResponse, error)
	// SubmitTx - user transactions
	SubmitTx(ctx context.Context, in *SubmitTxRequest, opts ...grpc.CallOption) (*SubmitTxResponse, error)
	// ExecuteOffChainTransaction - returns the result of executing the smart contract as a user, encrypted with the
	// viewing key corresponding to the `from` field
	ExecuteOffChainTransaction(ctx context.Context, in *OffChainRequest, opts ...grpc.CallOption) (*OffChainResponse, error)
	// GetTransactionCount - returns the nonce of the wallet with the given address.
	GetTransactionCount(ctx context.Context, in *GetTransactionCountRequest, opts ...grpc.CallOption) (*GetTransactionCountResponse, error)
	// Stop gracefully stops the enclave
	Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error)
	// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	// GetTransaction returns a transaction receipt given the transaction's signed hash, encrypted with the viewing key
	// corresponding to the original transaction submitter
	GetTransactionReceipt(ctx context.Context, in *GetTransactionReceiptRequest, opts ...grpc.CallOption) (*GetTransactionReceiptResponse, error)
	// AddViewingKey adds a viewing key to the enclave
	AddViewingKey(ctx context.Context, in *AddViewingKeyRequest, opts ...grpc.CallOption) (*AddViewingKeyResponse, error)
	// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to
	// the address
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	// GetCode returns the code stored at the given address in the state for the given rollup height or rollup hash
	GetCode(ctx context.Context, in *GetCodeRequest, opts ...grpc.CallOption) (*GetCodeResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeResponse, error)
	Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error)
	// EstimateGas returns the estimation of gas used for the given transactions
	EstimateGas(ctx context.Context, in *EstimateGasRequest, opts ...grpc.CallOption) (*EstimateGasResponse, error)
	GetLogs(ctx context.Context, in *GetLogsRequest, opts ...grpc.CallOption) (*GetLogsResponse, error)
	// HealthCheck returns the health status of enclave + db
	HealthCheck(ctx context.Context, in *EmptyArgs, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type enclaveProtoClient struct {
	cc grpc.ClientConnInterface
}

func NewEnclaveProtoClient(cc grpc.ClientConnInterface) EnclaveProtoClient {
	return &enclaveProtoClient{cc}
}

func (c *enclaveProtoClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) Attestation(ctx context.Context, in *AttestationRequest, opts ...grpc.CallOption) (*AttestationResponse, error) {
	out := new(AttestationResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Attestation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GenerateSecret(ctx context.Context, in *GenerateSecretRequest, opts ...grpc.CallOption) (*GenerateSecretResponse, error) {
	out := new(GenerateSecretResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GenerateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) InitEnclave(ctx context.Context, in *InitEnclaveRequest, opts ...grpc.CallOption) (*InitEnclaveResponse, error) {
	out := new(InitEnclaveResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/InitEnclave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) ProduceGenesis(ctx context.Context, in *ProduceGenesisRequest, opts ...grpc.CallOption) (*ProduceGenesisResponse, error) {
	out := new(ProduceGenesisResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/ProduceGenesis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) SubmitL1Block(ctx context.Context, in *SubmitBlockRequest, opts ...grpc.CallOption) (*SubmitBlockResponse, error) {
	out := new(SubmitBlockResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/SubmitL1Block", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) ProduceRollup(ctx context.Context, in *ProduceRollupRequest, opts ...grpc.CallOption) (*ProduceRollupResponse, error) {
	out := new(ProduceRollupResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/ProduceRollup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) SubmitTx(ctx context.Context, in *SubmitTxRequest, opts ...grpc.CallOption) (*SubmitTxResponse, error) {
	out := new(SubmitTxResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/SubmitTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) ExecuteOffChainTransaction(ctx context.Context, in *OffChainRequest, opts ...grpc.CallOption) (*OffChainResponse, error) {
	out := new(OffChainResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/ExecuteOffChainTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetTransactionCount(ctx context.Context, in *GetTransactionCountRequest, opts ...grpc.CallOption) (*GetTransactionCountResponse, error) {
	out := new(GetTransactionCountResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetTransactionCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error) {
	out := new(StopResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
	out := new(GetTransactionResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetTransactionReceipt(ctx context.Context, in *GetTransactionReceiptRequest, opts ...grpc.CallOption) (*GetTransactionReceiptResponse, error) {
	out := new(GetTransactionReceiptResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetTransactionReceipt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) AddViewingKey(ctx context.Context, in *AddViewingKeyRequest, opts ...grpc.CallOption) (*AddViewingKeyResponse, error) {
	out := new(AddViewingKeyResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/AddViewingKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetCode(ctx context.Context, in *GetCodeRequest, opts ...grpc.CallOption) (*GetCodeResponse, error) {
	out := new(GetCodeResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Subscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error) {
	out := new(UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) EstimateGas(ctx context.Context, in *EstimateGasRequest, opts ...grpc.CallOption) (*EstimateGasResponse, error) {
	out := new(EstimateGasResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/EstimateGas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) GetLogs(ctx context.Context, in *GetLogsRequest, opts ...grpc.CallOption) (*GetLogsResponse, error) {
	out := new(GetLogsResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/GetLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) HealthCheck(ctx context.Context, in *EmptyArgs, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnclaveProtoServer is the server API for EnclaveProto service.
// All implementations must embed UnimplementedEnclaveProtoServer
// for forward compatibility
type EnclaveProtoServer interface {
	// Status is used to check whether the server is ready for requests.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation(context.Context, *AttestationRequest) (*AttestationResponse, error)
	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret(context.Context, *GenerateSecretRequest) (*GenerateSecretResponse, error)
	// Init - initialise an enclave with a seed received by another enclave
	InitEnclave(context.Context, *InitEnclaveRequest) (*InitEnclaveResponse, error)
	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis(context.Context, *ProduceGenesisRequest) (*ProduceGenesisResponse, error)
	// Start - start speculative execution
	Start(context.Context, *StartRequest) (*StartResponse, error)
	// SubmitL1Block - Used for the host to submit blocks to the enclave, these may be:
	//
	//	a. historic block - if the enclave is behind and in the process of catching up with the L1 state
	//	b. the latest block published by the L1, to which the enclave should respond with a rollup
	//
	// It is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving ancestors of it, will result in it being ignored
	SubmitL1Block(context.Context, *SubmitBlockRequest) (*SubmitBlockResponse, error)
	// ProduceRollup creates a new rollup.
	ProduceRollup(context.Context, *ProduceRollupRequest) (*ProduceRollupResponse, error)
	// SubmitTx - user transactions
	SubmitTx(context.Context, *SubmitTxRequest) (*SubmitTxResponse, error)
	// ExecuteOffChainTransaction - returns the result of executing the smart contract as a user, encrypted with the
	// viewing key corresponding to the `from` field
	ExecuteOffChainTransaction(context.Context, *OffChainRequest) (*OffChainResponse, error)
	// GetTransactionCount - returns the nonce of the wallet with the given address.
	GetTransactionCount(context.Context, *GetTransactionCountRequest) (*GetTransactionCountResponse, error)
	// Stop gracefully stops the enclave
	Stop(context.Context, *StopRequest) (*StopResponse, error)
	// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
	GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error)
	// GetTransaction returns a transaction receipt given the transaction's signed hash, encrypted with the viewing key
	// corresponding to the original transaction submitter
	GetTransactionReceipt(context.Context, *GetTransactionReceiptRequest) (*GetTransactionReceiptResponse, error)
	// AddViewingKey adds a viewing key to the enclave
	AddViewingKey(context.Context, *AddViewingKeyRequest) (*AddViewingKeyResponse, error)
	// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to
	// the address
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	// GetCode returns the code stored at the given address in the state for the given rollup height or rollup hash
	GetCode(context.Context, *GetCodeRequest) (*GetCodeResponse, error)
	Subscribe(context.Context, *SubscribeRequest) (*SubscribeResponse, error)
	Unsubscribe(context.Context, *UnsubscribeRequest) (*UnsubscribeResponse, error)
	// EstimateGas returns the estimation of gas used for the given transactions
	EstimateGas(context.Context, *EstimateGasRequest) (*EstimateGasResponse, error)
	GetLogs(context.Context, *GetLogsRequest) (*GetLogsResponse, error)
	// HealthCheck returns the health status of enclave + db
	HealthCheck(context.Context, *EmptyArgs) (*HealthCheckResponse, error)
	mustEmbedUnimplementedEnclaveProtoServer()
}

// UnimplementedEnclaveProtoServer must be embedded to have forward compatible implementations.
type UnimplementedEnclaveProtoServer struct {
}

func (UnimplementedEnclaveProtoServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedEnclaveProtoServer) Attestation(context.Context, *AttestationRequest) (*AttestationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Attestation not implemented")
}
func (UnimplementedEnclaveProtoServer) GenerateSecret(context.Context, *GenerateSecretRequest) (*GenerateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateSecret not implemented")
}
func (UnimplementedEnclaveProtoServer) InitEnclave(context.Context, *InitEnclaveRequest) (*InitEnclaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitEnclave not implemented")
}
func (UnimplementedEnclaveProtoServer) ProduceGenesis(context.Context, *ProduceGenesisRequest) (*ProduceGenesisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceGenesis not implemented")
}
func (UnimplementedEnclaveProtoServer) Start(context.Context, *StartRequest) (*StartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedEnclaveProtoServer) SubmitL1Block(context.Context, *SubmitBlockRequest) (*SubmitBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitL1Block not implemented")
}
func (UnimplementedEnclaveProtoServer) ProduceRollup(context.Context, *ProduceRollupRequest) (*ProduceRollupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceRollup not implemented")
}
func (UnimplementedEnclaveProtoServer) SubmitTx(context.Context, *SubmitTxRequest) (*SubmitTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitTx not implemented")
}
func (UnimplementedEnclaveProtoServer) ExecuteOffChainTransaction(context.Context, *OffChainRequest) (*OffChainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteOffChainTransaction not implemented")
}
func (UnimplementedEnclaveProtoServer) GetTransactionCount(context.Context, *GetTransactionCountRequest) (*GetTransactionCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionCount not implemented")
}
func (UnimplementedEnclaveProtoServer) Stop(context.Context, *StopRequest) (*StopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedEnclaveProtoServer) GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
}
func (UnimplementedEnclaveProtoServer) GetTransactionReceipt(context.Context, *GetTransactionReceiptRequest) (*GetTransactionReceiptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionReceipt not implemented")
}
func (UnimplementedEnclaveProtoServer) AddViewingKey(context.Context, *AddViewingKeyRequest) (*AddViewingKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddViewingKey not implemented")
}
func (UnimplementedEnclaveProtoServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedEnclaveProtoServer) GetCode(context.Context, *GetCodeRequest) (*GetCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCode not implemented")
}
func (UnimplementedEnclaveProtoServer) Subscribe(context.Context, *SubscribeRequest) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedEnclaveProtoServer) Unsubscribe(context.Context, *UnsubscribeRequest) (*UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedEnclaveProtoServer) EstimateGas(context.Context, *EstimateGasRequest) (*EstimateGasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimateGas not implemented")
}
func (UnimplementedEnclaveProtoServer) GetLogs(context.Context, *GetLogsRequest) (*GetLogsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogs not implemented")
}
func (UnimplementedEnclaveProtoServer) HealthCheck(context.Context, *EmptyArgs) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedEnclaveProtoServer) mustEmbedUnimplementedEnclaveProtoServer() {}

// UnsafeEnclaveProtoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnclaveProtoServer will
// result in compilation errors.
type UnsafeEnclaveProtoServer interface {
	mustEmbedUnimplementedEnclaveProtoServer()
}

func RegisterEnclaveProtoServer(s grpc.ServiceRegistrar, srv EnclaveProtoServer) {
	s.RegisterService(&EnclaveProto_ServiceDesc, srv)
}

func _EnclaveProto_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_Attestation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttestationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Attestation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Attestation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Attestation(ctx, req.(*AttestationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GenerateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GenerateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GenerateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GenerateSecret(ctx, req.(*GenerateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_InitEnclave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitEnclaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).InitEnclave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/InitEnclave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).InitEnclave(ctx, req.(*InitEnclaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_ProduceGenesis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProduceGenesisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).ProduceGenesis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/ProduceGenesis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).ProduceGenesis(ctx, req.(*ProduceGenesisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_SubmitL1Block_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).SubmitL1Block(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/SubmitL1Block",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).SubmitL1Block(ctx, req.(*SubmitBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_ProduceRollup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProduceRollupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).ProduceRollup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/ProduceRollup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).ProduceRollup(ctx, req.(*ProduceRollupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_SubmitTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitTxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).SubmitTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/SubmitTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).SubmitTx(ctx, req.(*SubmitTxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_ExecuteOffChainTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OffChainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).ExecuteOffChainTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/ExecuteOffChainTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).ExecuteOffChainTransaction(ctx, req.(*OffChainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetTransactionCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetTransactionCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetTransactionCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetTransactionCount(ctx, req.(*GetTransactionCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Stop(ctx, req.(*StopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetTransaction(ctx, req.(*GetTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetTransactionReceipt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionReceiptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetTransactionReceipt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetTransactionReceipt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetTransactionReceipt(ctx, req.(*GetTransactionReceiptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_AddViewingKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddViewingKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).AddViewingKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/AddViewingKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).AddViewingKey(ctx, req.(*AddViewingKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetCode(ctx, req.(*GetCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Subscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Subscribe(ctx, req.(*SubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Unsubscribe(ctx, req.(*UnsubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_EstimateGas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EstimateGasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).EstimateGas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/EstimateGas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).EstimateGas(ctx, req.(*EstimateGasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_GetLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).GetLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/GetLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).GetLogs(ctx, req.(*GetLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).HealthCheck(ctx, req.(*EmptyArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// EnclaveProto_ServiceDesc is the grpc.ServiceDesc for EnclaveProto service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnclaveProto_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "generated.EnclaveProto",
	HandlerType: (*EnclaveProtoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _EnclaveProto_Status_Handler,
		},
		{
			MethodName: "Attestation",
			Handler:    _EnclaveProto_Attestation_Handler,
		},
		{
			MethodName: "GenerateSecret",
			Handler:    _EnclaveProto_GenerateSecret_Handler,
		},
		{
			MethodName: "InitEnclave",
			Handler:    _EnclaveProto_InitEnclave_Handler,
		},
		{
			MethodName: "ProduceGenesis",
			Handler:    _EnclaveProto_ProduceGenesis_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _EnclaveProto_Start_Handler,
		},
		{
			MethodName: "SubmitL1Block",
			Handler:    _EnclaveProto_SubmitL1Block_Handler,
		},
		{
			MethodName: "ProduceRollup",
			Handler:    _EnclaveProto_ProduceRollup_Handler,
		},
		{
			MethodName: "SubmitTx",
			Handler:    _EnclaveProto_SubmitTx_Handler,
		},
		{
			MethodName: "ExecuteOffChainTransaction",
			Handler:    _EnclaveProto_ExecuteOffChainTransaction_Handler,
		},
		{
			MethodName: "GetTransactionCount",
			Handler:    _EnclaveProto_GetTransactionCount_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _EnclaveProto_Stop_Handler,
		},
		{
			MethodName: "GetTransaction",
			Handler:    _EnclaveProto_GetTransaction_Handler,
		},
		{
			MethodName: "GetTransactionReceipt",
			Handler:    _EnclaveProto_GetTransactionReceipt_Handler,
		},
		{
			MethodName: "AddViewingKey",
			Handler:    _EnclaveProto_AddViewingKey_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _EnclaveProto_GetBalance_Handler,
		},
		{
			MethodName: "GetCode",
			Handler:    _EnclaveProto_GetCode_Handler,
		},
		{
			MethodName: "Subscribe",
			Handler:    _EnclaveProto_Subscribe_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _EnclaveProto_Unsubscribe_Handler,
		},
		{
			MethodName: "EstimateGas",
			Handler:    _EnclaveProto_EstimateGas_Handler,
		},
		{
			MethodName: "GetLogs",
			Handler:    _EnclaveProto_GetLogs_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _EnclaveProto_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "enclave.proto",
}
