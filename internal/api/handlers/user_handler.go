package handlers

import (
	"cli-project/internal/api/middleware"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if userMetaData.Username != username {
		logger.Logger.Errorw("Unauthorized access: token does not match requested user", "method", r.Method, "user", username, "time", time.Now())
		errs.JSONError(w, "Unauthorized access: token does not match requested user", http.StatusUnauthorized)
		return
	}

	user, err := u.userService.GetUserByID(r.Context(), userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Failed to fetch user by ID", "method", r.Method, "error", err, "time", time.Now())
		if errors.Is(err, errs.ErrUserNotFound) {
			errs.JSONError(w, "User not found", http.StatusNotFound)
		} else {
			errs.JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userResponse := UserResponse{
		Username:     user.StandardUser.Username,
		Name:         user.StandardUser.Name,
		Email:        user.StandardUser.Email,
		LeetcodeID:   user.LeetcodeID,
		Organisation: user.StandardUser.Organisation,
		Country:      user.StandardUser.Country,
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
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserProgress returns the user's progress.
func (u *UserHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if userMetaData.Username != username {
		logger.Logger.Errorw("Unauthorized access", "method", r.Method, "user", username, "time", time.Now())
		errs.JSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	leetcodeStats, err := u.userService.GetUserLeetcodeStats(userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Error fetching Leetcode stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Error fetching Leetcode stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	codesageStats, err := u.userService.GetUserCodesageStats(ctx, userMetaData.UserId.String())
	if err != nil {
		logger.Logger.Errorw("Error fetching Codesage stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Error fetching Codesage stats: "+err.Error(), http.StatusInternalServerError)
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
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUserProgress updates the user's progress.
func (u *UserHandler) UpdateUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		logger.Logger.Errorw("Could not retrieve user metadata", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	userID := userMetaData.UserId
	err := u.userService.UpdateUserProgress(r.Context(), userID)
	if err != nil {
		logger.Logger.Errorw("Failed to update user progress", "method", r.Method, "error", err, "time", time.Now())
		switch {
		case errors.Is(err, errs.ErrInvalidBodyError):
			errs.JSONError(w, "Invalid user ID", http.StatusBadRequest)
		case errors.Is(err, errs.ErrExternalAPI):
			errs.JSONError(w, "Error fetching data from LeetCode API", http.StatusInternalServerError)
		case errors.Is(err, errs.ErrDbError):
			errs.JSONError(w, "Database error updating user progress", http.StatusInternalServerError)
		default:
			errs.JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]any{"code": http.StatusOK, "message": "User progress updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
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
			errs.NewBadRequestError("Invalid limit: must be a positive number").ToJSON(w)
			return
		}
	}

	var offset int
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			logger.Logger.Errorw("Invalid offset parameter", "method", r.Method, "offset", offsetStr, "error", err, "time", time.Now())
			errs.NewBadRequestError("Invalid offset: must be a non-negative number").ToJSON(w)
			return
		}
	}

	ctx := r.Context()

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		logger.Logger.Errorw("Failed to fetch all users", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []dto.StandardUser{}
	}

	totalUsers := len(users)

	if limitStr == "" {
		jsonResponse := map[string]any{
			"code":    http.StatusOK,
			"message": "Fetched users successfully",
			"users":   users,
			"total":   totalUsers,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
			logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
			errs.JSONError(w, err.Error(), http.StatusInternalServerError)
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

	jsonResponse := map[string]any{
		"code":               http.StatusOK,
		"message":            "Fetched users successfully",
		"users":              paginatedUsers,
		"total":              totalUsers,
		"current_page_users": len(paginatedUsers),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetPlatformStats returns platform stats.
func (u *UserHandler) GetPlatformStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformStats, err := u.userService.GetPlatformStats(ctx)
	if err != nil {
		logger.Logger.Errorw("Failed to fetch platform stats", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to fetch platform stats", http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]any{
		"code":    http.StatusOK,
		"message": "Fetched platform stats successfully",
		"stats":   platformStats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateUserBanState updates a user's ban state.
func (u *UserHandler) UpdateUserBanState(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		logger.Logger.Errorw("Missing 'username' query parameter", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Bad Request: 'username' query parameter is required", http.StatusBadRequest)
		return
	}

	message, err := u.userService.UpdateUserBanState(r.Context(), username)
	if err != nil {
		logger.Logger.Errorw("Failed to update user ban state", "method", r.Method, "error", err, "time", time.Now())
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			errs.JSONError(w, "User not found", http.StatusNotFound)
		case errors.Is(err, errs.ErrInvalidParameterError):
			errs.JSONError(w, "Operation not allowed", http.StatusForbidden)
		default:
			errs.JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	jsonResponse := map[string]any{
		"code":    http.StatusOK,
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Failed to encode response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteUser handles the HTTP request for deleting a user
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		logger.Logger.Errorw("Missing username parameter", "method", r.Method, "time", time.Now())
		errs.JSONError(w, "Missing username parameter", http.StatusBadRequest)
		return
	}

	err := u.userService.DeleteUser(r.Context(), username)
	if err != nil {
		logger.Logger.Errorw("Failed to delete user", "method", r.Method, "error", err, "time", time.Now())
		if errors.Is(err, errs.ErrUserNotFound) {
			errs.JSONError(w, "User not found", http.StatusNotFound)
		} else {
			errs.JSONError(w, "Failed to delete user", http.StatusInternalServerError)
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
		errs.JSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
