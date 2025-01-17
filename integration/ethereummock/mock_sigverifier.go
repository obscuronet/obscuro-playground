package ethereummock

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type MockSignatureValidator struct{}

func NewMockSignatureValidator() *MockSignatureValidator {
	return &MockSignatureValidator{}
}

// CheckSequencerSignature - NO-OP
func (sigChecker *MockSignatureValidator) CheckSequencerSignature(_ gethcommon.Hash, _ []byte) error {
	return nil
}
