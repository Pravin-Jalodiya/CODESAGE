package service_test

import (
	"cli-project/internal/app/services"
	"cli-project/internal/config/roles"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/models"
	"cli-project/pkg/globals"
	"cli-project/pkg/utils"
	mocks "cli-project/tests/mocks/repository"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"testing"
	"time"
)

func TestUserService_Signup(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Successful Signup", func(t *testing.T) {
		mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)

		user := models.StandardUser{
			User: models.User{
				Username:     "TestUser",
				Email:        "testuser@example.com",
				Password:     "securepassword",
				Organisation: "test org",
				Country:      "test country",
			},
			LeetcodeID:      "testLeetcode",
			QuestionsSolved: []string{},
			LastSeen:        time.Now().UTC(),
		}

		err := userService.Signup(&user)
		assert.NoError(t, err)
		assert.Equal(t, strings.ToLower("TestUser"), user.Username)
		assert.Equal(t, strings.ToLower("testuser@example.com"), user.Email)
		assert.Equal(t, utils.CapitalizeWords("test org"), user.Organisation)
		assert.Equal(t, utils.CapitalizeWords("test country"), user.Country)
		assert.NotEmpty(t, user.ID)
		assert.NotEqual(t, "securepassword", user.Password)
		assert.Equal(t, "user", user.Role)
		assert.False(t, user.IsBanned)
		assert.Empty(t, user.QuestionsSolved)
		assert.WithinDuration(t, time.Now().UTC(), user.LastSeen, time.Second)
	})

	t.Run("Password Hash Error", func(t *testing.T) {
		originalHashString := services.HashString
		services.HashString = func(password string) (string, error) {
			return "", errors.New("hash error")
		}
		defer func() { services.HashString = originalHashString }()

		user := models.StandardUser{
			User: models.User{
				Username:     "TestUser",
				Email:        "testuser@example.com",
				Password:     "securepassword",
				Organisation: "test org",
				Country:      "test country",
			},
			LeetcodeID:      "testLeetcode",
			QuestionsSolved: []string{},
			LastSeen:        time.Now().UTC(),
		}

		err := userService.Signup(&user)
		assert.EqualError(t, err, "could not hash password")
	})

	t.Run("User Repository Create Error", func(t *testing.T) {
		mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("create error"))

		user := models.StandardUser{
			User: models.User{
				Username:     "TestUser",
				Email:        "testuser@example.com",
				Password:     "securepassword",
				Organisation: "test org",
				Country:      "test country",
			},
			LeetcodeID:      "testLeetcode",
			QuestionsSolved: []string{},
			LastSeen:        time.Now().UTC(),
		}

		err := userService.Signup(&user)
		assert.EqualError(t, err, "could not register user")
		assert.Error(t, err)
	})
}

