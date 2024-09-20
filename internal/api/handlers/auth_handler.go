package handlers

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	"cli-project/pkg/globals"
	"cli-project/pkg/logger"
	"cli-project/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"net/http"
	"time"
)

var validate *validator.Validate

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	validate = validator.New()
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.StandardUser

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errs.NewInvalidRequestBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Validate the user information
	if err := validate.Struct(user); err != nil {
		errs.NewBadRequestError("Invalid request parameters").ToJSON(w)
		logger.Logger.Errorw("Validation failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Call the service layer for signup
	err = a.authService.Signup(r.Context(), &user)
	if err != nil {
		if errors.Is(err, errs.ErrUserNameAlreadyExists) {
			errs.NewConflictError("User already exists").ToJSON(w)
		} else if errors.Is(err, errs.ErrEmailAlreadyExists) {
			errs.NewConflictError("Email already registered").ToJSON(w)
		} else if errors.Is(err, errs.ErrLeetcodeIDAlreadyExists) {
			errs.NewConflictError("LeetcodeID already registered").ToJSON(w)
		} else {
			errs.NewInternalServerError("Signup failed").ToJSON(w)
		}
		logger.Logger.Errorw("Signup failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"message": "User successfully registered",
		"user_info": map[string]string{
			"username":     user.StandardUser.Username,
			"organisation": user.StandardUser.Organisation,
			"country":      user.StandardUser.Country,
		},
	}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Signup Successful", "method", r.Method, "username", user.StandardUser.Username, "time", time.Now())
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Define the request body structure
	var requestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		errs.NewInvalidRequestBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Validate the request body
	err = validate.Struct(requestBody)
	if err != nil {
		errs.NewInvalidRequestBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Validation error", "method", r.Method, "error", err, "request", requestBody, "time", time.Now())
		return
	}

	// Call the service layer for login
	user, err := a.authService.Login(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) || errors.Is(err, errs.ErrInvalidPassword) {
			errs.NewAuthenticationError("Invalid username or password").ToJSON(w)
		} else {
			errs.NewInternalServerError("Login failed").ToJSON(w)
		}
		logger.Logger.Errorw("Authentication failed", "method", r.Method, "error", err, "username", requestBody.Username, "time", time.Now())
		return
	}

	// Create a JWT token
	token, err := utils.CreateJwtToken(user.StandardUser.Username, user.StandardUser.ID, user.StandardUser.Role)
	if err != nil {
		errs.NewInternalServerError("Failed to generate token").ToJSON(w)
		return
	}

	globals.ActiveUserID = user.StandardUser.ID

	// Return the token as a JSON
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]string{"token": token, "role": user.StandardUser.Role}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Login Successful", "method", r.Method, "request", requestBody, "time", time.Now())
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Call the service layer for logout
	err := a.authService.Logout(r.Context())
	if err != nil {
		errs.NewInternalServerError(err.Error()).ToJSON(w)
		logger.Logger.Errorw("Logout failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Invalidate the token (clear the token or blacklist it as needed)
	globals.ActiveUserID = ""

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]string{"message": "Logout successful"}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Logout Successful", "method", r.Method, "time", time.Now())
}
