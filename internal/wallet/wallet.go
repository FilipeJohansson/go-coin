package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/utxo"
	"github.com/FilipeJohansson/go-coin/pkg/common"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey `json:"-"`
	PublicKey  ecdsa.PublicKey  `json:"publicKey"`
	Address    string           `json:"address"`
}

func NewWallet() *Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		// err
	}

	publicKey := &privateKey.PublicKey

	wallet := &Wallet{
		PrivateKey: *privateKey,
		PublicKey:  *publicKey,
	}
	wallet.GetAddress()

	return wallet
}

func (w *Wallet) CreateTransaction(to string, amount float64, utxoSet *utxo.UTXOSet, msg ...string) (*transaction.Transaction, error) {
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

	tx, err := transaction.NewTransaction(w.Address, to, uint64(amount*common.COINS_PER_UNIT), utxoSet, w.PublicKey, message)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (w *Wallet) SignTransaction(tx *transaction.Transaction) {
	if tx == nil {
		return
	}

	for i := range tx.Inputs {
		signatureBytes, err := ecdsa.SignASN1(rand.Reader, &w.PrivateKey, tx.Inputs[i].GetHash())
		if err != nil {
			// err
		}
		tx.Inputs[i].Signature = hex.EncodeToString(signatureBytes)
	}
}

func (w *Wallet) GetAddress() string {
	w.Address = common.GetAddressFromPublicKey(w.PublicKey)
	return w.Address
}

func (w *Wallet) Print() string {
	json, err := json.Marshal(w)
	if err != nil {
		// error
	}

	return fmt.Sprintf("%s\n", json)
}

func ValidateTransactionSignature(tx transaction.Transaction) bool {
	if len(tx.Inputs) == 0 {
		return true // Coinbase transaction
	}

	for _, i := range tx.Inputs {
		if i.Signature == "" {
			return false
		}

		signature, err := hex.DecodeString(i.Signature)
		if err != nil {
			// err
			return false
		}

		if !ecdsa.VerifyASN1(&i.PublicKey, i.GetHash(), signature) {
			return false
		}
	}

	return true
}
