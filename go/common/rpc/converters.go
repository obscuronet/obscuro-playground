package rpc

import (
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/rpc/generated"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func ToAttestationReportMsg(report *common.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Report: report.Report, PubKey: report.PubKey, EnclaveID: report.EnclaveID.Bytes(), HostAddress: report.HostAddress}
}

func FromAttestationReportMsg(msg *generated.AttestationReportMsg) *common.AttestationReport {
	return &common.AttestationReport{
		Report:      msg.Report,
		PubKey:      msg.PubKey,
		EnclaveID:   gethcommon.BytesToAddress(msg.EnclaveID),
		HostAddress: msg.HostAddress,
	}
}

func ToBlockSubmissionResponseMsg(response *common.BlockSubmissionResponse) (*generated.BlockSubmissionResponseMsg, error) {
	if response == nil {
		return nil, fmt.Errorf("no response that could be converted to a message")
	}

	msg := &generated.BlockSubmissionResponseMsg{
		ProducedSecretResponses: ToSecretRespMsg(response.ProducedSecretResponses),
	}

	return msg, nil
}

func ToSecretRespMsg(responses []*common.ProducedSecretResponse) []*generated.SecretResponseMsg {
	respMsgs := make([]*generated.SecretResponseMsg, len(responses))

	for i, resp := range responses {
		msg := generated.SecretResponseMsg{
			Secret:      resp.Secret,
			RequesterID: resp.RequesterID.Bytes(),
			AttesterID:  resp.AttesterID.Bytes(),
			HostAddress: resp.HostAddress,
		}
		respMsgs[i] = &msg
	}

	return respMsgs
}

func FromSecretRespMsg(secretResponses []*generated.SecretResponseMsg) []*common.ProducedSecretResponse {
	respList := make([]*common.ProducedSecretResponse, len(secretResponses))

	for i, msgResp := range secretResponses {
		r := common.ProducedSecretResponse{
			Secret:      msgResp.Secret,
			RequesterID: gethcommon.BytesToAddress(msgResp.RequesterID),
			AttesterID:  gethcommon.BytesToAddress(msgResp.AttesterID),
			HostAddress: msgResp.HostAddress,
		}
		respList[i] = &r
	}
	return respList
}

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) (*common.BlockSubmissionResponse, error) {
	return &common.BlockSubmissionResponse{
		ProducedSecretResponses: FromSecretRespMsg(msg.ProducedSecretResponses),
	}, nil
}

func FromCrossChainMsgs(messages []*generated.CrossChainMsg) []MessageBus.StructsCrossChainMessage {
	outMessages := make([]MessageBus.StructsCrossChainMessage, 0)

	for _, message := range messages {
		outMessages = append(outMessages, MessageBus.StructsCrossChainMessage{
			Sender:   gethcommon.BytesToAddress(message.Sender),
			Sequence: message.Sequence,
			Nonce:    message.Nonce,
			Topic:    message.Topic,
			Payload:  message.Payload,
		})
	}

	return outMessages
}

func ToExtBatchMsg(batch *common.ExtBatch) generated.ExtBatchMsg {
	if batch == nil || batch.Header == nil {
		return generated.ExtBatchMsg{}
	}

	txHashBytes := make([][]byte, len(batch.TxHashes))
	for idx, txHash := range batch.TxHashes {
		txHashBytes[idx] = txHash.Bytes()
	}

	return generated.ExtBatchMsg{Header: ToBatchHeaderMsg(batch.Header), TxHashes: txHashBytes, Txs: batch.EncryptedTxBlob}
}

func ToBatchHeaderMsg(header *common.BatchHeader) *generated.BatchHeaderMsg {
	if header == nil {
		return nil
	}
	var headerMsg generated.BatchHeaderMsg

	baseFee := uint64(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee.Uint64()
	}
	headerMsg = generated.BatchHeaderMsg{
		ParentHash:       header.ParentHash.Bytes(),
		Proof:            header.L1Proof.Bytes(),
		Root:             header.Root.Bytes(),
		TxHash:           header.TxHash.Bytes(),
		Number:           header.Number.Uint64(),
		SequencerOrderNo: header.SequencerOrderNo.Uint64(),
		ReceiptHash:      header.ReceiptHash.Bytes(),
		Extra:            header.Extra,
		Signature:        header.Signature,
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Time:             header.Time,
		BaseFee:          baseFee,
		TransferTree:     header.CrossChainRoot.Bytes(),
		Coinbase:         header.Coinbase.Bytes(),
		CrossChainTree:   header.CrossChainTree,
	}

	return &headerMsg
}

