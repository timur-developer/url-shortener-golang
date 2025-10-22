package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go-to-do-checklist/internal/config"
	"go-to-do-checklist/internal/http-server/handlers/url"
	"go-to-do-checklist/internal/http-server/handlers/users"
	middleware2 "go-to-do-checklist/internal/http-server/middleware"
	mwLogger "go-to-do-checklist/internal/http-server/middleware/logger"
	"go-to-do-checklist/internal/lib/logger/handlers/slogpretty"
	"go-to-do-checklist/internal/lib/logger/sl"
	"go-to-do-checklist/internal/storage"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware2.AuthMiddleware(storage))

	// Публичные роуты для всех

	router.Get("/{alias}", url.Redirect(log, storage))
	router.Post("/url", url.New(log, storage))
	router.Post("/register", users.Register(cfg, log, storage))

	// Защищенные роуты только для авторизованных
	router.With(middleware2.RequireAuth).Get("/url", url.Get(log, storage))
	router.With(middleware2.RequireAuth).Get("/url/{alias}", url.Info(log, storage))
	router.With(middleware2.RequireAuth).Delete("/{alias}", url.Delete(log, storage))
	router.With(middleware2.RequireAuth).Patch("/{alias}", url.Patch(log, storage))
	router.With(middleware2.RequireAuth).Put("/{alias}", url.Put(log, storage))

	// Роуты только для админов
	router.With(middleware2.RequireAdmin).Get("/users", users.Get(log, storage))
	router.With(middleware2.RequireAdmin).Delete("/users/{username}", users.Delete(log, storage))
	router.With(middleware2.RequireAdmin).Get("/admin/url", url.AdminAccess(log, storage))
	router.With(middleware2.RequireAdmin).Get("/{alias}/{user_id}", url.AdminRedirect(log, storage, storage))
	router.With(middleware2.RequireAdmin).Delete("/{alias}/{user_id}", url.AdminDelete(log, storage))
	
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr: cfg.Address,

		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
