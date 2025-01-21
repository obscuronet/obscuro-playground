package common

import "fmt"

const (
	sequencer = "sequencer"
	validator = "validator"
	unknown   = "unknown"
)

// NodeType represents the type of the node.
type NodeType int

const (
	Validator NodeType = iota
	Sequencer
	Unknown
)

func (n NodeType) String() string {
	switch n {
	case Validator:
		return validator
	case Sequencer:
		return sequencer
	case Unknown:
		return unknown
	default:
		return unknown
	}
}

func (n *NodeType) UnmarshalText(text []byte) error {
	nodeType, err := ToNodeType(string(text))
	if err != nil {
		return err
	}
	*n = nodeType
	return nil
}

func ToNodeType(s string) (NodeType, error) {
	switch s {
	case validator:
		return Validator, nil
	case sequencer:
		return Sequencer, nil
	default:
		return Unknown, fmt.Errorf("string '%s' cannot be converted to a node type", s)
	}
}
