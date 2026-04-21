package handler

import (
	"html/template"
	"net/http"
)

type T_Error struct {
	Title       string
	StatusError int
	ErrorMsg    string
}

func ErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	title string,
	msg string,
) {
	data := T_Error{
		Title:       title,
		StatusError: status,
		ErrorMsg:    msg,
	}

	tmpl, err := template.ParseFiles("template/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}
