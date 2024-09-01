package interfaces

import (
	"cli-project/internal/domain/models"
)

type UserService interface {
	Signup(user *models.StandardUser) error
	Login(username, password string) error
	Logout() error
	GetAllUsers() (*[]models.StandardUser, error)
	ViewDashboard() error
	UpdateUserProgress(solvedQuestionID string) (bool, error)
	CountActiveUserInLast24Hours() (int64, error)
	GetUserByUsername(username string) (*models.StandardUser, error)
	GetUserByID(userID string) (*models.StandardUser, error)
	GetUserRole(userID string) (string, error)
	GetUserID(username string) (string, error)
	BanUser(username string) (bool, error)
	UnbanUser(username string) (bool, error)
	IsUserBanned(userID string) (bool, error)
	GetLeetcodeStats(userID string) (*models.LeetcodeStats, error)
}
