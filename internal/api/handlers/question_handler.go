package handlers

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"cli-project/pkg/validation"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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
		log.Printf("Error parsing form data: %v", err)
		errs.JSONError(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form
	file, _, err := r.FormFile("questions_file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		errs.JSONError(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a temporary file on the server
	tempFile, err := os.CreateTemp("", "upload-*.csv")
	if err != nil {
		log.Printf("Error creating temporary file: %v", err)
		errs.JSONError(w, "Error creating temporary file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read the file content into the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Printf("Error saving the file: %v", err)
		errs.JSONError(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the file path
	filePath := tempFile.Name()

	// Call the service method to process the file
	added, err := q.questionService.AddQuestionsFromFile(r.Context(), filePath)
	if err != nil {
		os.Remove(filePath) // Ensure the temporary file is deleted
		log.Printf("Error processing the file: %v", err)
		errs.JSONError(w, "Error processing the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the temporary file as it's no longer needed
	os.Remove(filePath)

	// Prepare the response message
	var message string
	if added {
		message = "Questions added successfully"
	} else {
		message = "No new questions found"
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]interface{}{
		"code":    http.StatusOK,
		"message": message,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		errs.JSONError(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (q *QuestionHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	questions, err := q.questionService.GetAllQuestions(ctx)
	if err != nil {
		errs.NewBadRequestError("Error fetching questions").ToJSON(w)
		logger.Logger.Errorw("Error fetching questions", "method", r.Method, "error", err, "time", time.Now())
		return
	}

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
	logger.Logger.Infow("Fetched questions successfully", "method", r.Method, "questionsCount", len(questions), "time", time.Now())
}

func (q *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Validate `limit` parameter since it's required
	if limitStr == "" {
		errs.NewBadRequestError("limit is a required query parameter").ToJSON(w)
		logger.Logger.Errorw("Limit query parameter missing", "method", r.Method, "time", time.Now())
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		errs.NewBadRequestError("Invalid limit: must be a positive number").ToJSON(w)
		logger.Logger.Errorw("Invalid limit value", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	var offset int
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			errs.NewBadRequestError("Invalid offset: must be a non-negative number").ToJSON(w)
			logger.Logger.Errorw("Invalid offset value", "method", r.Method, "error", err, "time", time.Now())
			return
		}
	}

	difficulty := r.URL.Query().Get("difficulty")
	company := r.URL.Query().Get("company")
	topic := r.URL.Query().Get("topic")

	if difficulty != "" {
		difficulty, err = validation.ValidateQuestionDifficulty(difficulty)
		if err != nil {
			errs.NewBadRequestError(err.Error()).ToJSON(w)
			logger.Logger.Errorw("Invalid difficulty level", "method", r.Method, "error", err, "time", time.Now())
			return
		}
	}

	ctx := r.Context()
	questions, err := q.questionService.GetQuestionsByFilters(ctx, difficulty, company, topic)
	if err != nil {
		errs.NewBadRequestError("Error fetching questions").ToJSON(w)
		logger.Logger.Errorw("Error fetching questions", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	// Handle slicing for pagination
	totalQuestions := len(questions)
	paginatedQuestions := make([]dto.Question, 0)

	if offset < totalQuestions {
		end := offset + limit
		if end > totalQuestions {
			end = totalQuestions
		}
		paginatedQuestions = questions[offset:end]
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse := map[string]any{
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
	json.NewEncoder(w).Encode(jsonResponse)
	logger.Logger.Infow("Fetched questions successfully", "method", r.Method, "questionsCount", len(paginatedQuestions), "time", time.Now())
}

func (q *QuestionHandler) RemoveQuestionById(w http.ResponseWriter, r *http.Request) {
	questionID := r.URL.Query().Get("id")

	valid, err := validation.ValidateQuestionID(questionID)
	if !valid {
		errs.NewBadRequestError("Invalid question ID").ToJSON(w)
		logger.Logger.Errorw("Invalid question ID", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	err = q.questionService.RemoveQuestionByID(r.Context(), questionID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			errs.NewNotFoundError("Question not found").ToJSON(w)
		} else if errors.Is(err, errs.ErrDatabaseConnection) {
			errs.NewInternalServerError("Failed to connect to database").ToJSON(w)
		} else if errors.Is(err, errs.ErrQueryExecution) {
			errs.NewInternalServerError("Failed to execute query").ToJSON(w)
		} else {
			errs.NewInternalServerError("Internal server error").ToJSON(w)
		}
		logger.Logger.Errorw("Error removing question", "method", r.Method, "error", err, "time", time.Now())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"code":    http.StatusOK,
		"message": "Question deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
	logger.Logger.Infow("Question deleted successfully", "method", r.Method, "questionID", questionID, "time", time.Now())
}
