package routes

//import (
//	"codesage/internal/api/handlers"
//	"codesage/internal/api/middleware"
//	"github.com/gorilla/mux"
//)

//func InitialiseUserRouter(r *mux.Router, userHandler *handlers.UserHandler) {
//	userRouter := r.PathPrefix("/users").Subrouter()
//	adminRouter := r.PathPrefix("/").Subrouter()
//	memberRouter := r.PathPrefix("/users").Subrouter()
//	userRouter.Use(middleware.JWTAuthMiddleware)
//	adminRouter.Use(middleware.JWTAuthMiddleware)
//	memberRouter.Use(middleware.JWTAuthMiddleware)
//	userRouter.Use(middleware.UserRoleMiddleware)
//	adminRouter.Use(middleware.AdminRoleMiddleware)
//	memberRouter.Use(middleware.MemeberRoleMiddleware)
//memberRouter.HandleFunc("/profile/{username}", userHandler.GetUserByID).Methods("GET")
//userRouter.HandleFunc("/progress/{username}", userHandler.GetUserProgress).Methods("GET")
//userRouter.HandleFunc("/progress-list", userHandler.GetProgressList).Methods("GET")
//userRouter.HandleFunc("/progress/{username}", userHandler.UpdateUserProgress).Methods("PATCH")
//memberRouter.HandleFunc("/update-profile", userHandler.UpdateUserProfile).Methods("PATCH")
//adminRouter.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
//adminRouter.HandleFunc("/platform-stats", userHandler.GetPlatformStats).Methods("GET")
//adminRouter.HandleFunc("/users/update-user-ban-state", userHandler.UpdateUserBanState).Methods("PATCH")
//adminRouter.HandleFunc("/users/delete", userHandler.DeleteUser).Methods("DELETE")
//}
