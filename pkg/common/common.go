package common

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

const COINS_PER_UNIT = 1000000
const MIN_FEE = 1000

const INITIAL_DIFFICULTY = 2
const TARGET_BLOCK_TIME = 1              // in seconds
const DIFFICULTY_ADJUSTMENT_INTERVAL = 5 // each n blocks

func GetAddressFromPublicKey(key ecdsa.PublicKey) string {
	data := append(key.X.Bytes(), key.Y.Bytes()...)

	// First hash - SHA256
	hasher := sha256.New()
	hasher.Write(data)
	hashBytes := hasher.Sum(nil)

	// Second hash - SHA256
	hasher = sha256.New()
	hasher.Write(hashBytes)
	addressBytes := hasher.Sum(nil)

	return base58.Encode(addressBytes)
}

func GetPublicKeyHash(key ecdsa.PublicKey) string {
	data := append(key.X.Bytes(), key.Y.Bytes()...)
	return base58.Encode([]byte(data))
}

func GetPrivateKeyHash(key ecdsa.PrivateKey) string {
	return base58.Encode(key.D.Bytes())
}

func GetPrivateKeyFromHash(encoded string) *ecdsa.PrivateKey {
	decoded := base58.Decode(encoded)
	if len(decoded) == 0 {
		return nil
	}

	if len(decoded) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(decoded):], decoded)
		decoded = padded
	}

	d := new(big.Int).SetBytes(decoded)
	curve := elliptic.P256()

	privateKey := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
	}

	ecdhKey, err := ecdh.P256().NewPrivateKey(decoded)
	if err != nil {
		// err
		return nil
	}

	ecdhPubKey := ecdhKey.PublicKey()
	pubKeyBytes := ecdhPubKey.Bytes()

	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 0x04 {
		// err
		return nil
	}

	privateKey.PublicKey.X = new(big.Int).SetBytes(pubKeyBytes[1:33])
	privateKey.PublicKey.Y = new(big.Int).SetBytes(pubKeyBytes[33:65])

	return privateKey
}

func BuildBox(lines ...string) string {
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	border := strings.Repeat("-", maxLen+4)

	var sb strings.Builder
	sb.WriteString(border + "\n")
	for _, l := range lines {
		sb.WriteString(fmt.Sprintf("| %-*s |\n", maxLen, l))
	}
	sb.WriteString(border + "\n")

	return sb.String()
}
