package models

import "time"

type User struct {
	ID       string `json:"id"`
	Role     string `json:"role,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Admin struct {
	Admin User
}

type StandardUser struct {
	StandardUser    User
	LeetcodeID      string    `json:"leetcode_id"`
	QuestionsSolved []int     `json:"questionsSolved"`
	LastSeenInHours time.Time `json:"last_seen"`
}
