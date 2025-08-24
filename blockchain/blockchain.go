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

func NewBlockchain(genesisWalletAddress string) *Blockchain {
	blockchain := &Blockchain{}

	genesis := block.NewBlock("", "Genesis block")
	genesisTx := blockchain.createMiningRewardTransaction(genesisWalletAddress)
	genesis.AddTransaction(genesisTx)
	genesis.Mine(2)

	chain := make([]*block.Block, 0)
	chain = append(chain, genesis)

	blockchain.Blocks = chain

	return blockchain
}

func (bc *Blockchain) AddTransaction(transaction *transaction.Transaction, w wallet.Wallet) {

	if !wallet.ValidateTransactionSignature(*transaction) {
		// err
		fmt.Printf("[INVALID] %s -> %s: %s (%s)\n", transaction.From, transaction.To, transaction.Amount, transaction.Message)
		return
	}

	txAmount, err := strconv.Atoi(transaction.Amount)
	if err != nil {
		// err
		return
	}
	walletFunds := bc.GetAddressFunds(w.Address)
	if walletFunds < txAmount {
		fmt.Printf("[INSUFFICIENT_FUNDS] %s (%d) -> %s: %s (%s)\n", transaction.From, walletFunds, transaction.To, transaction.Amount, transaction.Message)
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

	fmt.Printf("[VALID] %s (%d) -> %s: %s (%s)\n", transaction.From, walletFunds, transaction.To, transaction.Amount, transaction.Message)
	b.AddTransaction(transaction)
}

func (bc *Blockchain) MineBlock(minerAddress string) {
	if bc.PendingBlock == nil {
		return
	}

	miningRewardTx := bc.createMiningRewardTransaction(minerAddress)
	bc.PendingBlock.AddTransaction(miningRewardTx)

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

func (bc *Blockchain) GetAddressFunds(address string) int {
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

func (bc *Blockchain) Print() string {
	var formattedBlockchain string
	for _, b := range bc.Blocks {
		formattedBlockchain += b.Print()
	}

	return formattedBlockchain
}

func (bc *Blockchain) createMiningRewardTransaction(address string) *transaction.Transaction {
	return &transaction.Transaction{
		From:   "MINING_REWARD",
		To:     address,
		Amount: "50",
	}
}
