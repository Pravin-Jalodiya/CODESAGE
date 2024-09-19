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
func (s *AuthService) Signup(user *models.StandardUser) error {

	// Change username to lowercase for consistency
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)

	// Change email to lower for consistency
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)

	// Change org and country name to proper format
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

	// Set default role
	user.StandardUser.Role = "user"

	// Set default blocked status
	user.StandardUser.IsBanned = false

	// set question solved
	user.QuestionsSolved = []string{}

	// set last seen
	user.LastSeen = time.Now().UTC()

	// Register the user
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("could not register user")
	}
	return nil
}

// Login authenticates a user
// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.StandardUser, error) {
	// Change username to lowercase for consistency
	fmt.Println("Login attempt started for username:", username)
	username = utils.CleanString(username)
	fmt.Println("Username after cleaning and converting to lowercase:", username)

	// Retrieve the user by username
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		// Check if the error is because the user was not found (PostgreSQL)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found:", username)
			return nil, ErrUserNotFound // Return error if user doesn't exist
		}
		fmt.Println("Error fetching user from database:", err)
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	fmt.Println("User fetched successfully for username:", username)

	// Verify the password
	if !VerifyString(password, user.StandardUser.Password) {
		fmt.Println("Invalid password attempt for username:", username)
		return nil, ErrInvalidCredentials
	}
	fmt.Println("User authenticated successfully for username:", username)

	return user, nil
}

func (s *AuthService) Logout() error {
	// Get active user
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
	if err != nil {
		return errors.New("user not found")
	}

	//update last seen of user
	user.LastSeen = time.Now().UTC()

	// update data in db
	err = s.userRepo.UpdateUserDetails(user)
	if err != nil {
		return errors.New("could not update user details")
	}

	// clear Active user id
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
