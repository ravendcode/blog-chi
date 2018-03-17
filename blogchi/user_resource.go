package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// type password string

// func (password) MarshalJSON() ([]byte, error) {
// 	return []byte(`""`), nil
// }

// User model
type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	AccesToken string `json:"accessToken"`
}

// MarshalJSON custom fields output
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID         int    `json:"id,omitempty"`
		Username   string `json:"username,omitempty"`
		Email      string `json:"email,omitempty"`
		AccesToken string `json:"accessToken,omitempty"`
	}{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		AccesToken: u.AccesToken,
	})
}

// UnmarshalJSON custom fields add
func (u *User) UnmarshalJSON(b []byte) error {
	decoded := new(struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	})
	err := json.Unmarshal(b, decoded)
	if err == nil {
		u.Username = decoded.Username
		u.Email = decoded.Email
		u.Password = decoded.Password
	}
	return err
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
	newUser := new(User)
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		response.Send(w, 500, err).JSON()
		return
	}
	// var interfaceSlice = make([]interface{}, len(storage))
	// for i, u := range storage {
	// 	interfaceSlice[i] = u
	// }
	// validator.Unique("username", newUser, interfaceSlice)
	// validator.Unique("email", newUser, interfaceSlice)
	for _, u := range storage {
		if u.Username == newUser.Username {
			validator.AddError("username", "username is unique")
			break
		}
	}
	for _, u := range storage {
		if u.Email == newUser.Email {
			validator.AddError("email", "email is unique")
			break
		}
	}
	rules := map[string]interface{}{
		"username": "required|len(2,32)|forbiddenusernames",
		"email":    "required|email",
		"password": "required|len(6,32)",
	}
	if response.Validate(w, validator, rules, newUser) {
		newUser.ID = len(storage) + 1
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
	newUser := *user
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		response.Send(w, 500, err).JSON()
		return
	}
	for _, u := range storage {
		if u.Username == newUser.Username && u.ID != newUser.ID {
			validator.AddError("username", "username is unique")
			break
		}
	}
	for _, u := range storage {
		if u.Email == newUser.Email && u.ID != newUser.ID {
			validator.AddError("email", "email is unique")
			break
		}
	}
	rules := map[string]interface{}{
		"username": "required|len(2,32)|forbiddenusernames",
		"email":    "required|email",
		"password": "required|len(6,32)",
	}
	errors := validator.Validate(rules, &newUser)
	if errors != nil {
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
