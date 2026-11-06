package log

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"pkg/log/model"
)

// TestWithContextParamsTraceID проверяет, что лог в контексте спана получает traceID и spanID
func TestWithContextParamsTraceID(t *testing.T) {

	// Пишем логи в буфер
	var buf bytes.Buffer
	if err := Init(model.SystemInfo{}, NewJSONHandler(&buf, LevelDebug)); err != nil {
		t.Fatal(err)
	}

	// Контекст с активным спаном
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	}))

	cases := []struct {
		name        string
		ctx         context.Context
		wantTraceID string
		wantSpanID  string
	}{
		{name: "в контексте спана", ctx: ctx, wantTraceID: traceID.String(), wantSpanID: spanID.String()},
		{name: "без спана", ctx: context.Background(), wantTraceID: "", wantSpanID: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()

			// Пишем лог с параметрами контекста
			WithContextParams(c.ctx).Info("запрос")

			// Разбираем записанный JSON
			var got struct {
				TraceID *string `json:"traceID"`
				SpanID  *string `json:"spanID"`
			}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("лог не JSON: %v: %s", err, buf.String())
			}

			// Пустые ID в JSON не пишутся вовсе
			if c.wantTraceID == "" {
				if got.TraceID != nil || got.SpanID != nil {
					t.Fatalf("без спана не должно быть traceID/spanID: %s", buf.String())
				}
				return
			}
			if got.TraceID == nil || *got.TraceID != c.wantTraceID {
				t.Fatalf("traceID = %v, want %s", got.TraceID, c.wantTraceID)
			}
			if got.SpanID == nil || *got.SpanID != c.wantSpanID {
				t.Fatalf("spanID = %v, want %s", got.SpanID, c.wantSpanID)
			}
		})
	}
}
