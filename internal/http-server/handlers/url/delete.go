package url

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	middleware2 "go-to-do-checklist/internal/http-server/middleware"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

func Delete(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete"
		userIDValue := r.Context().Value(middleware2.UsersIdKey)
		if userIDValue == nil {
			log.Error("user_id not found in context")
			render.JSON(w, r, resp.Error("authentication required"))
			return
		}

		user_id, ok := userIDValue.(int)
		if !ok {
			log.Error("user_id has wrong type", slog.Any("type", fmt.Sprintf("%T", userIDValue)))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")

		if err := urlStorage.DeleteURL(alias, user_id); err != nil {
			if errors.Is(storage.ErrAliasNotFound, err) {
				log.Error(fmt.Sprintf("fail to find record with alias '%v'", alias))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find record with alias '%v'", alias)))
				return
			}
			log.Error("failed to delete record")
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to delete record with alias '%v'", alias)))
			return
		}

		w.WriteHeader(204)
	}
}

func AdminDelete(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.admin"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
		if err != nil {
			log.Error("fail to convert string to int", sl.Err(err))
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to convert string to int '%v'", user_id)))
			return
		}

		if err := urlStorage.DeleteURL(alias, user_id); err != nil {
			if errors.Is(storage.ErrAliasNotFound, err) {
				log.Error(fmt.Sprintf("fail to find record with alias '%v'", alias))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find record with alias '%v'", alias)))
				return
			}
			log.Error("failed to delete record")
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to delete record with alias '%v'", alias)))
			return
		}

		w.WriteHeader(204)
	}
}
