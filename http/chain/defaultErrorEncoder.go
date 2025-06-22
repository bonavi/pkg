package chain

import (
	"context"
	"encoding/json"
	"net/http"
	"pkg/errors"

	"pkg/log"
)

func DefaultErrorEncoder(ctx context.Context, w http.ResponseWriter, er error) {
	_, span := tracer.Start(ctx, "DefaultErrorEncoder")
	defer span.End()

	// Проверяем, что мы сюда попали из-за ошибки
	if er == nil {
		er = errors.Default.New("В функцию DefaultErrorEncoder передана пустая ошибка").SkipThisCall()
	}

	// Кастуем пришедшую ошибку
	err := errors.CastError(er)

	log.LogError(err)

	// Прописываем тип контента, который будем отправлять клиенту
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Прописываем HTTP-код
	w.WriteHeader(err.ErrorType.HTTPCode)

	// Сериализуем ошибку
	byt, er := json.Marshal(err)
	if er != nil {
		log.Error(er)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(er.Error()))
	}

	// Пишем ошибку
	if _, writeErr := w.Write(byt); writeErr != nil {
		log.Error(errors.Default.Wrap(writeErr))
	}
}
