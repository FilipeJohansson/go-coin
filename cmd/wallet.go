package cmd

import (
	"fmt"
	"os"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/FilipeJohansson/go-coin/pkg/common"
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:     "wallet",
	Aliases: []string{"w"},
	Short:   "Wallet operations",
	Long:    "Create and manage wallets",
}

var createWalletCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a new wallet",
	Long:    `Create a new public and private key`,
	Run:     createWallet,
}

var loadWalletCmd = &cobra.Command{
	Use:     "load",
	Aliases: []string{"l"},
	Short:   "Load a wallet",
	Long:    `Load a wallet passing the public and the private key`,
	Run:     loadWallet,
}

var balanceCmd = &cobra.Command{
	Use:     "balance",
	Aliases: []string{"b"},
	Short:   "See the balance from a wallet",
	Run:     getWalletBalance,
}

func init() {
	createWalletCmd.Flags().StringP("name", "n", "", "Name your wallet")
	createWalletCmd.Flags().BoolP("save", "s", false, "Save the wallet in a file")

	loadWalletCmd.Flags().StringP("private-key", "p", "", "Your wallet private key")

	balanceCmd.Flags().StringP("address", "a", "", "Wallet address to check balance")

	walletCmd.AddCommand(createWalletCmd)
	walletCmd.AddCommand(loadWalletCmd)
	walletCmd.AddCommand(balanceCmd)

	rootCmd.AddCommand(walletCmd)
}

func createWallet(cmd *cobra.Command, args []string) {
	save, _ := cmd.Flags().GetBool("save")
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		// err
		return
	}

	wallet := wallet.NewWallet()

	content := fmt.Sprintf("Wallet Name: %s\n%s", name, wallet.Print())
	fmt.Println(content)

	if save {
		file, err := os.OpenFile("wallet.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(content + "\n\n")
		if err != nil {
			panic(err)
		}

		fmt.Println("Wallet saved to wallet.txt")
	}
}

func loadWallet(cmd *cobra.Command, args []string) {
	privateKey, err := cmd.Flags().GetString("private-key")
	if err != nil {
		// err
		return
	}

	wallet := wallet.LoadWallet(privateKey)

	fmt.Printf("Wallet loaded:\n%s", wallet.Print())
}

func getWalletBalance(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")

	if address == "" {
		fmt.Println("Error: address is required")
		return
	}

	blockchain := blockchain.NewBlockchain("", blockchainFile)
	fmt.Printf("Wallet balance: %.2f", (float64(blockchain.UTXOSet.GetAddressBalance(address)) / common.COINS_PER_UNIT))
}
