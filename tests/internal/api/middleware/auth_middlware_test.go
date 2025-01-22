package middleware_test

import (
	"cli-project/internal/api/middleware"
	"cli-project/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testSecretKey = "mysecretkey"

func init() {
	config.SECRET_KEY = []byte(testSecretKey)

}

func TestJWTAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		claims         jwt.MapClaims
		expectedStatus int
	}{
		{
			name:           "No Authorization Header",
			claims:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed Authorization Header",
			claims:         jwt.MapClaims{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Token",
			claims:         jwt.MapClaims{"some": "claim"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Token",
			claims:         createValidClaims(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing userId",
			claims:         createClaimsWithout("userId"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid userId format",
			claims:         createClaimsWithInvalidUserID(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing role",
			claims:         createClaimsWithout("role"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid role value",
			claims:         createClaimsWithInvalidRole(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing username",
			claims:         createClaimsWithout("username"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.claims != nil {
				token := createTokenWithClaims(t, test.claims)
				req.Header.Set("Authorization", "Bearer "+token)
			}
			rr := httptest.NewRecorder()

			handler := middleware.JWTAuthMiddleware(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
		})
	}
}

func createTokenWithClaims(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.SECRET_KEY)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}
	return tokenString
}

func createValidClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"username": "testuser",
		"userId":   uuid.New().String(),
		"role":     "admin",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}
}

func createClaimsWithout(field string) jwt.MapClaims {
	claims := createValidClaims()
	delete(claims, field)
	return claims
}

func createClaimsWithInvalidUserID() jwt.MapClaims {
	claims := createValidClaims()
	claims["userId"] = "invalid-uuid"
	return claims
}

func createClaimsWithInvalidRole() jwt.MapClaims {
	claims := createValidClaims()
	claims["role"] = "invalid-role"
	return claims
}
