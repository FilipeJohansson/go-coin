package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/FilipeJohansson/go-coin/internal/transaction"
)

type Block struct {
	Timestamp     time.Time                  `json:"timestamp"`
	Transactions  []*transaction.Transaction `json:"transactions"`
	Message       string                     `json:"message"`
	PrevBlockHash string                     `json:"prevBlockHash"`
	BlockHash     string                     `json:"blockHash"`
	Nonce         int                        `json:"nonce"`
	Difficulty    int                        `json:"difficulty"`
}

func NewBlock(prevBlockHash string, msg ...string) *Block {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	return &Block{
		Timestamp:     time.Now(),
		PrevBlockHash: prevBlockHash,
		Message:       message,
	}
}

func (b *Block) AddTransaction(transaction *transaction.Transaction) {
	b.Transactions = append(b.Transactions, transaction)
}

func (b *Block) Mine(difficulty int) {
	b.Difficulty = difficulty
	target := strings.Repeat("0", difficulty)

	for {
		b.SaveBlockHash()
		if strings.HasPrefix(b.BlockHash, target) {
			break
		}
		b.Nonce++
	}
}

func (b *Block) SaveBlockHash() {
	b.BlockHash = b.GetHash()
}

func (b *Block) GetHash() string {
	data := fmt.Sprintf("%v%s%s%d",
		b.Timestamp.Unix(),
		b.FormatTransactions(),
		b.PrevBlockHash,
		b.Nonce)

	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to hexadecimal string
	return hex.EncodeToString(hashBytes)
}

func (b *Block) IsHashRight() bool {
	target := strings.Repeat("0", b.Difficulty)
	if !strings.HasPrefix(b.BlockHash, target) {
		return false
	}

	if b.GetHash() != b.BlockHash {
		return false
	}

	return true
}

func (b *Block) FormatTransactions() string {
	var formattedTransactions string
	for _, t := range b.Transactions {
		formattedTransactions += t.Json()
	}

	return formattedTransactions
}

func (b *Block) Json() string {
	json, err := json.Marshal(b)
	if err != nil {
		// error
		return ""
	}

	return fmt.Sprintf("%s\n", json)
}
