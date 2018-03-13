package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

func static(r *chi.Mux) {
	workDir, _ := os.Getwd()
	fmt.Println(1, http.Dir(filepath.Join(workDir, "node_modules")))
	fileServer(r, "/static", http.Dir(filepath.Join(workDir, "static")))
	fileServer(r, "/node_modules", http.Dir(filepath.Join(workDir, "node_modules")))
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if _, err := os.Stat(fmt.Sprintf("%s", root) + r.RequestURI); os.IsNotExist(err) {
		// 	notFound(w, r)
		// } else {
		// 	fs.ServeHTTP(w, r)
		// }
		fs.ServeHTTP(w, r)
	}))
}
