package interfaces

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"context"
)

type QuestionRepository interface {
	AddQuestions(context.Context, *[]models.Question) error
	RemoveQuestionByID(context.Context, string) error
	FetchQuestionByID(context.Context, string) (*models.Question, error)
	FetchAllQuestions(context.Context) (*[]dto.Question, error)
	FetchQuestionsByFilters(context.Context, string, string, string) ([]dto.Question, error)
	QuestionExistsByID(context.Context, string) (bool, error)
	QuestionExistsByTitleSlug(context.Context, string) (bool, error)
	CountQuestions(context.Context) (int, error)
}
