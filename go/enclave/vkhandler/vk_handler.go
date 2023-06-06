package vkhandler

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidAddressSignature = fmt.Errorf("invalid viewing key signature for requested address")
)

// ViewingKeySignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const ViewingKeySignedMsgPrefix = "vk"

// Used when the result to an eth_call is equal to nil. Attempting to encrypt then decrypt nil using ECIES throws an exception.
var placeholderResult = []byte("0x")

type VKHandler struct {
	publicViewingKey *ecies.PublicKey
}

func (m *VKHandler) Encrypt(bytes []byte) ([]byte, error) {
	if len(bytes) == 0 {
		bytes = placeholderResult
	}

	encryptedBytes, err := ecies.Encrypt(rand.Reader, m.publicViewingKey, bytes, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt with given public VK - %w", err)
	}

	return encryptedBytes, nil
}

// New creates a new viewing key handler
func New(targetAddress *gethcommon.Address, vkPubKeyBytes, accountSignatureHexBytes []byte) (*VKHandler, error) {
	// Recalculate the message signed by MetaMask.
	msgToSign := viewingkey.GenerateSignMessage(vkPubKeyBytes)

	// We recover the key based on the signed message and the signature.
	recoveredAccountPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), accountSignatureHexBytes)
	if err != nil {
		return nil, fmt.Errorf("viewing key but could not validate its signature - %w", err)
	}
	recoveredAccountAddress := crypto.PubkeyToAddress(*recoveredAccountPublicKey)

	if targetAddress.Hash() != recoveredAccountAddress.Hash() {
		return nil, ErrInvalidAddressSignature
	}

	// We decompress the viewing key and create the corresponding ECIES key.
	viewingKey, err := crypto.DecompressPubkey(vkPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not decompress viewing key bytes - %w", err)
	}

	return &VKHandler{
		publicViewingKey: ecies.ImportECDSAPublic(viewingKey),
	}, nil
}
