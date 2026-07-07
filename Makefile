BIN := bin/vivio
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)

# Release configuration
DIST_DIR := dist
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: fmt lint test coverage pre-commit clean clean-dist build install run release release-build release-checksums

## Format Go code
fmt:
	go fmt ./...

## Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin" && exit 1)
	golangci-lint run ./...

## Run tests for all packages with coverage
test:
	go test -v -coverprofile=cover.out ./...

## Generate HTML coverage report
coverage:
	@echo "Generating coverage report..."
	go tool cover -html=cover.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Run all pre-commit checks
pre-commit: fmt lint test
	@echo "✓ All pre-commit checks passed"

## Remove build artifacts
clean:
	rm -rf bin/
	rm -rf $(DIST_DIR)/
	rm -f cover.out coverage.html

## Remove binary artifacts
clean-bin:
	rm -rf bin/

## Remove release artifacts
clean-dist:
	rm -rf $(DIST_DIR)/

## Build the CLI binary with version info
build: clean-bin
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd
	@echo "Built: $(BIN)"
	@echo "Version: $(VERSION) (commit: $(COMMIT))"

## Install the binary to system (requires sudo on Linux/Mac)
install:
	@echo "Installing vivio to /usr/local/bin..."
	@sudo cp $(BIN) /usr/local/bin/vivio
	@echo "Installed: /usr/local/bin/vivio"

## Run the CLI directly (dev mode)
run:
	go run ./cmd $(ARGS)

## Build release binaries for all platforms
release: clean-dist release-build release-checksums
	@echo ""
	@echo "Release $(VERSION) built successfully!"
	@echo "Artifacts in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/

## Build binaries for all platforms
release-build:
	@echo "Building release $(VERSION) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=vivio-$(VERSION)-$${GOOS}-$${GOARCH}; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name=$${output_name}.exe; \
		fi; \
		echo "  Building $$output_name..."; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build \
			-ldflags "$(LDFLAGS)" \
			-o $(DIST_DIR)/$$output_name \
			./cmd || exit 1; \
	done
	@echo "  Compressing binaries..."
	@cd $(DIST_DIR) && for file in vivio-*; do \
		if [ -f "$$file" ]; then \
			tar czf "$${file}.tar.gz" "$$file" && rm "$$file"; \
		fi; \
	done

## Generate checksums for release artifacts
release-checksums:
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && sha256sum *.tar.gz > SHA256SUMS
	@echo "  ✓ SHA256SUMS created"

## Show help
help:
	@echo "Vivio CLI - Makefile targets:"
	@echo ""
	@echo "Development:"
	@echo "  make fmt         - Format Go code"
	@echo "  make lint        - Run linter (requires golangci-lint)"
	@echo "  make test        - Run all tests with coverage"
	@echo "  make coverage    - Generate HTML coverage report"
	@echo "  make pre-commit  - Run fmt + lint + test (use before committing)"
	@echo ""
	@echo "Build & Install:"
	@echo "  make build       - Build the CLI binary with version info"
	@echo "  make install     - Install binary to /usr/local/bin (requires sudo)"
	@echo "  make run         - Run CLI in dev mode (use ARGS='...' for arguments)"
	@echo ""
	@echo "Release:"
	@echo "  make release     - Build release for all platforms (linux/darwin/windows)"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean       - Remove all build artifacts and coverage files"
	@echo "  make clean-bin   - Remove binary artifacts only"
	@echo "  make clean-dist  - Remove release artifacts only"
	@echo ""
	@echo "Examples:"
	@echo "  make run ARGS='list channels --country=IN'"
	@echo "  make run ARGS='play 42'"
	@echo "  make pre-commit"
	@echo "  make release"
