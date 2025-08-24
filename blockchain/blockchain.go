package blockchain

import (
	"go-bitcoin/block"
	"go-bitcoin/transaction"
)

type Blockchain struct {
	Blocks       []*block.Block `json:"blocks"`
	PendingBlock *block.Block
}

func NewBlockchain() *Blockchain {
	genesis := block.NewBlock("")
	genesis.AddTransaction(transaction.NewTransaction("", "", "", "Genesis block"))
	genesis.Mine(2)

	chain := make([]*block.Block, 0)
	chain = append(chain, genesis)

	return &Blockchain{
		Blocks: chain,
	}
}

func (bc *Blockchain) AddTransaction(transaction *transaction.Transaction) {
	var bk *block.Block
	if bc.PendingBlock == nil {
		blockchainLen := len(bc.Blocks)
		lastBlockchainBlockHash := bc.Blocks[blockchainLen-1].BlockHash
		bk = block.NewBlock(lastBlockchainBlockHash)
		bc.PendingBlock = bk
	} else {
		bk = bc.PendingBlock
	}

	bk.AddTransaction(transaction)
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
