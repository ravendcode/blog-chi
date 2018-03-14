package main

import (
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	response.SendFile(w, r, "index.html")
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 405).JSON()
}

func echo(w http.ResponseWriter, r *http.Request) {
	response.SendFile(w, r, "echo.html")
}
