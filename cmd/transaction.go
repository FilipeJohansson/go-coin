package cmd

import (
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
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

func init() {
	sendCmd.Flags().StringP("to", "t", "", "Recipient address")
	sendCmd.Flags().StringP("private-key", "p", "", "The from address private key to autenticate")
	sendCmd.Flags().Float64P("amount", "a", 0, "Quantity to send from sender to recipient")
	sendCmd.Flags().Float64("fee", float64(common.MIN_FEE)/float64(common.COINS_PER_UNIT), "Optional miners fee")
	sendCmd.Flags().StringP("message", "m", "", "Optional message")

	transactionCmd.AddCommand(sendCmd)
	transactionCmd.AddCommand(listCmd)

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
