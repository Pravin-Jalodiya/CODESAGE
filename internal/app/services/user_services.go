package services

import (
	"bytes"
	"cli-project/internal/config"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	"cli-project/pkg/utils/data_cleaning"
	pwd "cli-project/pkg/utils/password"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("username or password incorrect")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo        interfaces.UserRepository
	questionService *QuestionService
	//userWG   *sync.WaitGroup
}

func NewUserService(userRepo interfaces.UserRepository, questionService *QuestionService) *UserService {
	return &UserService{
		userRepo:        userRepo,
		questionService: questionService,
		//userWG:   &sync.WaitGroup{},
	}
}

// SignUp creates a new user account
func (s *UserService) SignUp(user *models.StandardUser) error {

	// Change username to lowercase for consistency
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)

	// Change email to lower for consistency
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)

	// Change org and country name to proper format
	user.StandardUser.Organisation = data_cleaning.CapitalizeWords(user.StandardUser.Organisation)
	user.StandardUser.Country = data_cleaning.CapitalizeWords(user.StandardUser.Country)

	// Generate a new UUID for the user
	userID := utils.GenerateUUID()
	user.StandardUser.ID = userID

	// Hash the user password
	hashedPassword, err := pwd.HashPassword(user.StandardUser.Password)
	if err != nil {
		return fmt.Errorf("could not hash password")
	}
	user.StandardUser.Password = hashedPassword

	// Set default role
	user.StandardUser.Role = "user"

	// Set default blocked status
	user.StandardUser.IsBanned = false

	// set question solved
	user.QuestionsSolved = []string{}

	// set last seen
	user.LastSeen = time.Now().UTC()

	// Register the user
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("could not register user")
	}

	return nil
}

// Login authenticates a user
func (s *UserService) Login(username, password string) error {

	// Change username to lowercase for consistency
	username = data_cleaning.CleanString(username)

	// Retrieve the user by username
	user, err := s.userRepo.FetchUserByUsername(username)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound // Return error if user doesn't exist
		}
		return fmt.Errorf("%v", err)
	}

	// Verify the password
	if !pwd.VerifyPassword(password, user.StandardUser.Password) {
		return ErrInvalidCredentials
	}

	return nil
}

func (s *UserService) Logout() error {
	// Get active user
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
	if err != nil {
		return errors.New("user not found")
	}

	// update last seen of user
	user.LastSeen = time.Now().UTC()

	// update data in db
	err = s.userRepo.UpdateUserDetails(user)
	if err != nil {
		return errors.New("could not update user details")
	}

	// clear active user
	globals.ActiveUserID = ""

	return nil
}

func (s *UserService) GetAllUsers() (*[]models.StandardUser, error) {
	return s.userRepo.FetchAllUsers()
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard() error {
	// Placeholder implementation
	return nil
}

// UpdateUserProgress updates the user's progress by adding a solved question ID.
func (s *UserService) UpdateUserProgress(solvedQuestionID string) (bool, error) {
	// Fetch the current user from the repository
	user, err := s.userRepo.FetchUserByID(globals.ActiveUserID)
	if err != nil {
		return false, fmt.Errorf("could not fetch user: %v", err)
	}

	// Check if the question ID is already in the user's progress
	for _, id := range user.QuestionsSolved {
		if id == solvedQuestionID {
			return false, nil // No need to update if the question ID is already in the list
		}
	}

	// Check if the question ID exists in the questions repository
	exists, err := s.questionService.QuestionExists(solvedQuestionID)
	if err != nil {
		return false, fmt.Errorf("could not check if question exists: %v", err)
	}
	if !exists {
		return false, fmt.Errorf("question with ID %s does not exist", solvedQuestionID)
	}

	// Update the user's progress
	return true, s.userRepo.UpdateUserProgress(solvedQuestionID)
}

func (s *UserService) CountActiveUserInLast24Hours() (int64, error) {
	count, err := s.userRepo.CountActiveUsersInLast24Hours()
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.StandardUser, error) {

	if username == "" {
		return nil, errors.New("username is empty")
	}

	// Change userID to lowercase for consistency
	username = data_cleaning.CleanString(username)

	return s.userRepo.FetchUserByUsername(username)
}

func (s *UserService) GetUserByID(userID string) (*models.StandardUser, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	// Clean userID for consistency
	userID = data_cleaning.CleanString(userID)

	// Fetch user by ID from the repository
	return s.userRepo.FetchUserByID(userID)
}

func (s *UserService) GetUserRole(userID string) (string, error) {

	if userID == "" {
		return "", errors.New("userID is empty")
	}

	// Change userID to lowercase for consistency
	userID = data_cleaning.CleanString(userID)

	user, err := s.userRepo.FetchUserByID(userID)
	if err != nil {
		return "", err
	}

	return user.StandardUser.Role, nil
}

func (s *UserService) GetUserID(username string) (string, error) {
	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.StandardUser.ID, nil
}

func (s *UserService) BanUser(username string) (bool, error) {

	userID, err := s.GetUserID(username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(userID)
	if err != nil {
		return false, err
	}

	if alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.BanUser(userID)
}

func (s *UserService) UnbanUser(username string) (bool, error) {

	userID, err := s.GetUserID(username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(userID)
	if err != nil {
		return false, err
	}

	if !alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.UnbanUser(userID)
}

func (s *UserService) IsUserBanned(userID string) (bool, error) {

	user, err := s.userRepo.FetchUserByID(userID)
	if err != nil {
		return false, err
	}

	return user.StandardUser.IsBanned, nil
}

func (s *UserService) GetLeetCodeStats(userID string) (*models.LeetcodeStats, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	leetcodeID := user.LeetcodeID
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

//func (s *UserService) WaitForCompletion() {
//	s.userWG.Wait()
//}
