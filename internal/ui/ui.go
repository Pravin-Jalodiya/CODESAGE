package ui

import (
	"bufio"
	"cli-project/internal/app/services"
	"os"
)

// UI struct holds the UserService, bufio.Reader, and other dependencies
type UI struct {
	//authService    *services.authService
	userService     *services.UserService
	questionService *services.QuestionService
	reader          *bufio.Reader
}

// NewUI initializes the UI with the provided services and a bufio.Reader
func NewUI(userService *services.UserService, questionService *services.QuestionService) *UI {
	return &UI{
		userService:     userService,
		questionService: questionService,
		//authService:	authService,
		reader: bufio.NewReader(os.Stdin), // Initialize the reader to read from standard input
	}
}
