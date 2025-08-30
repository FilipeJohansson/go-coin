package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/FilipeJohansson/go-coin/internal/utxo"
	"github.com/FilipeJohansson/go-coin/pkg/common"
)

type CustomPublicKey struct {
	Curve elliptic.Curve `json:"-"`
	X     *big.Int       `json:"X"`
	Y     *big.Int       `json:"Y"`
}

type TransactionInput struct {
	TransactionID string          `json:"transactionID"`
	OutputIndex   uint            `json:"outputIndex"`
	Signature     string          `json:"signature"`
	PublicKey     CustomPublicKey `json:"publicKey"`
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
			PublicKey: CustomPublicKey{
				Curve: senderPublicKey.Curve,
				X:     senderPublicKey.X,
				Y:     senderPublicKey.Y,
			},
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

func (t *Transaction) Json() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
		return ""
	}

	return fmt.Sprintf("%s\n", json)
}

func (t *Transaction) Print() string {
	var inputs string
	for _, i := range t.Inputs {
		inputs += fmt.Sprintf("%s", i.Print())
	}

	var outputs string
	for _, o := range t.Outputs {
		outputs += fmt.Sprintf("%s", o.Print())
	}

	return fmt.Sprintf(`
Fee: %d
Message: %s
Inputs:
%s
Outputs:
%s`, t.Fee, t.Message, inputs, outputs)
}

func (t *TransactionInput) GetHash() []byte {
	data := fmt.Sprintf("%s%d", t.TransactionID, t.OutputIndex)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func (t *TransactionInput) Json() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}

func (t *TransactionInput) Print() string {
	return common.BuildBox(
		fmt.Sprintf("Transaction ID: %s", t.TransactionID),
		fmt.Sprintf("Output Index:   %d", t.OutputIndex),
		fmt.Sprintf("Public Key:     %s", common.GetAddressFromPublicKey(*t.PublicKey.GetPublicKey())),
		fmt.Sprintf("Signature:      %s", t.Signature),
	)
}

func (t *TransactionOutput) GetHash() []byte {
	data := fmt.Sprintf("%s%d", t.Address, t.Amount)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func (t *TransactionOutput) Json() string {
	json, err := json.Marshal(t)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}

func (t *TransactionOutput) Print() string {
	return common.BuildBox(
		fmt.Sprintf("Address: %s", t.Address),
		fmt.Sprintf("Amount:  %d", t.Amount),
	)
}

func (c *CustomPublicKey) GetPublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: c.Curve,
		X:     c.X,
		Y:     c.Y,
	}
}
