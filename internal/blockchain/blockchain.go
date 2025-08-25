package blockchain

import (
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/block"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
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
	fmt.Printf("%s -> %s:\n", transaction.From, transaction.To)
	fmt.Printf("Amount: %d\n", transaction.Amount)
	fmt.Printf("Message: %s\n", transaction.Message)

	if transaction.Amount <= 0 {
		fmt.Printf("[INVALID] Invalid amount: %d\n", transaction.Amount)
		return
	}

	if transaction.From == transaction.To {
		fmt.Printf("[INVALID] Cannot send to yourself\n")
		return
	}

	if !wallet.ValidateTransactionSignature(*transaction) {
		// err
		fmt.Print("[INVALID] Transaction signature invalid\n")
		return
	}

	walletFunds := bc.GetAddressFunds(w.Address)
	if walletFunds < transaction.Amount {
		fmt.Print("[INVALID] Insuficient funds\n")
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

	fmt.Printf("[VALID] %s (%d) -> %s: %d\n", transaction.From, walletFunds, transaction.To, transaction.Amount)
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

func (bc *Blockchain) GetAddressFunds(address string) uint64 {
	var funds uint64
	for _, b := range bc.Blocks {
		for _, tx := range b.Transactions {
			if tx.To == address {
				funds += tx.Amount
			} else if tx.From == address {
				funds -= tx.Amount
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
		Amount: 50 * common.COINS_PER_UNIT,
	}
}
