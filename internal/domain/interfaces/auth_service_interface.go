package interfaces

import (
	"cli-project/internal/domain/models"
	"context"
)

type AuthService interface {
	Signup(user *models.StandardUser) error
	Login(ctx context.Context, username, password string) (*models.StandardUser, error)
	Logout() error
	IsEmailUnique(email string) (bool, error)
	IsUsernameUnique(username string) (bool, error)
	IsLeetcodeIDUnique(LeetcodeID string) (bool, error)
	ValidateLeetcodeUsername(username string) (bool, error)
}
