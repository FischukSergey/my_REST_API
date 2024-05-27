package save

import (
	"dev/myrestapi/internal/logger/handlers/sl"
	resp "dev/myrestapi/lib/response"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Genre string `json:"genre"`
	Name  string `json:"name"` //validate:"required,name"` //TODO: вставить проверку валидности
	Size  string `json:"size"`
}

type RespSave struct {
	resp.Response
	Name string
}
type PictureSaver interface {
	SavePicture(r *Request) (int64, error)
}

func New(log *slog.Logger, pictureSaver PictureSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//TODO вставить проверку на валидность имени

		id, err := pictureSaver.SavePicture(&req) //реализуем метод из sqlite.go получаем назад id записи

		//TODO вставить проверку на уже существующее имя
		/*if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("name", req.name))

			render.JSON(w, r, resp.Error("name already exists"))

			return
		}
		*/

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add picture"))

			return
		}
		log.Info("picture added", slog.Int64("id", id))

		responseOK(w, r, req.Name)

	}
}
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, RespSave{
		Response: resp.OK(),
		Name:     alias,
	})
}
