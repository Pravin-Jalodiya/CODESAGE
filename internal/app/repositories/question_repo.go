package repositories

import (
	"cli-project/internal/config/queries"
	database "cli-project/internal/db"
	"cli-project/internal/domain/dto"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
)

type questionRepo struct{}

func NewQuestionRepo() interfaces.QuestionRepository {
	return &questionRepo{}
}

func (r *questionRepo) getDBConnection() (*sql.DB, error) {
	return database.ClientGetter()
}

func (r *questionRepo) AddQuestions(ctx context.Context, questions *[]models.Question) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
	}

	query := queries.QueryBuilder(queries.BaseInsert, map[string]string{
		"table":   "questions",
		"columns": "title_slug, id, title, difficulty, link, topic_tags, company_tags",
		"values":  "$1, $2, $3, $4, $5, $6, $7",
	}) + " ON CONFLICT (title_slug) DO NOTHING;"

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
	}

	for _, question := range *questions {
		_, err := tx.ExecContext(ctx, query,
			question.QuestionTitleSlug,
			question.QuestionID,
			question.QuestionTitle,
			question.Difficulty,
			question.QuestionLink,
			pq.Array(question.TopicTags),
			pq.Array(question.CompanyTags),
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
	}

	return nil
}

func (r *questionRepo) UpdateQuestion(ctx context.Context, question *models.Question) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
	}

	query := `UPDATE questions SET topic_tags = $1, company_tags = $2 WHERE id = $3`
	_, err = db.ExecContext(ctx, query, pq.Array(question.TopicTags), pq.Array(question.CompanyTags), question.QuestionID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDbOperation, err)
	}

	return nil
}

func (r *questionRepo) RemoveQuestionByID(ctx context.Context, questionID string) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseDelete, map[string]string{
		"table":      "questions",
		"conditions": "id = $1",
	})

	result, err := db.ExecContext(ctx, query, questionID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrCheckRowsAffected, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: question with ID %s not found", errs.ErrNoRows, questionID)
	}

	return nil
}

func (r *questionRepo) FetchQuestionByTitleSlug(ctx context.Context, titleSlug string) (*models.Question, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "title_slug, id, title, difficulty, link, topic_tags, company_tags",
		"table":      "questions",
		"conditions": "title_slug = $1",
	})

	row := db.QueryRowContext(ctx, query, titleSlug)

	var question models.Question
	var topicTags, companyTags []string

	err = row.Scan(
		&question.QuestionTitleSlug,
		&question.QuestionID,
		&question.QuestionTitle,
		&question.Difficulty,
		&question.QuestionLink,
		pq.Array(&topicTags),
		pq.Array(&companyTags),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: question with ID %s not found", errs.ErrNoRows, titleSlug)
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}

	question.TopicTags = topicTags
	question.CompanyTags = companyTags

	return &question, nil
}

func (r *questionRepo) FetchAllQuestions(ctx context.Context) ([]dto.Question, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelect, map[string]string{
		"columns": "id, title, difficulty, link, topic_tags, company_tags",
		"table":   "questions",
	})

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}
	defer rows.Close()

	var questions []dto.Question

	for rows.Next() {
		var (
			id          string
			title       string
			difficulty  string
			link        string
			topicTags   []string
			companyTags []string
		)

		err := rows.Scan(&id, &title, &difficulty, &link, pq.Array(&topicTags), pq.Array(&companyTags))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrRows, err)
	}

	return questions, nil
}

func (r *questionRepo) FetchQuestionsByFilters(ctx context.Context, difficulty, topic, company string) ([]dto.Question, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelect, map[string]string{
		"columns": "id, title, difficulty, link, topic_tags, company_tags",
		"table":   "questions",
	}) + " WHERE TRUE"

	var args []interface{}
	argIndex := 1

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

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}
	defer rows.Close()

	var questions []dto.Question

	for rows.Next() {
		var (
			id          string
			title       string
			difficulty  string
			link        string
			topicTags   []string
			companyTags []string
		)

		err = rows.Scan(&id, &title, &difficulty, &link, pq.Array(&topicTags), pq.Array(&companyTags))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrRows, err)
	}

	return questions, nil
}

func (r *questionRepo) CountQuestions(ctx context.Context) (int, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelect, map[string]string{
		"columns": "COUNT(*)",
		"table":   "questions",
	})

	var count int
	err = db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}

	return count, nil
}

func (r *questionRepo) QuestionExistsByID(ctx context.Context, questionID string) (bool, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectExistsWhere, map[string]string{
		"table":      "questions",
		"conditions": "id = $1",
	})

	var exists bool
	err = db.QueryRowContext(ctx, query, questionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}

	return exists, nil
}

func (r *questionRepo) QuestionExistsByTitleSlug(ctx context.Context, titleSlug string) (bool, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectExistsWhere, map[string]string{
		"table":      "questions",
		"conditions": "title_slug = $1",
	})

	var exists bool
	err = db.QueryRowContext(ctx, query, titleSlug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrQueryExecution, err)
	}

	return exists, nil
}
