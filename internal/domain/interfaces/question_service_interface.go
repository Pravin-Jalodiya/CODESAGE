package interfaces

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
)

type QuestionService interface {
	AddQuestionsFromFile(questionFilePath string) (bool, error)
	RemoveQuestionByID(questionID string) error
	GetQuestionByID(questionID string) (*models.Question, error)
	GetAllQuestions() (*[]dto.Question, error)
	GetQuestionsByFilters(difficulty, company, topic string) (*[]dto.Question, error)
	QuestionExistsByID(questionID string) (bool, error)
	QuestionExistsByTitleSlug(string) (bool, error)
	GetTotalQuestionsCount() (int64, error)
}
