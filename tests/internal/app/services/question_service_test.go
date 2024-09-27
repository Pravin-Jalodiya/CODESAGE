package service_test

import (
	"cli-project/internal/app/services"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuestionService_AddQuestionsFromFile(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "easy", "https://leetcode.com/question1", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		validQuestion := []models.Question{
			{
				QuestionTitleSlug: "slug1",
				QuestionID:        "1",
				QuestionTitle:     "title1",
				Difficulty:        "easy",
				QuestionLink:      "https://leetcode.com/question1",
				TopicTags:         []string{"tag1", "tag2"},
				CompanyTags:       []string{"comp1", "comp2"},
			},
		}

		gomock.InOrder(
			mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(false, nil),
			mockQuestionRepo.EXPECT().AddQuestions(&validQuestion).Return(nil),
		)

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.True(t, success)
		assert.NoError(t, err)
	})

	t.Run("CSV Reader Error", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return nil, errors.New("error reading CSV file")
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "error reading CSV file: error reading CSV file")
	})

	t.Run("Invalid CSV Format", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "easy", "https://leetcode.com/question1", "invalid"},
			}, nil
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "invalid CSV format, expected 7 columns")
	})

	t.Run("Invalid Question ID", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "invalid-id", "title1", "easy", "https://leetcode.com/question1", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		services.ValidateQuestionID = func(id string) (bool, error) {
			return false, errors.New("invalid question ID")
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "invalid question ID: invalid question ID")
	})

	t.Run("Invalid Difficulty", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "invalid-difficulty", "https://leetcode.com/question1", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		services.ValidateQuestionID = func(id string) (bool, error) {
			return true, nil
		}

		services.ValidateQuestionDifficulty = func(difficulty string) (string, error) {
			return "", errors.New("invalid difficulty")
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "invalid difficulty: invalid difficulty")
	})

	t.Run("Invalid Question Link", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "easy", "invalid-link", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		services.ValidateQuestionID = func(id string) (bool, error) {
			return true, nil
		}

		services.ValidateQuestionDifficulty = func(difficulty string) (string, error) {
			return difficulty, nil
		}

		services.ValidateQuestionLink = func(link string) (string, error) {
			return "", errors.New("invalid question link")
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "invalid question link: invalid question link")
	})

	t.Run("Error Checking Question Existence", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "easy", "https://leetcode.com/question1", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		gomock.InOrder(
			mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(false, errors.New("error")),
		)

		services.ValidateQuestionID = func(id string) (bool, error) {
			return true, nil
		}

		services.ValidateQuestionDifficulty = func(difficulty string) (string, error) {
			return difficulty, nil
		}

		services.ValidateQuestionLink = func(link string) (string, error) {
			return link, nil
		}

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "error checking if question exists: error")
	})

	t.Run("Error Adding Questions To Database", func(t *testing.T) {
		services.CSVReader = func(filePath string) ([][]string, error) {
			return [][]string{
				{"header1", "header2", "header3", "header4", "header5", "header6", "header7"},
				{"slug1", "1", "title1", "easy", "https://leetcode.com/question1", "tag1,tag2", "comp1,comp2"},
			}, nil
		}

		validQuestion := []models.Question{
			{
				QuestionTitleSlug: "slug1",
				QuestionID:        "1",
				QuestionTitle:     "title1",
				Difficulty:        "easy",
				QuestionLink:      "https://leetcode.com/question1",
				TopicTags:         []string{"tag1", "tag2"},
				CompanyTags:       []string{"comp1", "comp2"},
			},
		}

		services.ValidateQuestionID = func(id string) (bool, error) {
			return true, nil
		}

		services.ValidateQuestionDifficulty = func(difficulty string) (string, error) {
			return difficulty, nil
		}

		services.ValidateQuestionLink = func(link string) (string, error) {
			return link, nil
		}

		gomock.InOrder(
			mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(false, nil),
			mockQuestionRepo.EXPECT().AddQuestions(&validQuestion).Return(errors.New("database error")),
		)

		success, err := questionService.AddQuestionsFromFile("path/to/csv")
		assert.False(t, success)
		assert.EqualError(t, err, "error adding questions to the database: database error")
	})
}

func TestQuestionService_RemoveQuestionByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		gomock.InOrder(
			mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(true, nil),
			mockQuestionRepo.EXPECT().RemoveQuestionByID("1").Return(nil),
		)

		err := questionService.RemoveQuestionByID("1")
		assert.NoError(t, err)
	})

	t.Run("Question Not Found", func(t *testing.T) {
		mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(false, nil)

		err := questionService.RemoveQuestionByID("1")
		assert.Error(t, err)
		assert.EqualError(t, err, "question with ID 1 not found")
	})

	t.Run("Error Checking Question Existence", func(t *testing.T) {
		mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(false, errors.New("error"))

		err := questionService.RemoveQuestionByID("1")
		assert.Error(t, err)
		assert.EqualError(t, err, "error checking if question exists: error")
	})
}

