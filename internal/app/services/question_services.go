package services

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/utils"
	"cli-project/pkg/validation"
	"context"
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
	CSVReader                  = utils.ReadCSV
	ValidateQuestionID         = validation.ValidateQuestionID
	ValidateQuestionDifficulty = validation.ValidateQuestionDifficulty
	ValidateQuestionLink       = validation.ValidateQuestionLink
	ValidateTitleSlug          = validation.ValidateTitleSlug
)

func (s *QuestionService) AddQuestionsFromFile(ctx context.Context, questionFilePath string) (bool, error) {
	records, err := CSVReader(questionFilePath)
	if err != nil {
		return false, fmt.Errorf("error reading CSV file: %v", err)
	}

	var questions []models.Question
	newQuestionsAdded := false

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) != 7 {
			return false, errors.New("invalid CSV format, expected 7 columns")
		}

		titleSlug := utils.CleanString(record[0])
		questionID := utils.CleanString(record[1])
		questionTitle := utils.CleanString(record[2])
		difficulty := record[3]
		questionLink := record[4]
		topicTags := utils.CleanTags(record[5])
		companyTags := utils.CleanTags(record[6])

		valid, err := ValidateQuestionID(questionID)
		if !valid {
			return false, fmt.Errorf("invalid question ID: %v", err)
		}

		difficulty, err = ValidateQuestionDifficulty(difficulty)
		if err != nil {
			return false, fmt.Errorf("invalid difficulty: %v", err)
		}

		questionLink, err = ValidateQuestionLink(questionLink)
		if err != nil {
			return false, fmt.Errorf("invalid question link: %v", err)
		}

		question := models.Question{
			QuestionTitleSlug: titleSlug,
			QuestionID:        questionID,
			QuestionTitle:     questionTitle,
			Difficulty:        difficulty,
			QuestionLink:      questionLink,
			TopicTags:         topicTags,
			CompanyTags:       companyTags,
		}

		exists, err := s.QuestionExistsByID(ctx, questionID)
		if err != nil {
			return false, fmt.Errorf("error checking if question exists: %v", err)
		}

		if !exists {
			questions = append(questions, question)
			newQuestionsAdded = true
		}
	}

	if newQuestionsAdded {
		err = s.questionRepo.AddQuestions(ctx, &questions)
		if err != nil {
			return false, fmt.Errorf("error adding questions to the database: %v", err)
		}
	}

	return newQuestionsAdded, nil
}

func (s *QuestionService) RemoveQuestionByID(ctx context.Context, questionID string) error {
	exists, err := s.QuestionExistsByID(ctx, questionID)
	if err != nil {
		return fmt.Errorf("error checking if question exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("question with ID %s not found", questionID)
	}

	return s.questionRepo.RemoveQuestionByID(ctx, questionID)
}

func (s *QuestionService) GetQuestionByID(ctx context.Context, questionID string) (*models.Question, error) {
	exists, err := s.QuestionExistsByTitleSlug(ctx, questionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("question with title slug %s not found", questionID)
	}

	return s.questionRepo.FetchQuestionByID(ctx, questionID)
}

func (s *QuestionService) GetAllQuestions(ctx context.Context) (*[]dto.Question, error) {
	return s.questionRepo.FetchAllQuestions(ctx)
}

func (s *QuestionService) GetQuestionsByFilters(ctx context.Context, difficulty, topic, company string) (*[]dto.Question, error) {
	var validDifficulty string
	var err error

	if difficulty != "" && strings.ToLower(difficulty) != "any" {
		validDifficulty, err = ValidateQuestionDifficulty(difficulty)
		if err != nil {
			return nil, err
		}
	}

	cleanCompany := utils.CleanString(company)
	cleanTopic := utils.CleanString(topic)

	return s.questionRepo.FetchQuestionsByFilters(ctx, validDifficulty, cleanTopic, cleanCompany)
}

func (s *QuestionService) QuestionExistsByID(ctx context.Context, questionID string) (bool, error) {
	valid, err := validation.ValidateQuestionID(questionID)
	if !valid {
		return false, err
	}

	return s.questionRepo.QuestionExistsByID(ctx, questionID)
}

func (s *QuestionService) QuestionExistsByTitleSlug(ctx context.Context, titleSlug string) (bool, error) {
	valid, err := ValidateTitleSlug(titleSlug)
	if !valid {
		return false, err
	}

	return s.questionRepo.QuestionExistsByTitleSlug(ctx, titleSlug)
}

func (s *QuestionService) GetTotalQuestionsCount(ctx context.Context) (int, error) {
	return s.questionRepo.CountQuestions(ctx)
}
