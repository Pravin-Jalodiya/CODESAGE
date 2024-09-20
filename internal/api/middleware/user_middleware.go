package middleware

import (
	"cli-project/internal/config/roles"
	"net/http"
)

func UserRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userMetaData, ok := r.Context().Value("userMetaData").(UserMetaData)

		if !ok || userMetaData.Role != roles.USER {
			unauthorized(w, "You do not have permission to access this resource")
			return
		}

		next.ServeHTTP(w, r)
	})
}
