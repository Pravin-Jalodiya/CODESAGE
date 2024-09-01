package service_test

import (
	"cli-project/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cli-project/external/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestFetchData tests the fetchData function in the LeetcodeAPI implementation
func TestFetchData(t *testing.T) {
	// Set up a mock server to simulate API responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"allQuestionsCount": []interface{}{
					map[string]interface{}{"difficulty": "Easy", "count": 10},
					map[string]interface{}{"difficulty": "Medium", "count": 15},
					map[string]interface{}{"difficulty": "Hard", "count": 5},
				},
				"matchedUser": map[string]interface{}{
					"submitStatsGlobal": map[string]interface{}{
						"acSubmissionNum": []interface{}{
							map[string]interface{}{"difficulty": "Easy", "count": 8},
							map[string]interface{}{"difficulty": "Medium", "count": 10},
							map[string]interface{}{"difficulty": "Hard", "count": 3},
						},
					},
				},
			},
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBody)
	}))
	defer mockServer.Close()

	// Update the LeetcodeAPI to use the mock server URL
	originalURL := config.LEETCODE_API
	config.LEETCODE_API = mockServer.URL
	defer func() { config.LEETCODE_API = originalURL }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create an instance of LeetcodeAPI
	leetcodeAPI := api.NewLeetcodeAPI()

	// Test fetchData function
	query := `
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

	variables := map[string]interface{}{"username": "testuser"}

	data, err := leetcodeAPI.FetchData(query, variables)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Validate the structure of the response
	allQuestionsCount := data["allQuestionsCount"].([]interface{})
	assert.Len(t, allQuestionsCount, 3)

	submitStatsGlobal := data["matchedUser"].(map[string]interface{})["submitStatsGlobal"].(map[string]interface{})
	acSubmissionNum := submitStatsGlobal["acSubmissionNum"].([]interface{})
	assert.Len(t, acSubmissionNum, 3)
}

// TestGetUserProblemSolved tests the GetUserProblemSolved method in the LeetcodeAPI
func TestGetUserProblemSolved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	leetcodeAPI := api.NewLeetcodeAPI()

	// Mock the fetchData function
	mockFetchData := func(query string, variables map[string]interface{}) (map[string]interface{}, error) {
		response := map[string]interface{}{
			"allQuestionsCount": []interface{}{
				map[string]interface{}{"difficulty": "Easy", "count": 10},
				map[string]interface{}{"difficulty": "Medium", "count": 15},
				map[string]interface{}{"difficulty": "Hard", "count": 5},
			},
			"matchedUser": map[string]interface{}{
				"submitStatsGlobal": map[string]interface{}{
					"acSubmissionNum": []interface{}{
						map[string]interface{}{"difficulty": "Easy", "count": 8},
						map[string]interface{}{"difficulty": "Medium", "count": 10},
						map[string]interface{}{"difficulty": "Hard", "count": 3},
					},
				},
			},
		}
		return response, nil
	}

	// Replace the real fetchData with the mock
	leetcodeAPI.FetchData = mockFetchData

	// Test GetUserProblemSolved method
	stats, err := leetcodeAPI.GetUserProblemSolved("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Validate the stats
	assert.Equal(t, 10, stats.TotalEasyCount)
	assert.Equal(t, 15, stats.TotalMediumCount)
	assert.Equal(t, 5, stats.TotalHardCount)
	assert.Equal(t, 8, stats.EasyDoneCount)
	assert.Equal(t, 10, stats.MediumDoneCount)
	assert.Equal(t, 3, stats.HardDoneCount)
}

// TestGetRecentACSubmission tests the GetRecentACSubmission method in the LeetcodeAPI
func TestGetRecentACSubmission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	leetcodeAPI := api.NewLeetcodeAPI()

	// Mock the fetchData function
	mockFetchData := func(query string, variables map[string]interface{}) (map[string]interface{}, error) {
		response := map[string]interface{}{
			"recentAcSubmissionList": []interface{}{
				map[string]interface{}{"title": "Problem 1"},
				map[string]interface{}{"title": "Problem 2"},
				map[string]interface{}{"title": "Problem 3"},
			},
		}
		return response, nil
	}

	// Replace the real fetchData with the mock
	leetcodeAPI.FetchData = mockFetchData

	// Test GetRecentACSubmission method
	submissions, err := leetcodeAPI.GetRecentACSubmission("testuser", 3)
	assert.NoError(t, err)
	assert.Len(t, submissions, 3)
	assert.ElementsMatch(t, []string{"Problem 1", "Problem 2", "Problem 3"}, submissions)
}
