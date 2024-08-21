package interfaces

import "cli-project/internal/domain/models"

type UserRepository interface {
	RegisterUser(user models.StandardUser) error
	UpdateUserProgress(username string, questionID int) error
	FetchAllUsers() ([]models.StandardUser, error)
	FetchUser(username string) (models.StandardUser, error)
	CountActiveUserInLast24Hours() (int, error)
}
