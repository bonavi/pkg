package trace

// TracerConfig — настройки отправки трейсов по OTLP HTTP
type TracerConfig struct {
	IsEnabled    bool   `env:"OTLP_IS_ENABLED"`    // Включить отправку трейсов
	HTTPEndpoint string `env:"OTLP_HTTP_ENDPOINT"` // Адрес приёмника OTLP HTTP без схемы, например alloy.alloy.svc.cluster.local:4318
	Insecure     bool   `env:"OTLP_INSECURE"`      // Отправлять без TLS (внутри кластера)
}

// ServiceInfo — атрибуты сервиса, которые попадают в каждый трейс
type ServiceInfo struct {
	Name        string // service.name
	Version     string // service.version
	Environment string // deployment.environment.name
}
