package errors

import (
	"fmt"
	"net/http"
)

// Error represents a custom application error with a message, code, and optional wrapped error.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"` // HTTP status code or custom error code
	Err     error  `json:"-"`    // Original error, not exposed via API
}

// New creates a new custom error.
func New(message string, code int, err error) *Error {
	return &Error{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.Err
}

// Is checks if the target error is of the same type or wraps it.
func Is(err error, target error) bool {
	if customErr, ok := err.(*Error); ok {
		return customErr.Message == target.Error() || Is(customErr.Err, target)
	}
	return err == target
}

// Common application errors
var (
	ErrNotFound            = New("resource not found", http.StatusNotFound, nil)
	ErrUnauthorized        = New("unauthorized", http.StatusUnauthorized, nil)
	ErrForbidden           = New("forbidden", http.StatusForbidden, nil)
	ErrBadRequest          = New("bad request", http.StatusBadRequest, nil)
	ErrInternalServerError = New("internal server error", http.StatusInternalServerError, nil)
	ErrConflict            = New("conflict", http.StatusConflict, nil)
	ErrValidation          = New("validation error", http.StatusUnprocessableEntity, nil)
	ErrServiceUnavailable  = New("service unavailable", http.StatusServiceUnavailable, nil)
)
