package interfaces

import (
	"cli-project/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.StandardUser) error
	UpdateUserProgress(ctx context.Context, userID uuid.UUID, questionID []string) error
	FetchAllUsers(ctx context.Context) ([]models.StandardUser, error)
	FetchUserByID(ctx context.Context, id string) (*models.StandardUser, error)
	FetchUserByUsername(ctx context.Context, username string) (*models.StandardUser, error)
	FetchUserProgress(ctx context.Context, userID string) (*[]string, error)
	UpdateUserDetails(ctx context.Context, user *models.StandardUser) error
	BanUser(ctx context.Context, userID string) error
	UnbanUser(ctx context.Context, userID string) error
	CountActiveUsersInLast24Hours(ctx context.Context) (int, error)
	IsUsernameUnique(ctx context.Context, username string) (bool, error)
	IsEmailUnique(ctx context.Context, email string) (bool, error)
	IsLeetcodeIDUnique(ctx context.Context, leetcodeID string) (bool, error)
}
