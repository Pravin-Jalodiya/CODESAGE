package routes

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"github.com/gorilla/mux"
)

func InitialiseQuestionRouter(r *mux.Router, questionHandler *handlers.QuestionHandler) {
	r.HandleFunc("/question", questionHandler.GetQuestions).Methods("GET")
	questRouter := r.PathPrefix("/question").Subrouter()
	questRouter.Use(middleware.JWTAuthMiddleware)
	questRouter.Use(middleware.AdminRoleMiddleware)
	questRouter.HandleFunc("", questionHandler.RemoveQuestionById).Methods("DELETE")
}
