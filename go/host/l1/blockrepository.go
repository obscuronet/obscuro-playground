package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/go/common/host"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

var (
	// todo (@matt) make this configurable?
	_timeoutNoBlocks = 30 * time.Second
	one              = big.NewInt(1)
	ErrNoNextBlock   = errors.New("no next block")
)

type ContractType int

const (
	MgmtContract ContractType = iota
	MsgBus
)

// Repository is a host service for subscribing to new blocks and looking up L1 data
type Repository struct {
	blockSubscribers *subscription.Manager[host.L1BlockHandler]
	// this eth client should only be used by the repository, the repository may "reconnect" it at any time and don't want to interfere with other processes
	ethClient       ethadapter.EthClient
	logger          gethlog.Logger
	mgmtContractLib mgmtcontractlib.MgmtContractLib
	blobResolver    BlobResolver

	running           atomic.Bool
	head              gethcommon.Hash
	contractAddresses map[ContractType][]gethcommon.Address
}

func NewL1Repository(
	ethClient ethadapter.EthClient,
	logger gethlog.Logger,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	blobResolver BlobResolver,
	contractAddresses map[ContractType][]gethcommon.Address,
) *Repository {
	return &Repository{
		blockSubscribers:  subscription.NewManager[host.L1BlockHandler](),
		ethClient:         ethClient,
		running:           atomic.Bool{},
		logger:            logger,
		mgmtContractLib:   mgmtContractLib,
		blobResolver:      blobResolver,
		contractAddresses: contractAddresses,
	}
}

func (r *Repository) Start() error {
	r.running.Store(true)

	// Repository constantly streams new blocks and forwards them to subscribers
	go r.streamLiveBlocks()
	return nil
}

func (r *Repository) Stop() error {
	r.running.Store(false)
	return nil
}

func (r *Repository) HealthStatus(context.Context) host.HealthStatus {
	// todo (@matt) do proper health status based on last received block or something
	errMsg := ""
	if !r.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

// Subscribe will register a new block handler to receive new blocks as they arrive, returns unsubscribe func
func (r *Repository) Subscribe(handler host.L1BlockHandler) func() {
	return r.blockSubscribers.Subscribe(handler)
}

// FetchNextBlock calculates the next canonical block that should be sent to requester after a given hash.
// It returns the block and a bool for whether it is the latest known head
func (r *Repository) FetchNextBlock(prevBlockHash gethcommon.Hash) (*types.Block, bool, error) {
	if prevBlockHash == r.head {
		// prevBlock is the latest known head
		return nil, false, ErrNoNextBlock
	}

	if prevBlockHash == (gethcommon.Hash{}) {
		// prevBlock is empty, so we are starting from genesis
		blk, err := r.ethClient.BlockByNumber(big.NewInt(0))
		if err != nil {
			return nil, false, fmt.Errorf("could not find genesis block - %w", err)
		}
		return blk, false, nil
	}

	// the latestCanonAncestor will usually return the prevBlock itself but this step is necessary to walk back if there was a fork
	lca, err := r.latestCanonAncestor(prevBlockHash)
	if err != nil {
		return nil, false, err
	}
	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := r.ethClient.BlockByNumber(increment(lca.Number()))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, false, ErrNoNextBlock
		}
		return nil, false, fmt.Errorf("could not find block after latest canon ancestor, height=%s - %w", increment(lca.Number()), err)
	}

	return blk, blk.Hash() == r.head, nil
}

func (r *Repository) latestCanonAncestor(blkHash gethcommon.Hash) (*types.Block, error) {
	blk, err := r.ethClient.BlockByHash(blkHash)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block with hash=%s - %w", blkHash, err)
	}
	canonAtSameHeight, err := r.ethClient.BlockByNumber(blk.Number())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block at height=%d - %w", blk.Number(), err)
	}
	if blk.Hash() != canonAtSameHeight.Hash() {
		return r.latestCanonAncestor(blk.ParentHash())
	}
	return blk, nil
}

