package routes

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"github.com/gorilla/mux"
)

func InitialiseUserRouter(r *mux.Router, userHandler *handlers.UserHandler) {
	authRouter := r.PathPrefix("/user").Subrouter()
	authRouter.Use(middleware.JWTAuthMiddleware)
	authRouter.Use(middleware.UserRoleMiddleware)
	authRouter.HandleFunc("/profile/{username}", userHandler.GetUserByID).Methods("GET")
	authRouter.HandleFunc("/progress", userHandler.GetUserProgress).Methods("GET")
}
