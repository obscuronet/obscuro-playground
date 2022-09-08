package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

const (
	ws = "ws://"
)

// networkClient is a Client implementation that wraps Geth's rpc.Client to make calls to the obscuro node
type networkClient struct {
	rpcClient *rpc.Client
}

// NewEncNetworkClient returns a network RPC client with Viewing Key encryption/decryption
func NewEncNetworkClient(rpcAddress string, viewingKey *ViewingKey) (*EncRPCClient, error) {
	rpcClient, err := NewNetworkClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	encClient, err := NewEncRPCClient(rpcClient, viewingKey)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

// NewNetworkClient returns a client that can make RPC calls to an Obscuro node
func NewNetworkClient(address string) (Client, error) {
	rpcClient, err := rpc.Dial(ws + address)
	if err != nil {
		return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", ws+address, err)
	}

	return &networkClient{
		rpcClient: rpcClient,
	}, nil
}

// Call handles JSON rpc requests, delegating to the geth RPC client
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *networkClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(&result, method, args...)
}

func (c *networkClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.CallContext(ctx, result, method, args...)
}

func (c *networkClient) Stop() {
	c.rpcClient.Close()
}
