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
// func (s *AuthService) Signup(user *models.StandardUser) error {
//
//		// Change username to lowercase for consistency
//		user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)
//
//		// Change email to lower for consistency
//		user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)
//
//		// Change org and country name to proper format
//		user.StandardUser.Organisation = utils.CapitalizeWords(user.StandardUser.Organisation)
//		user.StandardUser.Country = utils.CapitalizeWords(user.StandardUser.Country)
//
//		// Generate a new UUID for the user
//		userID := utils.GenerateUUID()
//		user.StandardUser.ID = userID
//
//		// Hash the user password
//		hashedPassword, err := HashString(user.StandardUser.Password)
//		if err != nil {
//			return fmt.Errorf("could not hash password")
//		}
//		user.StandardUser.Password = hashedPassword
//
//		// Set default role
//		user.StandardUser.Role = "user"
//
//		// Set default blocked status
//		user.StandardUser.IsBanned = false
//
//		// set question solved
//		user.QuestionsSolved = []string{}
//
//		// set last seen
//		user.LastSeen = time.Now().UTC()
//
//		// Register the user
//		err = s.userRepo.CreateUser(user)
//		if err != nil {
//			return fmt.Errorf("could not register user")
//		}
//		return nil
//	}
//
// Signup creates a new user account
func (s *AuthService) Signup(user *models.StandardUser) error {
	// Normalize and sanitize user data
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)
	user.StandardUser.Organisation = utils.CapitalizeWords(user.StandardUser.Organisation)
	user.StandardUser.Country = utils.CapitalizeWords(user.StandardUser.Country)

	// Generate a new UUID for the user
	userID := utils.GenerateUUID()
	user.StandardUser.ID = userID

	// Hash the user password
	hashedPassword, err := HashString(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("could not hash password")
	}
	user.StandardUser.Password = hashedPassword

	// Set default role and status
	user.StandardUser.Role = "user"
	user.StandardUser.IsBanned = false

	// Initialize user-specific fields
	user.QuestionsSolved = []string{}
	user.LastSeen = time.Now().UTC()

	// Register the user
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("could not register user")
	}

	return nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.StandardUser, error) {
	// Change username to lowercase for consistency
	username = utils.CleanString(username)

	// Retrieve the user by username
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		// Check if the error is because the user was not found (PostgreSQL)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found:", username)
			return nil, ErrUserNotFound // Return error if user doesn't exist
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	// Verify the password
	if !VerifyString(password, user.StandardUser.Password) {
		return nil, ErrInvalidCredentials
	}

	globals.ActiveUserID = user.StandardUser.ID
	return user, nil
}

func (s *AuthService) Logout() error {
	// Get active user
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update last seen of the user
	user.LastSeen = time.Now().UTC()

	// Update user details in the database
	err = s.userRepo.UpdateUserDetails(user)
	if err != nil {
		return errors.New("could not update user details")
	}

	// Clear active user ID
	globals.ActiveUserID = ""

	return nil
}

func (s *AuthService) IsEmailUnique(email string) (bool, error) {
	return s.userRepo.IsEmailUnique(email)
}

func (s *AuthService) IsUsernameUnique(username string) (bool, error) {
	return s.userRepo.IsUsernameUnique(username)
}

func (s *AuthService) IsLeetcodeIDUnique(LeetcodeID string) (bool, error) {
	return s.userRepo.IsLeetcodeIDUnique(LeetcodeID)
}

// ValidateLeetcodeUsername checks if the provided Leetcode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	return s.LeetcodeAPI.ValidateUsername(username)
}
