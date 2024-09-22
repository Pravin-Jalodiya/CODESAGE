package middleware_test

import (
	"cli-project/internal/api/middleware"
	"cli-project/internal/config/roles"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminRoleMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		userMetaData   middleware.UserMetaData
		addMetaData    bool
		expectedStatus int
	}{
		{
			name:           "Admin Role",
			userMetaData:   middleware.UserMetaData{Role: roles.ADMIN},
			addMetaData:    true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-Admin Role",
			userMetaData:   middleware.UserMetaData{Role: roles.USER},
			addMetaData:    true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No User MetaData",
			addMetaData:    false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.addMetaData {
				ctx := context.WithValue(req.Context(), "userMetaData", test.userMetaData)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			handler := middleware.AdminRoleMiddleware(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)
		})
	}
}
