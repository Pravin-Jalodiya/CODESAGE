package interfaces

import (
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/models"
)

type UserService interface {
	GetAllUsers() (*[]models.StandardUser, error)
	ViewDashboard() error
	UpdateUserProgress() error
	GetUserProgress(userID string) (*[]string, error)
	CountActiveUserInLast24Hours() (int, error)
	GetUserByUsername(username string) (*models.StandardUser, error)
	GetUserByID(userID string) (*models.StandardUser, error)
	GetUserRole(userID string) (roles.Role, error)
	GetUserID(username string) (string, error)
	BanUser(username string) (bool, error)
	UnbanUser(username string) (bool, error)
	IsUserBanned(userID string) (bool, error)
	GetUserLeetcodeStats(userID string) (*models.LeetcodeStats, error)
	GetUserCodesageStats(userID string) (*models.CodesageStats, error)
	GetPlatformStats() (*models.PlatformStats, error)
}
