package service_test

import (
	"cli-project/internal/app/services"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	pwd "cli-project/pkg/utils/password"
	mocks "cli-project/tests/mocks/repository"
	mock_services "cli-project/tests/mocks/services"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestUserService_Signup(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	user := models.StandardUser{
		StandardUser: models.User{
			Username:     "testuser",
			Email:        "testuser@example.com",
			Password:     "password",
			Organisation: "TestOrg",
			Country:      "TestCountry",
		},
		LeetcodeID:      "testLeetcode",
		QuestionsSolved: []string{},
		LastSeen:        time.Now().UTC(),
	}

	// Mock the expected behavior
	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(nil).Times(1)

	err := userService.Signup(&user)
	assert.NoError(t, err)
}

func TestUserService_Signup_Error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	user := models.StandardUser{
		StandardUser: models.User{
			Username:     "testuser",
			Email:        "testuser@example.com",
			Password:     "password",
			Organisation: "TestOrg",
			Country:      "TestCountry",
		},
		LeetcodeID:      "testLeetcode",
		QuestionsSolved: []string{},
		LastSeen:        time.Now().UTC(),
	}

	// Simulate an error during user creation
	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("could not register user")).Times(1)

	err := userService.Signup(&user)
	assert.Error(t, err)
	assert.Equal(t, "could not register user", err.Error())
}

func TestUserService_Login(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "testuser"
	password := "password123"

	// Hash the password for testing purposes
	hashedPassword, err := pwd.HashPassword(password)
	assert.NoError(t, err)

	// Mock the user repository response
	mockUserRepo.EXPECT().FetchUserByUsername(gomock.Any()).Return(&models.StandardUser{
		StandardUser: models.User{
			Username: "testuser",
			Password: hashedPassword, // Use the hashed password
		},
	}, nil).Times(1)

	// Call the actual Login function
	err = userService.Login(username, password)
	assert.NoError(t, err)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "testuser"
	password := "password123"
	wrongPassword := "wrongpassword"

	// Hash the password for testing purposes
	hashedPassword, err := pwd.HashPassword(password)
	assert.NoError(t, err)

	// Mock the user repository response
	mockUserRepo.EXPECT().FetchUserByUsername(gomock.Any()).Return(&models.StandardUser{
		StandardUser: models.User{
			Username: "testuser",
			Password: hashedPassword, // Use the hashed password
		},
	}, nil).Times(1)

	// Call the actual Login function with the wrong password
	err = userService.Login(username, wrongPassword)
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCredentials, err)
}

func TestUserService_UpdateUserProgress(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	solvedQuestionID := "123"

	mockUserRepo.EXPECT().FetchUserByID(gomock.Any()).Return(&models.StandardUser{
		QuestionsSolved: []string{},
	}, nil).Times(1)

	mockQuestionService.EXPECT().QuestionExists(solvedQuestionID).Return(true, nil).Times(1)

	mockUserRepo.EXPECT().UpdateUserProgress(solvedQuestionID).Return(nil).Times(1)

	updated, err := userService.UpdateUserProgress(solvedQuestionID)
	assert.NoError(t, err)
	assert.True(t, updated)
}

func TestUserService_UpdateUserProgress_QuestionNotExist(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	solvedQuestionID := "123"

	mockUserRepo.EXPECT().FetchUserByID(gomock.Any()).Return(&models.StandardUser{
		QuestionsSolved: []string{},
	}, nil).Times(1)

	mockQuestionService.EXPECT().QuestionExists(solvedQuestionID).Return(false, nil).Times(1)

	updated, err := userService.UpdateUserProgress(solvedQuestionID)
	assert.Error(t, err)
	assert.False(t, updated)
}

func TestUserService_GetLeetcodeStats(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "12345"

	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
		LeetcodeID: "Leetcode_user",
	}, nil).Times(1)

	mockLeetcodeAPI.EXPECT().GetStats("Leetcode_user").Return(&models.LeetcodeStats{
		EasyDoneCount:           10,
		MediumDoneCount:         20,
		HardDoneCount:           5,
		TotalEasyCount:          500,
		TotalHardCount:          300,
		TotalMediumCount:        200,
		TotalQuestionsCount:     1000,
		TotalQuestionsDoneCount: 35,
	}, nil).Times(1)

	stats, err := userService.GetLeetcodeStats(userID)
	assert.NoError(t, err)
	assert.Equal(t, 10, stats.EasyDoneCount)
	assert.Equal(t, 20, stats.MediumDoneCount)
	assert.Equal(t, 5, stats.HardDoneCount)
	assert.Equal(t, 500, stats.TotalEasyCount)
	assert.Equal(t, 300, stats.TotalHardCount)
	assert.Equal(t, 200, stats.TotalMediumCount)
	assert.Equal(t, 35, stats.TotalQuestionsDoneCount)
	assert.Equal(t, 1000, stats.TotalQuestionsCount)
}

func TestUserService_GetUserByUsername(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "testuser"

	// Mock the user repository response
	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(&models.StandardUser{
		StandardUser: models.User{
			Username: username,
		},
	}, nil).Times(1)

	user, err := userService.GetUserByUsername(username)
	assert.NoError(t, err)
	assert.Equal(t, username, user.StandardUser.Username)
}

