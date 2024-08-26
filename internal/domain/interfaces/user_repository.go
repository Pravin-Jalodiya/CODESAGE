package interfaces

import "cli-project/internal/domain/models"

type UserRepository interface {
	RegisterUser(*models.StandardUser) error
	UpdateUserProgress(questionID string) error
	FetchAllUsers() (*[]models.StandardUser, error)
	FetchUserByID(string) (*models.StandardUser, error)
	FetchUserByUsername(string) (*models.StandardUser, error)
	UpdateUserDetails(*models.StandardUser) error
	CountActiveUsersInLast24Hours() (int64, error)
	IsUsernameUnique(string) (bool, error)
	IsEmailUnique(string) (bool, error)
	IsLeetcodeIDUnique(string) (bool, error)
}
