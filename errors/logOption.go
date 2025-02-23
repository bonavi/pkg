package errors

// LogOption - Перечисление, необходимое для конкретизации уровня логгирования ошибки
type LogOption int

const (
	LogAsError LogOption = iota + 1
	LogAsWarning
	LogAsDebug
	LogAsInfo
	LogNone
)

var errorTypeToLogOption = map[ErrorType]LogOption{
	BadRequest:     LogAsWarning,
	NotFound:       LogAsWarning,
	Teapot:         LogAsWarning,
	InternalServer: LogAsError,
	Forbidden:      LogAsWarning,
	Unauthorized:   LogAsWarning,
	Timeout:        LogAsWarning,
	BadGateway:     LogAsWarning,
}

func (et ErrorType) logOptionByDefault() LogOption {

	logOption, ok := errorTypeToLogOption[et]
	if !ok {
		return LogAsError
	}

	return logOption
}
