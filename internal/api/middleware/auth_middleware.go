package middleware

import (
	"cli-project/internal/config"
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

// Middleware to validate JWT
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w, "Missing Authorization header")
			return
		}

		// Extract the token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 || tokenParts[1] == "" {
			unauthorized(w, "Missing token in Authorization header")
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
			unauthorized(w, "Invalid or expired token")
			return
		}

		// Extract the claims from the token {userId, role}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			unauthorized(w, "Invalid Token")
			return
		}

		// Extract the userId from the token claims
		userIdStr, ok := claims["userId"].(string)
		if !ok {
			unauthorized(w, "Invalid Token")
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			unauthorized(w, "Invalid Token")
			return
		}

		// Extract the role from the token claims
		role, ok := claims["role"].(string)
		if !ok {
			unauthorized(w, "Invalid Token")
			return
		}

		// Extract the username from the token claims
		username, ok := claims["username"].(string)
		if !ok {
			unauthorized(w, "Invalid Token")
			return
		}

		// Create user metadata
		userMetaData := struct {
			Username string
			UserId   uuid.UUID
			Role     string
		}{
			Username: username,
			UserId:   userId,
			Role:     role,
		}

		// Attach user metadata to the request context
		ctx := context.WithValue(r.Context(), "userMetaData", userMetaData)
		r = r.WithContext(ctx)

		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// Helper to return unauthorized error response
func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	jsonResponse := map[string]string{
		"code":    "401",
		"message": message,
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
