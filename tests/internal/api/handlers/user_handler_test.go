package handlers_test

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	mocks "cli-project/tests/mocks/services"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			Role:     roles.USER,
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		user := &models.StandardUser{
			StandardUser: models.User{
				Username:     "testuser",
				Name:         "Test User",
				Email:        "testuser@gmail.com",
				Organisation: "TestOrg",
				Country:      "India",
			},
			LeetcodeID: "leetcodeTest",
			LastSeen:   time.Now().UTC(),
		}

		mockUserService.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}", userHandler.GetUserByID).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User profile retrieved successfully", response["message"])
	})

	t.Run("Unauthorized-Metadata Not Found", func(t *testing.T) {
		r := mux.NewRouter()
		r.HandleFunc("/user/{username}", userHandler.GetUserByID).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Could not retrieve user metadata", response["message"])
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			Role:     roles.USER,
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		mockUserService.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errs.ErrUserNotFound).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}", userHandler.GetUserByID).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["message"])
	})

	t.Run("Unauthorized-Different User", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "anotheruser",
		}

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}", userHandler.GetUserByID).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized access: token does not match requested user", response["message"])
	})
}

func TestGetUserProgress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		mockUserService.EXPECT().GetUserLeetcodeStats(userID).Return(&models.LeetcodeStats{}, nil).Times(1)
		mockUserService.EXPECT().GetUserCodesageStats(gomock.Any(), userID).Return(&models.CodesageStats{}, nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.GetUserProgress).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response handlers.UserProgressResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched user progress successfully", response.Message)
	})

	t.Run("Unauthorized-Metadata Not Found", func(t *testing.T) {
		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.GetUserProgress).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser/progress", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Could not retrieve user metadata", response["message"])
	})

	t.Run("Unauthorized-Username Mismatch", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "anotheruser",
		}

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.GetUserProgress).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized access", response["message"])
	})

	t.Run("Error Fetching Leetcode Stats", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		mockUserService.EXPECT().GetUserLeetcodeStats(userID).Return(nil, errors.New("leetcode error")).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.GetUserProgress).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Error fetching Leetcode stats: leetcode error", response["message"])
	})

	t.Run("Error Fetching Codesage Stats", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		mockUserService.EXPECT().GetUserLeetcodeStats(userID).Return(&models.LeetcodeStats{}, nil).Times(1)
		mockUserService.EXPECT().GetUserCodesageStats(gomock.Any(), userID).Return(nil, errors.New("codesage error")).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.GetUserProgress).Methods("GET")

		req := httptest.NewRequest("GET", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Error fetching Codesage stats: codesage error", response["message"])
	})
}

func TestUpdateUserProgress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}
		mockUserService.EXPECT().UpdateUserProgress(gomock.Any(), uuid.MustParse(userID)).Return(nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.UpdateUserProgress).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User progress updated successfully", response["message"])
	})

	t.Run("Unauthorized-Metadata Not Found", func(t *testing.T) {
		r := mux.NewRouter()
		r.HandleFunc("/user/{username}/progress", userHandler.UpdateUserProgress).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/testuser/progress", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Could not retrieve user metadata", response["message"])
	})

	t.Run("Error Updating User Progress", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "testuser",
		}

		// Mock the specific error that should be returned
		mockUserService.EXPECT().UpdateUserProgress(gomock.Any(), uuid.MustParse(userID)).Return(errors.New("Error fetching data from LeetCode API")).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/progress/{username}", userHandler.UpdateUserProgress).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/progress/testuser", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Error fetching data from LeetCode API", response["message"])
	})
}

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		users := []dto.StandardUser{{
			StandardUser: dto.User{
				Username:     "testuser",
				Name:         "Test User",
				Email:        "testuser@example.com",
				Organisation: "TestOrg",
				Country:      "TestCountry",
			},
			LeetcodeID: "leetcodeTest",
			LastSeen:   time.Now(),
		}}
		mockUserService.EXPECT().GetAllUsers(gomock.Any()).Return(users, nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched users successfully", response["message"])
	})
}

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		users := []dto.StandardUser{{
			StandardUser: dto.User{
				Username:     "testuser",
				Name:         "Test User",
				Email:        "testuser@example.com",
				Organisation: "TestOrg",
				Country:      "TestCountry",
			},
			LeetcodeID: "leetcodeTest",
			LastSeen:   time.Now(),
		}}

		mockUserService.EXPECT().GetAllUsers(gomock.Any()).Return(users, nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

		req := httptest.NewRequest("GET", "/users?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fetched users successfully", response["message"])
	})
}

