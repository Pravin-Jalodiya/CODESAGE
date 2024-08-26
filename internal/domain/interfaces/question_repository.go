package interfaces

import "cli-project/internal/domain/models"

type QuestionRepository interface {
	AddQuestionsByID(*[]string) error
	AddQuestions(*[]models.Question) error
	RemoveQuestionByID(string) error
	FetchQuestionByID(string) (*models.Question, error)
	FetchAllQuestions() (*[]models.Question, error)
	FetchQuestionsByFilters(string, string, string) (*[]models.Question, error)
	QuestionExists(string) (bool, error)
}
