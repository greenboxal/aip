package apimachinery

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type API struct {
	logger *zap.SugaredLogger
	server http.Server
	mux    *chi.Mux
}

func NewServer(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	mux *RootMux,
) *API {
	api := &API{}

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

func (a *API) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", ":30100")

	if err != nil {
		return err
	}

	go func() {
		if err := a.server.Serve(l); err != nil {
			a.logger.Error(err)
		}
	}()

	return nil
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
