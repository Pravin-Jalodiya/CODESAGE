package main

import (
	"bufio"
	"cli-project/external/api"
	"cli-project/internal/app/repositories"
	"cli-project/internal/app/services"
	"cli-project/internal/ui"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	defer repositories.CloseMongoClient()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s. Shutting down gracefully...", sig)

		// Call the function to close MongoDB client
		repositories.CloseMongoClient()

		os.Exit(0)
	}()

	// Initialize User Repository
	userRepo := repositories.NewUserRepo()
	if userRepo == nil {
		log.Fatal("Failed to initialize UserRepository")
	}

	// Initialize Question Repository
	questionRepo := repositories.NewQuestionRepo()
	if questionRepo == nil {
		log.Fatal("Failed to initialize QuestionRepository")
	}

	// Initialize Question Service
	questionService := services.NewQuestionService(questionRepo)
	if questionService == nil {
		log.Fatal("Failed to initialize QuestionService")
	}

	// Initialize Leetcode Service
	LeetcodeAPI := api.NewLeetcodeAPI()

	// Initialize User Service
	userService := services.NewUserService(userRepo, questionService, LeetcodeAPI)
	if userService == nil {
		log.Fatal("Failed to initialize UserService")
	}

	// Initialize Auth Service
	authService := services.NewAuthService(userRepo, LeetcodeAPI)
	if authService == nil {
		log.Fatal("Failed to initialize AuthService")
	}

	// Initialize UI
	newUI := ui.NewUI(authService, userService, questionService, bufio.NewReader(os.Stdin))
	if newUI == nil {
		log.Fatal("Failed to initialize UI")
	}

	// Show Main Menu
	newUI.ShowMainMenu()
}
