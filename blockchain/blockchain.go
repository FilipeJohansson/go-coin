package blockchain

import (
	"fmt"
	"go-bitcoin/block"
	"go-bitcoin/transaction"
	"go-bitcoin/wallet"
	"strconv"
)

type Blockchain struct {
	Blocks       []*block.Block `json:"blocks"`
	PendingBlock *block.Block
}

func NewBlockchain() *Blockchain {
	genesis := block.NewBlock("", "Genesis block")
	genesis.Mine(2)

	chain := make([]*block.Block, 0)
	chain = append(chain, genesis)

	return &Blockchain{
		Blocks: chain,
	}
}

func (bc *Blockchain) AddTransaction(transaction *transaction.Transaction, w wallet.Wallet) {
	if !wallet.ValidateTransactionSignature(*transaction) {
		// err
		fmt.Printf("[INVALID] %s -> %s: %s (%s)\n", transaction.From, transaction.To, transaction.Amount, transaction.Message)
		return
	}

	amount, err := strconv.Atoi(transaction.Amount)
	if err != nil {
		// err
		return
	}
	walletFunds := bc.getAddressFunds(w.Address)
	if walletFunds < amount {
		fmt.Printf("[INSUFFICIENT_FUNDS] %s -> %s: %s (%s)\n", transaction.From, transaction.To, transaction.Amount, transaction.Message)
		return
	}

	var b *block.Block
	if bc.PendingBlock == nil {
		blockchainLen := len(bc.Blocks)
		lastBlockchainBlockHash := bc.Blocks[blockchainLen-1].BlockHash
		b = block.NewBlock(lastBlockchainBlockHash)
		bc.PendingBlock = b
	} else {
		b = bc.PendingBlock
	}

	fmt.Printf("[VALID] %s -> %s: %s (%s)\n", transaction.From, transaction.To, transaction.Amount, transaction.Message)
	b.AddTransaction(transaction)
}

func (bc *Blockchain) MineBlock() {
	if bc.PendingBlock == nil {
		return
	}

	bc.PendingBlock.Mine(2)
	bc.Blocks = append(bc.Blocks, bc.PendingBlock)

	bc.PendingBlock = nil
}

func (bc *Blockchain) IsBlockchainValid() bool {
	for i, b := range bc.Blocks {
		if i != 0 {
			if b.PrevBlockHash != bc.Blocks[i-1].BlockHash {
				return false
			}
		} else if b.PrevBlockHash != "" {
			return false
		}

		if !b.IsHashRight() {
			return false
		}

		for _, tx := range b.Transactions {
			if !wallet.ValidateTransactionSignature(*tx) {
				return false
			}
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

func (bc *Blockchain) getAddressFunds(address string) int {
	var funds int
	for _, b := range bc.Blocks {
		for _, tx := range b.Transactions {
			if tx.To == address {
				amount, err := strconv.Atoi(tx.Amount)
				if err != nil {
					// err
					continue
				}

				funds += amount
			} else if tx.From == address {
				amount, err := strconv.Atoi(tx.Amount)
				if err != nil {
					// err
					continue
				}

				funds -= amount
			}
		}
	}

	return funds
}
