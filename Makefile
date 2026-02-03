# Makefile for things-cli
# Provides targets for building, testing, and installing the things command-line tool

.PHONY: help build install clean test lint

# Variables
BINARY_NAME = things
VERSION ?= 1.0.0
GO_LDFLAGS = -ldflags="-X main.Version=$(VERSION)"
INSTALL_PATH = /usr/local/bin
BUILD_DIR = bin

# Color output for readability
COLOR_RESET = \033[0m
COLOR_BLUE = \033[34m
COLOR_GREEN = \033[32m

help: ## Display this help message
	@echo "$(COLOR_BLUE)things-cli - Makefile targets$(COLOR_RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "Example usage:"
	@echo "  make build        # Build the binary"
	@echo "  make install      # Install to /usr/local/bin"
	@echo "  make clean        # Remove build artifacts"

build: ## Compile the things executable
	@echo "$(COLOR_BLUE)Building things CLI...$(COLOR_RESET)"
	mkdir -p $(BUILD_DIR)
	go build $(GO_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "$(COLOR_GREEN)✓ Build complete: ./$(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

install: build ## Build and install to /usr/local/bin/things
	@echo "$(COLOR_BLUE)Installing things to $(INSTALL_PATH)/$(BINARY_NAME)...$(COLOR_RESET)"
	# Check if we need sudo (if /usr/local/bin is not writable by current user)
	@if [ -w $(INSTALL_PATH) ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
		chmod +x $(INSTALL_PATH)/$(BINARY_NAME); \
		echo "$(COLOR_GREEN)✓ Installation complete$(COLOR_RESET)"; \
		echo "$(COLOR_GREEN)✓ Run 'things --help' to get started$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_BLUE)Requires sudo to install to $(INSTALL_PATH)$(COLOR_RESET)"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
		sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME); \
		echo "$(COLOR_GREEN)✓ Installation complete$(COLOR_RESET)"; \
		echo "$(COLOR_GREEN)✓ Run 'things --help' to get started$(COLOR_RESET)"; \
	fi

install-user: build ## Install to ~/.local/bin/things (no sudo required)
	@echo "$(COLOR_BLUE)Installing things to ~/.local/bin/$(BINARY_NAME)...$(COLOR_RESET)"
	mkdir -p ~/.local/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	chmod +x ~/.local/bin/$(BINARY_NAME)
	@echo "$(COLOR_GREEN)✓ Installation complete$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)Note: Make sure ~/.local/bin is in your PATH$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)✓ Run 'things --help' to get started$(COLOR_RESET)"

clean: ## Remove build artifacts
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)/
	go clean
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

test: ## Run Go tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	go test -v ./...

lint: ## Run golangci-lint if available
	@if command -v golangci-lint > /dev/null; then \
		echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"; \
		golangci-lint run ./...; \
	else \
		echo "$(COLOR_BLUE)golangci-lint not found, running go vet instead...$(COLOR_RESET)"; \
		go vet ./...; \
	fi

fmt: ## Format code with go fmt
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	go fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

deps: ## Download and verify dependencies
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	go mod download
	go mod verify
	@echo "$(COLOR_GREEN)✓ Dependencies verified$(COLOR_RESET)"

dev: ## Build in development mode (with verbose output)
	@echo "$(COLOR_BLUE)Building in development mode...$(COLOR_RESET)"
	mkdir -p $(BUILD_DIR)
	go build -v -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "$(COLOR_GREEN)✓ Development build complete: ./$(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

run: build ## Build and run the things command with help
	./$(BUILD_DIR)/$(BINARY_NAME) --help

version: ## Show version information
	@echo "things-cli version $(VERSION)"
	@go version

all: clean deps build ## Run clean, deps, and build

.DEFAULT_GOAL := help
