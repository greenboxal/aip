package api

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggest/jsonrpc"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v3cdn"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type API struct {
	logger *zap.SugaredLogger
	server http.Server
	mux    *chi.Mux
}

func NewApi(lc fx.Lifecycle, logger *zap.SugaredLogger, rpc *jsonrpc.Handler) *API {

	api := &API{}

	api.mux = chi.NewMux()

	api.mux.Use(middleware.RealIP)
	api.mux.Use(middleware.Logger)
	api.mux.Use(middleware.Recoverer)
	api.mux.Use(cors.AllowAll().Handler)
	api.mux.Use(middleware.AllowContentEncoding("gzip"))
	api.mux.Use(middleware.Heartbeat("/ping"))

	api.mux.Mount("/rpc/v1", rpc)

	api.mux.Method(http.MethodGet, "/docs/openapi.json", rpc.OpenAPI)

	api.mux.Mount("/docs", v3cdn.NewHandlerWithConfig(swgui.Config{
		Title:       "AIP SUPD",
		SwaggerJSON: "/docs/openapi.json",
		BasePath:    "/docs",
		SettingsUI:  jsonrpc.SwguiSettings(nil, "/rpc/v1"),
	}))

	api.logger = logger.Named("api")
	api.server.Handler = api.mux

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
