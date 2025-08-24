package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	From      string          `json:"from"`
	To        string          `json:"to"`
	Amount    string          `json:"amount"`
	Signature string          `json:"signature"`
	PublicKey ecdsa.PublicKey `json:"publicKey"`
	Message   string          `json:"message,omitempty"`
}

func NewTransaction(from string, to string, amount string, publicKey ecdsa.PublicKey, msg ...string) *Transaction {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	return &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Message:   message,
		PublicKey: publicKey,
	}
}

func (t *Transaction) GetHash() []byte {
	data := fmt.Sprintf("%s%s%s", t.From, t.To, t.Amount)

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
