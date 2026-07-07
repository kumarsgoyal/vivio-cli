package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	channelsURL = "https://iptv-org.github.io/api/channels.json"
	streamsURL  = "https://iptv-org.github.io/api/streams.json"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func fetchChannels() ([]apiChannel, error) {
	resp, err := httpClient.Get(channelsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch channels: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch channels: status %d", resp.StatusCode)
	}
	var result []apiChannel
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode channels: %w", err)
	}
	return result, nil
}

func fetchStreams() ([]apiStream, error) {
	resp, err := httpClient.Get(streamsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch streams: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch streams: status %d", resp.StatusCode)
	}
	var result []apiStream
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode streams: %w", err)
	}
	return result, nil
}

// logoURL constructs the channel logo URL from the iptv-org database CDN.
func logoURL(channelID string) string {
	return "https://cdn.jsdelivr.net/gh/iptv-org/database@latest/data/logos/" + channelID + ".png"
}

// merge joins channel metadata with stream URLs.
// Stores all streams per channel, with primary being the highest-quality one.
// Skips channels with no matched stream and NSFW channels.
func merge(channels []apiChannel, streams []apiStream) []*Channel {
	// Group all streams by channel ID
	streamsByChannel := make(map[string][]apiStream, len(streams))
	for _, s := range streams {
		if s.Channel == nil || *s.Channel == "" || s.URL == "" {
			continue
		}
		id := *s.Channel
		streamsByChannel[id] = append(streamsByChannel[id], s)
	}

	result := make([]*Channel, 0, len(channels))
	for _, ch := range channels {
		if ch.IsNSFW {
			continue
		}
		channelStreams, ok := streamsByChannel[ch.ID]
		if !ok || len(channelStreams) == 0 {
			continue
		}

		// Sort streams by quality (highest first)
		sortStreamsByQuality(channelStreams)

		// Convert to Stream objects
		streams := make([]Stream, 0, len(channelStreams))
		for _, s := range channelStreams {
			streams = append(streams, Stream{
				URL:       s.URL,
				Quality:   s.Quality,
				Label:     s.Label,
				UserAgent: s.UserAgent,
				Referrer:  s.Referrer,
			})
		}

		// Primary stream is the first (highest quality)
		primary := channelStreams[0]

		result = append(result, &Channel{
			ID:        ch.ID,
			Name:      ch.Name,
			Logo:      logoURL(ch.ID),
			Country:   ch.Country,
			Category:  first(ch.Categories),
			StreamURL: primary.URL,
			Quality:   primary.Quality,
			Streams:   streams,
		})
	}
	return result
}

func sortStreamsByQuality(streams []apiStream) {
	// Simple bubble sort by quality (in-place)
	for i := 0; i < len(streams)-1; i++ {
		for j := i + 1; j < len(streams); j++ {
			if higherQuality(streams[j].Quality, streams[i].Quality) {
				streams[i], streams[j] = streams[j], streams[i]
			}
		}
	}
}

func higherQuality(a, b string) bool {
	rank := map[string]int{
		"2160p": 5, "1080p": 4, "1080i": 3, "720p": 3,
		"576p": 2, "480p": 1, "360p": 0,
	}
	return rank[a] > rank[b]
}

func first(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}
