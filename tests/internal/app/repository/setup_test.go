package repositories_test

import (
	"cli-project/internal/app/repositories"
	db2 "cli-project/internal/db"
	"cli-project/internal/domain/interfaces"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

var (
	db           *sql.DB
	mock         sqlmock.Sqlmock
	userRepo     interfaces.UserRepository
	questionRepo interfaces.QuestionRepository
)

func setup(t *testing.T) func() {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Inject the mock database connection
	db2.UseDBClient(func() (*sql.DB, error) {
		return db, nil
	})

	questionRepo = repositories.NewQuestionRepo()

	// Initialize the UserRepo
	userRepo = repositories.NewUserRepo()

	return func() {
		db.Close()
	}
}
