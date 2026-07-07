# CLI Reference

Complete command reference for Vivio CLI.

---

## Commands

### `vivio info <channel>`

Show detailed information about a channel including all available streams.

**Syntax:**
```bash
vivio info <channel number or name>
```

**Examples:**
```bash
vivio info 3736           # by channel number
vivio info "BBC News"     # by name
```

**Output:**
- Channel ID, Name, Country, Category
- Logo URL  
- All available stream URLs with quality indicators
- Stream labels (Geo-blocked, Not 24/7, etc.)
- Primary stream marked with ★

This is useful to:
- See all alternative streams before playing
- Check stream quality and availability
- Identify geo-blocked streams

---

### `vivio list channels`

List all channels or filter by country/category.

**Syntax:**
```bash
vivio list channels [--country=CODE] [--category=NAME]
```

**Flags:**
- `--country=CODE` — 2-letter ISO code (IN, US, GB, AU, etc.)
- `--category=NAME` — category (sports, news, music, movies, etc.)

**Examples:**
```bash
vivio list channels                                # all channels
vivio list channels --country=IN                   # India only
vivio list channels --category=sports              # sports worldwide
vivio list channels --country=IN --category=news   # India news
```

**Output:**
```
NO    COUNTRY  NAME                    CATEGORY  QUALITY
3736  IN       10 TV                   news      720p
3737  IN       22Scope News            news      1080p
```

---

### `vivio search channels <query>`

Search channels by name or use advanced filters.

**Syntax:**
```bash
vivio search channels "query"
vivio search channels "filter:value filter:value query"
```

**Filters:**
- `country:CODE` — Filter by country (e.g., `country:IN`)
- `category:NAME` — Filter by category (e.g., `category:sports`)
- `quality:RES` — Filter by quality (e.g., `quality:1080p`)
- `language:CODE` — Filter by language (future)

Use `|` for OR logic: `country:US|UK|IN`

**Examples:**
```bash
# Simple name search
vivio search channels "BBC"
vivio search channels "star sports"

# Advanced filters
vivio search channels "country:IN category:sports"
vivio search channels "country:US|UK quality:1080p news"
vivio search channels "category:sports quality:1080p"
vivio search channels "country:IN 1080p"
```

**Output:**
```
NO    COUNTRY  NAME               CATEGORY       QUALITY
7936  UK       BBC Alba           general        
7937  UK       BBC Arabic         news           720p
```

---

### `vivio play <NO or "name">`

Play a channel by number or name with automatic fallback to alternative streams.

**Syntax:**
```bash
vivio play <number>
vivio play "channel name"
vivio play <NO or "name"> --player=PLAYER
```

**Flags:**
- `--player=mpv` — force mpv player
- `--player=ffplay` — force ffplay player

**Examples:**
```bash
vivio play 3736                    # play channel #3736
vivio play "BBC News"              # play by name
vivio play 42 --player=ffplay      # force ffplay
```

**Auto Fallback:**
If a channel has multiple streams, Vivio automatically tries alternatives if the primary fails:
1. Tries primary stream (highest quality)
2. On error (403/404), automatically tries next stream
3. Continues until one works or all fail
4. Shows which alternative is being tried

**Player controls:**
| Key | Action |
|-----|--------|
| `space` | pause/play |
| `f` | fullscreen |
| `q` | quit |
| `← →` | seek (ffplay) |
| `9 / 0` | volume (mpv) |

---

### `vivio countries list`

List all countries with active channels (sorted A→Z).

**Syntax:**
```bash
vivio countries list
```

**Output:**
```
CODE   COUNTRY
AF     Afghanistan
IN     India
US     United States
```

Use CODE with `--country=` flag.

---

### `vivio categories list`

List all available categories (sorted A→Z).

**Syntax:**
```bash
vivio categories list
```

**Output:**
```
animation
business
comedy
news
sports
```

Use with `--category=` flag.

---

### `vivio cache clear`

Delete local cache. Next command fetches fresh data.

**Syntax:**
```bash
vivio cache clear
```

**When to use:**
- Stream URL is dead (404/403)
- Want latest channels
- Cache corrupted

**Cache details:**
- Location: `~/.cache/vivio/channels.json`
- TTL: 6 hours
- Auto-refreshes in background

---

## Complete Examples

### Find Indian sports channels

```bash
# See all countries
vivio countries list

# List Indian sports
vivio list channels --country=IN --category=sports

# Output shows:
# NO    COUNTRY  NAME              CATEGORY  QUALITY
# 3845  IN       DD Sports         sports    1080p
# 3902  IN       Star Sports 1     sports    1080p

# Play one
vivio play 3902
```

---

### Search and play BBC News

```bash
# Search
vivio search channels "BBC News"

# Output:
# NO    COUNTRY  NAME         CATEGORY  QUALITY
# 7943  UK       BBC News     news      1080p

# Play by number
vivio play 7943

# Or by name
vivio play "BBC News"
```

---

### Browse all news worldwide

```bash
vivio categories list           # see all categories
vivio list channels --category=news
```

---

### Switch players

```bash
# Auto-detect (mpv preferred)
vivio play 42

# Force mpv
vivio play 42 --player=mpv

# Force ffplay
vivio play 42 --player=ffplay
```

---

## How Numbering Works

**Channel numbers (NO) are globally stable:**
```bash
vivio list channels               # 10 TV is NO 3736
vivio list channels --country=IN  # 10 TV is still NO 3736
vivio play 3736                   # always plays 10 TV
```

Numbers come from the sorted full list (by Country → Name). Filters hide channels but preserve numbering, so the number you see is always the number you play.

---

## See Also

- [Installation Guide](INSTALL.md)
- [Troubleshooting](TROUBLESHOOTING.md)
