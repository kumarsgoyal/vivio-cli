# Development Guide

**For contributors and developers** — build from source, create packages, contribute code.

---

## Prerequisites

### Install Go 1.22+

**Linux (Fedora/RHEL):**
```bash
sudo dnf install golang
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt install golang-go
```

**Linux (Arch):**
```bash
sudo pacman -S go base-devel
```

**macOS:**
```bash
brew install go
```

**Windows:**
Download from https://go.dev/dl/ and run the installer.

**Verify:**
```bash
go version    # should show 1.22+
```

---

## Clone and Build

```bash
# Clone
git clone https://github.com/viviotv/vivio-cli.git
cd vivio

# Build
make build

# Binary at: bin/vivio
```

**Without Make:**
```bash
go build -o bin/vivio ./cmd
```

**Windows:**
```powershell
go build -o bin\vivio.exe .\cmd
```

---

## Development Workflow

### Run without installing

```bash
# Using make (recommended)
make run ARGS='list channels --country=IN'
make run ARGS='play "BBC News"'

# Or directly with go run
go run ./cmd list channels --country=IN
go run ./cmd play "BBC News"
```

---

### Build for multiple platforms

**Easy way (recommended):**
```bash
make release
# Builds for all platforms and creates dist/ with compressed binaries
```

**Manual builds:**
```bash
# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/vivio-linux-amd64 ./cmd

# macOS Intel
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/vivio-darwin-amd64 ./cmd

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/vivio-darwin-arm64 ./cmd

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/vivio-windows-amd64.exe ./cmd
```

---

### Test

```bash
# Run all tests with coverage
make test

# Or manually
go test -v ./...
```

---

### Format code

```bash
# Format all Go code
make fmt

# Or manually
go fmt ./...
goimports -w ./cmd ./pkg
```

---

## Project Structure

```
vivio/
├── cmd/                # CLI application
│   ├── main.go         # Entry point
│   └── commands/       # Cobra commands (package commands)
│       ├── root.go     # Root command & client init
│       ├── play.go     # Play command
│       ├── list.go     # List command
│       ├── search.go   # Search command
│       ├── info.go     # Info command
│       ├── version.go  # Version command
│       └── ...
├── pkg/                # Reusable library (package core)
│   ├── channel.go      # Data models
│   ├── client.go       # CoreClient (main API)
│   ├── fetcher.go      # API fetching from iptv-org
│   ├── cache.go        # Disk cache with TTL
│   ├── search.go       # Filtering logic
│   ├── advanced_search.go  # Advanced query parser
│   └── countries.go    # Country code → name map
├── examples/           # Usage example scripts
├── docs/               # Documentation
├── bin/                # Build output (local dev)
├── dist/               # Release output (all platforms)
├── go.mod              # Single Go module
└── Makefile            # Build automation
```

---

## Creating Release Binaries

```bash
# Build for all platforms (Linux, macOS, Windows - amd64 + arm64)
make release

# Output in dist/:
# vivio-VERSION-linux-amd64.tar.gz
# vivio-VERSION-linux-arm64.tar.gz
# vivio-VERSION-darwin-amd64.tar.gz      (Intel Mac)
# vivio-VERSION-darwin-arm64.tar.gz      (M chip Mac)
# vivio-VERSION-windows-amd64.exe.tar.gz
# vivio-VERSION-windows-arm64.exe.tar.gz
# SHA256SUMS                              (checksums)
```

Features:
- Static binaries (CGO_ENABLED=0)
- Version from git tags (e.g., v1.0.0)
- Compressed archives
- SHA256 checksums for verification

---

## Creating Distribution Packages

### RPM (Fedora/RHEL)

1. Install tools:
```bash
sudo dnf install rpm-build rpmdevtools
```

2. Create spec file: `vivio.spec`

3. Build:
```bash
rpmbuild -ba vivio.spec
```

See [Fedora Packaging Guidelines](https://docs.fedoraproject.org/en-US/packaging-guidelines/)

---

### DEB (Ubuntu/Debian)

1. Install tools:
```bash
sudo apt install build-essential debhelper
```

2. Create `debian/` directory structure

3. Build:
```bash
dpkg-buildpackage -us -uc
```

See [Debian Packaging](https://www.debian.org/doc/manuals/maint-guide/)

---

### Homebrew (macOS)

Create a tap repository and formula:

```ruby
class Vivio < Formula
  desc "Free live TV streaming CLI"
  homepage "https://github.com/viviotv/vivio-cli"
  url "https://github.com/viviotv/vivio-cli/archive/v1.0.0.tar.gz"
  sha256 "..."

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"vivio", "./cmd"
  end
end
```

See [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)

---

### Scoop (Windows)

Create a manifest in a bucket:

```json
{
    "version": "1.0.0",
    "description": "Free live TV streaming",
    "homepage": "https://github.com/viviotv/vivio-cli",
    "license": "Apache-2.0",
    "architecture": {
        "64bit": {
            "url": "https://github.com/viviotv/vivio-cli/releases/download/v1.0.0/vivio-windows-amd64.exe",
            "bin": [["vivio-windows-amd64.exe", "vivio"]]
        }
    }
}
```

See [Scoop Documentation](https://github.com/ScoopInstaller/Scoop/wiki)

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make changes and test
4. Commit: `git commit -m "Add feature"`
5. Push: `git push origin feature-name`
6. Open a Pull Request

---

## Code Style

- Follow Go conventions: `go fmt`
- Keep functions short and focused
- Document exported functions
- Use meaningful variable names
- Write tests for new features

---

## Adding New Commands

1. Create new file in `cmd/commands/`
2. Define cobra command
3. Register in `init()` function

Example:
```go
// cmd/commands/mycommand.go
package commands

import (
    "github.com/spf13/cobra"
    core "github.com/viviotv/vivio/pkg"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    Long:  `Detailed description with examples.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Use the global client
        if err := initClient(); err != nil {
            return err
        }
        if err := client.Load(); err != nil {
            return err
        }
        
        // Your implementation here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

---

## Debugging

```bash
# Enable verbose logging
vivio --verbose list channels

# Print cache location
echo ~/.cache/vivio/

# Check cache contents
cat ~/.cache/vivio/channels.json | jq
```

---

## Next Steps

- Read [CLI.md](CLI.md) for command reference
- Check existing issues on GitHub
- Join discussions for feature ideas
