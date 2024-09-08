package repositories_test

import (
	"cli-project/internal/app/repositories"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// to ignore the logs
	log.SetOutput(os.Stdout)
	os.Exit(m.Run())
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true)) // Enable ping monitoring
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db, mock
}

func TestGetPostgresClient_NewConnection(t *testing.T) {
	_, mock := setupMockDB(t)

	// If simulating the replacement of an expired or otherwise closed connection,
	// you must ensure that your test scenario aligns with this:
	if repositories.Db != nil {
		mock.ExpectClose() // Expect the old connection to close if it exists and being replaced
	}

	mock.ExpectPing().WillReturnError(nil) // Expect a successful ping for a new connection setup

	repositories.Db = nil                                              // simulate no existing connection
	repositories.ConnectedAt = time.Now().Add(-2 * repositories.DbTTL) // ensure expiration

	client, err := repositories.GetPostgresClient()

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestGetPostgresClient_ExistingConnection(t *testing.T) {
	db, mock := setupMockDB(t)

	mock.ExpectPing().WillReturnError(nil) // Expect a successful ping if the connection is reused

	repositories.Db = db
	repositories.ConnectedAt = time.Now() // non-expired

	client, err := repositories.GetPostgresClient()

	assert.NoError(t, err)
	assert.Equal(t, db, client)
}

func TestClosePostgresClient(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectClose()

	repositories.Db = db // set the mocked db

	repositories.ClosePostgresClient()

	assert.Nil(t, repositories.Db)
	assert.NoError(t, mock.ExpectationsWereMet())
}
