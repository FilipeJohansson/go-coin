package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Timestamp     time.Time `json:"timestamp"`
	Data          string    `json:"data"`
	PrevBlockHash string    `json:"prevBlockHash"`
	BlockHash     string    `json:"blockHash"`
	Nonce         int       `json:"nonce"`
}

func NewBlock(data string, prevBlockHash string) *Block {
	return &Block{
		Timestamp:     time.Now(),
		Data:          data,
		PrevBlockHash: prevBlockHash,
	}
}

func (b *Block) CalcBlockHash() {
	data := fmt.Sprintf("%v%s%s%d",
		b.Timestamp.Unix(),
		b.Data,
		b.PrevBlockHash,
		b.Nonce)

	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to hexadecimal string
	b.BlockHash = hex.EncodeToString(hashBytes)
}

func (b *Block) Mine(difficulty int) {
	target := strings.Repeat("0", difficulty)

	for {
		b.CalcBlockHash()
		if strings.HasPrefix(b.BlockHash, target) {
			break
		}
		b.Nonce++
	}
}

func (b *Block) Print() {
	json, err := json.Marshal(b)
	if err != nil {
		// error
	}

	fmt.Printf("%s", json)
}
