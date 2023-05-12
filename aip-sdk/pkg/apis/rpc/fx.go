package rpc

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-sdk/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-sdk/pkg/utils"
)

var Module = fx.Module(
	"apis/rpc",

	utils.WithBindingRegistry[RpcServiceBinding]("rpc-service-bindings"),

	apimachinery.ProvideHttpService[*RpcService](NewRpcService, "/v1/rpc"),
	apimachinery.ProvideHttpService[*Docs](NewDocs, "/v1/docs"),

	fx.Invoke(func(rpcsrv *RpcService, bindings utils.BindingRegistry[RpcServiceBinding]) {
		for _, m := range bindings.Bindings() {
			m.Bind(rpcsrv)
		}
	}),
)
