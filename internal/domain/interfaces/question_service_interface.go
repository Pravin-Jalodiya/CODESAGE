package interfaces

import (
	"cli-project/internal/domain/models"
)

type QuestionService interface {
	AddQuestionsFromFile(questionFilePath string) (bool, error)
	RemoveQuestionByID(questionID string) error
	GetQuestionByID(questionID string) (*models.Question, error)
	GetAllQuestions() (*[]models.Question, error)
	GetQuestionsByFilters(difficulty, company, topic string) (*[]models.Question, error)
	QuestionExists(questionID string) (bool, error)
	GetTotalQuestionsCount() (int64, error)
}
