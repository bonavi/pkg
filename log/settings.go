package log

type LoggerSettingsEnv struct {
	LogLevel  LogLevel  `env:"LOG_LEVEL"`
	LogFormat LogFormat `env:"LOG_FORMAT"`
}
