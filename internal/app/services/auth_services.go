package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	LeetcodeAPI interfaces2.LeetcodeAPI
}

func NewAuthService(userRepo interfaces.UserRepository, LeetcodeAPI interfaces2.LeetcodeAPI) interfaces.AuthService {
	return &AuthService{
		userRepo:    userRepo,
		LeetcodeAPI: LeetcodeAPI,
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(ctx context.Context, user *models.StandardUser) error {
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)
	user.StandardUser.Organisation = utils.CapitalizeWords(user.StandardUser.Organisation)
	user.StandardUser.Country = utils.CapitalizeWords(user.StandardUser.Country)
	user.StandardUser.ID = utils.GenerateUUID()

	// Check if the email is unique
	emailUnique, err := s.userRepo.IsEmailUnique(ctx, user.StandardUser.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !emailUnique {
		return fmt.Errorf("%w: %v", errs.ErrEmailAlreadyExists, user.StandardUser.Email)
	}

	// Check if the username is unique
	usernameUnique, err := s.userRepo.IsUsernameUnique(ctx, user.StandardUser.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !usernameUnique {
		return fmt.Errorf("%w: %v", errs.ErrUserNameAlreadyExists, user.StandardUser.Username)
	}

	// Check if the Leetcode ID is unique
	leetcodeIDUnique, err := s.userRepo.IsLeetcodeIDUnique(ctx, user.LeetcodeID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	if !leetcodeIDUnique {
		return fmt.Errorf("%w: %v", errs.ErrLeetcodeIDAlreadyExists, user.LeetcodeID)
	}

	hashedPassword, err := HashString(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("%w: could not hash password", errs.ErrInternalServerError)
	}
	user.StandardUser.Password = hashedPassword
	user.StandardUser.Role = "user"
	user.StandardUser.IsBanned = false
	user.QuestionsSolved = []string{}
	user.LastSeen = time.Now().UTC()

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: could not register user", errs.ErrDbError)
	}
	return nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.StandardUser, error) {
	username = utils.CleanString(username)

	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}

	if !VerifyString(password, user.StandardUser.Password) {
		return nil, errs.ErrInvalidPassword
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	user, err := s.userRepo.FetchUserByID(ctx, globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}

	user.LastSeen = time.Now().UTC()

	err = s.userRepo.UpdateUserDetails(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}

	globals.ActiveUserID = ""
	return nil
}

func (s *AuthService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	unique, err := s.userRepo.IsEmailUnique(ctx, email)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

func (s *AuthService) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	unique, err := s.userRepo.IsUsernameUnique(ctx, username)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

func (s *AuthService) IsLeetcodeIDUnique(ctx context.Context, LeetcodeID string) (bool, error) {
	unique, err := s.userRepo.IsLeetcodeIDUnique(ctx, LeetcodeID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return unique, nil
}

// ValidateLeetcodeUsername checks if the provided Leetcode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	valid, err := s.LeetcodeAPI.ValidateUsername(username)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrExternalAPI, err)
	}
	return valid, nil
}
