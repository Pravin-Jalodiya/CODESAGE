package routes

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"github.com/gorilla/mux"
)

func InitialiseQuestionRouter(r *mux.Router, questionHandler *handlers.QuestionHandler) {
	memberQuestionRouter := r.PathPrefix("/questions").Subrouter()
	memberQuestionRouter.Use(middleware.JWTAuthMiddleware)
	memberQuestionRouter.Use(middleware.MemeberRoleMiddleware)
	memberQuestionRouter.HandleFunc("", questionHandler.GetQuestions).Methods("GET")
	//get a specific question
	questionsRouter := r.PathPrefix("/questions").Subrouter()
	questionsRouter.Use(middleware.JWTAuthMiddleware)
	questionsRouter.Use(middleware.AdminRoleMiddleware)
	questionRouter := r.PathPrefix("/question").Subrouter()
	questionRouter.Use(middleware.JWTAuthMiddleware)
	questionRouter.Use(middleware.AdminRoleMiddleware)
	questionRouter.HandleFunc("", questionHandler.RemoveQuestionById).Methods("DELETE")
	questionsRouter.HandleFunc("", questionHandler.AddQuestions).Methods("POST")
}
