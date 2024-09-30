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
	ErrLeetcodeValidationFailed   = errors.New("Leetcode validation failed")
	ErrLeetcodeUsernameInvalid    = errors.New("Leetcode username is invalid")
	ErrFetchingQuestion           = errors.New("error fetching question")
	ErrDeletingUserFailed         = errors.New("failed to delete question")
	ErrTransactionStart           = errors.New("could not start transaction")
	ErrTransactionCommit          = errors.New("could not commit transaction")
	ErrQueryExecution             = errors.New("could not execute query")
	ErrCheckRowsAffected          = errors.New("could not get rows affected")
	ErrNoRows                     = errors.New("no rows found")
	ErrRows                       = errors.New("rows error")
	ErrInternalServerError        = errors.New("internal server error")
	ErrInvalidCSVFormat           = errors.New("invalid CSV format, expected 7 columns")
	ErrInvalidQuestionID          = errors.New("invalid question ID")
	ErrInvalidQuestionDifficulty  = errors.New("invalid difficulty")
	ErrInvalidQuestionLink        = errors.New("invalid question link")
	ErrQuestionExists             = errors.New("question already exists in the database")
	ErrReadingCSVFile             = errors.New("error reading CSV file")
	ErrCSVFileOpening             = errors.New("error opening CSV file")
	ErrDbOperation                = errors.New("error in database operations")
)

var (
	CodeInvalidRequest   = 1100
	CodePermissionDenied = 2200
	CodeValidationError  = 3300
	CodeDbError          = 4400
	CodeUnexpectedError  = 9900
)

// JSONError creates a JSON error response
func JSONError(w http.ResponseWriter, message string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getHttpStatus(errorCode))
	jsonResponse := map[string]interface{}{"error_code": errorCode, "message": message}
	json.NewEncoder(w).Encode(jsonResponse)
}

// getHttpStatus maps custom error codes to HTTP status codes
func getHttpStatus(errorCode int) int {
	switch errorCode {
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodePermissionDenied:
		return http.StatusForbidden
	case CodeValidationError:
		return http.StatusUnprocessableEntity
	case CodeDbError:
		return http.StatusInternalServerError
	case CodeUnexpectedError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// AppError is a custom error type that includes an HTTP status code and a message
type AppError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func NewAppError(errorCode int, message string) *AppError {
	return &AppError{
		ErrorCode: errorCode,
		Message:   message,
	}
}

// ToJSON sends the AppError as a JSON response with the appropriate HTTP status code
func (e *AppError) ToJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(getHttpStatus(e.ErrorCode))
	json.NewEncoder(w).Encode(e)
}
