package main

import (
	"cli-project/internal/app/repositories"
	"cli-project/internal/app/services"
	"cli-project/internal/ui"
	"log"
)

func main() {
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

	// Initialize Auth Service
	authService := services.NewAuthService(userRepo)
	if authService == nil {
		log.Fatal("Failed to initialize AuthService")
	}

	// Initialize Question Service
	questionService := services.NewQuestionService(questionRepo)
	if questionService == nil {
		log.Fatal("Failed to initialize QuestionService")
	}

	// Initialize User Service
	userService := services.NewUserService(userRepo, questionService)
	if userService == nil {
		log.Fatal("Failed to initialize UserService")
	}

	// Initialize UI
	newUI := ui.NewUI(authService, userService, questionService)
	if newUI == nil {
		log.Fatal("Failed to initialize UI")
	}

	// Show Main Menu
	newUI.ShowMainMenu()
}
