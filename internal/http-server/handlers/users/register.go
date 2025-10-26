package users

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-to-do-checklist/internal/config"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/jwt"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AdminPass string `json:"admin_pass,omitempty"`
}

func Register(cfg *config.Config, log *slog.Logger, userStorage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.register.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RegisterRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		role := "user"
		fmt.Println(cfg.AdminPass, req.AdminPass)
		if req.AdminPass == cfg.AdminPass {
			role = "admin"
		}

		if err := userStorage.SaveUser(req.Username, req.Password, role); err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Info("user already exists", slog.String("url", req.Username))
				render.JSON(w, r, resp.Error("user already exists"))
				return
			}
			log.Error("failed to add user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}

		record, err := userStorage.GetRecordByUsername(req.Username)
		if err != nil {
			log.Error("failed to get record", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get record"))
			return
		}

		tokens, err := jwt.GenerateTokenPair(record.ID, record.Username, record.UserRole)
		if err != nil {
			log.Error("failed to generate tokens", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to generate tokens"))
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status": "user registered successfully",
			"tokens": tokens,
		})
	}
}
