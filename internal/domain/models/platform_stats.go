package models

type PlatformStats struct {
	ActiveUserInLast24Hours      int
	TotalQuestionsCount          int
	DifficultyWiseQuestionsCount map[string]int
	TopicWiseQuestionsCount      map[string]int
	CompanyWiseQuestionsCount    map[string]int
}
