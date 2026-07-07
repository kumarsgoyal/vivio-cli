package core

import (
	regexp "regexp"
	strings "strings"
)

// SearchQuery represents parsed advanced search criteria
type SearchQuery struct {
	Countries  []string // e.g. ["IN", "US", "UK"]
	Categories []string // e.g. sports, news
	Languages  []string
	Qualities  []string
	NameQuery  string // free-text search for channel name
}

// ParseSearchQuery parses advanced search syntax like:
// "country:IN category:sports quality:1080p BBC"
// "country:US|UK language:en news"
func ParseSearchQuery(query string) *SearchQuery {
	sq := &SearchQuery{}

	// Regex patterns for filters
	countryRe := regexp.MustCompile(`country:([a-zA-Z|]+)`)
	categoryRe := regexp.MustCompile(`category:([a-zA-Z|]+)`)
	languageRe := regexp.MustCompile(`language:([a-zA-Z|]+)`)
	qualityRe := regexp.MustCompile(`quality:([0-9a-zA-Z|]+)`)

	// Extract countries
	if match := countryRe.FindStringSubmatch(query); match != nil {
		sq.Countries = strings.Split(match[1], "|")
		query = countryRe.ReplaceAllString(query, "")
	}

	// Extract categories
	if match := categoryRe.FindStringSubmatch(query); match != nil {
		sq.Categories = strings.Split(match[1], "|")
		query = categoryRe.ReplaceAllString(query, "")
	}

	// Extract languages
	if match := languageRe.FindStringSubmatch(query); match != nil {
		sq.Languages = strings.Split(match[1], "|")
		query = languageRe.ReplaceAllString(query, "")
	}

	// Extract qualities
	if match := qualityRe.FindStringSubmatch(query); match != nil {
		sq.Qualities = strings.Split(match[1], "|")
		query = qualityRe.ReplaceAllString(query, "")
	}

	// Remaining text is name query
	sq.NameQuery = strings.TrimSpace(query)

	return sq
}

// AdvancedFilter filters channels based on advanced search query
func AdvancedFilter(channels []*Channel, query string) []*Channel {
	sq := ParseSearchQuery(query)
	result := make([]*Channel, 0)

	for _, ch := range channels {
		// Check country filter
		if len(sq.Countries) > 0 && !matchesAny(ch.Country, sq.Countries) {
			continue
		}

		// Check category filter
		if len(sq.Categories) > 0 && !matchesAny(ch.Category, sq.Categories) {
			continue
		}

		// Check language filter (not yet available in current data)
		// Future: when language field is populated
		if len(sq.Languages) > 0 && ch.Language != "" && !matchesAny(ch.Language, sq.Languages) {
			continue
		}

		// Check quality filter
		if len(sq.Qualities) > 0 && !matchesAny(ch.Quality, sq.Qualities) {
			continue
		}

		// Check name query
		if sq.NameQuery != "" && !strings.Contains(strings.ToLower(ch.Name), strings.ToLower(sq.NameQuery)) {
			continue
		}

		result = append(result, ch)
	}

	return result
}

// matchesAny checks if value matches any of the options (case-insensitive)
func matchesAny(value string, options []string) bool {
	for _, opt := range options {
		if strings.EqualFold(opt, value) {
			return true
		}
	}
	return false
}
