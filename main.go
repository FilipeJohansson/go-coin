package main

import (
	"fmt"
	"go-bitcoin/blockchain"
)

func main() {
	blockchain := blockchain.NewBlockchain()
	blockchain.NewBlock("SECOND BLOCK, YEAH!")
	blockchain.NewBlock("Wow, thats the third block :o")

	blockchain.Print()
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())

	for i, b := range blockchain.Blocks {
		if i == 1 {
			b.Data = "DATA CHANGED"
		}
	}

	blockchain.Print()
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())
}
