package hook

import (
	"strings"

	model "gr-tr/src/models"
)

type Filters struct {
	Search       string
	FirstAlbum   string
	CreationYear string
	Members      string
}

func SearchArtists(artists []model.Artist, search string) []model.Artist {
	if search == ""{
		return artists
	}
	filtered := make([]model.Artist, 0, len(artists))

	for _, artist := range artists {
		if !matchSearch(artist, search) {
			continue
		}

		filtered = append(filtered, artist)
	}

	return filtered
}

func matchSearch(artist model.Artist, search string) bool {
	search = strings.TrimSpace(strings.ToLower(search))
	if search == "" {
		return true
	}

	if strings.Contains(strings.ToLower(artist.Name), search) {
		return true
	}
	if strings.Contains(artist.FirstAlbum, search) {
		return true
	}
	if strings.Contains(artist.ConcertDates, search) {
		return true
	}
	for _, v := range artist.Members {
		if strings.Contains(strings.ToLower(v), search) {
			return true
		}
	}
	var locations model.Locations
	if err := ApiCall(artist.Locations, &locations); err != nil {
		return false
	}

	for _, location := range locations.Locations {
		normalized := strings.ToLower(strings.ReplaceAll(location, "_", " "))
		if strings.Contains(normalized, search) {
			return true
		}
	}

	return false
}
