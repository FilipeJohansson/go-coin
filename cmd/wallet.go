package cmd

import (
	"fmt"
	"os"

	"github.com/FilipeJohansson/go-coin/internal/wallet"
	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet operations",
	Long:  "Create and manage wallets",
}

var createWalletCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Long:  `Create a new public and private key`,
	Run:   createWallet,
}

func init() {
	walletCmd.AddCommand(createWalletCmd)
	createWalletCmd.Flags().StringP("name", "n", "", "Name your wallet")
	createWalletCmd.Flags().BoolP("save", "s", false, "Save the wallet in a file")

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
