package services

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils"
	"cli-project/pkg/utils/password"
	"fmt"
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
func (s *UserService) SignUp(user models.StandardUser) error {
	// Generate a new UUID for the user
	userID := utils.GenerateUUID()
	user.StandardUser.ID = userID

	// Hash the user password
	hashedPassword, err := password.HashPassword(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err)
	}
	user.StandardUser.Password = hashedPassword

	// Set default role
	user.StandardUser.Role = "user"

	// Register the user
	err = s.userRepo.RegisterUser(user)
	if err != nil {
		return fmt.Errorf("could not register user: %v", err)
	}

	return nil
}

// Login authenticates a user
func (s *UserService) Login(username string, password string) error {
	// Placeholder implementation
	return nil
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard() error {
	// Placeholder implementation
	return nil
}

// UpdateProgress updates the user's progress in some context
func (s *UserService) UpdateProgress(userId string, progressData interface{}) error {
	// Placeholder implementation
	return nil
}

func (s *UserService) CountActiveUserInLast24Hours() (int64, error) {
	count, err := s.userRepo.CountActiveUsersInLast24Hours()
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *UserService) Logout() {
	// update lastseen of user
	// logout the user
}

func (us *UserService) IsEmailUnique(email string) (bool, error) {
	return us.userRepo.FindUserByEmail(email)
}

func (us *UserService) IsUsernameUnique(username string) (bool, error) {
	return us.userRepo.FindUserByUsername(username)
}

func (us *UserService) IsLeetCodeIDUnique(leetcodeID string) (bool, error) {
	return us.userRepo.FindUserByLeetcodeID(leetcodeID)
}

//func (s *UserService) WaitForCompletion() {
//	s.userWG.Wait()
//}
