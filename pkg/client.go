package core

import (
	sync "sync"
	time "time"
)

const channelTTL = 6 * time.Hour

// CoreClient is the main entry point for all platforms.
// Exported for gomobile — all method signatures use gomobile-compatible types.
type CoreClient struct {
	cache    *diskCache
	mu       sync.RWMutex
	channels []*Channel
}

// NewCoreClient creates a CoreClient with a disk cache at cacheDir.
// cacheDir is created automatically if it does not exist.
func NewCoreClient(cacheDir string) *CoreClient {
	return &CoreClient{
		cache: newDiskCache(cacheDir),
	}
}

// Load populates the channel list from disk cache or network.
// Safe to call on app startup — returns immediately if cache is warm.
func (c *CoreClient) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.cache.load(); ok {
		c.channels = cached
		return nil
	}
	return c.refresh()
}

// Refresh forces a fresh fetch from iptv-org regardless of TTL.
func (c *CoreClient) Refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.refresh()
}

func (c *CoreClient) refresh() error {
	type chanResult struct {
		channels []apiChannel
		err      error
	}
	type streamResult struct {
		streams []apiStream
		err     error
	}

	chCh := make(chan chanResult, 1)
	stCh := make(chan streamResult, 1)

	go func() {
		v, err := fetchChannels()
		chCh <- chanResult{v, err}
	}()
	go func() {
		v, err := fetchStreams()
		stCh <- streamResult{v, err}
	}()

	chRes := <-chCh
	stRes := <-stCh

	if chRes.err != nil {
		return chRes.err
	}
	if stRes.err != nil {
		return stRes.err
	}

	merged := merge(chRes.channels, stRes.streams)
	c.channels = merged
	return c.cache.save(merged, channelTTL)
}

// GetAll returns all available channels.
func (c *CoreClient) GetAll() *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(c.channels)
}

// GetByCountry returns channels for the given country code (e.g. "IN", "US", "GB").
func (c *CoreClient) GetByCountry(country string) *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(filterChannels(c.channels, country, "", ""))
}

// GetByCategory returns channels for the given category (e.g. "sports", "news", "music").
func (c *CoreClient) GetByCategory(category string) *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(filterChannels(c.channels, "", category, ""))
}

// Search returns channels whose name contains the query string (case-insensitive).
func (c *CoreClient) Search(query string) *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(filterChannels(c.channels, "", "", query))
}

// Filter returns channels matching all non-empty parameters.
// Pass an empty string to skip any filter.
func (c *CoreClient) Filter(country, category, query string) *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(filterChannels(c.channels, country, category, query))
}

// Countries returns a deduplicated list of all country codes present in the data.
func (c *CoreClient) Countries() *StringList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return uniqueStrings(c.channels, func(ch *Channel) string { return ch.Country })
}

// Categories returns a deduplicated list of all categories present in the data.
func (c *CoreClient) Categories() *StringList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return uniqueStrings(c.channels, func(ch *Channel) string { return ch.Category })
}

// AdvancedSearch performs search with advanced syntax like:
// "country:IN category:sports quality:1080p"
// "country:US|UK language:en news"
func (c *CoreClient) AdvancedSearch(query string) *ChannelList {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return newChannelList(AdvancedFilter(c.channels, query))
}

// GetChannelByID returns a single channel by its ID, or nil if not found.
func (c *CoreClient) GetChannelByID(id string) *Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, ch := range c.channels {
		if ch.ID == id {
			return ch
		}
	}
	return nil
}

// TotalChannels returns the count of loaded channels.
func (c *CoreClient) TotalChannels() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.channels)
}

// ClearCache deletes the on-disk channel cache.
func (c *CoreClient) ClearCache() error {
	return c.cache.clear()
}
