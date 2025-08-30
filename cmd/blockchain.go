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

var mineCmd = &cobra.Command{
	Use:     "mine",
	Aliases: []string{"m"},
	Short:   "Mine a new block",
	Run:     mineBlock,
}

var blocksCmd = &cobra.Command{
	Use:     "blocks",
	Aliases: []string{"b"},
	Short:   "List all the blocks",
	Run:     listBlocks,
}

func init() {
	mineCmd.Flags().StringP("miner", "m", "", "Wallet address to receive coinbase")

	blockchainCmd.AddCommand(mineCmd)
	blockchainCmd.AddCommand(validateCmd)
	blockchainCmd.AddCommand(blocksCmd)

	rootCmd.AddCommand(blockchainCmd)
}

func validateBlockchain(cmd *cobra.Command, args []string) {
	blockchain := blockchain.NewBlockchain("", blockchainFile)
	fmt.Printf("Is Blockchain valid: %t", blockchain.IsBlockchainValid())
}

func mineBlock(cmd *cobra.Command, args []string) {
	minerAddress, _ := cmd.Flags().GetString("miner")
	if minerAddress == "" {
		// err
		return
	}

	blockchain := blockchain.NewBlockchain("", blockchainFile)
	if len(blockchain.Mempool.PendingTransactions) < 1 {
		fmt.Println("No transactions pending")
		return
	}

	blockchain.MineBlock(minerAddress)
	blockchain.SaveToFile(blockchainFile)
}

func listBlocks(cmd *cobra.Command, args []string) {
	blockchain := blockchain.NewBlockchain("", blockchainFile)
	fmt.Println(blockchain.Print())
}
