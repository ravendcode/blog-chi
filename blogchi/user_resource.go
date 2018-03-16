package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var x = "aaaa"

type password string

// func (password) MarshalJSON() ([]byte, error) {
// 	return []byte(`""`), nil
// }

// User model
type User struct {
	ID         int      `json:"id,omitempty"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Password   password `json:"password"`
	AccesToken string   `json:"accessToken"`
}

// Users slice User
type Users []*User

var storage = Users{{
	ID:         1,
	Username:   "root",
	Email:      "root@email.com",
	Password:   "qwerty",
	AccesToken: "",
}}

type userResource struct{}

func (rs userResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..
	r.Get("/", rs.getAll)
	r.Post("/", rs.createOne)
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Use(rs.userCtx)
		r.Get("/", rs.findOne)
		r.Patch("/", rs.updateOne)
		r.Delete("/", rs.deleteOne)
	})
	return r
}

func (rs userResource) userCtx(next http.Handler) http.Handler {
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

func (rs userResource) getAll(w http.ResponseWriter, r *http.Request) {
	response.Send(w, 200, "getAll").JSON(storage)
}

func (rs userResource) createOne(w http.ResponseWriter, r *http.Request) {
	rules := map[string]string{
		"username": "required|len(2,32)|forbiddenusernames",
		"email":    "required|email",
		"password": "required|len(6,32)",
	}
	newUser := new(User)
	if response.Validate(w, r, rules, newUser) {
		storage = append(storage, newUser)
		response.Send(w, 201, "createOne").JSON(newUser)
	}
}

func (rs userResource) findOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		response.Send(w, 422).JSON()
		return
	}
	response.Send(w, 200, "findOne").JSON(user)
}

func (rs userResource) updateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		response.Send(w, 422).JSON()
		return
	}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.Send(w, 500, err).JSON()
		return
	}
	rules := map[string]string{
		"username": "required|len(2,32)|forbiddenusernames",
		"email":    "required|email",
		"password": "required|len(6,32)",
	}
	ok, errors := validator.Validate(rules, user)
	if !ok {
		response.Send(w, 400, "Validation Error", errors).JSON()
		return
	}

	response.Send(w, 200, "updateOne").JSON(user)
}

func (rs userResource) deleteOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userKey).(*User)
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
	response.Send(w, 200, "deleteOne").JSON(user)
}
