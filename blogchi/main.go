package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	c "github.com/ravendcode/blog-chi/blogchi/config"
	"github.com/ravendcode/blog-chi/blogchi/utils"
)

var (
	config    *c.Config
	response  utils.Response
	validator utils.Validator
	logger    utils.Logger
)

func init() {
	godotenv.Load()
	config = c.NewConfig()
	response = utils.NewResponse()
	validator = utils.NewValidator()
	logger = utils.NewLogger()
	logger.Yellow().Info("ENV is", config.Env)
}

// func loggerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	static(r)

	r.Get("/", notFoundHandler)
	r.Get("/echo", echoHandler)
	r.Get("/echows", echoWS)

	r.Mount("/api/user", userResource{}.Routes())

	logger.Yellow().Infof("Server is listening on http://localhost:%s", config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