func TestUserService_Login(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Successful Login", func(t *testing.T) {
		username := "testuser"
		password := "securepassword"
		lowercaseUsername := strings.ToLower(username)

		user := &models.StandardUser{
			User: models.User{
				Username: lowercaseUsername,
				Password: password,
			},
		}

		mockUserRepo.EXPECT().FetchUserByUsername(lowercaseUsername).Return(user, nil)
		mockVerifyString := services.VerifyString
		services.VerifyString = func(inputPassword, storedPassword string) bool {
			return inputPassword == storedPassword
		}
		defer func() { services.VerifyString = mockVerifyString }()

		err := userService.Login(username, password)
		assert.NoError(t, err)
	})

	t.Run("User Not Found", func(t *testing.T) {
		username := "nonexistent"
		password := "password"
		lowercaseUsername := strings.ToLower(username)

		mockUserRepo.EXPECT().FetchUserByUsername(lowercaseUsername).Return(nil, sql.ErrNoRows)

		err := userService.Login(username, password)
		assert.EqualError(t, err, services.ErrUserNotFound.Error())
	})

	t.Run("Error Fetching User", func(t *testing.T) {
		username := "erroruser"
		password := "password"
		lowercaseUsername := strings.ToLower(username)

		mockUserRepo.EXPECT().FetchUserByUsername(lowercaseUsername).Return(nil, errors.New("db error"))

		err := userService.Login(username, password)
		assert.Contains(t, err.Error(), "error fetching user: db error")
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		username := "testuser"
		password := "wrongpassword"
		lowercaseUsername := strings.ToLower(username)

		user := &models.StandardUser{
			User: models.User{
				Username: lowercaseUsername,
				Password: "securepassword",
			},
		}

		mockUserRepo.EXPECT().FetchUserByUsername(lowercaseUsername).Return(user, nil)
		mockVerifyString := services.VerifyString
		services.VerifyString = func(inputPassword, storedPassword string) bool {
			return inputPassword == storedPassword
		}
		defer func() { services.VerifyString = mockVerifyString }()

		err := userService.Login(username, password)
		assert.EqualError(t, err, services.ErrInvalidCredentials.Error())
		assert.Error(t, err)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Successful FetchAllUsers", func(t *testing.T) {
		users := []models.StandardUser{
			{
				User: models.User{
					Username: "testuser",
					Email:    "testuser@example.com",
				},
			},
		}
		mockUserRepo.EXPECT().FetchAllUsers().Return(&users, nil)

		result, err := userService.GetAllUsers()
		assert.NoError(t, err)
		assert.Equal(t, &users, result)
	})

	t.Run("Error Fetching Users", func(t *testing.T) {
		mockUserRepo.EXPECT().FetchAllUsers().Return(nil, errors.New("fetch error"))

		result, err := userService.GetAllUsers()
		assert.Nil(t, result)
		assert.EqualError(t, err, "fetch error")
	})
}

func TestUserService_ViewDashboard(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("ViewDashboard Placeholder", func(t *testing.T) {
		err := userService.ViewDashboard()
		assert.NoError(t, err)
	})
}

func TestUserService_UpdateUserProgress_InvalidUUID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// Invalid UUID
	globals.ActiveUserID = "invalid-uuid"

	err := userService.UpdateUserProgress()
	assert.Error(t, err)
	assert.Equal(t, "invalid user ID: invalid UUID length: 12", err.Error())
}

