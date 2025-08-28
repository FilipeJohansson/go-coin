package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/utxo"
)

type TransactionInput struct {
	TransactionID string          `json:"transactionID"`
	OutputIndex   uint            `json:"outputIndex"`
	Signature     string          `json:"signature"`
	PublicKey     ecdsa.PublicKey `json:"publicKey"`
}

type TransactionOutput struct {
	Address string `json:"address"` // Recipient
	Amount  uint64 `json:"amount"`
}

type Transaction struct {
	Inputs  []TransactionInput  `json:"inputs"`
	Outputs []TransactionOutput `json:"outputs"`
	Fee     uint64              `json:"fee"`
	Message string              `json:"message,omitempty"`
}

func NewTransaction(senderAddress string, recipientAddress string, amount uint64, fee uint64, utxoSet *utxo.UTXOSet, senderPublicKey ecdsa.PublicKey, msg ...string) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if recipientAddress == "" {
		return nil, errors.New("recipient address cannot be empty")
	}

	if senderAddress == recipientAddress {
		return nil, errors.New("sender and recipient cannot be the same")
	}

	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	spendableUTXOs, err := utxoSet.FindSpendableUTXOsForAddress(senderAddress, amount+fee)
	if err != nil {
		return nil, err
	}

	var spendableUTXOsAmount uint64 = 0
	inputs := make([]TransactionInput, 0)
	for _, u := range spendableUTXOs {
		inputs = append(inputs, TransactionInput{
			TransactionID: u.TransactionID,
			OutputIndex:   u.OutputIndex,
			PublicKey:     senderPublicKey,
		})
		spendableUTXOsAmount += u.Amount
	}

	outputs := make([]TransactionOutput, 0)
	outputs = append(outputs, TransactionOutput{
		Address: recipientAddress,
		Amount:  amount,
	})

	if spendableUTXOsAmount > (amount + fee) {
		outputs = append(outputs, TransactionOutput{
			Address: senderAddress,
			Amount:  spendableUTXOsAmount - amount - fee,
		})
	}

	return &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Message: message,
	}, nil
}

func NewCoinbaseTransaction(recipientAddress string, amount uint64) *Transaction {
	return &Transaction{
		Inputs: []TransactionInput{},
		Outputs: []TransactionOutput{
			{
				Address: recipientAddress,
				Amount:  amount,
			},
		},
		Message: "Coinbase reward",
	}
}

func (t *Transaction) GetHash() []byte {
	var data string
	for _, tx := range t.Inputs {
		data = fmt.Sprintf("%s%s", data, tx.GetHash())
	}
	for _, tx := range t.Outputs {
		data = fmt.Sprintf("%s%s", data, tx.GetHash())
	}
	data = fmt.Sprintf("%s%s%d", data, t.Message, t.Fee)

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

func (t *TransactionInput) GetHash() []byte {
	data := fmt.Sprintf("%s%d", t.TransactionID, t.OutputIndex)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func (t *TransactionInput) Print() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}

func (t *TransactionOutput) GetHash() []byte {
	data := fmt.Sprintf("%s%d", t.Address, t.Amount)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func (t *TransactionOutput) Print() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}
