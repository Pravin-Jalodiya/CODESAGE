package service_test

import (
	"cli-project/internal/app/services"
	"cli-project/internal/domain/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuestionService_RemoveQuestionByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock expectations
	mockQuestionRepo.EXPECT().QuestionExists("q1").Return(true, nil)
	mockQuestionRepo.EXPECT().RemoveQuestionByID("q1").Return(nil)

	// Execute
	err := questionService.RemoveQuestionByID("q1")

	// Assert
	assert.Nil(t, err)
}

func TestQuestionService_GetQuestionByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	questionID := "q1"
	mockQuestion := &models.Question{
		QuestionID:    "q1",
		QuestionTitle: "Title1",
		Difficulty:    "easy",
		QuestionLink:  "http://example.com/1",
		TopicTags:     []string{"topic1"},
		CompanyTags:   []string{"company1"},
	}

	mockQuestionRepo.EXPECT().QuestionExists(questionID).Return(true, nil).AnyTimes()
	mockQuestionRepo.EXPECT().FetchQuestionByID(questionID).Return(mockQuestion, nil)

	// Execute
	question, err := questionService.GetQuestionByID(questionID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, mockQuestion, question)
}

func TestQuestionService_GetAllQuestions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	mockQuestions := []models.Question{
		{QuestionID: "q1", QuestionTitle: "Title1", Difficulty: "easy"},
		{QuestionID: "q2", QuestionTitle: "Title2", Difficulty: "medium"},
	}

	mockQuestionRepo.EXPECT().FetchAllQuestions().Return(&mockQuestions, nil)

	// Execute
	questions, err := questionService.GetAllQuestions()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, &mockQuestions, questions)
}

func TestQuestionService_GetQuestionsByFilters(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	difficulty := "easy"
	company := "company1"
	topic := "topic1"

	mockQuestions := []models.Question{
		{QuestionID: "q1", QuestionTitle: "Title1", Difficulty: "easy", CompanyTags: []string{company}, TopicTags: []string{topic}},
	}

	mockQuestionRepo.EXPECT().FetchQuestionsByFilters(difficulty, company, topic).Return(&mockQuestions, nil)

	// Execute
	questions, err := questionService.GetQuestionsByFilters(difficulty, company, topic)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, &mockQuestions, questions)
}

func TestQuestionService_QuestionExists(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Mock expectations
	mockQuestionRepo.EXPECT().QuestionExists("q1").Return(true, nil)

	// Execute
	exists, err := questionService.QuestionExists("q1")

	// Assert
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestQuestionService_GetTotalQuestionsCount(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	count := int64(10)

	mockQuestionRepo.EXPECT().CountQuestions().Return(count, nil)

	// Execute
	totalCount, err := questionService.GetTotalQuestionsCount()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, count, totalCount)
}

func TestQuestionService_AddQuestionsFromFile(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Define the CSV data
	//csvData := [][]string{
	//	{"QuestionID", "QuestionTitle", "Difficulty", "QuestionLink", "TopicTags", "CompanyTags"},
	//	{"q1", "Title1", "easy", "http://example.com/1", "topic1", "company1"},
	//	{"q2", "Title2", "medium", "http://example.com/2", "topic2", "company2"},
	//}

	// Mock the CSV reader function
	//mockReadCSV := func(filePath string) ([][]string, error) {
	//	return csvData, nil
	//}

	// Create the service with the mock reader function
	questionService = services.NewQuestionService(mockQuestionRepo)

	// Set expectations for the repository
	mockQuestionRepo.EXPECT().QuestionExists("q1").Return(false, nil)
	mockQuestionRepo.EXPECT().QuestionExists("q2").Return(false, nil)
	mockQuestionRepo.EXPECT().AddQuestions(gomock.Any()).Return(nil)

	// Execute
	newQuestionsAdded, err := questionService.AddQuestionsFromFile("dummy/path")

	// Assert
	assert.Nil(t, err)
	assert.True(t, newQuestionsAdded)

	// Test with a CSV entry for an existing question
	mockQuestionRepo.EXPECT().QuestionExists("q1").Return(true, nil)
	mockQuestionRepo.EXPECT().AddQuestions(gomock.Any()).Return(nil)

	newQuestionsAdded, err = questionService.AddQuestionsFromFile("dummy/path")
	assert.Nil(t, err)
	assert.False(t, newQuestionsAdded)
}
