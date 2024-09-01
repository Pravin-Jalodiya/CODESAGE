package service_test

import (
	"cli-project/external/api"
	interfaces2 "cli-project/external/domain/interfaces"
	"cli-project/internal/app/services"
	"cli-project/internal/domain/interfaces"
	mock_interfaces "cli-project/tests/mocks/repository"
	mock_services "cli-project/tests/mocks/services"
	"github.com/golang/mock/gomock"
	"testing"
)

var (
	ctrl                *gomock.Controller
	mockUserRepo        *mock_interfaces.MockUserRepository
	mockQuestionRepo    *mock_interfaces.MockQuestionRepository
	mockUserService     *mock_services.MockUserService
	mockQuestionService *mock_services.MockQuestionService
	mockAuthService     *mock_services.MockAuthService
	mockLeetcodeAPI     *mock_services.MockLeetcodeAPI
	userService         interfaces.UserService
	questionService     interfaces.QuestionService
	authService         interfaces.AuthService
	LeetcodeAPI         interfaces2.LeetcodeAPI
)

func setup(t *testing.T) func() {
	// Set up the go mock controller
	ctrl = gomock.NewController(t)

	// Create mock repositories
	mockUserRepo = mock_interfaces.NewMockUserRepository(ctrl)
	mockQuestionRepo = mock_interfaces.NewMockQuestionRepository(ctrl)

	// Create mock services
	mockUserService = mock_services.NewMockUserService(ctrl)
	mockQuestionService = mock_services.NewMockQuestionService(ctrl)
	mockAuthService = mock_services.NewMockAuthService(ctrl)
	mockLeetcodeAPI = mock_services.NewMockLeetcodeAPI(ctrl)
	LeetcodeAPI = mock_services.NewMockLeetcodeAPI(ctrl)

	// Create Genuine Services
	userService = services.NewUserService(mockUserRepo, mockQuestionService, mockLeetcodeAPI)
	questionService = services.NewQuestionService(mockQuestionRepo)
	authService = services.NewAuthService(mockUserRepo, mockLeetcodeAPI)
	LeetcodeAPI = api.NewLeetcodeAPI()

	// Return a cleanup function to be called at the end of the test
	return func() {
		ctrl.Finish()
	}
}
