package errors

import (
	"net/http"
)

type RestError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *RestError) Error() string {
	return e.Message
}

func NewDatabaseError() *RestError {
	return NewInternalServerError("database error")
}
func NewBadRequestError(msg string) *RestError {
	return &RestError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) *RestError {
	return &RestError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func NewInternalServerError(msg string) *RestError {
	return &RestError{
		Message: msg,
		Code:    http.StatusInternalServerError,
	}
}
