# Arktree CLI Makefile

# Variables
BINARY_NAME=arktree
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go
	@echo "Build complete: ./${BINARY_NAME}"

# Install the binary to /usr/local/bin
.PHONY: install
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	sudo cp ${BINARY_NAME} /usr/local/bin/
	sudo chmod +x /usr/local/bin/${BINARY_NAME}
	@echo "Installation complete! You can now run '${BINARY_NAME}' from anywhere."

# Install to user's home directory (no sudo required)
.PHONY: install-user
install-user: build
	@echo "Installing ${BINARY_NAME} to ~/.local/bin..."
	mkdir -p ~/.local/bin
	cp ${BINARY_NAME} ~/.local/bin/
	chmod +x ~/.local/bin/${BINARY_NAME}
	@echo "Installation complete!"
	@echo "Add 'export PATH=~/.local/bin:\$$PATH' to your shell profile if not already done."

# Uninstall the binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	sudo rm -f /usr/local/bin/${BINARY_NAME}
	rm -f ~/.local/bin/${BINARY_NAME}
	@echo "Uninstallation complete!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f ${BINARY_NAME}
	@echo "Clean complete!"

# Run the application directly
.PHONY: run
run:
	go run main.go

# Run with specific command
.PHONY: run-generate
run-generate:
	go run main.go generate 5

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  install      - Install to /usr/local/bin (requires sudo)"
	@echo "  install-user - Install to ~/.local/bin (no sudo required)"
	@echo "  uninstall    - Remove installed binary"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Run the application directly"
	@echo "  run-generate - Run with generate command (5 leaves)"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make install      # Install system-wide"
	@echo "  make install-user # Install for current user"
	@echo "  make run-generate # Test the application"

# Test the build
.PHONY: test-build
test-build: build
	@echo "Testing the build..."
	./${BINARY_NAME} --help
	@echo "Build test successful!"

# Show version info
.PHONY: version
version:
	@echo "Version: ${VERSION}"
	@echo "Build Time: ${BUILD_TIME}" 