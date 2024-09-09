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
	Db          *sql.DB
	DbErr       error
	DbMutex     sync.Mutex // Mutex to handle connection expiration logic
	ConnectedAt time.Time
	DbTTL       = 1 * time.Hour
	ConnStr     = "postgres://postgres:password@localhost:5432/codesage?sslmode=disable" // pass username, password, Db name and port for env variables or config file
)

// CreateContext creates a context with a timeout for database operations.
func CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func GetPostgresClient() (*sql.DB, error) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	// If the connection has expired or does not exist, create a new one.
	if Db == nil || time.Since(ConnectedAt) > DbTTL { // Example: 1-hour expiration
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
