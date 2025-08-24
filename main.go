package main

import (
	"fmt"
	"go-bitcoin/blockchain"
	"go-bitcoin/transaction"
	"go-bitcoin/wallet"
)

func main() {
	createTestWallet()

	blockchain := blockchain.NewBlockchain()
	addTestTransactions(blockchain)

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())

	addMaliciousChange(blockchain)

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())
}

func addTestTransactions(blockchain *blockchain.Blockchain) {
	blockchain.AddTransaction(
		transaction.NewTransaction(
			"Filipe",
			"Anabela",
			"10",
			"To you, baby",
		),
	)
	blockchain.MineBlock()

	blockchain.AddTransaction(
		transaction.NewTransaction(
			"Anabela",
			"Maxwell",
			"5",
			"That's your payment",
		),
	)
	blockchain.MineBlock()

	blockchain.AddTransaction(
		transaction.NewTransaction(
			"Maxwell",
			"Duda",
			"2",
		),
	)
	blockchain.AddTransaction(
		transaction.NewTransaction(
			"Maxwell",
			"Filipe",
			"2",
			"Paying you that thing",
		),
	)
	blockchain.MineBlock()
}

func addMaliciousChange(blockchain *blockchain.Blockchain) {
	for i, b := range blockchain.Blocks {
		if i == 1 {
			txs := make([]*transaction.Transaction, 0)
			txs = append(txs, &transaction.Transaction{
				From:    "Anabela",
				To:      "Filipe",
				Amount:  "5",
				Message: "I STOLE!",
			})
			b.Transactions = txs
		}
	}
}

func createTestWallet() {
	wallet := wallet.NewWallet()
	fmt.Print(wallet.Print())
}
