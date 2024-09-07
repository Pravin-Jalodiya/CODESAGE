package models

type LeetcodeStats struct {
	TotalQuestionsCount          int
	TotalQuestionsDoneCount      int
	TotalEasyCount               int
	TotalMediumCount             int
	TotalHardCount               int
	EasyDoneCount                int
	MediumDoneCount              int
	HardDoneCount                int
	RecentACSubmissionTitles     []string `json:"recent_ac_submission_title"`
	RecentACSubmissionTitleSlugs []string `json:"recent_ac_submissions_title_slugs"`
}

type CodesageStats struct {
	TotalQuestionsCount     int
	TotalQuestionsDoneCount int
	TotalEasyCount          int
	TotalMediumCount        int
	TotalHardCount          int
	EasyDoneCount           int
	MediumDoneCount         int
	HardDoneCount           int
	CompanyWiseStats        map[string]int
	TopicWiseStats          map[string]int
}
