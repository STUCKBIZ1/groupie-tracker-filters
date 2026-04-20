package handler

import (
	"html/template"
	"net/http"
	"strconv"

	// "strings"

	hook "gr-tr/src/hook"
	model "gr-tr/src/models"
)

const ARTIST_URL = "https://groupietrackers.herokuapp.com/api/artists"

var suggestions []string

func Suggestion(artists []model.Artist) []string {
	var d_lacations model.Locations
	var suggestions []string

	for _, artist := range artists {
		suggestions = append(suggestions, artist.Name)
		for _, a := range artist.Members {
			suggestions = append(suggestions, a)
		}
		suggestions = append(suggestions, artist.FirstAlbum)
		suggestions = append(suggestions, strconv.Itoa(artist.CreationDate))
		hook.ApiCall(artist.Locations, &d_lacations)
		for _, v := range d_lacations.Locations {
			suggestions = append(suggestions, v)
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
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "500 - Template Error", "Unable to load page template.")
		return
	}
	if len(suggestions) == 0 {
		suggestions = Suggestion(artists)
	}
	search := r.URL.Query().Get("search")
	ArSu := model.AS{
		Artists:    hook.SearchArtists(artists, search),
		Suggestion: suggestions,
	}

	if err = tmpl.Execute(w, ArSu); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "500 - Render Error", "Failed to render page.")
		return
	}
}
