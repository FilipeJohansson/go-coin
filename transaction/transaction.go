package transaction

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
	Message string `json:"message,omitempty"`
}

func NewTransaction(from string, to string, amount string, msg ...string) *Transaction {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	return &Transaction{
		From:    from,
		To:      to,
		Amount:  amount,
		Message: message,
	}
}

func (t *Transaction) Print() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}
