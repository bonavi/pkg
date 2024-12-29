package fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"pkg/contextKeys"
)

type ServerSettingsEnv struct {
	Port                    string        `env:"FIBER_SERVER_REST_PORT,required"`
	ReadTimeout             time.Duration `env:"FIBER_SERVER_READ_TIMEOUT,required"`
	WriteTimeout            time.Duration `env:"FIBER_SERVER_WRITE_TIMEOUT,required"`
	IdleTimeout             time.Duration `env:"FIBER_SERVER_IDLE_TIMEOUT,required"`
	EnableTrustedProxyCheck bool          `env:"FIBER_SERVER_TRUSTED_PROXY_CHECK,required"`
	TrustedProxies          []string      `env:"FIBER_SERVER_TRUSTED_PROXIES,required"`
	ProxyHeader             string        `env:"FIBER_SERVER_PROXY_HEADER,required"`
}

func GetDefaultServer(
	appName string,
	serverCfg ServerSettingsEnv,
	readyIndicator chan struct{},
) (*fiber.App, error) {

	app := fiber.New(fiber.Config{
		Prefork:                      false,
		ServerHeader:                 "",
		StrictRouting:                false,
		CaseSensitive:                false,
		Immutable:                    false,
		UnescapePath:                 false,
		ETag:                         false,
		BodyLimit:                    0,
		Concurrency:                  0,
		Views:                        nil,
		ViewsLayout:                  "",
		PassLocalsToViews:            false,
		ReadTimeout:                  serverCfg.ReadTimeout,
		WriteTimeout:                 serverCfg.WriteTimeout,
		IdleTimeout:                  serverCfg.IdleTimeout,
		ReadBufferSize:               0,
		WriteBufferSize:              0,
		CompressedFileSuffix:         "",
		ProxyHeader:                  serverCfg.ProxyHeader,
		GETOnly:                      false,
		ErrorHandler:                 nil,
		DisableKeepalive:             false,
		DisableDefaultDate:           false,
		DisableDefaultContentType:    false,
		DisableHeaderNormalizing:     false,
		DisableStartupMessage:        false,
		AppName:                      appName,
		StreamRequestBody:            false,
		DisablePreParseMultipartForm: false,
		ReduceMemoryUsage:            false,
		JSONEncoder:                  nil,
		JSONDecoder:                  nil,
		XMLEncoder:                   nil,
		Network:                      "",
		EnableTrustedProxyCheck:      serverCfg.EnableTrustedProxyCheck,
		TrustedProxies:               serverCfg.TrustedProxies,
		EnableIPValidation:           false,
		EnablePrintRoutes:            false,
		ColorScheme: fiber.Colors{
			Black:   "",
			Red:     "",
			Green:   "",
			Yellow:  "",
			Blue:    "",
			Magenta: "",
			Cyan:    "",
			White:   "",
			Reset:   "",
		},
		RequestMethods:           nil,
		EnableSplittingOnParsers: false,
	})

	app.Use(requestid.New(requestid.Config{
		Next:       nil,
		Header:     "X-Request-ID",
		Generator:  nil,
		ContextKey: contextKeys.XRequestIDKey,
	}))
	app.Use(pprof.New())
	app.Use(recover.New(recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: RecoverHandler,
	}))
	app.Use(healthcheck.New(healthcheck.Config{
		Next:              nil,
		LivenessProbe:     nil,
		LivenessEndpoint:  "",
		ReadinessProbe:    NewReadyHandler(readyIndicator),
		ReadinessEndpoint: "",
	}))

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return app, nil
}
