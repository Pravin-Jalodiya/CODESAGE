package interfaces

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"context"
)

type QuestionService interface {
	AddQuestionsFromFile(ctx context.Context, questionFilePath string) (bool, error)
	RemoveQuestionByID(ctx context.Context, questionID string) error
	GetQuestionByID(ctx context.Context, questionID string) (*models.Question, error)
	GetAllQuestions(context.Context) ([]dto.Question, error)
	GetQuestionsByFilters(ctx context.Context, difficulty, company, topic string) ([]dto.Question, error)
	QuestionExistsByID(ctx context.Context, questionID string) (bool, error)
	QuestionExistsByTitleSlug(context.Context, string) (bool, error)
	GetTotalQuestionsCount(context.Context) (int, error)
}
