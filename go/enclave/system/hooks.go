package system

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ten-protocol/go-ten/contracts/generated/TransactionsAnalyzer"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/wallet"
)

var (
	transactionsAnalyzerABI, _ = abi.JSON(strings.NewReader(TransactionsAnalyzer.TransactionsAnalyzerMetaData.ABI))
)

type SystemContractCallbacks interface {
	GetOwner() gethcommon.Address
	Initialize(batch *core.Batch, receipts types.Receipts) error
	Load() error
	CreateOnBatchEndTransaction(ctx context.Context, l2State *state.StateDB, batch *core.Batch, receipts common.L2Receipts) (*common.L2Tx, error)
	TransactionAnalyzerAddress() *gethcommon.Address
}

type systemContractCallbacks struct {
	transactionsAnalyzerAddress *gethcommon.Address
	ownerWallet                 wallet.Wallet
	storage                     storage.Storage

	logger gethlog.Logger
}

func NewSystemContractCallbacks(ownerWallet wallet.Wallet, logger gethlog.Logger) SystemContractCallbacks {
	return &systemContractCallbacks{
		transactionsAnalyzerAddress: nil,
		ownerWallet:                 ownerWallet,
		logger:                      logger,
		storage:                     nil,
	}
}

func (s *systemContractCallbacks) TransactionAnalyzerAddress() *gethcommon.Address {
	return s.transactionsAnalyzerAddress
}

func (s *systemContractCallbacks) GetOwner() gethcommon.Address {
	return s.ownerWallet.Address()
}

func (s *systemContractCallbacks) Load() error {
	s.logger.Info("Load: Initializing system contracts")

	if s.storage == nil {
		s.logger.Error("Load: Storage is not set")
		return fmt.Errorf("storage is not set")
	}

	batchSeqNo := uint64(1)
	s.logger.Debug("Load: Fetching batch", "batchSeqNo", batchSeqNo)
	batch, err := s.storage.FetchBatchBySeqNo(context.Background(), batchSeqNo)
	if err != nil {
		s.logger.Error("Load: Failed fetching batch", "batchSeqNo", batchSeqNo, "error", err)
		return fmt.Errorf("failed fetching batch %w", err)
	}

	if len(batch.Transactions) < 2 {
		s.logger.Error("Load: Genesis batch does not have enough transactions", "batchSeqNo", batchSeqNo, "transactionCount", len(batch.Transactions))
		return fmt.Errorf("genesis batch does not have enough transactions")
	}

	receipt, err := s.storage.GetTransactionReceipt(context.Background(), batch.Transactions[1].Hash())
	if err != nil {
		s.logger.Error("Load: Failed fetching receipt", "transactionHash", batch.Transactions[1].Hash().Hex(), "error", err)
		return fmt.Errorf("failed fetching receipt %w", err)
	}

	addresses, err := DeriveAddresses(receipt)
	if err != nil {
		s.logger.Error("Load: Failed deriving addresses", "error", err, "receiptHash", receipt.TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) initializeRequiredAddresses(addresses SystemContractAddresses) error {
	if addresses["TransactionsAnalyzer"] == nil {
		return fmt.Errorf("required contract address TransactionsAnalyzer is nil")
	}

	s.transactionsAnalyzerAddress = addresses["TransactionsAnalyzer"]

	return nil
}

func (s *systemContractCallbacks) Initialize(batch *core.Batch, receipts types.Receipts) error {
	s.logger.Info("Initialize: Starting initialization of system contracts", "batchSeqNo", batch.SeqNo)

	if len(receipts) < 2 {
		s.logger.Error("Initialize: Genesis batch does not have enough receipts", "expected", 2, "got", len(receipts))
		return fmt.Errorf("genesis batch does not have enough receipts")
	}

	receiptIndex := 1
	s.logger.Debug("Initialize: Deriving addresses from receipt", "receiptIndex", receiptIndex, "transactionHash", receipts[receiptIndex].TxHash.Hex())
	addresses, err := DeriveAddresses(receipts[receiptIndex])
	if err != nil {
		s.logger.Error("Initialize: Failed deriving addresses", "error", err, "receiptHash", receipts[receiptIndex].TxHash.Hex())
		return fmt.Errorf("failed deriving addresses %w", err)
	}

	s.logger.Info("Initialize: Initializing required addresses", "addresses", addresses)
	return s.initializeRequiredAddresses(addresses)
}

func (s *systemContractCallbacks) CreateOnBatchEndTransaction(ctx context.Context, l2State *state.StateDB, batch *core.Batch, receipts common.L2Receipts) (*common.L2Tx, error) {
	if s.transactionsAnalyzerAddress == nil {
		s.logger.Debug("CreateOnBatchEndTransaction: TransactionsAnalyzerAddress is nil, skipping transaction creation")
		return nil, nil
	}

	s.logger.Info("CreateOnBatchEndTransaction: Creating transaction on batch end", "batchSeqNo", batch.SeqNo)

	nonceForSyntheticTx := l2State.GetNonce(s.GetOwner())
	s.logger.Debug("CreateOnBatchEndTransaction: Retrieved nonce for synthetic transaction", "nonce", nonceForSyntheticTx)

	blockTransactions := TransactionsAnalyzer.TransactionsAnalyzerBlockTransactions{
		Transactions: make([][]byte, 0),
	}
	for _, tx := range batch.Transactions {
		encodedBytes, err := rlp.EncodeToBytes(tx)
		if err != nil {
			s.logger.Error("CreateOnBatchEndTransaction: Failed encoding transaction", "transactionHash", tx.Hash().Hex(), "error", err)
			return nil, fmt.Errorf("failed encoding transaction for onBlock %w", err)
		}

		blockTransactions.Transactions = append(blockTransactions.Transactions, encodedBytes)
		s.logger.Debug("CreateOnBatchEndTransaction: Encoded transaction", "transactionHash", tx.Hash().Hex())
	}

	data, err := transactionsAnalyzerABI.Pack("onBlock", blockTransactions)
	if err != nil {
		s.logger.Error("CreateOnBatchEndTransaction: Failed packing onBlock data", "error", err)
		return nil, fmt.Errorf("failed packing onBlock() %w", err)
	}

	tx := &types.LegacyTx{
		Nonce:    nonceForSyntheticTx,
		Value:    gethcommon.Big0,
		Gas:      500_000_000,
		GasPrice: gethcommon.Big0, // Synthetic transactions are on the house. Or the house.
		Data:     data,
		To:       s.transactionsAnalyzerAddress,
	}

	s.logger.Debug("CreateOnBatchEndTransaction: Signing transaction", "to", s.transactionsAnalyzerAddress.Hex(), "nonce", nonceForSyntheticTx)
	signedTx, err := s.ownerWallet.SignTransaction(tx)
	if err != nil {
		s.logger.Error("CreateOnBatchEndTransaction: Failed signing transaction", "error", err)
		return nil, fmt.Errorf("failed signing transaction %w", err)
	}

	s.logger.Info("CreateOnBatchEndTransaction: Successfully created signed transaction", "transactionHash", signedTx.Hash().Hex())
	return signedTx, nil
}
