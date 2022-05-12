package host

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"
)

// An in-memory implementation of `clientserver.Client` that speaks directly to the node.
type inMemObscuroClient struct {
	nodeID     common.Address
	obscuroAPI ObscuroAPI
}

func NewInMemObscuroClient(nodeID int64, p2p *P2P, db *DB) obscuroclient.Client {
	return &inMemObscuroClient{
		obscuroAPI: *NewObscuroAPI(p2p, db),
		nodeID:     common.BigToAddress(big.NewInt(nodeID)),
	}
}

func (c *inMemObscuroClient) ID() common.Address {
	return c.nodeID
}

// Call bypasses RPC, and invokes methods on the node directly.
func (c *inMemObscuroClient) Call(result interface{}, method string, args ...interface{}) error {
	switch method {
	case obscuroclient.RPCSendTransactionEncrypted:
		// TODO - Extract this checking logic as the set of RPC operations grows.
		if len(args) != 1 {
			return fmt.Errorf("expected 1 arg to %s, got %d", obscuroclient.RPCSendTransactionEncrypted, len(args))
		}
		tx, ok := args[0].(nodecommon.EncryptedTx)
		if !ok {
			return fmt.Errorf("arg to %s was not of expected type EncryptedTx", obscuroclient.RPCSendTransactionEncrypted)
		}

		c.obscuroAPI.SendTransactionEncrypted(tx)

	case obscuroclient.RPCGetCurrentBlockHeadHeight:
		*result.(*int64) = c.obscuroAPI.GetCurrentBlockHeadHeight()

	case obscuroclient.RPCGetCurrentRollupHead:
		*result.(**nodecommon.Header) = c.obscuroAPI.GetCurrentRollupHead()
	}

	// todo - joel - return error if no match
	return nil
}

func (c *inMemObscuroClient) Stop() {
	// There is no RPC connection to close.
}
