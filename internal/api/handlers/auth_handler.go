package handlers

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/globals"
	"cli-project/pkg/logger"
	"cli-project/pkg/utils"
	"cli-project/pkg/validation"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator"
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
	var requestBody struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		Organisation string `json:"organisation"`
		Country      string `json:"country"`
		LeetcodeID   string `json:"leetcode_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		errs.NewAppError(errs.CodeInvalidRequest, "Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	user := models.StandardUser{
		User: models.User{
			Username:     requestBody.Username,
			Password:     requestBody.Password,
			Name:         requestBody.Name,
			Email:        requestBody.Email,
			Organisation: requestBody.Organisation,
			Country:      requestBody.Country,
		},
		LeetcodeID:      requestBody.LeetcodeID,
		QuestionsSolved: []string{},
		LastSeen:        time.Time{},
	}

	// Perform custom validations
	if !validation.ValidateUsername(user.Username) {
		errs.NewAppError(errs.CodeValidationError, "Invalid username").ToJSON(w)
		logger.Logger.Errorw("Invalid username", "method", r.Method, "username", user.Username, "time", time.Now())
		return
	}
	if !validation.ValidatePassword(user.Password) {
		errs.NewAppError(errs.CodeValidationError, "Invalid password").ToJSON(w)
		logger.Logger.Errorw("Invalid password", "method", r.Method, "time", time.Now())
		return
	}
	if !validation.ValidateName(user.Name) {
		errs.NewAppError(errs.CodeValidationError, "Invalid name").ToJSON(w)
		logger.Logger.Errorw("Invalid name", "method", r.Method, "time", time.Now())
		return
	}

	isEmailValid, isReputable := validation.ValidateEmail(user.Email)
	if !isEmailValid {
		errs.NewAppError(errs.CodeValidationError, "Invalid email format").ToJSON(w)
		logger.Logger.Errorw("Invalid email format", "method", r.Method, "email", user.Email, "time", time.Now())
		return
	} else if !isReputable {
		errs.NewAppError(errs.CodeValidationError, "Unsupported email domain (use gmail, hotmail, outlook, watchguard or icloud)").ToJSON(w)
		logger.Logger.Errorw("Unsupported email domain", "method", r.Method, "email", user.Email, "time", time.Now())
		return
	}

	isOrgValid, orgErr := validation.ValidateOrganizationName(user.Organisation)
	if !isOrgValid {
		errs.NewAppError(errs.CodeValidationError, orgErr.Error()).ToJSON(w)
		logger.Logger.Errorw("Invalid organization name", "method", r.Method, "organisation", user.Organisation, "time", time.Now())
		return
	}

	isCountryValid, countryErr := validation.ValidateCountryName(user.Country)
	if !isCountryValid {
		errs.NewAppError(errs.CodeValidationError, countryErr.Error()).ToJSON(w)
		logger.Logger.Errorw("Invalid country name", "method", r.Method, "country", user.Country, "time", time.Now())
		return
	}

	err = a.authService.Signup(r.Context(), &user)
	if err != nil {
		if errors.Is(err, errs.ErrUserNameAlreadyExists) {
			errs.NewAppError(errs.CodeValidationError, "User already exists").ToJSON(w)
		} else if errors.Is(err, errs.ErrEmailAlreadyExists) {
			errs.NewAppError(errs.CodeValidationError, "Email already registered").ToJSON(w)
		} else if errors.Is(err, errs.ErrLeetcodeIDAlreadyExists) {
			errs.NewAppError(errs.CodeValidationError, "LeetcodeID already registered").ToJSON(w)
		} else if errors.Is(err, errs.ErrLeetcodeUsernameInvalid) {
			errs.NewAppError(errs.CodeValidationError, "Invalid leetcode id").ToJSON(w)
		} else {
			errs.NewAppError(errs.CodeUnexpectedError, "Signup failed").ToJSON(w)
		}
		logger.Logger.Errorw("Signup failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"message": "User successfully registered",
		"code":    http.StatusOK,
		"user_info": map[string]any{
			"username": user.Username,
			"role":     user.Role,
		},
	}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Signup Successful", "method", r.Method, "username", user.Username, "time", time.Now())
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		errs.NewAppError(errs.CodeInvalidRequest, "Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Validate the request body
	if !validation.ValidateUsername(requestBody.Username) {
		errs.NewAppError(errs.CodeValidationError, "Invalid username").ToJSON(w)
		logger.Logger.Errorw("Invalid username", "method", r.Method, "username", requestBody.Username, "time", time.Now())
		return
	}
	if !validation.ValidatePassword(requestBody.Password) {
		errs.NewAppError(errs.CodeValidationError, "Invalid username or password").ToJSON(w)
		logger.Logger.Errorw("Invalid password", "method", r.Method, "username", requestBody.Username, "time", time.Now())
		return
	}

	user, err := a.authService.Login(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPassword) {
			errs.NewAppError(errs.CodeInvalidRequest, "Invalid username or password").ToJSON(w)
		} else if errors.Is(err, errs.ErrUserNotFound) {
			errs.NewAppError(errs.CodeInvalidRequest, "User not found").ToJSON(w)
		} else {
			errs.NewAppError(errs.CodeUnexpectedError, "Login failed").ToJSON(w)
		}
		logger.Logger.Errorw("Authentication failed", "method", r.Method, "error", err, "username", requestBody.Username, "time", time.Now())
		return
	}

	token, err := utils.CreateJwtToken(user.Username, user.ID, user.Role, user.IsBanned)
	if err != nil {
		errs.NewAppError(errs.CodeUnexpectedError, "Failed to generate token").ToJSON(w)
		logger.Logger.Errorw("Failed to generate token", "method", r.Method, "username", requestBody.Username, "error", err, "time", time.Now())
		return
	}

	globals.ActiveUserID = user.ID

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{"code": http.StatusOK, "message": "Login successful", "token": token, "role": user.Role}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Login Successful", "method", r.Method, "username", requestBody.Username, "time", time.Now())
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := a.authService.Logout(r.Context())
	if err != nil {
		errs.NewAppError(errs.CodeUnexpectedError, err.Error()).ToJSON(w)
		logger.Logger.Errorw("Logout failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	globals.ActiveUserID = ""

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{"code": http.StatusOK, "message": "Logout successful"}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Logout Successful", "method", r.Method, "time", time.Now())
}
