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
	collection *mongo.Collection
}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{
		collection: client.Database(config.DB_NAME).Collection(config.QUESTION_COLLECTION),
	}
}

func (r *questionRepo) AddQuestionsByID(questionID *[]string) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) AddQuestions(questions *[]models.Question) error {

	ctx, cancel := CreateContext()
	defer cancel()

	var documents []interface{} = make([]interface{}, len(*questions))
	for i, question := range *questions {
		documents[i] = question
	}

	_, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("could not insert questions: %v", err)
	}

	return nil
}

func (r *questionRepo) RemoveQuestionByID(questionID string) error {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"question_id": questionID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete question: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("question with ID %s not found", questionID)
	}

	return nil
}

func (r *questionRepo) FetchQuestionByID(questionID string) (*models.Question, error) {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"questions_id": questionID}

	var question models.Question
	err := r.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &models.Question{}, fmt.Errorf("question with ID %s not found", questionID)
		}
		return &models.Question{}, fmt.Errorf("could not fetch question: %v", err)
	}

	return &question, nil
}

func (r *questionRepo) FetchAllQuestions() (*[]models.Question, error) {
	ctx, cancel := CreateContext()
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
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

	fmt.Printf("Constructed filter: %+v\n", filter) // debug

	cursor, err := r.collection.Find(ctx, filter)
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

	fmt.Printf("Number of questions found: %d\n", len(questions)) // Debug: Print number of questions

	return &questions, nil
}

//func (r *questionRepo) FetchQuestionsByFilters(difficulty, company, topic string) (*[]models.Question, error) {
//	ctx, cancel := CreateContext()
//	defer cancel()
//
//	filter := bson.M{}
//	if difficulty != "" {
//		filter["difficulty"] = difficulty
//	}
//	if company != "" {
//		filter["company_tags"] = company
//	}
//	if topic != "" {
//		filter["topic_tags"] = topic
//	}
//
//	cursor, err := r.collection.Find(ctx, filter)
//	if err != nil {
//		return nil, fmt.Errorf("could not fetch questions by filters: %v", err)
//	}
//
//	defer func(cursor *mongo.Cursor, ctx context.Context) {
//		err := cursor.Close(ctx)
//		if err != nil {
//
//		}
//	}(cursor, ctx)
//
//	var questions []models.Question
//	for cursor.Next(ctx) {
//		var question models.Question
//		if err := cursor.Decode(&question); err != nil {
//			return nil, fmt.Errorf("could not decode question: %v", err)
//		}
//		questions = append(questions, question)
//	}
//
//	if err := cursor.Err(); err != nil {
//		return nil, fmt.Errorf("cursor error: %v", err)
//	}
//
//	return &questions, nil
//}

func (r *questionRepo) CountQuestions() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("could not count questions: %v", err)
	}

	return count, nil
}

func (r *questionRepo) QuestionExists(questionID string) (bool, error) {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"question_id": questionID}
	var existingQuestion models.Question
	err := r.collection.FindOne(ctx, filter).Decode(&existingQuestion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
