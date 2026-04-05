package hook

import (
	"strconv"
	"strings"
	"time"

	model "gr-tr/src/models"
)

type Filters struct {
	Search       string
	FirstAlbum   string
	CreationYear string
	Members      string
}

func FilterArtists(artists []model.Artist, filters Filters) []model.Artist {
	filtered := make([]model.Artist, 0, len(artists))

	for _, artist := range artists {
		if !matchSearch(artist, filters.Search) {
			continue
		}
		if !matchFirstAlbum(artist.FirstAlbum, filters.FirstAlbum) {
			continue
		}
		if !matchCreationYear(artist.CreationDate, filters.CreationYear) {
			continue
		}
		if !matchMembers(len(artist.Members), filters.Members) {
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

func matchFirstAlbum(firstAlbum string, selected string) bool {
	selected = strings.TrimSpace(selected)
	if selected == "" {
		return true
	}

	parsedDate, err := time.Parse("2006-01-02", selected)
	if err != nil {
		return true
	}

	return firstAlbum == parsedDate.Format("02-01-2006")
}

func matchCreationYear(creationDate int, selected string) bool {
	selected = strings.TrimSpace(selected)
	if selected == "" {
		return true
	}

	selectedYear, err := strconv.Atoi(selected)
	if err != nil {
		return true
	}

	return creationDate >= selectedYear
}

func matchMembers(memberCount int, selected string) bool {
	selected = strings.TrimSpace(selected)
	if selected == "" {
		return true
	}

	if selected == "8" {
		return memberCount >= 8
	}

	return selected == intToString(memberCount)
}

func intToString(value int) string {
	return strconv.Itoa(value)
}
