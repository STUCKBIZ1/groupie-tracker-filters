package cmd

import (
	"fmt"
	"net/http"

	"gr-tr/src/config"
	route "gr-tr/src/routes"
)

func Server() {
	route.Route()
	config.RegisterAssets()
	fmt.Println("Listning to http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Error: Runing server")
	}
}