//func TestUserService_UpdateUserProgress_LeetcodeAPIError(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID but simulate Leetcode API error
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(nil, errors.New("leetcode API error"))
//
//	err := userService.UpdateUserProgress()
//	assert.Error(t, err)
//	assert.Equal(t, "could not fetch stats from LeetCode API: Leetcode API error", err.Error())
//}
//
//func TestUserService_UpdateUserProgress_NoRecentSubmissions(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID with no recent submissions
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(&models.LeetcodeStats{
//		RecentACSubmissionTitleSlugs: []string{},
//	}, nil)
//
//	err := userService.UpdateUserProgress()
//	assert.NoError(t, err)
//}
//
//func TestUserService_UpdateUserProgress_QuestionExistenceCheckError(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID with recent submissions, but error checking question existence
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().FetchUserStats("leetcode_user").Return(&models.LeetcodeStats{}, nil)
//	mockLeetcodeAPI.EXPECT().FetchRecentSubmissions("leetcode_user", 10).Return([]map[string]string{}, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(&models.LeetcodeStats{
//		RecentACSubmissionTitleSlugs: []string{"slug1", "slug2"},
//	}, nil)
//
//	mockQuestionService.EXPECT().QuestionExistsByTitleSlug("slug1").Return(false, errors.New("question existence check error"))
//
//	err := userService.UpdateUserProgress()
//	assert.Error(t, err)
//	assert.Equal(t, "could not check if question exists: question existence check error", err.Error())
//}
//
//func TestUserService_UpdateUserProgress_NoValidSlugs(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID with recent submissions, but no valid slugs
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(&models.LeetcodeStats{
//		RecentACSubmissionTitleSlugs: []string{"slug1", "slug2"},
//	}, nil)
//
//	mockQuestionService.EXPECT().QuestionExistsByTitleSlug("slug1").Return(false, nil)
//	mockQuestionService.EXPECT().QuestionExistsByTitleSlug("slug2").Return(false, nil)
//
//	err := userService.UpdateUserProgress()
//	assert.NoError(t, err)
//}
//
//func TestUserService_UpdateUserProgress_UpdateError(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID with valid slugs, but error updating progress
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(&models.LeetcodeStats{
//		RecentACSubmissionTitleSlugs: []string{"slug1"},
//	}, nil)
//
//	mockQuestionService.EXPECT().QuestionExistsByTitleSlug("slug1").Return(true, nil)
//	mockUserRepo.EXPECT().UpdateUserProgress(gomock.Any(), []string{"slug1"}).Return(errors.New("update error"))
//
//	err := userService.UpdateUserProgress()
//	assert.Error(t, err)
//	assert.Equal(t, "could not update user progress: update error", err.Error())
//}
//
//func TestUserService_UpdateUserProgress_Success(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Valid UUID with valid slugs, update success
//	validUUID := uuid.New().String()
//	globals.ActiveUserID = validUUID
//
//	mockUser := &models.StandardUser{
//		User: models.User{
//			ID: validUUID,
//		},
//		LeetcodeID: "leetcode_user",
//	}
//
//	mockUserRepo.EXPECT().FetchUserByID(validUUID).Return(mockUser, nil)
//	mockLeetcodeAPI.EXPECT().GetStats(mockUser.LeetcodeID).Return(&models.LeetcodeStats{
//		RecentACSubmissionTitleSlugs: []string{"slug1"},
//	}, nil)
//
//	mockQuestionService.EXPECT().QuestionExistsByTitleSlug("slug1").Return(true, nil)
//	mockUserRepo.EXPECT().UpdateUserProgress(gomock.Any(), []string{"slug1"}).Return(nil)
//
//	err := userService.UpdateUserProgress()
//	assert.NoError(t, err)
//}

func TestUserService_GetUserByUsername_NotFound(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	username := "nonexistentuser"

	// Mock the user repository response
	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, mongo.ErrNoDocuments).Times(1)

	user, err := userService.GetUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestUserService_Logout(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Successful Logout", func(t *testing.T) {
		globals.ActiveUserID = "active-user-id"

		user := &models.StandardUser{
			User: models.User{
				ID: globals.ActiveUserID,
			},
			LastSeen: time.Now().UTC(),
		}

		mockUserRepo.EXPECT().FetchUserByID(globals.ActiveUserID).Return(user, nil)
		mockUserRepo.EXPECT().UpdateUserDetails(user).Return(nil)

		err := userService.Logout()
		assert.NoError(t, err)
		assert.Equal(t, "", globals.ActiveUserID)
		assert.WithinDuration(t, time.Now().UTC(), user.LastSeen, time.Second)
	})

	t.Run("User Not Found", func(t *testing.T) {
		globals.ActiveUserID = "nonexistent-user-id"

		mockUserRepo.EXPECT().FetchUserByID(globals.ActiveUserID).Return(nil, errors.New("user not found"))

		err := userService.Logout()
		assert.EqualError(t, err, "user not found")
	})

	t.Run("Error Updating User Details", func(t *testing.T) {
		globals.ActiveUserID = "active-user-id"

		user := &models.StandardUser{
			User: models.User{
				ID: globals.ActiveUserID,
			},
			LastSeen: time.Now().UTC(),
		}

		mockUserRepo.EXPECT().FetchUserByID(globals.ActiveUserID).Return(user, nil)
		mockUserRepo.EXPECT().UpdateUserDetails(user).Return(errors.New("update error"))

		err := userService.Logout()
		assert.EqualError(t, err, "could not update user details")
	})
}

func TestUserService_GetUserRole(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Empty UserID", func(t *testing.T) {
		role, err := userService.GetUserRole("")
		assert.EqualError(t, err, "userID is empty")
		assert.Equal(t, roles.Role(-1), role)
	})

	t.Run("Fetch User Success", func(t *testing.T) {
		mockUser := &models.StandardUser{StandardUser: models.User{ID: "userid", Role: roles.ADMIN.String()}}
		cleanUserID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(cleanUserID).Return(mockUser, nil)

		role, err := userService.GetUserRole(cleanUserID)
		assert.NoError(t, err)
		assert.Equal(t, roles.ADMIN, role)
	})

	t.Run("Fetch User Error", func(t *testing.T) {
		cleanUserID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(cleanUserID).Return(nil, errors.New("fetch error"))

		role, err := userService.GetUserRole(cleanUserID)
		assert.EqualError(t, err, "fetch error")
		assert.Equal(t, roles.Role(-1), role)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Empty UserID", func(t *testing.T) {
		user, err := userService.GetUserByID("")
		assert.EqualError(t, err, "user ID is empty")
		assert.Nil(t, user)
	})

	t.Run("Fetch User Success", func(t *testing.T) {
		mockUser := &models.StandardUser{StandardUser: models.User{ID: "userid"}}
		cleanUserID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(cleanUserID).Return(mockUser, nil)

		user, err := userService.GetUserByID(" UserID ")
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
	})

	t.Run("Fetch User Error", func(t *testing.T) {
		cleanUserID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(cleanUserID).Return(nil, errors.New("fetch error"))

		user, err := userService.GetUserByID(" UserID ")
		assert.EqualError(t, err, "fetch error")
		assert.Nil(t, user)
	})
}

func TestUserService_CountActiveUserInLast24Hours(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Count Successful", func(t *testing.T) {
		mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(5, nil)

		count, err := userService.CountActiveUserInLast24Hours()
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("Count Error", func(t *testing.T) {
		mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(0, errors.New("count error"))

		count, err := userService.CountActiveUserInLast24Hours()
		assert.EqualError(t, err, "count error")
		assert.Equal(t, 0, count)
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Empty Username", func(t *testing.T) {
		user, err := userService.GetUserByUsername("")
		assert.EqualError(t, err, "username is empty")
		assert.Nil(t, user)
	})

	t.Run("Fetch User Success", func(t *testing.T) {
		mockUser := &models.StandardUser{StandardUser: models.User{Username: "testuser"}}
		cleanUsername := "testuser"
		mockUserRepo.EXPECT().FetchUserByUsername(cleanUsername).Return(mockUser, nil)

		user, err := userService.GetUserByUsername(" TestUser ")
		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
	})

	t.Run("Fetch User Error", func(t *testing.T) {
		cleanUsername := "testuser"
		mockUserRepo.EXPECT().FetchUserByUsername(cleanUsername).Return(nil, errors.New("fetch error"))

		user, err := userService.GetUserByUsername(" TestUser ")
		assert.EqualError(t, err, "fetch error")
		assert.Nil(t, user)
	})
}

func TestUserService_Signup_UserAlreadyExists(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	user := models.StandardUser{
		StandardUser: models.User{
			Username: "existinguser",
			Email:    "existinguser@example.com",
		},
	}

	// Simulate an error due to a user already existing
	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("could not register user")).Times(1)

	err := userService.Signup(&user)
	assert.Error(t, err)
	assert.Equal(t, "could not register user", err.Error())
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockUserRepo, nil, nil)

	userID := "user-id"

	// Simulate an error during user retrieval
	mockUserRepo.EXPECT().FetchUserByID(userID).Return(nil, errors.New("user not found")).Times(1)

	result, err := userService.GetUserByID(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
}

//func TestUserService_GetUserLeetcodeStats_Error(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	userID := "12345"
//
//	// Ensure mockUserRepo.FetchUserByID is called correctly
//	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
//		LeetcodeID: "Leetcode_user",
//	}, nil).Times(1)
//
//	// Simulate an error while fetching stats from Leetcode API
//	mockLeetcodeAPI.EXPECT().GetStats("Leetcode_user").Return(nil, errors.New("leetcode API error")).Times(1)
//
//	// Call the method
//	stats, err := userService.GetUserLeetcodeStats(userID)
//	assert.Error(t, err)
//	assert.Nil(t, stats)
//	assert.Equal(t, "Leetcode API error", err.Error())
//}

func TestUserService_BanUser_Error(t *testing.T) {
	// Set up the mock environment and cleanup
	cleanup := setup(t)
	defer cleanup()

	username := "testuser"
	userID := "user-id"

	// Mock FetchUserByUsername to return a user with a valid, non-admin role
	mockUserRepo.EXPECT().FetchUserByUsername(username).Return(&models.StandardUser{
		StandardUser: models.User{
			ID:   userID,
			Role: roles.USER.String(), // Ensure it's a valid non-admin role
		},
	}, nil).Times(2)

	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
		StandardUser: models.User{
			ID:       userID,
			IsBanned: false,
		},
	}, nil).Times(1)

	// Simulate an error while banning the user
	mockUserRepo.EXPECT().BanUser(userID).Return(errors.New("ban user error")).Times(1)

	// Call BanUser and check the result
	banned, err := userService.BanUser(username)

	// Validate the results
	assert.Error(t, err)
	assert.False(t, banned)
	assert.Equal(t, "ban user error", err.Error())
}

