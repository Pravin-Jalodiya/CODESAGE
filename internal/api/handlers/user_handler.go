package handlers

import (
	"cli-project/internal/api/middleware"
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
	LeetcodeStats *models.LeetcodeStats `json:"leetcodeStats"`
	CodesageStats *models.CodesageStats `json:"codesageStats"`
}

// UserHandler is a struct for handling user-related requests.
type UserHandler struct {
	userService interfaces.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// GetUserByID handles the GET request for user profile by verifying the user ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	// Check if URL username matches token username
	if userMetaData.Username != username {
		errs.JSONError(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userMetaData.UserId.String())
	if err != nil {
		// If the user is not found (assuming the service returns a specific error for this case):
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUserProgress handles the GET request to fetch user's progress.
func (h *UserHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	userMetaData, ok := r.Context().Value("userMetaData").(middleware.UserMetaData)
	if !ok {
		errs.JSONError(w, "Could not retrieve user metadata", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	leetcodeStats, err := h.userService.GetUserLeetcodeStats(userMetaData.UserId.String())
	if err != nil {
		errs.JSONError(w, "Error fetching Leetcode stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	codesageStats, err := h.userService.GetUserCodesageStats(ctx, userMetaData.UserId.String())
	if err != nil {
		errs.JSONError(w, "Error fetching Codesage stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	progressResponse := UserProgressResponse{
		LeetcodeStats: leetcodeStats,
		CodesageStats: codesageStats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(progressResponse); err != nil {
		errs.JSONError(w, err.Error(), http.StatusInternalServerError)
	}
}
