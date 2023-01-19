package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/contracts/generated/MessageBus"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
)

// Functions to convert classes that need to be sent between the host and the enclave to and from their equivalent
// Protobuf message classes.

func ToAttestationReportMsg(report *common.AttestationReport) generated.AttestationReportMsg {
	return generated.AttestationReportMsg{Report: report.Report, PubKey: report.PubKey, Owner: report.Owner.Bytes(), HostAddress: report.HostAddress}
}

func FromAttestationReportMsg(msg *generated.AttestationReportMsg) *common.AttestationReport {
	return &common.AttestationReport{
		Report:      msg.Report,
		PubKey:      msg.PubKey,
		Owner:       gethcommon.BytesToAddress(msg.Owner),
		HostAddress: msg.HostAddress,
	}
}

func ToBlockSubmissionResponseMsg(response *common.BlockSubmissionResponse) (generated.BlockSubmissionResponseMsg, error) {
	subscribedLogBytes, err := json.Marshal(response.SubscribedLogs)
	if err != nil {
		return generated.BlockSubmissionResponseMsg{}, fmt.Errorf("could not marshal subscribed logs to JSON. Cause: %w", err)
	}

	producedBatchMsg := ToExtBatchMsg(response.ProducedBatch)
	producedRollupMsg := ToExtRollupMsg(response.ProducedRollup)

	return generated.BlockSubmissionResponseMsg{
		ProducedBatch:           &producedBatchMsg,
		ProducedRollup:          &producedRollupMsg,
		SubscribedLogs:          subscribedLogBytes,
		ProducedSecretResponses: ToSecretRespMsg(response.ProducedSecretResponses),
	}, nil
}

func ToBlockSubmissionRejectionMsg(rejectError *common.BlockRejectError) (generated.BlockSubmissionResponseMsg, error) {
	errMsg := &generated.BlockSubmissionErrorMsg{
		Cause:  rejectError.Wrapped.Error(),
		L1Head: rejectError.L1Head.Bytes(),
	}
	return generated.BlockSubmissionResponseMsg{
		Error: errMsg,
	}, nil
}

func ToSecretRespMsg(responses []*common.ProducedSecretResponse) []*generated.SecretResponseMsg {
	respMsgs := make([]*generated.SecretResponseMsg, len(responses))

	for i, resp := range responses {
		msg := generated.SecretResponseMsg{
			Secret:      resp.Secret,
			RequesterID: resp.RequesterID.Bytes(),
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
			HostAddress: msgResp.HostAddress,
		}
		respList[i] = &r
	}
	return respList
}

func FromBlockSubmissionResponseMsg(msg *generated.BlockSubmissionResponseMsg) (*common.BlockSubmissionResponse, error) {
	if msg.Error != nil {
		return nil, &common.BlockRejectError{
			L1Head:  gethcommon.BytesToHash(msg.Error.L1Head),
			Wrapped: errors.New(msg.Error.Cause),
		}
	}
	var subscribedLogs map[rpc.ID][]byte
	if err := json.Unmarshal(msg.SubscribedLogs, &subscribedLogs); err != nil {
		return nil, fmt.Errorf("could not unmarshal subscribed logs from submission response JSON. Cause: %w", err)
	}
	return &common.BlockSubmissionResponse{
		ProducedBatch:           FromExtBatchMsg(msg.ProducedBatch),
		ProducedRollup:          FromExtRollupMsg(msg.ProducedRollup),
		SubscribedLogs:          subscribedLogs,
		ProducedSecretResponses: FromSecretRespMsg(msg.ProducedSecretResponses),
	}, nil
}

