package apimachinery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

type MountOption func(opts *MountOptions)

func WithStripPrefix() MountOption {
	return func(opts *MountOptions) {
		opts.StripPrefix = true
	}
}

type MountOptions struct {
	StripPrefix bool
}

var Module = fx.Module(
	"apimachinery",

	fx.Provide(NewServer),
	fx.Provide(NewRootMux),

	ProvideHttpService[*RpcService](NewRpcService, "/v1/rpc"),
	ProvideHttpService[*Docs](NewDocs, "/v1/docs"),
)

func NewMountOptions(opts ...MountOption) *MountOptions {
	options := &MountOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func MountHttpService[T http.Handler](path string, options ...MountOption) fx.Option {
	opts := NewMountOptions(options...)

	return fx.Invoke(func(mux *RootMux, handler T) {
		if opts.StripPrefix {
			mux.Route(path, func(r chi.Router) {
				r.Use(middleware.PathRewrite(path, ""))
				r.Mount("/", handler)
			})
		} else {
			mux.Mount(path, handler)
		}
	})
}

func ProvideHttpService[T http.Handler](constructor any, path string) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		MountHttpService[T](path),
	)
}

func ProvideRpcService[T any](constructor any, name string) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		fx.Invoke(func(handler *RpcService, svc T) {
			mustRegister(handler.Handler, name, svc)
		}),
	)
}
