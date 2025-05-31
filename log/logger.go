package log

import (
	"os"
	"time"

	"pkg/errors"
	"pkg/log/model"
)

// loggerSettings - конфигурация логгера
type loggerSettings struct {

	// Массив обработчиков лога
	handlers []Handler

	// Дополнительные параметры, которые добавляются в каждый тег и настраиваются при инициализации
	systemInfo model.SystemInfo
}

// logger - синглтон переменная логгера
var logger = &loggerSettings{
	systemInfo: model.SystemInfo{
		BuildDate:   "",
		Hostname:    "",
		Version:     "",
		ServiceName: "",
		Build:       "",
		Env:         "",
	},
	handlers: []Handler{NewTextHandler(os.Stdout, LevelDebug)},
}

// Init конфигурирует логгер
func Init(
	systemInfo model.SystemInfo,
	handlers ...Handler,
) error {
	logger = &loggerSettings{
		systemInfo: systemInfo,
		handlers:   handlers,
	}
	return nil
}

func ChangeLogLevel(level LogLevel) {
	for _, handler := range logger.handlers {
		handler.SetLogLevel(level)
	}
}

func GetLogLevel() LogLevel {
	if len(logger.handlers) == 0 {
		return ""
	}
	return logger.handlers[0].GetLogLevel()
}

func GetSystemInfo() model.SystemInfo {
	return logger.systemInfo
}

// Error логгирует сообщения для ошибок системы
func Error(log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(LevelError, log, opts...)
	}
}

// Warning логгирует сообщения для ошибок пользователя
func Warning(log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(LevelWarning, log, opts...)
	}
}

// Info логгирует сообщения для информации
func Info(log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(LevelInfo, log, opts...)
	}
}

// Fatal логгирует сообщения для фатальных ошибок и завершает работу программы
func Fatal(log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(LevelFatal, log, opts...)
	}
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// Debug логгирует сообщения для дебага
func Debug(log any, opts ...Option) {
	for _, handler := range logger.handlers {
		handler.handle(LevelDebug, log, opts...)
	}
}

func LogError(err error, opts ...Option) {

	customErr := errors.CastError(err)

	switch customErr.LogAs {
	case errors.LogAsError:
		Error(err, opts...)
	case errors.LogAsWarning:
		Warning(err, opts...)
	case errors.LogAsDebug:
		Debug(err, opts...)
	case errors.LogAsInfo:
		Info(err, opts...)
	case errors.LogNone:
		break
	}
}
