package middleware

import (
	"cli-project/internal/config/roles"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/utils"
	"net/http"
)

func AdminRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userMetaData, ok := r.Context().Value("userMetaData").(UserMetaData)
		if !ok || userMetaData.Role != roles.ADMIN {
			utils.Unauthorized(w, "Unauthorized access", errs.CodePermissionDenied)
			return
		}
		next.ServeHTTP(w, r)
	})
}
