# gx - Go Version Manager
# Build and Release Makefile

# Version management
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build configuration
BINARY_NAME := gx
MAIN_PACKAGE := ./cmd/gx
BUILD_DIR := build
DIST_DIR := dist

# Go build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -s -w"
GO_BUILD := go build $(LDFLAGS)

# Supported platforms
PLATFORMS := \
	windows/amd64 \
	windows/386 \
	linux/amd64 \
	linux/386 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64

.PHONY: all
all: clean build

.PHONY: help
help:
	@echo "gx Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build for current platform"
	@echo "  make build-all      - Build for all supported platforms"
	@echo "  make release        - Create release packages for all platforms"
	@echo "  make install        - Install gx to GOPATH/bin"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make version        - Show version information"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"

.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) \
		$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME)$(shell GOOS=$(shell go env GOHOSTOS) go env GOEXE) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(shell GOOS=$(shell go env GOHOSTOS) go env GOEXE)"

.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@$(MAKE) $(foreach platform,$(PLATFORMS),build-$(subst /,-,$(platform)))

# Generate build targets for each platform
define build-platform
.PHONY: build-$(subst /,-,$(1))
build-$(subst /,-,$(1)):
	@echo "Building for $(1)..."
	@mkdir -p $(BUILD_DIR)/$(subst /,_,$(1))
	GOOS=$(word 1,$(subst /, ,$(1))) GOARCH=$(word 2,$(subst /, ,$(1))) \
		$(GO_BUILD) -o $(BUILD_DIR)/$(subst /,_,$(1))/$(BINARY_NAME)$(if $(findstring windows,$(1)),.exe,) $(MAIN_PACKAGE)
endef

$(foreach platform,$(PLATFORMS),$(eval $(call build-platform,$(platform))))

.PHONY: release
release: clean build-all
	@echo "Creating release packages..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		platform_dir=$$(echo $$platform | tr '/' '_'); \
		archive_name="$(BINARY_NAME)-$(VERSION)-$$os-$$arch"; \
		echo "Packaging $$archive_name..."; \
		if [ "$$os" = "windows" ]; then \
			cd $(BUILD_DIR)/$$platform_dir && \
			zip -q ../../$(DIST_DIR)/$$archive_name.zip $(BINARY_NAME).exe && \
			cd ../..; \
		else \
			tar -czf $(DIST_DIR)/$$archive_name.tar.gz -C $(BUILD_DIR)/$$platform_dir $(BINARY_NAME); \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/

.PHONY: release-checksums
release-checksums:
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "Checksums saved to $(DIST_DIR)/checksums.txt"

.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PACKAGE)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Clean complete"

.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Test complete"

.PHONY: test-coverage
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	@echo "Dependencies ready"

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi
