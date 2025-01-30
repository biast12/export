package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Err         error
	StatusCode  int
	UserMessage string
}

func NewError(err error, statusCode int, userMessage string) *Error {
	return &Error{
		Err:         err,
		StatusCode:  statusCode,
		UserMessage: userMessage,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Write(w http.ResponseWriter) {
	response := ErrorResponse{
		Error: e.UserMessage,
	}

	w.WriteHeader(e.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
