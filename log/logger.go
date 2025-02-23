package log

import (
	"context"
	"os"
	"time"

	"pkg/errors"
)

// settings - конфигурация логгера
type settings struct {

	// Массив обработчиков лога
	handlers []Handler

	// Дополнительные параметры, которые добавляются в каждый тег и настраиваются при инициализации
	systemInfo any

	userInfoContextKey any
}

// logger - синглтон переменная логгера
var logger = &settings{
	handlers:   []Handler{},
	systemInfo: nil,
}

// Init конфигурирует логгер
func Init(
	systemInfo any,
	userInfoContextKey any,
	handlers ...Handler,
) {
	logger = &settings{
		systemInfo:         systemInfo,
		userInfoContextKey: userInfoContextKey,
		handlers:           handlers,
	}
}

func Off() {
	logger = new(settings)
}

// Error логгирует сообщения для ошибок системы
func Error(ctx context.Context, log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(ctx, LevelError, log, opts...)
	}
}

// Warning логгирует сообщения для ошибок пользователя
func Warning(ctx context.Context, log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(ctx, LevelWarning, log, opts...)
	}
}

// Info логгирует сообщения для информации
func Info(ctx context.Context, log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(ctx, LevelInfo, log, opts...)
	}
}

// Fatal логгирует сообщения для фатальных ошибок и завершает работу программы
func Fatal(ctx context.Context, log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(ctx, LevelFatal, log, opts...)
	}
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// Debug логгирует сообщения для дебага
func Debug(ctx context.Context, log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(ctx, LevelDebug, log, opts...)
	}
}

func LogError(ctx context.Context, err error) {

	customErr := errors.CastError(ctx, err)

	switch customErr.LogAs {
	case errors.LogAsError:
		Error(ctx, err)
	case errors.LogAsWarning:
		Warning(ctx, err)
	case errors.LogAsDebug:
		Debug(ctx, err)
	case errors.LogAsInfo:
		Info(ctx, err)
	case errors.LogNone:
		break
	}
}