func TestUserService_GetUserProgress(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Fetch Progress Success", func(t *testing.T) {
		userID := "userid"
		progress := []string{"question1", "question2"}
		mockUserRepo.EXPECT().FetchUserProgress(userID).Return(&progress, nil)

		userProgress, err := userService.GetUserProgress(userID)
		assert.NoError(t, err)
		assert.Equal(t, &progress, userProgress)
	})

	t.Run("Fetch Progress Error", func(t *testing.T) {
		userID := "userid"
		mockUserRepo.EXPECT().FetchUserProgress(userID).Return(nil, errors.New("fetch error"))

		userProgress, err := userService.GetUserProgress(userID)
		assert.EqualError(t, err, "fetch error")
		assert.Nil(t, userProgress)
	})
}

func TestUserService_GetUserID(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Fetch User Success", func(t *testing.T) {
		username := "TestUser"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				ID: "userid",
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)

		userID, err := userService.GetUserID(username)
		assert.NoError(t, err)
		assert.Equal(t, "userid", userID)
	})

	t.Run("Fetch User Error", func(t *testing.T) {
		username := "TestUser"
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, errors.New("fetch error"))

		userID, err := userService.GetUserID(username)
		assert.EqualError(t, err, "fetch error")
		assert.Empty(t, userID)
	})
}

