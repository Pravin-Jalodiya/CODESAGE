package handlers

import (
	"cli-project/internal/api/middleware"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	errs "cli-project/pkg/errors"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"net/http"
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
	validate    *validator.Validate
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// GetUserByID returns the user's profile.
func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract user metadata from the request context
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	// Extract the username from the request parameters
	vars := mux.Vars(r)
	username := vars["username"]

	// Fetch user by ID from the service layer
	user, err := u.userService.GetUserByID(r.Context(), userMetaData.UserId.String())
	if err != nil {
		// If user is not found, return the appropriate error
		if errors.Is(err, errs.ErrUserNotFound) {
			errs.JSONError(w, "User not found", http.StatusNotFound)
		} else {
			errs.JSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Proceed only if the token username matches the requested username
	if userMetaData.Username != username {
		errs.JSONError(w, "Unauthorized access: token does not match requested user", http.StatusUnauthorized)
		return
	}

	// Create the user response object
	userResponse := UserResponse{
		Username:     user.StandardUser.Username,
		Name:         user.StandardUser.Name,
		Email:        user.StandardUser.Email,
		LeetcodeID:   user.LeetcodeID,
		Organisation: user.StandardUser.Organisation,
		Country:      user.StandardUser.Country,
	}

	// Create a response wrapper
	response := struct {
		Code        int          `json:"code"`
		Message     string       `json:"message"`
		UserProfile UserResponse `json:"user_profile"`
	}{
		Code:        http.StatusOK,
		Message:     "User profile retrieved successfully",
		UserProfile: userResponse,
	}

	// Respond with the user profile data in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserProgress returns the user's progress.
func (u *UserHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if userMetaData.Username != username {
		errs.JSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	leetcodeStats, err := u.userService.GetUserLeetcodeStats(userMetaData.UserId.String())
	if err != nil {
		errs.JSONError(w, "Error fetching Leetcode stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	codesageStats, err := u.userService.GetUserCodesageStats(ctx, userMetaData.UserId.String())
	if err != nil {
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
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateUserProgress updates the user's progress.
func (u *UserHandler) UpdateUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	userID := userMetaData.UserId
	err := u.userService.UpdateUserProgress(r.Context(), userID)
	if err != nil {
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
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAllUsers returns all users.
func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []dto.StandardUser{}
	}

	jsonResponse := map[string]any{
		"code":    http.StatusOK,
		"message": "Fetched users successfully",
		"users":   users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

// GetPlatformStats returns platform stats.
func (u *UserHandler) GetPlatformStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	platformStats, err := u.userService.GetPlatformStats(ctx)
	if err != nil {
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
		errs.JSONError(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateUserBanState updates a user's ban state.
func (u *UserHandler) UpdateUserBanState(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		errs.JSONError(w, "Bad Request: 'username' query parameter is required", http.StatusBadRequest)
		return
	}

	message, err := u.userService.UpdateUserBanState(r.Context(), username)
	if err != nil {
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
	json.NewEncoder(w).Encode(jsonResponse)
}
