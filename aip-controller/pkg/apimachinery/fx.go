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

	fx.Invoke(func(mux *RootMux, bindings utils.BindingRegistry[HttpServiceMount]) {
		for _, m := range bindings.Bindings() {
			m.Install(mux)
		}
	}),
)
