package ui

import (
	"bufio"
	"cli-project/internal/app/services"
	"os"
)

// UI struct holds the UserService, bufio.Reader, and other dependencies
type UI struct {
	authService     *services.AuthService
	userService     *services.UserService
	questionService *services.QuestionService
	reader          *bufio.Reader
}

// NewUI initializes the UI with the provided services and a bufio.Reader
func NewUI(authService *services.AuthService, userService *services.UserService, questionService *services.QuestionService) *UI {
	return &UI{
		authService:     authService,
		userService:     userService,
		questionService: questionService,
		reader:          bufio.NewReader(os.Stdin), // Initialize the reader to read from standard input
	}
}
