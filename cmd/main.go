package main

import (
	"log/slog"
	"os"
	"rest_urlshort_postgres/internal/config"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Debug("Launching app ", slog.String("env", cfg.Env))

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
