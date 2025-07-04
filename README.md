# 🌳 Arktree CLI

A command-line tool for generating Ark trees.

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

# Show command help
go run main.go generate --help
```

### Example Output

```
🌳 Ark Tree Generator
===================================================
📊 Generating Ark tree with 5 leaves...

🔧 Initializing random data... ✅
🍃 Generating 5 leaves... ✅
🌿 Building Vtxo tree... ✅
📈 Calculating tree statistics... ✅

────────────────────────────────────────────────────────────
📊 TREE STATISTICS
────────────────────────────────────────────────────────────
🌳 Total Transactions:           9
🍃 Number of Leaves:             5
📏 Biggest Branch Size:          4 tx
📊 Average Branch Size:        1.8 tx
📊 Median Branch Size:         4.0 tx
────────────────────────────────────────────────────────────

🌿 BRANCH DETAILS:
────────────────────────────────────────
 1 branch  with  2 tx
 4 branches with  4 tx

============================================================
🎉 Successfully generated Ark tree with 5 leaves!
============================================================
``` 