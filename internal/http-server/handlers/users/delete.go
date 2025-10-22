package users

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
)

func Delete(log *slog.Logger, userStorage UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		username := chi.URLParam(r, "username")

		if err := userStorage.DeleteUser(username); err != nil {
			if errors.Is(storage.ErrUserNotFound, err) {
				log.Error(fmt.Sprintf("fail to find user '%v'", username))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find user '%v'", username)))
				return
			}
			log.Error("failed to delete user")
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to delete user '%v'", username)))
			return
		}

		w.WriteHeader(204)
	}
}
