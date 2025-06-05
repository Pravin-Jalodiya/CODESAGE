package main

import (
	"codesage/external/api"
	"codesage/internal/api/handlers"
	"codesage/internal/api/routes"
	"codesage/internal/app/repositories"
	"codesage/internal/app/services"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func createRouter() *mux.Router {

	// Initialize User Repository
	userRepo := repositories.NewUserRepo()
	if userRepo == nil {
		log.Fatal("Failed to initialize UserRepository")
	}

	// Initialize Question Repository
	//questionRepo := repositories.NewQuestionRepo()
	//if questionRepo == nil {
	//	log.Fatal("Failed to initialize QuestionRepository")
	//}

	// Initialize Question Service
	//questionService := services.NewQuestionService(questionRepo)
	//if questionService == nil {
	//	log.Fatal("Failed to initialize QuestionService")
	//}

	// Initialize Leetcode Service
	LeetcodeAPI := api.NewLeetcodeAPI()

	// Initialize User Service
	//userService := services.NewUserService(userRepo, questionService, LeetcodeAPI)
	//if userService == nil {
	//	log.Fatal("Failed to initialize UserService")
	//}

	// Initialize Auth Service
	authService := services.NewAuthService(userRepo, LeetcodeAPI)
	if authService == nil {
		log.Fatal("Failed to initialize AuthService")
	}

	// Define routes
	r := mux.NewRouter()

	authHandler := handlers.NewAuthHandler(authService)
	//userHandler := handlers.NewUserHandler(userService)
	//questionHandler := handlers.NewQuestionHandler(questionService)

	//r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

	routes.InitialiseAuthRouter(r, authHandler)
	//routes.InitialiseUserRouter(r, userHandler)
	//routes.InitialiseQuestionRouter(r, questionHandler)

	return r
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create the mux router
	router := createRouter()

	// Create the HTTP request from the API Gateway event
	req, err := http.NewRequestWithContext(ctx, request.HTTPMethod, request.Path, strings.NewReader(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Add headers from API Gateway event to the request
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	// Create an in-memory response recorder
	rr := httptest.NewRecorder()

	// Serve the HTTP request using the mux router
	router.ServeHTTP(rr, req)

	// Return the API Gateway response
	return events.APIGatewayProxyResponse{
		StatusCode: rr.Code,
		Body:       rr.Body.String(),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "*",
		},
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
