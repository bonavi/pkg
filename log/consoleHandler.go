package log

import (
	"fmt"
	"io"
	"os"
	"pkg/errors"
	"sync/atomic"

	"pkg/log/buffer/buffer"
	"pkg/maps"
)

type consoleLog struct {
	Level   LogLevel
	Message string
	Path    string
	Params  map[string]any
}

var _ Handler = new(ConsoleHandler)

func (h *ConsoleHandler) SetLogLevel(level LogLevel) {
	h.logLevel.Store(level)
}

func (h *ConsoleHandler) GetLogLevel() LogLevel {
	logLevel, ok := h.logLevel.Load().(LogLevel)
	if !ok {
		return ""
	}
	return logLevel
}

// ConsoleHandler - это версия обработчика журналов для печати
// человекочитаемого формата в w.
type ConsoleHandler struct {
	logLevel atomic.Value
	w        io.Writer
}

// NewTextHandler возвращает новый экземпляр ConsoleHandler.
func NewTextHandler(w io.Writer, level LogLevel) *ConsoleHandler {
	h := &ConsoleHandler{
		w:        w,
		logLevel: atomic.Value{},
	}
	h.logLevel.Store(level)
	return h
}

// handle реализует интерфейс Handler.
func (h *ConsoleHandler) handle(log Log) {

	if h.GetLogLevel().GreaterThan(log.level) {
		return
	}

	state := newTextState(buffer.New())
	defer state.buf.Free()

	var logStruct consoleLog

	switch v := log.content.(type) {
	case error:
		customErr := errors.CastError(v)
		var path string
		if len(customErr.StackTrace) > 0 {
			path = customErr.StackTrace[0]
		} else {
			path = ""
		}

		var message string
		if customErr.Err != nil {
			message = customErr.Error()
		} else {
			message = "unknown error"
		}

		logStruct = consoleLog{
			Level:   log.level,
			Message: message,
			Path:    path,
			Params:  maps.Join(log.params, customErr.Params),
		}
	default:
		var path string
		if len(log.stackTrace) > 0 {
			path = log.stackTrace[0]
		} else {
			path = ""
		}
		logStruct = consoleLog{
			Level:   log.level,
			Message: fmt.Sprintf("%v", v),
			Path:    path,
			Params:  log.params,
		}
	}

	var delimer byte = ' '

	state.buf.WriteString(getColor(log.level))

	state.buf.WriteString(log.level.ToUpper())
	state.buf.WriteString(colorReset)
	state.buf.WriteByte(delimer)

	state.buf.WriteString(logStruct.Path)
	state.buf.WriteByte(delimer)
	state.buf.WriteString(logStruct.Message)

	for key, value := range logStruct.Params {
		state.buf.WriteByte(' ')
		state.buf.WriteString(getColor(log.level))
		state.buf.WriteString(key)
		state.buf.WriteString(colorReset)
		state.buf.WriteByte('=')
		state.buf.WriteString(fmt.Sprintf("%+v", value))
	}

	state.buf.WriteByte('\n')

	_, err := state.buf.WriteTo(h.w)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not write log: %s\n", err)
	}
}

type textState struct {
	buf *buffer.Buffer
}

func newTextState(b *buffer.Buffer) textState {
	return textState{buf: b}
}

func getColor(level LogLevel) string {
	switch level {
	case LevelFatal:
		return colorMagenta
	case LevelError:
		return colorRed
	case LevelWarning:
		return colorYellow
	case LevelInfo:
		return colorBlue
	case LevelDebug:
		return colorCyan
	default:
		return colorWhite
	}
}
