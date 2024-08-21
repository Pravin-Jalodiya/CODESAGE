package interfaces

import "cli-project/internal/domain/models"

type QuestionRepository interface {
	AddQuestionsByID(questionID []int) error
	AddQuestionsFromFile(filePath string) error
	RemoveQuestionsByID(questionID []int) error
	FetchQuestionByID(questionID int) (models.Question, error)
	FetchAllQuestions() ([]models.Question, error)
	FetchQuestionsByFilters(difficulty, company, topic string) ([]models.Question, error)
}
