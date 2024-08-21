package repositories

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
)

type questionRepo struct {
}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{}
}

func (r *questionRepo) AddQuestionsByID(questionID []int) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) AddQuestionsFromFile(filePath string) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) RemoveQuestionsByID(questionID []int) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) FetchQuestionByID(questionID int) (models.Question, error) {
	// Placeholder implementation
	return models.Question{}, nil
}

func (r *questionRepo) FetchAllQuestions() ([]models.Question, error) {
	// Placeholder implementation
	return nil, nil
}

func (r *questionRepo) FetchQuestionsByFilters(difficulty, company, topic string) ([]models.Question, error) {
	// Placeholder implementation
	return nil, nil
}
