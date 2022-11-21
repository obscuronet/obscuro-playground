package enclave

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/params"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/enclave/bridge"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/enclave/events"
	"github.com/obscuronet/go-obscuro/go/enclave/mempool"
	"github.com/obscuronet/go-obscuro/go/enclave/rollupchain"
	"github.com/obscuronet/go-obscuro/go/enclave/rpc"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
)

type enclaveImpl struct {
	config               config.EnclaveConfig
	storage              db.Storage
	blockResolver        db.BlockResolver
	mempool              mempool.Manager
	l1Blockchain         *gethcore.BlockChain
	rpcEncryptionManager rpc.EncryptionManager
	bridge               *bridge.Bridge
	subscriptionManager  *events.SubscriptionManager

	chain *rollupchain.RollupChain

	txCh   chan *common.L2Tx
	exitCh chan bool

	// Todo - disabled temporarily until TN1 is released
	// speculativeWorkInCh  chan bool
	// speculativeWorkOutCh chan speculativeWork

	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	erc20ContractLib    erc20contractlib.ERC20ContractLib
	attestationProvider AttestationProvider // interface for producing attestation reports and verifying them

	enclaveKey    *ecdsa.PrivateKey // this is a key specific to this enclave, which is included in the Attestation. Used for signing rollups and for encryption of the shared secret.
	enclavePubKey []byte            // the public key of the above

	transactionBlobCrypto crypto.TransactionBlobCrypto
	profiler              *profiler.Profiler
	logger                gethlog.Logger
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(config config.EnclaveConfig, mgmtContractLib mgmtcontractlib.MgmtContractLib, erc20ContractLib erc20contractlib.ERC20ContractLib, logger gethlog.Logger) common.Enclave {
	if len(config.ERC20ContractAddresses) < 2 {
		logger.Crit("failed to initialise enclave. At least two ERC20 contract addresses are required - the HOC " +
			"ERC20 address and the POC ERC20 address")
	}

	// todo - add the delay: N hashes

	var prof *profiler.Profiler
	// don't run a profiler on an attested enclave
	if !config.WillAttest && config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultEnclavePort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler", log.ErrKey, err)
		}
	}

	// Initialise the database
	backingDB, err := db.CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	chainConfig := params.ChainConfig{
		ChainID:             big.NewInt(config.ObscuroChainID),
		HomesteadBlock:      gethcommon.Big0,
		DAOForkBlock:        gethcommon.Big0,
		EIP150Block:         gethcommon.Big0,
		EIP155Block:         gethcommon.Big0,
		EIP158Block:         gethcommon.Big0,
		ByzantiumBlock:      gethcommon.Big0,
		ConstantinopleBlock: gethcommon.Big0,
		PetersburgBlock:     gethcommon.Big0,
		IstanbulBlock:       gethcommon.Big0,
		MuirGlacierBlock:    gethcommon.Big0,
		BerlinBlock:         gethcommon.Big0,
		LondonBlock:         gethcommon.Big0,
	}
	storage := db.NewStorage(backingDB, &chainConfig, logger)

	// Initialise the Ethereum "Blockchain" structure that will allow us to validate incoming blocks
	// Todo - check the minimum difficulty parameter
	var l1Blockchain *gethcore.BlockChain
	if config.ValidateL1Blocks {
		if config.GenesisJSON == nil {
			logger.Crit("enclave is configured to validate blocks, but genesis JSON is nil")
		}
		l1Blockchain = rollupchain.NewL1Blockchain(config.GenesisJSON, logger)
	} else {
		logger.Info("validateBlocks is set to false. L1 blocks will not be validated.")
	}

	// Todo- make sure the enclave cannot be started in production with WillAttest=false
	var attestationProvider AttestationProvider
	if config.WillAttest {
		attestationProvider = &EgoAttestationProvider{}
	} else {
		logger.Info("WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &DummyAttestationProvider{}
	}

	// todo - this has to be read from the database when the node restarts.
	// first time the node starts we derive the obscuro key from the master seed received after the shared secret exchange
	logger.Info("Generating the Obscuro key")

	// todo - save this to the db
	enclaveKey, err := gethcrypto.GenerateKey()
	if err != nil {
		logger.Crit("Failed to generate enclave key.", log.ErrKey, err)
	}
	serializedEnclavePubKey := gethcrypto.CompressPubkey(&enclaveKey.PublicKey)
	logger.Info(fmt.Sprintf("Generated public key %s", gethcommon.Bytes2Hex(serializedEnclavePubKey)))

	obscuroKey := crypto.GetObscuroKey(logger)
	rpcEncryptionManager := rpc.NewEncryptionManager(ecies.ImportECDSA(obscuroKey))

	transactionBlobCrypto := crypto.NewTransactionBlobCryptoImpl(logger)

	obscuroBridge := bridge.New(
		config.ERC20ContractAddresses[0],
		config.ERC20ContractAddresses[1],
		mgmtContractLib,
		erc20ContractLib,
		transactionBlobCrypto,
		config.ObscuroChainID,
		config.L1ChainID,
		logger,
	)
	memp := mempool.New(config.ObscuroChainID)

	subscriptionManager := events.NewSubscriptionManager(&rpcEncryptionManager, storage, logger)
	chain := rollupchain.New(config.HostID, config.NodeType, storage, l1Blockchain, obscuroBridge, subscriptionManager, transactionBlobCrypto, memp, rpcEncryptionManager, enclaveKey, config.L1ChainID, &chainConfig, logger)

	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Enclave service created with following config", log.CfgKey, string(jsonConfig))
	return &enclaveImpl{
		config:                config,
		storage:               storage,
		blockResolver:         storage,
		mempool:               memp,
		l1Blockchain:          l1Blockchain,
		rpcEncryptionManager:  rpcEncryptionManager,
		bridge:                obscuroBridge,
		subscriptionManager:   subscriptionManager,
		chain:                 chain,
		txCh:                  make(chan *common.L2Tx),
		exitCh:                make(chan bool),
		mgmtContractLib:       mgmtContractLib,
		erc20ContractLib:      erc20ContractLib,
		attestationProvider:   attestationProvider,
		enclaveKey:            enclaveKey,
		enclavePubKey:         serializedEnclavePubKey,
		transactionBlobCrypto: transactionBlobCrypto,
		profiler:              prof,
		logger:                logger,
	}
}

