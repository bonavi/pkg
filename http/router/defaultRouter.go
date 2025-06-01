package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"pkg/http/middleware"
)

func NewRouter() *chi.Mux {

	r := chi.NewRouter()

	// middlewares
	r.Use(
		middleware.PanicRecover,
		middleware.ResponseTime,
		middleware.RequestID,
		middleware.Logger,
	)

	// prometheus
	r.Handle("/metrics", promhttp.Handler())

	// healthCheck
	r.Get("/health", HealthCheck)

	return r
}
