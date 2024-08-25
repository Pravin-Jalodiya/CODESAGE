package models

import "time"

type User struct {
	ID       string `bson:"id"`
	Role     string `bson:"role"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
}

type Admin struct {
	Admin User
}

type StandardUser struct {
	StandardUser    User      `bson:",inline"`
	LeetcodeID      string    `bson:"leetcode_id"`
	QuestionsSolved []string  `bson:"questions_solved"`
	LastSeen        time.Time `bson:"last_seen"`
}
