package users

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/jwt"
	"go-to-do-checklist/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func Refresh(log *slog.Logger, userStorage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.refresh"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RefreshRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		claims, err := jwt.ValidateToken(req.RefreshToken)
		if err != nil {
			log.Error("invalid token", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid token"))
			return
		}

		if claims["type"] != "refresh" {
			log.Error("not a refresh token")
			render.JSON(w, r, resp.Error("not a refresh token"))
			return
		}

		username := claims["username"].(string)

		user, err := userStorage.GetRecordByUsername(username)
		if err != nil {
			log.Error("fail to find user", sl.Err(err))
			render.JSON(w, r, resp.Error("fail to find user"))
			return
		}

		newTokens, err := jwt.GenerateTokenPair(user.ID, user.Username, user.UserRole)
		if err != nil {
			log.Error("fail to generate tokens", sl.Err(err))
			render.JSON(w, r, resp.Error("fail to generate tokens"))
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status": "tokens refreshed",
			"tokens": newTokens,
		})
	}
}
