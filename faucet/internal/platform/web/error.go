package web

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Err    error
	Status int
}

type HTMLError struct {
	Err    error
	Status int
}

type ErrorResponse struct {
	Error string `json:"message,omitempty"`
}

func (err *Error) Error() string {
	return err.Err.Error()
}

func (err *HTMLError) Error() string {
	return err.Err.Error()
}

func NewRequestError(err error, status int) error {
	return &Error{err, status}
}

func NewResponseError(err error, status int) error {
	return &Error{err, status}
}

func NewHTMLError(err error, status int) error {
	return &HTMLError{err, status}
}

func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
