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
	config   *c.Config
	response utils.Response
	logger   utils.Logger
)

func init() {
	godotenv.Load()
	config = c.NewConfig()
	response = utils.NewResponse()
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

	r.NotFound(notFound)
	r.MethodNotAllowed(methodNotAllowed)

	static(r)

	r.Get("/", notFound)
	r.Get("/echo", echoWS)

	r.Mount("/api/user", userAPIRouter())

	logger.Yellow().Infof("Server is listening on http://localhost:%s", config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
