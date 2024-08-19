package models

type UserData struct {
	UserID          int   `json:"user_id"`
	QuestionsSolved []int `json:"questionsSolved"`
	LastSeenInHours int   `json:"last_seen"`
}
