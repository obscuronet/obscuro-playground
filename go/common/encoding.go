package common

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeRollup(r *ExtRollup) (EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}

func DecodeRollup(encoded EncodedRollup) (*ExtRollup, error) {
	r := new(ExtRollup)
	err := rlp.DecodeBytes(encoded, r)
	return r, err
}

func EncodeAttestation(att *AttestationReport) (EncodedAttestationReport, error) {
	return rlp.EncodeToBytes(att)
}

func DecodeAttestation(encoded EncodedAttestationReport) (*AttestationReport, error) {
	att := new(AttestationReport)
	err := rlp.DecodeBytes(encoded, att)
	return att, err
}
