package middleware

import (
	"cli-project/internal/config/roles"
	errs "cli-project/pkg/errors"
	"net/http"
)

func UserRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userMetaData, ok := r.Context().Value("userMetaData").(UserMetaData)

		if !ok || userMetaData.Role != roles.USER || userMetaData.BanState == true {
			unauthorized(w, "Unauthorized access", errs.CodePermissionDenied)
			return
		}

		next.ServeHTTP(w, r)
	})
}
