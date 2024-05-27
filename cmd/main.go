package main

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"os"
	"rest_urlshort_postgres/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"rest_urlshort_postgres/internal/config"
	"rest_urlshort_postgres/internal/lib/logger/sl"

	"rest_urlshort_postgres/internal/storage/sql"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {

	//::::::  CONFIG LOAD
	cfg := config.MustLoad()

	//::::::::LOGGER  SETUP
	log := setupLogger(cfg.Env)

	log.Info("Helo")
	log.Debug("Launching app ", slog.String("env", cfg.Env))

	//:::::::STORAGE SETUP
	err, storage := sql.New(cfg.StoragePath)
	if err != nil {
		log.Error("error while new storage :::", sl.Error(err))
		os.Exit(1)
	}
	err = storage.Ping()
	if err != nil {
		log.Error("error while ping db :::  ", err)
	}

	router := chi.NewRouter()

	router.With(middleware.Logger)
	router.With(middleware.RequestID)
	router.With(middleware.Recoverer)
	router.With(middleware.URLFormat)

	router.With(logger.New(log))
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log

}