func TestUserService_GetUserByUsername_NotFound(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "nonexistentuser"

	// Mock the user repository response
	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, mongo.ErrNoDocuments).Times(1)

	user, err := userService.GetUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestUserService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initialize mocks
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockLeetcodeAPI := mock_services.NewMockLeetcodeAPI(ctrl)
	mockQuestionService := mock_services.NewMockQuestionService(ctrl)

	// Create the UserService instance with mocks
	userService := services.NewUserService(mockUserRepo, mockQuestionService, mockLeetcodeAPI)

	hashedPassword, err := pwd.HashPassword("password123")

	// Create a standard user with the necessary fields
	standardUser := &models.StandardUser{
		StandardUser: models.User{
			ID:       "user-id",
			Username: "testuser",
			Password: hashedPassword,
			Email:    "test@example.com",
			Name:     "Test User",
			Role:     "user",
		},
		LeetcodeID:      "leet123",
		QuestionsSolved: []string{},
		LastSeen:        time.Time{},
	}

	// Set the active user ID globally
	globals.ActiveUserID = "user-id"

	// Set up expectations for the mock methods
	mockUserRepo.EXPECT().FetchUserByID("user-id").Return(standardUser, nil)
	mockUserRepo.EXPECT().UpdateUserDetails(gomock.Any()).Return(nil).Times(1)

	// Call the Logout method
	err = userService.Logout()

	// Assert that no errors occurred
	assert.NoError(t, err)

	// Assert that the LastSeen field was updated
	assert.NotEqual(t, time.Time{}, standardUser.LastSeen)

	// Assert that the ActiveUserID was cleared
	assert.Equal(t, "", globals.ActiveUserID)
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockUserRepo, nil, nil)

	userID := "user-id"

	hashedPassword, err := pwd.HashPassword("password123")

	// Create a standard user with the necessary fields
	standardUser := &models.StandardUser{
		StandardUser: models.User{
			ID:       "user-id",
			Username: "testuser",
			Password: hashedPassword,
			Email:    "test@example.com",
			Name:     "Test User",
			Role:     "user",
		},
		LeetcodeID:      "leet123",
		QuestionsSolved: []string{},
		LastSeen:        time.Time{},
	}

	mockUserRepo.EXPECT().FetchUserByID(userID).Return(standardUser, nil).Times(1)

	result, err := userService.GetUserByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, standardUser, result)
}

func TestUserService_GetUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockUserRepo, nil, nil)

	userID := "user-id"
	role := "user"
	hashedPassword, err := pwd.HashPassword("password123")

	// Create a standard user with the necessary fields
	standardUser := &models.StandardUser{
		StandardUser: models.User{
			ID:       "user-id",
			Username: "testuser",
			Password: hashedPassword,
			Email:    "test@example.com",
			Name:     "Test User",
			Role:     "user",
		},
		LeetcodeID:      "leet123",
		QuestionsSolved: []string{},
		LastSeen:        time.Time{},
	}

	mockUserRepo.EXPECT().FetchUserByID(userID).Return(standardUser, nil).Times(1)

	result, err := userService.GetUserRole(userID)

	assert.NoError(t, err)
	assert.Equal(t, role, result)
}

func TestUserService_BanUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockUserRepo, nil, nil)

	username := "testuser"
	userID := "user-id"

	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(&models.StandardUser{
		StandardUser: models.User{
			ID: userID,
		},
	}, nil).Times(1)

	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
		StandardUser: models.User{
			ID:       userID,
			IsBanned: false,
		},
	}, nil).Times(1)
	mockUserRepo.EXPECT().BanUser(userID).Return(nil).Times(1)

	banned, err := userService.BanUser(username)

	assert.NoError(t, err)
	assert.False(t, banned)
}

func TestUserService_UnbanUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockUserRepo, nil, nil)

	username := "testuser"
	userID := "awe1231"

	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(&models.StandardUser{
		StandardUser: models.User{
			ID: userID,
		},
	}, nil).Times(1)
	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
		StandardUser: models.User{
			ID:       userID,
			IsBanned: true,
		},
	}, nil).Times(1)
	mockUserRepo.EXPECT().UnbanUser(userID).Return(nil).Times(1)

	unbanned, err := userService.UnbanUser(username)

	assert.NoError(t, err)
	assert.False(t, unbanned)
}

func TestUserService_GetAllUsers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock the expected response
	mockUsers := []models.StandardUser{
		{StandardUser: models.User{Username: "user1"}},
		{StandardUser: models.User{Username: "user2"}},
	}
	mockUserRepo.EXPECT().FetchAllUsers().Return(&mockUsers, nil).Times(1)

	// Call the GetAllUsers method
	users, err := userService.GetAllUsers()

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 2, len(*users))
}

func TestUserService_ViewDashboard(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Placeholder test as ViewDashboard has no implementation
	err := userService.ViewDashboard()

	// Assert no errors (since it's a placeholder)
	assert.NoError(t, err)
}

func TestUserService_CountActiveUserInLast24Hours(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock the expected response
	mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(int64(5), nil).Times(1)

	// Call the CountActiveUserInLast24Hours method
	count, err := userService.CountActiveUserInLast24Hours()

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestUserService_GetUserID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "testuser"
	userID := "user-id"

	// Mock the expected response
	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(&models.StandardUser{
		StandardUser: models.User{
			ID: userID,
		},
	}, nil).Times(1)

	// Call the GetUserID method
	result, err := userService.GetUserID(username)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, userID, result)
}

func TestUserService_IsUserBanned(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	userID := "user-id"

	// Mock the expected response
	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
		StandardUser: models.User{
			IsBanned: true,
		},
	}, nil).Times(1)

	// Call the IsUserBanned method
	isBanned, err := userService.IsUserBanned(userID)

	// Assert the results
	assert.NoError(t, err)
	assert.True(t, isBanned)
}
