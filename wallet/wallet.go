package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FilipeJohansson/go-coin/transaction"

	"github.com/btcsuite/btcutil/base58"
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

func (w *Wallet) CreateTransaction(to string, amount float64, msg ...string) (*transaction.Transaction, error) {
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

	tx, err := transaction.NewTransaction(w.Address, to, amount, w.PublicKey, message)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (w *Wallet) SignTransaction(tx *transaction.Transaction) {
	signatureBytes, err := ecdsa.SignASN1(rand.Reader, &w.PrivateKey, tx.GetHash())
	if err != nil {
		// err
	}
	tx.Signature = hex.EncodeToString(signatureBytes)
}

func (w *Wallet) GetAddress() string {
	data := fmt.Sprintf("%s%s", w.PublicKey.X, w.PublicKey.Y)

	// First hash - SHA256
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Second hash - SHA256
	hasher = sha256.New()
	hasher.Write(hashBytes)
	addressBytes := hasher.Sum(nil)

	w.Address = base58.Encode(addressBytes)

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
	if tx.From == "MINING_REWARD" {
		return true
	}

	if tx.Signature == "" {
		return false
	}

	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		// err
		return false
	}
	return ecdsa.VerifyASN1(&tx.PublicKey, tx.GetHash(), signature)
}
