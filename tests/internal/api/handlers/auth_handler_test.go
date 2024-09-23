package handlers_test

import (
	"bytes"
	"cli-project/internal/api/handlers"
	"cli-project/internal/domain/models"
	errs "cli-project/pkg/errors"
	mocks "cli-project/tests/mocks/services"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSignupHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)

	handler := handlers.NewAuthHandler(mockAuthService)

	t.Run("Successful Signup", func(t *testing.T) {
		user := &models.StandardUser{
			StandardUser: models.User{
				Username:     "Dummy",
				Password:     "Dummy@123",
				Name:         "Dummy User",
				Email:        "dummy@gmail.com",
				Organisation: "Dummy Org",
				Country:      "India",
			},
			LeetcodeID:      "rabbit1",
			QuestionsSolved: []string{},
			LastSeen:        time.Time{},
		}

		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, u *models.StandardUser) error {
				assert.Equal(t, user.StandardUser.Username, u.StandardUser.Username)
				assert.Equal(t, user.StandardUser.Password, u.StandardUser.Password)
				assert.Equal(t, user.StandardUser.Name, u.StandardUser.Name)
				assert.Equal(t, user.StandardUser.Email, u.StandardUser.Email)
				assert.Equal(t, user.StandardUser.Organisation, u.StandardUser.Organisation)
				assert.Equal(t, user.StandardUser.Country, u.StandardUser.Country)
				assert.Equal(t, user.LeetcodeID, u.LeetcodeID)
				return nil
			},
		)

		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`

		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		payload := `{"invalid_json":`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		payload := `{
		    "standard_user": {
		        "username": "",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// New test case for invalid password
	t.Run("Invalid Password", func(t *testing.T) {
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "short",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// New test case for invalid email
	t.Run("Invalid Email", func(t *testing.T) {
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@invalid.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// New test case for already existing username
	t.Run("Username Already Exists", func(t *testing.T) {
		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(fmt.Errorf("%w", errs.ErrUserNameAlreadyExists))
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	// New test case for already existing email
	t.Run("Email Already Exists", func(t *testing.T) {
		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(fmt.Errorf("%w", errs.ErrEmailAlreadyExists))
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	// New test case for already existing LeetcodeID
	t.Run("LeetcodeID Already Exists", func(t *testing.T) {
		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(fmt.Errorf("%w", errs.ErrLeetcodeIDAlreadyExists))
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	// New test case for invalid LeetcodeID
	t.Run("Invalid LeetcodeID", func(t *testing.T) {
		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(fmt.Errorf("%w", errs.ErrLeetcodeUsernameInvalid))
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "InvalidLeetcodeID"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// New test case for internal server error
	t.Run("Internal Server Error", func(t *testing.T) {
		mockAuthService.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))
		payload := `{
		    "standard_user": {
		        "username": "Dummy",
		        "password": "Dummy@123",
		        "name": "Dummy User",
		        "email": "dummy@gmail.com",
		        "organisation": "Dummy Org",
		        "country": "India"
		    },
		    "leetcode_id": "rabbit1"
		}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.SignupHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)

	handler := handlers.NewAuthHandler(mockAuthService)

	t.Run("Successful Login", func(t *testing.T) {
		user := &models.StandardUser{
			StandardUser: models.User{
				ID:       "user1",
				Username: "testuser",
				Role:     "user",
			},
		}

		mockAuthService.EXPECT().Login(gomock.Any(), "testuser", "Password@123").Return(user, nil)

		payload := `{"username": "testuser", "password": "Password@123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		payload := `{"invalid_json":`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Invalid Username or Password", func(t *testing.T) {
		payload := `{"username": "", "password": "Password@123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockAuthService.EXPECT().Login(gomock.Any(), "pravin", "Pravin@12").Return(nil, fmt.Errorf("%w", errs.ErrInvalidPassword))

		payload := `{"username": "pravin", "password": "Pravin@12"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// New test case for user not found
	t.Run("User Not Found", func(t *testing.T) {
		mockAuthService.EXPECT().Login(gomock.Any(), "unknownuser", "Unknown@123").Return(nil, fmt.Errorf("%w", errs.ErrUserNotFound))

		payload := `{"username": "unknownuser", "password": "Unknown@123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	// New test case for internal server error during login
	t.Run("Internal Server Error", func(t *testing.T) {
		mockAuthService.EXPECT().Login(gomock.Any(), "testuser", "Password@123").Return(nil, errors.New("internal error"))

		payload := `{"username": "testuser", "password": "Password@123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.LoginHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestLogoutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)

	handler := handlers.NewAuthHandler(mockAuthService)

	t.Run("Successful Logout", func(t *testing.T) {
		mockAuthService.EXPECT().Logout(gomock.Any()).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rr := httptest.NewRecorder()

		handler.LogoutHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Logout Failure", func(t *testing.T) {
		mockAuthService.EXPECT().Logout(gomock.Any()).Return(errors.New("logout failure"))

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rr := httptest.NewRecorder()

		handler.LogoutHandler(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
