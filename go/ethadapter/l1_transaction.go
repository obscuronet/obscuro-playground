package ethadapter

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// L1Transaction is an abstraction that transforms an Ethereum transaction into a format that can be consumed more
// easily by TEN.
type L1Transaction interface{}

type L1RollupTx struct {
	Rollup common.EncodedRollup
}

type L1RollupHashes struct {
	BlobHashes []gethcommon.Hash
}

type L1DepositTx struct {
	Amount        *big.Int            // Amount to be deposited
	To            *gethcommon.Address // Address the ERC20 Transfer was made to (always be the Management Contract Addr)
	Sender        *gethcommon.Address // Address that issued the ERC20, the token holder or tx.origin
	TokenContract *gethcommon.Address // Address of the ERC20 Contract address that was executed
}

type L1RespondSecretTx struct {
	Secret      []byte
	RequesterID gethcommon.Address
	AttesterID  gethcommon.Address
	AttesterSig []byte
}

type L1SetImportantContractsTx struct {
	Key        string
	NewAddress gethcommon.Address
}

// Sign signs the payload with a given private key
func (l *L1RespondSecretTx) Sign(privateKey *ecdsa.PrivateKey) *L1RespondSecretTx {
	var data []byte
	data = append(data, l.AttesterID.Bytes()...)
	data = append(data, l.RequesterID.Bytes()...)
	data = append(data, string(l.Secret)...)

	ethereumMessageHash := func(data []byte) []byte {
		prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(data))
		return crypto.Keccak256([]byte(prefix), data)
	}

	hashedData := ethereumMessageHash(data)
	// sign the hash
	signedHash, err := crypto.Sign(hashedData, privateKey)
	if err != nil {
		return nil
	}

	// set recovery id to 27; prevent malleable signatures
	signedHash[64] += 27
	l.AttesterSig = signedHash
	return l
}

type L1RequestSecretTx struct {
	Attestation common.EncodedAttestationReport
}

type L1InitializeSecretTx struct {
	EnclaveID     *gethcommon.Address
	InitialSecret []byte
	Attestation   common.EncodedAttestationReport
}

type L1TxType int

const (
	RollupTx L1TxType = iota
	SecretRequestTx
	InitialiseSecretTx
	CrossChainMessageTx
	CrossChainValueTranserTx
	SequencerAddedTx
	SetImportantContractsTx
)

// ProcessedL1Data is submitted to the enclave by the guardian
type ProcessedL1Data struct {
	BlockHeader *types.Header
	Events      map[L1TxType][]*L1TxData
}

// L1TxData represents an L1 transaction that's relevant to us
type L1TxData struct {
	Type               *L1Transaction
	Transaction        *types.Transaction
	Receipt            *types.Receipt
	Blobs              []*kzg4844.Blob             // Only populated for blob transactions
	CrossChainMessages *common.CrossChainMessages  // Only populated for xchain messages
	ValueTransfers     *common.ValueTransferEvents // Only populated for xchain transfers
}
