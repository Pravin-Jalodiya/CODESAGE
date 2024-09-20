package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
}

func NewUserService(
	userRepo interfaces.UserRepository,
	questionService interfaces.QuestionService,
	LeetcodeAPI interfaces2.LeetcodeAPI,
) interfaces.UserService {
	return &UserService{
		userRepo:        userRepo,
		questionService: questionService,
		LeetcodeAPI:     LeetcodeAPI,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) (*[]models.StandardUser, error) {
	return s.userRepo.FetchAllUsers(ctx)
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard(ctx context.Context) error {
	return nil
}

// UpdateUserProgress updates the user's progress by adding a solved question ID.
func (s *UserService) UpdateUserProgress(ctx context.Context) error {
	userUUID, err := uuid.Parse(globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	stats, err := s.GetUserLeetcodeStats(globals.ActiveUserID)
	if err != nil {
		return fmt.Errorf("could not fetch stats from LeetCode API: %v", err)
	}

	recentSlugs := stats.RecentACSubmissionTitleSlugs
	var validSlugs []string
	for _, slug := range recentSlugs {
		exists, err := s.questionService.QuestionExistsByTitleSlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("could not check if question exists: %v", err)
		}
		if exists {
			validSlugs = append(validSlugs, slug)
		}
	}

	if len(validSlugs) > 0 {
		err := s.userRepo.UpdateUserProgress(ctx, userUUID, validSlugs)
		if err != nil {
			return fmt.Errorf("could not update user progress: %v", err)
		}
	}
	return nil
}

func (s *UserService) CountActiveUserInLast24Hours(ctx context.Context) (int, error) {
	return s.userRepo.CountActiveUsersInLast24Hours(ctx)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.StandardUser, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

	username = utils.CleanString(username)
	return s.userRepo.FetchUserByUsername(ctx, username)
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.StandardUser, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	userID = utils.CleanString(userID)
	return s.userRepo.FetchUserByID(ctx, userID)
}

func (s *UserService) GetUserRole(ctx context.Context, userID string) (roles.Role, error) {
	if userID == "" {
		return -1, errors.New("userID is empty")
	}

	userID = utils.CleanString(userID)
	user, err := s.userRepo.FetchUserByID(ctx, userID)
	if err != nil {
		return -1, err
	}

	role, err := roles.ParseRole(user.StandardUser.Role)
	if err != nil {
		return -1, err
	}

	return role, nil
}

func (s *UserService) GetUserProgress(ctx context.Context, userID string) (*[]string, error) {
	return s.userRepo.FetchUserProgress(ctx, userID)
}

func (s *UserService) GetUserID(ctx context.Context, username string) (string, error) {
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return user.StandardUser.ID, nil
}

func (s *UserService) BanUser(ctx context.Context, username string) (bool, error) {
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	role, err := roles.ParseRole(user.StandardUser.Role)
	if err != nil {
		return false, err
	}

	if role == roles.ADMIN {
		return false, errors.New("ban operation on admin not allowed")
	}

	userID, err := s.GetUserID(ctx, username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(ctx, userID)
	if err != nil {
		return false, err
	}

	if alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.BanUser(ctx, userID)
}

func (s *UserService) UnbanUser(ctx context.Context, username string) (bool, error) {
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	role, err := roles.ParseRole(user.StandardUser.Role)
	if err != nil {
		return false, err
	}

	if role == roles.ADMIN {
		return false, errors.New("unban operation on admin not allowed")
	}

	userID, err := s.GetUserID(ctx, username)
	if err != nil {
		return false, err
	}

	alreadyBanned, err := s.IsUserBanned(ctx, userID)
	if err != nil {
		return false, err
	}

	if !alreadyBanned {
		return true, nil
	}

	return false, s.userRepo.UnbanUser(ctx, userID)
}

func (s *UserService) IsUserBanned(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FetchUserByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.StandardUser.IsBanned, nil
}

func (s *UserService) GetUserLeetcodeStats(userID string) (*models.LeetcodeStats, error) {
	user, err := s.GetUserByID(context.Background(), userID) // use default context as this is external API call
	if err != nil {
		return nil, err
	}

	LeetcodeID := user.LeetcodeID
	return s.LeetcodeAPI.GetStats(LeetcodeID)
}

func (s *UserService) GetUserCodesageStats(ctx context.Context, userID string) (*models.CodesageStats, error) {
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %v", err)
	}

	totalQuestionsDoneCount := len(*userProgress)
	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total questions count: %v", err)
	}

	var easyDoneCount, mediumDoneCount, hardDoneCount int
	topicWiseStats := make(map[string]int)
	companyWiseStats := make(map[string]int)

	for _, titleSlug := range *userProgress {
		question, err := s.questionService.GetQuestionByID(ctx, titleSlug)
		if err != nil {
			return nil, fmt.Errorf("failed to get question details for %s: %v", titleSlug, err)
		}

		switch question.Difficulty {
		case "easy":
			easyDoneCount++
		case "medium":
			mediumDoneCount++
		case "hard":
			hardDoneCount++
		}

		for _, tag := range question.TopicTags {
			topicWiseStats[tag]++
		}

		for _, company := range question.CompanyTags {
			companyWiseStats[company]++
		}
	}

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

func (s *UserService) GetPlatformStats(ctx context.Context) (*models.PlatformStats, error) {
	activeUsersInLast24Hours, err := s.CountActiveUserInLast24Hours(ctx)
	if err != nil {
		return nil, err
	}

	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount(ctx)
	if err != nil {
		return nil, err
	}

	allQuestions, err := s.questionService.GetAllQuestions(ctx)
	if err != nil {
		return nil, err
	}

	difficultyWiseCount := make(map[string]int)
	topicWiseCount := make(map[string]int)
	companyWiseCount := make(map[string]int)

	for _, question := range *allQuestions {
		difficultyWiseCount[question.Difficulty]++

		for _, topic := range question.TopicTags {
			topicWiseCount[topic]++
		}

		for _, company := range question.CompanyTags {
			companyWiseCount[company]++
		}
	}

	platformStats := &models.PlatformStats{
		ActiveUserInLast24Hours:      activeUsersInLast24Hours,
		TotalQuestionsCount:          totalQuestionsCount,
		DifficultyWiseQuestionsCount: difficultyWiseCount,
		TopicWiseQuestionsCount:      topicWiseCount,
		CompanyWiseQuestionsCount:    companyWiseCount,
	}

	return platformStats, nil
}
