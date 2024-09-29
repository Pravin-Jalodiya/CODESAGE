package handlers

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
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
	var RequestBody struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		Organization string `json:"organisation"`
		Country      string `json:"country"`
		LeetcodeID   string `json:"leetcode_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		errs.NewInvalidRequestBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	user := models.StandardUser{
		StandardUser: models.User{
			Username:     RequestBody.Username,
			Password:     RequestBody.Password,
			Name:         RequestBody.Name,
			Email:        RequestBody.Email,
			Organisation: RequestBody.Organization,
			Country:      RequestBody.Country,
		},
		LeetcodeID:      RequestBody.LeetcodeID,
		QuestionsSolved: []string{},
		LastSeen:        time.Time{},
	}

	// Perform custom validations
	if !validation.ValidateUsername(user.StandardUser.Username) {
		errs.NewBadRequestError("Invalid username").ToJSON(w)
		logger.Logger.Errorw("Invalid username", "method", r.Method, "username", user.StandardUser.Username, "time", time.Now())
		return
	}
	if !validation.ValidatePassword(user.StandardUser.Password) {
		errs.NewBadRequestError("Invalid password").ToJSON(w)
		logger.Logger.Errorw("Invalid password", "method", r.Method, "time", time.Now())
		return
	}
	if !validation.ValidateName(user.StandardUser.Name) {
		errs.NewBadRequestError("Invalid name").ToJSON(w)
		logger.Logger.Errorw("Invalid name", "method", r.Method, "time", time.Now())
		return
	}

	isEmailValid, isReputable := validation.ValidateEmail(user.StandardUser.Email)
	if !isEmailValid {
		errs.NewBadRequestError("Invalid email format").ToJSON(w)
		logger.Logger.Errorw("Invalid email format", "method", r.Method, "email", user.StandardUser.Email, "time", time.Now())
		return
	} else if !isReputable {
		errs.NewBadRequestError("Unsupported email domain (use gmail, hotmail, outlook, watchguard or icloud)").ToJSON(w)
		logger.Logger.Errorw("Unsupported email domain", "method", r.Method, "email", user.StandardUser.Email, "time", time.Now())
		return
	}

	isOrgValid, orgErr := validation.ValidateOrganizationName(user.StandardUser.Organisation)
	if !isOrgValid {
		errs.NewBadRequestError(orgErr.Error()).ToJSON(w)
		logger.Logger.Errorw("Invalid organization name", "method", r.Method, "organization", user.StandardUser.Organisation, "time", time.Now())
		return
	}

	isCountryValid, countryErr := validation.ValidateCountryName(user.StandardUser.Country)
	if !isCountryValid {
		errs.NewBadRequestError(countryErr.Error()).ToJSON(w)
		logger.Logger.Errorw("Invalid country name", "method", r.Method, "country", user.StandardUser.Country, "time", time.Now())
		return
	}

	err = a.authService.Signup(r.Context(), &user)
	if err != nil {
		if errors.Is(err, errs.ErrUserNameAlreadyExists) {
			errs.NewConflictError("User already exists").ToJSON(w)
		} else if errors.Is(err, errs.ErrEmailAlreadyExists) {
			errs.NewConflictError("Email already registered").ToJSON(w)
		} else if errors.Is(err, errs.ErrLeetcodeIDAlreadyExists) {
			errs.NewConflictError("LeetcodeID already registered").ToJSON(w)
		} else if errors.Is(err, errs.ErrLeetcodeUsernameInvalid) {
			errs.NewBadRequestError("Invalid leetcode id").ToJSON(w)
		} else {
			errs.NewInternalServerError("Signup failed").ToJSON(w)
		}
		logger.Logger.Errorw("Signup failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"message": "User successfully registered",
		"code":    http.StatusOK,
		"user_info": map[string]any{
			"username": user.StandardUser.Username,
		},
	}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Signup Successful", "method", r.Method, "username", user.StandardUser.Username, "time", time.Now())
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		errs.NewInvalidRequestBodyError("Invalid request body").ToJSON(w)
		logger.Logger.Errorw("Error decoding request body", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Validate the request body
	if !validation.ValidateUsername(requestBody.Username) {
		errs.NewBadRequestError("Invalid username").ToJSON(w)
		logger.Logger.Errorw("Invalid username", "method", r.Method, "username", requestBody.Username, "time", time.Now())
		return
	}
	if !validation.ValidatePassword(requestBody.Password) {
		errs.NewBadRequestError("Invalid username or password").ToJSON(w)
		logger.Logger.Errorw("Invalid password", "method", r.Method, "username", requestBody.Username, "time", time.Now())
		return
	}

	user, err := a.authService.Login(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPassword) {
			errs.NewInvalidRequestBodyError("Invalid username or password").ToJSON(w)
			w.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(err, errs.ErrUserNotFound) {
			errs.NewNotFoundError("User not found").ToJSON(w)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			errs.NewInternalServerError("Login failed").ToJSON(w)
		}
		logger.Logger.Errorw("Authentication failed", "method", r.Method, "error", err, "username", requestBody.Username, "time", time.Now())
		return
	}

	token, err := utils.CreateJwtToken(user.StandardUser.Username, user.StandardUser.ID, user.StandardUser.Role)
	if err != nil {
		errs.NewInternalServerError("Failed to generate token").ToJSON(w)
		logger.Logger.Errorw("Failed to generate token", "method", r.Method, "username", requestBody.Username, "error", err, "time", time.Now())
		return
	}

	globals.ActiveUserID = user.StandardUser.ID

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{"code": http.StatusOK, "message": "Login successful", "token": token, "role": user.StandardUser.Role}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Login Successful", "method", r.Method, "username", requestBody.Username, "time", time.Now())
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := a.authService.Logout(r.Context())
	if err != nil {
		errs.NewInternalServerError(err.Error()).ToJSON(w)
		logger.Logger.Errorw("Logout failed", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	globals.ActiveUserID = ""

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{"code": http.StatusOK, "message": "Logout successful"}
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Logout Successful", "method", r.Method, "time", time.Now())
}
