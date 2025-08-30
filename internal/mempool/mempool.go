package mempool

import (
	"encoding/hex"
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/transaction"
)

type Mempool struct {
	PendingTransactions []*transaction.Transaction `json:"pendingTransactions"`
}

func NewMempool() *Mempool {
	return &Mempool{
		PendingTransactions: make([]*transaction.Transaction, 0),
	}
}

func (m *Mempool) AddTransaction(tx *transaction.Transaction) {
	m.PendingTransactions = append(m.PendingTransactions, tx)
}

func (m *Mempool) GetTransactions() []*transaction.Transaction {
	return m.PendingTransactions[:]
}

func (m *Mempool) CleanProcessedTransactions(processedTxs []*transaction.Transaction) {
	processedHashes := make(map[string]bool)
	for _, tx := range processedTxs {
		processedHashes[hex.EncodeToString(tx.GetHash())] = true
	}

	remaining := make([]*transaction.Transaction, 0)
	for _, tx := range m.PendingTransactions {
		if !processedHashes[hex.EncodeToString(tx.GetHash())] {
			remaining = append(remaining, tx)
		}
	}

	m.PendingTransactions = remaining
}

func (m *Mempool) Size() int {
	return len(m.PendingTransactions)
}

func (m *Mempool) Contains(tx *transaction.Transaction) bool {
	targetHash := hex.EncodeToString(tx.GetHash())
	for _, pendingTx := range m.PendingTransactions {
		if hex.EncodeToString(pendingTx.GetHash()) == targetHash {
			return true
		}
	}

	return false
}

func (m *Mempool) Print() string {
	var content string

	for _, tx := range m.PendingTransactions {
		content += fmt.Sprintf("===[ Pending Transaction ]===%s", tx.Print())
	}

	return content
}
