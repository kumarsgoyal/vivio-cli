package commands

import (
	sort "sort"
	strings "strings"

	core "github.com/viviotv/vivio/pkg"
)

// getAvailableQualities returns unique qualities from channel streams
func getAvailableQualities(ch *core.Channel) []string {
	seen := make(map[string]bool)
	var qualities []string

	for _, stream := range ch.Streams {
		if stream.Quality != "" && !seen[stream.Quality] {
			seen[stream.Quality] = true
			qualities = append(qualities, stream.Quality)
		}
	}

	return qualities
}

// joinQualities formats qualities list for display
func joinQualities(qualities []string) string {
	if len(qualities) == 0 {
		return "N/A"
	}

	result := ""
	for i, q := range qualities {
		if i > 0 {
			result += ", "
		}
		result += q
	}
	return result
}

// sortedChannels returns channels sorted by Country then Name for consistent numbering.
func sortedChannels(list *core.ChannelList) []*core.Channel {
	result := make([]*core.Channel, list.Size())
	for i := range result {
		result[i] = list.Get(i)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Country != result[j].Country {
			return result[i].Country < result[j].Country
		}
		return result[i].Name < result[j].Name
	})
	return result
}

// isValidCountryCode checks if a country code exists in the data
func isValidCountryCode(code string) bool {
	if client == nil {
		return true // Can't validate without client
	}

	countries := client.Countries()
	for i := 0; i < countries.Size(); i++ {
		if strings.EqualFold(countries.Get(i), code) {
			return true
		}
	}
	return false
}

// isValidCategory checks if a category exists in the data
func isValidCategory(category string) bool {
	if client == nil {
		return true // Can't validate without client
	}

	categories := client.Categories()
	for i := 0; i < categories.Size(); i++ {
		if strings.EqualFold(categories.Get(i), category) {
			return true
		}
	}
	return false
}

// suggestCountries returns a helpful list of common countries
func suggestCountries() string {
	return "Common countries: IN (India), US (USA), UK (United Kingdom), CA (Canada), AU (Australia)"
}

// suggestCategories returns a helpful list of common categories
func suggestCategories() string {
	return "Common categories: news, sports, movies, music, kids, entertainment"
}
