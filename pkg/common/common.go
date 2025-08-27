package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

const COINS_PER_UNIT = 1000000

func GetAddressFromPublicKey(key ecdsa.PublicKey) string {
	data := fmt.Sprintf("%s%s", key.X, key.Y)

	// First hash - SHA256
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Second hash - SHA256
	hasher = sha256.New()
	hasher.Write(hashBytes)
	addressBytes := hasher.Sum(nil)

	return base58.Encode(addressBytes)
}
