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
)

type questionRepo struct {
	collection *mongo.Collection
}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{
		collection: client.Database(config.DB_NAME).Collection(config.QUESTION_COLLECTION),
	}
}

func (r *questionRepo) AddQuestionsByID(questionID []int) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) AddQuestions(questions []models.Question) error {
	ctx, cancel := CreateContext()
	defer cancel()

	for _, question := range questions {
		// Check if the question already exists
		filter := bson.M{"questions_id": question.QuestionID}
		var existingQuestion models.Question
		err := r.collection.FindOne(ctx, filter).Decode(&existingQuestion)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("could not check if question exists: %v", err)
		}

		// If the question does not exist, add it to the database
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, err := r.collection.InsertOne(ctx, question)
			if err != nil {
				return fmt.Errorf("could not insert question: %v", err)
			}
		}
	}
	return nil
}

func (r *questionRepo) RemoveQuestionByID(questionID int) error {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"questions_id": questionID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete question: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("question with ID %d not found", questionID)
	}

	return nil
}

func (r *questionRepo) FetchQuestionByID(questionID int) (models.Question, error) {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{"questions_id": questionID}

	var question models.Question
	err := r.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Question{}, fmt.Errorf("question with ID %d not found", questionID)
		}
		return models.Question{}, fmt.Errorf("could not fetch question: %v", err)
	}

	return question, nil
}

func (r *questionRepo) FetchAllQuestions() ([]models.Question, error) {
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

	return questions, nil
}

func (r *questionRepo) FetchQuestionsByFilters(difficulty, company, topic string) ([]models.Question, error) {
	ctx, cancel := CreateContext()
	defer cancel()

	filter := bson.M{}
	if difficulty != "" {
		filter["difficulty"] = difficulty
	}
	if company != "" {
		filter["company_tags"] = company
	}
	if topic != "" {
		filter["topic_tags"] = topic
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions by filters: %v", err)
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

	return questions, nil
}