// Status is only implemented by the RPC wrapper
func (e *enclaveImpl) Status() (common.Status, error) {
	_, found := e.storage.FetchSecret()
	if !found {
		return common.AwaitingSecret, nil
	}
	return common.Running, nil // The enclave is local so it is always ready
}

// StopClient is only implemented by the RPC wrapper
func (e *enclaveImpl) StopClient() error {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) error {
	// todo - reinstate after TN1
	/*	if e.config.SpeculativeExecution {
			//start the speculative rollup execution loop on its own go routine
			go e.start(block)
		}
	*/
	return nil
}

func (e *enclaveImpl) ProduceGenesis(blkHash gethcommon.Hash) (*common.BlockSubmissionResponse, error) {
	rolGenesis, b := e.chain.ProduceGenesis(blkHash)
	return &common.BlockSubmissionResponse{
		ProducedRollup: rolGenesis.ToExtRollup(e.transactionBlobCrypto),
		BlockHeader:    b.Header(),
	}, nil
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitBlock(block types.Block, isLatest bool) (*common.BlockSubmissionResponse, error) {
	bsr, err := e.chain.SubmitBlock(block, isLatest)
	if err != nil {
		e.logger.Trace("SubmitBlock failed", "blk", block.Number(), "blkHash", block.Hash(), "err", err)
		return nil, err
	}
	e.logger.Trace("SubmitBlock successful", "blk", block.Number(), "blkHash", block.Hash())

	if bsr.IngestedRollupHeader != nil {
		hr, f := e.storage.FetchRollup(bsr.IngestedRollupHeader.Hash())
		if !f {
			e.logger.Crit("This should not happen because this rollup was just processed.")
		}
		e.mempool.RemoveMempoolTxs(hr, e.storage)
	}

	bsr.ProducedSecretResponses = e.processNetworkSecretMsgs(block)

	return bsr, nil
}

func (e *enclaveImpl) SubmitTx(tx common.EncryptedTx) (common.EncryptedResponseSendRawTx, error) {
	encodedTx, err := e.rpcEncryptionManager.DecryptBytes(tx)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_sendRawTransaction request. Cause: %w", err)
	}

	decryptedTx, err := rpc.ExtractTx(encodedTx)
	if err != nil {
		e.logger.Info("could not decrypt transaction. ", log.ErrKey, err)
		return nil, fmt.Errorf("could not decrypt transaction. Cause: %w", err)
	}

	err = e.checkGas(decryptedTx)
	if err != nil {
		e.logger.Info("", log.ErrKey, err.Error())
		return nil, err
	}

	err = e.mempool.AddMempoolTx(decryptedTx)
	if err != nil {
		return nil, err
	}

	if e.config.SpeculativeExecution {
		e.txCh <- decryptedTx
	}

	viewingKeyAddress, err := rpc.GetSender(decryptedTx)
	if err != nil {
		return nil, fmt.Errorf("could not recover viewing key address to encrypt eth_sendRawTransaction response. Cause: %w", err)
	}

	txHashBytes := []byte(decryptedTx.Hash().Hex())
	encryptedResult, err := e.rpcEncryptionManager.EncryptWithViewingKey(viewingKeyAddress, txHashBytes)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_sendRawTransaction request. Cause: %w", err)
	}

	return encryptedResult, nil
}

