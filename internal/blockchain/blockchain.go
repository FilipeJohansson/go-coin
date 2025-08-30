package blockchain

import (
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/FilipeJohansson/go-coin/internal/block"
	"github.com/FilipeJohansson/go-coin/internal/mempool"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/utxo"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
)

type Blockchain struct {
	Blocks  []*block.Block   `json:"blocks"`
	UTXOSet *utxo.UTXOSet    `json:"-"`
	Mempool *mempool.Mempool `json:"mempool"`
}

func NewBlockchain(genesisWalletAddress string, filename ...string) *Blockchain {
	if len(filename) > 0 && filename[0] != "" {
		blockchain, err := LoadFromFile(filename[0])
		if err != nil {
			// err
		} else {
			return blockchain
		}
	}

	blockchain := &Blockchain{
		UTXOSet: utxo.NewUTXOSet(),
		Mempool: mempool.NewMempool(),
	}

	blockchain.MineBlock(genesisWalletAddress)

	return blockchain
}

func (bc *Blockchain) AddTransaction(tx *transaction.Transaction) {
	if tx == nil {
		return
	}

	var from string
	to := tx.Outputs[0].Address
	amount := float64(tx.Outputs[0].Amount) / common.COINS_PER_UNIT
	if len(tx.Inputs) == 0 {
		from = "Coinbase"
	} else {
		from = common.GetAddressFromPublicKey(*tx.Inputs[0].PublicKey.GetPublicKey())
	}

	fmt.Printf("%s -> %s:\n", from, to)
	fmt.Printf("Amount: %.2f\n", amount)
	fmt.Printf("Message: %s\n", tx.Message)

	if len(tx.Inputs) == 0 {
		// Coinbase
		bc.Mempool.AddTransaction(tx)
		return
	}

	if tx.Fee < common.MIN_FEE {
		// err
		return
	}

	var totalOutputs uint64
	for _, o := range tx.Outputs {
		totalOutputs += o.Amount
	}

	if amount <= 0 {
		fmt.Printf("[INVALID] Invalid amount: %.2f\n", amount)
		return
	}

	if from == to {
		fmt.Printf("[INVALID] Cannot send to yourself\n")
		return
	}

	if !wallet.ValidateTransactionSignature(*tx) {
		// err
		fmt.Print("[INVALID] Transaction signature invalid\n")
		return
	}

	var totalInputs uint64
	for _, input := range tx.Inputs {
		if !bc.UTXOSet.UTXOExists(input.TransactionID, input.OutputIndex) {
			fmt.Print("[INVALID] Input UTXO does not exist\n")
			return
		}

		utxo := bc.UTXOSet.GetUTXO(input.TransactionID, input.OutputIndex)
		if utxo.Address != from {
			fmt.Print("[INVAID] Input UTXO does not belong to sender\n")
			return
		}
		totalInputs += utxo.Amount
	}

	if totalInputs < totalOutputs+tx.Fee {
		fmt.Print("[INVALID] Insuficient funds\n")
		return
	}

	fmt.Printf("[VALID] %s -> %s: %.2f\n", from, to, amount)

	bc.Mempool.AddTransaction(tx)
}

func (bc *Blockchain) MineBlock(minerAddress string) {
	if bc.Mempool.Size() == 0 && len(bc.Blocks) > 0 {
		return
	}

	var prevHash string
	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].BlockHash
	}

	newBlock := block.NewBlock(prevHash)

	var totalFees uint64
	transactions := bc.Mempool.GetTransactions()
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Fee > transactions[j].Fee
	})
	for _, tx := range transactions {
		if len(newBlock.Transactions) == 10 {
			break
		}
		newBlock.AddTransaction(tx)
		totalFees += tx.Fee
	}

	coinbaseTx := bc.createCoinbaseTransaction(minerAddress, totalFees)
	newBlock.Transactions = append([]*transaction.Transaction{coinbaseTx}, newBlock.Transactions...)

	newBlock.Mine(bc.calculateDifficulty())
	fmt.Printf("[NEW MINED BLOCK]\n%s", newBlock.Print())

	for _, tx := range newBlock.Transactions {
		bc.updateUTXOSet(tx)
	}

	bc.Mempool.CleanProcessedTransactions(newBlock.Transactions)

	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) IsBlockchainValid() bool {
	tempUTXOSet := utxo.NewUTXOSet()

	for i, b := range bc.Blocks {
		if !b.IsHashRight() {
			return false
		}

		if i != 0 {
			if b.PrevBlockHash != bc.Blocks[i-1].BlockHash {
				return false
			}
		} else if b.PrevBlockHash != "" {
			return false
		}

		for _, tx := range b.Transactions {
			if !bc.validateTransactionInContext(tx, tempUTXOSet) {
				return false
			}

			bc.applyTransactionToUTXOSet(tx, tempUTXOSet)
		}
	}

	return true
}

