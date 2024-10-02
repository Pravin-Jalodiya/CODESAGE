package services

import (
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	"cli-project/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	HashString   = utils.HashString
	VerifyString = utils.VerifyString
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

func (s *UserService) GetAllUsers(ctx context.Context) ([]dto.StandardUser, error) {
	users, err := s.userRepo.FetchAllUsers(ctx)

	var dtoUsers []dto.StandardUser

	for _, user := range users {
		role, _ := roles.ParseRole(user.Role)
		if role == roles.ADMIN {
			continue
		}

		dtoUser := dto.StandardUser{
			User: dto.User{
				Username:     user.Username,
				Role:         user.Role,
				Name:         user.Name,
				Email:        user.Email,
				Organisation: user.Organisation,
				Country:      user.Country,
				IsBanned:     user.IsBanned,
			},
			LeetcodeID: user.LeetcodeID,
			LastSeen:   user.LastSeen,
		}

		dtoUsers = append(dtoUsers, dtoUser)

	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return dtoUsers, nil
}

// ViewDashboard retrieves the dashboard for the active user
func (s *UserService) ViewDashboard(ctx context.Context) error {
	// Implementation for viewing the dashboard
	return nil
}

// UpdateUserProgress updates the user's progress by adding a solved question ID.
func (s *UserService) UpdateUserProgress(ctx context.Context, userID uuid.UUID) error {

	stats, err := s.GetUserLeetcodeStats(userID.String())
	if err != nil {
		return fmt.Errorf("%w: could not fetch stats from LeetCode API", err)
	}

	recentSlugs := stats.RecentACSubmissionTitleSlugs
	var validSlugs []string
	for _, slug := range recentSlugs {
		exists, err := s.questionService.QuestionExistsByTitleSlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("%w: could not check if question exists", err)
		}
		if exists {
			validSlugs = append(validSlugs, slug)
		}
	}

	err = s.userRepo.UpdateUserProgress(ctx, userID, validSlugs)
	if err != nil {
		return fmt.Errorf("%w: could not update user progress", err)
	}
	return nil
}

func (s *UserService) CountActiveUserInLast24Hours(ctx context.Context) (int, error) {
	count, err := s.userRepo.CountActiveUsersInLast24Hours(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return count, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.StandardUser, error) {
	if username == "" {
		return nil, fmt.Errorf("%w: username is empty", errs.ErrInvalidParameterError)
	}

	username = utils.CleanString(username)
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.StandardUser, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: user ID is empty", errs.ErrInvalidParameterError)
	}

	userID = utils.CleanString(userID)
	user, err := s.userRepo.FetchUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}
	return user, nil
}

func (s *UserService) GetUserRole(ctx context.Context, userID string) (roles.Role, error) {
	if userID == "" {
		return -1, fmt.Errorf("%w: userID is empty", errs.ErrInvalidParameterError)
	}

	userID = utils.CleanString(userID)
	user, err := s.userRepo.FetchUserByID(ctx, userID)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}

	role, err := roles.ParseRole(user.Role)
	if err != nil {
		return -1, fmt.Errorf("%w: %v", errs.ErrInvalidParameterError, err)
	}

	return role, nil
}

func (s *UserService) GetUserProgress(ctx context.Context, userID string) (*[]string, error) {
	progress, err := s.userRepo.FetchUserProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}
	return progress, nil
}

func (s *UserService) GetUserID(ctx context.Context, username string) (string, error) {
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}
	return user.ID, nil
}

func (s *UserService) UpdateUserBanState(ctx context.Context, username string) (string, error) {
	// Fetch the user by username
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}

	// Check the role to ensure we're not operating on an admin
	role, err := roles.ParseRole(user.Role)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrInvalidParameterError, err)
	}

	if role == roles.ADMIN {
		return "", fmt.Errorf("%w: ban/unban operation on admin not allowed", errs.ErrInvalidParameterError)
	}

	// Fetch user ID by username
	userID, err := s.GetUserID(ctx, username)
	if err != nil {
		return "", err
	}

	// Check if the user is already banned
	alreadyBanned, err := s.IsUserBanned(ctx, userID)
	if err != nil {
		return "", err
	}

	// If the user is banned, unban them
	if alreadyBanned {
		err = s.userRepo.UnbanUser(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("%w: %v", errs.ErrDbError, err)
		}
		return "User has been unbanned successfully", nil
	}

	// If the user is not banned, ban them
	err = s.userRepo.BanUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errs.ErrDbError, err)
	}

	return "User has been banned successfully", nil
}

// DeleteUser deletes a user by the given ID
func (s *UserService) DeleteUser(ctx context.Context, username string) error {

	username = utils.CleanString(username)
	user, err := s.userRepo.FetchUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		} else {
			return fmt.Errorf("%w: %v", errs.ErrFetchingUserFailed, err)
		}
	}
	err = s.userRepo.DeleteUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return fmt.Errorf("delete user failed: %w", errs.ErrUserNotFound)
		}
		return fmt.Errorf("delete user failed: %w", err)
	}
	return nil
}

// UpdateUser updates the user's profile.
func (s *UserService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	if username, ok := updates["username"].(string); ok {
		exists, err := s.userRepo.IsUsernameUnique(ctx, username)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w: %v", errs.ErrUserNameAlreadyExists, err)
		}
	}

	if email, ok := updates["email"].(string); ok {
		exists, err := s.userRepo.IsEmailUnique(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w: %v", errs.ErrEmailAlreadyExists, err)
		}
	}

	err := s.userRepo.UpdateUserProfile(ctx, userID, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) IsUserBanned(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FetchUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
	}

	return user.IsBanned, nil
}

func (s *UserService) GetUserLeetcodeStats(userID string) (*models.LeetcodeStats, error) {
	user, err := s.GetUserByID(context.Background(), userID) // use default context as this is external API call
	if err != nil {
		return nil, err
	}

	LeetcodeID := user.LeetcodeID
	stats, err := s.LeetcodeAPI.GetStats(LeetcodeID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrExternalAPI, err)
	}

	return stats, nil
}

func (s *UserService) GetUserCodesageStats(ctx context.Context, userID string) (*models.CodesageStats, error) {
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get user progress", err)
	}

	totalQuestionsDoneCount := len(*userProgress)
	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get total questions count", err)
	}

	var easyDoneCount, mediumDoneCount, hardDoneCount int
	topicWiseStats := make(map[string]int)
	companyWiseStats := make(map[string]int)

	for _, titleSlug := range *userProgress {
		question, err := s.questionService.GetQuestionByID(ctx, titleSlug)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get question details for %s", err, titleSlug)
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
		return nil, fmt.Errorf("%w: failed to count active users in the last 24 hours: %v", errs.ErrDbError, err)
	}

	totalQuestionsCount, err := s.questionService.GetTotalQuestionsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get total questions count: %v", errs.ErrDbError, err)
	}

	allQuestions, err := s.questionService.GetAllQuestions(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch all questions: %v", errs.ErrDbError, err)
	}

	// Initialize counts
	difficultyWiseCount := make(map[string]int)
	topicWiseCount := make(map[string]int)
	companyWiseCount := make(map[string]int)

	// Count the questions by difficulty, topic, and company
	for _, question := range allQuestions {
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
