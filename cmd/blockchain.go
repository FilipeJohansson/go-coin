package cmd

import (
	"fmt"
	"strings"
	"time"

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

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Mine all pending transactions continuously",
	Long:    "Continuously mine blocks until all pending transactions are processed",
	Run:     runContinuousMining,
}

func init() {
	mineCmd.Flags().StringP("miner", "m", "", "Wallet address to receive coinbase")

	runCmd.Flags().StringP("miner", "m", "", "Wallet address to receive coinbase rewards")
	runCmd.Flags().BoolP("verbose", "v", false, "Show detailed mining progress")
	runCmd.Flags().IntP("delay", "d", 0, "Delay in seconds between blocks (default: 0)")

	blockchainCmd.AddCommand(mineCmd)
	blockchainCmd.AddCommand(validateCmd)
	blockchainCmd.AddCommand(blocksCmd)
	blockchainCmd.AddCommand(runCmd)

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

func runContinuousMining(cmd *cobra.Command, args []string) {
	minerAddress, _ := cmd.Flags().GetString("miner")
	if minerAddress == "" {
		fmt.Println("Error: miner address is required")
		return
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	delay, _ := cmd.Flags().GetInt("delay")

	blockchain := blockchain.NewBlockchain("", blockchainFile)

	totalTransactions := len(blockchain.Mempool.PendingTransactions)
	if totalTransactions == 0 {
		fmt.Println("No pending transactions to mine")
		return
	}

	fmt.Printf("Starting continuous mining with %d pending transactions...\n", totalTransactions)
	fmt.Printf("Miner address: %s\n", minerAddress)
	if delay > 0 {
		fmt.Printf("Block delay: %d seconds\n", delay)
	}
	fmt.Println(strings.Repeat("=", 50))

	blocksMinedCount := 0
	totalTransactionsProcessed := 0

	for {
		pendingCount := len(blockchain.Mempool.PendingTransactions)
		if pendingCount == 0 {
			break
		}

		if verbose {
			fmt.Printf("\n[Block %d] Mining block with %d pending transactions...\n",
				blocksMinedCount+1, pendingCount)
		}

		// Store count before mining to calculate processed transactions
		transactionsBeforeMining := pendingCount

		// Mine the block
		blockchain.MineBlock(minerAddress)

		// Calculate how many transactions were processed
		transactionsAfterMining := len(blockchain.Mempool.PendingTransactions)
		transactionsProcessedThisBlock := transactionsBeforeMining - transactionsAfterMining
		totalTransactionsProcessed += transactionsProcessedThisBlock

		blocksMinedCount++

		if verbose {
			fmt.Printf("✓ Block mined successfully!")
			fmt.Printf(" Processed %d transactions\n", transactionsProcessedThisBlock)
			fmt.Printf("  Remaining: %d transactions\n", transactionsAfterMining)
		} else {
			// Show progress without verbose details
			fmt.Printf("Block %d mined: %d/%d transactions processed\n",
				blocksMinedCount, totalTransactionsProcessed, totalTransactions)
		}

		// Save blockchain after each block
		err := blockchain.SaveToFile(blockchainFile)
		if err != nil {
			fmt.Printf("Warning: Failed to save blockchain: %v\n", err)
		}

		// Add delay if specified
		if delay > 0 && len(blockchain.Mempool.PendingTransactions) > 0 {
			if verbose {
				fmt.Printf("Waiting %d seconds before next block...\n", delay)
			}
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✓ Mining completed!\n")
	fmt.Printf("Blocks mined: %d\n", blocksMinedCount)
	fmt.Printf("Total transactions processed: %d\n", totalTransactionsProcessed)
	fmt.Printf("Blockchain saved to: %s\n", blockchainFile)
}
