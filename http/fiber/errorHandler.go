package fiber

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
)

var DefaultErrorHandler = func(ctx *fiber.Ctx, err error) error {

	// Если ошибки нет, а мы сюда попали, значит какие-то проблемы в чейне вызовов HTTP-сервера
	if err == nil {
		err = errors.InternalServer.New("В функцию DefaultErrorHandler передана пустая ошибка", []errors.Option{
			errors.SkipThisCallOption(),
		}...)
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		if _, fiberWSErr := ctx.Status(fiberError.Code).WriteString(fiberError.Error()); fiberWSErr != nil {
			log.Error(ctx.Context(), errors.InternalServer.Wrap(fiberWSErr))
		}
		return nil
	}

	// Кастуем ошибку, если она не обернута, внутри оборачиваем ее в InternalServer,
	// на вызовы начнет алертить система оповещения и разработчик быстро пофиксит проблему
	customErr := errors.CastError(err)

	// Проставляем HumanText по дефолту, если его не определили при создании ошибки
	if customErr.HumanText != "" {
		customErr.HumanText = errors.HumanTextByLevel[customErr.ErrorType]
	}
	customErr.DeveloperText = customErr.Err.Error()

	// Получаем идентификатор запроса
	requestID := ctx.GetRespHeader("X-Request-Id")
	if requestID != "" {
		if customErr.Params == nil {
			customErr.Params = make(map[string]string, 1)
		}
		customErr.Params["request_id"] = requestID
	}

	customErr.SystemInfo = log.GetSystemInfo()

	// Логгируем ошибку
	log.LogError(context.Background(), customErr)

	// Возвращаем ошибку согласно статусу
	return ctx.Status(int(customErr.ErrorType)).JSON(customErr)
}
