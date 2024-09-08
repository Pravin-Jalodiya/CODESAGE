package services

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils/data_cleaning"
	"cli-project/pkg/utils/readers"
	"cli-project/pkg/validation"
	"errors"
	"fmt"
	"strings"
)

type QuestionService struct {
	questionRepo interfaces.QuestionRepository
}

func NewQuestionService(questionRepo interfaces.QuestionRepository) interfaces.QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

var (
	CSVReader                  = readers.ReadCSV
	ValidateQuestionID         = validation.ValidateQuestionID
	ValidateQuestionDifficulty = validation.ValidateQuestionDifficulty
	ValidateQuestionLink       = validation.ValidateQuestionLink
	ValidateTitleSlug          = validation.ValidateTitleSlug
)

func (s *QuestionService) AddQuestionsFromFile(questionFilePath string) (bool, error) {

	// Read the CSV file
	records, err := CSVReader(questionFilePath)
	if err != nil {
		return false, fmt.Errorf("error reading CSV file: %v", err)
	}

	var questions []models.Question
	newQuestionsAdded := false

	// Loop through the records (skip header row)
	for i, record := range records {
		if i == 0 {
			continue // Skip the header
		}

		// Ensure CSV has the correct number of fields (7 columns expected)
		if len(record) != 7 {
			return false, errors.New("invalid CSV format, expected 7 columns")
		}

		// Clean and validate the fields
		titleSlug := data_cleaning.CleanString(record[0])
		questionID := data_cleaning.CleanString(record[1])
		questionTitle := data_cleaning.CleanString(record[2])
		difficulty := record[3]
		questionLink := record[4]
		topicTags := data_cleaning.CleanTags(record[5])
		companyTags := data_cleaning.CleanTags(record[6])

		// Validate question ID
		valid, err := ValidateQuestionID(questionID)
		if !valid {
			return false, fmt.Errorf("invalid question ID: %v", err)
		}

		// Validate difficulty
		difficulty, err = ValidateQuestionDifficulty(difficulty)
		if err != nil {
			return false, fmt.Errorf("invalid difficulty: %v", err)
		}

		// Validate question link
		questionLink, err = ValidateQuestionLink(questionLink)
		if err != nil {
			return false, fmt.Errorf("invalid question link: %v", err)
		}

		// Build the question struct
		question := models.Question{
			QuestionTitleSlug: titleSlug,
			QuestionID:        questionID,
			QuestionTitle:     questionTitle,
			Difficulty:        difficulty,
			QuestionLink:      questionLink,
			TopicTags:         topicTags,
			CompanyTags:       companyTags,
		}

		// Check if the question already exists
		exists, err := s.QuestionExistsByID(questionID)
		if err != nil {
			return false, fmt.Errorf("error checking if question exists: %v", err)
		}

		// If the question doesn't exist, append to the questions slice
		if !exists {
			questions = append(questions, question)
			newQuestionsAdded = true
		}
	}

	// If new questions were added, insert them into the database
	if newQuestionsAdded {
		err = s.questionRepo.AddQuestions(&questions)
		if err != nil {
			return false, fmt.Errorf("error adding questions to the database: %v", err)
		}
	}

	return newQuestionsAdded, nil
}

func (s *QuestionService) RemoveQuestionByID(questionID string) error {
	// Check if the question exists in the database
	exists, err := s.QuestionExistsByID(questionID)
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
	exists, err := s.QuestionExistsByTitleSlug(questionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("question with title slug %s not found", questionID)
	}

	// Fetch the question from the repository
	question, err := s.questionRepo.FetchQuestionByID(questionID)
	if err != nil {
		return &models.Question{}, err
	}

	return question, nil
}

func (s *QuestionService) GetAllQuestions() (*[]dto.Question, error) {
	return s.questionRepo.FetchAllQuestions()
}

func (s *QuestionService) GetQuestionsByFilters(difficulty, topic, company string) (*[]dto.Question, error) {
	// Validate and clean the difficulty level
	var validDifficulty string
	var err error

	if difficulty != "" && strings.ToLower(difficulty) != "any" {
		validDifficulty, err = ValidateQuestionDifficulty(difficulty)
		if err != nil {
			return nil, err
		}
	}

	// Clean company and topic strings
	cleanCompany := data_cleaning.CleanString(company)
	cleanTopic := data_cleaning.CleanString(topic)

	// Fetch questions by filters from the repository
	return s.questionRepo.FetchQuestionsByFilters(validDifficulty, cleanTopic, cleanCompany)
}

func (s *QuestionService) QuestionExistsByID(questionID string) (bool, error) {
	// Validate the question ID
	valid, err := validation.ValidateQuestionID(questionID)
	if !valid {
		return false, err
	}

	return s.questionRepo.QuestionExistsByID(questionID)
}

func (s *QuestionService) QuestionExistsByTitleSlug(titleSlug string) (bool, error) {
	// Validate the title slug
	valid, err := ValidateTitleSlug(titleSlug)
	if !valid {
		return false, err
	}

	return s.questionRepo.QuestionExistsByTitleSlug(titleSlug)
}

func (s *QuestionService) GetTotalQuestionsCount() (int, error) {
	return s.questionRepo.CountQuestions()
}
