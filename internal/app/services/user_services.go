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
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("username or password incorrect")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo        interfaces.UserRepository
	questionService *QuestionService
	//userWG   *sync.WaitGroup
}

func NewUserService(userRepo interfaces.UserRepository, questionService *QuestionService) *UserService {
	return &UserService{
		userRepo:        userRepo,
		questionService: questionService,
		//userWG:   &sync.WaitGroup{},
	}
}

// SignUp creates a new user account
func (s *UserService) SignUp(user *models.StandardUser) error {

	// Change username to lowercase for consistency
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)

	// Change email to lower for consistency
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)

	// Change org and country name to lowercase
	user.StandardUser.Organisation = strings.ToLower(user.StandardUser.Organisation)
	user.StandardUser.Country = strings.ToLower(user.StandardUser.Country)

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

	// Set default blocked status
	user.StandardUser.IsBanned = false

	// set question solved
	user.QuestionsSolved = []string{}

	// set last seen
	user.LastSeen = time.Now().UTC()

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
	user, err := s.userRepo.FetchUserByUsername(username)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound // Return error if user doesn't exist
		}
		return fmt.Errorf("%v", err)
	}

	// Verify the password
	if !pwd.VerifyPassword(password, user.StandardUser.Password) {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *UserService) Logout() error {
	// Get active user
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
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
	globals.ActiveUserID = ""

	return nil
}

func (s *UserService) GetAllUsers() (*[]models.StandardUser, error) {
	return s.userRepo.FetchAllUsers()
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard() error {
	// Placeholder implementation
	return nil
}

// UpdateUserProgress updates the user's progress by adding a solved question ID.
func (s *UserService) UpdateUserProgress(solvedQuestionID string) (bool, error) {
	// Fetch the current user from the repository
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
	if err != nil {
		return false, fmt.Errorf("could not fetch user: %v", err)
	}

	// Check if the question ID is already in the user's progress
	for _, id := range user.QuestionsSolved {
		if id == solvedQuestionID {
			return false, nil // No need to update if the question ID is already in the list
		}
	}

	// Check if the question ID exists in the questions repository
	exists, err := s.questionService.QuestionExists(solvedQuestionID)
	if err != nil {
		return false, fmt.Errorf("could not check if question exists: %v", err)
	}
	if !exists {
		return false, fmt.Errorf("question with ID %s does not exist", solvedQuestionID)
	}

	// Update the user's progress
	return true, s.userRepo.UpdateUserProgress(solvedQuestionID)
}

func (s *UserService) CountActiveUserInLast24Hours() (int64, error) {
	count, err := s.userRepo.CountActiveUsersInLast24Hours()
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.StandardUser, error) {

	if username == "" {
		return nil, errors.New("username is empty")
	}

	// Change userID to lowercase for consistency
	username = data_cleaning.CleanString(username)

	return s.userRepo.FetchUserByUsername(username)
}

func (s *UserService) GetUserByID(userID string) (*models.StandardUser, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	// Clean userID for consistency
	userID = data_cleaning.CleanString(userID)

	// Fetch user by ID from the repository
	return s.userRepo.FetchUserByID(userID)
}

func (s *UserService) GetUserRole(userID string) (string, error) {

	if userID == "" {
		return "", errors.New("userID is empty")
	}

	// Change userID to lowercase for consistency
	userID = data_cleaning.CleanString(userID)

	user, err := s.userRepo.FetchUserByID(userID)
	if err != nil {
		return "", err
	}

	return user.StandardUser.Role, nil
}

func (s *UserService) GetUserID(username string) (string, error) {
	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.StandardUser.ID, nil
}

func (s *UserService) BanUser(username string) (bool, error) {

	userID, err := s.GetUserID(username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(userID)
	if err != nil {
		return false, err
	}

	if alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.BanUser(userID)
}

func (s *UserService) UnbanUser(username string) (bool, error) {

	userID, err := s.GetUserID(username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(userID)
	if err != nil {
		return false, err
	}

	if !alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.UnbanUser(userID)
}

func (s *UserService) IsUserBanned(userID string) (bool, error) {

	user, err := s.userRepo.FetchUserByID(userID)
	if err != nil {
		return false, err
	}

	return user.StandardUser.IsBanned, nil
}

//func (s *UserService) WaitForCompletion() {
//	s.userWG.Wait()
//}
