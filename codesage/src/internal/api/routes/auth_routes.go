package routes

import (
	"codesage/internal/api/handlers"
	"codesage/internal/api/middleware"
	"github.com/gorilla/mux"
)

func InitialiseAuthRouter(r *mux.Router, authHandler *handlers.AuthHandler) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authMemberRouter := r.PathPrefix("/auth/member").Subrouter()
	authMemberRouter.Use(middleware.JWTAuthMiddleware, middleware.MemeberRoleMiddleware)
	//authMemberRouter.HandleFunc("/role", authHandler.GetRole).Methods("GET")
	authRouter.HandleFunc("/signup", authHandler.SignupHandler).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	//authRouter.HandleFunc("/forgot-password", authHandler.ForgotPasswordHandler).Methods("POST")
	//authRouter.HandleFunc("/reset-password", authHandler.ResetPasswordHandler).Methods("POST")
	authMemberRouter.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
}
