package url

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	middleware2 "go-to-do-checklist/internal/http-server/middleware"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

func Info(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.info"
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

		record, err := urlStorage.GetRecordByAlias(alias, user_id)
		if err != nil {
			log.Error("failed to get record", sl.Err(err))
			render.JSON(w, r, resp.Error(fmt.Sprintf("failed to get record for alias '%v'", alias)))
			return
		}

		render.JSON(w, r, record)
	}
}
