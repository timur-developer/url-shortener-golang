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
)

type recordToUpdate struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

func Patch(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Patch handler called")
		const op = "handlers.url.update"

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
		var record recordToUpdate

		err := render.DecodeJSON(r.Body, &record)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		updatedRecord, err := urlStorage.UpdateRecordPartly(alias, record.Alias, record.URL, user_id)
		if err != nil {
			if errors.Is(storage.ErrIdNotFound, err) {
				log.Error(fmt.Sprintf("failed to find record with alias '%v'", alias), sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find record with alias '%v'", alias)))
				return
			} else if errors.Is(storage.ErrWrongPatchRequest, err) {
				log.Error("using patch instead of put request to make complete update of record", sl.Err(err))
				render.JSON(w, r, resp.Error("using patch instead of put request to make complete update of record"))
				return
			}
			log.Error("failed to update record", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update record"))
			return
		}

		render.JSON(w, r, updatedRecord)
	}

}

func Put(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update"

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
		var record recordToUpdate

		err := render.DecodeJSON(r.Body, &record)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		updatedRecord, err := urlStorage.UpdateRecordCompletely(alias, record.Alias, record.URL, user_id)
		if err != nil {
			if errors.Is(storage.ErrIdNotFound, err) {
				log.Error(fmt.Sprintf("failed to find record with alias '%v'", alias), sl.Err(err))
				render.JSON(w, r, resp.Error(fmt.Sprintf("failed to find record with alias '%v'", alias)))
				return
			}
			log.Error("failed to update record", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update record"))
			return
		}

		render.JSON(w, r, updatedRecord)
	}
}
