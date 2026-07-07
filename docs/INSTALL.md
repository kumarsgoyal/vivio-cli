# Installation Guide

**For end users** — simple one-command install.

---

## Quick Install

> **Note:** Package manager distribution (dnf, apt, brew, scoop) is coming soon! For now, use the manual installation below.

**Coming Soon:**
- `sudo dnf install vivio` (Fedora/RHEL via COPR)
- `sudo apt install vivio` (Ubuntu/Debian via PPA)
- `brew install vivio` (macOS via Homebrew tap)
- `scoop install vivio` (Windows)

For now, see [Manual Installation](#manual-installation) below.

---

## Manual Installation (Current Method)

### 1. Download Binary

**Linux (x64):**
```bash
curl -Lo vivio https://github.com/viviotv/vivio-cli/releases/latest/download/vivio-linux-amd64
chmod +x vivio
sudo mv vivio /usr/local/bin/
```

**macOS (Intel):**
```bash
curl -Lo vivio https://github.com/viviotv/vivio-cli/releases/latest/download/vivio-darwin-amd64
chmod +x vivio
sudo mv vivio /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -Lo vivio https://github.com/viviotv/vivio-cli/releases/latest/download/vivio-darwin-arm64
chmod +x vivio
sudo mv vivio /usr/local/bin/
```

**Windows (x64):**
1. Download: https://github.com/viviotv/vivio-cli/releases/latest/download/vivio-windows-amd64.exe
2. Rename to `vivio.exe`
3. Move to `C:\Windows\System32\` (requires admin)

Or add to your PATH instead.

---

### 2. Install Media Player

You need **at least one player** to watch streams. Install **mpv** (recommended), **ffplay**, or both.

#### Option A: mpv (Recommended — better UI and features)

**Linux:**
```bash
# Fedora/RHEL/CentOS
sudo dnf install mpv

# Ubuntu/Debian
sudo apt install mpv

# Arch
sudo pacman -S mpv
```

**macOS:**
```bash
brew install mpv
```

**Windows:**
```powershell
scoop install mpv
```

**Verify:**
```bash
mpv --version
```

---

#### Option B: ffplay (Lightweight alternative, comes with FFmpeg)

**Linux:**
```bash
# Fedora/RHEL/CentOS
sudo dnf install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Arch
sudo pacman -S ffmpeg
```

**macOS:**
```bash
brew install ffmpeg
```

**Windows:**
```powershell
scoop install ffmpeg
```

**Verify:**
```bash
ffplay -version
```

---

#### Install Both (Recommended)

For best experience, install both players. Vivio will auto-detect mpv first, fall back to ffplay if needed.


---

## Optional: Shell Completion

> **Note:** When installed via package manager (dnf/apt/brew/scoop), completion works automatically. This section is only for manual binary installation.

**Bash:**
```bash
echo 'eval "$(vivio completion bash)"' >> ~/.bashrc
source ~/.bashrc
```

**Zsh:**
```bash
echo 'eval "$(vivio completion zsh)"' >> ~/.zshrc
source ~/.zshrc
```

**Fish:**
```bash
vivio completion fish > ~/.config/fish/completions/vivio.fish
```

**PowerShell (Windows):**
```powershell
vivio completion powershell >> $PROFILE
. $PROFILE
```

---

## Troubleshooting

### `vivio: command not found`

Binary is not in PATH.

**Fix (Linux/macOS):**
```bash
# Use full path
/usr/local/bin/vivio list channels

# Or reinstall to /usr/local/bin/
```

**Fix (Windows):**
Move `vivio.exe` to `C:\Windows\System32\` or add its location to PATH.

---

### `No player found`

Install mpv or ffplay (see step 2 above).

---

### Stream won't play (404/403 error)

External stream is offline. Try different channels:
```bash
vivio search channels "BBC"
vivio search channels "CNN"
```

---

## For Developers

Want to build from source or contribute?  
See [DEVELOPMENT.md](DEVELOPMENT.md)
