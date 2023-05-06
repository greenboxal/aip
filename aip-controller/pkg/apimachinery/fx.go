package apimachinery

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

var Module = fx.Module(
	"apimachinery",

	fx.Provide(NewServer),
	fx.Provide(NewRootMux),

	utils.WithBindingRegistry[HttpServiceMount]("http-service-mounts"),
	utils.WithBindingRegistry[RpcServiceBinding]("rpc-service-bindings"),

	fx.Invoke(func(mux *RootMux, bindings utils.BindingRegistry[HttpServiceMount]) {
		for _, m := range bindings.Bindings() {
			m.Install(mux)
		}
	}),

	fx.Invoke(func(rpcsrv *RpcService, bindings utils.BindingRegistry[RpcServiceBinding]) {
		for _, m := range bindings.Bindings() {
			m.Bind(rpcsrv)
		}
	}),
)
