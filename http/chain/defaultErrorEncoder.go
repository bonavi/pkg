package chain

import (
	"context"
	"net/http"

	"pkg/errors"
	"pkg/log"
)

func DefaultErrorEncoder(ctx context.Context, w http.ResponseWriter, er error) {
	_, span := tracer.Start(ctx, "DefaultErrorEncoder")
	defer span.End()

	// Проверяем, что мы сюда попали из-за ошибки
	if er == nil {
		er = errors.InternalServer.New(ctx, "В функцию DefaultErrorEncoder передана пустая ошибка",
			errors.SkipThisCallOption(),
		)
	}

	// Кастуем пришедшую ошибку
	err := errors.CastError(ctx, er)

	// Если человекочитаемый текст не написан, заполняем шаблонным
	if err.HumanText == "" {
		err.HumanText = humanTextByLevel[err.ErrorType]
	}

	log.LogError(ctx, err)

	// Прописываем тип контента, который будем отправлять клиенту
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Прописываем HTTP-код
	w.WriteHeader(err.ErrorType.HTTPCode())

	// Сериализуем ошибку
	byt, er := errors.JSON(err)
	if er != nil {
		log.Error(ctx, er)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(er.Error()))
	}

	// Пишем ошибку
	if _, writeErr := w.Write(byt); writeErr != nil {
		log.Error(ctx, errors.InternalServer.Wrap(ctx, writeErr))
	}
}

// Сопоставление типа ошибки и дефолтной человекочитаемой ошибки
var humanTextByLevel = map[errors.ErrorType]string{
	errors.BadRequest:     "Введены неверные данные",
	errors.InternalServer: "Произошла непредвиденная ошибка",
	errors.NotFound:       "Данные не найдены",
	errors.Forbidden:      "Доступ запрещен",
	errors.Teapot:         "Разработчик забыл написать текст ошибки",
	errors.BadGateway:     "Произошла ошибка на сервере внешнего сервиса",
	errors.Unauthorized:   "Пользователь не авторизован",
	errors.Timeout:        "Клиент отказался принимать данные",
}
