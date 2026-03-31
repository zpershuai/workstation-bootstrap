GO := /usr/local/go/bin/go
BINARY_NAME := dwell
BUILD_DIR := bin
VERSION := 0.1.0
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all build clean test install

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
	@echo "Installing to /usr/local/bin..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed! Run 'dwell --help' to get started."

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

dev:
	$(GO) run $(LDFLAGS) ./internal/cmd