// FetchObscuroReceipts returns all obscuro-relevant receipts for an L1 block
func (r *Repository) FetchObscuroReceipts(block *common.L1Block) (types.Receipts, error) {
	receipts := make([]*types.Receipt, len(block.Transactions()))
	if len(block.Transactions()) == 0 {
		return receipts, nil
	}

	blkHash := block.Hash()
	// we want to send receipts for any transactions that produced obscuro-relevant log events
	var allAddresses []gethcommon.Address
	allAddresses = append(allAddresses, r.contractAddresses[MgmtContract]...)
	allAddresses = append(allAddresses, r.contractAddresses[MsgBus]...)
	logs, err := r.ethClient.GetLogs(ethereum.FilterQuery{BlockHash: &blkHash, Addresses: allAddresses})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch logs for L1 block - %w", err)
	}
	// make a lookup map of the relevant tx hashes which need receipts
	relevantTx := make(map[gethcommon.Hash]bool)
	for _, l := range logs {
		relevantTx[l.TxHash] = true
	}

	for idx, transaction := range block.Transactions() {
		if !relevantTx[transaction.Hash()] && !r.isObscuroTransaction(transaction) {
			// put in a dummy receipt so that the index matches the transaction index
			// (the receipts list maintains the indexes of the transactions, it is a sparse list)
			receipts[idx] = &types.Receipt{Status: types.ReceiptStatusFailed}
			continue
		}
		receipt, err := r.ethClient.TransactionReceipt(transaction.Hash())

		if err != nil || receipt == nil {
			r.logger.Error("Problem with retrieving the receipt on the host!", log.ErrKey, err, log.CmpKey, log.CrossChainCmp)
			continue
		}

		r.logger.Trace("Adding receipt", "status", receipt.Status, log.TxKey, transaction.Hash(),
			log.BlockHashKey, blkHash, log.CmpKey, log.CrossChainCmp)

		receipts[idx] = receipt
	}

	return receipts, nil
}

// ExtractTenTransactions does all the filtering of txs to find all the transaction types we care about on the L2. These
// are pulled from the data in the L1 blocks and then submitted to the enclave for processing
func (r *Repository) ExtractTenTransactions(block *common.L1Block) (*common.ProcessedL1Data, error) {
	processed := &common.ProcessedL1Data{
		BlockHeader: block.Header(),
		Events:      []common.L1Event{},
	}

	// Get all logs for our contracts
	blkHash := block.Hash()
	var allAddresses []gethcommon.Address
	allAddresses = append(allAddresses, r.contractAddresses[MgmtContract]...)
	allAddresses = append(allAddresses, r.contractAddresses[MsgBus]...)

	logs, err := r.ethClient.GetLogs(ethereum.FilterQuery{BlockHash: &blkHash, Addresses: allAddresses})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch logs for L1 block - %w", err)
	}

	// FIXME clean this monster function up & look for redundant block.transactions() loops on enclave side
	// group logs by transaction hash and topic
	type logGroup struct {
		crossChainLogs     []types.Log
		valueTransferLogs  []types.Log
		sequencerLogs      []types.Log
		secretRequestLogs  []types.Log
		secretResponseLogs []types.Log
		rollupAddedLogs    []types.Log
	}

	logsByTx := make(map[gethcommon.Hash]*logGroup)

	// filter logs by topic
	for _, l := range logs {
		if _, exists := logsByTx[l.TxHash]; !exists {
			logsByTx[l.TxHash] = &logGroup{}
		}

		switch l.Topics[0] {
		case crosschain.CrossChainEventID:
			logsByTx[l.TxHash].crossChainLogs = append(logsByTx[l.TxHash].crossChainLogs, l)
		case crosschain.ValueTransferEventID:
			logsByTx[l.TxHash].valueTransferLogs = append(logsByTx[l.TxHash].valueTransferLogs, l)
		case crosschain.SequencerEnclaveGrantedEventID:
			logsByTx[l.TxHash].sequencerLogs = append(logsByTx[l.TxHash].sequencerLogs, l)
		case crosschain.NetworkSecretRequestedID:
			logsByTx[l.TxHash].secretRequestLogs = append(logsByTx[l.TxHash].secretRequestLogs, l)
		case crosschain.NetworkSecretRespondedID:
			logsByTx[l.TxHash].secretResponseLogs = append(logsByTx[l.TxHash].secretResponseLogs, l)
		case crosschain.RollupAddedID:
			logsByTx[l.TxHash].rollupAddedLogs = append(logsByTx[l.TxHash].rollupAddedLogs, l)
		}
	}

	// Process each transaction once
	for txHash, txLogs := range logsByTx {
		tx, _, err := r.ethClient.TransactionByHash(txHash)
		if err != nil {
			r.logger.Error("Error fetching transaction by hash", txHash, err)
			continue
		}

		receipt, err := r.ethClient.TransactionReceipt(txHash)
		if err != nil {
			r.logger.Error("Error fetching transaction receipt", txHash, err)
			continue
		}

		txData := &common.L1TxData{
			Transaction:        tx,
			Receipt:            receipt,
			CrossChainMessages: &common.CrossChainMessages{},
			ValueTransfers:     &common.ValueTransferEvents{},
		}

		if len(txLogs.crossChainLogs) > 0 {
			if messages, err := crosschain.ConvertLogsToMessages(txLogs.crossChainLogs, crosschain.CrossChainEventName, crosschain.MessageBusABI); err == nil {
				txData.CrossChainMessages = &messages
				processed.AddEvent(common.CrossChainMessageTx, txData)
			}
		}

		if len(txLogs.valueTransferLogs) > 0 {
			if transfers, err := crosschain.ConvertLogsToValueTransfers(txLogs.valueTransferLogs, crosschain.ValueTransferEventName, crosschain.MessageBusABI); err == nil {
				txData.ValueTransfers = &transfers
				processed.AddEvent(common.CrossChainValueTranserTx, txData)
			}
		}

		if len(txLogs.sequencerLogs) > 0 {
			if enclaveIDs, err := crosschain.ConvertLogsToSequencerEnclaves(txLogs.sequencerLogs, crosschain.SequencerEnclaveGrantedEventName, crosschain.MgmtContractABI); err == nil {
				txData.SequencerEnclaveIDs = enclaveIDs
				processed.AddEvent(common.SequencerAddedTx, txData)
			}
		}

		if len(txLogs.secretRequestLogs) > 0 {
			processed.AddEvent(common.SecretRequestTx, txData)
		}

		if len(txLogs.secretResponseLogs) > 0 {
			processed.AddEvent(common.SecretResponseTx, txData)
		}

		if decodedTx := r.mgmtContractLib.DecodeTx(tx); decodedTx != nil {
			switch t := decodedTx.(type) {
			case *common.L1InitializeSecretTx:
				processed.AddEvent(common.InitialiseSecretTx, txData)
			case *common.L1SetImportantContractsTx:
				processed.AddEvent(common.SetImportantContractsTx, txData)
			case *common.L1RollupHashes:
				if blobs, err := r.blobResolver.FetchBlobs(context.Background(), block.Header(), t.BlobHashes); err == nil {
					txData.Blobs = blobs
					processed.AddEvent(common.RollupTx, txData)
				}
			}
		}
	}

	return processed, nil
}