func (e *enclaveImpl) ExecuteOffChainTransaction(encryptedParams common.EncryptedParamsCall) (common.EncryptedResponseCall, error) {
	resp, err := e.chain.ExecuteOffChainTransaction(encryptedParams)
	if err != nil {
		e.logger.Info("Could not execute off chain call.", log.ErrKey, err)
	}
	return resp, err
}

func (e *enclaveImpl) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (common.EncryptedResponseGetTxCount, error) {
	var nonce uint64
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, err
	}

	address, err := rpc.ExtractAddress(paramBytes)
	if err != nil {
		return nil, err
	}
	hs := e.storage.FetchHeadState()
	if hs != nil {
		// todo: we should return an error when head state is not available, but for current test situations with race
		// 		conditions we allow it to return zero while head state is uninitialized
		s := e.storage.CreateStateDB(hs.HeadRollup)
		nonce = s.GetNonce(address)
	}

	encCount, err := e.rpcEncryptionManager.EncryptWithViewingKey(address, []byte(hexutil.EncodeUint64(nonce)))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getTransactionCount request. Cause: %w", err)
	}
	return encCount, nil
}

func (e *enclaveImpl) GetTransaction(encryptedParams common.EncryptedParamsGetTxByHash) (common.EncryptedResponseGetTxByHash, error) {
	hashBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted RPC request params. Cause: %w", err)
	}
	var paramList []string
	err = json.Unmarshal(hashBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPC request params from JSON. Cause: %w", err)
	}
	if len(paramList) == 0 {
		return nil, fmt.Errorf("required at least one param, but received zero")
	}
	txHash := gethcommon.HexToHash(paramList[0])

	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := e.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, db.ErrTxNotFound) {
			return nil, nil
		}
		return nil, err
	}

	viewingKeyAddress, err := rpc.GetSender(tx)
	if err != nil {
		return nil, fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionByHash response. Cause: %w", err)
	}

	// Unlike in the Geth impl, we hardcode the use of a London signer.
	// TODO - Once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`.
	signer := types.NewLondonSigner(tx.ChainId())
	rpcTx := newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer)

	txBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, fmt.Errorf("could not marshal transaction to JSON. Cause: %w", err)
	}
	return e.rpcEncryptionManager.EncryptWithViewingKey(viewingKeyAddress, txBytes)
}

func (e *enclaveImpl) GetTransactionReceipt(encryptedParams common.EncryptedParamsGetTxReceipt) (common.EncryptedResponseGetTxReceipt, error) {
	// We decrypt the transaction bytes.
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in eth_getTransactionReceipt request. Cause: %w", err)
	}
	txHash, err := rpc.ExtractTxHash(paramBytes)
	if err != nil {
		return nil, err
	}

	// We retrieve the transaction.
	tx, txRollupHash, txRollupHeight, _, err := e.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, db.ErrTxNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Only return receipts for transactions included in the canonical chain.
	r, f := e.storage.FetchRollupByHeight(txRollupHeight)
	if !f {
		return nil, fmt.Errorf("transaction not included in the canonical chain")
	}
	if !bytes.Equal(r.Hash().Bytes(), txRollupHash.Bytes()) {
		return nil, fmt.Errorf("transaction not included in the canonical chain")
	}

	// We retrieve the transaction receipt.
	txReceipt, err := e.storage.GetTransactionReceipt(txHash)
	if err != nil {
		if errors.Is(err, db.ErrTxNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	// We retrieve the sender's address.
	sender, err := rpc.GetSender(tx)
	if err != nil {
		return nil, fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionReceipt response. Cause: %w", err)
	}

	// We filter out irrelevant logs.
	txReceipt.Logs = e.subscriptionManager.FilterLogs(txReceipt.Logs, txRollupHash, &sender, &filters.FilterCriteria{})

	// We marshal the receipt to JSON.
	txReceiptBytes, err := txReceipt.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal transaction receipt to JSON in eth_getTransactionReceipt request. Cause: %w", err)
	}

	// We encrypt the receipt.
	encryptedTxReceipt, err := e.rpcEncryptionManager.EncryptWithViewingKey(sender, txReceiptBytes)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getTransactionReceipt request. Cause: %w", err)
	}

	return encryptedTxReceipt, nil
}

