package main

import "go-bitcoin/blockchain"

func main() {
	blockchain := blockchain.NewBlockchain()
	blockchain.Print()
}
