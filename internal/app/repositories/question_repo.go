package repositories

import (
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
	"strings"
)

type questionRepo struct {
}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{}
}

func (r *questionRepo) getTableName() string {
	return "questions"
}

// getDBConnection returns a PostgreSQL client connection and handles errors.
func (r *questionRepo) getDBConnection() (*sql.DB, error) {
	db, err := GetPostgresClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}
	return db, nil
}

func (r *questionRepo) AddQuestionsByID(questionID *[]string) error {
	// Placeholder implementation
	return nil
}

func (r *questionRepo) AddQuestions(questions *[]models.Question) error {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Prepare SQL query for bulk insertion
	query := `
		INSERT INTO questions 
		(title_slug, id, title, difficulty, link, topic_tags, company_tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (title_slug) DO NOTHING;
	`

	tx, err := db.Begin() // Start a transaction
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	for _, question := range *questions {
		_, err := tx.Exec(query,
			question.QuestionTitleSlug,
			question.QuestionID,
			question.QuestionTitle,
			question.Difficulty,
			question.QuestionLink,
			question.TopicTags,
			question.CompanyTags,
		)
		if err != nil {
			tx.Rollback() // Rollback transaction on error
			return fmt.Errorf("could not insert question: %v", err)
		}
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

func (r *questionRepo) RemoveQuestionByID(questionID string) error {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Define the SQL query to delete a question by its title slug
	query := `DELETE FROM questions WHERE id = $1`

	// Execute the query
	result, err := db.ExecContext(context.Background(), query, questionID)
	if err != nil {
		return fmt.Errorf("could not delete question: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("question with title slug %s not found", questionID)
	}

	return nil
}

func (r *questionRepo) FetchQuestionByID(questionID string) (*models.Question, error) {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Define the SQL query to fetch a question by its title slug
	query := `
		SELECT title_slug, id, title, difficulty, link, topic_tags, company_tags
		FROM questions
		WHERE id = $1
	`

	// Execute the query
	row := db.QueryRowContext(context.Background(), query, questionID)

	var question models.Question
	err = row.Scan(
		&question.QuestionTitleSlug,
		&question.QuestionID,
		&question.QuestionTitle,
		&question.Difficulty,
		&question.QuestionLink,
		&question.TopicTags,
		&question.CompanyTags,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("question with title slug %s not found", questionID)
		}
		return nil, fmt.Errorf("could not fetch question: %v", err)
	}

	return &question, nil
}

func (r *questionRepo) FetchAllQuestions() (*[]dto.Question, error) {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Define the SQL query to fetch all questions excluding title_slug
	query := `
		SELECT id, title, difficulty, link, topic_tags, company_tags
		FROM questions
	`

	// Execute the query
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions: %v", err)
	}
	defer rows.Close()

	var questions []dto.Question

	// Iterate over the rows
	for rows.Next() {
		var (
			id          string
			title       string
			difficulty  string
			link        string
			topicTags   []string
			companyTags []string
		)

		// Use the appropriate scan method for arrays
		err := rows.Scan(&id, &title, &difficulty, &link, pq.Array(&topicTags), pq.Array(&companyTags))
		if err != nil {
			return nil, fmt.Errorf("could not scan question: %v", err)
		}

		questions = append(questions, dto.Question{
			QuestionID:    id,
			QuestionTitle: title,
			Difficulty:    difficulty,
			QuestionLink:  link,
			TopicTags:     topicTags,
			CompanyTags:   companyTags,
		})
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return &questions, nil
}

func (r *questionRepo) FetchQuestionsByFilters(difficulty, topic, company string) (*[]dto.Question, error) {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Initialize the base query
	query := `SELECT id, title, difficulty, link, topic_tags, company_tags
	          FROM questions
	          WHERE TRUE`

	var args []interface{}
	argIndex := 1

	// Add filters based on user input
	if difficulty != "" && strings.ToLower(difficulty) != "any" {
		query += fmt.Sprintf(" AND difficulty = $%d", argIndex)
		args = append(args, difficulty)
		argIndex++
	}

	if topic != "" && strings.ToLower(topic) != "any" {
		query += fmt.Sprintf(" AND topic_tags @> $%d::varchar[]", argIndex)
		args = append(args, pq.Array([]string{topic}))
		argIndex++
	}

	if company != "" && strings.ToLower(company) != "any" {
		query += fmt.Sprintf(" AND $%d = ANY(company_tags::varchar[])", argIndex)
		args = append(args, company)
		argIndex++
	}

	// Log the query and arguments
	fmt.Printf("Executing query: %s\n", query)
	fmt.Printf("With args: %v\n", args)

	// Execute the query
	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions by filters: %v", err)
	}
	defer rows.Close()

	var questions []dto.Question

	// Iterate over the rows
	for rows.Next() {
		var (
			id          string
			title       string
			difficulty  string
			link        string
			topicTags   []string
			companyTags []string
		)

		err := rows.Scan(
			&id,
			&title,
			&difficulty,
			&link,
			pq.Array(&topicTags),
			pq.Array(&companyTags),
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan question: %v", err)
		}

		questions = append(questions, dto.Question{
			QuestionID:    id,
			QuestionTitle: title,
			Difficulty:    difficulty,
			QuestionLink:  link,
			TopicTags:     topicTags,
			CompanyTags:   companyTags,
		})
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return &questions, nil
}

func (r *questionRepo) CountQuestions() (int64, error) {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return 0, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Define the SQL query to count the number of questions
	query := `SELECT COUNT(*) FROM questions`

	// Execute the query
	var count int64
	err = db.QueryRowContext(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not count questions: %v", err)
	}

	return count, nil
}

func (r *questionRepo) QuestionExists(questionID string) (bool, error) {
	// Get the PostgreSQL connection
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}

	// Prepare the SQL query to check if the question exists by its title slug
	query := `SELECT EXISTS (SELECT 1 FROM questions WHERE id = $1)`

	var exists bool
	err = db.QueryRow(query, questionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if question exists: %v", err)
	}

	return exists, nil
}
