package common

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// L1TenEventType represents different types of L1 transactions we monitor for
type L1TenEventType uint8 // Change to uint8 for RLP serialization

const (
	RollupTx L1TenEventType = iota
	InitialiseSecretTx
	SecretRequestTx
	SecretResponseTx
	CrossChainMessageTx
	CrossChainValueTranserTx
	SequencerAddedTx
	SequencerRevokedTx
	SetImportantContractsTx
)

// ProcessedL1Data is submitted to the enclave by the guardian
type ProcessedL1Data struct {
	BlockHeader *types.Header
	Events      []L1Event
}

// L1Event represents a single event type and its associated transactions
type L1Event struct {
	Type uint8
	Txs  []*L1TxData
}

// L1TxData represents an L1 transaction that are relevant to us
type L1TxData struct {
	Transaction        *types.Transaction
	Receipt            *types.Receipt
	Blobs              []*kzg4844.Blob     // Only populated for blob transactions
	SequencerEnclaveID gethcommon.Address  // Only non-zero when a new enclave is added as a sequencer
	CrossChainMessages CrossChainMessages  // Only populated for xchain messages
	ValueTransfers     ValueTransferEvents // Only populated for xchain transfers
	Proof              []byte              // Some merkle proof TBC
}

// HasSequencerEnclaveID helper method to check if SequencerEnclaveID is set to avoid custom RLP when we send over grpc
func (tx *L1TxData) HasSequencerEnclaveID() bool {
	return tx.SequencerEnclaveID != (gethcommon.Address{})
}

func (p *ProcessedL1Data) AddEvent(tenEventType L1TenEventType, tx *L1TxData) {
	eventType := uint8(tenEventType)

	for i := range p.Events {
		if p.Events[i].Type != eventType {
			continue
		}

		txHash := tx.Transaction.Hash()

		// check for duplicate transaction
		for _, existingTx := range p.Events[i].Txs {
			if existingTx.Transaction.Hash() == txHash {
				return // Skip duplicate transaction
			}
		}

		p.Events[i].Txs = append(p.Events[i].Txs, tx)
		return
	}

	p.Events = append(p.Events, L1Event{
		Type: eventType,
		Txs:  []*L1TxData{tx},
	})
}

func (p *ProcessedL1Data) GetEvents(txType L1TenEventType) []*L1TxData {
	if p == nil || len(p.Events) == 0 {
		return nil
	}

	for _, event := range p.Events {
		if event.Type == uint8(txType) {
			if event.Txs == nil {
				return nil
			}
			return event.Txs
		}
	}
	return nil
}
