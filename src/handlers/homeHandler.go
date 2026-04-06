package handler

import (
	"html/template"
	"net/http"
	"strings"

	hook "gr-tr/src/hook"
	model "gr-tr/src/models"
)

const ARTIST_URL = "https://groupietrackers.herokuapp.com/api/artists"

// Suggestion holds a search suggestion value and its type label.
type Suggestion struct {
	Value string // the text that fills the input
	Label string // the type shown beside it (e.g. "artist", "member", "location")
}

type HomePageData struct {
	Artists     []model.Artist
	Filters     hook.Filters
	Suggestions []Suggestion
}

// buildSuggestions creates datalist entries from all artists.
// It deduplicates by (value, label) pairs.
func buildSuggestions(artists []model.Artist) []Suggestion {
	seen := make(map[string]bool)
	var suggestions []Suggestion

	add := func(value, label string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		key := strings.ToLower(value) + "|" + label
		if !seen[key] {
			seen[key] = true
			suggestions = append(suggestions, Suggestion{Value: value, Label: label})
		}
	}

	for _, a := range artists {
		// Artist / band name
		add(a.Name, "artist/band")

		// Members
		for _, m := range a.Members {
			add(m, "member")
		}

		// First album date
		add(a.FirstAlbum, "first album")

		// Locations (fetched individually per artist)
		var locs model.Locations
		if err := hook.ApiCall(a.Locations, &locs); err == nil {
			for _, loc := range locs.Locations {
				// normalise underscores to spaces for readability
				display := strings.ReplaceAll(loc, "_", " ")
				add(display, "location")
			}
		}
	}

	return suggestions
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "405 - Method Not Allowed", "Only GET method is allowed for this endpoint.")
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound, "404 - Page Not Found", "The page you are looking for does not exist.")
		return
	}

	var artists []model.Artist
	if err := hook.ApiCall(ARTIST_URL, &artists); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "500 - Server Error", "Failed to load artists data.")
		return
	}

	filters := hook.Filters{
		Search:       r.URL.Query().Get("search"),
		FirstAlbum:   r.URL.Query().Get("first_album"),
		CreationYear: r.URL.Query().Get("creation_year"),
		Members:      r.URL.Query().Get("members"),
	}

	pageData := HomePageData{
		Artists:     hook.FilterArtists(artists, filters),
		Filters:     filters,
		Suggestions: buildSuggestions(artists), // always build from full list
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "500 - Template Error", "Unable to load page template.")
		return
	}

	if err = tmpl.Execute(w, pageData); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "500 - Render Error", "Failed to render page.")
		return
	}
}
