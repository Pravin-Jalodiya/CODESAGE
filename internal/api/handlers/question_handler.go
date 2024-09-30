package handlers

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"cli-project/pkg/validation"
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// QuestionHandler handles question-related requests.
type QuestionHandler struct {
	questionService interfaces.QuestionService
}

// NewQuestionHandler creates a new QuestionHandler instance.
func NewQuestionHandler(questionService interfaces.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// AddQuestions handles the file upload for questions CSV.
func (q *QuestionHandler) AddQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form, specifying a max memory buffer
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Logger.Errorw("Error parsing form data", "error", err, "time", time.Now())
		errs.JSONError(w, "Error parsing form data: "+err.Error(), errs.CodeInvalidRequest)
		return
	}

	// Retrieve the file from the form
	file, _, err := r.FormFile("questions_file")
	if err != nil {
		logger.Logger.Errorw("Error retrieving the file", "error", err, "time", time.Now())
		errs.JSONError(w, "Error retrieving the file: "+err.Error(), errs.CodeInvalidRequest)
		return
	}
	defer file.Close()

	// Read the CSV content directly from the uploaded file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		logger.Logger.Errorw("Error reading the CSV file", "error", err, "time", time.Now())
		errs.JSONError(w, "Error reading the CSV file: "+err.Error(), errs.CodeInvalidRequest)
		return
	}

	// Call the service method to process the records
	newQuestionsAdded, existingQuestionsUpdated, err := q.questionService.AddQuestionsFromRecords(r.Context(), records)
	if err != nil {
		logger.Logger.Errorw("Error processing the records", "error", err, "time", time.Now())
		errs.JSONError(w, "Error processing the records: "+err.Error(), errs.CodeDbError)
		return
	}

	// Prepare the response message
	var message string
	switch {
	case newQuestionsAdded && existingQuestionsUpdated:
		message = "Questions added and updated successfully"
	case newQuestionsAdded:
		message = "Questions added successfully"
	case existingQuestionsUpdated:
		message = "Questions updated successfully"
	default:
		message = "No new questions found"
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"code":    http.StatusOK,
		"message": message,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Error encoding response", "error", err, "time", time.Now())
		errs.JSONError(w, "Error encoding response: "+err.Error(), errs.CodeUnexpectedError)
	}
}

func (q *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	difficulty := r.URL.Query().Get("difficulty")
	company := r.URL.Query().Get("company")
	topic := r.URL.Query().Get("topic")

	var limit, offset int
	var err error

	// Parse and validate limit if provided
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			errs.NewAppError(errs.CodeInvalidRequest, "Invalid limit: must be a positive number").ToJSON(w)
			logger.Logger.Errorw("Invalid limit value", "method", r.Method, "error", err, "time", time.Now())
			return
		}
	}

	// Parse and validate offset if provided
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			errs.NewAppError(errs.CodeInvalidRequest, "Invalid offset: must be a non-negative number").ToJSON(w)
			logger.Logger.Errorw("Invalid offset value", "method", r.Method, "error", err, "time", time.Now())
			return
		}
	}

	// Validate difficulty level if provided
	if difficulty != "" {
		_, err = validation.ValidateQuestionDifficulty(difficulty)
		if err != nil {
			errs.NewAppError(errs.CodeValidationError, err.Error()).ToJSON(w)
			logger.Logger.Errorw("Invalid difficulty level", "method", r.Method, "error", err, "time", time.Now())
			return
		}
	}

	ctx := r.Context()
	questions, err := q.questionService.GetQuestionsByFilters(ctx, difficulty, topic, company)
	if err != nil {
		errs.JSONError(w, "Error fetching questions: "+err.Error(), errs.CodeDbError)
		logger.Logger.Errorw("Error fetching questions", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Handle slicing for pagination
	totalQuestions := len(questions)
	paginatedQuestions := make([]dto.Question, 0)

	if limitStr != "" { // If limit is provided, do pagination
		if offset < totalQuestions {
			end := offset + limit
			if end > totalQuestions {
				end = totalQuestions
			}
			paginatedQuestions = questions[offset:end]
		}
	} else { // If no limit is provided, return all questions
		paginatedQuestions = questions
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Fetched questions successfully",
		"questions": func() []dto.Question {
			if paginatedQuestions == nil {
				return []dto.Question{}
			}
			return paginatedQuestions
		}(),
		"total": totalQuestions, // Include total count for client-side pagination handling
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		logger.Logger.Errorw("Error encoding response", "method", r.Method, "error", err, "time", time.Now())
		errs.JSONError(w, "Error encoding response: "+err.Error(), errs.CodeUnexpectedError)
	}
	logger.Logger.Infow("Fetched questions successfully", "method", r.Method, "questionsCount", len(paginatedQuestions), "time", time.Now())
}

func (q *QuestionHandler) RemoveQuestionById(w http.ResponseWriter, r *http.Request) {
	questionID := r.URL.Query().Get("id")

	valid, err := validation.ValidateQuestionID(questionID)
	if !valid {
		errs.NewAppError(errs.CodeInvalidRequest, "Invalid question ID").ToJSON(w)
		logger.Logger.Errorw("Invalid question ID", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	err = q.questionService.RemoveQuestionByID(r.Context(), questionID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			errs.NewAppError(errs.CodeInvalidRequest, "Question not found").ToJSON(w)
		} else if errors.Is(err, errs.ErrDatabaseConnection) {
			errs.NewAppError(errs.CodeDbError, "Failed to connect to database").ToJSON(w)
		} else if errors.Is(err, errs.ErrQueryExecution) {
			errs.NewAppError(errs.CodeDbError, "Failed to execute query").ToJSON(w)
		} else {
			errs.NewAppError(errs.CodeUnexpectedError, "Internal server error").ToJSON(w)
		}
		logger.Logger.Errorw("Error removing question", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Question deleted successfully",
	}
	_ = json.NewEncoder(w).Encode(response)
	logger.Logger.Infow("Question deleted successfully", "method", r.Method, "questionID", questionID, "time", time.Now())
}
