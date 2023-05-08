package jobs

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

var Module = fx.Module(
	"jobs",

	fx.Provide(NewManager),
	fx.Provide(NewReconciler),
	fx.Provide(NewSupervisor),

	utils.WithBindingRegistry[JobHandlerBinding]("job-handler-bindings"),
)

func ProvideJobHandler[T JobHandler](name BasicHandlerID, constructor interface{}) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		RegisterJobHandler[T](name),
	)
}

func RegisterJobHandler[T JobHandler](name BasicHandlerID) fx.Option {
	return utils.WithBinding[JobHandlerBinding](
		"job-handler-bindings",
		func(handler T) JobHandlerBinding {
			return &jobHandlerBinding{
				id:      name,
				handler: handler,
			}
		},
	)
}
