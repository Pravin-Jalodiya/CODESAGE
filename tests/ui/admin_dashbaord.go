package ui

import (
	"bytes"
	"cli-project/pkg/utils/formatting"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CountActiveUserInLast24Hours() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// MockQuestionService is a mock implementation of the QuestionService interface
type MockQuestionService struct {
	mock.Mock
}

func (m *MockQuestionService) GetTotalQuestionsCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// MockReader is a mock implementation of the Reader interface
type MockReader struct {
	mock.Mock
}

func (m *MockReader) ReadString(delim byte) (string, error) {
	args := m.Called(delim)
	return args.String(0), args.Error(1)
}

func TestShowAdminDashboard(t *testing.T) {
	mockUserService := new(MockUserService)
	mockQuestionService := new(MockQuestionService)
	mockReader := new(MockReader)

	// Redirect stdout to capture output
	var out bytes.Buffer
	fmt.SetOut(&out)

	ui := &UI{
		userService:     mockUserService,
		questionService: mockQuestionService,
		reader:          mockReader,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserService.On("CountActiveUserInLast24Hours").Return(10, nil)
		mockQuestionService.On("GetTotalQuestionsCount").Return(500, nil)
		mockReader.On("ReadString", byte('\n')).Return("", nil)

		ui.ShowAdminDashboard()

		output := out.String()

		assert.Contains(t, output, formatting.Colorize("Active Users (Last 24 Hours): 10", "cyan", "bold"))
		assert.Contains(t, output, formatting.Colorize("Total Questions on the Platform: 500", "cyan", "bold"))
		mockUserService.AssertExpectations(t)
		mockQuestionService.AssertExpectations(t)
		mockReader.AssertExpectations(t)
	})

	t.Run("Error fetching active users", func(t *testing.T) {
		mockUserService.On("CountActiveUserInLast24Hours").Return(0, errors.New("some error"))
		mockQuestionService.On("GetTotalQuestionsCount").Return(500, nil)
		mockReader.On("ReadString", byte('\n')).Return("", nil)

		ui.ShowAdminDashboard()

		output := out.String()

		assert.Contains(t, output, formatting.Colorize("Error fetching active users count: ", "red", "bold"))
		mockUserService.AssertExpectations(t)
		mockQuestionService.AssertExpectations(t)
		mockReader.AssertExpectations(t)
	})

	t.Run("Error fetching total questions", func(t *testing.T) {
		mockUserService.On("CountActiveUserInLast24Hours").Return(10, nil)
		mockQuestionService.On("GetTotalQuestionsCount").Return(0, errors.New("some error"))
		mockReader.On("ReadString", byte('\n')).Return("", nil)

		ui.ShowAdminDashboard()

		output := out.String()

		assert.Contains(t, output, formatting.Colorize("Error fetching total questions count: ", "red", "bold"))
		mockUserService.AssertExpectations(t)
		mockQuestionService.AssertExpectations(t)
		mockReader.AssertExpectations(t)
	})
}