func ToCrossChainMsgs(messages []MessageBus.StructsCrossChainMessage) []*generated.CrossChainMsg {
	generatedMessages := make([]*generated.CrossChainMsg, 0)

	for _, message := range messages {
		generatedMessages = append(generatedMessages, &generated.CrossChainMsg{
			Sender:   message.Sender.Bytes(),
			Sequence: message.Sequence,
			Nonce:    message.Nonce,
			Topic:    message.Topic,
			Payload:  message.Payload,
		})
	}

	return generatedMessages
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
	withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
	for _, withdrawal := range header.Withdrawals {
		withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount.Bytes(), Recipient: withdrawal.Recipient.Bytes(), Contract: withdrawal.Contract.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	diff := uint64(0)
	if header.Difficulty != nil {
		diff = header.Difficulty.Uint64()
	}
	baseFee := uint64(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee.Uint64()
	}
	headerMsg = generated.BatchHeaderMsg{
		ParentHash:                  header.ParentHash.Bytes(),
		Node:                        header.Agg.Bytes(),
		Nonce:                       []byte{},
		Proof:                       header.L1Proof.Bytes(),
		Root:                        header.Root.Bytes(),
		BodyHash:                    header.BodyHash.Bytes(),
		Number:                      header.Number.Uint64(),
		Bloom:                       header.Bloom.Bytes(),
		ReceiptHash:                 header.ReceiptHash.Bytes(),
		Extra:                       header.Extra,
		R:                           header.R.Bytes(),
		S:                           header.S.Bytes(),
		Withdrawals:                 withdrawalMsgs,
		UncleHash:                   header.UncleHash.Bytes(),
		Coinbase:                    header.Coinbase.Bytes(),
		Difficulty:                  diff,
		GasLimit:                    header.GasLimit,
		GasUsed:                     header.GasUsed,
		Time:                        header.Time,
		MixDigest:                   header.MixDigest.Bytes(),
		BaseFee:                     baseFee,
		CrossChainMessages:          ToCrossChainMsgs(header.CrossChainMessages),
		LatestInboundCrossChainHash: header.LatestInboudCrossChainHash.Bytes(),
	}

	if header.LatestInboundCrossChainHeight != nil {
		headerMsg.LatestInboundCrossChainHeight = header.LatestInboundCrossChainHeight.Bytes()
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

func FromBatchHeaderMsg(header *generated.BatchHeaderMsg) *common.BatchHeader { //nolint:dupl
	if header == nil {
		return nil
	}
	withdrawals := make([]common.Withdrawal, 0)
	for _, withdrawalMsg := range header.Withdrawals {
		recipient := gethcommon.BytesToAddress(withdrawalMsg.Recipient)
		contract := gethcommon.BytesToAddress(withdrawalMsg.Contract)
		amount := big.NewInt(0).SetBytes(withdrawalMsg.Amount)
		withdrawal := common.Withdrawal{Amount: amount, Recipient: recipient, Contract: contract}
		withdrawals = append(withdrawals, withdrawal)
	}

	r := &big.Int{}
	s := &big.Int{}
	return &common.BatchHeader{
		ParentHash:                    gethcommon.BytesToHash(header.ParentHash),
		Agg:                           gethcommon.BytesToAddress(header.Node),
		Nonce:                         types.EncodeNonce(big.NewInt(0).SetBytes(header.Nonce).Uint64()),
		L1Proof:                       gethcommon.BytesToHash(header.Proof),
		Root:                          gethcommon.BytesToHash(header.Root),
		BodyHash:                      gethcommon.BytesToHash(header.BodyHash),
		Number:                        big.NewInt(int64(header.Number)),
		Bloom:                         types.BytesToBloom(header.Bloom),
		ReceiptHash:                   gethcommon.BytesToHash(header.ReceiptHash),
		Extra:                         header.Extra,
		R:                             r.SetBytes(header.R),
		S:                             s.SetBytes(header.S),
		Withdrawals:                   withdrawals,
		UncleHash:                     gethcommon.BytesToHash(header.UncleHash),
		Coinbase:                      gethcommon.BytesToAddress(header.Coinbase),
		Difficulty:                    big.NewInt(int64(header.Difficulty)),
		GasLimit:                      header.GasLimit,
		GasUsed:                       header.GasUsed,
		Time:                          header.Time,
		MixDigest:                     gethcommon.BytesToHash(header.MixDigest),
		BaseFee:                       big.NewInt(int64(header.BaseFee)),
		CrossChainMessages:            FromCrossChainMsgs(header.CrossChainMessages),
		LatestInboudCrossChainHash:    gethcommon.BytesToHash(header.LatestInboundCrossChainHash),
		LatestInboundCrossChainHeight: big.NewInt(0).SetBytes(header.LatestInboundCrossChainHeight),
	}
}

func ToExtRollupMsg(rollup *common.ExtRollup) generated.ExtRollupMsg {
	if rollup == nil || rollup.Header == nil {
		return generated.ExtRollupMsg{}
	}

	txHashBytes := make([][]byte, len(rollup.TxHashes))
	for idx, txHash := range rollup.TxHashes {
		txHashBytes[idx] = txHash.Bytes()
	}

	return generated.ExtRollupMsg{Header: ToRollupHeaderMsg(rollup.Header), TxHashes: txHashBytes, Txs: rollup.EncryptedTxBlob}
}

func ToRollupHeaderMsg(header *common.RollupHeader) *generated.RollupHeaderMsg {
	if header == nil {
		return nil
	}
	var headerMsg generated.RollupHeaderMsg
	withdrawalMsgs := make([]*generated.WithdrawalMsg, 0)
	for _, withdrawal := range header.Withdrawals {
		withdrawalMsg := generated.WithdrawalMsg{Amount: withdrawal.Amount.Bytes(), Recipient: withdrawal.Recipient.Bytes(), Contract: withdrawal.Contract.Bytes()}
		withdrawalMsgs = append(withdrawalMsgs, &withdrawalMsg)
	}

	diff := uint64(0)
	if header.Difficulty != nil {
		diff = header.Difficulty.Uint64()
	}
	baseFee := uint64(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee.Uint64()
	}
	headerMsg = generated.RollupHeaderMsg{
		ParentHash:                  header.ParentHash.Bytes(),
		Node:                        header.Agg.Bytes(),
		Nonce:                       []byte{},
		Proof:                       header.L1Proof.Bytes(),
		Root:                        header.Root.Bytes(),
		BatchHash:                   header.BatchHash.Bytes(),
		Number:                      header.Number.Uint64(),
		Bloom:                       header.Bloom.Bytes(),
		ReceiptHash:                 header.ReceiptHash.Bytes(),
		Extra:                       header.Extra,
		R:                           header.R.Bytes(),
		S:                           header.S.Bytes(),
		Withdrawals:                 withdrawalMsgs,
		UncleHash:                   header.UncleHash.Bytes(),
		Coinbase:                    header.Coinbase.Bytes(),
		Difficulty:                  diff,
		GasLimit:                    header.GasLimit,
		GasUsed:                     header.GasUsed,
		Time:                        header.Time,
		MixDigest:                   header.MixDigest.Bytes(),
		BaseFee:                     baseFee,
		CrossChainMessages:          ToCrossChainMsgs(header.CrossChainMessages),
		LatestInboundCrossChainHash: header.LatestInboudCrossChainHash.Bytes(),
	}

	if header.LatestInboundCrossChainHeight != nil {
		headerMsg.LatestInboundCrossChainHeight = header.LatestInboundCrossChainHeight.Bytes()
	}

	return &headerMsg
}

func FromExtRollupMsg(msg *generated.ExtRollupMsg) *common.ExtRollup {
	if msg.Header == nil {
		return &common.ExtRollup{
			Header: nil,
		}
	}

	// We recreate the transaction hashes.
	txHashes := make([]gethcommon.Hash, len(msg.TxHashes))
	for idx, bytes := range msg.TxHashes {
		txHashes[idx] = gethcommon.BytesToHash(bytes)
	}

	return &common.ExtRollup{
		Header:          FromRollupHeaderMsg(msg.Header),
		TxHashes:        txHashes,
		EncryptedTxBlob: msg.Txs,
	}
}

func FromRollupHeaderMsg(header *generated.RollupHeaderMsg) *common.RollupHeader { //nolint:dupl
	if header == nil {
		return nil
	}
	withdrawals := make([]common.Withdrawal, 0)
	for _, withdrawalMsg := range header.Withdrawals {
		recipient := gethcommon.BytesToAddress(withdrawalMsg.Recipient)
		contract := gethcommon.BytesToAddress(withdrawalMsg.Contract)
		amount := big.NewInt(0).SetBytes(withdrawalMsg.Amount)
		withdrawal := common.Withdrawal{Amount: amount, Recipient: recipient, Contract: contract}
		withdrawals = append(withdrawals, withdrawal)
	}

	r := &big.Int{}
	s := &big.Int{}
	return &common.RollupHeader{
		ParentHash:                    gethcommon.BytesToHash(header.ParentHash),
		Agg:                           gethcommon.BytesToAddress(header.Node),
		Nonce:                         types.EncodeNonce(big.NewInt(0).SetBytes(header.Nonce).Uint64()),
		L1Proof:                       gethcommon.BytesToHash(header.Proof),
		Root:                          gethcommon.BytesToHash(header.Root),
		BatchHash:                     gethcommon.BytesToHash(header.BatchHash),
		Number:                        big.NewInt(int64(header.Number)),
		Bloom:                         types.BytesToBloom(header.Bloom),
		ReceiptHash:                   gethcommon.BytesToHash(header.ReceiptHash),
		Extra:                         header.Extra,
		R:                             r.SetBytes(header.R),
		S:                             s.SetBytes(header.S),
		Withdrawals:                   withdrawals,
		UncleHash:                     gethcommon.BytesToHash(header.UncleHash),
		Coinbase:                      gethcommon.BytesToAddress(header.Coinbase),
		Difficulty:                    big.NewInt(int64(header.Difficulty)),
		GasLimit:                      header.GasLimit,
		GasUsed:                       header.GasUsed,
		Time:                          header.Time,
		MixDigest:                     gethcommon.BytesToHash(header.MixDigest),
		BaseFee:                       big.NewInt(int64(header.BaseFee)),
		CrossChainMessages:            FromCrossChainMsgs(header.CrossChainMessages),
		LatestInboudCrossChainHash:    gethcommon.BytesToHash(header.LatestInboundCrossChainHash),
		LatestInboundCrossChainHeight: big.NewInt(0).SetBytes(header.LatestInboundCrossChainHeight),
	}
}
