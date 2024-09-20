package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUserNameAlreadyExists      = errors.New("username already exists")
	ErrLeetcodeIDAlreadyExists    = errors.New("leetcode ID already exists")
	ErrDbError                    = errors.New("db error")
	ErrInvalidParameterError      = errors.New("invalid parameter")
	ErrInvalidBodyError           = errors.New("invalid request body")
	ErrExternalAPI                = errors.New("external API error")
	ErrUserNotFound               = errors.New("user not found")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrEmailAlreadyExists         = errors.New("email already exists")
	ErrUserCreationFailed         = errors.New("failed to create user")
	ErrFetchingUserFailed         = errors.New("failed to fetch user")
	ErrFetchingUsersFailed        = errors.New("failed to fetch users")
	ErrUpdatingUserProgressFailed = errors.New("failed to update user progress")
	ErrUpdatingUserDetailsFailed  = errors.New("failed to update user details")
	ErrBanningUserFailed          = errors.New("failed to ban user")
	ErrUnbanningUserFailed        = errors.New("failed to unban user")
	ErrDatabaseConnection         = errors.New("failed to connect to database")
	ErrTransactionStart           = errors.New("could not start transaction")
	ErrTransactionCommit          = errors.New("could not commit transaction")
	ErrQueryExecution             = errors.New("could not execute query")
	ErrCheckRowsAffected          = errors.New("could not get rows affected")
	ErrNoRows                     = errors.New("no rows found")
	ErrRows                       = errors.New("rows error")
	ErrInternalServerError        = errors.New("internal server error")
)

// AppError is a custom error type that includes an HTTP status code and a message
type AppError struct {
	error   `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// Error constructors
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
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

func NewInvalidRequestBodyError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// ToJSON sends the AppError as a JSON response with the appropriate HTTP status code
func (e *AppError) ToJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
}
