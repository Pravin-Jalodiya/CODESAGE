package ui

import (
	"bufio"
	"cli-project/internal/domain/interfaces"
)

// UI struct holds the UserService, bufio.Reader, and other dependencies
type UI struct {
	authService     interfaces.AuthService
	userService     interfaces.UserService
	questionService interfaces.QuestionService
	reader          *bufio.Reader
}

// NewUI initializes the UI with the provided services and a bufio.Reader
func NewUI(authService interfaces.AuthService, userService interfaces.UserService, questionService interfaces.QuestionService, reader *bufio.Reader) *UI {
	return &UI{
		authService:     authService,
		userService:     userService,
		questionService: questionService,
		reader:          reader, // Initialize the reader to read from standard input
	}
}
