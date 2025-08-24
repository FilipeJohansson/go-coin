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
	Difficulty    int       `json:"difficulty"`
}

func NewBlock(data string, prevBlockHash string) *Block {
	return &Block{
		Timestamp:     time.Now(),
		Data:          data,
		PrevBlockHash: prevBlockHash,
	}
}

func (b *Block) Mine(difficulty int) {
	b.Difficulty = difficulty
	target := strings.Repeat("0", difficulty)

	for {
		b.SaveBlockHash()
		if strings.HasPrefix(b.BlockHash, target) {
			break
		}
		b.Nonce++
	}
}

func (b *Block) SaveBlockHash() {
	b.BlockHash = b.GetHash()
}

func (b *Block) GetHash() string {
	data := fmt.Sprintf("%v%s%s%d",
		b.Timestamp.Unix(),
		b.Data,
		b.PrevBlockHash,
		b.Nonce)

	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to hexadecimal string
	return hex.EncodeToString(hashBytes)
}

func (b *Block) IsHashRight() bool {
	target := strings.Repeat("0", b.Difficulty)
	if !strings.HasPrefix(b.BlockHash, target) {
		return false
	}

	if b.GetHash() != b.BlockHash {
		return false
	}

	return true
}

func (b *Block) Print() {
	json, err := json.Marshal(b)
	if err != nil {
		// error
	}

	fmt.Printf("%s\n", json)
}
