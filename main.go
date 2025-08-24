package main

import (
	"fmt"
	"go-bitcoin/blockchain"
	"go-bitcoin/transaction"
)

func main() {
	blockchain := blockchain.NewBlockchain()
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

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())

	// Malicious change
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

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())
}