func (bc *Blockchain) SaveToFile(filename string) error {
	json, err := json.MarshalIndent(bc, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, json, 0644)
}

func (bc *Blockchain) Print() string {
	var formattedBlockchain string
	for _, b := range bc.Blocks {
		formattedBlockchain += b.Print()
	}

	return formattedBlockchain
}

func (bc *Blockchain) createCoinbaseTransaction(address string, totalFees uint64) *transaction.Transaction {
	return transaction.NewCoinbaseTransaction(address, (50*common.COINS_PER_UNIT)+totalFees)
}

func (bc *Blockchain) rebuildUTXOSet() {
	bc.UTXOSet = utxo.NewUTXOSet()
	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			bc.updateUTXOSet(tx)
		}
	}
}

func (bc *Blockchain) updateUTXOSet(tx *transaction.Transaction) {
	for _, i := range tx.Inputs {
		bc.UTXOSet.RemoveUTXOByID(i.TransactionID, i.OutputIndex)
	}

	txID := hex.EncodeToString(tx.GetHash())
	for i, o := range tx.Outputs {
		newUTXO := &utxo.UTXO{
			TransactionID: txID,
			OutputIndex:   uint(i),
			Address:       o.Address,
			Amount:        o.Amount,
		}
		bc.UTXOSet.AddUTXO(newUTXO)
	}
}

func (bc *Blockchain) validateTransactionInContext(tx *transaction.Transaction, tempUTXOSet *utxo.UTXOSet) bool {
	if len(tx.Inputs) == 0 {
		if len(tx.Outputs) != 1 {
			return false
		}

		return true
	}

	if tx.Fee < common.MIN_FEE {
		return false
	}

	if !wallet.ValidateTransactionSignature(*tx) {
		return false
	}

	var totalInputs uint64
	for _, input := range tx.Inputs {
		if !tempUTXOSet.UTXOExists(input.TransactionID, input.OutputIndex) {
			return false
		}

		utxo := tempUTXOSet.GetUTXO(input.TransactionID, input.OutputIndex)

		from := common.GetAddressFromPublicKey(*input.PublicKey.GetPublicKey())
		if utxo.Address != from {
			return false
		}

		totalInputs += utxo.Amount
	}

	var totalOutputs uint64
	for _, output := range tx.Outputs {
		totalOutputs += output.Amount
	}

	if totalInputs < totalOutputs+tx.Fee {
		return false
	}

	return true
}

func (bc *Blockchain) applyTransactionToUTXOSet(tx *transaction.Transaction, tempUTXOSet *utxo.UTXOSet) {
	for _, input := range tx.Inputs {
		tempUTXOSet.RemoveUTXOByID(input.TransactionID, input.OutputIndex)
	}

	txID := hex.EncodeToString(tx.GetHash())
	for i, output := range tx.Outputs {
		newUTXO := &utxo.UTXO{
			TransactionID: txID,
			OutputIndex:   uint(i),
			Address:       output.Address,
			Amount:        output.Amount,
		}
		tempUTXOSet.AddUTXO(newUTXO)
	}
}

func (bc *Blockchain) calculateDifficulty() int {
	if len(bc.Blocks) < common.DIFFICULTY_ADJUSTMENT_INTERVAL {
		return common.INITIAL_DIFFICULTY
	}

	startIndex := len(bc.Blocks) - common.DIFFICULTY_ADJUSTMENT_INTERVAL
	lastBlocks := bc.Blocks[startIndex:]

	var totalTime time.Duration
	for i := 1; i < len(lastBlocks); i++ {
		blockTime := lastBlocks[i].Timestamp.Sub(lastBlocks[i-1].Timestamp)
		totalTime += blockTime
	}

	avgTime := totalTime / time.Duration(len(lastBlocks)-1)

	currentDifficulty := lastBlocks[len(lastBlocks)-1].Difficulty
	targetTime := time.Duration(common.TARGET_BLOCK_TIME) * time.Second

	if avgTime < targetTime {
		return currentDifficulty + 1
	} else if avgTime > targetTime {
		if currentDifficulty > 1 {
			return currentDifficulty - 1
		}
		return 1
	} else {
		return currentDifficulty
	}
}

func (bc *Blockchain) fixPublicKeyCurves() {
	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			for i := range tx.Inputs {
				tx.Inputs[i].PublicKey.Curve = elliptic.P256()
			}
		}
	}
}

func LoadFromFile(filename string) (*Blockchain, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var blockchain Blockchain
	err = json.Unmarshal(content, &blockchain)
	if err != nil {
		return nil, err
	}

	blockchain.rebuildUTXOSet()
	blockchain.fixPublicKeyCurves()

	if !blockchain.IsBlockchainValid() {
		return nil, errors.New("Blockchain saved is invalid")
	}

	return &blockchain, nil
}
