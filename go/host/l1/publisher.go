package l1

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/pkg/errors"
)

type Publisher struct {
	hostData        host.Identity
	hostWallet      wallet.Wallet // Wallet used to issue ethereum transactions
	ethClient       ethadapter.EthClient
	mgmtContractLib mgmtcontractlib.MgmtContractLib // Library to handle Management Contract lib operations

	repository host.L1BlockRepository
	logger     gethlog.Logger

	hostStopper *stopcontrol.StopControl

	maxWaitForL1Receipt       time.Duration
	retryIntervalForL1Receipt time.Duration
}

func NewL1Publisher(
	hostData host.Identity,
	hostWallet wallet.Wallet,
	client ethadapter.EthClient,
	mgmtContract mgmtcontractlib.MgmtContractLib,
	repository host.L1BlockRepository,
	hostStopper *stopcontrol.StopControl,
	logger gethlog.Logger,
	maxWaitForL1Receipt time.Duration,
	retryIntervalForL1Receipt time.Duration,
) *Publisher {
	return &Publisher{
		hostData:                  hostData,
		hostWallet:                hostWallet,
		ethClient:                 client,
		mgmtContractLib:           mgmtContract,
		repository:                repository,
		hostStopper:               hostStopper,
		logger:                    logger,
		maxWaitForL1Receipt:       maxWaitForL1Receipt,
		retryIntervalForL1Receipt: retryIntervalForL1Receipt,
	}
}

func (p *Publisher) Start() error {
	return nil
}

func (p *Publisher) Stop() error {
	return nil
}

func (p *Publisher) HealthStatus() host.HealthStatus {
	// todo (@matt) do proper health status based on failed transactions or something
	errMsg := ""
	if p.hostStopper.IsStopping() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

func (p *Publisher) InitializeSecret(attestation *common.AttestationReport, encSecret common.EncryptedSharedEnclaveSecret) error {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &ethadapter.L1InitializeSecretTx{
		AggregatorID:  &p.hostData.ID,
		Attestation:   encodedAttestation,
		InitialSecret: encSecret,
		HostAddress:   p.hostData.P2PPublicAddress,
	}
	initialiseSecretTx := p.mgmtContractLib.CreateInitializeSecret(l1tx)
	// we block here until we confirm a successful receipt. It is important this is published before the initial rollup.
	return p.publishTransaction(initialiseSecretTx)
}

func (p *Publisher) RequestSecret(attestation *common.AttestationReport) (gethcommon.Hash, error) {
	encodedAttestation, err := common.EncodeAttestation(attestation)
	if err != nil {
		return gethcommon.Hash{}, errors.Wrap(err, "could not encode attestation")
	}
	l1tx := &ethadapter.L1RequestSecretTx{
		Attestation: encodedAttestation,
	}
	// record the L1 head height before we submit the secret request, so we know which block to watch from
	l1Head, err := p.ethClient.FetchHeadBlock()
	if err != nil {
		err = p.ethClient.ReconnectIfClosed()
		if err != nil {
			panic(errors.Wrap(err, "could not reconnect to eth client"))
		}
		l1Head, err = p.ethClient.FetchHeadBlock()
		if err != nil {
			panic(errors.Wrap(err, "could not fetch head block"))
		}
	}
	requestSecretTx := p.mgmtContractLib.CreateRequestSecret(l1tx)
	// we wait until the secret req transaction has succeeded before we start polling for the secret
	err = p.publishTransaction(requestSecretTx)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	return l1Head.Hash(), nil
}

func (p *Publisher) PublishSecretResponse(secretResponse *common.ProducedSecretResponse) error {
	l1tx := &ethadapter.L1RespondSecretTx{
		Secret:      secretResponse.Secret,
		RequesterID: secretResponse.RequesterID,
		AttesterID:  p.hostData.ID,
		HostAddress: secretResponse.HostAddress,
	}
	// todo (#1624) - l1tx.Sign(a.attestationPubKey) doesn't matter as the waitSecret will process a tx that was reverted
	respondSecretTx := p.mgmtContractLib.CreateRespondSecret(l1tx, false)
	p.logger.Info("Broadcasting secret response L1 tx.", "requester", secretResponse.RequesterID)

	// fire-and-forget (track the receipt asynchronously)
	go func() {
		err := p.publishTransaction(respondSecretTx)
		if err != nil {
			p.logger.Error("could not broadcast secret response L1 tx", log.ErrKey, err)
		}
	}()

	return nil
}

func (p *Publisher) ExtractSecretResponses(block *types.Block) []*ethadapter.L1RespondSecretTx {
	var secretRespTxs []*ethadapter.L1RespondSecretTx
	for _, tx := range block.Transactions() {
		t := p.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if scrtTx, ok := t.(*ethadapter.L1RespondSecretTx); ok {
			secretRespTxs = append(secretRespTxs, scrtTx)
		}
	}
	return secretRespTxs
}

func (p *Publisher) ExtractRollupTxs(block *types.Block) []*ethadapter.L1RollupTx {
	var rollupTxs []*ethadapter.L1RollupTx
	for _, tx := range block.Transactions() {
		t := p.mgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}
		if rollupTx, ok := t.(*ethadapter.L1RollupTx); ok {
			rollupTxs = append(rollupTxs, rollupTx)
		}
	}
	return rollupTxs
}

func (p *Publisher) FetchLatestSeqNo() (*big.Int, error) {
	return p.ethClient.FetchLastBatchSeqNo(*p.mgmtContractLib.GetContractAddr())
}

