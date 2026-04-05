package handler

import (
	"fmt"
	"html/template"
	"net/http"

	hook "gr-tr/src/hook"
	model "gr-tr/src/models"
)

type T_ArtistData struct {
	Artist    model.Artist
	Relation  model.Relation
	Locations model.Locations
	Dates     model.Dates
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
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
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		ErrorHandler(
			w,
			r,
			http.StatusBadRequest,
			"400 - Bad Request",
			"Artist ID is missing.",
		)
		return
	}

	artistURL := fmt.Sprintf("%s/%s", ARTIST_URL, artistID)
	var d_artist model.Artist
	var d_relation model.Relation
	var d_dates model.Dates
	var d_lacations model.Locations
	err := hook.ApiCall(artistURL, &d_artist)

	err1 := hook.ApiCall(d_artist.Relations, &d_relation)
	err2 := hook.ApiCall(d_artist.ConcertDates, &d_dates)
	err3 := hook.ApiCall(d_artist.Locations, &d_lacations)
	if err != nil || err1 != nil || err2 != nil || err3 != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"Internal Server Error",
			"Failed to process the API request due to invalid or missing data.",
		)
		return
	}
	dArtistData := T_ArtistData{
		Artist:    d_artist,
		Dates:     d_dates,
		Locations: d_lacations,
		Relation:  d_relation,
	}

	tmpl, err := template.ParseFiles("template/artist.html")
	if err != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"500 - Template Error",
			"Unable to load artist page.",
		)
		return
	}

	if err := tmpl.Execute(w, dArtistData); err != nil {
		ErrorHandler(
			w,
			r,
			http.StatusInternalServerError,
			"500 - Render Erelations  ror",
			"Unable to render artist page.",
		)
		return
	}
}
