package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type user struct {
	ID         int    `json:"id,omitempty"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	AccesToken string `json:"accessToken"`
}

type users []*user

var storage = users{{
	ID:         1,
	Username:   "root",
	Email:      "root@email.com",
	Password:   "qwerty",
	AccesToken: "",
}}

func userAPIRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", userGetAll)
	r.Post("/", userCreateOne)
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Use(userCtx)
		r.Get("/", userFindOne)
		r.Patch("/", userUpdateOne)
		r.Delete("/", userDeleteOne)
	})
	return r
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stringID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(stringID)
		if err != nil {
			response.Send(w, 500, err).JSON()
			return
		}
		for _, user := range storage {
			if user.ID == id {
				ctx := context.WithValue(r.Context(), userKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		response.Send(w, 404, "Not Found User").JSON()
	})
}

func userGetAll(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 200, "userGetAll").JSON(storage)
}

func userCreateOne(w http.ResponseWriter, r *http.Request) {
	user := &user{len(storage) + 1, "John", "john@email.com", "qwerty", ""}
	storage = append(storage, user)
	response.Send(w, 201, "userCreateOne").JSON(user)
}

func userFindOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*user)
	if !ok {
		response.Send(w, 422).JSON()
		return
	}
	response.Send(w, 200, "userFindOne").JSON(user)
}

func userUpdateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*user)
	if !ok {
		response.Send(w, 422).JSON()
		return
	}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.Send(w, 500, err).JSON()
		return
	}
	response.Send(w, 200, "userUpdateOne").JSON(user)
}

func userDeleteOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*user)
	if !ok {
		// http.Error(w, http.StatusText(422), 422)
		response.Send(w, 422).JSON()
		return
	}
	for index, userStorage := range storage {
		if user.ID == userStorage.ID {
			storage = append(storage[:index], storage[index+1:]...)
			break
		}
	}
	response.Send(w, 200, "userDeleteOne").JSON(user)
}
