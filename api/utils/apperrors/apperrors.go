package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

// set of valid error types
const (
	Authorization Type = "AUTHORIZATION"
	BadRequest    Type = "BADREQUEST"
	Conflict      Type = "CONFLICT"
	NotFound      Type = "NOTFOUND"
	Internal      Type = "INTERNAL"
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Error satisfies stander error interface
/*
 The error type in Go is implemented as the following interface:

 type error interface {
     Error() string
 }
*/

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case NotFound:
		return http.StatusNotFound
	case Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Status Check the run time error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

// Factories

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
	}
}
func NewConflict(reason string) *Error {
	return &Error{
		Type:    Conflict,
		Message: reason,
	}
}
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value %v not found", name, value),
	}
}
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("Internal Server Error"),
	}
}
