package components

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/signature"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type SignatureValidator struct {
	SequencerID gethcommon.Address
	attestedKey *ecdsa.PublicKey
	storage     storage.Storage
}

func NewSignatureValidator(seqID gethcommon.Address, storage storage.Storage) (*SignatureValidator, error) {
	// todo (#718) - sequencer identities should be retrieved from the L1 management contract
	return &SignatureValidator{
		SequencerID: seqID,
		storage:     storage,
		attestedKey: nil,
	}, nil
}

// CheckSequencerSignature - verifies the signature against the registered sequencer
func (sigChecker *SignatureValidator) CheckSequencerSignature(headerHash gethcommon.Hash, sig []byte) error {
	if sig == nil {
		return fmt.Errorf("missing signature on batch")
	}

	// todo (@matt) disabling sequencer signature verification for now while we transition to EnclaveIDs
	// This must be re-enabled once sequencer enclaveIDs are available from the management contract

	//if sigChecker.attestedKey == nil {
	//	attestedKey, err := sigChecker.storage.FetchAttestedKey(sigChecker.SequencerID)
	//	if err != nil {
	//		return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", sigChecker.SequencerID, err)
	//	}
	//	sigChecker.attestedKey = attestedKey
	//}
	//
	// return signature.VerifySignature(sigChecker.attestedKey, headerHash.Bytes(), sig)

	return nil
}