func TestUserService_BanUser(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Fetch User Error", func(t *testing.T) {
		username := "TestUser"
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, errors.New("fetch error"))

		success, err := userService.BanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("Admin Ban Not Allowed", func(t *testing.T) {
		username := "TestAdmin"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testadmin",
				Role:     roles.ADMIN.String(), // Ensure it's an admin role
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)

		success, err := userService.BanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "ban operation on admin not allowed")
	})

	t.Run("GetUserID Error", func(t *testing.T) {
		username := "TestUser"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				Role:     roles.USER.String(),
			},
		}

		// Expect the call for FetchUserByUsername and return mockUser
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil).Times(1)

		// Expect the second call and return error
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, errors.New("fetch error")).Times(1)

		success, err := userService.BanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("IsUserBanned Error", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				Role:     roles.USER.String(),
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil) // For GetUserID
		mockUserRepo.EXPECT().FetchUserByID(userID).Return(nil, errors.New("fetch error"))

		success, err := userService.BanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("Already Banned User", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				IsBanned: true,
				Role:     roles.USER.String(),
			},
		}

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil),
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil), // For GetUserID
			mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil),
		)

		success, err := userService.BanUser(username)
		assert.True(t, success) // User is already banned
		assert.NoError(t, err)
	})

	t.Run("Ban User Success", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				Role:     roles.USER.String(),
			},
		}

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil),
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil), // For GetUserID
			mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil),
			mockUserRepo.EXPECT().BanUser(userID).Return(nil),
		)

		success, err := userService.BanUser(username)
		assert.False(t, success)
		assert.NoError(t, err)
	})
}