func (p *Publisher) PublishRollup(producedRollup *common.ExtRollup) {
	encRollup, err := common.EncodeRollup(producedRollup)
	if err != nil {
		p.logger.Crit("could not encode rollup.", log.ErrKey, err)
	}
	tx := &ethadapter.L1RollupTx{
		Rollup: encRollup,
	}
	p.logger.Info("Publishing rollup", "size", len(encRollup)/1024, log.RollupHashKey, producedRollup.Hash())

	p.logger.Trace("Sending transaction to publish rollup", "rollup_header",
		gethlog.Lazy{Fn: func() string {
			header, err := json.MarshalIndent(producedRollup.Header, "", "   ")
			if err != nil {
				return err.Error()
			}

			return string(header)
		}}, log.RollupHashKey, producedRollup.Header.Hash(), "batches_len", len(producedRollup.BatchPayloads))

	rollupTx := p.mgmtContractLib.CreateRollup(tx)

	err = p.publishTransaction(rollupTx)
	if err != nil {
		p.logger.Error("could not issue rollup tx", log.ErrKey, err)
	} else {
		p.logger.Info("Rollup included in L1", log.RollupHashKey, producedRollup.Hash())
	}
}

func (p *Publisher) FetchLatestPeersList() ([]string, error) {
	msg, err := p.mgmtContractLib.GetHostAddresses()
	if err != nil {
		return nil, err
	}
	response, err := p.ethClient.CallContract(msg)
	if err != nil {
		return nil, err
	}
	decodedResponse, err := p.mgmtContractLib.DecodeCallResponse(response)
	if err != nil {
		return nil, err
	}
	hostAddresses := decodedResponse[0]

	// We remove any duplicate addresses and our own address from the retrieved peer list
	var filteredHostAddresses []string
	uniqueHostKeys := make(map[string]bool) // map to track addresses we've seen already
	for _, hostAddress := range hostAddresses {
		// We exclude our own address.
		if hostAddress == p.hostData.P2PPublicAddress {
			continue
		}
		if _, found := uniqueHostKeys[hostAddress]; !found {
			uniqueHostKeys[hostAddress] = true
			filteredHostAddresses = append(filteredHostAddresses, hostAddress)
		}
	}

	return filteredHostAddresses, nil
}

// publishTransaction will keep trying unless the L1 seems to be unavailable or the tx is otherwise rejected
// It is responsible for keeping the nonce accurate, according to the following rules:
// - Caller should not increment the wallet nonce before this method is called
// - This method will increment the wallet nonce only if the transaction is successfully broadcast
// - This method will continue to resend the tx using latest gas price until it is successfully broadcast or the L1 is unavailable/this service is shutdown
// - **ONLY** the L1 publisher service is publishing transactions for this wallet (to avoid nonce conflicts)
func (p *Publisher) publishTransaction(tx types.TxData) error {
	// the nonce to be used for this tx attempt
	nonce := p.hostWallet.GetNonceAndIncrement()

	// while the publisher service is still alive we keep trying to get the transaction into the L1
	for !p.hostStopper.IsStopping() {
		// make sure an earlier tx hasn't been abandoned
		if nonce > p.hostWallet.GetNonce() {
			return errors.New("earlier transaction has failed to complete, we need to abort this transaction")
		}
		// update the tx gas price before each attempt
		tx, err := p.ethClient.PrepareTransactionToSend(tx, p.hostWallet.Address(), nonce)
		if err != nil {
			p.hostWallet.SetNonce(nonce) // revert the wallet nonce because we failed to complete the transaction
			return errors.Wrap(err, "could not estimate gas/gas price for L1 tx")
		}

		signedTx, err := p.hostWallet.SignTransaction(tx)
		if err != nil {
			p.hostWallet.SetNonce(nonce) // revert the wallet nonce because we failed to complete the transaction
			return errors.Wrap(err, "could not sign L1 tx")
		}

		p.logger.Info("Host issuing l1 tx", log.TxKey, signedTx.Hash(), "size", signedTx.Size()/1024)
		err = p.ethClient.SendTransaction(signedTx)
		if err != nil {
			p.hostWallet.SetNonce(nonce) // revert the wallet nonce because we failed to complete the transaction
			return errors.Wrap(err, "could not broadcast L1 tx")
		}
		p.logger.Info("Successfully submitted tx to L1", "txHash", signedTx.Hash())

		var receipt *types.Receipt
		// retry until receipt is found
		err = retry.Do(
			func() error {
				receipt, err = p.ethClient.TransactionReceipt(signedTx.Hash())
				if err != nil {
					return fmt.Errorf("could not get receipt for L1 tx=%s: %w", signedTx.Hash(), err)
				}
				return err
			},
			retry.NewTimeoutStrategy(p.maxWaitForL1Receipt, p.retryIntervalForL1Receipt),
		)
		if err != nil {
			p.logger.Info("Receipt not found for transaction, we will re-attempt", log.ErrKey, err)
			continue // try again on the same nonce, with updated gas price
		}

		if err == nil && receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("unsuccessful receipt found for published L1 transaction, status=%d", receipt.Status)
		}

		p.logger.Debug("L1 transaction successful receipt found.", log.TxKey, signedTx.Hash(),
			log.BlockHeightKey, receipt.BlockNumber, log.BlockHashKey, receipt.BlockHash)
		break
	}
	return nil
}
