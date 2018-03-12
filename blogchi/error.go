package main

import (
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
	// response.Send(w, 404).JSON()
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 405).JSON()
}