func TestUserService_UnbanUser(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Fetch User Error", func(t *testing.T) {
		username := "TestUser"
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, errors.New("fetch error"))

		success, err := userService.UnbanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("Admin Unban Not Allowed", func(t *testing.T) {
		username := "TestAdmin"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testadmin",
				Role:     roles.ADMIN.String(), // Ensure it's an admin role
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)

		success, err := userService.UnbanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "unban operation on admin not allowed")
	})

	t.Run("GetUserID Error", func(t *testing.T) {
		username := "TestUser"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				Role:     roles.USER.String(),
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(nil, errors.New("fetch error"))

		success, err := userService.UnbanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("IsUserBanned Error", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				Role:     roles.USER.String(),
			},
		}
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil)
		mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil) // For GetUserID
		mockUserRepo.EXPECT().FetchUserByID(userID).Return(nil, errors.New("fetch error"))

		success, err := userService.UnbanUser(username)
		assert.False(t, success)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("Already Unbanned User", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				IsBanned: false,
				Role:     roles.USER.String(),
			},
		}

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil),
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil), // For GetUserID
			mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil),
		)

		success, err := userService.UnbanUser(username)
		assert.True(t, success) // User is already unbanned
		assert.NoError(t, err)
	})

	t.Run("Unban User Success", func(t *testing.T) {
		username := "TestUser"
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				Username: "testuser",
				ID:       userID,
				Role:     roles.USER.String(),
			},
		}

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil),
			mockUserRepo.EXPECT().FetchUserByUsername(username).Return(mockUser, nil), // For GetUserID
			mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil),
		)

		success, err := userService.UnbanUser(username)
		assert.True(t, success)
		assert.NoError(t, err)
	})
}

func TestUserService_IsUserBanned(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Fetch User Error", func(t *testing.T) {
		userID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(userID).Return(nil, errors.New("fetch error"))

		banned, err := userService.IsUserBanned(userID)
		assert.False(t, banned)
		assert.EqualError(t, err, "fetch error")
	})

	t.Run("Fetch User Success", func(t *testing.T) {
		userID := "userid"
		mockUser := &models.StandardUser{
			StandardUser: models.User{
				ID:       userID,
				IsBanned: true,
			},
		}
		mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil)

		banned, err := userService.IsUserBanned(userID)
		assert.True(t, banned)
		assert.NoError(t, err)
	})
}

func TestUserService_GetUserLeetcodeStats(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("GetUserByID Error", func(t *testing.T) {
		userID := "userid"
		mockUserRepo.EXPECT().FetchUserByID(userID).Return(nil, errors.New("fetch error"))

		stats, err := userService.GetUserLeetcodeStats(userID)
		assert.Nil(t, stats)
		assert.EqualError(t, err, "fetch error")
	})

	//t.Run("Fetch Stats Success", func(t *testing.T) {
	//	userID := "userid"
	//	mockUser := &models.StandardUser{
	//		StandardUser: models.User{
	//			ID: userID,
	//		},
	//		LeetcodeID: "leetcode_user",
	//	}
	//	mockUserRepo.EXPECT().FetchUserByID(userID).Return(mockUser, nil)
	//	stats := &models.LeetcodeStats{RecentACSubmissionTitleSlugs: []string{"slug1", "slug2"}}
	//	mockLeetcodeAPI.EXPECT().GetStats("leetcode_user").Return(stats, nil)
	//
	//	result, err := userService.GetUserLeetcodeStats(userID)
	//	assert.NoError(t, err)
	//	assert.Equal(t, stats, result)
	//})
}

