package interfaces

import "cli-project/internal/domain/models"

type LeetcodeAPI interface {
	GetStats(LeetcodeID string) (*models.LeetcodeStats, error)
	ValidateUsername(username string) (bool, error)
}
