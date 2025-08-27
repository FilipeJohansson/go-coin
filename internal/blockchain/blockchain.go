package blockchain

import (
	"encoding/hex"
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/block"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/utxo"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
)

type Blockchain struct {
	Blocks       []*block.Block `json:"blocks"`
	UTXOSet      *utxo.UTXOSet  `json:"utxoSet"`
	PendingBlock *block.Block
}

func NewBlockchain(genesisWalletAddress string) *Blockchain {
	blockchain := &Blockchain{
		UTXOSet:      utxo.NewUTXOSet(),
		PendingBlock: block.NewBlock("", "Genesis block"),
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
		from = common.GetAddressFromPublicKey(tx.Inputs[0].PublicKey)
	}

	fmt.Printf("%s -> %s:\n", from, to)
	fmt.Printf("Amount: %.2f\n", amount)
	fmt.Printf("Message: %s\n", tx.Message)

	if len(tx.Inputs) == 0 {
		bc.addTxToPendingBlock(tx)
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

	if totalInputs < totalOutputs {
		fmt.Print("[INVALID] Insuficient funds\n")
		return
	}

	fmt.Printf("[VALID] %s -> %s: %.2f\n", from, to, amount)

	bc.addTxToPendingBlock(tx)
}

func (bc *Blockchain) MineBlock(minerAddress string) {
	if bc.PendingBlock == nil {
		return
	}

	miningRewardTx := bc.createCoinbaseTransaction(minerAddress)
	bc.AddTransaction(miningRewardTx)
	lastTx := bc.PendingBlock.Transactions[len(bc.PendingBlock.Transactions)-1]
	bc.PendingBlock.Transactions = append(
		[]*transaction.Transaction{lastTx},
		bc.PendingBlock.Transactions[:len(bc.PendingBlock.Transactions)-1]...,
	)

	bc.PendingBlock.Mine(2)

	for _, tx := range bc.PendingBlock.Transactions {
		bc.updateUTXOSet(tx)
	}

	bc.Blocks = append(bc.Blocks, bc.PendingBlock)

	bc.PendingBlock = nil
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

func (bc *Blockchain) Print() string {
	var formattedBlockchain string
	for _, b := range bc.Blocks {
		formattedBlockchain += b.Print()
	}

	return formattedBlockchain
}

func (bc *Blockchain) addTxToPendingBlock(tx *transaction.Transaction) {
	var b *block.Block
	if bc.PendingBlock == nil {
		blockchainLen := len(bc.Blocks)
		lastBlockchainBlockHash := bc.Blocks[blockchainLen-1].BlockHash
		b = block.NewBlock(lastBlockchainBlockHash)
		bc.PendingBlock = b
	} else {
		b = bc.PendingBlock
	}

	b.AddTransaction(tx)
}

func (bc *Blockchain) createCoinbaseTransaction(address string) *transaction.Transaction {
	return transaction.NewCoinbaseTransaction(address, 50*common.COINS_PER_UNIT)
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

	if !wallet.ValidateTransactionSignature(*tx) {
		return false
	}

	var totalInputs uint64
	for _, input := range tx.Inputs {
		if !tempUTXOSet.UTXOExists(input.TransactionID, input.OutputIndex) {
			return false
		}

		utxo := tempUTXOSet.GetUTXO(input.TransactionID, input.OutputIndex)

		from := common.GetAddressFromPublicKey(input.PublicKey)
		if utxo.Address != from {
			return false
		}

		totalInputs += utxo.Amount
	}

	var totalOutputs uint64
	for _, output := range tx.Outputs {
		totalOutputs += output.Amount
	}

	if totalInputs < totalOutputs {
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
