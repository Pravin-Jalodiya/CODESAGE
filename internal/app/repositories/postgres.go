package repositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"sync"
	"time"
)

var (
	db          *sql.DB
	dbErr       error
	dbMutex     sync.Mutex // Mutex to handle connection expiration logic
	connectedAt time.Time
	dbTTL       = 1 * time.Hour
	connStr     = "postgres://username:password@localhost:5432/codesage?sslmode=disable"
)

// CreateContext creates a context with a timeout for database operations.
func CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func GetPostgresClient() (*sql.DB, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// If the connection has expired or does not exist, create a new one.
	if db == nil || time.Since(connectedAt) > dbTTL { // Example: 1-hour expiration
		// Close the old client if it exists
		if db != nil {
			if err := db.Close(); err != nil {
				log.Printf("Failed to close old PostgreSQL client: %v", err)
			}
		}

		db, dbErr = sql.Open("pgx", connStr)
		if dbErr != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", dbErr)
		}

		// Ping the database to ensure connection is successful
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
		}

		connectedAt = time.Now() // Update the connection time
	}
	return db, dbErr
}

// ClosePostgresClient closes the PostgreSQL client connection gracefully.
func ClosePostgresClient() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Failed to disconnect PostgreSQL: %v", err)
		}
		db = nil
		log.Println("PostgreSQL connection closed.")
	}
}