//func TestUserService_GetLeetcodeStats(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	userID := "12345"
//
//	mockUserRepo.EXPECT().FetchUserByID(userID).Return(&models.StandardUser{
//		LeetcodeID: "Leetcode_user",
//	}, nil).Times(1)
//
//	mockLeetcodeAPI.EXPECT().GetStats("Leetcode_user").Return(&models.LeetcodeStats{
//		EasyDoneCount:           10,
//		MediumDoneCount:         20,
//		HardDoneCount:           5,
//		TotalEasyCount:          500,
//		TotalHardCount:          300,
//		TotalMediumCount:        200,
//		TotalQuestionsCount:     1000,
//		TotalQuestionsDoneCount: 35,
//	}, nil).Times(1)
//
//	stats, err := userService.GetUserLeetcodeStats(userID)
//	assert.NoError(t, err)
//	assert.Equal(t, 10, stats.EasyDoneCount)
//	assert.Equal(t, 20, stats.MediumDoneCount)
//	assert.Equal(t, 5, stats.HardDoneCount)
//	assert.Equal(t, 500, stats.TotalEasyCount)
//	assert.Equal(t, 300, stats.TotalHardCount)
//	assert.Equal(t, 200, stats.TotalMediumCount)
//	assert.Equal(t, 35, stats.TotalQuestionsDoneCount)
//	assert.Equal(t, 1000, stats.TotalQuestionsCount)
//}

