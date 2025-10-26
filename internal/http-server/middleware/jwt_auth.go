package middleware

import (
	"context"
	"go-to-do-checklist/internal/lib/jwt"
	"net/http"
	"strings"
)

func JWTAuthMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID := 0
			username := ""
			role := "anonymous"

			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
					return
				}
				tokenString := parts[1]
				if claims, err := jwt.ValidateToken(tokenString); err == nil {
					userID = int(claims["user_id"].(float64))
					username = claims["username"].(string)
					role = claims["role"].(string)
				} else {
					http.Error(w, "fail to validate token", http.StatusInternalServerError)
					return
				}
			}

			ctx = context.WithValue(ctx, UsersIdKey, userID)
			ctx = context.WithValue(ctx, UsersNameKey, username)
			ctx = context.WithValue(ctx, UsersRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
