package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/mailru/easyjson"

	"pkg/errors"
	"pkg/log/buffer/buffer"
	"pkg/log/model"
	"pkg/maps"
	"pkg/stackTrace"
)

// jsonLog - Структура лога
//easyjson:json
type jsonLog struct {
	Level      string            `json:"level"`
	Message    string            `json:"message"`
	StackTrace []string          `json:"stackTrace"`
	Params     map[string]string `json:"params,omitempty"`
	SystemInfo model.SystemInfo  `json:"systemInfo"`
}

var _ Handler = new(JSONHandler)

func (h *JSONHandler) SetLogLevel(level LogLevel) {
	h.logLevel.Store(level)
}

func (h *JSONHandler) GetLogLevel() LogLevel {
	logLevel, ok := h.logLevel.Load().(LogLevel)
	if !ok {
		return ""
	}
	return logLevel
}

// JSONHandler - это версия обработчика журналов для печати json в w.
type JSONHandler struct {
	logLevel atomic.Value
	w        io.Writer
}

// NewJSONHandler возвращает новый экземпляр JSONHandler.
func NewJSONHandler(w io.Writer, level LogLevel) *JSONHandler {
	h := &JSONHandler{
		w:        w,
		logLevel: atomic.Value{},
	}
	h.logLevel.Store(level)
	return h
}

// handle реализует интерфейс Handler.
func (h *JSONHandler) handle(level LogLevel, log any, opts ...Option) {

	if h.GetLogLevel().GreaterThan(level) {
		return
	}

	state := newJSONState(buffer.New())
	defer state.buf.Free()

	var logStruct jsonLog

	optsStruct := mergeOptions(opts...)

	switch v := log.(type) {
	case error:
		customErr := errors.CastError(v)
		logStruct = jsonLog{
			Level:      level.String(),
			Message:    customErr.Error(),
			StackTrace: customErr.StackTrace,
			Params:     maps.Join(optsStruct.params, customErr.Params),
			SystemInfo: logger.systemInfo,
		}
	default:
		logStruct = jsonLog{
			Level:      level.String(),
			Message:    fmt.Sprintf("%v", v),
			StackTrace: stackTrace.GetStackTrace(errors.SkipPreviousCaller),
			Params:     optsStruct.params,
			SystemInfo: logger.systemInfo,
		}
	}

	json, err := easyjson.Marshal(logStruct)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not generate json jsonLog: %s\n", err)
		return
	}

	_, _ = state.buf.Write(json)

	state.buf.WriteByte('\n')

	_, err = state.buf.WriteTo(h.w)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logging: could not write jsonLog: %s\n", err)
	}
}

type jsonState struct {
	buf *buffer.Buffer
}

func newJSONState(b *buffer.Buffer) jsonState {
	return jsonState{buf: b}
}
