package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FilipeJohansson/go-coin/common"
)

type Transaction struct {
	From      string          `json:"from"`
	To        string          `json:"to"`
	Amount    uint64          `json:"amount"`
	Signature string          `json:"signature"`
	PublicKey ecdsa.PublicKey `json:"publicKey"`
	Message   string          `json:"message,omitempty"`
}

func NewTransaction(from string, to string, amount float64, publicKey ecdsa.PublicKey, msg ...string) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if to == "" {
		return nil, errors.New("recipient address cannot be empty")
	}

	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	return &Transaction{
		From:      from,
		To:        to,
		Amount:    uint64(amount * common.COINS_PER_UNIT),
		Message:   message,
		PublicKey: publicKey,
	}, nil
}

func (t *Transaction) GetHash() []byte {
	data := fmt.Sprintf("%s%s%d", t.From, t.To, t.Amount)

	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func (t *Transaction) Print() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}
