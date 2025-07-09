# 🌳 Arktree CLI

A command-line tool for generating and analyzing Ark trees with comprehensive statistics.

## 🚀 Quick Start

```bash
# Install dependencies
go mod download

# Run directly
go run main.go generate 5

# Or build and run
go build -o arktree main.go
./arktree generate 5
```

## 🔧 Installation

```bash
# Install system-wide (requires sudo)
make install

# Install for current user only (no sudo)
make install-user

# See all available commands
make help
```

## 📖 Usage

```bash
# Show help
go run main.go --help

# Generate a tree with N leaves
go run main.go generate 5
go run main.go generate 100

# Generate charts across a range of leaf counts
go run main.go chart 2 10 2    # From 2 to 10 leaves, step by 2
go run main.go chart 5 25 5    # From 5 to 25 leaves, step by 5

# Export charts to CSV for further analysis
go run main.go chart 2 10 2 --csv    # Generate chart and export to CSV



# Show command help
go run main.go generate --help
go run main.go chart --help
```

## 📊 Features

The Arktree CLI generates Ark trees and provides detailed statistics including:

- **Tree Structure**: Total transactions and number of leaves
- **Branch Analysis**: Size distribution of branches (how many transactions each branch contains)
- **Broadcast Analysis**: Weight distribution showing how many transactions each user needs to broadcast
- **Performance Metrics**: Average, median, and distribution statistics for both size and broadcast weight
- **Comparative Analysis**: Chart generation across different leaf counts to analyze scaling characteristics
- **CSV Export**: Export chart results to CSV files for further analysis in spreadsheet applications


### Example Output

```
🌳 Ark Tree Generator
===================================================
📊 Generating Ark tree with 5 leaves...

🔧 Initializing random data... ✅
🍃 Generating 5 leaves... ✅
🌿 Building Vtxo tree... ✅ (683.75µs)
📈 Calculating tree statistics... ✅

────────────────────────────────────────────────────────────
📊 TREE STATISTICS
────────────────────────────────────────────────────────────
🌳 Total Transactions:           9
🍃 Number of Leaves:             5
📏 Biggest Branch Size:          4 tx
📊 Average Branch Size:        1.8 tx
📊 Median Branch Size:         4.0 tx
📡 Most Tx to Broadcast:        1.83
📊 Avg Tx to Broadcast:        1.67
📊 Median Tx to Broadcast:     1.83
────────────────────────────────────────────────────────────

🌿 BRANCH SIZE DETAILS:
────────────────────────────────────────
 1 branch  with  2 tx
 4 branches with  4 tx

📡 BROADCAST WEIGHT DETAILS:
────────────────────────────────────────
 1 branch  with 1.33 tx to broadcast
 4 branches with 1.83 tx to broadcast

============================================================
🎉 Successfully generated Ark tree with 5 leaves!
============================================================

### Chart Command Example

```
📊 Ark Tree Chart Generator
===================================================
📈 Generating charts from 2 to 6 leaves (step: 2)...

🔧 Processing 2 leaves... ✅
🔧 Processing 4 leaves... ✅
🔧 Processing 6 leaves... ✅

================================================================================
📊 CHART RESULTS
================================================================================
Leaves   | Total Tx | Biggest  | Avg Size | Median   | Most Tx  | Avg Tx   | Median Tx
         |          | Branch   |          | Size     | Broadcast | Broadcast | Broadcast
────────────────────────────────────────────────────────────────────────────────
2        | 3        | 2        | 2.0      | 2.0      | 1.50     | 1.50     | 1.50    
4        | 7        | 3        | 3.0      | 3.0      | 1.75     | 1.75     | 1.75    
6        | 11       | 4        | 3.7      | 4.0      | 1.92     | 1.83     | 1.92    
────────────────────────────────────────────────────────────────────────────────

📈 SUMMARY STATISTICS:
────────────────────────────────────────
Total Transactions: 3 → 11 (range: 8)
Biggest Branch:     2 → 4 (range: 2)
Most Tx Broadcast:  1.50 → 1.92 (range: 0.42)

================================================================================
🎉 Chart generation complete! Analyzed 3 different leaf counts.
================================================================================
```

## 🔍 Understanding the Statistics

### Branch Size
- **Biggest Branch Size**: The maximum number of transactions in any single branch
- **Average/Median Branch Size**: Statistical measures of branch transaction counts
- **Branch Size Details**: Distribution showing how many branches have each transaction count

### Broadcast Weight
- **Most Tx to Broadcast**: The maximum number of transactions any user needs to broadcast
- **Avg/Median Tx to Broadcast**: Statistical measures of broadcast burden per user
- **Broadcast Weight Details**: Distribution showing how many branches require each broadcast count

The broadcast weight represents the computational and network burden on each cosigner. For example, if a transaction is shared by 3 cosigners, each cosigner broadcasts 1/3 of the transaction (weight = 1/3).

## 🛠️ Development

```bash
# Build the project
go build

# Run tests
go test

# Format code
go fmt

# Run linter
golangci-lint run
```

## 📦 Dependencies

- [ark-network/ark](https://github.com/ark-network/ark) - Core Ark tree functionality
- [btcsuite/btcd](https://github.com/btcsuite/btcd) - Bitcoin protocol implementation
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [decred/dcrd](https://github.com/decred/dcrd) - Cryptographic primitives 