package route

import (
	"net/http"

	handlers "gr-tr/src/handlers"
)

func Route() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)
}
