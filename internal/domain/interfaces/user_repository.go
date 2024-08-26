package interfaces

import "cli-project/internal/domain/models"

type UserRepository interface {
	RegisterUser(user *models.StandardUser) error
	UpdateUserProgress(username string, questionID string) error
	FetchAllUsers() (*[]models.StandardUser, error)
	FetchUser(username string) (*models.StandardUser, error)
	UpdateUserDetails(user *models.StandardUser) error
	CountActiveUsersInLast24Hours() (int64, error)
	IsUsernameUnique(username string) (bool, error)
	IsEmailUnique(email string) (bool, error)
	IsLeetcodeIDUnique(leetcodeID string) (bool, error)
}
