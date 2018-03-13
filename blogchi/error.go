package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "static/index.html")
	workDir, _ := os.Getwd()
	fmt.Println(111, fmt.Sprintf("%s/index.html", http.Dir(filepath.Join(workDir, "static"))))
	http.ServeFile(w, r, fmt.Sprintf("%s/index.html", http.Dir(filepath.Join(workDir, "static"))))
	// response.Send(w, 404).JSON()
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 405).JSON()
}
