package core

// Stream represents a single stream URL for a channel with its metadata.
type Stream struct {
	URL       string
	Quality   string
	Label     string // e.g. "Geo-blocked", "Not 24/7", ""
	UserAgent string
	Referrer  string
}

// Channel is the unified model combining iptv-org channel metadata and its stream URLs.
type Channel struct {
	ID        string
	Name      string
	Logo      string
	Country   string
	Language  string
	Category  string
	StreamURL string   // Primary stream (highest quality)
	Quality   string   // Quality of primary stream
	Streams   []Stream // All available streams for fallback
	Website   string   // Channel website
	Network   string   // Broadcasting network
}

// ChannelList wraps a slice of Channel for gomobile compatibility.
// gomobile cannot export slice types directly — use Get/Size accessors.
type ChannelList struct {
	items []*Channel
}

func newChannelList(channels []*Channel) *ChannelList {
	return &ChannelList{items: channels}
}

func (l *ChannelList) Size() int {
	return len(l.items)
}

func (l *ChannelList) Get(index int) *Channel {
	if index < 0 || index >= len(l.items) {
		return nil
	}
	return l.items[index]
}

// StringList wraps a string slice for gomobile compatibility.
type StringList struct {
	items []string
}

func (l *StringList) Size() int {
	return len(l.items)
}

func (l *StringList) Get(index int) string {
	if index < 0 || index >= len(l.items) {
		return ""
	}
	return l.items[index]
}
