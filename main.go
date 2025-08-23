package main

import "go-bitcoin/blockchain"

func main() {
	blockchain := blockchain.NewBlockchain()
	blockchain.NewBlock("SECOND BLOCK, YEAH!")
	blockchain.NewBlock("Wow, thats the third block :o")
	blockchain.Print()
}
