package middleware

import (
	"cli-project/internal/config"
	"cli-project/internal/config/roles"
	errs "cli-project/pkg/errors"
	"cli-project/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type UserMetaData struct {
	Username string
	UserId   uuid.UUID
	Role     roles.Role
	BanState bool
}

// Middleware to validate JWT
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Logger.Errorw("Missing Authorization header", "error", nil, "time", time.Now())
			unauthorized(w, "Missing Authorization header", errs.CodeInvalidRequest)
			return
		}

		// Extract the token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 || tokenParts[1] == "" {
			logger.Logger.Errorw("Missing token in Authorization header", "tokenParts", tokenParts, "time", time.Now())
			unauthorized(w, "Missing token in Authorization header", errs.CodeInvalidRequest)
			return
		}
		tokenString := tokenParts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return config.SECRET_KEY, nil
		})

		// Handle token validation errors
		if err != nil || !token.Valid {
			logger.Logger.Errorw("Invalid token", "error", err, "time", time.Now())
			unauthorized(w, "Invalid or expired token", errs.CodePermissionDenied)
			return
		}

		// Extract the claims from the token {userId, role}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.Logger.Errorw("Invalid Token Claims", "token", tokenString, "claims", claims, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		// Extract the userId from the token claims
		userIdStr, ok := claims["userId"].(string)
		if !ok {
			logger.Logger.Errorw("userId not found in token claims", "claims", claims, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			logger.Logger.Errorw("Invalid userId format", "userId", userIdStr, "error", err, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		// Extract the role from the token claims
		roleStr, ok := claims["role"].(string)
		if !ok {
			logger.Logger.Errorw("role not found in token claims", "claims", claims, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		role, err := roles.ParseRole(roleStr)
		if err != nil {
			logger.Logger.Errorw("Invalid role value", "role", roleStr, "error", err, "time", time.Now())
			unauthorized(w, "Invalid Role", errs.CodeValidationError)
			return
		}

		// Extract the username from the token claims
		username, ok := claims["username"].(string)
		if !ok {
			logger.Logger.Errorw("username not found in token claims", "claims", claims, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		// Extract the banState from the token claims
		banState, ok := claims["banState"].(bool)
		if !ok {
			logger.Logger.Errorw("banState not found in token claims", "claims", claims, "time", time.Now())
			unauthorized(w, "Invalid Token", errs.CodeInvalidRequest)
			return
		}

		// Create user metadata
		userMetaData := UserMetaData{
			Username: username,
			UserId:   userId,
			Role:     role,
			BanState: banState,
		}

		// Attach user metadata to the request context
		ctx := context.WithValue(r.Context(), "userMetaData", userMetaData)
		r = r.WithContext(ctx)

		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// Helper to return unauthorized error response
func unauthorized(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	jsonResponse := map[string]interface{}{
		"error_code": code,
		"message":    message,
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
