package url

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go-to-do-checklist/internal/http-server/handlers/users"
	middleware2 "go-to-do-checklist/internal/http-server/middleware"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

func Redirect(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect"

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
		log.Info("redirect attempt", slog.Int("user_id", user_id), slog.String("alias", alias))
		url, err := urlStorage.GetURL(alias, user_id)
		if err != nil {
			if errors.Is(storage.ErrAliasNotFound, err) {
				log.Error("alias not found", sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to found record with alias '%v'", alias)))
				return
			}
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to get url for alias '%v'", alias)))
			return
		}

		http.Redirect(w, r, url, 301)
	}
}

func AdminRedirect(log *slog.Logger, urlStorage URLStorage, userStorage users.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.admin"

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
		
		log.Info("redirect attempt", slog.String("alias", alias))
		url, err := urlStorage.GetURL(alias, user_id)
		if err != nil {
			if errors.Is(storage.ErrAliasNotFound, err) {
				log.Error("alias not found", sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to found record with alias '%v'", alias)))
				return
			}
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to get url for alias '%v'", alias)))
			return
		}

		http.Redirect(w, r, url, 301)
	}
}
