package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
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

	hashedPassword, err := HashString(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("could not hash password")
	}
	user.StandardUser.Password = hashedPassword
	user.StandardUser.Role = "user"
	user.StandardUser.IsBanned = false
	user.QuestionsSolved = []string{}
	user.LastSeen = time.Now().UTC()

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("could not register user")
	}
	return nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) error {
	username = utils.CleanString(username)

	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("error fetching user: %v", err)
	}

	if !VerifyString(password, user.StandardUser.Password) {
		return ErrInvalidCredentials
	}
	return nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	user, err := s.userRepo.FetchUserByID(ctx, globals.ActiveUserID)
	if err != nil {
		return errors.New("user not found")
	}

	user.LastSeen = time.Now().UTC()

	err = s.userRepo.UpdateUserDetails(ctx, user)
	if err != nil {
		return errors.New("could not update user details")
	}

	globals.ActiveUserID = ""

	return nil
}

func (s *AuthService) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	return s.userRepo.IsEmailUnique(ctx, email)
}

func (s *AuthService) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	return s.userRepo.IsUsernameUnique(ctx, username)
}

func (s *AuthService) IsLeetcodeIDUnique(ctx context.Context, LeetcodeID string) (bool, error) {
	return s.userRepo.IsLeetcodeIDUnique(ctx, LeetcodeID)
}

// ValidateLeetcodeUsername checks if the provided Leetcode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	return s.LeetcodeAPI.ValidateUsername(username)
}
