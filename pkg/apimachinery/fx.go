package apimachinery

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"apimachinery",

	fx.Provide(NewServer),
	fx.Provide(NewRootMux),

	ProvideHttpService[*RpcService](NewRpcService, "/v1/rpc"),
	ProvideHttpService[*Docs](NewDocs, "/v1/docs"),

	fx.Invoke(
		fx.Annotate(
			func(mux *RootMux, mounts []HttpServiceMount) {
				for _, m := range mounts {
					m.Install(mux)
				}
			},
			fx.ParamTags(``, `group:"http-service-mounts"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			func(server *RpcService, mounts []RpcServiceBinding) {
				for _, m := range mounts {
					m.Bind(server)
				}
			},
			fx.ParamTags(``, `group:"rpc-service-bindings"`),
		),
	),
)

type rpcServiceBinding[T any] struct {
	Name    string
	Handler T
}

func (r *rpcServiceBinding[T]) Bind(server *RpcService) {
	mustRegister(server.Handler, r.Name, r.Handler)
}

type RpcServiceBinding interface {
	Bind(server *RpcService)
}

func BindRpcService[T any](name string) fx.Option {
	return fx.Provide(fx.Annotate(func(handler T) *rpcServiceBinding[T] {
		return &rpcServiceBinding[T]{
			Name:    name,
			Handler: handler,
		}
	}, fx.As((*RpcServiceBinding)(nil)), fx.ResultTags(`group:"rpc-service-bindings"`)))
}

func ProvideRpcService[T any](constructor any, name string) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		BindRpcService[T](name),
	)
}
