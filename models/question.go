package models

type Question struct {
	QuestionsID  string   `json:"questions_id"`
	QuestionLink string   `json:"question_link"`
	CompanyTags  []string `json:"company_tags"`
	TopicTags    []string `json:"topic_tags"`
	Difficulty   string   `json:"difficulty"`
}