func (e *enclaveImpl) GetRollup(rollupHash common.L2RootHash) (*common.ExtRollup, error) {
	rollup, found := e.storage.FetchRollup(rollupHash)
	if !found {
		return nil, nil //nolint:nilnil
	}
	extRollup := rollup.ToExtRollup(e.transactionBlobCrypto)
	return &extRollup, nil
}

func (e *enclaveImpl) Attestation() (*common.AttestationReport, error) {
	if e.enclavePubKey == nil {
		e.logger.Error("public key not initialized, we can't produce the attestation report")
		return nil, fmt.Errorf("public key not initialized, we can't produce the attestation report")
	}
	report, err := e.attestationProvider.GetReport(e.enclavePubKey, e.config.HostID, e.config.HostAddress)
	if err != nil {
		e.logger.Error("could not produce remote report")
		return nil, fmt.Errorf("could not produce remote report")
	}
	return report, nil
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() (common.EncryptedSharedEnclaveSecret, error) {
	secret := crypto.GenerateEntropy(e.logger)
	e.storage.StoreSecret(secret)
	encSec, err := crypto.EncryptSecret(e.enclavePubKey, secret, e.logger)
	if err != nil {
		e.logger.Error("failed to encrypt secret.", log.ErrKey, err)
		return nil, fmt.Errorf("failed to encrypt secret. Cause: %w", err)
	}
	return encSec, nil
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s common.EncryptedSharedEnclaveSecret) error {
	secret, err := crypto.DecryptSecret(s, e.enclaveKey)
	if err != nil {
		return err
	}
	e.storage.StoreSecret(*secret)
	e.logger.Trace(fmt.Sprintf("Secret decrypted and stored. Secret: %v", secret))
	return nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (e *enclaveImpl) verifyAttestationAndEncryptSecret(att *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := e.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, fmt.Errorf("unable to verify report - %w", err)
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = VerifyIdentity(data, att); err != nil {
		return nil, fmt.Errorf("unable to verify identity - %w", err)
	}
	e.logger.Info(fmt.Sprintf("Successfully verified attestation and identity. Owner: %s", att.Owner))

	secret, found := e.storage.FetchSecret()
	if !found {
		return nil, errors.New("secret was nil, no secret to share - this shouldn't happen")
	}
	return crypto.EncryptSecret(att.PubKey, *secret, e.logger)
}

func (e *enclaveImpl) AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error {
	return e.rpcEncryptionManager.AddViewingKey(encryptedViewingKeyBytes, signature)
}

// storeAttestation stores the attested keys of other nodes so we can decrypt their rollups
func (e *enclaveImpl) storeAttestation(att *common.AttestationReport) error {
	e.logger.Info(fmt.Sprintf("Store attestation. Owner: %s", att.Owner))
	// Store the attestation
	key, err := gethcrypto.DecompressPubkey(att.PubKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key %w", err)
	}
	e.storage.StoreAttestedKey(att.Owner, key)
	return nil
}

func (e *enclaveImpl) GetBalance(encryptedParams common.EncryptedParamsGetBalance) (common.EncryptedResponseGetBalance, error) {
	return e.chain.GetBalance(encryptedParams)
}

func (e *enclaveImpl) GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, error) {
	return e.storage.CreateStateDB(*rollupHash).GetCode(address), nil
}

func (e *enclaveImpl) Subscribe(id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	return e.subscriptionManager.AddSubscription(id, encryptedSubscription)
}

func (e *enclaveImpl) Unsubscribe(id gethrpc.ID) error {
	e.subscriptionManager.RemoveSubscription(id)
	return nil
}

func (e *enclaveImpl) Stop() error {
	if e.config.SpeculativeExecution {
		e.exitCh <- true
	}

	if e.profiler != nil {
		return e.profiler.Stop()
	}

	return nil
}

// EstimateGas decrypts CallMsg data, runs the gas estimation for the data.
// Using the callMsg.From Viewing Key, returns the encrypted gas estimation
func (e *enclaveImpl) EstimateGas(encryptedParams common.EncryptedParamsEstimateGas) (common.EncryptedResponseEstimateGas, error) {
	// decrypt the input with the enclave PK
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt params in EstimateGas request. Cause: %w", err)
	}

	// extract params from byte slice to array of strings
	var paramList []interface{}
	err = json.Unmarshal(paramBytes, &paramList)
	if err != nil {
		return nil, fmt.Errorf("unable to decode EthCall params - %w", err)
	}

	// params are [callMsg, block number (optional) ]
	if len(paramList) < 1 {
		return nil, fmt.Errorf("required at least 1 params, but received %d", len(paramList))
	}

	callMsg, err := gethencoding.ExtractEthCall(paramList[0])
	if err != nil {
		return nil, fmt.Errorf("unable to decode EthCall Params - %w", err)
	}

	// TODO hook the correct blockNumber from the API call (paramList[1])
	gasEstimate, err := e.DoEstimateGas(callMsg, gethrpc.BlockNumber(0), e.chain.GlobalGasCap)
	if err != nil {
		return nil, fmt.Errorf("unable to estimate transaction - %w", err)
	}

	// encrypt the gas cost with the callMsg.From viewing key
	encryptedGasCost, err := e.rpcEncryptionManager.EncryptWithViewingKey(*callMsg.From, []byte(hexutil.EncodeUint64(uint64(gasEstimate))))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_estimateGas request. Cause: %w", err)
	}
	return encryptedGasCost, nil
}

