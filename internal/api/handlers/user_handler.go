package handlers

import (
	"cli-project/internal/api/middleware"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"cli-project/pkg/validation"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	LeetcodeID   string `json:"leetcodeId"`
	Organisation string `json:"organisation"`
	Country      string `json:"country"`
}

type UserProgressResponse struct {
	Code          int                   `json:"code"`
	Message       string                `json:"message"`
	LeetcodeStats *models.LeetcodeStats `json:"leetcodeStats"`
	CodesageStats *models.CodesageStats `json:"codesageStats"`
}

// UserHandler handles user-related requests.
type UserHandler struct {
	userService interfaces.UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserByID returns the user's profile.
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", errs.CodeInvalidRequest)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if userMetaData.Username != username {
		logger.Logger.Errorw("Unauthorized access: token does not match requested user", "method", r.Method, "user", username, "time", time.Now())
		errs.JSONError(w, "Unauthorized access: token does not match requested user", errs.CodePermissionDenied)
		return
	}

	user, err := u.userService.GetUserByID(r.Context(), userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Failed to fetch user by ID", "method", r.Method, "error", err, "time", time.Now())
		if errors.Is(err, errs.ErrUserNotFound) {
			errs.JSONError(w, "User not found", errs.CodeInvalidRequest)
		} else {
			errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
		}
		return
	}

	userResponse := UserResponse{
		Username:     user.Username,
		Name:         user.Name,
		Email:        user.Email,
		LeetcodeID:   user.LeetcodeID,
		Organisation: user.Organisation,
		Country:      user.Country,
	}

	response := struct {
		Code        int          `json:"code"`
		Message     string       `json:"message"`
		UserProfile UserResponse `json:"user_profile"`
	}{
		Code:        http.StatusOK,
		Message:     "User profile retrieved successfully",
		UserProfile: userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
	}
}

// GetUserProgress returns the user's progress.
func (u *UserHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", errs.CodeInvalidRequest)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if userMetaData.Username != username {
		logger.Logger.Errorw("Unauthorized access", "method", r.Method, "user", username, "time", time.Now())
		errs.JSONError(w, "Unauthorized access", errs.CodePermissionDenied)
		return
	}

	ctx := r.Context()
	leetcodeStats, err := u.userService.GetUserLeetcodeStats(userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Error fetching Leetcode stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Error fetching Leetcode stats: "+err.Error(), errs.CodeDbError)
		return
	}

	codesageStats, err := u.userService.GetUserCodesageStats(ctx, userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Error fetching Codesage stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Error fetching Codesage stats: "+err.Error(), errs.CodeDbError)
		return
	}

	progressResponse := UserProgressResponse{
		Code:          http.StatusOK,
		Message:       "Fetched user progress successfully",
		LeetcodeStats: leetcodeStats,
		CodesageStats: codesageStats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(progressResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
	}
}

