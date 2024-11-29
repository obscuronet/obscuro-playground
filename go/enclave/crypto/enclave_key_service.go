package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/log"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/signature"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common"
)

// enclavePrivateKey - the private key that was attested
type enclavePrivateKey struct {
	privateKey     *ecdsa.PrivateKey // generated by the enclave at startup, used to sign messages
	enclaveID      common.EnclaveID  // the enclave's ID, derived from the public key
	publicKeyBytes []byte
}

// EnclaveAttestedKeyService manages the attestation key - including
type EnclaveAttestedKeyService struct {
	logger     gethlog.Logger
	enclaveKey *enclavePrivateKey
}

func NewEnclaveAttestedKeyService(logger gethlog.Logger) *EnclaveAttestedKeyService {
	return &EnclaveAttestedKeyService{logger: logger}
}

func (eks *EnclaveAttestedKeyService) Sign(payload gethcommon.Hash) ([]byte, error) {
	return signature.Sign(payload.Bytes(), eks.enclaveKey.privateKey)
}

func (eks *EnclaveAttestedKeyService) EnclaveID() common.EnclaveID {
	return eks.enclaveKey.enclaveID
}

func (eks *EnclaveAttestedKeyService) PublicKey() *ecdsa.PublicKey {
	return &eks.enclaveKey.privateKey.PublicKey
}

func (eks *EnclaveAttestedKeyService) PublicKeyBytes() []byte {
	return eks.enclaveKey.publicKeyBytes
}

func (eks *EnclaveAttestedKeyService) Decrypt(encBytes []byte) ([]byte, error) {
	return decryptWithPrivateKey(encBytes, eks.enclaveKey.privateKey)
}

func (eks *EnclaveAttestedKeyService) Encrypt(encBytes []byte) ([]byte, error) {
	return encryptWithPublicKey(encBytes, &eks.enclaveKey.privateKey.PublicKey)
}

func (eks *EnclaveAttestedKeyService) GenerateEnclaveKey() ([]byte, error) {
	privKey, err := ecdsa.GenerateKey(gethcrypto.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not generate enclave key. Cause: %w", err)
	}
	return gethcrypto.FromECDSA(privKey), nil
}

func (eks *EnclaveAttestedKeyService) SetEnclaveKey(keyBytes []byte) {
	ecdsaKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		eks.logger.Crit("could not parse enclave key", log.ErrKey, err)
	}
	pubKey := gethcrypto.CompressPubkey(&ecdsaKey.PublicKey)
	enclaveID := gethcrypto.PubkeyToAddress(ecdsaKey.PublicKey)
	eks.enclaveKey = &enclavePrivateKey{
		privateKey:     ecdsaKey,
		publicKeyBytes: pubKey,
		enclaveID:      enclaveID,
	}
}
