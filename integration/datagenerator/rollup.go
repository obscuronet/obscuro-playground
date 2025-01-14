package datagenerator

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// RandomRollup - block is needed in order to pass the smart contract check
// when submitting cross chain messages.
func RandomRollup(_ *types.Header) common.ExtRollup {
	extRollup := common.ExtRollup{
		Header: &common.RollupHeader{},
	}

	return extRollup
}
