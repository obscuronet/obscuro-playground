package ethereummock

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/host/l1"

	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/integration/datagenerator"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
)

var (
	// addresses used to simulate different methods on the mgmt contract
	depositTxAddr          = datagenerator.RandomAddress()
	rollupTxAddr           = datagenerator.RandomAddress()
	storeSecretTxAddr      = datagenerator.RandomAddress()
	requestSecretTxAddr    = datagenerator.RandomAddress()
	initializeSecretTxAddr = datagenerator.RandomAddress()
	grantSeqTxAddr         = datagenerator.RandomAddress()

	messageBusAddr = datagenerator.RandomAddress()

	// ContractAddresses maps contract types to their addresses
	ContractAddresses = map[l1.ContractType][]gethcommon.Address{
		l1.MgmtContract: {
			depositTxAddr,
			rollupTxAddr,
			storeSecretTxAddr,
			requestSecretTxAddr,
			initializeSecretTxAddr,
			grantSeqTxAddr,
		},
		l1.MsgBus: {
			messageBusAddr,
		},
	}
)

// mockContractLib is an implementation of the mgmtcontractlib.MgmtContractLib
// it creates ethereum mocked transactions from common.L1Transaction
// and converts ethereum mocked transactions to common.L1Transaction
type mockContractLib struct{}

func NewMgmtContractLibMock() mgmtcontractlib.MgmtContractLib {
	return &mockContractLib{}
}

func (m *mockContractLib) IsMock() bool {
	return true
}

func (m *mockContractLib) GetContractAddr() *gethcommon.Address {
	return &rollupTxAddr
}

func (m *mockContractLib) DecodeTx(tx *types.Transaction) common.L1TenTransaction {
	// Do not decode erc20 transactions, this is the responsibility
	// of the erc20 contract lib.
	if tx.To().Hex() == depositTxAddr.Hex() {
		return nil
	}

	if tx.To().Hex() == rollupTxAddr.Hex() {
		return &common.L1RollupHashes{
			BlobHashes: tx.BlobHashes(),
		}
	}
	return decodeTx(tx)
}

// TODO: Ziga - fix this mock implementation later if needed
func (m *mockContractLib) CreateBlobRollup(t *common.L1RollupTx) (types.TxData, error) {
	var err error
	blobs, err := ethadapter.EncodeBlobs(t.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to convert rollup to blobs: %w", err)
	}

	var blobHashes []gethcommon.Hash
	var sidecar *types.BlobTxSidecar
	if sidecar, blobHashes, err = ethadapter.MakeSidecar(blobs); err != nil {
		return nil, fmt.Errorf("failed to make sidecar: %w", err)
	}

	hashesTx := common.L1RollupHashes{BlobHashes: blobHashes}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(hashesTx); err != nil {
		panic(err)
	}

	return &types.BlobTx{
		To:         rollupTxAddr,
		Data:       buf.Bytes(),
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}, nil
}

func (m *mockContractLib) CreateRequestSecret(tx *common.L1RequestSecretTx) types.TxData {
	return encodeTx(tx, requestSecretTxAddr)
}

func (m *mockContractLib) CreateRespondSecret(tx *common.L1RespondSecretTx, _ bool) types.TxData {
	return encodeTx(tx, storeSecretTxAddr)
}

func (m *mockContractLib) CreateInitializeSecret(tx *common.L1InitializeSecretTx) types.TxData {
	return encodeTx(tx, initializeSecretTxAddr)
}

func (m *mockContractLib) GetHostAddressesMsg() (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeHostAddressesResponse([]byte) ([]string, error) {
	return []string{""}, nil
}

func (m *mockContractLib) GetImportantContractKeysMsg() (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeImportantContractKeysResponse([]byte) ([]string, error) {
	return []string{""}, nil
}

func (m *mockContractLib) SetImportantContractMsg(string, gethcommon.Address) (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) GetImportantAddressCallMsg(string) (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeImportantAddressResponse([]byte) (gethcommon.Address, error) {
	return gethcommon.Address{}, nil
}

func MockGrantSeqTxAddress() gethcommon.Address {
	return grantSeqTxAddr
}

func decodeTx(tx *types.Transaction) common.L1TenTransaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify the L1 operation (rollup/storesecret/requestsecret)
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t common.L1TenTransaction
	switch tx.To().Hex() {
	case storeSecretTxAddr.Hex():
		t = &common.L1RespondSecretTx{}
	case depositTxAddr.Hex():
		t = &common.L1DepositTx{}
	case requestSecretTxAddr.Hex():
		t = &common.L1RequestSecretTx{}
	case initializeSecretTxAddr.Hex():
		t = &common.L1InitializeSecretTx{}
	case grantSeqTxAddr.Hex():
		// this tx is empty and entirely mocked, no need to decode
		return &common.L1PermissionSeqTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}

	return t
}

func encodeTx(tx common.L1TenTransaction, opType gethcommon.Address) types.TxData {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	// the mock implementation does not process contract calls
	// this uses the To address to distinguish between different contract calls / different l1 transactions
	return &types.LegacyTx{
		Data: buf.Bytes(),
		To:   &opType,
	}
}
