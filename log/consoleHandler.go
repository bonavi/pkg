package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"pkg/errors"
	"pkg/log/buffer/buffer"
	"pkg/maps"
	"pkg/stackTrace"
)

type consoleLog struct {
	Level   LogLevel
	Message string
	Path    string
	Params  map[string]string
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
func (h *ConsoleHandler) handle(level LogLevel, log any, opts ...Option) {

	if h.GetLogLevel().GreaterThan(level) {
		return
	}

	state := newTextState(buffer.New())
	defer state.buf.Free()

	var logStruct consoleLog

	optsStruct := mergeOptions(opts...)

	switch v := log.(type) {
	case error:
		customErr := errors.CastError(v)
		var path string
		if len(customErr.StackTrace) > 0 {
			path = customErr.StackTrace[0]
		} else {
			path = ""
		}
		logStruct = consoleLog{
			Level:   level,
			Message: customErr.Error(),
			Path:    path,
			Params:  maps.Join(optsStruct.params, customErr.Params),
		}
	default:
		stackTrace := stackTrace.GetStackTrace(errors.SkipPreviousCaller)
		var path string
		if len(stackTrace) > 0 {
			path = stackTrace[0]
		} else {
			path = ""
		}
		logStruct = consoleLog{
			Level:   level,
			Message: fmt.Sprintf("%v", v),
			Path:    path,
			Params:  optsStruct.params,
		}
	}

	var delimer byte = ' '

	state.buf.WriteString(getColor(level))

	state.buf.WriteString(level.ToUpper())
	state.buf.WriteString(colorReset)
	state.buf.WriteByte(delimer)

	state.buf.WriteString(logStruct.Path)
	state.buf.WriteByte(delimer)
	state.buf.WriteString(logStruct.Message)

	for key, value := range logStruct.Params {
		state.buf.WriteByte(' ')
		state.buf.WriteString(getColor(level))
		state.buf.WriteString(key)
		state.buf.WriteString(colorReset)
		state.buf.WriteByte('=')
		state.buf.WriteString(value)
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
