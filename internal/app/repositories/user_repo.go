package repositories

import (
	"cli-project/internal/config"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepo() interfaces.UserRepository {
	return &userRepo{
		collection: client.Database(config.DB_NAME).Collection(config.USER_COLLECTION),
	}
}

func (r *userRepo) RegisterUser(user models.StandardUser) error {
	// Convert the user model to BSON format
	userBson := bson.M{
		"id":               user.StandardUser.ID,
		"username":         user.StandardUser.Username,
		"password":         user.StandardUser.Password,
		"name":             user.StandardUser.Name,
		"email":            user.StandardUser.Email,
		"role":             user.StandardUser.Role,
		"leetcode_id":      user.LeetcodeID,
		"questions_solved": user.QuestionsSolved,
		"last_seen":        user.LastSeen,
	}

	// Insert the user document into the collection
	_, err := r.collection.InsertOne(context.TODO(), userBson)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}

	return nil
}

func (r *userRepo) UpdateUserProgress(username string, questionID int) error {
	return nil
}

func (r *userRepo) FetchAllUsers() ([]models.StandardUser, error) {
	return nil, nil
}

func (r *userRepo) FetchUser(username string) (models.StandardUser, error) {
	// Set a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}

	var user models.StandardUser

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments // User not found
		}
		return user, err
	}

	// Return the found user
	return user, nil
}

func (r *userRepo) CountActiveUsersInLast24Hours() (int64, error) {

	now := time.Now().UTC()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)

	filter := bson.M{
		"last_seen_in_hours": bson.M{
			"$gte": twentyFourHoursAgo,
		},
	}

	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, fmt.Errorf("could not count active users: %v", err)
	}

	return count, nil
}

func (r *userRepo) IsEmailUnique(email string) (bool, error) {
	var result models.StandardUser
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (r *userRepo) IsUsernameUnique(username string) (bool, error) {
	var result models.StandardUser
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (r *userRepo) IsLeetcodeIDUnique(leetcodeID string) (bool, error) {
	var result models.StandardUser
	err := r.collection.FindOne(context.Background(), bson.M{"leetcode_id": leetcodeID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
