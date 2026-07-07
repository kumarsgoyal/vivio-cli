package core

import (
	json "encoding/json"
	os "os"
	filepath "path/filepath"
	time "time"
)

type cacheEntry struct {
	FetchedAt int64      `json:"fetched_at"`
	TTL       int64      `json:"ttl_seconds"`
	Channels  []*Channel `json:"channels"`
}

type diskCache struct {
	dir string
}

func newDiskCache(dir string) *diskCache {
	_ = os.MkdirAll(dir, 0o755)
	return &diskCache{dir: dir}
}

func (c *diskCache) filePath() string {
	return filepath.Join(c.dir, "channels.json")
}

// load reads cached channels from disk.
// Returns (channels, true) if cache exists and is not expired.
// Returns (nil, false) if missing, corrupt, or expired.
func (c *diskCache) load() ([]*Channel, bool) {
	data, err := os.ReadFile(c.filePath())
	if err != nil {
		return nil, false
	}
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	if time.Now().Unix()-entry.FetchedAt > entry.TTL {
		return nil, false
	}
	return entry.Channels, true
}

func (c *diskCache) save(channels []*Channel, ttl time.Duration) error {
	entry := cacheEntry{
		FetchedAt: time.Now().Unix(),
		TTL:       int64(ttl.Seconds()),
		Channels:  channels,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath(), data, 0o644)
}

func (c *diskCache) clear() error {
	err := os.Remove(c.filePath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
