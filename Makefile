GO := /usr/local/go/bin/go
BINARY_NAME := dwell
BUILD_DIR := bin
LOCAL_DIR := local
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
	rm -rf $(LOCAL_DIR)

test:
	$(GO) test -v ./...

install: build
	@echo "Installing to $(LOCAL_DIR)/bin..."
	@mkdir -p $(LOCAL_DIR)/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(LOCAL_DIR)/bin/
	@echo "Installed! Add $(LOCAL_DIR)/bin to your PATH:"
	@echo "  export PATH=$(PWD)/$(LOCAL_DIR)/bin:\$$PATH"
	@echo "Or run: ./$(LOCAL_DIR)/bin/$(BINARY_NAME) --help"

uninstall:
	@echo "Removing $(LOCAL_DIR)..."
	rm -rf $(LOCAL_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

dev:
	$(GO) run $(LDFLAGS) ./internal/cmd

init-config: build
	@./$(BUILD_DIR)/$(BINARY_NAME) init
