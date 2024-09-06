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
