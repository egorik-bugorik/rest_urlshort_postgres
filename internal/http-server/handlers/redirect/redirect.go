package redirect

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"rest_urlshort_postgres/internal/lib/api/resp"
	"rest_urlshort_postgres/internal/lib/logger/sl"
	"rest_urlshort_postgres/internal/storage"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func Redirect(log *slog.Logger, s UrlGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.redirect.Redirect"

		log = log.With(slog.String("op", op), slog.String("request id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {

			log.Error("Error to get alias for redirect")

			render.JSON(w, r, resp.Error("invalid requesy"))

			return
		}

		url, err := s.GetUrl(alias)
		if err != nil {

			if err == storage.ErrUrlNotFound {

				log.Error("Url not found", slog.String("alias", alias))

				render.JSON(w, r, resp.Error("Url not found"))

				return

			}

			log.Error("coudn't get url", sl.Error(err))

			render.JSON(w, r, resp.Error("Internal error "))

			return

		}

		log.Info("Gotta url", "url", url)

		http.Redirect(w, r, url, http.StatusFound)
	}

}
