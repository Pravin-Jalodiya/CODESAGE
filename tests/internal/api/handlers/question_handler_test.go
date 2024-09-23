package handlers_test

import (
	"bytes"
	"cli-project/internal/api/handlers"
	"cli-project/internal/domain/dto"
	errs "cli-project/pkg/errors"
	mocks "cli-project/tests/mocks/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMultipartRequest(uri string, params map[string]string, paramName, path string, fileContents string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(fileContents))
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestGetAllQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQuestionService := mocks.NewMockQuestionService(ctrl)
	questionHandler := handlers.NewQuestionHandler(mockQuestionService)

	t.Run("Success", func(t *testing.T) {
		mockQuestionService.EXPECT().GetAllQuestions(gomock.Any()).Return([]dto.Question{}, nil).Times(1)

		req := httptest.NewRequest("GET", "/question/all", nil)
		w := httptest.NewRecorder()

		questionHandler.GetAllQuestions(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched questions successfully", response["message"])
	})

	t.Run("Bad Request Error", func(t *testing.T) {
		mockQuestionService.EXPECT().GetAllQuestions(gomock.Any()).Return(nil, errors.New("bad request")).Times(1)

		req := httptest.NewRequest("GET", "/question/all", nil)
		w := httptest.NewRecorder()

		questionHandler.GetAllQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Error fetching questions", response["message"])
	})
}

func TestAddQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQuestionService := mocks.NewMockQuestionService(ctrl)
	questionHandler := handlers.NewQuestionHandler(mockQuestionService)

	t.Run("Success", func(t *testing.T) {
		added := true
		mockQuestionService.EXPECT().AddQuestionsFromFile(gomock.Any(), gomock.Any()).Return(added, nil).Times(1)

		params := map[string]string{}
		fileContents := "sample,data,for,csv"
		req, err := createMultipartRequest("/question/add", params, "questions_file", "test.csv", fileContents)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		questionHandler.AddQuestions(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Questions added successfully", response["message"])
	})

	t.Run("Error parsing form data", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/question/add", nil)
		w := httptest.NewRecorder()
		questionHandler.AddQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Error parsing form data")
	})

	t.Run("Error retrieving the file", func(t *testing.T) {
		params := map[string]string{}
		req, _ := createMultipartRequest("/question/add", params, "invalid_param", "test.csv", "sample,data")
		w := httptest.NewRecorder()
		questionHandler.AddQuestions(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Error retrieving the file")
	})

	//t.Run("Error creating temporary file", func(t *testing.T) {
	//	originalCreateTemp := os.CreateTemp
	//	os.CreateTemp = func(dir, pattern string) (*os.File, error) {
	//		return nil, errors.New("temp file error")
	//	}
	//	defer func() { osCreateTemp = originalCreateTemp }()
	//
	//	params := map[string]string{}
	//	req, _ := createMultipartRequest("/question/add", params, "questions_file", "test.csv", "sample,data")
	//	w := httptest.NewRecorder()
	//	questionHandler.AddQuestions(w, req)
	//
	//	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	//	var response map[string]interface{}
	//	err := json.NewDecoder(w.Body).Decode(&response)
	//	assert.NoError(t, err)
	//	assert.Contains(t, response["message"], "Error creating temporary file")
	//})

	t.Run("Error processing the file", func(t *testing.T) {
		mockQuestionService.EXPECT().AddQuestionsFromFile(gomock.Any(), gomock.Any()).Return(false, errors.New("processing error")).Times(1)

		params := map[string]string{}
		req, _ := createMultipartRequest("/question/add", params, "questions_file", "test.csv", "sample,data")
		w := httptest.NewRecorder()
		questionHandler.AddQuestions(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "Error processing the file")
	})

	//t.Run("Error saving the file", func(t *testing.T) {
	//	originalCopy := ioCopy
	//	ioCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
	//		return 0, errors.New("saving error")
	//	}
	//	defer func() { ioCopy = originalCopy }()
	//
	//	params := map[string]string{}
	//	req, _ := createMultipartRequest("/question/add", params, "questions_file", "test.csv", "sample,data")
	//	w := httptest.NewRecorder()
	//	questionHandler.AddQuestions(w, req)
	//
	//	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	//	var response map[string]interface{}
	//	err := json.NewDecoder(w.Body).Decode(&response)
	//	assert.NoError(t, err)
	//	assert.Contains(t, response["message"], "Error saving the file")
	//})
}

func TestGetQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQuestionService := mocks.NewMockQuestionService(ctrl)
	questionHandler := handlers.NewQuestionHandler(mockQuestionService)

	t.Run("Success", func(t *testing.T) {
		mockQuestionService.EXPECT().GetQuestionsByFilters(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]dto.Question{}, nil).Times(1)

		req := httptest.NewRequest("GET", "/question?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched questions successfully", response["message"])
	})

	t.Run("Missing limit parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/question", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "limit is a required query parameter", response["message"])
	})

	t.Run("Invalid limit parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/question?limit=invalid", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit: must be a positive number", response["message"])
	})

	t.Run("Invalid offset parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/question?limit=10&offset=invalid", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid offset: must be a non-negative number", response["message"])
	})

	t.Run("Invalid difficulty parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/question?limit=10&offset=0&difficulty=invalid", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid difficulty level: must be 'easy', 'medium', or 'hard'", response["message"])
	})

	t.Run("Error fetching questions", func(t *testing.T) {
		mockQuestionService.EXPECT().GetQuestionsByFilters(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("%w", errs.ErrFetchingQuestion)).Times(1)

		req := httptest.NewRequest("GET", "/question?limit=100", nil)
		w := httptest.NewRecorder()
		questionHandler.GetQuestions(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Error fetching questions", response["message"])
	})
}

func TestRemoveQuestionById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQuestionService := mocks.NewMockQuestionService(ctrl)
	questionHandler := handlers.NewQuestionHandler(mockQuestionService)

	t.Run("Success", func(t *testing.T) {
		mockQuestionService.EXPECT().RemoveQuestionByID(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		req := httptest.NewRequest("DELETE", "/question?id=123", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Question deleted successfully", response["message"])
	})

	t.Run("Invalid question ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/question?id=invalid", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid question ID", response["message"])
	})

	t.Run("Question not found", func(t *testing.T) {
		mockQuestionService.EXPECT().RemoveQuestionByID(gomock.Any(), gomock.Any()).Return(fmt.Errorf("%w", errs.ErrNoRows)).Times(1)

		req := httptest.NewRequest("DELETE", "/question?id=123", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Question not found", response["message"])
	})

	t.Run("Database connection error", func(t *testing.T) {
		mockQuestionService.EXPECT().RemoveQuestionByID(gomock.Any(), gomock.Any()).Return(errors.New("database connection error")).Times(1)

		req := httptest.NewRequest("DELETE", "/question?id=123", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["message"])
	})

	t.Run("Query execution error", func(t *testing.T) {
		mockQuestionService.EXPECT().RemoveQuestionByID(gomock.Any(), gomock.Any()).Return(errors.New("query execution error")).Times(1)

		req := httptest.NewRequest("DELETE", "/question?id=123", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["message"])
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockQuestionService.EXPECT().RemoveQuestionByID(gomock.Any(), gomock.Any()).Return(errors.New("internal server error")).Times(1)

		req := httptest.NewRequest("DELETE", "/question?id=123", nil)
		w := httptest.NewRecorder()
		questionHandler.RemoveQuestionById(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["message"])
	})
}