// UpdateUserProfile updates the user's profile.
func (u *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", errs.CodeInvalidRequest)
		return
	}

	// Parse the update request
	var updateReq map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		logger.Logger.Errorw("Invalid request body", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Invalid request body", errs.CodeInvalidRequest)
		return
	}

	// Check if the body is empty
	if len(updateReq) == 0 {
		logger.Logger.Warnw("Empty request body", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Empty request body", errs.CodeInvalidRequest)
		return
	}

	// Prepare allowable updates
	allowedUpdates := map[string]bool{
		"username":     true,
		"password":     true,
		"name":         true,
		"email":        true,
		"organisation": true,
		"country":      true,
	}

	// Prepare restricted fields
	restrictedFieldsMap := map[string]bool{
		"id":          true,
		"role":        true,
		"last_seen":   true,
		"leetcode_id": true,
		"is_banned":   true,
	}

	// Collect updates
	updates := map[string]interface{}{}
	restrictedFields := []string{}
	invalidFields := []string{}

	for key, value := range updateReq {
		if allowedUpdates[key] {
			updates[key] = value
		} else if restrictedFieldsMap[key] {
			restrictedFields = append(restrictedFields, key)
		} else {
			invalidFields = append(invalidFields, key)
		}
	}

	// If there are invalid fields, return a bad request error
	if len(invalidFields) > 0 {
		message := "Found invalid field(s): " + strings.Join(invalidFields, ", ")
		logger.Logger.Warnw(message, "method", r.Method, "time", time.Now())
		errs.JSONError(w, message, errs.CodeInvalidRequest)
		return
	}

	// If there are restricted fields, return a forbidden error
	if len(restrictedFields) > 0 {
		message := "Cannot update restricted field(s): " + strings.Join(restrictedFields, ", ")
		logger.Logger.Warnw(message, "method", r.Method, "user", userMetaData.Username, "time", time.Now())
		errs.JSONError(w, message, errs.CodePermissionDenied)
		return
	}

	// Perform custom validations
	if username, ok := updates["username"].(string); ok {
		if !validation.ValidateUsername(username) {
			errs.NewAppError(errs.CodeValidationError, "Invalid username").ToJSON(w)
			logger.Logger.Errorw("Invalid username", "method", r.Method, "username", username, "time", time.Now())
			return
		}
	}
	if password, ok := updates["password"].(string); ok {
		if !validation.ValidatePassword(password) {
			errs.NewAppError(errs.CodeValidationError, "Invalid password").ToJSON(w)
			logger.Logger.Errorw("Invalid password", "method", r.Method, "time", time.Now())
			return
		}
	}
	if name, ok := updates["name"].(string); ok {
		if !validation.ValidateName(name) {
			errs.NewAppError(errs.CodeValidationError, "Invalid name").ToJSON(w)
			logger.Logger.Errorw("Invalid name", "method", r.Method, "time", time.Now())
			return
		}
	}
	if email, ok := updates["email"].(string); ok {
		isEmailValid, isReputable := validation.ValidateEmail(email)
		if !isEmailValid {
			errs.NewAppError(errs.CodeValidationError, "Invalid email format").ToJSON(w)
			logger.Logger.Errorw("Invalid email format", "method", r.Method, "email", email, "time", time.Now())
			return
		} else if !isReputable {
			errs.NewAppError(errs.CodeValidationError, "Unsupported email domain (use gmail, hotmail, outlook, watchguard or icloud)").ToJSON(w)
			logger.Logger.Errorw("Unsupported email domain", "method", r.Method, "email", email, "time", time.Now())
			return
		}
	}
	if organisation, ok := updates["organisation"].(string); ok {
		isOrgValid, orgErr := validation.ValidateOrganizationName(organisation)
		if !isOrgValid {
			errs.NewAppError(errs.CodeValidationError, orgErr.Error()).ToJSON(w)
			logger.Logger.Errorw("Invalid organization name", "method", r.Method, "organisation", organisation, "time", time.Now())
			return
		}
	}
	if country, ok := updates["country"].(string); ok {
		isCountryValid, countryErr := validation.ValidateCountryName(country)
		if !isCountryValid {
			errs.NewAppError(errs.CodeValidationError, countryErr.Error()).ToJSON(w)
			logger.Logger.Errorw("Invalid country name", "method", r.Method, "country", country, "time", time.Now())
			return
		}
	}

	// Ensure there is at least one valid field to update
	if len(updates) == 0 {
		logger.Logger.Warnw("No valid fields to update", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "No valid fields to update", errs.CodeInvalidRequest)
		return
	}

	// Perform the update
	err := u.userService.UpdateUser(r.Context(), userMetaData.UserId.String(), updates)
	if err != nil {
		logger.Logger.Errorw("Failed to update user", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to update user: "+err.Error(), errs.CodeDbError)
		return
	}

	response := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "User profile updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
	}
}

