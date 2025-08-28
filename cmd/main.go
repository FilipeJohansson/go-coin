package main

import (
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
)

func main() {
	wallet1 := createTestWallet()
	wallet2 := createTestWallet()
	wallet3 := createTestWallet()

	blockchain := blockchain.NewBlockchain(wallet1.Address)
	addTestTransactions(blockchain, wallet1, wallet2, wallet3)

	fmt.Printf("[BLOCKCHAIN]\n%s", blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())

	checkWalletFunds(blockchain, wallet1.Address)
	checkWalletFunds(blockchain, wallet2.Address)
	checkWalletFunds(blockchain, wallet3.Address)

	addMaliciousChange(blockchain)

	fmt.Printf("[BLOCKCHAIN]\n%s", blockchain.Print())
	fmt.Printf("Is Blockchain valid: %t\n", blockchain.IsBlockchainValid())
}

func addTestTransactions(blockchain *blockchain.Blockchain, wallet1 *wallet.Wallet, wallet2 *wallet.Wallet, wallet3 *wallet.Wallet) {
	checkWalletFunds(blockchain, wallet1.Address)
	tx, _ := wallet1.CreateTransaction(
		wallet2.Address,
		10,
		0.0015,
		blockchain.UTXOSet,
		"To you, baby",
	)
	wallet1.SignTransaction(tx)
	blockchain.AddTransaction(tx)
	blockchain.MineBlock(wallet1.Address)

	checkWalletFunds(blockchain, wallet2.Address)
	tx, _ = wallet2.CreateTransaction(
		wallet3.Address,
		4.5,
		0.0015,
		blockchain.UTXOSet,
		"That's your payment",
	)
	wallet2.SignTransaction(tx)
	blockchain.AddTransaction(tx)
	blockchain.MineBlock(wallet1.Address)

	checkWalletFunds(blockchain, wallet3.Address)
	tx, _ = wallet3.CreateTransaction(
		wallet1.Address,
		2,
		0.0015,
		blockchain.UTXOSet,
	)
	wallet3.SignTransaction(tx)
	blockchain.AddTransaction(tx)
	tx, _ = wallet3.CreateTransaction(
		wallet2.Address,
		1,
		0.0015,
		blockchain.UTXOSet,
		"Paying you that thing (not signed)",
	)
	// wallet3.SignTransaction(tx) <- its propositally commented so I can see the transaction being rejected
	blockchain.AddTransaction(tx)
	blockchain.MineBlock(wallet3.Address)

	checkWalletFunds(blockchain, wallet3.Address)
	tx, err := wallet3.CreateTransaction(
		wallet2.Address,
		10000,
		0.0015,
		blockchain.UTXOSet,
		"Paying you that thing",
	)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	wallet3.SignTransaction(tx)
	blockchain.AddTransaction(tx)
	blockchain.MineBlock(wallet1.Address)
}

func addMaliciousChange(blockchain *blockchain.Blockchain) {
	for i := range blockchain.Blocks {
		if i == 1 {
			maliciousTx := &transaction.Transaction{
				Inputs:  []transaction.TransactionInput{},
				Outputs: []transaction.TransactionOutput{{Address: "Hacker", Amount: 1000000}},
				Message: "I STOLE!",
			}
			blockchain.Blocks[1].Transactions = []*transaction.Transaction{maliciousTx}
		}
	}
}

func createTestWallet() *wallet.Wallet {
	wallet := wallet.NewWallet()
	fmt.Print(wallet.Print())
	return wallet
}

func checkWalletFunds(blockchain *blockchain.Blockchain, address string) {
	fmt.Printf("\n[%s] Funds: %.2f\n", address, float64(blockchain.UTXOSet.GetAddressBalance(address))/common.COINS_PER_UNIT)
}
