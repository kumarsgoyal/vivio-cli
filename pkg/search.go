package core

import "strings"

func filterChannels(channels []*Channel, country, category, query string) []*Channel {
	country = strings.ToLower(strings.TrimSpace(country))
	category = strings.ToLower(strings.TrimSpace(category))
	query = strings.ToLower(strings.TrimSpace(query))

	result := make([]*Channel, 0)
	for _, ch := range channels {
		if country != "" && !strings.EqualFold(ch.Country, country) {
			continue
		}
		if category != "" && !strings.EqualFold(ch.Category, category) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(ch.Name), query) {
			continue
		}
		result = append(result, ch)
	}
	return result
}

func uniqueStrings(channels []*Channel, field func(*Channel) string) *StringList {
	seen := make(map[string]struct{})
	items := make([]string, 0)
	for _, ch := range channels {
		v := field(ch)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			items = append(items, v)
		}
	}
	return &StringList{items: items}
}
