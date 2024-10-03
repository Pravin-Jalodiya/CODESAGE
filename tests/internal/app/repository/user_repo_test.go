package repositories_test

import (
	db2 "cli-project/internal/db"
	"cli-project/internal/domain/models"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserRepo_CreateUser(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New()
	user := &models.StandardUser{
		User: models.User{
			ID:           userID.String(),
			Username:     "testuser",
			Password:     "password",
			Name:         "Test User",
			Email:        "test@example.com",
			Role:         "user",
			Organisation: "Test Org",
			Country:      "Test Country",
			IsBanned:     false,
		},
		LeetcodeID: "leetcode123",
		LastSeen:   time.Now().UTC(),
	}

	}

	// Case 1: Successfully creating a user
	mock.ExpectExec(`INSERT INTO Users \(\s*id,\s*username,\s*password,\s*name,\s*email,\s*role,\s*last_seen,\s*organisation,\s*country,\s*leetcode_id,\s*is_banned\s*\) VALUES \(\s*\$1,\s*\$2,\s*\$3,\s*\$4,\s*\$5,\s*\$6,\s*\$7,\s*\$8,\s*\$9,\s*\$10,\s*\$11\s*\)`).
		WithArgs(
			user.ID,
			user.Username,
			user.Password,
			user.Name,
			user.Email,
			user.Role,
			user.LastSeen,
			user.Organisation,
			user.Country,
			user.LeetcodeID,
			user.IsBanned,
		).WillReturnResult(sqlmock.NewResult(0, 1)) // LastInsertId should be 0

	err := userRepo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Failure case due to unique constraint violation
	mock.ExpectExec(`INSERT INTO Users \(\s*id,\s*username,\s*password,\s*name,\s*email,\s*role,\s*last_seen,\s*organisation,\s*country,\s*leetcode_id,\s*is_banned\s*\) VALUES \(\s*\$1,\s*\$2,\s*\$3,\s*\$4,\s*\$5,\s*\$6,\s*\$7,\s*\$8,\s*\$9,\s*\$10,\s*\$11\s*\)`).
		WithArgs(
			user.ID,
			user.Username,
			user.Password,
			user.Name,
			user.Email,
			user.Role,
			user.LastSeen,
			user.Organisation,
			user.Country,
			user.LeetcodeID,
			user.IsBanned,
		).WillReturnError(fmt.Errorf("ERROR: duplicate key value violates unique constraint \"users_new_username_key\" (SQLSTATE 23505)"))

	err = userRepo.CreateUser(user)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	err = userRepo.CreateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_UpdateUserProgress(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New()
	newSlugs := []string{"slug1", "slug2"}

	// Case 1: Inserting new progress when no existing progress
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO users_progress \\(user_id, title_slugs\\) VALUES \\(\\$1, \\$2\\)").
		WithArgs(userID, pq.Array(newSlugs)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := userRepo.UpdateUserProgress(userID, newSlugs)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Updating existing progress when there are new slugs
	mock.ExpectBegin()
	currentSlugs := pq.StringArray{"slug1"}
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"title_slugs"}).AddRow(currentSlugs))
	mock.ExpectExec("UPDATE users_progress SET title_slugs = array\\(SELECT DISTINCT unnest\\(title_slugs\\) \\|\\| unnest\\(\\$1::text\\[\\]\\)\\) WHERE user_id = \\$2").
		WithArgs(pq.Array([]string{"slug2"}), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = userRepo.UpdateUserProgress(userID, newSlugs)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: Error fetching user's progress
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db error"))
	mock.ExpectRollback()

	err = userRepo.UpdateUserProgress(userID, newSlugs)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: No new slugs to add
	mock.ExpectBegin()
	currentSlugs = pq.StringArray{"slug1", "slug2"}
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"title_slugs"}).AddRow(currentSlugs))
	mock.ExpectCommit()

	err = userRepo.UpdateUserProgress(userID, newSlugs)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 5: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	err = userRepo.UpdateUserProgress(userID, newSlugs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_FetchAllUsers(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	expectedUsers := []models.StandardUser{
		{
			User: models.User{
				ID:           uuid.New().String(),
				Username:     "user1",
				Password:     "pass1",
				Name:         "User One",
				Email:        "user1@example.com",
				Role:         "user",
				Organisation: "Org1",
				Country:      "Country1",
				IsBanned:     false,
			},
			LeetcodeID: "lc1",
			LastSeen:   time.Now().UTC(),
		},
		{
			User: models.User{
				ID:           uuid.New().String(),
				Username:     "user2",
				Password:     "pass2",
				Name:         "User Two",
				Email:        "user2@example.com",
				Role:         "user",
				Organisation: "Org2",
				Country:      "Country2",
				IsBanned:     false,
			},
			LeetcodeID: "lc2",
			LastSeen:   time.Now().UTC(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "password", "name", "email", "role", "last_seen", "organisation", "country", "leetcode_id", "is_banned",
	}).AddRow(
		expectedUsers[0].StandardUser.ID, expectedUsers[0].StandardUser.Username, expectedUsers[0].StandardUser.Password, expectedUsers[0].StandardUser.Name, expectedUsers[0].StandardUser.Email,
		expectedUsers[0].StandardUser.Role, expectedUsers[0].LastSeen, expectedUsers[0].StandardUser.Organisation, expectedUsers[0].StandardUser.Country, expectedUsers[0].LeetcodeID, expectedUsers[0].StandardUser.IsBanned,
	).AddRow(
		expectedUsers[1].StandardUser.ID, expectedUsers[1].StandardUser.Username, expectedUsers[1].StandardUser.Password, expectedUsers[1].StandardUser.Name, expectedUsers[1].StandardUser.Email,
		expectedUsers[1].StandardUser.Role, expectedUsers[1].LastSeen, expectedUsers[1].StandardUser.Organisation, expectedUsers[1].StandardUser.Country, expectedUsers[1].LeetcodeID, expectedUsers[1].StandardUser.IsBanned,
	)

	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users").
		WillReturnRows(rows)

	users, err := userRepo.FetchAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, &expectedUsers, users)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	users, err = userRepo.FetchAllUsers()
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_FetchUserByID(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New().String()
	expectedUser := &models.StandardUser{
		StandardUser: models.User{
			ID:           userID,
			Username:     "testuser",
			Password:     "password",
			Name:         "Test User",
			Email:        "test@example.com",
			Role:         "user",
			Organisation: "Test Org",
			Country:      "Test Country",
			IsBanned:     false,
		},
		LeetcodeID: "leetcode123",
		LastSeen:   time.Now().UTC(),
	}

	// Case 1: Successfully fetching a user by ID
	rows := sqlmock.NewRows([]string{
		"id", "username", "password", "name", "email", "role", "last_seen", "organisation", "country", "leetcode_id", "is_banned",
	}).AddRow(
		expectedUser.StandardUser.ID, expectedUser.StandardUser.Username, expectedUser.StandardUser.Password, expectedUser.StandardUser.Name, expectedUser.StandardUser.Email,
		expectedUser.StandardUser.Role, expectedUser.LastSeen, expectedUser.StandardUser.Organisation, expectedUser.StandardUser.Country, expectedUser.LeetcodeID, expectedUser.StandardUser.IsBanned,
	)

	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := userRepo.FetchUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: No user found (sql.ErrNoRows)
	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err = userRepo.FetchUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: General SQL error
	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some SQL error"))

	user, err = userRepo.FetchUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "could not fetch user")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	user, err = userRepo.FetchUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_FetchUserByUsername(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	username := "testuser"
	expectedUser := &models.StandardUser{
		StandardUser: models.User{
			ID:           uuid.New().String(),
			Username:     username,
			Password:     "password",
			Name:         "Test User",
			Email:        "test@example.com",
			Role:         "user",
			Organisation: "Test Org",
			Country:      "Test Country",
			IsBanned:     false,
		},
		LeetcodeID: "leetcode123",
		LastSeen:   time.Now().UTC(),
	}

	// Case 1: Successfully fetching a user by username
	rows := sqlmock.NewRows([]string{
		"id", "username", "password", "name", "email", "role", "last_seen", "organisation", "country", "leetcode_id", "is_banned",
	}).AddRow(
		expectedUser.StandardUser.ID, expectedUser.StandardUser.Username, expectedUser.StandardUser.Password, expectedUser.StandardUser.Name, expectedUser.StandardUser.Email,
		expectedUser.StandardUser.Role, expectedUser.LastSeen, expectedUser.StandardUser.Organisation, expectedUser.StandardUser.Country, expectedUser.LeetcodeID, expectedUser.StandardUser.IsBanned,
	)

	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := userRepo.FetchUserByUsername(username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: No user found (sql.ErrNoRows)
	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	user, err = userRepo.FetchUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: General SQL error
	mock.ExpectQuery("SELECT id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned FROM Users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(fmt.Errorf("some SQL error"))

	user, err = userRepo.FetchUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "could not fetch user")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	user, err = userRepo.FetchUserByUsername(username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_FetchUserProgress(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New().String()
	expectedProgress := []string{"slug1", "slug2"}

	// Case 1: Successfully fetching user progress
	rows := sqlmock.NewRows([]string{"title_slugs"}).AddRow(pq.StringArray(expectedProgress))

	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	progress, err := userRepo.FetchUserProgress(userID)
	assert.NoError(t, err)
	assert.Equal(t, &expectedProgress, progress)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: No user progress found (sql.ErrNoRows)
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	progress, err = userRepo.FetchUserProgress(userID)
	assert.Error(t, err)
	assert.Nil(t, progress)
	assert.Equal(t, "user progress not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: General SQL error
	mock.ExpectQuery("SELECT title_slugs FROM users_progress WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some SQL error"))

	progress, err = userRepo.FetchUserProgress(userID)
	assert.Error(t, err)
	assert.Nil(t, progress)
	assert.Contains(t, err.Error(), "could not fetch user progress")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	progress, err = userRepo.FetchUserProgress(userID)
	assert.Error(t, err)
	assert.Nil(t, progress)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_UpdateUserDetails(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	user := &models.StandardUser{
		StandardUser: models.User{
			ID:           uuid.New().String(),
			Username:     "testuser",
			Password:     "password",
			Name:         "Test User",
			Email:        "test@example.com",
			Role:         "user",
			Organisation: "Test Org",
			Country:      "Test Country",
			IsBanned:     false,
		},
		LeetcodeID: "leetcode123",
		LastSeen:   time.Now().UTC(),
	}

	// Case 1: Successfully updating user details
	mock.ExpectExec("UPDATE Users SET username = \\$1, email = \\$2, password = \\$3, name = \\$4, organisation = \\$5, country = \\$6, leetcode_id = \\$7, last_seen = \\$8 WHERE id = \\$9").
		WithArgs(
			user.StandardUser.Username, user.StandardUser.Email, user.StandardUser.Password, user.StandardUser.Name,
			user.StandardUser.Organisation, user.StandardUser.Country, user.LeetcodeID, user.LastSeen, user.StandardUser.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.UpdateUserDetails(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Missing user ID
	user.StandardUser.ID = ""
	err = userRepo.UpdateUserDetails(user)
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())

	// Case 3: General SQL error
	user.StandardUser.ID = uuid.New().String()
	mock.ExpectExec("UPDATE Users SET username = \\$1, email = \\$2, password = \\$3, name = \\$4, organisation = \\$5, country = \\$6, leetcode_id = \\$7, last_seen = \\$8 WHERE id = \\$9").
		WithArgs(
			user.StandardUser.Username, user.StandardUser.Email, user.StandardUser.Password, user.StandardUser.Name,
			user.StandardUser.Organisation, user.StandardUser.Country, user.LeetcodeID, user.LastSeen, user.StandardUser.ID,
		).WillReturnError(fmt.Errorf("some SQL error"))

	err = userRepo.UpdateUserDetails(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not update user details")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	err = userRepo.UpdateUserDetails(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_BanUser(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New().String()

	// Case 1: Successfully banning a user
	mock.ExpectExec("UPDATE Users SET is_banned = TRUE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.BanUser(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: User ID is empty
	err = userRepo.BanUser("")
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())

	// Case 3: SQL error when banning a user
	mock.ExpectExec("UPDATE Users SET is_banned = TRUE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some SQL error"))

	err = userRepo.BanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not ban user")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: No rows affected (user not found)
	mock.ExpectExec("UPDATE Users SET is_banned = TRUE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = userRepo.BanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("user with ID %s not found", userID))
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 5: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	err = userRepo.BanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_UnbanUser(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	userID := uuid.New().String()

	// Case 1: Successfully unbanning a user
	mock.ExpectExec("UPDATE Users SET is_banned = FALSE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := userRepo.UnbanUser(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: User ID is empty
	err = userRepo.UnbanUser("")
	assert.Error(t, err)
	assert.Equal(t, "user ID is required", err.Error())

	// Case 3: SQL error when unbanning a user
	mock.ExpectExec("UPDATE Users SET is_banned = FALSE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("some SQL error"))

	err = userRepo.UnbanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not unban user")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: No rows affected (user not found)
	mock.ExpectExec("UPDATE Users SET is_banned = FALSE WHERE id = \\$1 and role = 'user'").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = userRepo.UnbanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("user with ID %s not found", userID))
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 5: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	err = userRepo.UnbanUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_CountActiveUsersInLast24Hours(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	// Define the time range with consistent precision
	now := time.Now().UTC().Truncate(time.Millisecond) // Truncate to milliseconds for precision match
	twentyFourHoursAgo := now.Add(-24 * time.Hour)
	expectedCount := 5

	// Case 1: Successfully counting active users
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE last_seen >= \\$1").
		WithArgs(twentyFourHoursAgo).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	count, err := userRepo.CountActiveUsersInLast24Hours()
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: SQL error when counting active users
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE last_seen >= \\$1").
		WithArgs(twentyFourHoursAgo).
		WillReturnError(fmt.Errorf("some SQL error"))

	count, err = userRepo.CountActiveUsersInLast24Hours()
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "could not count active users")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	count, err = userRepo.CountActiveUsersInLast24Hours()
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_IsEmailUnique(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	email := "test@example.com"

	// Case 1: Successfully checking email uniqueness when email is unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isUnique, err := userRepo.IsEmailUnique(email)
	assert.NoError(t, err)
	assert.True(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Successfully checking email uniqueness when email is not unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isUnique, err = userRepo.IsEmailUnique(email)
	assert.NoError(t, err)
	assert.False(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: SQL error when checking email uniqueness
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(fmt.Errorf("some SQL error"))

	isUnique, err = userRepo.IsEmailUnique(email)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "could not check email uniqueness")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	isUnique, err = userRepo.IsEmailUnique(email)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_IsUsernameUnique(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	username := "testuser"

	// Case 1: Username is unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isUnique, err := userRepo.IsUsernameUnique(username)
	assert.NoError(t, err)
	assert.True(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: Username is not unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isUnique, err = userRepo.IsUsernameUnique(username)
	assert.NoError(t, err)
	assert.False(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: SQL error when checking username uniqueness
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(fmt.Errorf("some SQL error"))

	isUnique, err = userRepo.IsUsernameUnique(username)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "could not check username uniqueness")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	isUnique, err = userRepo.IsUsernameUnique(username)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}

func TestUserRepo_IsLeetcodeIDUnique(t *testing.T) {
	cleanup := setup(t)
	defer cleanup()

	leetcodeID := "leetcode123"

	// Case 1: LeetcodeID is unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE leetcode_id = \\$1").
		WithArgs(leetcodeID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isUnique, err := userRepo.IsLeetcodeIDUnique(leetcodeID)
	assert.NoError(t, err)
	assert.True(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 2: LeetcodeID is not unique
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE leetcode_id = \\$1").
		WithArgs(leetcodeID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	isUnique, err = userRepo.IsLeetcodeIDUnique(leetcodeID)
	assert.NoError(t, err)
	assert.False(t, isUnique)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 3: SQL error when checking LeetcodeID uniqueness
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Users WHERE leetcode_id = \\$1").
		WithArgs(leetcodeID).
		WillReturnError(fmt.Errorf("some SQL error"))

	isUnique, err = userRepo.IsLeetcodeIDUnique(leetcodeID)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "could not check LeetcodeID uniqueness")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Case 4: Error when failing to get DB connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return nil, fmt.Errorf("failed to get DB connection")
	})

	isUnique, err = userRepo.IsLeetcodeIDUnique(leetcodeID)
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "failed to get DB connection")
}
