package components

import (
	"fmt"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

type SharedSecretProcessor struct {
	mgmtContractLib     mgmtcontractlib.MgmtContractLib
	attestationProvider AttestationProvider // interface for producing attestation reports and verifying them
	storage             storage.Storage
	logger              gethlog.Logger
}

func NewSharedSecretProcessor(mgmtcontractlib mgmtcontractlib.MgmtContractLib, attestationProvider AttestationProvider, storage storage.Storage, logger gethlog.Logger) *SharedSecretProcessor {
	return &SharedSecretProcessor{
		mgmtContractLib:     mgmtcontractlib,
		attestationProvider: attestationProvider,
		storage:             storage,
		logger:              logger,
	}
}

// ProcessNetworkSecretMsgs we watch for all messages that are requesting or receiving the secret and we store the nodes attested keys
func (ssp *SharedSecretProcessor) ProcessNetworkSecretMsgs(br *common.BlockAndReceipts) []*common.ProducedSecretResponse {
	var responses []*common.ProducedSecretResponse
	transactions := br.SuccessfulTransactions()
	block := br.Block
	for _, tx := range *transactions {
		t := ssp.mgmtContractLib.DecodeTx(tx)

		// this transaction is for a node that has joined the network and needs to be sent the network secret
		if scrtReqTx, ok := t.(*ethadapter.L1RequestSecretTx); ok {
			ssp.logger.Info("Process shared secret request.", log.BlockHeightKey, block.Number(), log.BlockHashKey, block.Hash(), log.TxKey, tx.Hash())
			resp, err := ssp.processSecretRequest(scrtReqTx)
			if err != nil {
				ssp.logger.Error("Failed to process shared secret request.", log.ErrKey, err)
				continue
			}
			responses = append(responses, resp)
		}

		// this transaction was created by the genesis node, we need to store their attested key to decrypt their rollup
		if initSecretTx, ok := t.(*ethadapter.L1InitializeSecretTx); ok {
			// todo (#1580) - ensure that we don't accidentally skip over the real `L1InitializeSecretTx` message. Otherwise
			//  our node will never be able to speak to other nodes.
			// there must be a way to make sure that this transaction can only be sent once.
			att, err := common.DecodeAttestation(initSecretTx.Attestation)
			if err != nil {
				ssp.logger.Error("Could not decode attestation report", log.ErrKey, err)
			}

			err = ssp.storeAttestation(att)
			if err != nil {
				ssp.logger.Error("Could not store the attestation report.", log.ErrKey, err)
			}
		}
	}
	return responses
}

func (ssp *SharedSecretProcessor) processSecretRequest(req *ethadapter.L1RequestSecretTx) (*common.ProducedSecretResponse, error) {
	att, err := common.DecodeAttestation(req.Attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation - %w", err)
	}

	ssp.logger.Info("received attestation", "attestation", att)
	secret, err := ssp.verifyAttestationAndEncryptSecret(att)
	if err != nil {
		return nil, fmt.Errorf("secret request failed, no response will be published - %w", err)
	}

	// Store the attested key only if the attestation process succeeded.
	err = ssp.storeAttestation(att)
	if err != nil {
		return nil, fmt.Errorf("could not store attestation, no response will be published. Cause: %w", err)
	}

	ssp.logger.Trace("Processed secret request.", "owner", att.Owner)
	return &common.ProducedSecretResponse{
		Secret:      secret,
		RequesterID: att.Owner,
		HostAddress: att.HostAddress,
	}, nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (ssp *SharedSecretProcessor) verifyAttestationAndEncryptSecret(att *common.AttestationReport) (common.EncryptedSharedEnclaveSecret, error) {
	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := ssp.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, fmt.Errorf("unable to verify report - %w", err)
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = VerifyIdentity(data, att); err != nil {
		return nil, fmt.Errorf("unable to verify identity - %w", err)
	}
	ssp.logger.Info(fmt.Sprintf("Successfully verified attestation and identity. Owner: %s", att.Owner))

	secret, err := ssp.storage.FetchSecret()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve secret; this should not happen. Cause: %w", err)
	}
	return crypto.EncryptSecret(att.PubKey, *secret, ssp.logger)
}

// storeAttestation stores the attested keys of other nodes so we can decrypt their rollups
func (ssp *SharedSecretProcessor) storeAttestation(att *common.AttestationReport) error {
	ssp.logger.Info(fmt.Sprintf("Store attestation. Owner: %s", att.Owner))
	// Store the attestation
	key, err := gethcrypto.DecompressPubkey(att.PubKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key %w", err)
	}
	err = ssp.storage.StoreAttestedKey(att.Owner, key)
	if err != nil {
		return fmt.Errorf("could not store attested key. Cause: %w", err)
	}
	return nil
}
