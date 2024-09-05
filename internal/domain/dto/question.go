package dto

type Question struct {
	QuestionID    string   `json:"question_id" db:"id"`
	QuestionTitle string   `json:"question_title" db:"title"`
	Difficulty    string   `json:"difficulty" db:"difficulty"`
	QuestionLink  string   `json:"question_link" db:"link"`
	TopicTags     []string `json:"topic_tags" db:"topic_tags"`
	CompanyTags   []string `json:"company_tags" db:"company_tags"`
}
