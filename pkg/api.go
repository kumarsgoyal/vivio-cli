package core

// internal JSON types mirroring the current iptv-org API responses

type apiChannel struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Country    string   `json:"country"`
	Categories []string `json:"categories"`
	IsNSFW     bool     `json:"is_nsfw"`
}

type apiStream struct {
	Channel   *string `json:"channel"` // null for unmatched streams
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	Quality   string  `json:"quality"`
	Label     string  `json:"label"` // e.g. "Geo-blocked", "Not 24/7", ""
	UserAgent string  `json:"user_agent"`
	Referrer  string  `json:"referrer"`
}
