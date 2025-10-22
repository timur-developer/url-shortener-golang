package users

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
)

func Get(log *slog.Logger, userStorage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.get.all"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var allUsers *[]storage.User

		allUsers, err := userStorage.GetAllUsers()
		if err != nil {
			log.Error("failed to get urls", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get urls"))
			return
		}

		log.Info("users gotten")
		render.JSON(w, r, allUsers)
	}
}
