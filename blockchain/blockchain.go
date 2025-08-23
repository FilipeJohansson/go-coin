package blockchain

import "go-bitcoin/block"

type Blockchain struct {
	Blocks []*block.Block `json:"blocks"`
}

func NewBlockchain() *Blockchain {
	genesis := block.NewBlock("Genesis block data", "")
	genesis.Mine(2)

	chain := make([]*block.Block, 0)
	chain = append(chain, genesis)

	return &Blockchain{
		Blocks: chain,
	}
}

func (bc *Blockchain) Print() {
	for _, b := range bc.Blocks {
		b.Print()
	}
}
