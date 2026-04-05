package handler

import (
	"html/template"
	"net/http"

	hook "gr-tr/src/hook"
	model "gr-tr/src/models"
)

const ARTIST_URL = "https://groupietrackers.herokuapp.com/api/artists"

type HomePageData struct {
	Artists []model.Artist
	Filters hook.Filters
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(
			w,
			r,
			http.StatusMethodNotAllowed,
			"405 - Method Not Allowed",
			"Only GET method is allowed for this endpoint.",
		)
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(
			w,
			r,
			http.StatusNotFound,
			"404 - Page Not Found",
			"The page you are looking for does not exist.",
		)
		return
	}

	var artists []model.Artist

	err := hook.ApiCall(ARTIST_URL, &artists)
	if err != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"500 - Server Error",
			"Failed to load artists data.",
		)
		return
	}

	filters := hook.Filters{
		Search:       r.URL.Query().Get("search"),
		FirstAlbum:   r.URL.Query().Get("first_album"),
		CreationYear: r.URL.Query().Get("creation_year"),
		Members:      r.URL.Query().Get("members"),
	}

	pageData := HomePageData{
		Artists: hook.FilterArtists(artists, filters),
		Filters: filters,
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"500 - Template Error",
			"Unable to load page template.",
		)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"500 - Render Error",
			"Failed to render page.",
		)
		return
	}
}
