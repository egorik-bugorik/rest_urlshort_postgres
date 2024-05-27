package logger

import (
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		sloger := log.With(slog.String("component", "middleware/logegr"))

		sloger.Info("middleware logger is enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {

			entry := sloger.With(
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
				slog.String("method", r.Method),
				slog.String("urlPath", r.URL.Path),
				slog.String("request ID", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			tim := time.Now()
			defer func() {
				entry.Info(
					"request completed",
					slog.Int("bytes", ww.BytesWritten()),
					slog.Int("status", ww.Status()),
					slog.String("duration", time.Since(tim).String()),
				)

			}()
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}
}
