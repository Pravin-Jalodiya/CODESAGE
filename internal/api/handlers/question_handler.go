package handlers

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/pkg/errors"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"net/http"
)

// QuestionHandler handles question-related requests.
type QuestionHandler struct {
	questionService interfaces.QuestionService
	validate        *validator.Validate
}

// NewQuestionHandler creates a new QuestionHandler instance.
func NewQuestionHandler(questionService interfaces.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		validate:        validator.New(),
	}
}

func (q *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")
	topic := r.URL.Query().Get("topic")
	company := r.URL.Query().Get("company")

	ctx := r.Context()

	var questions []dto.Question
	var err error

	if questions, err = q.questionService.GetQuestionsByFilters(ctx, difficulty, topic, company); err != nil {
		errs.NewBadRequestError("Error fetching questions").ToJSON(w)
		return
	}

	// Return JSON response on success
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{
		"code":    http.StatusOK,
		"message": "Fetched questions successfully",
		"questions": func() []dto.Question {
			if questions == nil {
				return []dto.Question{}
			}
			return questions
		}(),
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

// RemoveQuestionById handles the DELETE request to remove a question by its ID.
func (q *QuestionHandler) RemoveQuestionById(w http.ResponseWriter, r *http.Request) {
	// Extract question ID from query parameters
	questionID := r.URL.Query().Get("id")

	if questionID == "" {
		errs.NewBadRequestError("Question ID is required").ToJSON(w)
		return
	}

	// Call service to remove the question
	err := q.questionService.RemoveQuestionByID(r.Context(), questionID)
	if err != nil {
		// Handle different errors with corresponding JSON responses
		if errors.Is(err, errs.ErrNoRows) {
			errs.NewNotFoundError("Question not found").ToJSON(w)
		} else if errors.Is(err, errs.ErrDatabaseConnection) {
			errs.NewInternalServerError("Failed to connect to database").ToJSON(w)
		} else if errors.Is(err, errs.ErrQueryExecution) {
			errs.NewInternalServerError("Failed to execute query").ToJSON(w)
		} else {
			errs.NewInternalServerError("Internal server error").ToJSON(w)
		}
		return
	}

	// If successful, return success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"code":    http.StatusOK,
		"message": "Question deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
