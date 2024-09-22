package handlers

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/pkg/errors"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"io"
	"net/http"
	"os"
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

// AddQuestions handles the file upload for questions CSV.
func (q *QuestionHandler) AddQuestions(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form, specifying a max memory buffer
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errs.JSONError(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form
	file, _, err := r.FormFile("questions_file")
	if err != nil {
		errs.JSONError(w, "Error retrieving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a temporary file on the server
	tempFile, err := os.CreateTemp("", "upload-*.csv")
	if err != nil {
		errs.JSONError(w, "Error creating temporary file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read the file content into the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		errs.JSONError(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the file path
	filePath := tempFile.Name()

	// Call the service method to process the file
	added, err := q.questionService.AddQuestionsFromFile(r.Context(), filePath)
	if err != nil {
		os.Remove(filePath) // Ensure the temporary file is deleted
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
	json.NewEncoder(w).Encode(jsonResponse)
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
