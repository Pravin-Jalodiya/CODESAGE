package routes

import (
	"cli-project/internal/api/handlers"
	"github.com/gorilla/mux"
)

func InitialiseAuthRouter(r *mux.Router, authHandler *handlers.AuthHandler) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authHandler.SignupHandler).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	authRouter.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
}
