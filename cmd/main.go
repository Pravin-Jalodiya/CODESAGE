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

	// Initialize User Service
	userService := services.NewUserService(userRepo)
	if userService == nil {
		log.Fatal("Failed to initialize UserService")
	}

	// Initialize UI
	ui := ui.NewUI(userService)
	if ui == nil {
		log.Fatal("Failed to initialize UI")
	}

	// Show Main Menu
	for {
		ui.ShowMainMenu()
	}
}
