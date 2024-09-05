package interfaces

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
)

type QuestionRepository interface {
	AddQuestionsByID(*[]string) error
	AddQuestions(*[]models.Question) error
	RemoveQuestionByID(string) error
	FetchQuestionByID(string) (*models.Question, error)
	FetchAllQuestions() (*[]dto.Question, error)
	FetchQuestionsByFilters(string, string, string) (*[]dto.Question, error)
	QuestionExists(string) (bool, error)
	CountQuestions() (int64, error)
}
