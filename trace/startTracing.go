package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"pkg/errors"
)

// shutdownTimeout — сколько ждём отправки оставшихся спанов при остановке
const shutdownTimeout = 5 * time.Second

// StartTracing настраивает глобальный TracerProvider с отправкой трейсов по OTLP HTTP.
// Возвращает функцию остановки: её нужно вызвать при завершении сервиса, чтобы досылать
// накопленные спаны. Если трейсинг выключен, глобальный провайдер остаётся no-op,
// а функция остановки ничего не делает.
func StartTracing(ctx context.Context, cfg TracerConfig, service ServiceInfo) (func(), error) {

	// Трейсинг выключен — спаны остаются no-op
	if !cfg.IsEnabled {
		return func() {}, nil
	}

	// Настраиваем экспортёр OTLP HTTP
	options := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.HTTPEndpoint)}
	if cfg.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return nil, errors.Default.Wrap(err)
	}

	// Описываем сервис, от имени которого идут спаны
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(service.Name),
		semconv.ServiceVersion(service.Version),
		semconv.DeploymentEnvironment(service.Environment),
	))
	if err != nil {
		return nil, errors.Default.Wrap(err)
	}

	// Регистрируем провайдер с пакетной отправкой
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	// Передаём контекст трейса между сервисами в заголовках W3C traceparent
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Функция остановки досылает накопленные спаны
	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = provider.Shutdown(shutdownCtx)
	}, nil
}
