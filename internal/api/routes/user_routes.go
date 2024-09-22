package routes

import (
	"cli-project/internal/api/handlers"
	"cli-project/internal/api/middleware"
	"github.com/gorilla/mux"
)

func InitialiseUserRouter(r *mux.Router, userHandler *handlers.UserHandler) {
	userRouter := r.PathPrefix("/users").Subrouter()
	adminRouter := r.PathPrefix("/").Subrouter()
	userRouter.Use(middleware.JWTAuthMiddleware)
	adminRouter.Use(middleware.JWTAuthMiddleware)
	userRouter.Use(middleware.UserRoleMiddleware)
	adminRouter.Use(middleware.AdminRoleMiddleware)
	userRouter.HandleFunc("/profile/{username}", userHandler.GetUserByID).Methods("GET")
	userRouter.HandleFunc("/progress/{username}", userHandler.GetUserProgress).Methods("GET")
	userRouter.HandleFunc("/progress/{username}", userHandler.UpdateUserProgress).Methods("PATCH")
	adminRouter.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	adminRouter.HandleFunc("/platform-stats", userHandler.GetPlatformStats).Methods("GET")
	adminRouter.HandleFunc("/users/update-user-ban-state", userHandler.UpdateUserBanState).Methods("PATCH")
}
