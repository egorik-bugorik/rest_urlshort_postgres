package save

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"rest_urlshort_postgres/internal/lib/api/resp"
	"rest_urlshort_postgres/internal/lib/logger/sl"
	"rest_urlshort_postgres/internal/lib/random"
	"rest_urlshort_postgres/internal/storage"
)

type UrlSaver interface {
	SaveUrl(urlTYoSave string, alias string) (int64, error)
}

type Request struct {
	Alias string `json:"alias,omitempty"`
	Url   string `json:"url" validate:"required,url"`
}

func New(log *slog.Logger, s UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.saveUrl.New"

		log = log.With(slog.String("op", op), slog.String("request id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Error("request body is empty!!!")

			render.JSON(w, r, resp.Error("request is empty!!!"))

			return
		}
		if err != nil {

			log.Error("error while decode requset")

			render.JSON(w, r, resp.Error("request is invalid"))

			return
		}

		log.Info("reqest body decoede", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validErrerr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Error(err))

			render.JSON(w, r, resp.ValidateError(validErrerr))

			return

		}

		alias := req.Alias
		if alias == "" {

			alias = random.NewAlias(6)

		}

		id, err := s.SaveUrl(req.Url, alias)

		if err == storage.ErrUrlExist {

			log.Error("Url already exist ", slog.String("url", req.Url))

			render.JSON(w, r, resp.Error("url already exists!"))

			return

		}
		if err != nil {

			log.Error("Fail to save url", sl.Error(err))

			render.JSON(w, r, resp.Error("fail to save url"))

			return

		}
		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, resp.OK(alias))
	}
}
