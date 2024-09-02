package repositories

//
//import (
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"log"
//	"sync"
//	"time"
//)
//
//var (
//	client      *mongo.Client
//	clientErr   error
//	mongoMutex  sync.Mutex // Mutex to handle connection expiration logic
//	connectedAt time.Time
//	clientTTL   = 1 * time.Hour
//)
//
//// CreateContext creates a context with a timeout for database operations.
//func CreateContext() (context.Context, context.CancelFunc) {
//	return context.WithTimeout(context.Background(), 10*time.Second)
//}
//
//func GetMongoClient() (*mongo.Client, error) {
//	mongoMutex.Lock()
//	defer mongoMutex.Unlock()
//
//	// If the connection has expired or does not exist, create a new one.
//	if client == nil || time.Since(connectedAt) > clientTTL { // Example: 1-hour expiration
//		// Close the old client if it exists
//		if client != nil {
//			if err := client.Disconnect(context.TODO()); err != nil {
//				log.Printf("Failed to disconnect old MongoDB client: %v", err)
//			}
//		}
//
//		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
//		client, clientErr = mongo.Connect(context.TODO(), clientOptions)
//		if clientErr != nil {
//			log.Fatalf("Failed to connect to MongoDB: %v", clientErr)
//		}
//
//		if err := client.Ping(context.TODO(), nil); err != nil {
//			return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
//		}
//
//		connectedAt = time.Now() // Update the connection time
//	}
//	return client, clientErr
//}
//
//// CloseMongoClient closes the MongoDB client connection gracefully.
//func CloseMongoClient() {
//	mongoMutex.Lock()
//	defer mongoMutex.Unlock()
//
//	if client != nil {
//		if err := client.Disconnect(context.TODO()); err != nil {
//			log.Printf("Failed to disconnect MongoDB: %v", err)
//		}
//		client = nil
//		log.Println("MongoDB connection closed.")
//	}
//}
