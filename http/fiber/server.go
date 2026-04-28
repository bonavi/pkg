package fiber

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"pkg/contextKeys"
)

func GetDefaultServer(
	appName string,
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
		ReadTimeout:                  15 * time.Second,
		WriteTimeout:                 30 * time.Second,
		IdleTimeout:                  60 * time.Second,
		ReadBufferSize:               0,
		WriteBufferSize:              0,
		CompressedFileSuffix:         "",
		ProxyHeader:                  "X-Forwarded-For",
		GETOnly:                      false,
		ErrorHandler:                 nil,
		DisableKeepalive:             false,
		DisableDefaultDate:           false,
		DisableDefaultContentType:    false,
		DisableHeaderNormalizing:     false,
		DisableStartupMessage:        true,
		AppName:                      appName,
		StreamRequestBody:            false,
		DisablePreParseMultipartForm: false,
		ReduceMemoryUsage:            false,
		JSONEncoder:                  nil,
		JSONDecoder:                  nil,
		XMLEncoder:                   nil,
		Network:                      "",
		EnableTrustedProxyCheck:      true,
		TrustedProxies:               []string{"127.0.0.1"},
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
