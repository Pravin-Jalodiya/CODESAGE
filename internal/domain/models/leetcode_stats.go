package models

type LeetcodeStats struct {
	TotalQuestionsCount     int
	TotalQuestionsDoneCount int
	TotalEasyCount          int
	TotalMediumCount        int
	TotalHardCount          int
	EasyDoneCount           int
	MediumDoneCount         int
	HardDoneCount           int
	RecentACSubmissions     []string `json:"recent_ac_submissions"`
}
