package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
	"github.com/FilipeJohansson/go-coin/internal/transaction"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
	"github.com/spf13/cobra"
)

var transactionCmd = &cobra.Command{
	Use:     "transaction",
	Aliases: []string{"tx"},
	Short:   "Transaction operations",
}

var sendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "Send coins to another wallet",
	Run:     sendTransaction,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List pending transactions",
	Run:     listPendingTransactions,
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "g"},
	Short:   "Generate multiple transactions for testing",
	Long:    "Generate a specified number of transactions between randomly created wallets",
	Run:     generateTransactions,
}

func init() {
	sendCmd.Flags().StringP("to", "t", "", "Recipient address")
	sendCmd.Flags().StringP("private-key", "p", "", "The from address private key to autenticate")
	sendCmd.Flags().Float64P("amount", "a", 0, "Quantity to send from sender to recipient")
	sendCmd.Flags().Float64("fee", float64(common.MIN_FEE)/float64(common.COINS_PER_UNIT), "Optional miners fee")
	sendCmd.Flags().StringP("message", "m", "", "Optional message")

	generateCmd.Flags().IntP("count", "c", 10, "Number of transactions to generate")
	generateCmd.Flags().IntP("wallets", "w", 5, "Number of wallets to create and use")
	generateCmd.Flags().Float64("min-amount", 0.1, "Minimum transaction amount")
	generateCmd.Flags().Float64("max-amount", 10.0, "Maximum transaction amount")
	generateCmd.Flags().Float64("fee", float64(common.MIN_FEE)/float64(common.COINS_PER_UNIT), "Transaction fee")
	generateCmd.Flags().Bool("fund-wallets", true, "Create funding transactions for wallets")

	transactionCmd.AddCommand(sendCmd)
	transactionCmd.AddCommand(listCmd)
	transactionCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(transactionCmd)
}

func sendTransaction(cmd *cobra.Command, args []string) {
	to, _ := cmd.Flags().GetString("to")
	if to == "" {
		fmt.Println("Error: recipient address is required")
		return
	}

	privateKey, _ := cmd.Flags().GetString("private-key")
	if privateKey == "" {
		fmt.Println("Error: private key is required")
		return
	}

	amount, _ := cmd.Flags().GetFloat64("amount")
	fee, _ := cmd.Flags().GetFloat64("fee")
	message, _ := cmd.Flags().GetString("message")

	wallet := wallet.LoadWallet(privateKey)

	blockchain := blockchain.NewBlockchain("", blockchainFile)

	tx, err := wallet.CreateTransaction(
		to,
		amount,
		fee,
		blockchain.UTXOSet,
		message,
	)
	if err != nil {
		fmt.Printf("Error to create transaction: %s", err.Error())
		return
	}
	wallet.SignTransaction(tx)
	blockchain.AddTransaction(tx)

	err = blockchain.SaveToFile(blockchainFile)
	if err != nil {
		fmt.Printf("Error to save Blockchain: %v\n", err)
	}
}

func listPendingTransactions(cmd *cobra.Command, args []string) {
	blockchain := blockchain.NewBlockchain("", blockchainFile)
	fmt.Println(blockchain.Mempool.Print())
}

func generateTransactions(cmd *cobra.Command, args []string) {
	count, _ := cmd.Flags().GetInt("count")
	walletCount, _ := cmd.Flags().GetInt("wallets")
	minAmount, _ := cmd.Flags().GetFloat64("min-amount")
	maxAmount, _ := cmd.Flags().GetFloat64("max-amount")
	fee, _ := cmd.Flags().GetFloat64("fee")
	fundWallets, _ := cmd.Flags().GetBool("fund-wallets")

	if count <= 0 {
		fmt.Println("Error: count must be positive")
		return
	}

	if walletCount < 2 {
		fmt.Println("Error: need at least 2 wallets")
		return
	}

	blockchain := blockchain.NewBlockchain("", blockchainFile)

	// Create wallets for testing
	fmt.Printf("Creating %d test wallets...\n", walletCount)
	wallets := make([]*wallet.Wallet, walletCount)
	for i := 0; i < walletCount; i++ {
		wallets[i] = wallet.NewWallet()
		fmt.Printf("Wallet %d: %s\n", i+1, wallets[i].Address)
	}

	// Fund wallets if requested
	if fundWallets {
		fmt.Printf("\nFunding wallets with initial coins...\n")
		for i, w := range wallets {
			// Create coinbase-like transaction to fund each wallet
			fundingTx := transaction.NewCoinbaseTransaction(w.Address, 1000*common.COINS_PER_UNIT) // 1000 coins each
			blockchain.AddTransaction(fundingTx)
			fmt.Printf("Funded wallet %d with 1000 coins\n", i+1)
		}

		// Mine funding transactions first
		fmt.Println("Mining funding transactions...")
		blockchain.MineBlock(wallets[0].Address) // Use first wallet as miner
		err := blockchain.SaveToFile(blockchainFile)
		if err != nil {
			fmt.Printf("Error saving blockchain: %v\n", err)
			return
		}
		fmt.Printf("Funding complete. Wallets now have funds.\n\n")
	}

	// Generate random transactions
	fmt.Printf("Generating %d random transactions...\n", count)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	successCount := 0
	failCount := 0

	for i := 0; i < count; i++ {
		// Pick random sender and receiver (must be different)
		senderIdx := rand.Intn(walletCount)
		receiverIdx := rand.Intn(walletCount)
		for receiverIdx == senderIdx {
			receiverIdx = rand.Intn(walletCount)
		}

		sender := wallets[senderIdx]
		receiver := wallets[receiverIdx]

		// Random amount between min and max
		amount := minAmount + rand.Float64()*(maxAmount-minAmount)

		// Create transaction
		tx, err := sender.CreateTransaction(
			receiver.Address,
			amount,
			fee,
			blockchain.UTXOSet,
			fmt.Sprintf("Test transaction #%d", i+1),
		)

		if err != nil {
			failCount++
			if fundWallets {
				// Only show detailed errors if we funded wallets (unexpected failures)
				fmt.Printf("Transaction %d failed: %s\n", i+1, err.Error())
			}
			continue
		}

		// Sign and add transaction
		sender.SignTransaction(tx)
		blockchain.AddTransaction(tx)
		successCount++

		if (i+1)%10 == 0 || i == count-1 {
			fmt.Printf("Progress: %d/%d transactions created (%d successful, %d failed)\n",
				i+1, count, successCount, failCount)
		}
	}

	// Save blockchain with all pending transactions
	err := blockchain.SaveToFile(blockchainFile)
	if err != nil {
		fmt.Printf("Error saving blockchain: %v\n", err)
		return
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Printf("✓ Transaction generation completed!\n")
	fmt.Printf("Total transactions created: %d\n", successCount)
	fmt.Printf("Failed transactions: %d\n", failCount)
	fmt.Printf("Pending transactions in mempool: %d\n", len(blockchain.Mempool.PendingTransactions))

	if successCount > 0 {
		fmt.Printf("\nTo mine all transactions, run:\n")
		fmt.Printf("go-coin blockchain run --miner \"%s\"\n", wallets[0].Address)
	}
}
