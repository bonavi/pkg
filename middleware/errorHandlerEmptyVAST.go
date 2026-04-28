package middleware

import (
	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
	"pkg/vast"
)

var ErrorHandlerEmptyVAST = func(c *fiber.Ctx, err error) error {

	// Если ошибки нет, а мы сюда попали, значит какие-то проблемы в чейне вызовов HTTP-сервера
	if err == nil {
		err = errors.Default.New("В функцию ExchangeErrorHandler передана пустая ошибка").SkipThisCall()
	}

	// Кастуем ошибку, если она не обернута, внутри оборачиваем ее в Default,
	// на вызовы начнет алертить система оповещения и разработчик быстро пофиксит проблему
	customErr := errors.CastError(err)

	// Логгируем ошибку
	log.LogError(customErr)

	c.Set("Content-Type", "text/xml")

	// Отправляем VAST-ответ и 200 код
	if err = c.Status(fiber.StatusOK).Send([]byte(vast.EmptyVAST)); err != nil {
		log.Error(errors.Default.Wrap(err))
	}

	// И игнорируем ошибку для fiber, чтобы он ее не писал в
	return nil
}
