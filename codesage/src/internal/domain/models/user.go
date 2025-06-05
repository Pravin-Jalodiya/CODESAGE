package models

import "time"

type User struct {
	ID           string `json:"id" dynamodbav:"id"`
	Role         string `json:"role" dynamodbav:"role"`
	Username     string `json:"username" dynamodbav:"username"`
	Password     string `json:"password" dynamodbav:"password"`
	Name         string `json:"name" dynamodbav:"name"`
	Email        string `json:"email" dynamodbav:"email"`
	Organisation string `json:"organisation" dynamodbav:"organisation"`
	Country      string `json:"country" dynamodbav:"country"`
	IsBanned     bool   `json:"is_banned" dynamodbav:"is_banned"`
	Avatar       string `json:"avatar" dynamodbav:"avatar"`
}

type Admin struct {
	Admin User `json:"admin" dynamodbav:"admin"`
}

type StandardUser struct {
	User            `json:",inline" dynamodbav:",inline"`
	LeetcodeID      string    `json:"leetcode_id" dynamodbav:"leetcode_id"`
	QuestionsSolved []string  `json:"questions_solved" dynamodbav:"questions_solved"`
	LastSeen        time.Time `json:"last_seen" dynamodbav:"last_seen"`
}
