// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
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
	// IsReady is used to check whether the server is ready for requests.
	IsReady(ctx context.Context, in *IsReadyRequest, opts ...grpc.CallOption) (*IsReadyResponse, error)
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation(ctx context.Context, in *AttestationRequest, opts ...grpc.CallOption) (*AttestationResponse, error)
	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret(ctx context.Context, in *GenerateSecretRequest, opts ...grpc.CallOption) (*GenerateSecretResponse, error)
	// ShareSecret - return the shared secret encrypted with the key from the attestation
	ShareSecret(ctx context.Context, in *FetchSecretRequest, opts ...grpc.CallOption) (*ShareSecretResponse, error)
	// Init - initialise an enclave with a seed received by another enclave
	InitEnclave(ctx context.Context, in *InitEnclaveRequest, opts ...grpc.CallOption) (*InitEnclaveResponse, error)
	// IsInitialised - true if the shared secret is available
	IsInitialised(ctx context.Context, in *IsInitialisedRequest, opts ...grpc.CallOption) (*IsInitialisedResponse, error)
	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis(ctx context.Context, in *ProduceGenesisRequest, opts ...grpc.CallOption) (*ProduceGenesisResponse, error)
	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(ctx context.Context, in *IngestBlocksRequest, opts ...grpc.CallOption) (*IngestBlocksResponse, error)
	// Start - start speculative execution
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(ctx context.Context, in *SubmitBlockRequest, opts ...grpc.CallOption) (*SubmitBlockResponse, error)
	// SubmitRollup - receive gossiped rollups
	SubmitRollup(ctx context.Context, in *SubmitRollupRequest, opts ...grpc.CallOption) (*SubmitRollupResponse, error)
	// SubmitTx - user transactions
	SubmitTx(ctx context.Context, in *SubmitTxRequest, opts ...grpc.CallOption) (*SubmitTxResponse, error)
	// Balance - returns the balance of an address with a block delay
	Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error)
	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(ctx context.Context, in *RoundWinnerRequest, opts ...grpc.CallOption) (*RoundWinnerResponse, error)
	// Stop gracefully stops the enclave
	Stop(ctx context.Context, in *StopRequest, opts ...grpc.CallOption) (*StopResponse, error)
	// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
}

type enclaveProtoClient struct {
	cc grpc.ClientConnInterface
}

func NewEnclaveProtoClient(cc grpc.ClientConnInterface) EnclaveProtoClient {
	return &enclaveProtoClient{cc}
}

