package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var blockchainFile string

var rootCmd = &cobra.Command{
	Use:   "go-coin",
	Short: "A simple blockchain implementation in Go",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&blockchainFile, "blockchain-file", "f", "blockchain.json", "Path to blockchain file")
}
