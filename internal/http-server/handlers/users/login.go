package users

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/jwt"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(log *slog.Logger, userStorage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.login"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req LoginRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		user, err := userStorage.AuthenticateUser(req.Username, req.Password)
		if err != nil {
			if errors.Is(storage.ErrUserNotFound, err) {
				log.Error(fmt.Sprintf("failed to find user with username %v", req.Username), sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find user with username %v", req.Username)))
				return
			} else if errors.Is(storage.ErrIncorrectPassowrd, err) {
				log.Error(fmt.Sprintf("incorrect password for user %v", req.Username), sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("incorrect password for user %v", req.Username)))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		tokens, err := jwt.GenerateTokenPair(user.ID, user.Username, user.UserRole)
		if err != nil {
			log.Error("failed to generate tokens", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to generate tokens"))
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status": "login succesful",
			"tokens": tokens,
		})
	}
}
