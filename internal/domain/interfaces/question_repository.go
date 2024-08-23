package interfaces

import "cli-project/internal/domain/models"

type QuestionRepository interface {
	AddQuestionsByID(questionID []int) error
	AddQuestions([]models.Question) error
	RemoveQuestionByID(questionID string) error
	FetchQuestionByID(questionID int) (models.Question, error)
	FetchAllQuestions() ([]models.Question, error)
	FetchQuestionsByFilters(difficulty, company, topic string) ([]models.Question, error)
}
