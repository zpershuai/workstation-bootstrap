GO := /usr/local/go/bin/go
BINARY_NAME := dwell
BUILD_DIR := bin
INSTALL_DIR := $(HOME)/.local/bin
VERSION := 0.1.0
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all build clean test install uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./internal/cmd

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

test:
	$(GO) test -v ./...

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installed! Make sure ~/.local/bin is in your PATH:"
	@echo "  export PATH=$(INSTALL_DIR):\$$PATH"
	@echo "Or add to your shell config (.bashrc/.zshrc):"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""

uninstall:
	@echo "Removing $(INSTALL_DIR)/$(BINARY_NAME)..."
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

dev:
	$(GO) run $(LDFLAGS) ./internal/cmd

init-config: build
	@./$(BUILD_DIR)/$(BINARY_NAME) init
