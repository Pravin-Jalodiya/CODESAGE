package dto

import "time"

type User struct {
	Role         string `json:"role"`
	Username     string `json:"username"`
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
	User       `json:",inline"`
	LeetcodeID string    `json:"leetcode_id"`
	LastSeen   time.Time `json:"last_seen"`
}
