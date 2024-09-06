package interfaces

import (
	"cli-project/internal/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(*models.StandardUser) error
	UpdateUserProgress(userID uuid.UUID, questionID []string) error
	FetchAllUsers() (*[]models.StandardUser, error)
	FetchUserByID(string) (*models.StandardUser, error)
	FetchUserByUsername(string) (*models.StandardUser, error)
	UpdateUserDetails(*models.StandardUser) error
	BanUser(string) error
	UnbanUser(string) error
	CountActiveUsersInLast24Hours() (int64, error)
	IsUsernameUnique(string) (bool, error)
	IsEmailUnique(string) (bool, error)
	IsLeetcodeIDUnique(string) (bool, error)
}
