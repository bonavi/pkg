package middleware

import (
	"encoding/json"
	"net/http"
	"pkg/errors"
	"pkg/log"

	"github.com/gofiber/fiber/v2"
)

type versionRes struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	BuildDate string `json:"buildDate"`
	Hostname  string `json:"hostname"`
}

func NewVersionHandler(version, build, buildDate, hostname string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Сериализуем ответ в JSON
		message, _ := json.Marshal(versionRes{
			Version:   version,
			Build:     build,
			BuildDate: buildDate,
			Hostname:  hostname,
		})

		// Пишем 200 ответ
		w.WriteHeader(fiber.StatusOK)

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Пишем ответ в HTTP-ответ
		if _, err := w.Write(message); err != nil {
			log.Error(errors.InternalServer.Wrap(err))
		}
	}
}
