package main

import (
	"log/slog"
	"os"
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

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Helo")
	log.Debug("Launching app ", slog.String("env", cfg.Env))

	err, storage := sql.New(cfg.StoragePath)
	if err != nil {
		log.Error("error while new storage :::", sl.Error(err))
		os.Exit(1)
	}
	err = storage.Ping()
	if err != nil {
		log.Error("error while ping db :::  ", err)
	}

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