func (c *enclaveProtoClient) IsReady(ctx context.Context, in *IsReadyRequest, opts ...grpc.CallOption) (*IsReadyResponse, error) {
	out := new(IsReadyResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/IsReady", in, out, opts...)
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

func (c *enclaveProtoClient) ShareSecret(ctx context.Context, in *FetchSecretRequest, opts ...grpc.CallOption) (*ShareSecretResponse, error) {
	out := new(ShareSecretResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/ShareSecret", in, out, opts...)
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

func (c *enclaveProtoClient) IsInitialised(ctx context.Context, in *IsInitialisedRequest, opts ...grpc.CallOption) (*IsInitialisedResponse, error) {
	out := new(IsInitialisedResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/IsInitialised", in, out, opts...)
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

func (c *enclaveProtoClient) IngestBlocks(ctx context.Context, in *IngestBlocksRequest, opts ...grpc.CallOption) (*IngestBlocksResponse, error) {
	out := new(IngestBlocksResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/IngestBlocks", in, out, opts...)
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

func (c *enclaveProtoClient) SubmitBlock(ctx context.Context, in *SubmitBlockRequest, opts ...grpc.CallOption) (*SubmitBlockResponse, error) {
	out := new(SubmitBlockResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/SubmitBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) SubmitRollup(ctx context.Context, in *SubmitRollupRequest, opts ...grpc.CallOption) (*SubmitRollupResponse, error) {
	out := new(SubmitRollupResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/SubmitRollup", in, out, opts...)
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

func (c *enclaveProtoClient) Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error) {
	out := new(BalanceResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/Balance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enclaveProtoClient) RoundWinner(ctx context.Context, in *RoundWinnerRequest, opts ...grpc.CallOption) (*RoundWinnerResponse, error) {
	out := new(RoundWinnerResponse)
	err := c.cc.Invoke(ctx, "/generated.EnclaveProto/RoundWinner", in, out, opts...)
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

// EnclaveProtoServer is the server API for EnclaveProto service.
// All implementations must embed UnimplementedEnclaveProtoServer
// for forward compatibility
type EnclaveProtoServer interface {
	// IsReady is used to check whether the server is ready for requests.
	IsReady(context.Context, *IsReadyRequest) (*IsReadyResponse, error)
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation(context.Context, *AttestationRequest) (*AttestationResponse, error)
	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret(context.Context, *GenerateSecretRequest) (*GenerateSecretResponse, error)
	// ShareSecret - return the shared secret encrypted with the key from the attestation
	ShareSecret(context.Context, *FetchSecretRequest) (*ShareSecretResponse, error)
	// Init - initialise an enclave with a seed received by another enclave
	InitEnclave(context.Context, *InitEnclaveRequest) (*InitEnclaveResponse, error)
	// IsInitialised - true if the shared secret is available
	IsInitialised(context.Context, *IsInitialisedRequest) (*IsInitialisedResponse, error)
	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis(context.Context, *ProduceGenesisRequest) (*ProduceGenesisResponse, error)
	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(context.Context, *IngestBlocksRequest) (*IngestBlocksResponse, error)
	// Start - start speculative execution
	Start(context.Context, *StartRequest) (*StartResponse, error)
	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(context.Context, *SubmitBlockRequest) (*SubmitBlockResponse, error)
	// SubmitRollup - receive gossiped rollups
	SubmitRollup(context.Context, *SubmitRollupRequest) (*SubmitRollupResponse, error)
	// SubmitTx - user transactions
	SubmitTx(context.Context, *SubmitTxRequest) (*SubmitTxResponse, error)
	// Balance - returns the balance of an address with a block delay
	Balance(context.Context, *BalanceRequest) (*BalanceResponse, error)
	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(context.Context, *RoundWinnerRequest) (*RoundWinnerResponse, error)
	// Stop gracefully stops the enclave
	Stop(context.Context, *StopRequest) (*StopResponse, error)
	// GetTransaction returns a transaction given its Signed Hash, returns nil, false when Transaction is unknown
	GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error)
	mustEmbedUnimplementedEnclaveProtoServer()
}

// UnimplementedEnclaveProtoServer must be embedded to have forward compatible implementations.
type UnimplementedEnclaveProtoServer struct {
}

func (UnimplementedEnclaveProtoServer) IsReady(context.Context, *IsReadyRequest) (*IsReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedEnclaveProtoServer) Attestation(context.Context, *AttestationRequest) (*AttestationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Attestation not implemented")
}
func (UnimplementedEnclaveProtoServer) GenerateSecret(context.Context, *GenerateSecretRequest) (*GenerateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateSecret not implemented")
}
func (UnimplementedEnclaveProtoServer) ShareSecret(context.Context, *FetchSecretRequest) (*ShareSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareSecret not implemented")
}
func (UnimplementedEnclaveProtoServer) InitEnclave(context.Context, *InitEnclaveRequest) (*InitEnclaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitEnclave not implemented")
}
func (UnimplementedEnclaveProtoServer) IsInitialised(context.Context, *IsInitialisedRequest) (*IsInitialisedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsInitialised not implemented")
}
func (UnimplementedEnclaveProtoServer) ProduceGenesis(context.Context, *ProduceGenesisRequest) (*ProduceGenesisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProduceGenesis not implemented")
}
func (UnimplementedEnclaveProtoServer) IngestBlocks(context.Context, *IngestBlocksRequest) (*IngestBlocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IngestBlocks not implemented")
}
func (UnimplementedEnclaveProtoServer) Start(context.Context, *StartRequest) (*StartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedEnclaveProtoServer) SubmitBlock(context.Context, *SubmitBlockRequest) (*SubmitBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBlock not implemented")
}
func (UnimplementedEnclaveProtoServer) SubmitRollup(context.Context, *SubmitRollupRequest) (*SubmitRollupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitRollup not implemented")
}
func (UnimplementedEnclaveProtoServer) SubmitTx(context.Context, *SubmitTxRequest) (*SubmitTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitTx not implemented")
}
func (UnimplementedEnclaveProtoServer) Balance(context.Context, *BalanceRequest) (*BalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
}
func (UnimplementedEnclaveProtoServer) RoundWinner(context.Context, *RoundWinnerRequest) (*RoundWinnerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RoundWinner not implemented")
}
func (UnimplementedEnclaveProtoServer) Stop(context.Context, *StopRequest) (*StopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedEnclaveProtoServer) GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
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

func _EnclaveProto_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).IsReady(ctx, req.(*IsReadyRequest))
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

func _EnclaveProto_ShareSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).ShareSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/ShareSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).ShareSecret(ctx, req.(*FetchSecretRequest))
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

func _EnclaveProto_IsInitialised_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsInitialisedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).IsInitialised(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/IsInitialised",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).IsInitialised(ctx, req.(*IsInitialisedRequest))
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

func _EnclaveProto_IngestBlocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestBlocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).IngestBlocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/IngestBlocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).IngestBlocks(ctx, req.(*IngestBlocksRequest))
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

func _EnclaveProto_SubmitBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).SubmitBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/SubmitBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).SubmitBlock(ctx, req.(*SubmitBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_SubmitRollup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitRollupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).SubmitRollup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/SubmitRollup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).SubmitRollup(ctx, req.(*SubmitRollupRequest))
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

func _EnclaveProto_Balance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).Balance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/Balance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).Balance(ctx, req.(*BalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnclaveProto_RoundWinner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoundWinnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnclaveProtoServer).RoundWinner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.EnclaveProto/RoundWinner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnclaveProtoServer).RoundWinner(ctx, req.(*RoundWinnerRequest))
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

// EnclaveProto_ServiceDesc is the grpc.ServiceDesc for EnclaveProto service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnclaveProto_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "generated.EnclaveProto",
	HandlerType: (*EnclaveProtoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _EnclaveProto_IsReady_Handler,
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
			MethodName: "ShareSecret",
			Handler:    _EnclaveProto_ShareSecret_Handler,
		},
		{
			MethodName: "InitEnclave",
			Handler:    _EnclaveProto_InitEnclave_Handler,
		},
		{
			MethodName: "IsInitialised",
			Handler:    _EnclaveProto_IsInitialised_Handler,
		},
		{
			MethodName: "ProduceGenesis",
			Handler:    _EnclaveProto_ProduceGenesis_Handler,
		},
		{
			MethodName: "IngestBlocks",
			Handler:    _EnclaveProto_IngestBlocks_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _EnclaveProto_Start_Handler,
		},
		{
			MethodName: "SubmitBlock",
			Handler:    _EnclaveProto_SubmitBlock_Handler,
		},
		{
			MethodName: "SubmitRollup",
			Handler:    _EnclaveProto_SubmitRollup_Handler,
		},
		{
			MethodName: "SubmitTx",
			Handler:    _EnclaveProto_SubmitTx_Handler,
		},
		{
			MethodName: "Balance",
			Handler:    _EnclaveProto_Balance_Handler,
		},
		{
			MethodName: "RoundWinner",
			Handler:    _EnclaveProto_RoundWinner_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _EnclaveProto_Stop_Handler,
		},
		{
			MethodName: "GetTransaction",
			Handler:    _EnclaveProto_GetTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "enclave.proto",
}
