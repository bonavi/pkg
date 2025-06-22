package log

import (
	"os"
	"pkg/stackTrace"
	"time"

	"pkg/errors"
	"pkg/log/model"
)

type Log struct {
	level      LogLevel
	content    any
	params     map[string]any
	stackTrace []string
}

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

func handle(log Log) {
	for _, handler := range logger.handlers {
		handler.handle(log)
	}
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

func emptyLog() Log {
	return Log{
		level:      LevelError,
		content:    nil,
		params:     make(map[string]any),
		stackTrace: stackTrace.GetStackTrace(errors.SkipPreviousCaller),
	}
}

func (l Log) ChangeLog(level LogLevel, content any) Log {

	if len(l.stackTrace) == 0 {
		l.stackTrace = stackTrace.GetStackTrace(errors.Skip2PreviousCallers)
	}

	l.level = level
	l.content = content

	if l.params == nil {
		l.params = make(map[string]any)
	}

	return l
}

func (l Log) Error(content any) {
	handle(l.ChangeLog(LevelError, content))
}

func (l Log) Info(content any) {
	handle(l.ChangeLog(LevelInfo, content))
}

func (l Log) Warning(content any) {
	handle(l.ChangeLog(LevelWarning, content))
}

func (l Log) Debug(content any) {
	handle(l.ChangeLog(LevelDebug, content))
}

func (l Log) Fatal(content any) {
	handle(l.ChangeLog(LevelFatal, content))

	time.Sleep(1 * time.Second)
	os.Exit(1)
}

func (l Log) LogError(err error) {
	customErr := errors.CastError(err)

	switch customErr.ErrorType.LogAs {
	case errors.LogAsError:
		l.Error(err)
	case errors.LogAsWarning:
		l.Warning(err)
	case errors.LogAsDebug:
		l.Debug(err)
	case errors.LogAsInfo:
		l.Info(err)
	case errors.LogNone:
		break
	}
	return
}

func Error(content any)   { emptyLog().Error(content) }
func Warning(content any) { emptyLog().Warning(content) }
func Info(content any)    { emptyLog().Info(content) }
func Fatal(content any)   { emptyLog().Fatal(content) }
func Debug(content any)   { emptyLog().Debug(content) }
func LogError(err error)  { emptyLog().LogError(err) }
