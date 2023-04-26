package api

import "go.uber.org/fx"

var Module = fx.Module(
	"api",

	fx.Provide(NewApi),
	fx.Provide(NewRpcServer),
	fx.Provide(NewSupervisorApi),
)
