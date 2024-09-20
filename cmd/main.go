package main

import (
	"cli-project/external/api"
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/routes"
	"cli-project/internal/app/repositories"
	"cli-project/internal/app/services"
	"cli-project/internal/config"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	defer repositories.ClosePostgresClient()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %s. Shutting down gracefully...", sig)

		// Call the function to close Postgres client
		repositories.ClosePostgresClient()

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

	r := mux.NewRouter()
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	routes.InitialiseAuthRouter(r, authHandler)
	routes.InitialiseUserRouter(r, userHandler)
	http.Handle("/", r)
	fmt.Println("server is running on port:", config.PORT)
	log.Fatal(http.ListenAndServe(config.PORT, nil))

	//// Initialize UI
	//newUI := ui.NewUI(authService, userService, questionService, bufio.NewReader(os.Stdin))
	//if newUI == nil {
	//	log.Fatal("Failed to initialize UI")
	//}
	//
	//// Show Main Menu
	//newUI.ShowMainMenu()
}
