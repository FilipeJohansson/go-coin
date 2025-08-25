package main

import (
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
)

func main() {
	wallet1 := createTestWallet()
	wallet2 := createTestWallet()
	wallet3 := createTestWallet()

	blockchain := blockchain.NewBlockchain(wallet1.Address)
	addTestTransactions(blockchain, wallet1, wallet2, wallet3)

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())

	checkWalletFunds(blockchain, wallet1.Address)
	checkWalletFunds(blockchain, wallet2.Address)
	checkWalletFunds(blockchain, wallet3.Address)

	addMaliciousChange(blockchain)

	fmt.Print(blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())
}

func addTestTransactions(blockchain *blockchain.Blockchain, wallet1 *wallet.Wallet, wallet2 *wallet.Wallet, wallet3 *wallet.Wallet) {
	tx, _ := wallet1.CreateTransaction(
		wallet2.Address,
		10,
		"To you, baby",
	)
	wallet1.SignTransaction(tx)
	blockchain.AddTransaction(tx, *wallet1)
	blockchain.MineBlock(wallet1.Address)

	tx, _ = wallet2.CreateTransaction(
		wallet3.Address,
		4.5,
		"That's your payment",
	)
	wallet2.SignTransaction(tx)
	blockchain.AddTransaction(tx, *wallet2)
	blockchain.MineBlock(wallet1.Address)

	tx, _ = wallet3.CreateTransaction(
		wallet1.Address,
		2,
	)
	wallet3.SignTransaction(tx)
	blockchain.AddTransaction(tx, *wallet3)
	tx, _ = wallet3.CreateTransaction(
		wallet2.Address,
		1,
		"Paying you that thing",
	)
	// wallet3.SignTransaction(tx)
	blockchain.AddTransaction(tx, *wallet3)
	blockchain.MineBlock(wallet1.Address)

	tx, _ = wallet3.CreateTransaction(
		wallet2.Address,
		10000,
		"Paying you that thing",
	)
	wallet3.SignTransaction(tx)
	blockchain.AddTransaction(tx, *wallet3)
	blockchain.MineBlock(wallet1.Address)
}

func addMaliciousChange(blockchain *blockchain.Blockchain) {
	for i, b := range blockchain.Blocks {
		if i == 1 {
			txs := make([]*transaction.Transaction, 0)
			txs = append(txs, &transaction.Transaction{
				From:    "Anabela",
				To:      "Filipe",
				Amount:  5,
				Message: "I STOLE!",
			})
			b.Transactions = txs
		}
	}
}

func createTestWallet() *wallet.Wallet {
	wallet := wallet.NewWallet()
	fmt.Print(wallet.Print())
	return wallet
}

func checkWalletFunds(blockchain *blockchain.Blockchain, address string) {
	fmt.Printf("[%s] Funds: %d\n", address, blockchain.GetAddressFunds(address))
}