// UpdateUserProgress updates the user's progress.
func (u *UserHandler) UpdateUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", errs.CodeInvalidRequest)
		return
	}

	userID := userMetaData.UserId
	err := u.userService.UpdateUserProgress(r.Context(), userID)
	if err != nil {
		logger.Logger.Errorw("Failed to update user progress", "method", r.Method, "error", err, "time", time.Now())
		switch {
		case errors.Is(err, errs.ErrInvalidBodyError):
			errs.JSONError(w, "Invalid user ID", errs.CodeValidationError)
		case errors.Is(err, errs.ErrExternalAPI):
			errs.JSONError(w, "Error fetching data from LeetCode API", errs.CodeUnexpectedError)
		case errors.Is(err, errs.ErrDbError):
			errs.JSONError(w, "Database error updating user progress", errs.CodeDbError)
		default:
			errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
		}
		return
	}

	response := map[string]interface{}{"code": http.StatusOK, "message": "User progress updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
	}
}

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit int
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			logger.Logger.Errorw("Invalid limit parameter", "method", r.Method, "limit", limitStr, "error", err, "time", time.Now())
			errs.NewAppError(errs.CodeInvalidRequest, "Invalid limit: must be a positive number").ToJSON(w)
			return
		}
	}

	var offset int
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			logger.Logger.Errorw("Invalid offset parameter", "method", r.Method, "offset", offsetStr, "error", err, "time", time.Now())
			errs.NewAppError(errs.CodeInvalidRequest, "Invalid offset: must be a non-negative number").ToJSON(w)
			return
		}
	}

	ctx := r.Context()

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		logger.Logger.Errorw("Failed to fetch all users", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeDbError)
		return
	}

	if users == nil {
		users = []dto.StandardUser{}
	}

	totalUsers := len(users)

	if limitStr == "" {
		jsonResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Fetched users successfully",
			"users":   users,
			"total":   totalUsers,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
			logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
			errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
		}
		return
	}

	totalUsers = len(users)
	var paginatedUsers []dto.StandardUser

	if offset < totalUsers {
		end := offset + limit
		if end > totalUsers {
			end = totalUsers
		}
		paginatedUsers = users[offset:end]
	} else {
		paginatedUsers = []dto.StandardUser{}
	}

	jsonResponse := map[string]interface{}{
		"code":               http.StatusOK,
		"message":            "Fetched users successfully",
		"users":              paginatedUsers,
		"total":              totalUsers,
		"current_page_users": len(paginatedUsers),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), errs.CodeUnexpectedError)
	}
}

// GetPlatformStats returns platform stats.
func (u *UserHandler) GetPlatformStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformStats, err := u.userService.GetPlatformStats(ctx)
	if err != nil {
		logger.Logger.Errorw("Failed to fetch platform stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to fetch platform stats", errs.CodeDbError)
		return
	}

	jsonResponse := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Fetched platform stats successfully",
		"stats":   platformStats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to encode response", errs.CodeUnexpectedError)
	}
}

// UpdateUserBanState updates a user's ban state.
func (u *UserHandler) UpdateUserBanState(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		logger.Logger.Errorw("Missing 'username' query parameter", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Bad Request: 'username' query parameter is required", errs.CodeInvalidRequest)
		return
	}

	message, err := u.userService.UpdateUserBanState(r.Context(), username)
	if err != nil {
		logger.Logger.Errorw("Failed to update user ban state", "method", r.Method, "error", err, "time", time.Now())
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			errs.JSONError(w, "User not found", errs.CodeInvalidRequest)
		case errors.Is(err, errs.ErrInvalidParameterError):
			errs.JSONError(w, "Operation not allowed", errs.CodePermissionDenied)
		default:
			errs.JSONError(w, "Internal Server Error", errs.CodeUnexpectedError)
		}
		return
	}

	jsonResponse := map[string]interface{}{
		"code":    http.StatusOK,
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to encode response", errs.CodeUnexpectedError)
	}
}

// DeleteUser handles the HTTP request for deleting a user.
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		logger.Logger.Errorw("Missing username parameter", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Missing username parameter", errs.CodeInvalidRequest)
		return
	}

	err := u.userService.DeleteUser(r.Context(), username)
	if err != nil {
		logger.Logger.Errorw("Failed to delete user", "method", r.Method, "error", err, "time", time.Now())
		if errors.Is(err, errs.ErrUserNotFound) {
			errs.JSONError(w, "User not found", errs.CodeInvalidRequest)
		} else {
			errs.JSONError(w, "Failed to delete user", errs.CodeUnexpectedError)
		}
		return
	}

	response := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to encode response", errs.CodeUnexpectedError)
	}
}
