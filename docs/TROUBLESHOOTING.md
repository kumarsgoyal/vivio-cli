# Troubleshooting

Common issues and solutions.

---

## Stream Errors

### 404 Not Found

**Error:**
```
Failed to open https://example.com/stream.m3u8
HTTP Error 404: Not Found
```

**Cause:** External stream server is offline or URL changed.

**Solution:**
```bash
# Get fresh URLs from iptv-org
vivio cache clear
vivio list channels --country=IN

# Try a different channel
vivio search channels "BBC"
vivio play 7943
```

**Reality:** Only 30-50% of iptv-org channels work at any time. Try multiple channels. Large broadcasters (BBC, CNN, Al Jazeera, DW, France 24) are most reliable.

---

### 403 Forbidden

**Error:**
```
HTTP error 403 Forbidden
Server returned 403 Forbidden (access denied)
```

**Causes:**
- Geo-blocking (stream blocked in your region)
- Authentication required (stream needs referrer/token)
- Anti-scraping protection

**Solution:**
Try a different channel. Some streams are region-restricted or require authentication Vivio can't provide.

---

## Player Issues

### No player found

**Error:**
```
No player found. Install mpv: sudo dnf install mpv
```

**Solution:**
```bash
# Install mpv (recommended)
sudo dnf install mpv     # Fedora/RHEL
brew install mpv         # macOS

# Or use ffplay (usually pre-installed)
vivio play 42 --player=ffplay
```

---

### Player window doesn't open

**Check:**
```bash
# Verify DISPLAY is set (Linux)
echo $DISPLAY

# Try forcing a player
vivio play 42 --player=mpv

# Check if player works standalone
mpv https://test.com/stream.m3u8
```

---

## Cache Issues

### Cache seems stale

**Solution:**
```bash
vivio cache clear
vivio list channels
```

---

### Cache location

**Default:** `~/.cache/vivio/channels.json`

**Manual inspection:**
```bash
ls -lh ~/.cache/vivio/
cat ~/.cache/vivio/channels.json | head
```

**Delete manually:**
```bash
rm -rf ~/.cache/vivio/
```

---

## Command Issues

### `vivio: command not found`

**Cause:** Binary not in PATH.

**Solution:**
```bash
# Use full path
/home/user/vivio/bin/vivio list channels

# Or add to PATH
export PATH="$PATH:/home/user/vivio/bin"
echo 'export PATH="$PATH:$HOME/vivio/bin"' >> ~/.bashrc
```

---

### Channel number doesn't match

**Issue:** Filtered list shows different numbers than full list.

**Expected behavior:** Numbers are global and stable.

```bash
vivio list channels                # 10 TV is NO 3736
vivio list channels --country=IN   # 10 TV still shows NO 3736
vivio play 3736                    # plays 10 TV correctly
```

Filters hide channels but preserve numbering so `vivio play <NO>` always works.

---

## Build Issues

### Build fails

```bash
# Check Go version (needs 1.22+)
go version

# Clean build
cd vivio
rm -rf bin/
make build

# Or manually
go build -o bin/vivio ./cli/
```

---

### Missing dependencies

```bash
cd vivio/cli
go mod tidy
go build -o ../bin/vivio .
```

---

## Performance

### Slow channel list loading

**First run after cache clear is slow** (~3-5 seconds) — fetching 5MB from iptv-org.

Subsequent runs are instant (loads from `~/.cache/vivio/`).

**To speed up:**
Cache auto-refreshes in background after 6 hours. You rarely need to clear it manually.

---

## Getting Help

### Check version

```bash
vivio --help
```

### Debug mode

Not implemented yet. For now, check:
```bash
# Cache state
ls -lh ~/.cache/vivio/

# Test player
mpv --version
ffplay -version
```

---

## Still Having Issues?

1. Clear cache: `vivio cache clear`
2. Try a stable channel: `vivio play "BBC News"`
4. File an issue on GitHub with error output