// stream blocks from L1 as they arrive and forward them to subscribers, no guarantee of perfect ordering or that there won't be gaps.
// If streaming is interrupted it will carry on from latest, it won't try to replay missed blocks.
func (r *Repository) streamLiveBlocks() {
	liveStream, streamSub := r.resetLiveStream()
	for r.running.Load() {
		select {
		case header := <-liveStream:
			r.head = header.Hash()
			block, err := r.ethClient.BlockByHash(header.Hash())
			if err != nil {
				r.logger.Error("Error fetching new block", log.BlockHashKey, header.Hash(),
					log.BlockHeightKey, header.Number, log.ErrKey, err)
				continue
			}
			for _, handler := range r.blockSubscribers.Subscribers() {
				go handler.HandleBlock(block)
			}
		case <-time.After(_timeoutNoBlocks):
			r.logger.Warn("no new blocks received since timeout", "timeout", _timeoutNoBlocks)
			// reset stream to ensure it has not died
			liveStream, streamSub = r.resetLiveStream()
		}
	}

	if streamSub != nil {
		streamSub.Unsubscribe()
	}
}

func (r *Repository) resetLiveStream() (chan *types.Header, ethereum.Subscription) {
	err := retry.Do(func() error {
		if !r.running.Load() {
			// break out of the loop if repository has stopped
			return retry.FailFast(errors.New("repository is stopped"))
		}
		err := r.ethClient.ReconnectIfClosed()
		if err != nil {
			r.logger.Warn("failed to reconnect to L1", log.ErrKey, err)
			return err
		}
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{100 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))
	if err != nil {
		// this should only happen if repository has been stopped, because we retry reconnecting forever
		r.logger.Warn("unable to reconnect to L1", log.ErrKey, err)
		return nil, nil
	}
	return r.ethClient.BlockListener()
}

func (r *Repository) FetchBlockByHeight(height *big.Int) (*types.Block, error) {
	return r.ethClient.BlockByNumber(height)
}

// isObscuroTransaction will look at the 'to' address of the transaction, we are only interested in management contract and bridge transactions
func (r *Repository) isObscuroTransaction(transaction *types.Transaction) bool {
	var allAddresses []gethcommon.Address
	allAddresses = append(allAddresses, r.contractAddresses[MgmtContract]...)
	allAddresses = append(allAddresses, r.contractAddresses[MsgBus]...)
	for _, address := range allAddresses {
		if transaction.To() != nil && *transaction.To() == address {
			return true
		}
	}
	return false
}

//func (r *Repository) getRelevantTxReceiptsAndBlobs(block *common.L1Block) ([]*common.TxAndReceiptAndBlobs, error) {
//	// Create a slice that will only contain valid transactions
//	var txsWithReceipts []*common.TxAndReceiptAndBlobs
//
//	receipts, err := r.FetchObscuroReceipts(block)
//	if err != nil {
//		return nil, fmt.Errorf("failed to fetch receipts: %w", err)
//	}
//
//	for i, tx := range block.Transactions() {
//		// skip unsuccessful txs
//		if receipts[i].Status == types.ReceiptStatusFailed {
//			continue
//		}
//
//		txWithReceipt := &common.TxAndReceiptAndBlobs{
//			Tx:      tx,
//			Receipt: receipts[i],
//		}
//
//		if tx.Type() == types.BlobTxType {
//			txBlobs := tx.BlobHashes()
//			blobs, err := r.blobResolver.FetchBlobs(context.Background(), block.Header(), txBlobs)
//			if err != nil {
//				if errors.Is(err, ethereum.NotFound) {
//					r.logger.Crit("Blobs were not found on beacon chain or archive service", "block", block.Hash(), "error", err)
//				} else {
//					r.logger.Crit("could not fetch blobs", log.ErrKey, err)
//				}
//				continue
//			}
//			txWithReceipt.Blobs = blobs
//		}
//
//		// Append only valid transactions
//		txsWithReceipts = append(txsWithReceipts, txWithReceipt)
//	}
//
//	return txsWithReceipts, nil
//}

//func (r *Repository) getCrossChainMessages(receipt *types.Receipt) (common.CrossChainMessages, error) {
//	logsForReceipt, err := crosschain.FilterLogsFromReceipt(receipt, &r.contractAddresses[MsgBus][0], &crosschain.CrossChainEventID)
//	if err != nil {
//		r.logger.Error("Error encountered when filtering receipt logs for cross chain messages.", log.ErrKey, err)
//		return make(common.CrossChainMessages, 0), err
//	}
//	messages, err := crosschain.ConvertLogsToMessages(logsForReceipt, crosschain.CrossChainEventName, crosschain.MessageBusABI)
//	if err != nil {
//		r.logger.Error("Error encountered converting the extracted relevant logs to messages", log.ErrKey, err)
//		return make(common.CrossChainMessages, 0), err
//	}
//
//	return messages, nil
//}
//
//func (r *Repository) getValueTransferEvents(receipt *types.Receipt) (common.ValueTransferEvents, error) {
//	logsForReceipt, err := crosschain.FilterLogsFromReceipt(receipt, &r.contractAddresses[MsgBus][0], &crosschain.ValueTransferEventID)
//	if err != nil {
//		r.logger.Error("Error encountered when filtering receipt logs for value transfers.", log.ErrKey, err)
//		return make(common.ValueTransferEvents, 0), err
//	}
//	transfers, err := crosschain.ConvertLogsToValueTransfers(logsForReceipt, crosschain.CrossChainEventName, crosschain.MessageBusABI)
//	if err != nil {
//		r.logger.Error("Error encountered converting the extracted relevant logs to messages", log.ErrKey, err)
//		return make(common.ValueTransferEvents, 0), err
//	}
//
//	return transfers, nil
//}

//func (r *Repository) getSequencerEventLogs(receipt *types.Receipt) ([]types.Log, error) {
//	sequencerLogs, err := crosschain.FilterLogsFromReceipt(receipt, &r.contractAddresses[MgmtContract][0], &crosschain.SequencerEnclaveGrantedEventID)
//	if err != nil {
//		r.logger.Error("Error filtering sequencer logs", log.ErrKey, err)
//		return []types.Log{}, err
//	}
//
//	// TODO convert to add sequencer?
//
//	return sequencerLogs, nil
//}

func (r *Repository) getRequestSecretEventLogs(receipt *types.Receipt) ([]types.Log, error) {
	sequencerLogs, err := crosschain.FilterLogsFromReceipt(receipt, &r.contractAddresses[MgmtContract][0], &crosschain.NetworkSecretRequestedID)
	if err != nil {
		r.logger.Error("Error filtering sequencer logs", log.ErrKey, err)
		return []types.Log{}, err
	}
	return sequencerLogs, nil
}

func (r *Repository) getSecretResponseLogs(receipt *types.Receipt) ([]types.Log, error) {
	sequencerLogs, err := crosschain.FilterLogsFromReceipt(receipt, &r.contractAddresses[MgmtContract][0], &crosschain.NetworkSecretRespondedID)
	if err != nil {
		r.logger.Error("Error filtering sequencer logs", log.ErrKey, err)
		return []types.Log{}, err
	}
	return sequencerLogs, nil
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
