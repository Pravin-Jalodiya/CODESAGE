package db

import (
	"cli-project/internal/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"sync"
	"time"
)

var (
	Db                  *sql.DB
	DbErr               error
	DbMutex             sync.Mutex // Mutex to handle connection expiration logic
	ConnectedAt         time.Time
	DbTTL               = 6 * time.Hour
	ConnStr             = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	MaxOpenConns        = 50
	MaxIdleConns        = 10
	ConnMaxLifetime     = 30 * time.Minute
	IdleConnMaxLifetime = 10 * time.Minute
	ClientGetter        = defaultGetPostgresClient
)

// CreateContext creates a context with a timeout for database operations.
func CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func defaultGetPostgresClient() (*sql.DB, error) {
	db, err := GetPostgresClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}
	return db, nil
}

// UseDBClient allows for injecting a custom DB client getter function (used in tests).
func UseDBClient(getter func() (*sql.DB, error)) {
	DbMutex.Lock()
	defer DbMutex.Unlock()
	ClientGetter = getter
}

func GetPostgresClient() (*sql.DB, error) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	// If the connection has expired or does not exist, create a new one.
	if Db == nil || time.Since(ConnectedAt) > DbTTL {
		// Close the old client if it exists
		if Db != nil {
			if err := Db.Close(); err != nil {
				log.Printf("Failed to close old PostgreSQL client: %v", err)
			}
		}

		Db, DbErr = sql.Open("pgx", ConnStr)
		if DbErr != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", DbErr)
		}
		Db.SetMaxOpenConns(MaxOpenConns)
		Db.SetMaxIdleConns(MaxIdleConns)
		Db.SetConnMaxLifetime(ConnMaxLifetime)
		Db.SetConnMaxIdleTime(IdleConnMaxLifetime)

		// Ping the database to ensure connection is successful
		if err := Db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
		}

		ConnectedAt = time.Now() // Update the connection time
	}
	return Db, DbErr
}

// ClosePostgresClient closes the PostgreSQL client connection gracefully.
func ClosePostgresClient() {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	if Db != nil {
		if err := Db.Close(); err != nil {
			log.Printf("Failed to disconnect PostgreSQL: %v", err)
		}
		Db = nil
		log.Println("PostgreSQL connection closed.")
	}
}
