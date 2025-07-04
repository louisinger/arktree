package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ark-network/ark/common"
	"github.com/ark-network/ark/common/tree"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arktree",
	Short: "A CLI tool for generating Ark trees",
	Long:  `Arktree is a command-line tool for generating and working with Ark trees.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate [number-of-leaves]",
	Short: "Generate an Ark tree with the specified number of leaves",
	Long:  `Generate an Ark tree with the specified number of leaves. The number of leaves must be a positive integer.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		numLeaves, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid number of leaves: %s\n", args[0])
			os.Exit(1)
		}

		if numLeaves <= 0 {
			fmt.Println("Error: Number of leaves must be a positive integer")
			os.Exit(1)
		}

		// Print header with styling
		fmt.Println("🌳 Ark Tree Generator")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("📊 Generating Ark tree with %d leaves...\n\n", numLeaves)

		// Generate random data
		fmt.Print("🔧 Initializing random data... ")
		randomSweepTreeRoot := make([]byte, 32)
		rand.Read(randomSweepTreeRoot)

		randomTxid := make([]byte, 32)
		rand.Read(randomTxid)
		fmt.Println("✅")

		// Generate leaves
		fmt.Printf("🍃 Generating %d leaves... ", numLeaves)
		leaves := make([]tree.Leaf, numLeaves)

		for i := 0; i < numLeaves; i++ {
			randomScript := make([]byte, 34)
			rand.Read(randomScript)

			randomPrivkey, err := secp256k1.GeneratePrivateKey()
			if err != nil {
				fmt.Printf("\n❌ Error: Failed to generate private key: %s\n", err)
				os.Exit(1)
			}

			randomPubkey := randomPrivkey.PubKey()

			leaves[i] = tree.Leaf{
				Amount:              1000,
				Script:              hex.EncodeToString(randomScript),
				CosignersPublicKeys: []string{hex.EncodeToString(randomPubkey.SerializeCompressed())},
			}
		}
		fmt.Println("✅")

		// Build tree
		fmt.Print("🌿 Building Vtxo tree... ")
		start := time.Now()
		txtree, err := tree.BuildVtxoTree(
			&wire.OutPoint{
				Hash:  chainhash.Hash(randomTxid),
				Index: 0,
			},
			leaves,
			randomSweepTreeRoot,
			common.RelativeLocktime{Value: 100, Type: common.LocktimeTypeBlock},
		)
		if err != nil {
			fmt.Printf("\n❌ Error: Failed to build tree: %s\n", err)
			os.Exit(1)
		}
		elapsed := time.Since(start)
		fmt.Printf("✅ (%s)\n", elapsed)

		// Calculate statistics
		fmt.Print("📈 Calculating tree statistics... ")
		totalSize, err := numberOfNodes(txtree)
		if err != nil {
			fmt.Printf("\n❌ Error: Failed to get total size: %s\n", err)
			os.Exit(1)
		}

		branchSizes, err := sizeOfBranches(txtree)
		if err != nil {
			fmt.Printf("\n❌ Error: Failed to get size of branches: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("✅")

		// Find biggest branch
		biggestBranch := 0
		for _, size := range branchSizes {
			if size > biggestBranch {
				biggestBranch = size
			}
		}

		// Print results with beautiful formatting
		fmt.Println("\n" + strings.Repeat("─", 60))
		fmt.Println("📊 TREE STATISTICS")
		fmt.Println(strings.Repeat("─", 60))

		fmt.Printf("🌳 Total Transactions:    %8d\n", totalSize)
		fmt.Printf("🍃 Number of Leaves:      %8d\n", numLeaves)
		fmt.Printf("📏 Biggest Branch Size:   %8d tx\n", biggestBranch)

		// Calculate average and median branch size
		if len(branchSizes) > 0 {
			avgSize := calculateAverage(branchSizes)
			fmt.Printf("📊 Average Branch Size:   %8.1f tx\n", avgSize)

			// Calculate median
			medianSize := calculateMedian(branchSizes)
			fmt.Printf("📊 Median Branch Size:    %8.1f tx\n", medianSize)
		}

		fmt.Println(strings.Repeat("─", 60))

		// Group branches by size
		sizeCount := make(map[int]int)
		for _, size := range branchSizes {
			sizeCount[size]++
		}

		// Print branch details grouped by size
		fmt.Println("\n🌿 BRANCH DETAILS:")
		fmt.Println(strings.Repeat("─", 40))

		// Sort sizes for consistent output
		var sizes []int
		for size := range sizeCount {
			sizes = append(sizes, size)
		}

		// Simple sort (bubble sort for small arrays)
		for i := 0; i < len(sizes)-1; i++ {
			for j := 0; j < len(sizes)-i-1; j++ {
				if sizes[j] > sizes[j+1] {
					sizes[j], sizes[j+1] = sizes[j+1], sizes[j]
				}
			}
		}

		for _, size := range sizes {
			count := sizeCount[size]
			if count == 1 {
				fmt.Printf("%2d branch  with %2d tx\n", count, size)
			} else {
				fmt.Printf("%2d branches with %2d tx\n", count, size)
			}
		}

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Printf("🎉 Successfully generated Ark tree with %d leaves!\n", numLeaves)
		fmt.Println(strings.Repeat("=", 60))

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func sizeOfBranches(g *tree.TxGraph) ([]int, error) {
	leaves := g.Leaves()

	branchSizes := make([]int, 0, len(leaves))

	for _, leaf := range leaves {
		branch, err := g.SubGraph([]string{leaf.UnsignedTx.TxID()})
		if err != nil {
			return nil, err
		}

		count, err := numberOfNodes(branch)
		if err != nil {
			return nil, err
		}

		branchSizes = append(branchSizes, count)
	}

	return branchSizes, nil
}

func numberOfNodes(g *tree.TxGraph) (int, error) {
	count := 0
	if err := g.Apply(func(tx *tree.TxGraph) (bool, error) {
		count++
		return true, nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

func calculateAverage(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0
	for _, value := range values {
		sum += value
	}
	return float64(sum) / float64(len(values))
}

func calculateMedian(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]int, len(values))
	copy(sorted, values)

	// Sort the values
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Calculate median
	n := len(sorted)
	if n%2 == 0 {
		// Even number of elements - average of two middle values
		return float64(sorted[n/2-1]+sorted[n/2]) / 2.0
	} else {
		// Odd number of elements - middle value
		return float64(sorted[n/2])
	}
}
