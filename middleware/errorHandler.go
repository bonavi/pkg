package middleware

import (
	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
)

var DefaultErrorHandler = func(ctx *fiber.Ctx, err error) error {

	// Если ошибки нет, а мы сюда попали, значит какие-то проблемы в чейне вызовов HTTP-сервера
	if err == nil {
		err = errors.Default.New("В функцию DefaultErrorHandler передана пустая ошибка").SkipThisCall()
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		if _, fiberWSErr := ctx.Status(fiberError.Code).WriteString(fiberError.Error()); fiberWSErr != nil {
			log.Error(errors.Default.Wrap(fiberWSErr))
		}
		return nil
	}

	// Кастуем ошибку, если она не обернута, внутри оборачиваем ее в Default,
	// на вызовы начнет алертить система оповещения и разработчик быстро пофиксит проблему
	customErr := errors.CastError(err)

	// Проставляем HumanText по дефолту, если его не определили при создании ошибки
	if customErr.HumanText != "" {
		customErr.HumanText = customErr.ErrorType.HumanText
	}
	customErr.DeveloperText = customErr.Err.Error()

	customErr.SystemInfo = log.GetSystemInfo()

	// Логгируем ошибку
	log.LogError(customErr)

	// Возвращаем ошибку согласно статусу
	return ctx.Status(customErr.ErrorType.HTTPCode).JSON(customErr)
}
