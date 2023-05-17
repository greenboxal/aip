package apimachinery

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
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

type httpServiceMount[T any] struct {
	Path    string
	Handler http.Handler
	Options MountOptions
}

type HttpServiceMount interface {
	Install(mux *RootMux)
}

func (m *httpServiceMount[T]) Install(mux *RootMux) {
	handler := m.Handler

	if m.Options.StripPrefix {
		handler = http.StripPrefix(m.Path, handler)
	}

	mux.Mount(m.Path, handler)
}

func NewMountOptions(opts ...MountOption) MountOptions {
	options := MountOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func MountHttpService[T http.Handler](path string, options ...MountOption) fx.Option {
	opts := NewMountOptions(options...)

	return utils.WithBinding[HttpServiceMount]("http-service-mounts", func(handler T) HttpServiceMount {
		return &httpServiceMount[T]{
			Path:    path,
			Handler: handler,
			Options: opts,
		}
	})
}

func ProvideHttpService[T http.Handler](constructor any, path string, options ...MountOption) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		MountHttpService[T](path, options...),
	)
}
