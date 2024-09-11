package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("username or password incorrect")
	ErrUserNotFound       = errors.New("user not found")
	HashString            = utils.HashString
	VerifyString          = utils.VerifyString
)

type UserService struct {
	userRepo        interfaces.UserRepository
	questionService interfaces.QuestionService
	LeetcodeAPI     interfaces2.LeetcodeAPI
	//userWG   *sync.WaitGroup
}

func NewUserService(userRepo interfaces.UserRepository, questionService interfaces.QuestionService, LeetcodeAPI interfaces2.LeetcodeAPI) interfaces.UserService {
	return &UserService{
		userRepo:        userRepo,
		questionService: questionService,
		LeetcodeAPI:     LeetcodeAPI,
		//userWG:   &sync.WaitGroup{},
	}
}

// Signup creates a new user account
func (s *UserService) Signup(user *models.StandardUser) error {

	// Change username to lowercase for consistency
	user.StandardUser.Username = strings.ToLower(user.StandardUser.Username)

	// Change email to lower for consistency
	user.StandardUser.Email = strings.ToLower(user.StandardUser.Email)

	// Change org and country name to proper format
	user.StandardUser.Organisation = utils.CapitalizeWords(user.StandardUser.Organisation)
	user.StandardUser.Country = utils.CapitalizeWords(user.StandardUser.Country)

	// Generate a new UUID for the user
	userID := utils.GenerateUUID()
	user.StandardUser.ID = userID

	// Hash the user password
	hashedPassword, err := HashString(user.StandardUser.Password)
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
	username = utils.CleanString(username)

	// Retrieve the user by username
	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		// Check if the error is because the user was not found (PostgreSQL)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound // Return error if user doesn't exist
		}
		return fmt.Errorf("error fetching user: %v", err)
	}
	// Verify the password
	if !VerifyString(password, user.StandardUser.Password) {
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

	//update last seen of user
	user.LastSeen = time.Now().UTC()

	// update data in db
	err = s.userRepo.UpdateUserDetails(user)
	if err != nil {
		return errors.New("could not update user details")
	}

	// clear Active user id
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
func (s *UserService) UpdateUserProgress() error {
	// Validate if the ActiveUserID is a valid UUID
	userUUID, err := uuid.Parse(globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}
	// Fetch recent submissions from LeetCode API using GetStats
	stats, err := s.GetUserLeetcodeStats(globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("could not fetch stats from LeetCode API: %v", err)
	}
	// Extract title slugs from recent submissions
	recentSlugs := stats.RecentACSubmissionTitleSlugs

	// Check which of the recent slugs are in the questions table
	var validSlugs []string
	for _, slug := range recentSlugs {
		exists, err := s.questionService.QuestionExistsByTitleSlug(slug)
		if err != nil {
			return fmt.Errorf("could not check if question exists: %v", err)
		}
		if exists {
			validSlugs = append(validSlugs, slug)
		}
	}
	// Update user's progress with valid slugs
	if len(validSlugs) > 0 {
		err := s.userRepo.UpdateUserProgress(userUUID, validSlugs)
		if err != nil {
			return fmt.Errorf("could not update user progress: %v", err)
		}
	}
	return nil
}

func (s *UserService) CountActiveUserInLast24Hours() (int, error) {
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
	username = utils.CleanString(username)

	return s.userRepo.FetchUserByUsername(username)
}

func (s *UserService) GetUserByID(userID string) (*models.StandardUser, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	// Clean userID for consistency
	userID = utils.CleanString(userID)

	// Fetch user by ID from the repository
	return s.userRepo.FetchUserByID(userID)
}

func (s *UserService) GetUserRole(userID string) (string, error) {

	if userID == "" {
		return "", errors.New("userID is empty")
	}

	// Change userID to lowercase for consistency
	userID = utils.CleanString(userID)

	user, err := s.userRepo.FetchUserByID(userID)
	if err != nil {
		return "", err
	}

	return user.StandardUser.Role, nil
}

func (s *UserService) GetUserProgress(userID string) (*[]string, error) {

	userProgress, err := s.userRepo.FetchUserProgress(userID)
	if err != nil {
		return nil, err
	}

	return userProgress, nil

}

func (s *UserService) GetUserID(username string) (string, error) {
	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.StandardUser.ID, nil
}

func (s *UserService) BanUser(username string) (bool, error) {

	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		return false, err
	}

	role := user.StandardUser.Role

	if role == roles.ADMIN {
		return false, errors.New("ban operation on admin not allowed")
	}

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

	user, err := s.userRepo.FetchUserByUsername(username)
	if err != nil {
		return false, err
	}

	role := user.StandardUser.Role

	if role == roles.ADMIN {
		return false, errors.New("unban operation on admin not allowed")
	}

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

func (s *UserService) GetUserLeetcodeStats(userID string) (*models.LeetcodeStats, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	LeetcodeID := user.LeetcodeID

	return s.LeetcodeAPI.GetStats(LeetcodeID)
}

func (s *UserService) GetUserCodesageStats(userID string) (*models.CodesageStats, error) {
	// Get user progress
	userProgress, err := s.GetUserProgress(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %v", err)
	}

	// Calculate the total number of questions the user has done
	totalQuestionsDoneCount := len(*userProgress)

	// Get the total number of questions on the platform
	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total questions count: %v", err)
	}

	// Initialize variables for difficulty counts and tag-wise stats
	var easyDoneCount, mediumDoneCount, hardDoneCount int
	topicWiseStats := make(map[string]int)
	companyWiseStats := make(map[string]int)

	// Loop through the user's completed questions
	for _, titleSlug := range *userProgress {
		// Get the question details by title_slug
		question, err := s.questionService.GetQuestionByID(titleSlug)
		if err != nil {
			return nil, fmt.Errorf("failed to get question details for %s: %v", titleSlug, err)
		}

		// Increment the difficulty count
		switch question.Difficulty {
		case "easy":
			easyDoneCount++
		case "medium":
			mediumDoneCount++
		case "hard":
			hardDoneCount++
		}

		// Accumulate topic tags
		for _, tag := range question.TopicTags {
			topicWiseStats[tag]++
		}

		// Accumulate company tags
		for _, company := range question.CompanyTags {
			companyWiseStats[company]++
		}
	}

	// Create the final CodesageStats struct
	stats := &models.CodesageStats{
		TotalQuestionsCount:     totalQuestionsCount,
		TotalQuestionsDoneCount: totalQuestionsDoneCount,
		TotalEasyCount:          easyDoneCount,
		TotalMediumCount:        mediumDoneCount,
		TotalHardCount:          hardDoneCount,
		EasyDoneCount:           easyDoneCount,
		MediumDoneCount:         mediumDoneCount,
		HardDoneCount:           hardDoneCount,
		CompanyWiseStats:        companyWiseStats,
		TopicWiseStats:          topicWiseStats,
	}

	return stats, nil
}