func TestQuestionService_GetQuestionByID(t *testing.T) {
	teardown := setup(t)

	defer teardown()

	t.Run("Success", func(t *testing.T) {
		question := &models.Question{
			QuestionID:        "1",
			QuestionTitleSlug: "slug1",
		}

		services.ValidateTitleSlug = func(id string) (bool, error) {
			return true, nil
		}

		mockQuestionRepo.EXPECT().QuestionExistsByTitleSlug("slug1").Return(true, nil)

		mockQuestionRepo.EXPECT().FetchQuestionByTitleSlug("slug1").Return(question, nil)

		result, err := questionService.GetQuestionByID("slug1")
		assert.NoError(t, err)
		assert.Equal(t, question, result)
	})

	t.Run("Question Not Found", func(t *testing.T) {

		mockQuestionRepo.EXPECT().QuestionExistsByTitleSlug("slug").Return(false, nil)

		result, err := questionService.GetQuestionByID("slug")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "question with title slug slug not found")
	})

	t.Run("Ertror Checking Quesion Existence", func(t *testing.T) {

		mockQuestionRepo.EXPECT().QuestionExistsByTitleSlug("slug").Return(false, errors.New("error"))

		result, err := questionService.GetQuestionByID("slug")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestQuestionService_GetAllQuestions(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		questions := &[]dto.Question{
			{QuestionID: "1"},
		}
		mockQuestionRepo.EXPECT().FetchAllQuestions().Return(questions, nil)

		result, err := questionService.GetAllQuestions()
		assert.NoError(t, err)
		assert.Equal(t, questions, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockQuestionRepo.EXPECT().FetchAllQuestions().Return(nil, errors.New("fetch error"))

		result, err := questionService.GetAllQuestions()
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestQuestionService_GetQuestionsByFilters(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		questions := &[]dto.Question{
			{QuestionID: "1"},
		}
		mockQuestionRepo.EXPECT().FetchQuestionsByFilters("easy", "topic", "company").Return(questions, nil)

		result, err := questionService.GetQuestionsByFilters("easy", "topic", "company")
		assert.NoError(t, err)
		assert.Equal(t, questions, result)
	})

	t.Run("Validation Error", func(t *testing.T) {
		services.ValidateQuestionDifficulty = func(difficulty string) (string, error) {
			return "", errors.New("invalid difficulty")
		}

		defer func() {
			services.ValidateQuestionDifficulty = originalValidateQuestionDifficulty
		}()

		result, err := questionService.GetQuestionsByFilters("invalid-difficulty", "topic", "company")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestQuestionService_QuestionExistsByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		mockQuestionRepo.EXPECT().QuestionExistsByID("1").Return(true, nil)

		result, err := questionService.QuestionExistsByID("1")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("Validation Error", func(t *testing.T) {
		services.ValidateQuestionID = func(id string) (bool, error) {
			return false, errors.New("invalid question ID")
		}

		defer func() {
			services.ValidateQuestionID = originalValidateQuestionID
		}()

		result, err := questionService.QuestionExistsByID("invalid!")
		assert.Error(t, err)
		assert.False(t, result)
	})
}

func TestQuestionService_QuestionExistsByTitleSlug(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {

		services.ValidateTitleSlug = func(id string) (bool, error) {
			return true, nil
		}

		mockQuestionRepo.EXPECT().QuestionExistsByTitleSlug("slug").Return(true, nil)

		result, err := questionService.QuestionExistsByTitleSlug("slug")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("Validation Error", func(t *testing.T) {

		services.ValidateTitleSlug = func(id string) (bool, error) {
			return false, errors.New("invalid title slug")
		}

		result, err := questionService.QuestionExistsByTitleSlug("invalid_slug")
		assert.Error(t, err)
		assert.False(t, result)
	})
}

func TestQuestionService_GetTotalQuestionsCount(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		mockQuestionRepo.EXPECT().CountQuestions().Return(10, nil)

		count, err := questionService.GetTotalQuestionsCount()
		assert.NoError(t, err)
		assert.Equal(t, 10, count)
	})

	t.Run("Error", func(t *testing.T) {
		mockQuestionRepo.EXPECT().CountQuestions().Return(0, errors.New("count error"))

		count, err := questionService.GetTotalQuestionsCount()
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})
}
