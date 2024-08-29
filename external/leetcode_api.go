package external

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

type leetcodeAPI struct{}

func NewLeetcodeAPI() interfaces.LeetcodeAPI {
	return &leetcodeAPI{}
}

// GetLeetcodeStats makes the API call to fetch user leetcode stats
func (api *leetcodeAPI) GetStats(leetcodeID string) (*models.LeetcodeStats, error) {

	recentLimit := 10

	// Updated GraphQL query
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

	recentSubmissionsQuery := `
	query recentAcSubmissions($username: String!, $limit: Int!) {
		recentAcSubmissionList(username: $username, limit: $limit) {
			title
		}
	}`

	// Function to perform GraphQL request
	fetchData := func(query string, variables map[string]interface{}) (map[string]interface{}, error) {
		requestBody := map[string]interface{}{
			"query":     query,
			"variables": variables,
		}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request body: %v", err)
		}

		resp, err := http.Post(config.LEETCODE_API, "application/json", bytes.NewBuffer(jsonBody))
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
	statsData, err := fetchData(userStatsQuery, map[string]interface{}{"username": leetcodeID})
	if err != nil {
		return nil, err
	}

	// Parse stats data
	stats := &models.LeetcodeStats{
		RecentACSubmissions: []string{},
	}

	if allQuestionsCount, ok := statsData["allQuestionsCount"].([]interface{}); ok {
		for _, item := range allQuestionsCount {
			countInfo := item.(map[string]interface{})
			difficulty := countInfo["difficulty"].(string)
			count := int(countInfo["count"].(float64))
			switch difficulty {
			case "All":
				// These are totals for all difficulties combined, not used for individual counts
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

	// Fetch recent accepted submissions
	submissionsData, err := fetchData(recentSubmissionsQuery, map[string]interface{}{"username": leetcodeID, "limit": recentLimit})
	if err != nil {
		return nil, err
	}

	if recentSubmissions, ok := submissionsData["recentAcSubmissionList"].([]interface{}); ok {
		for _, item := range recentSubmissions {
			submission := item.(map[string]interface{})
			if title, ok := submission["title"].(string); ok {
				stats.RecentACSubmissions = append(stats.RecentACSubmissions, title)
			}
		}
	}

	return stats, nil
}

func (api *leetcodeAPI) ValidateUsername(username string) (bool, error) {

	// GraphQL query to check if a LeetCode user exists
	const userQuery = `
	query getUserProfile($username: String!) {
  	matchedUser(username: $username) {
    username
  }
}
`

	// Construct the GraphQL query with variables
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

	// Make the HTTP request
	resp, err := http.Post(config.LEETCODE_API, "application/json", bytes.NewBuffer(jsonBody))
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

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("could not decode response: %v", err)
	}

	// Check if the user exists
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("invalid response format")
	}
	matchedUser, ok := data["matchedUser"].(map[string]interface{})
	if !ok || matchedUser == nil {
		return false, nil // User does not exist
	}

	// Check if the username matches
	return matchedUser["username"] == username, nil
}