func FromExtBatchMsg(msg *generated.ExtBatchMsg) *common.ExtBatch {
	if msg.Header == nil {
		return &common.ExtBatch{
			Header: nil,
		}
	}

	// We recreate the transaction hashes.
	txHashes := make([]gethcommon.Hash, len(msg.TxHashes))
	for idx, bytes := range msg.TxHashes {
		txHashes[idx] = gethcommon.BytesToHash(bytes)
	}

	return &common.ExtBatch{
		Header:          FromBatchHeaderMsg(msg.Header),
		TxHashes:        txHashes,
		EncryptedTxBlob: msg.Txs,
	}
}

func FromBatchHeaderMsg(header *generated.BatchHeaderMsg) *common.BatchHeader {
	if header == nil {
		return nil
	}

	return &common.BatchHeader{
		ParentHash:       gethcommon.BytesToHash(header.ParentHash),
		L1Proof:          gethcommon.BytesToHash(header.Proof),
		Root:             gethcommon.BytesToHash(header.Root),
		TxHash:           gethcommon.BytesToHash(header.TxHash),
		Number:           big.NewInt(int64(header.Number)),
		SequencerOrderNo: big.NewInt(int64(header.SequencerOrderNo)),
		ReceiptHash:      gethcommon.BytesToHash(header.ReceiptHash),
		Extra:            header.Extra,
		Signature:        header.Signature,
		GasLimit:         header.GasLimit,
		GasUsed:          header.GasUsed,
		Time:             header.Time,
		CrossChainRoot:   gethcommon.BytesToHash(header.TransferTree),
		BaseFee:          big.NewInt(0).SetUint64(header.BaseFee),
		Coinbase:         gethcommon.BytesToAddress(header.Coinbase),
		CrossChainTree:   header.CrossChainTree,
	}
}

func ToExtRollupMsg(rollup *common.ExtRollup) generated.ExtRollupMsg {
	if rollup == nil || rollup.Header == nil {
		return generated.ExtRollupMsg{}
	}

	return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), BatchPayloads: rollup.BatchPayloads, CalldataRollupHeader: rollup.CalldataRollupHeader}
}

func ToRollupHeaderMsg(header *common.RollupHeader) *generated.RollupHeaderMsg {
	if header == nil {
		return nil
	}
	headerMsg := generated.RollupHeaderMsg{
		BlockHash:      header.BlockHash.Bytes(),
		Signature:      header.Signature,
		LastBatchSeqNo: header.LastBatchSeqNo,
	}

	return &headerMsg
}

func FromExtRollupMsg(msg *generated.ExtRollupMsg) *common.ExtRollup {
	if msg.Header == nil {
		return &common.ExtRollup{
			Header: nil,
		}
	}

	return &common.ExtRollup{
		Header:               FromRollupHeaderMsg(msg.Header),
		BatchPayloads:        msg.BatchPayloads,
		CalldataRollupHeader: msg.CalldataRollupHeader,
	}
}

func FromRollupHeaderMsg(header *generated.RollupHeaderMsg) *common.RollupHeader {
	if header == nil {
		return nil
	}

	return &common.RollupHeader{
		BlockHash:      gethcommon.BytesToHash(header.BlockHash),
		Signature:      header.Signature,
		LastBatchSeqNo: header.LastBatchSeqNo,
	}
}

func ToRollupDataMsg(rollupData *common.PublicRollupMetadata) generated.PublicRollupDataMsg {
	if rollupData == nil {
		return generated.PublicRollupDataMsg{}
	}

	return generated.PublicRollupDataMsg{StartSeq: rollupData.FirstBatchSequence.Uint64(), Timestamp: rollupData.StartTime}
}

func FromRollupDataMsg(msg *generated.PublicRollupDataMsg) (*common.PublicRollupMetadata, error) {
	if msg.Timestamp == 0 {
		return nil, fmt.Errorf("timestamp on the rollup can not be zero")
	}

	if msg.StartSeq == 0 {
		return &common.PublicRollupMetadata{
			FirstBatchSequence: nil,
			StartTime:          msg.Timestamp,
		}, nil
	}

	return &common.PublicRollupMetadata{
		FirstBatchSequence: big.NewInt(int64(msg.StartSeq)),
		StartTime:          msg.Timestamp,
	}, nil
}
