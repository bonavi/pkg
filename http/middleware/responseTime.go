package middleware

import (
	"net/http"
	"pkg/http/responseWriter"
	"pkg/log"
	"strings"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Оборачиваем Writer в кастомный для получения статус-кода
		writerWithStatus := &responseWriter.ResponseWriterWithStatus{
			ResponseWriter: w,
			Status:         nil,
		}
		
		// Сохраняем статус ответа после выполнения следующего в стеке хандлера
		next.ServeHTTP(writerWithStatus, r)

		if writerWithStatus.Status == nil {
			log.Error("Cannot get status code from handler")
			return
		}
	})
}

// preparePath заменяет пути /rtb/ssp на rtb_ssp
func preparePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return strings.ReplaceAll(path, "/", "_")
}
