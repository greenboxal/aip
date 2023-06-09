package apimachinery

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	server http.Server
	mux    *chi.Mux
}

func NewServer(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	mux *RootMux,
) *Server {
	api := &Server{}

	api.logger = logger.Named("api")
	api.server.Handler = mux

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return api.Start(ctx)
		},

		OnStop: func(ctx context.Context) error {
			return api.Shutdown(ctx)
		},
	})

	return api
}

func (a *Server) Start(ctx context.Context) error {
	endpoint := os.Getenv("AIP_LISTEN_ENDPOINT")

	if endpoint == "" {
		endpoint = "0.0.0.0:30100"
	}

	l, err := net.Listen("tcp", endpoint)

	if err != nil {
		return err
	}

	a.logger.Infow("Server is listening", "endpoint", endpoint)

	go func() {
		if err := a.server.Serve(l); err != nil {
			a.logger.Error(err)
		}
	}()

	return nil
}

func (a *Server) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
