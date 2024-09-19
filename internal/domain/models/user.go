package models

import "time"

type User struct {
	ID           string `json:"id"`
	Role         string `json:"role"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organisation string `json:"organisation"`
	Country      string `json:"country"`
	IsBanned     bool   `json:"is_banned"`
}

type Admin struct {
	Admin User
}

type StandardUser struct {
	StandardUser    User      `json:",inline"`
	LeetcodeID      string    `json:"leetcode_id"`
	QuestionsSolved []string  `json:"questions_solved"`
	LastSeen        time.Time `json:"last_seen"`
}
