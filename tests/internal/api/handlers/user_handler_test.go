package handlers_test

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	mocks "cli-project/tests/mocks/services"
	"context"
	"encoding/json"
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

		// Create the mux router and register the required route
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

	t.Run("Unauthorized-Different User", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			Role:     roles.USER,
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

	t.Run("Unauthorized-Different User", func(t *testing.T) {
		userID := uuid.New().String()
		userMetaData := middleware.UserMetaData{
			UserId:   uuid.MustParse(userID),
			Username: "anotheruser",
		}
		req := httptest.NewRequest("GET", "/user/testuser", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		userHandler.GetUserByID(w, req)

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

		// Create the mux router and register the required route
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

		req := httptest.NewRequest("PUT", "/user/testuser/progress", nil)
		ctx := context.WithValue(req.Context(), "userMetaData", userMetaData)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		userHandler.UpdateUserProgress(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User progress updated successfully", response["message"])
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

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		userHandler.GetAllUsers(w, req)

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

		req := httptest.NewRequest("GET", "/users?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		userHandler.GetUsers(w, req)

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

		req := httptest.NewRequest("PUT", "/user/ban?username=testuser", nil)
		w := httptest.NewRecorder()

		userHandler.UpdateUserBanState(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "User banned successfully", response["message"])
	})
}
