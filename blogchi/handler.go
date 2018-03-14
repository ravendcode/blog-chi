package main

import (
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.SendFile(w, r, "index.html")
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 405).JSON()
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	response.SendFile(w, r, "echo.html")
}
