package middleware

import (
	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
)

var ErrorHandler204Code = func(ctx *fiber.Ctx, err error) error {

	// Если ошибки нет, а мы сюда попали, значит какие-то проблемы в чейне вызовов HTTP-сервера
	if err == nil {
		err = errors.Default.New("В функцию ExchangeErrorHandler передана пустая ошибка").SkipThisCall()
	}

	// Кастуем ошибку, если она не обернута, внутри оборачиваем ее в Default,
	// на вызовы начнет алертить система оповещения и разработчик быстро пофиксит проблему
	customErr := errors.CastError(err)

	// Логгируем ошибку
	log.LogError(customErr)

	// Для вызывающей SSP возвращаем 204 HTTP-код
	ctx.Status(fiber.StatusNoContent)

	// И игнорируем ошибку для fiber, чтобы он ее не писал в
	return nil
}
