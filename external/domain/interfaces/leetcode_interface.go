package interfaces

import "cli-project/internal/domain/models"

type LeetcodeAPI interface {
	GetStats(LeetcodeID string) (*models.LeetcodeStats, error)
	ValidateUsername(username string) (bool, error)
	FetchData(query string, variables map[string]interface{}) (map[string]interface{}, error)
	FetchUserStats(username string) (*models.LeetcodeStats, error)
	FetchRecentSubmissions(username string, limit int) ([]map[string]string, error)
}
