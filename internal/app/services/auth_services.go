package services

import (
	"bytes"
	"cli-project/internal/domain/interfaces"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthService struct {
	userRepo interfaces.UserRepository
}

func NewAuthService(userRepo interfaces.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) IsEmailUnique(email string) (bool, error) {
	return s.userRepo.IsEmailUnique(email)
}

func (s *AuthService) IsUsernameUnique(username string) (bool, error) {
	return s.userRepo.IsUsernameUnique(username)
}

func (s *AuthService) IsLeetcodeIDUnique(leetcodeID string) (bool, error) {
	return s.userRepo.IsLeetcodeIDUnique(leetcodeID)
}

// ValidateLeetcodeUsername checks if the provided LeetCode username exists
func (s *AuthService) ValidateLeetcodeUsername(username string) (bool, error) {

	// GraphQL query to check if a LeetCode user exists
	const userQuery = `
	query getUserProfile($username: String!) {
  	matchedUser(username: $username) {
    username
  }
}
`

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
