package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a next handler that writes 200
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Generate a JWT exactly as util.VarifyJwt expects
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "u1",
		"role":    float64(roles.Employee),
	})
	tokenString, _ := token.SignedString(util.JwtSecret) // use your util secret

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{"missing token", "", http.StatusUnauthorized},
		{"invalid format", "Token abc", http.StatusUnauthorized},
		{"invalid token", "Bearer invalidtoken", http.StatusUnauthorized},
		{"valid token", "Bearer " + tokenString, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler := AuthMiddleWare(nextHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}


func TestRequireRole(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name         string
		contextRole  string
		requiredRole string
		expectedCode int
	}{
		{"matching role", "Admin", "Admin", http.StatusOK},
		{"non-matching role", "Employee", "Manager", http.StatusUnauthorized},
		{"no role in context", "", "Admin", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Add role to context if provided
			ctx := req.Context()
			if tt.contextRole != "" {
				ctx = context.WithValue(ctx, "role", tt.contextRole)
			}
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := RequireRole(tt.requiredRole)(nextHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}