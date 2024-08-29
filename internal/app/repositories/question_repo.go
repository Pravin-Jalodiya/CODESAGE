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
	"strings"
	"time"
)

type questionRepo struct {
}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{}
}

func (r *questionRepo) getCollection() (*mongo.Collection, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	return client.Database(config.DB_NAME).Collection(config.QUESTION_COLLECTION), nil
}

func (r *questionRepo) AddQuestionsByID(questionID *[]string) error {
	// Placeholder implementation

	return nil
}

func (r *questionRepo) AddQuestions(questions *[]models.Question) error {

	collection, err := r.getCollection()
	if err != nil {
		return fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	var documents []interface{} = make([]interface{}, len(*questions))
	for i, question := range *questions {
		documents[i] = question
	}

	_, err = collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("could not insert questions: %v", err)
	}

	return nil
}

func (r *questionRepo) RemoveQuestionByID(questionID string) error {

	collection, err := r.getCollection()
	if err != nil {
		return fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"question_id": questionID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete question: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("question with ID %s not found", questionID)
	}

	return nil
}

func (r *questionRepo) FetchQuestionByID(questionID string) (*models.Question, error) {

	collection, err := r.getCollection()
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"questions_id": questionID}

	var question models.Question
	err = collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.Question{}, fmt.Errorf("question with ID %s not found", questionID)
		}
		return &models.Question{}, fmt.Errorf("could not fetch question: %v", err)
	}

	return &question, nil
}

func (r *questionRepo) FetchAllQuestions() (*[]models.Question, error) {

	collection, err := r.getCollection()
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions: %v", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var questions []models.Question
	for cursor.Next(ctx) {
		var question models.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, fmt.Errorf("could not decode question: %v", err)
		}
		questions = append(questions, question)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &questions, nil
}

func (r *questionRepo) FetchQuestionsByFilters(difficulty, company, topic string) (*[]models.Question, error) {

	collection, err := r.getCollection()
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{}

	// Apply filters only if parameters are not "any"
	if difficulty != "" && strings.ToLower(difficulty) != "any" {
		filter["difficulty"] = difficulty
	}
	if company != "" && strings.ToLower(company) != "any" {
		filter["company_tags"] = company
	}
	if topic != "" && strings.ToLower(topic) != "any" {
		filter["topic_tags"] = topic
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions by filters: %v", err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println("could not close cursor")
		}
	}(cursor, ctx)

	var questions []models.Question
	for cursor.Next(ctx) {
		var question models.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, fmt.Errorf("could not decode question: %v", err)
		}
		questions = append(questions, question)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &questions, nil
}

func (r *questionRepo) CountQuestions() (int64, error) {

	collection, err := r.getCollection()
	if err != nil {
		return 0, fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("could not count questions: %v", err)
	}

	return count, nil
}

func (r *questionRepo) QuestionExists(questionID string) (bool, error) {

	collection, err := r.getCollection()
	if err != nil {
		return false, fmt.Errorf("failed to get collection: %v", err)
	}

	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"question_id": questionID}
	var existingQuestion models.Question
	err = collection.FindOne(ctx, filter).Decode(&existingQuestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
