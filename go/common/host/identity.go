package host

import (
	"github.com/ten-protocol/go-ten/go/common"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
)

type Identity struct {
	ID               string
	P2PPublicAddress string
	IsGenesis        bool
	IsSequencer      bool
}

func NewIdentity(cfg *hostconfig.HostConfig) Identity {
	return Identity{
		ID:               cfg.ID,
		P2PPublicAddress: cfg.P2PPublicAddress,
		IsGenesis:        cfg.IsGenesis,
		IsSequencer:      cfg.NodeType == common.Sequencer,
	}
}
