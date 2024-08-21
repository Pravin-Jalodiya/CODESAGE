package interfaces

import "cli-project/internal/domain/models"

type UserRepository interface {
	RegisterUser(user models.StandardUser) error
	UpdateUserProgress(username string, questionID int) error
	FetchAllUsers() ([]models.StandardUser, error)
	FetchUser(username string) (models.StandardUser, error)
	CountActiveUsersInLast24Hours() (int64, error)
	FindUserByUsername(username string) (bool, error)
	FindUserByEmail(email string) (bool, error)
	FindUserByLeetcodeID(leetcodeID string) (bool, error)
}
