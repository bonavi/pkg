package log

// Handler - это интерфейс обработчика журналов.
type Handler interface {
	handle(level LogLevel, log any, opts ...Option)
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel
}
