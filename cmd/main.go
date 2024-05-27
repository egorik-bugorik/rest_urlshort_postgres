package main

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
	"rest_urlshort_postgres/internal/http-server/handlers/redirect"
	"rest_urlshort_postgres/internal/http-server/handlers/url/save"
	"rest_urlshort_postgres/internal/http-server/middleware/logger"
	"rest_urlshort_postgres/internal/lib/logger/handler/slogpretty"

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

	//:::::SETUP CHI
	router := chi.NewRouter()

	//*** CHI MIDLEWARE
	router.With(middleware.Logger)
	router.With(middleware.RequestID)
	router.With(middleware.Recoverer)
	router.With(middleware.URLFormat)
	router.With(logger.New(log))

	//:::: CHI ENDPOINTS
	router.Post("/", save.New(log, storage))
	router.Get("/{alias}", redirect.Redirect(log, storage))

	//	:::::CREATE SERVER

	log.Info("START server...", slog.String("HTTPaddress", cfg.HttpServer.Address))

	server := http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Error while starting server", sl.Error(err))

	}

	log.Error("Fail to server run")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvLocal:
		log = setupPrettySlog()
	case EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
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
