package handlers_test

import (
	"bytes"
	"cli-project/internal/api/handlers"
	"cli-project/internal/domain/dto"
	mocks "cli-project/tests/mocks/services"
	"encoding/json"
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
		// Ensure the mock expectations are set properly
		mockQuestionService.EXPECT().GetAllQuestions(gomock.Any()).Return([]dto.Question{}, nil).Times(1)

		// Create a new HTTP request and recorder
		req := httptest.NewRequest("GET", "/question/all", nil)
		w := httptest.NewRecorder()

		// Call the handler method
		questionHandler.GetAllQuestions(w, req)

		// Validate the response
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched questions successfully", response["message"])
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

		// Create a multipart request with file contents
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
	})
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
}
