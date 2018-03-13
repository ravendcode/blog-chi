package utils

import (
	"encoding/json"
	"net/http"
)

// Response interface for response
// response.Send(w, status, message).JSON(data)
type Response interface {
	Send(w http.ResponseWriter, statusCode int, message ...string) *response
	JSON(data ...interface{})
}

type response struct {
	StatusCode     int                 `json:"statusCode,omitempty"`
	Status         string              `json:"status,omitempty"`
	Message        string              `json:"message,omitempty"`
	Data           interface{}         `json:"data,omitempty"`
	ResponseWriter http.ResponseWriter `json:"-"`
}

func (r *response) Send(w http.ResponseWriter, statusCode int, message ...string) *response {
	r.ResponseWriter = w
	r.StatusCode = statusCode

	// switch r.StatusCode {
	// case 200:
	// 	r.Status = "OK"
	// case 201:
	// 	r.Status = "Created"
	// case 400:
	// 	r.Status = "BadRequest"
	// case 401:
	// 	r.Status = "Unauthorized"
	// case 403:
	// 	r.Status = "Forbidden"
	// case 404:
	// 	r.Status = "NotFound"
	// case 422:
	// 	r.Status = "UnprocessableEntity"
	// case 500:
	// 	r.Status = "InternalServer"
	// }

	r.Status = http.StatusText(statusCode)
	if len(message) != 0 {
		r.Message = message[0]
	} else {
		r.Message = http.StatusText(statusCode)
	}
	return r
}

func (r *response) JSON(data ...interface{}) {
	if len(data) != 0 {
		r.Data = data[0]
	}

	r.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.ResponseWriter.WriteHeader(r.StatusCode)

	if err := json.NewEncoder(r.ResponseWriter).Encode(r); err != nil {
		r.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		http.Error(r.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
	r.clean()
}

func (r *response) clean() {
	r.StatusCode = 0
	r.Status = ""
	r.Message = ""
	r.Data = nil
	r.ResponseWriter = nil
}

// NewResponse create new Response
func NewResponse() Response {
	return new(response)
}