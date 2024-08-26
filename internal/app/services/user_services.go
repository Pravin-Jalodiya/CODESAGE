package services

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	"cli-project/pkg/utils/data_cleaning"
	pwd "cli-project/pkg/utils/password"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("username or password incorrect")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo interfaces.UserRepository
	//userWG   *sync.WaitGroup
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		//userWG:   &sync.WaitGroup{},
	}
}

// SignUp creates a new user account
func (s *UserService) SignUp(user *models.StandardUser) error {

	// Change username to lowercase for consistency
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Name)

	// Change email to lower for consistency
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)

	// Generate a new UUID for the user
	userID := utils.GenerateUUID()
	user.StandardUser.ID = userID

	// Hash the user password
	hashedPassword, err := pwd.HashPassword(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("could not hash password")
	}
	user.StandardUser.Password = hashedPassword

	// Set default role
	user.StandardUser.Role = "user"

	// Register the user
	err = s.userRepo.RegisterUser(user)
	if err != nil {
		return fmt.Errorf("could not register user")
	}

	return nil
}

// Login authenticates a user
func (s *UserService) Login(username, password string) error {

	// Change username to lowercase for consistency
	username = data_cleaning.CleanString(username)

	// Retrieve the user by username
	user, err := s.userRepo.FetchUser(username)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound // Return error if user doesn't exist
		}
		return fmt.Errorf("could not retrieve user: %v", err)
	}

	// Verify the password
	if !pwd.VerifyPassword(password, user.StandardUser.Password) {
		return ErrInvalidCredentials
	}

	return nil
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard() error {
	// Placeholder implementation
	return nil
}

// UpdateUserProgress updates the user's progress in some context
func (s *UserService) UpdateUserProgress(username string, solvedQuestionID string) error {

	// Change username to lowercase for consistency
	username = data_cleaning.CleanString(username)

	// Check if question ID is valid or not
	_, err := strconv.Atoi(solvedQuestionID)
	if err != nil {
		return fmt.Errorf("invalid question id")
	}

	return s.userRepo.UpdateUserProgress(username, solvedQuestionID)
}

func (s *UserService) CountActiveUserInLast24Hours() (int64, error) {
	count, err := s.userRepo.CountActiveUsersInLast24Hours()
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *UserService) Logout() error {
	// Get active user
	user, err := s.userRepo.FetchUser(globals.ActiveUser)
	if err != nil {
		return errors.New("user not found")
	}

	// update last seen of user
	user.LastSeen = time.Now().UTC()

	// update data in db
	err = s.userRepo.UpdateUserDetails(user)
	if err != nil {
		return errors.New("could not update user details")
	}

	// clear active user
	globals.ActiveUser = ""

	return nil
}

func (s *UserService) GetUserByUsername(username string) (*models.StandardUser, error) {

	if username == "" {
		return nil, errors.New("username is empty")
	}

	// Change username to lowercase for consistency
	username = data_cleaning.CleanString(username)

	return s.userRepo.FetchUser(username)
}

func (s *UserService) GetUserRole(username string) (string, error) {

	if username == "" {
		return "", errors.New("username is empty")
	}

	// Change username to lowercase for consistency
	username = data_cleaning.CleanString(username)

	user, err := s.userRepo.FetchUser(username)
	if err != nil {
		return "", err
	}

	return user.StandardUser.Role, nil
}

//func (s *UserService) WaitForCompletion() {
//	s.userWG.Wait()
//}
