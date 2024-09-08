package api_test

import (
	"cli-project/external/api"
	"cli-project/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchData_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"allQuestionsCount": []interface{}{
					map[string]interface{}{"difficulty": "All", "count": float64(100)},
				},
				"matchedUser": map[string]interface{}{
					"submitStatsGlobal": map[string]interface{}{
						"acSubmissionNum": []interface{}{
							map[string]interface{}{"difficulty": "All", "count": float64(50)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	api := api.NewLeetcodeAPI()
	query := `query userProblemsSolved($username: String!) { allQuestionsCount { difficulty count } matchedUser(username: $username) { submitStatsGlobal { acSubmissionNum { difficulty count } } } }`
	variables := map[string]interface{}{"username": "testuser"}

	data, err := api.FetchData(query, variables)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 100, int(data["allQuestionsCount"].([]interface{})[0].(map[string]interface{})["count"].(float64)))
}

func TestFetchData_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	api := api.NewLeetcodeAPI()
	query := `query userProblemsSolved($username: String!) { allQuestionsCount { difficulty count } matchedUser(username: $username) { submitStatsGlobal { acSubmissionNum { difficulty count } } } }`
	variables := map[string]interface{}{"username": "testuser"}

	data, err := api.FetchData(query, variables)
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestFetchUserStats_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"allQuestionsCount": []interface{}{
					map[string]interface{}{"difficulty": "All", "count": float64(100)},
					map[string]interface{}{"difficulty": "Easy", "count": float64(60)},
					map[string]interface{}{"difficulty": "Medium", "count": float64(30)},
					map[string]interface{}{"difficulty": "Hard", "count": float64(10)},
				},
				"matchedUser": map[string]interface{}{
					"submitStatsGlobal": map[string]interface{}{
						"acSubmissionNum": []interface{}{
							map[string]interface{}{"difficulty": "All", "count": float64(50)},
							map[string]interface{}{"difficulty": "Easy", "count": float64(30)},
							map[string]interface{}{"difficulty": "Medium", "count": float64(15)},
							map[string]interface{}{"difficulty": "Hard", "count": float64(5)},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	stats, err := leetcodeAPI.FetchUserStats("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 100, stats.TotalQuestionsCount)
	assert.Equal(t, 60, stats.TotalEasyCount)
	assert.Equal(t, 30, stats.TotalMediumCount)
	assert.Equal(t, 10, stats.TotalHardCount)
	assert.Equal(t, 50, stats.TotalQuestionsDoneCount)
	assert.Equal(t, 30, stats.EasyDoneCount)
	assert.Equal(t, 15, stats.MediumDoneCount)
	assert.Equal(t, 5, stats.HardDoneCount)
}

func TestFetchUserStats_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	stats, err := leetcodeAPI.FetchUserStats("testuser")
	assert.Error(t, err)
	assert.Nil(t, stats)
}

func TestFetchRecentSubmissions_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"recentAcSubmissionList": []interface{}{
					map[string]interface{}{
						"title":     "Question 1",
						"titleSlug": "question-1",
					},
					map[string]interface{}{
						"title":     "Question 2",
						"titleSlug": "question-2",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	submissions, err := leetcodeAPI.FetchRecentSubmissions("testuser", 2)
	assert.NoError(t, err)
	assert.Len(t, submissions, 2)
	assert.Equal(t, map[string]string{"title": "Question 1", "titleSlug": "question-1"}, submissions[0])
	assert.Equal(t, map[string]string{"title": "Question 2", "titleSlug": "question-2"}, submissions[1])
}

func TestFetchRecentSubmissions_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	submissions, err := leetcodeAPI.FetchRecentSubmissions("testuser", 2)
	assert.Error(t, err)
	assert.Nil(t, submissions)
}

func TestLeetcodeAPI_ValidateUsername_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"matchedUser": map[string]interface{}{
					"username": "testuser",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	valid, err := leetcodeAPI.ValidateUsername("testuser")
	assert.NoError(t, err)
	assert.True(t, valid)
}

func TestLeetcodeAPI_ValidateUsername_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	valid, err := leetcodeAPI.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestGetStats_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"allQuestionsCount": []interface{}{
					map[string]interface{}{"difficulty": "All", "count": float64(100)},
					map[string]interface{}{"difficulty": "Easy", "count": float64(60)},
					map[string]interface{}{"difficulty": "Medium", "count": float64(30)},
					map[string]interface{}{"difficulty": "Hard", "count": float64(10)},
				},
				"matchedUser": map[string]interface{}{
					"submitStatsGlobal": map[string]interface{}{
						"acSubmissionNum": []interface{}{
							map[string]interface{}{"difficulty": "All", "count": float64(50)},
							map[string]interface{}{"difficulty": "Easy", "count": float64(30)},
							map[string]interface{}{"difficulty": "Medium", "count": float64(15)},
							map[string]interface{}{"difficulty": "Hard", "count": float64(5)},
						},
					},
				},
				"recentAcSubmissionList": []interface{}{
					map[string]interface{}{
						"title":     "Question 1",
						"titleSlug": "question-1",
					},
					map[string]interface{}{
						"title":     "Question 2",
						"titleSlug": "question-2",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL
	config.RECENT_SUBMISSION_LIMIT = 2

	leetcodeAPI := api.NewLeetcodeAPI()
	stats, err := leetcodeAPI.GetStats("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 100, stats.TotalQuestionsCount)
	assert.Equal(t, 60, stats.TotalEasyCount)
	assert.Equal(t, 30, stats.TotalMediumCount)
	assert.Equal(t, 10, stats.TotalHardCount)
	assert.Equal(t, 50, stats.TotalQuestionsDoneCount)
	assert.Equal(t, 30, stats.EasyDoneCount)
	assert.Equal(t, 15, stats.MediumDoneCount)
	assert.Equal(t, 5, stats.HardDoneCount)
	assert.ElementsMatch(t, []string{"Question 1", "Question 2"}, stats.RecentACSubmissionTitles)
	assert.ElementsMatch(t, []string{"question-1", "question-2"}, stats.RecentACSubmissionTitleSlugs)
}

func TestGetStats_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	stats, err := leetcodeAPI.GetStats("testuser")
	assert.Error(t, err)
	assert.Nil(t, stats)
}

func TestValidateUsername_HTTPRequestError(t *testing.T) {
	// Saving old Leetcode API URL to restore later
	oldLeetcodeAPI := config.Leetcode_API
	defer func() { config.Leetcode_API = oldLeetcodeAPI }()

	config.Leetcode_API = "\n" // Invalid URL to trigger error

	leetcodeAPI := api.NewLeetcodeAPI()
	_, err := leetcodeAPI.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestValidateUsername_UnexpectedStatusCode(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	oldLeetcodeAPI := config.Leetcode_API
	defer func() { config.Leetcode_API = oldLeetcodeAPI }()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()

	_, err := leetcodeAPI.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestValidateUsername_JSONDecodeError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer testServer.Close()

	oldLeetcodeAPI := config.Leetcode_API
	defer func() { config.Leetcode_API = oldLeetcodeAPI }()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	_, err := leetcodeAPI.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not decode response")
}

func TestValidateUsername_InvalidResponseFormat(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalidKey": "invalidValue"}`))
	}))
	defer testServer.Close()

	oldLeetcodeAPI := config.Leetcode_API
	defer func() { config.Leetcode_API = oldLeetcodeAPI }()

	config.Leetcode_API = testServer.URL

	leetcodeAPI := api.NewLeetcodeAPI()
	_, err := leetcodeAPI.ValidateUsername("testuser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid response format")
}
