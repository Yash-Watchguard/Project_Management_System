package middleware

import (
	"context"
	"net/http"
	"strings"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			response.ErrorResponse(w, http.StatusUnauthorized, "Missing token", 1002)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token format", 1002)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		token, err := util.VarifyJwt(tokenString)
		if err != nil {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}

		roleFloat, ok := claims["role"].(float64)
		if !ok {
			response.ErrorResponse(w, http.StatusUnauthorized, "Invalid token", 1002)
			return
		}
		role := roles.Role(int(roleFloat)) // convert back ✅

		ctx := context.WithValue(r.Context(), ContextKey.UserId, userID)
		ctx = context.WithValue(ctx, ContextKey.UserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(string)
			if !ok || userRole != role {
				response.ErrorResponse(w, http.StatusUnauthorized, "Invalid role", 1007)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