func TestUserService_GetUserCodesageStats(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		userID := "userid"
		userProgress := &[]string{"question1", "question2", "question3"}
		totalQuestionsCount := 100

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserProgress(userID).Return(userProgress, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(totalQuestionsCount, nil),
			mockQuestionService.EXPECT().GetQuestionByID("question1").Return(&models.Question{
				Difficulty:  "easy",
				TopicTags:   []string{"arrays"},
				CompanyTags: []string{"company1"},
			}, nil),
			mockQuestionService.EXPECT().GetQuestionByID("question2").Return(&models.Question{
				Difficulty:  "hard",
				TopicTags:   []string{"strings"},
				CompanyTags: []string{"company2"},
			}, nil),
			mockQuestionService.EXPECT().GetQuestionByID("question3").Return(&models.Question{
				Difficulty:  "medium",
				TopicTags:   []string{"arrays"},
				CompanyTags: []string{"company3"},
			}, nil),
		)

		stats, err := userService.GetUserCodesageStats(userID)
		require.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, totalQuestionsCount, stats.TotalQuestionsCount)
		assert.Equal(t, 3, stats.TotalQuestionsDoneCount)
		assert.Equal(t, 1, stats.EasyDoneCount)
		assert.Equal(t, 1, stats.MediumDoneCount)
		assert.Equal(t, 1, stats.HardDoneCount)
		assert.Equal(t, 2, stats.TopicWiseStats["arrays"])
		assert.Equal(t, 1, stats.TopicWiseStats["strings"])
		assert.Equal(t, 1, stats.CompanyWiseStats["company1"])
		assert.Equal(t, 1, stats.CompanyWiseStats["company2"])
		assert.Equal(t, 1, stats.CompanyWiseStats["company3"])
	})

	t.Run("Error in GetUserProgress", func(t *testing.T) {
		userID := "userid"
		mockUserRepo.EXPECT().FetchUserProgress(userID).Return(nil, errors.New("error in GetUserProgress"))

		stats, err := userService.GetUserCodesageStats(userID)
		require.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("Error in GetTotalQuestionsCount", func(t *testing.T) {
		userID := "userid"
		userProgress := &([]string{"question1", "question2"})

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserProgress(userID).Return(userProgress, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(0, errors.New("error in GetTotalQuestionsCount")),
		)

		stats, err := userService.GetUserCodesageStats(userID)
		require.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("Error in GetQuestionByID", func(t *testing.T) {
		userID := "userid"
		userProgress := &([]string{"question1", "question2"})
		totalQuestionsCount := 100

		gomock.InOrder(
			mockUserRepo.EXPECT().FetchUserProgress(userID).Return(userProgress, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(totalQuestionsCount, nil),
			mockQuestionService.EXPECT().GetQuestionByID("question1").Return(nil, errors.New("error in GetQuestionByID")),
		)

		stats, err := userService.GetUserCodesageStats(userID)
		require.Error(t, err)
		assert.Nil(t, stats)
	})
}

func TestUserService_GetPlatformStats(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	t.Run("Success", func(t *testing.T) {
		activeUsersInLast24Hours := 50
		totalQuestionsCount := 100
		allQuestions := &[]dto.Question{
			{
				Difficulty:  "easy",
				TopicTags:   []string{"arrays"},
				CompanyTags: []string{"company1"},
			},
			{
				Difficulty:  "hard",
				TopicTags:   []string{"strings"},
				CompanyTags: []string{"company2"},
			},
			{
				Difficulty:  "medium",
				TopicTags:   []string{"strings"},
				CompanyTags: []string{"company3"},
			},
		}

		gomock.InOrder(
			mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(activeUsersInLast24Hours, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(totalQuestionsCount, nil),
			mockQuestionService.EXPECT().GetAllQuestions().Return(allQuestions, nil),
		)

		stats, err := userService.GetPlatformStats()
		require.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, activeUsersInLast24Hours, stats.ActiveUserInLast24Hours)
		assert.Equal(t, totalQuestionsCount, stats.TotalQuestionsCount)
		assert.Equal(t, 1, stats.DifficultyWiseQuestionsCount["easy"])
		assert.Equal(t, 1, stats.DifficultyWiseQuestionsCount["hard"])
		assert.Equal(t, 1, stats.DifficultyWiseQuestionsCount["medium"])
		assert.Equal(t, 1, stats.TopicWiseQuestionsCount["arrays"])
		assert.Equal(t, 2, stats.TopicWiseQuestionsCount["strings"])
		assert.Equal(t, 1, stats.CompanyWiseQuestionsCount["company1"])
		assert.Equal(t, 1, stats.CompanyWiseQuestionsCount["company2"])
		assert.Equal(t, 1, stats.CompanyWiseQuestionsCount["company3"])
	})

	t.Run("Error in CountActiveUserInLast24Hours", func(t *testing.T) {
		mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(0, errors.New("error in CountActiveUserInLast24Hours"))

		stats, err := userService.GetPlatformStats()
		require.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("Error in GetTotalQuestionsCount", func(t *testing.T) {
		activeUsersInLast24Hours := 50

		gomock.InOrder(
			mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(activeUsersInLast24Hours, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(0, errors.New("error in GetTotalQuestionsCount")),
		)

		stats, err := userService.GetPlatformStats()
		require.Error(t, err)
		assert.Nil(t, stats)
	})

	t.Run("Error in GetAllQuestions", func(t *testing.T) {
		activeUsersInLast24Hours := 50
		totalQuestionsCount := 100

		gomock.InOrder(
			mockUserRepo.EXPECT().CountActiveUsersInLast24Hours().Return(activeUsersInLast24Hours, nil),
			mockQuestionService.EXPECT().GetTotalQuestionsCount().Return(totalQuestionsCount, nil),
			mockQuestionService.EXPECT().GetAllQuestions().Return(nil, errors.New("error in GetAllQuestions")),
		)

		stats, err := userService.GetPlatformStats()
		require.Error(t, err)
		assert.Nil(t, stats)
	})
}