func TestUpdateUserBanState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		mockUserService.EXPECT().UpdateUserBanState(gomock.Any(), "testuser").Return("User banned successfully", nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/ban", userHandler.UpdateUserBanState).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/ban?username=testuser", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User banned successfully", response["message"])
	})

	t.Run("BadRequest-Username Not Provided", func(t *testing.T) {
		r := mux.NewRouter()
		r.HandleFunc("/user/ban", userHandler.UpdateUserBanState).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/ban", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request: 'username' query parameter is required", response["message"])
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserService.EXPECT().UpdateUserBanState(gomock.Any(), "testuser").Return("", errs.ErrUserNotFound).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/ban", userHandler.UpdateUserBanState).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/ban?username=testuser", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["message"])
	})

	t.Run("Operation Not Allowed", func(t *testing.T) {
		mockUserService.EXPECT().UpdateUserBanState(gomock.Any(), "testuser").Return("", errs.ErrInvalidParameterError).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/ban", userHandler.UpdateUserBanState).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/ban?username=testuser", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Operation not allowed", response["message"])
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockUserService.EXPECT().UpdateUserBanState(gomock.Any(), "testuser").Return("", errors.New("internal server error")).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/user/ban", userHandler.UpdateUserBanState).Methods("PUT")

		req := httptest.NewRequest("PUT", "/user/ban?username=testuser", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal Server Error", response["message"])
	})
}

func TestGetPlatformStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := mocks.NewMockUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		platformStats := &models.PlatformStats{
			ActiveUserInLast24Hours:      800,
			TotalQuestionsCount:          1000,
			DifficultyWiseQuestionsCount: map[string]int{"easy": 300, "medium": 500, "hard": 200},
			TopicWiseQuestionsCount:      map[string]int{"arrays": 200, "strings": 150, "trees": 100},
			CompanyWiseQuestionsCount:    map[string]int{"google": 300, "amazon": 400},
		}
		mockUserService.EXPECT().GetPlatformStats(gomock.Any()).Return(platformStats, nil).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/platform-stats", userHandler.GetPlatformStats).Methods("GET")

		req := httptest.NewRequest("GET", "/platform-stats", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		// Extract the stats as a PlatformStats struct from the response
		stats := response["stats"].(map[string]interface{})
		difficultyWise := make(map[string]int)
		for k, v := range stats["DifficultyWiseQuestionsCount"].(map[string]interface{}) {
			difficultyWise[k] = int(v.(float64)) // type assertion may vary depending on JSON deserialization
		}
		topicWise := make(map[string]int)
		for k, v := range stats["TopicWiseQuestionsCount"].(map[string]interface{}) {
			topicWise[k] = int(v.(float64))
		}
		companyWise := make(map[string]int)
		for k, v := range stats["CompanyWiseQuestionsCount"].(map[string]interface{}) {
			companyWise[k] = int(v.(float64))
		}

		actualStats := &models.PlatformStats{
			ActiveUserInLast24Hours:      int(stats["ActiveUserInLast24Hours"].(float64)),
			TotalQuestionsCount:          int(stats["TotalQuestionsCount"].(float64)),
			DifficultyWiseQuestionsCount: difficultyWise,
			TopicWiseQuestionsCount:      topicWise,
			CompanyWiseQuestionsCount:    companyWise,
		}

		assert.Equal(t, platformStats, actualStats)
		assert.Equal(t, "Fetched platform stats successfully", response["message"])
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockUserService.EXPECT().GetPlatformStats(gomock.Any()).Return(nil, errors.New("failed to fetch platform stats")).Times(1)

		r := mux.NewRouter()
		r.HandleFunc("/platform-stats", userHandler.GetPlatformStats).Methods("GET")

		req := httptest.NewRequest("GET", "/platform-stats", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to fetch platform stats", response["message"])
	})
}
