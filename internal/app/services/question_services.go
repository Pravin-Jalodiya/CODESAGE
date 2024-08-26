package services

import (
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils/data_cleaning"
	"cli-project/pkg/validation"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type QuestionService struct {
	questionRepo interfaces.QuestionRepository
}

func NewQuestionService(questionRepo interfaces.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

func (s *QuestionService) AddQuestionsFromFile(questionFilePath string) (bool, error) {
	file, err := os.Open(questionFilePath)
	if err != nil {
		return false, errors.New("error opening question file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file")
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	var questions []models.Question
	newQuestionsAdded := false

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) != 6 {
			return false, errors.New("invalid CSV format")
		}

		questionID := data_cleaning.CleanString(record[0])
		valid, err := validation.ValidateQuestionID(questionID)
		if !valid {
			return false, err
		}

		questionTitle := data_cleaning.CleanString(record[1])
		difficulty := data_cleaning.CleanString(record[2])
		difficulty, err = validation.ValidateDifficulty(difficulty)
		if err != nil {
			return false, err
		}

		questionLink := data_cleaning.CleanString(record[3])
		questionLink, err = validation.ValidateQuestionLink(questionLink)
		if err != nil {
			return false, err
		}

		topicTags := data_cleaning.CleanTags(record[4])
		companyTags := data_cleaning.CleanTags(record[5])

		question := models.Question{
			QuestionID:    questionID,
			QuestionTitle: questionTitle,
			Difficulty:    difficulty,
			QuestionLink:  questionLink,
			TopicTags:     topicTags,
			CompanyTags:   companyTags,
		}

		exists, err := s.QuestionExists(questionID)
		if err != nil {
			return false, err
		}

		if !exists {
			questions = append(questions, question)
			newQuestionsAdded = true
		}
	}

	if newQuestionsAdded {
		err = s.questionRepo.AddQuestions(&questions)
		if err != nil {
			return false, err
		}
	}

	return newQuestionsAdded, nil
}

func (s *QuestionService) RemoveQuestionByID(questionID string) error {
	// Check if the question exists in the database
	exists, err := s.QuestionExists(questionID)
	if err != nil {
		return fmt.Errorf("error checking if question exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("question with ID %s not found", questionID)
	}

	// Call repository to remove the question
	return s.questionRepo.RemoveQuestionByID(questionID)
}

func (s *QuestionService) GetQuestionByID(questionID string) (*models.Question, error) {
	// Check if the question exists
	exists, err := s.QuestionExists(questionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("question with ID %s not found", questionID)
	}

	// Fetch the question from the repository
	question, err := s.questionRepo.FetchQuestionByID(questionID)
	if err != nil {
		return &models.Question{}, err
	}

	return question, nil
}

func (s *QuestionService) GetAllQuestions() (*[]models.Question, error) {
	return s.questionRepo.FetchAllQuestions()
}

func (s *QuestionService) GetQuestionsByFilters(difficulty, company, topic string) (*[]models.Question, error) {
	// Validate and clean the difficulty level
	validDifficulty, err := validation.ValidateDifficulty(difficulty)
	if err != nil {
		return nil, err
	}

	// Clean company and topic strings
	cleanCompany := data_cleaning.CleanString(company)
	cleanTopic := data_cleaning.CleanString(topic)

	// Fetch questions by filters from the repository
	return s.questionRepo.FetchQuestionsByFilters(validDifficulty, cleanCompany, cleanTopic)
}

func (s *QuestionService) QuestionExists(questionID string) (bool, error) {
	// Validate the question ID
	valid, err := validation.ValidateQuestionID(questionID)
	if !valid {
		return false, err
	}

	return s.questionRepo.QuestionExists(questionID)
}
