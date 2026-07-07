# Vivio - Quick Start Guide

## Basic Commands

### 1. Browse All Channels

```bash
vivio list channels
```

Shows all 10,000+ channels with numbers, countries, names, categories, quality.

---

### 2. Filter by Country

```bash
vivio list channels --country=IN    # India
vivio list channels --country=US    # USA
vivio list channels --country=UK    # United Kingdom
```

---

### 3. Filter by Category

```bash
vivio list channels --category=sports
vivio list channels --category=news
vivio list channels --category=movies
```

---

### 4. Combine Filters

```bash
vivio list channels --country=IN --category=sports
```

---

### 5. Search by Name

```bash
vivio search channels "BBC"
vivio search channels "CNN"
vivio search channels "star sports"
```

---

### 6. Advanced Search with Filters

```bash
# Sports channels from India in 1080p
vivio search channels "country:IN category:sports quality:1080p"

# News from US or UK in high quality
vivio search channels "country:US|UK quality:1080p news"

# All 1080p sports channels
vivio search channels "category:sports quality:1080p"
```

**Filter syntax:**
- `country:CODE` - Filter by country (e.g., `country:IN`)
- `category:NAME` - Filter by category (e.g., `category:sports`)
- `quality:RES` - Filter by quality (e.g., `quality:1080p`)
- Use `|` for OR - `country:US|UK|IN`

---

### 7. Get Channel Info

```bash
# By number
vivio info 3736

# By name
vivio info "BBC News"
```

**Shows:**
- Available qualities (1080p, 720p, etc.)
- All stream URLs
- Geo-blocking info
- Logo, category, etc.

---

### 8. Play a Channel

```bash
# By number
vivio play 3736

# By name
vivio play "BBC News"

# Choose specific quality
vivio play "BBC News" --quality=720p

# Force specific player
vivio play "BBC News" --player=mpv
vivio play "BBC News" --player=ffplay
```

**Auto features:**
- Tries best quality by default
- Auto-retries alternative streams if primary fails
- Shows which stream is being tried

---

### 9. Discover What's Available

```bash
# See all countries
vivio countries list

# See all categories
vivio categories list
```

---

### 10. Cache Management

```bash
# Clear cache (fetch fresh data)
vivio cache clear

# Cache auto-refreshes every 6 hours
```

---

## Common Workflows

### Find and play Indian sports channel

```bash
# 1. List Indian sports channels
vivio list channels --country=IN --category=sports

# Output:
# NO    COUNTRY  NAME              CATEGORY  QUALITY
# 3882  IN       DD Sports         sports    1080p
# 3908  IN       Star Sports 1     sports    1080p

# 2. Play by number
vivio play 3882
```

---

### Search and play BBC News in 720p

```bash
# 1. Search for BBC
vivio search channels "BBC News"

# Output:
# NO    COUNTRY  NAME       CATEGORY  QUALITY
# 7943  UK       BBC News   news      1080p

# 2. Check available qualities
vivio info 7943

# Output:
# Available Qualities: 1080p, 720p, 576p, 540p

# 3. Play in 720p
vivio play 7943 --quality=720p
```

---

### Find high-quality news channels

```bash
# Search for 1080p news from US or UK
vivio search channels "country:US|UK quality:1080p category:news"

# Play first result
vivio play 8150
```


## Player Controls

### mpv
- `Space` - Pause/Play
- `f` - Fullscreen
- `q` - Quit
- `9/0` - Volume down/up

### ffplay
- `Space` - Pause/Play
- `f` - Fullscreen
- `q` - Quit
- `←/→` - Seek backward/forward

---

## Getting Help

```bash
vivio --help                    # Main help
vivio list --help               # List command help
vivio play --help               # Play command help
vivio search channels --help    # Search help
```

---

## Troubleshooting

**Stream won't play (403/404)?**
- Channel may be geo-blocked or offline
- Try a different channel (only 30-50% work at any time)
- Large channels (BBC, CNN, Al Jazeera) are most reliable

**No player found?**
```bash
sudo dnf install mpv    # Fedora/RHEL
sudo apt install mpv    # Ubuntu/Debian
brew install mpv        # macOS
```

**Cache seems old?**
```bash
vivio cache clear
```

---

## Next Steps

- Full documentation: [CLI.md](docs/CLI.md)
- Troubleshooting: [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)
