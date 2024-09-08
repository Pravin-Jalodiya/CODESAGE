package repositories_test

//func TestAddQuestionsFailure(t *testing.T) {
//	teardown := setup(t)
//	defer teardown()
//
//	// Define a slice of Question models
//	questions := []models.Question{
//		{
//			QuestionTitleSlug: "example-question-fail",
//			QuestionID:        "2",
//			QuestionTitle:     "Example Question Failure",
//			Difficulty:        "Hard",
//			QuestionLink:      "http://examplefail.com",
//			TopicTags:         []string{"trees", "graphs"},
//			CompanyTags:       []string{"facebook", "amazon"},
//		},
//	}
//
//	// Define expected SQL operations
//	query := `
//		INSERT INTO questions
//		(title_slug, id, title, difficulty, link, topic_tags, company_tags)
//		VALUES ($1, $2, $3, $4, $5, $6, $7)
//		ON CONFLICT (title_slug) DO NOTHING;
//	`
//
//	// Begin transaction
//	mock.ExpectBegin()
//
//	// Expect question insertion to fail
//	mock.ExpectExec(query).WithArgs(
//		questions[0].QuestionTitleSlug,
//		questions[0].QuestionID,
//		questions[0].QuestionTitle,
//		questions[0].Difficulty,
//		questions[0].QuestionLink,
//		pq.Array(questions[0].TopicTags),
//		pq.Array(questions[0].CompanyTags),
//	).WillReturnError(sql.ErrConnDone) // Simulate a connection error
//
//	// Expect transaction to be rolled back due to failure
//	mock.ExpectRollback()
//
//	// Call the method under test
//	err := questionRepo.AddQuestions(&questions)
//	if err == nil {
//		t.Errorf("Expected error, got none")
//	}
//
//	// Verify that all expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("Not all expectations were met: %s", err)
//	}
//}
