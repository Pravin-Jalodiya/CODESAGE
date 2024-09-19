package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Predefined errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// AppError is a custom error type that includes an HTTP status code and a message
type AppError struct {
	error   `json:"error,omitempty"`
	Code    int    `json:"status_code,omitempty"`
	Message string `json:"message"`
}

// Error constructors
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewInvalidParameterError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewInvalidBodyError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewInvalidRequestMethodError(message string) *AppError {
	return &AppError{
		Code:    http.StatusMethodNotAllowed,
		Message: message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// ToJSON sends the AppError as a JSON response with the appropriate HTTP status code
func (e *AppError) ToJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		return
	}
}