func (e *enclaveImpl) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (common.EncryptedResponseGetLogs, error) {
	// We decrypt the params.
	paramBytes, err := e.rpcEncryptionManager.DecryptBytes(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt params in GetLogs request. Cause: %w", err)
	}

	// We extract the arguments from the param bytes.
	filter, forAddress, err := extractGetLogsParams(paramBytes)
	if err != nil {
		return nil, err
	}

	// We retrieve the relevant logs that match the filter.
	filteredLogs, err := e.subscriptionManager.GetFilteredLogs(forAddress, filter)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
	}

	// We encode and encrypt the logs with the viewing key for the requester's address.
	logBytes, err := json.Marshal(filteredLogs)
	if err != nil {
		return nil, fmt.Errorf("could not marshal logs to JSON. Cause: %w", err)
	}
	encryptedLogs, err := e.rpcEncryptionManager.EncryptWithViewingKey(*forAddress, logBytes)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to GetLogs request. Cause: %w", err)
	}
	return encryptedLogs, nil
}

// DoEstimateGas returns the estimation of minimum gas required to execute transaction
// This is a copy of https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/api.go#L1055
// there's a high complexity to the method due to geth business rules (which is mimic'd here)
// once the work of obscuro gas mechanics is established this method should be simplified
func (e *enclaveImpl) DoEstimateGas(args *gethapi.TransactionArgs, blockNr gethrpc.BlockNumber, gasCap uint64) (hexutil.Uint64, error) { //nolint: gocognit
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  = params.TxGas - 1
		hi  uint64
		cap uint64 //nolint:predeclared
	)
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(gethcommon.Address)
	}
	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= params.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// TODO review this with the gas mechanics/tokenomics work
		/*
			//Retrieve the block to act as the gas ceiling
			block, err := b.BlockByNumberOrHash(ctx, blockNrOrHash)
			if err != nil {
				return 0, err
			}
			if block == nil {
				return 0, errors.New("block not found")
			}
			hi = block.GasLimit()
		*/
		hi = e.chain.GlobalGasCap
	}
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = gethcommon.Big0
	}
	// Recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 { //nolint:nestif
		balance, err := e.chain.GetBalanceAtBlock(*args.From, blockNr)
		if err != nil {
			return 0, fmt.Errorf("unable to fetch account balance - %w", err)
		}

		available := new(big.Int).Set(balance.ToInt())
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, feeCap)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			e.logger.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer.ToInt(), "maxFeePerGas", feeCap, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if gasCap != 0 && hi > gasCap {
		e.logger.Warn("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap
	}
	cap = hi

	// Execute the binary search and hone in on an isGasEnough gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := e.isGasEnough(args, mid)
		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap { //nolint:nestif
		failed, result, err := e.isGasEnough(args, hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas { //nolint: errorlint
				if len(result.Revert()) > 0 {
					return 0, newRevertError(result)
				}
				return 0, result.Err
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hexutil.Uint64(hi), nil
}

// HealthCheck returns whether the enclave is deemed healthy
func (e *enclaveImpl) HealthCheck() (bool, error) {
	// check the storage health
	storageHealthy, err := e.storage.HealthCheck()
	if err != nil {
		// simplest iteration, log the error and just return that it's not healthy
		e.logger.Error("unable to HealthCheck enclave storage", "err", err)
		return false, nil
	}
	// TODO enclave healthcheck operations
	enclaveHealthy := true
	return storageHealthy && enclaveHealthy, nil
}

// Create a helper to check if a gas allowance results in an executable transaction
// isGasEnough returns whether the gaslimit should be raised, lowered, or if it was impossible to execute the message
func (e *enclaveImpl) isGasEnough(args *gethapi.TransactionArgs, gas uint64) (bool, *gethcore.ExecutionResult, error) {
	args.Gas = (*hexutil.Uint64)(&gas)
	result, err := e.chain.ExecuteOffChainTransactionAtBlock(args, gethrpc.BlockNumber(0))
	if err != nil {
		if errors.Is(err, gethcore.ErrIntrinsicGas) {
			return true, nil, nil // Special case, raise gas limit
		}
		return true, nil, err // Bail out
	}
	return result.Failed(), result, nil
}

func newRevertError(result *gethcore.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompasses an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

func (e *enclaveImpl) checkGas(tx *types.Transaction) error {
	txGasPrice := tx.GasPrice()
	if txGasPrice == nil {
		return fmt.Errorf("rejected transaction %s. No gas price was set", tx.Hash())
	}
	minGasPrice := e.config.MinGasPrice
	if txGasPrice.Cmp(minGasPrice) == -1 {
		return fmt.Errorf("rejected transaction %s. Gas price was only %d, wanted at least %d", tx.Hash(), txGasPrice, minGasPrice)
	}
	return nil
}

// processNetworkSecretMsgs we watch for all messages that are requesting or receiving the secret and we store the nodes attested keys
func (e *enclaveImpl) processNetworkSecretMsgs(block types.Block) []*common.ProducedSecretResponse {
	var responses []*common.ProducedSecretResponse
	for _, tx := range block.Transactions() {
		t := e.mgmtContractLib.DecodeTx(tx)

		// this transaction is for a node that has joined the network and needs to be sent the network secret
		if scrtReqTx, ok := t.(*ethadapter.L1RequestSecretTx); ok {
			e.logger.Info(fmt.Sprintf("Process shared secret request. Block: %d. TxKey: %d",
				block.NumberU64(), common.ShortHash(tx.Hash())))
			resp, err := e.processSecretRequest(scrtReqTx)
			if err != nil {
				e.logger.Error("Failed to process shared secret request.", log.ErrKey, err)
				continue
			}
			responses = append(responses, resp)
		}

		// this transaction was created by the genesis node, we need to store their attested key to decrypt their rollup
		if initSecretTx, ok := t.(*ethadapter.L1InitializeSecretTx); ok {
			// TODO - Ensure that we don't accidentally skip over the real `L1InitializeSecretTx` message. Otherwise
			//  our node will never be able to speak to other nodes.
			// there must be a way to make sure that this transaction can only be sent once.
			att, err := common.DecodeAttestation(initSecretTx.Attestation)
			if err != nil {
				e.logger.Error("Could not decode attestation report", log.ErrKey, err)
			}

			err = e.storeAttestation(att)
			if err != nil {
				e.logger.Error("Could not store the attestation report.", log.ErrKey, err)
			}
		}
	}
	return responses
}

func (e *enclaveImpl) processSecretRequest(req *ethadapter.L1RequestSecretTx) (*common.ProducedSecretResponse, error) {
	att, err := common.DecodeAttestation(req.Attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation - %w", err)
	}

	e.logger.Info("received attestation", "attestation", att)
	secret, err := e.verifyAttestationAndEncryptSecret(att)
	if err != nil {
		return nil, fmt.Errorf("secret request failed, no response will be published - %w", err)
	}

	// Store the attested key only if the attestation process succeeded.
	err = e.storeAttestation(att)
	if err != nil {
		return nil, fmt.Errorf("could not store attestation, no response will be published. Cause: %w", err)
	}

	e.logger.Trace("Processed secret request.", "owner", att.Owner)
	return &common.ProducedSecretResponse{
		Secret:      secret,
		RequesterID: att.Owner,
		HostAddress: att.HostAddress,
	}, nil
}

// Returns the params extracted from an eth_getLogs request.
func extractGetLogsParams(paramBytes []byte) (*filters.FilterCriteria, *gethcommon.Address, error) {
	// We verify the params.
	var paramsList []interface{}
	err := json.Unmarshal(paramBytes, &paramsList)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode GetLogs params. Cause: %w", err)
	}
	if len(paramsList) != 2 {
		return nil, nil, fmt.Errorf("expected 2 params in GetLogs request, but received %d", len(paramsList))
	}

	// We extract the first param, the filter for the logs.
	// We marshal the filter criteria from a map to JSON, then back from JSON into a FilterCriteria. This is
	// because the filter criteria arrives as a map, and there is no way to convert it to a map directly into a
	// FilterCriteria.
	filterJSON, err := json.Marshal(paramsList[0])
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}
	filter := filters.FilterCriteria{}
	err = filter.UnmarshalJSON(filterJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal filter criteria from JSON. Cause: %w", err)
	}

	// We extract the second param, the address the logs are for.
	forAddressHex, ok := paramsList[1].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected second argument in GetLogs request to be of type string, but got %T", paramsList[0])
	}
	forAddress := gethcommon.HexToAddress(forAddressHex)
	return &filter, &forAddress, nil
}

// Todo - reinstate speculative execution after TN1
/*
// internal structure to pass information.
type speculativeWork struct {
	found bool
	r     *obscurocore.Rollup
	s     *state.StateDB
	h     *nodecommon.Header
	txs   []*nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup              // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header               // the header of the new rollup
	processedTxs    []*nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]*nodecommon.L2Tx // structure used to prevent duplicates
	state           *state.StateDB                   // the state as calculated from the previous rollup and the processed transactions
}
*/
/*
func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]*nodecommon.L2Tx)}
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		env.headRollup, _ = e.storage.FetchRollup(blockState.HeadRollup)
		if env.headRollup != nil {
			env.state = e.storage.CreateStateDB(env.headRollup.Hash())
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:
			hash := winnerRollup.Hash()
			env.header = obscurocore.NewHeader(&hash, winnerRollup.Header.Number.Uint64()+1, e.config.HostID)
			env.headRollup = winnerRollup
			env.state = e.storage.CreateStateDB(winnerRollup.Hash())
			common.TraceWithID(e.nodeShortID, "Create new speculative env  r_%d(%d).",
				obscurocommon.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			))

			// determine the transactions that were not yet included
			env.processedTxs = currentTxs(winnerRollup, e.mempool.FetchMempoolTxs(), e.storage)
			env.processedTxsMap = obscurocore.MakeMap(env.processedTxs)

			// calculate the State after executing them
			evm.ExecuteTransactions(env.processedTxs, env.state, env.headRollup.Header, e.storage, e.config.ObscuroChainID, 0)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					evm.ExecuteTransactions([]*nodecommon.L2Tx{tx}, env.state, env.header, e.storage, e.config.ObscuroChainID, 0)
				}
			}

		case <-e.speculativeWorkInCh:
			if env.header == nil {
				e.speculativeWorkOutCh <- speculativeWork{found: false}
			} else {
				b := make([]*nodecommon.L2Tx, 0, len(env.processedTxs))
				b = append(b, env.processedTxs...)
				e.speculativeWorkOutCh <- speculativeWork{
					found: true,
					r:     env.headRollup,
					s:     env.state,
					h:     env.header,
					txs:   b,
				}
			}

		case <-e.exitCh:
			return
		}
	}
}
*/
