package interfaces

import (
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]dto.StandardUser, error)
	ViewDashboard(ctx context.Context) error
	UpdateUserProgress(ctx context.Context, userID uuid.UUID) error
	GetUserProgress(ctx context.Context, userID string) (*[]string, error)
	CountActiveUserInLast24Hours(ctx context.Context) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*models.StandardUser, error)
	GetUserByID(ctx context.Context, userID string) (*models.StandardUser, error)
	GetUserRole(ctx context.Context, userID string) (roles.Role, error)
	GetUserID(ctx context.Context, username string) (string, error)
	UpdateUserBanState(ctx context.Context, username string) (string, error)
	IsUserBanned(ctx context.Context, userID string) (bool, error)
	GetUserLeetcodeStats(userID string) (*models.LeetcodeStats, error)
	GetUserCodesageStats(ctx context.Context, userID string) (*models.CodesageStats, error)
	GetPlatformStats(ctx context.Context) (*models.PlatformStats, error)
	DeleteUser(ctx context.Context, username string) error
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
}
