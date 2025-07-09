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

// TreeStatistics holds all the statistics for an Ark tree
type TreeStatistics struct {
	// Basic tree info
	TotalTransactions int
	NumberOfLeaves    int

	// Size statistics
	BiggestBranchSize      int
	AverageBranchSize      float64
	MedianBranchSize       float64
	BranchSizeDistribution map[int]int // size -> count

	// Broadcast weight statistics
	MostTxToBroadcast     float64
	AverageTxToBroadcast  float64
	MedianTxToBroadcast   float64
	BroadcastDistribution map[float64]int // weight -> count

	// Performance
	BuildTime time.Duration
}

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
		stats, err := calculateTreeStatistics(txtree, numLeaves, elapsed)
		if err != nil {
			fmt.Printf("\n❌ Error: Failed to calculate statistics: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("✅")

		// Print results with beautiful formatting
		printTreeStatistics(stats)

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Printf("🎉 Successfully generated Ark tree with %d leaves!\n", numLeaves)
		fmt.Println(strings.Repeat("=", 60))

	},
}

var chartCmd = &cobra.Command{
	Use:   "chart [start-leaves] [end-leaves] [step]",
	Short: "Generate charts comparing Ark trees across a range of leaf counts",
	Long:  `Generate charts comparing Ark tree statistics across a range of leaf counts. Requires start, end, and step values.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		exportCSV, _ := cmd.Flags().GetBool("csv")
		startLeaves, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid start number of leaves: %s\n", args[0])
			os.Exit(1)
		}

		endLeaves, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: Invalid end number of leaves: %s\n", args[1])
			os.Exit(1)
		}

		step, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Error: Invalid step value: %s\n", args[2])
			os.Exit(1)
		}

		if startLeaves <= 0 || endLeaves <= 0 || step <= 0 {
			fmt.Println("Error: All values must be positive integers")
			os.Exit(1)
		}

		if startLeaves > endLeaves {
			fmt.Println("Error: Start value must be less than or equal to end value")
			os.Exit(1)
		}

		// Print header with styling
		fmt.Println("📊 Ark Tree Chart Generator")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("📈 Generating charts from %d to %d leaves (step: %d)...\n\n", startLeaves, endLeaves, step)

		// Generate statistics for each leaf count
		var allStats []*TreeStatistics
		var leafCounts []int

		for leaves := startLeaves; leaves <= endLeaves; leaves += step {
			fmt.Printf("🔧 Processing %d leaves... ", leaves)

			stats, err := generateTreeWithStats(leaves)
			if err != nil {
				fmt.Printf("❌ Error: %s\n", err)
				continue
			}

			allStats = append(allStats, stats)
			leafCounts = append(leafCounts, leaves)
			fmt.Println("✅")
		}

		if len(allStats) == 0 {
			fmt.Println("❌ No valid statistics generated")
			os.Exit(1)
		}

		// Print chart results
		printChartResults(leafCounts, allStats)

		// Export to CSV if requested
		if exportCSV {
			filename := fmt.Sprintf("arktree_chart_%d_to_%d_step_%d.csv", startLeaves, endLeaves, step)
			err := exportChartToCSV(filename, leafCounts, allStats)
			if err != nil {
				fmt.Printf("❌ Error exporting CSV: %s\n", err)
			} else {
				fmt.Printf("📊 CSV chart exported to: %s\n", filename)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	chartCmd.Flags().Bool("csv", false, "Export chart results to CSV file")
	rootCmd.AddCommand(chartCmd)
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

func weightOfBranches(g *tree.TxGraph) ([]float64, error) {
	leaves := g.Leaves()

	branchWeights := make([]float64, 0, len(leaves))

	for _, leaf := range leaves {
		branch, err := g.SubGraph([]string{leaf.UnsignedTx.TxID()})
		if err != nil {
			return nil, err
		}

		weight, err := computeBroadcastWeight(branch)
		if err != nil {
			return nil, err
		}

		branchWeights = append(branchWeights, weight)
	}

	return branchWeights, nil
}

func calculateAverageFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func calculateMedianFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]float64, len(values))
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
		return (sorted[n/2-1] + sorted[n/2]) / 2.0
	} else {
		// Odd number of elements - middle value
		return sorted[n/2]
	}
}

func calculateTreeStatistics(g *tree.TxGraph, numLeaves int, buildTime time.Duration) (*TreeStatistics, error) {
	// Calculate total size
	totalSize, err := numberOfNodes(g)
	if err != nil {
		return nil, err
	}

	// Calculate branch sizes
	branchSizes, err := sizeOfBranches(g)
	if err != nil {
		return nil, err
	}

	// Calculate branch weights
	branchWeights, err := weightOfBranches(g)
	if err != nil {
		return nil, err
	}

	// Find biggest branch
	biggestBranch := 0
	for _, size := range branchSizes {
		if size > biggestBranch {
			biggestBranch = size
		}
	}

	// Find heaviest branch
	heaviestBranch := 0.0
	for _, weight := range branchWeights {
		if weight > heaviestBranch {
			heaviestBranch = weight
		}
	}

	// Calculate averages and medians
	var avgSize, medianSize float64
	if len(branchSizes) > 0 {
		avgSize = calculateAverage(branchSizes)
		medianSize = calculateMedian(branchSizes)
	}

	var avgWeight, medianWeight float64
	if len(branchWeights) > 0 {
		avgWeight = calculateAverageFloat(branchWeights)
		medianWeight = calculateMedianFloat(branchWeights)
	}

	// Group branches by size
	sizeCount := make(map[int]int)
	for _, size := range branchSizes {
		sizeCount[size]++
	}

	// Group branches by weight (rounded to 2 decimal places)
	weightCount := make(map[float64]int)
	for _, weight := range branchWeights {
		roundedWeight := float64(int(weight*100)) / 100 // Round to 2 decimal places
		weightCount[roundedWeight]++
	}

	return &TreeStatistics{
		TotalTransactions:      totalSize,
		NumberOfLeaves:         numLeaves,
		BiggestBranchSize:      biggestBranch,
		AverageBranchSize:      avgSize,
		MedianBranchSize:       medianSize,
		BranchSizeDistribution: sizeCount,
		MostTxToBroadcast:      heaviestBranch,
		AverageTxToBroadcast:   avgWeight,
		MedianTxToBroadcast:    medianWeight,
		BroadcastDistribution:  weightCount,
		BuildTime:              buildTime,
	}, nil
}

func printTreeStatistics(stats *TreeStatistics) {
	// Print results with beautiful formatting
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Println("📊 TREE STATISTICS")
	fmt.Println(strings.Repeat("─", 60))

	fmt.Printf("🌳 Total Transactions:    %8d\n", stats.TotalTransactions)
	fmt.Printf("🍃 Number of Leaves:      %8d\n", stats.NumberOfLeaves)
	fmt.Printf("📏 Biggest Branch Size:   %8d tx\n", stats.BiggestBranchSize)

	if stats.AverageBranchSize > 0 {
		fmt.Printf("📊 Average Branch Size:   %8.1f tx\n", stats.AverageBranchSize)
		fmt.Printf("📊 Median Branch Size:    %8.1f tx\n", stats.MedianBranchSize)
	}

	fmt.Printf("📡 Most Tx to Broadcast:    %8.2f\n", stats.MostTxToBroadcast)

	if stats.AverageTxToBroadcast > 0 {
		fmt.Printf("📊 Avg Tx to Broadcast:    %8.2f\n", stats.AverageTxToBroadcast)
		fmt.Printf("📊 Median Tx to Broadcast: %8.2f\n", stats.MedianTxToBroadcast)
	}

	fmt.Println(strings.Repeat("─", 60))

	// Print branch size details
	fmt.Println("\n🌿 BRANCH SIZE DETAILS:")
	fmt.Println(strings.Repeat("─", 40))

	// Sort sizes for consistent output
	var sizes []int
	for size := range stats.BranchSizeDistribution {
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
		count := stats.BranchSizeDistribution[size]
		if count == 1 {
			fmt.Printf("%2d branch  with %2d tx\n", count, size)
		} else {
			fmt.Printf("%2d branches with %2d tx\n", count, size)
		}
	}

	// Print broadcast weight details
	fmt.Println("\n📡 BROADCAST WEIGHT DETAILS:")
	fmt.Println(strings.Repeat("─", 40))

	// Sort weights for consistent output
	var weights []float64
	for weight := range stats.BroadcastDistribution {
		weights = append(weights, weight)
	}

	// Simple sort (bubble sort for small arrays)
	for i := 0; i < len(weights)-1; i++ {
		for j := 0; j < len(weights)-i-1; j++ {
			if weights[j] > weights[j+1] {
				weights[j], weights[j+1] = weights[j+1], weights[j]
			}
		}
	}

	for _, weight := range weights {
		count := stats.BroadcastDistribution[weight]
		if count == 1 {
			fmt.Printf("%2d branch  with %.2f tx to broadcast\n", count, weight)
		} else {
			fmt.Printf("%2d branches with %.2f tx to broadcast\n", count, weight)
		}
	}
}

func generateTreeWithStats(numLeaves int) (*TreeStatistics, error) {
	// Generate random data
	randomSweepTreeRoot := make([]byte, 32)
	rand.Read(randomSweepTreeRoot)

	randomTxid := make([]byte, 32)
	rand.Read(randomTxid)

	// Generate leaves
	leaves := make([]tree.Leaf, numLeaves)

	for i := 0; i < numLeaves; i++ {
		randomScript := make([]byte, 34)
		rand.Read(randomScript)

		randomPrivkey, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate private key: %s", err)
		}

		randomPubkey := randomPrivkey.PubKey()

		leaves[i] = tree.Leaf{
			Amount:              1000,
			Script:              hex.EncodeToString(randomScript),
			CosignersPublicKeys: []string{hex.EncodeToString(randomPubkey.SerializeCompressed())},
		}
	}

	// Build tree
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
		return nil, fmt.Errorf("failed to build tree: %s", err)
	}
	elapsed := time.Since(start)

	// Calculate statistics
	return calculateTreeStatistics(txtree, numLeaves, elapsed)
}

func printChartResults(leafCounts []int, allStats []*TreeStatistics) {
	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Println("📊 CHART RESULTS")
	fmt.Println(strings.Repeat("=", 90))

	// Print header
	fmt.Printf("%-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s\n",
		"Leaves", "Total Tx", "Biggest", "Avg Size", "Median", "Most Tx", "Avg Tx", "Median Tx", "Build")
	fmt.Printf("%-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s\n",
		"", "", "Branch", "", "Size", "Broadcast", "Broadcast", "Broadcast", "Time(µs)")
	fmt.Println(strings.Repeat("─", 90))

	// Print data rows
	for i, stats := range allStats {
		fmt.Printf("%-8d | %-8d | %-8d | %-8.1f | %-8.1f | %-8.2f | %-8.2f | %-8.2f | %-8d\n",
			leafCounts[i],
			stats.TotalTransactions,
			stats.BiggestBranchSize,
			stats.AverageBranchSize,
			stats.MedianBranchSize,
			stats.MostTxToBroadcast,
			stats.AverageTxToBroadcast,
			stats.MedianTxToBroadcast,
			stats.BuildTime.Microseconds())
	}

	fmt.Println(strings.Repeat("─", 90))

	// Print summary statistics
	fmt.Println("\n📈 SUMMARY STATISTICS:")
	fmt.Println(strings.Repeat("─", 40))

	// Find min/max values
	minTotalTx, maxTotalTx := allStats[0].TotalTransactions, allStats[0].TotalTransactions
	minBiggest, maxBiggest := allStats[0].BiggestBranchSize, allStats[0].BiggestBranchSize
	minBroadcast, maxBroadcast := allStats[0].MostTxToBroadcast, allStats[0].MostTxToBroadcast

	for _, stats := range allStats {
		if stats.TotalTransactions < minTotalTx {
			minTotalTx = stats.TotalTransactions
		}
		if stats.TotalTransactions > maxTotalTx {
			maxTotalTx = stats.TotalTransactions
		}
		if stats.BiggestBranchSize < minBiggest {
			minBiggest = stats.BiggestBranchSize
		}
		if stats.BiggestBranchSize > maxBiggest {
			maxBiggest = stats.BiggestBranchSize
		}
		if stats.MostTxToBroadcast < minBroadcast {
			minBroadcast = stats.MostTxToBroadcast
		}
		if stats.MostTxToBroadcast > maxBroadcast {
			maxBroadcast = stats.MostTxToBroadcast
		}
	}

	fmt.Printf("Total Transactions: %d → %d (range: %d)\n", minTotalTx, maxTotalTx, maxTotalTx-minTotalTx)
	fmt.Printf("Biggest Branch:     %d → %d (range: %d)\n", minBiggest, maxBiggest, maxBiggest-minBiggest)
	fmt.Printf("Most Tx Broadcast:  %.2f → %.2f (range: %.2f)\n", minBroadcast, maxBroadcast, maxBroadcast-minBroadcast)

	fmt.Println("\n" + strings.Repeat("=", 90))
	fmt.Printf("🎉 Chart generation complete! Analyzed %d different leaf counts.\n", len(allStats))
	fmt.Println(strings.Repeat("=", 90))
}

func exportChartToCSV(filename string, leafCounts []int, allStats []*TreeStatistics) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %s", err)
	}
	defer file.Close()

	// Write CSV header
	header := "Leaves,Total_Transactions,Biggest_Branch_Size,Average_Branch_Size,Median_Branch_Size,Most_Tx_To_Broadcast,Average_Tx_To_Broadcast,Median_Tx_To_Broadcast,Build_Time_Microseconds\n"
	_, err = file.WriteString(header)
	if err != nil {
		return fmt.Errorf("failed to write CSV header: %s", err)
	}

	// Write data rows
	for i, stats := range allStats {
		row := fmt.Sprintf("%d,%d,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%d\n",
			leafCounts[i],
			stats.TotalTransactions,
			stats.BiggestBranchSize,
			stats.AverageBranchSize,
			stats.MedianBranchSize,
			stats.MostTxToBroadcast,
			stats.AverageTxToBroadcast,
			stats.MedianTxToBroadcast,
			stats.BuildTime.Microseconds())

		_, err = file.WriteString(row)
		if err != nil {
			return fmt.Errorf("failed to write CSV row: %s", err)
		}
	}

	return nil
}

// weight = the part of the tx a user has to broadcast
// if a tx is shared by 3 cosigners, each cosigner has to broadcast 1/3 of the tx
func computeBroadcastWeight(branch *tree.TxGraph) (float64, error) {
	var totalWeight float64
	if err := branch.Apply(func(g *tree.TxGraph) (bool, error) {
		cosignerKeys, err := tree.GetCosignerKeys(g.Root.Inputs[0])
		if err != nil {
			return false, err
		}

		totalWeight += 1 / float64(len(cosignerKeys))
		return true, nil
	}); err != nil {
		return 0, err
	}

	return totalWeight, nil
}
