package handlers

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"cli-project/pkg/utils"
	"encoding/json"
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
	// Signup logic here...
}

func (a *AuthHandler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Define the request body structure
	var requestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Validate the request body
	err = validate.Struct(requestBody)
	if err != nil {
		errs.NewInvalidBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Validation error", "method", r.Method, "error", err, "request", requestBody, "time", time.Now())
		return
	}

	// Call the service layer for login
	user, err := a.authService.Login(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		errs.NewAuthenticationError("Invalid username or password").ToJSON(w)
		logger.Logger.Errorw("Authentication failed", "method", r.Method, "error", err, "username", requestBody.Username, "time", time.Now())
		return
	}

	// create a jwt
	token, err := utils.CreateJwtToken(user.StandardUser.Username, user.StandardUser.ID, user.StandardUser.Role)
	if err != nil {
		errs.NewInternalServerError("Failed to generate token").ToJSON(w)
		return
	}
	// Return the token as a JSON
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]string{"token": token, "role": user.StandardUser.Role}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Login Successful", "method", r.Method, "request", requestBody, "time", time.Now())
}

func (a *AuthHandler) AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Admin login logic here...
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Logout logic here...
}
