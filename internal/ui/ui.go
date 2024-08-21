package ui

import (
	"bufio"
	"cli-project/internal/app/services"
	"os"
)

// UI struct holds the UserService, bufio.Reader, and other dependencies
type UI struct {
	userService *services.UserService
	reader      *bufio.Reader
}

// NewUI initializes the UI with the provided services and a bufio.Reader
func NewUI(userService *services.UserService) *UI {
	return &UI{
		userService: userService,
		reader:      bufio.NewReader(os.Stdin), // Initialize the reader to read from standard input
	}
}
