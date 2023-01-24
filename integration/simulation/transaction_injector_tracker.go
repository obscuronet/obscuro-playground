package simulation

import (
	"sync"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/go/common"
)

type txInjectorTracker struct {
	l1TransactionsLock                sync.RWMutex
	L1Transactions                    []ethadapter.L1Transaction
	l2TransactionsLock                sync.RWMutex
	TransferL2Transactions            []*common.L2Tx
	NativeValueTransferL2Transactions []*common.L2Tx
	WithdrawalL2Transactions          []*common.L2Tx
}

func newCounter() *txInjectorTracker {
	return &txInjectorTracker{
		l1TransactionsLock:       sync.RWMutex{},
		L1Transactions:           []ethadapter.L1Transaction{},
		l2TransactionsLock:       sync.RWMutex{},
		TransferL2Transactions:   []*common.L2Tx{},
		WithdrawalL2Transactions: []*common.L2Tx{},
	}
}

// trackL1Tx adds an L1Tx to the internal list
func (m *txInjectorTracker) trackL1Tx(tx ethadapter.L1Transaction) {
	m.l1TransactionsLock.Lock()
	defer m.l1TransactionsLock.Unlock()
	m.L1Transactions = append(m.L1Transactions, tx)
}

func (m *txInjectorTracker) trackTransferL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.TransferL2Transactions = append(m.TransferL2Transactions, tx)
}

func (m *txInjectorTracker) trackNativeValueTransferL2Tx(tx *common.L2Tx) {
	m.l2TransactionsLock.Lock()
	defer m.l2TransactionsLock.Unlock()
	m.NativeValueTransferL2Transactions = append(m.TransferL2Transactions, tx)
}

// GetL1Transactions returns all generated L1 L2Txs
func (m *txInjectorTracker) GetL1Transactions() []ethadapter.L1Transaction {
	return m.L1Transactions
}

// GetL2Transactions returns all generated transactions, excluding prefund and ERC20 deploy transactions.
func (m *txInjectorTracker) GetL2Transactions() ([]*common.L2Tx, []*common.L2Tx, []*common.L2Tx) {
	return m.TransferL2Transactions, m.WithdrawalL2Transactions, m.NativeValueTransferL2Transactions
}
