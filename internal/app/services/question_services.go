package services

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

type QuestionService struct {
	questionRepo interfaces.QuestionRepository
}

func NewQuestionService(questionRepo interfaces.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

func (s *QuestionService) AddQuestionsFromFile(questionFilePath string) error {
	// Open the CSV file
	file, err := os.Open(questionFilePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var questions []models.Question

	for _, record := range records {
		// Assuming the CSV has the format: ID, Title, Difficulty, Link, TopicTags, CompanyTags
		if len(record) != 6 {
			return errors.New("invalid CSV format")
		}

		questionID := record[0]
		questionTitle := record[1]
		difficulty := record[2]
		questionLink := record[3]
		topicTags := strings.Split(record[4], ",")
		companyTags := strings.Split(record[5], ",")

		// Create a Question object
		question := models.Question{
			QuestionID:    questionID,
			QuestionTitle: questionTitle,
			Difficulty:    difficulty,
			QuestionLink:  questionLink,
			TopicTags:     topicTags,
			CompanyTags:   companyTags,
		}

		questions = append(questions, question)
	}

	// Pass the list of questions to the repository layer
	err = s.questionRepo.AddQuestions(questions)
	if err != nil {
		return err
	}

	return nil
}

func (s *QuestionService) RemoveQuestionByID(questionID string) error {
	return s.questionRepo.RemoveQuestionByID(questionID)
}

func (s *QuestionService) GetQuestionByID(questionID string) (models.Question, error) {
	return s.questionRepo.FetchQuestionByID(questionID)
}

func (s *QuestionService) GetAllQuestions() ([]models.Question, error) {
	return s.questionRepo.FetchAllQuestions()
}

func (s *QuestionService) GetQuestionsByFilters(difficulty, company, topic string) ([]models.Question, error) {
	return s.questionRepo.FetchQuestionsByFilters(difficulty, company, topic)
}
