package server

import (
	"context"
	"net/http"
	"pkg/errors"
	"time"
)

const (
	readHeaderTimeout = 10 * time.Second
)

type Server struct {
	server *http.Server
}

func GetDefaultServer(
	addr string,
	router http.Handler,
) (*Server, error) {

	// Проверяем, передали ли адрес
	if addr == "" {
		return nil, errors.Default.New("Переменная окружения LISTEN_HTTP не задана").SkipThisCall()
	}

	return &Server{
		&http.Server{
			Addr:                         addr,
			Handler:                      router,
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  0,
			ReadHeaderTimeout:            readHeaderTimeout,
			WriteTimeout:                 0,
			IdleTimeout:                  0,
			MaxHeaderBytes:               0,
			TLSNextProto:                 nil,
			ConnState:                    nil,
			ErrorLog:                     nil,
			BaseContext:                  nil,
			ConnContext:                  nil,
		},
	}, nil
}

func (s *Server) Serve() error {

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return errors.Default.Wrap(err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	_ = s.server.Shutdown(ctx)
}
