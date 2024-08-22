package models

type Question struct {
	QuestionID    string   `bson:"question_id"`
	QuestionTitle string   `bson:"question_title"`
	Difficulty    string   `bson:"difficulty"`
	QuestionLink  string   `bson:"question_link"`
	TopicTags     []string `bson:"topic_tags"`
	CompanyTags   []string `bson:"company_tags"`
}
