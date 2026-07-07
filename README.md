# Vivio

**Free live TV streaming — 10,000+ channels worldwide, zero cost, no account.**

Stream public IPTV channels from No ads, no signups, no video hosting.

See [Installation Guide](docs/INSTALL.md) for setup instructions.

---

## Features

- 🌍 **10,000+ channels** from 176 countries
- 📺 **28 categories** — news, sports, movies, music, kids, and more
- 🔍 **Advanced search** — filter by country, category, quality with OR logic
- 🔢 **Numbered channels** — `vivio play 42` plays channel #42
- 📋 **Channel info** — see all available streams with quality indicators
- 🔄 **Auto fallback** — retry alternative streams if primary fails
- 🎬 **Multiple players** — mpv, ffplay, or any IPTV player
- 💾 **Offline cache** — instant browsing with 6-hour cache
- 🆓 **Zero cost** — no servers, streams go directly from CDN to your device

---

## Quick Start

```bash
# List channels
vivio list channels --country=IN --category=sports

# Search with advanced filters
vivio search channels "country:US|UK quality:1080p news"

# Check channel info and available qualities
vivio info "BBC News"

# Play (auto-retries if stream fails)
vivio play "BBC News" --quality=720p
```


---

## Usage

```bash
# Browse channels
vivio list channels
vivio list channels --country=IN --category=sports

# Advanced search with filters
vivio search channels "country:IN category:sports quality:1080p"
vivio search channels "country:US|UK news"

# Channel info (shows all streams)
vivio info "BBC News"
vivio info 3736

# Play (auto-tries fallback streams if primary fails)
vivio play 42              # by number
vivio play "BBC News"      # by name

# Discovery
vivio countries list       # see all countries
vivio categories list      # see all categories
```

See [Usage Guide](docs/QUICK_START.md) for usage instructions.

See [CLI Documentation](docs/CLI.md) for full command reference.

---

## Platforms

| Platform | Status | Docs |
|----------|--------|------|
| CLI (Linux/Mac/Windows) | ✅ Ready | [CLI.md](docs/CLI.md) |
| Android / Android TV | 🚧 Planned | — |
| iOS / Apple TV | 🚧 Planned | — |
| Web UI | 🚧 Planned | — |
---

## How It Works

Vivio fetches channel metadata from iptv-org's free API, caches it locally, and passes stream URLs to your player. No video passes through our code — streams go directly from external CDNs to your device.

```
iptv-org (free API) → Vivio (cache + directory) → mpv/ffplay → External CDN
```


---

## Documentation

**Getting Started:**
- [Quick Start Guide](docs/QUICK_START.md) — learn all commands in 5 minutes

**User Guides:**
- [Installation Guide](docs/INSTALL.md) — detailed setup for all platforms
- [CLI Reference](docs/CLI.md) — all commands, flags, examples
- [Troubleshooting](docs/TROUBLESHOOTING.md) — common issues and fixes

**Developer Guides:**
- [Development Guide](docs/DEVELOPMENT.md) — build from source, create packages

---

## Project Structure

```
vivio/
├── cmd/          # CLI binary (main package)
├── pkg/          # Reusable library code
├── examples/     # Usage examples
├── docs/         # Documentation
└── bin/          # Build output directory
```

---

## License

Apache 2.0 — see [LICENSE](LICENSE)

---

## Credits

- Channel data: [iptv-org/iptv](https://github.com/iptv-org/iptv)
- Built with: [cobra](https://github.com/spf13/cobra)
