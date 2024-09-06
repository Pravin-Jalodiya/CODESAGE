package api

import (
	"bytes"
	"cli-project/external/domain/interfaces"
	"cli-project/internal/config"
	"cli-project/internal/domain/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LeetcodeAPI struct{}

func NewLeetcodeAPI() interfaces.LeetcodeAPI {
	return &LeetcodeAPI{}
}

// Function to perform GraphQL request
func fetchData(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %v", err)
	}

	resp, err := http.Post(config.Leetcode_API, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	return data, nil
}

// Fetch user stats
func (api *LeetcodeAPI) fetchUserStats(username string) (*models.LeetcodeStats, error) {
	userStatsQuery := `
	query userProblemsSolved($username: String!) {
		allQuestionsCount {
			difficulty
			count
		}
		matchedUser(username: $username) {
			submitStatsGlobal {
				acSubmissionNum {
					difficulty
					count
				}
			}
		}
	}`

	statsData, err := fetchData(userStatsQuery, map[string]interface{}{"username": username})
	if err != nil {
		return nil, err
	}

	stats := &models.LeetcodeStats{
		RecentACSubmissionTitles:     []string{},
		RecentACSubmissionTitleSlugs: []string{},
	}

	if allQuestionsCount, ok := statsData["allQuestionsCount"].([]interface{}); ok {
		for _, item := range allQuestionsCount {
			countInfo := item.(map[string]interface{})
			difficulty := countInfo["difficulty"].(string)
			count := int(countInfo["count"].(float64))
			switch difficulty {
			case "All":
				stats.TotalQuestionsCount = count
			case "Easy":
				stats.TotalEasyCount = count
			case "Medium":
				stats.TotalMediumCount = count
			case "Hard":
				stats.TotalHardCount = count
			}
		}
	}

	if matchedUser, ok := statsData["matchedUser"].(map[string]interface{}); ok {
		if submitStatsGlobal, ok := matchedUser["submitStatsGlobal"].(map[string]interface{}); ok {
			if acSubmissionNum, ok := submitStatsGlobal["acSubmissionNum"].([]interface{}); ok {
				for _, item := range acSubmissionNum {
					difficultyCount := item.(map[string]interface{})
					difficulty := difficultyCount["difficulty"].(string)
					count := int(difficultyCount["count"].(float64))
					switch difficulty {
					case "All":
						stats.TotalQuestionsDoneCount = count
					case "Easy":
						stats.EasyDoneCount = count
					case "Medium":
						stats.MediumDoneCount = count
					case "Hard":
						stats.HardDoneCount = count
					}
				}
			}
		}
	}

	return stats, nil
}

// Fetch recent accepted submissions
func (api *LeetcodeAPI) fetchRecentSubmissions(username string, limit int) ([]map[string]string, error) {
	recentSubmissionsQuery := `
	query recentAcSubmissions($username: String!, $limit: Int!) {
		recentAcSubmissionList(username: $username, limit: $limit) {
			title
			titleSlug
		}
	}`

	submissionsData, err := fetchData(recentSubmissionsQuery, map[string]interface{}{"username": username, "limit": limit})
	if err != nil {
		return nil, err
	}

	var submissions []map[string]string
	if recentSubmissions, ok := submissionsData["recentAcSubmissionList"].([]interface{}); ok {
		for _, item := range recentSubmissions {
			submission := item.(map[string]interface{})
			if title, ok := submission["title"].(string); ok {
				if titleSlug, ok := submission["titleSlug"].(string); ok {
					submissions = append(submissions, map[string]string{
						"title":     title,
						"titleSlug": titleSlug,
					})
				}
			}
		}
	}

	return submissions, nil
}

// GetStats combines fetchUserStats and fetchRecentSubmissions
func (api *LeetcodeAPI) GetStats(LeetcodeID string) (*models.LeetcodeStats, error) {
	recentLimit := config.RECENT_SUBMISSION_LIMIT

	stats, err := api.fetchUserStats(LeetcodeID)
	if err != nil {
		return nil, err
	}

	recentSubmissions, err := api.fetchRecentSubmissions(LeetcodeID, recentLimit)
	if err != nil {
		return nil, err
	}

	// Initialize slices for titles and titleSlugs
	var submissionTitles []string
	var submissionSlugs []string

	for _, submission := range recentSubmissions {
		if title, ok := submission["title"]; ok {
			submissionTitles = append(submissionTitles, title)
		}
		if slug, ok := submission["titleSlug"]; ok {
			submissionSlugs = append(submissionSlugs, slug)
		}
	}

	// Set both slices in the stats
	stats.RecentACSubmissionTitles = submissionTitles
	stats.RecentACSubmissionTitleSlugs = submissionSlugs

	return stats, nil
}

func (api *LeetcodeAPI) ValidateUsername(username string) (bool, error) {
	const userQuery = `
	query getUserProfile($username: String!) {
  	matchedUser(username: $username) {
    username
  }
}
`

	query := userQuery
	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"username": username,
		},
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("could not marshal request body: %v", err)
	}

	resp, err := http.Post(config.Leetcode_API, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("request failed: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("could not close response body")
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("could not decode response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid response format")
	}
	matchedUser, ok := data["matchedUser"].(map[string]interface{})
	if !ok || matchedUser == nil {
		return false, nil // User does not exist
	}

	return matchedUser["username"] == username, nil
}
