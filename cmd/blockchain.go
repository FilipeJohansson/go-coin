package cmd

import (
	"fmt"

	"github.com/FilipeJohansson/go-coin/internal/blockchain"
	"github.com/spf13/cobra"
)

var blockchainCmd = &cobra.Command{
	Use:     "blockchain",
	Aliases: []string{"bc"},
	Short:   "Blockchain operations",
}

var validateCmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"v"},
	Short:   "Check if the blockchain is valid",
	Run:     validateBlockchain,
}

func init() {
	blockchainCmd.AddCommand(validateCmd)

	rootCmd.AddCommand(blockchainCmd)
}

func validateBlockchain(cmd *cobra.Command, args []string) {
	blockchain := blockchain.NewBlockchain("", blockchainFile)
	fmt.Printf("Is Blockchain valid: %t", blockchain.IsBlockchainValid())
}
