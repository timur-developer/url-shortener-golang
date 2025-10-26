package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go-to-do-checklist/internal/http-server/handlers/users"
	"go-to-do-checklist/internal/storage"
	"net/http"
	"strings"
)

type contextKey string

const (
	UsersIdKey   contextKey = "user_id"
	UsersNameKey contextKey = "username"
	UsersRoleKey contextKey = "user_role"
)

func AuthMiddleware(userStorage users.UserStorage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				ctx = context.WithValue(ctx, UsersIdKey, 0)
				ctx = context.WithValue(ctx, UsersRoleKey, "anonymous")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			b, err := base64.StdEncoding.DecodeString(strings.Split(authHeader, " ")[1])
			if err != nil {
				http.Error(w, "fail to decode string", http.StatusInternalServerError)
				return
			}
			data := (strings.Split(string(b), ":"))
			username, password := data[0], data[1]
			userGotten, err := userStorage.AuthenticateUser(username, password)
			if err != nil {
				if errors.Is(storage.ErrUserNotFound, err) {
					http.Error(w, fmt.Sprintf("fail to found user '%v'", username), http.StatusNotFound)
					return
				} else if errors.Is(storage.ErrIncorrectPassowrd, err) {
					http.Error(w, storage.ErrIncorrectPassowrd.Error(), http.StatusNotFound)
					return
				} else {
					http.Error(w, fmt.Sprintf("fail to authenticate user '%v'", username), http.StatusNotFound)
					return
				}
			}
			ctx = context.WithValue(ctx, UsersIdKey, userGotten.ID)
			ctx = context.WithValue(ctx, UsersNameKey, userGotten.Username)
			ctx = context.WithValue(ctx, UsersRoleKey, userGotten.UserRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
