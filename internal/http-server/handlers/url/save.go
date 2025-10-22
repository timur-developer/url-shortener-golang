package url

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	middleware2 "go-to-do-checklist/internal/http-server/middleware"
	resp "go-to-do-checklist/internal/lib/api/response"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/lib/random"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
)

const aliasLength = 4

type SaveRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type SaveResponse struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlStorage URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.new"

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

		var req SaveRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = urlStorage.SaveURL(req.URL, alias, user_id)

		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		recordCreated, err := urlStorage.GetRecordByURL(req.URL, user_id)
		if err != nil {
			log.Error("failed to get record", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get record"))
			return
		}

		render.JSON(w, r, recordCreated)
	}
}
