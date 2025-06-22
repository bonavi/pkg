package log

// Handler - это интерфейс обработчика журналов.
type Handler interface {
	handle(log Log)
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel
}
