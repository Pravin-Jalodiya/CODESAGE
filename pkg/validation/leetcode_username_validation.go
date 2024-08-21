package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphQL query to check if a LeetCode user exists
const userQuery = `
query getUserProfile($username: String!) {
  matchedUser(username: $username) {
    username
  }
}
`

// ValidateLeetCodeUsername checks if the provided LeetCode username exists
func ValidateLeetCodeUsername(username string) (bool, error) {
	apiURL := "https://leetcode.com/graphql/"

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
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

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
