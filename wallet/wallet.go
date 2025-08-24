package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
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

func (w *Wallet) GetAddress() string {
	data := fmt.Sprintf("%s%s", w.PublicKey.X, w.PublicKey.Y)

	// First hash - SHA256
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Second hash - SHA256
	hasher.Write(hashBytes)
	hashBytes = hasher.Sum(nil)

	// Address - RIPEMD160
	hasher = ripemd160.New()
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
