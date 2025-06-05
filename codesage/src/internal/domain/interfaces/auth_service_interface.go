package interfaces

import (
	"codesage/internal/domain/models"
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, user *models.StandardUser) error
	Login(ctx context.Context, username, password string) (*models.StandardUser, error)
	Logout(ctx context.Context) error
	IsEmailUnique(ctx context.Context, email string) (bool, error)
	IsUsernameUnique(ctx context.Context, username string) (bool, error)
	IsLeetcodeIDUnique(ctx context.Context, LeetcodeID string) (bool, error)
	//ValidateLeetcodeUsername(username string) (bool, error)
	//UpdateUserPassword(ctx context.Context, email string, password string) error
	//GenerateAndSendOtp(email string) error
}
