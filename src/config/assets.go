package config

import (
	"net/http"

	handlers "gr-tr/src/handlers"
)

func RegisterAssets() {
	fs := http.Dir("assets")
	fileServer := http.FileServer(fs)

	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		// Block "/assets" and "/assets/"
		if r.URL.Path == "/assets" || r.URL.Path == "/assets/" {
			handlers.ErrorHandler(
				w,
				r,
				http.StatusNotFound,
				"404 - Page Not Found",
				"The page you are looking for does not exist.",
			)
			return
		}

		// Serve the file
		http.StripPrefix("/assets/", fileServer).ServeHTTP(w, r)
	})
}
