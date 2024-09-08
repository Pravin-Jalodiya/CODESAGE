package repositories_test

import (
	"cli-project/internal/domain/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"strings"
	"testing"
)

func TestFetchQuestionByID(t *testing.T) {
	defer setup(t)()

	questionID := "1"
	expectedQuestion := &models.Question{
		QuestionTitleSlug: "test-slug",
		QuestionID:        questionID,
		QuestionTitle:     "Test Question",
		Difficulty:        "Medium",
		QuestionLink:      "http://localhost",
		TopicTags:         []string{"tag1", "tag2"},
		CompanyTags:       []string{"Company1", "Company2"},
	}

	query := `SELECT title_slug, id, title, difficulty, link, topic_tags, company_tags FROM questions WHERE title_slug \= \$1`
	rows := sqlmock.NewRows([]string{
		"title_slug", "id", "title", "difficulty", "link", "topic_tags", "company_tags",
	}).AddRow(
		expectedQuestion.QuestionTitleSlug,
		expectedQuestion.QuestionID,
		expectedQuestion.QuestionTitle,
		expectedQuestion.Difficulty,
		expectedQuestion.QuestionLink,
		pq.Array(expectedQuestion.TopicTags),
		pq.Array(expectedQuestion.CompanyTags),
	)

	mock.ExpectQuery(query).WithArgs(questionID).WillReturnRows(rows)

	question, err := questionRepo.FetchQuestionByID(questionID)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if question.QuestionID != expectedQuestion.QuestionID {
		t.Errorf("expected %v, got %v", expectedQuestion, question)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRemoveQuestionByID(t *testing.T) {
	defer setup(t)()

	questionID := "1"

	mock.ExpectExec(`DELETE FROM questions WHERE id = \$1`).
		WithArgs(questionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := questionRepo.RemoveQuestionByID(questionID)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	// Verify all expectations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestRemoveQuestionByIDNotFound(t *testing.T) {
	defer setup(t)()

	questionID := "1"

	mock.ExpectExec(`DELETE FROM questions WHERE id = \$1`).
		WithArgs(questionID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // no rows affected

	err := questionRepo.RemoveQuestionByID(questionID)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestFetchQuestionByIDNoResult(t *testing.T) {
	defer setup(t)()

	questionID := "1"
	mock.ExpectQuery(`SELECT title_slug, id, title, difficulty, link, topic_tags, company_tags FROM questions WHERE title_slug = \$1`).
		WithArgs(questionID).
		WillReturnError(sql.ErrNoRows)

	_, err := questionRepo.FetchQuestionByID(questionID)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, did not receive it")
	}
}

func TestFetchAllQuestions(t *testing.T) {
	defer setup(t)()

	rows := sqlmock.NewRows([]string{"id", "title", "difficulty", "link", "topic_tags", "company_tags"}).
		AddRow("1", "Question 1", "Hard", "http://example.com/q1", pq.Array([]string{"tag1"}), pq.Array([]string{"comp1"})).
		AddRow("2", "Question 2", "Easy", "http://example.com/q2", pq.Array([]string{"tag2"}), pq.Array([]string{"comp2"}))

	mock.ExpectQuery(`SELECT id, title, difficulty, link, topic_tags, company_tags FROM questions`).
		WillReturnRows(rows)

	questions, err := questionRepo.FetchAllQuestions()
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if len(*questions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(*questions))
	}
}

func TestFetchQuestionsByFilters(t *testing.T) {
	defer setup(t)()

	rows := sqlmock.NewRows([]string{"id", "title", "difficulty", "link", "topic_tags", "company_tags"}).
		AddRow("1", "Question 1", "Easy", "http://example.com/q1", pq.Array([]string{"tag1"}), pq.Array([]string{"comp1"})).
		AddRow("2", "Question 2", "Medium", "http://example.com/q2", pq.Array([]string{"tag2"}), pq.Array([]string{"comp2"}))

	mock.ExpectQuery(`SELECT id, title, difficulty, link, topic_tags, company_tags FROM questions WHERE TRUE AND difficulty = \$1 AND topic_tags @> \$2::varchar\[] AND \$3 = ANY\(company_tags::varchar\[]\)`).
		WithArgs("Easy", pq.Array([]string{"tag1"}), "comp1").
		WillReturnRows(rows)

	questions, err := questionRepo.FetchQuestionsByFilters("Easy", "tag1", "comp1")
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if len(*questions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(*questions))
	}
}

func TestCountQuestions(t *testing.T) {
	defer setup(t)()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM questions`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := questionRepo.CountQuestions()
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if count != 5 {
		t.Errorf("expected count to be 5, got %d", count)
	}
}

func TestQuestionExistsByID(t *testing.T) {
	defer setup(t)()

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM questions WHERE id = \$1\)`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := questionRepo.QuestionExistsByID("1")
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if !exists {
		t.Errorf("expected exists to be true, got false")
	}
}

func TestQuestionExistsByTitleSlug(t *testing.T) {
	defer setup(t)()

	mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM questions WHERE title_slug = \$1\)`).
		WithArgs("test-slug").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := questionRepo.QuestionExistsByTitleSlug("test-slug")
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if !exists {
		t.Errorf("expected exists to be true, got false")
	}
}