func (s *UserService) GetPlatformStats() (*models.PlatformStats, error) {
	// Get active users in the last 24 hours
	activeUsersInLast24Hours, err := s.CountActiveUserInLast24Hours()
	if err != nil {
		return nil, err
	}

	// Get total questions count
	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount()
	if err != nil {
		return nil, err
	}

	// Get all questions
	allQuestions, err := s.questionService.GetAllQuestions()
	if err != nil {
		return nil, err
	}

	// Initialize maps to hold counts
	difficultyWiseCount := make(map[string]int)
	topicWiseCount := make(map[string]int)
	companyWiseCount := make(map[string]int)

	// Populate maps with counts from allQuestions
	for _, question := range *allQuestions {
		// Count difficulty-wise questions
		difficultyWiseCount[question.Difficulty]++

		// Count topic-wise questions
		for _, topic := range question.TopicTags {
			topicWiseCount[topic]++
		}

		// Count company-wise questions
		for _, company := range question.CompanyTags {
			companyWiseCount[company]++
		}
	}

	// Generate platform stats
	platformStats := &models.PlatformStats{
		ActiveUserInLast24Hours:      activeUsersInLast24Hours,
		TotalQuestionsCount:          totalQuestionsCount,
		DifficultyWiseQuestionsCount: difficultyWiseCount,
		TopicWiseQuestionsCount:      topicWiseCount,
		CompanyWiseQuestionsCount:    companyWiseCount,
	}

	// make sure all three difficulties are there as keys in map or else set missing difficulty key with value 0

	return platformStats, nil
}

//func (s *UserService) WaitForCompletion() {
//	s.userWG.Wait()
//}
