package api

import "go.uber.org/fx"

var Module = fx.Module(
	"api",

	fx.Provide(NewServer),
	fx.Provide(NewRpcServer),

	fx.Provide(NewResourcesAPI),

	ProvideRpcService[*SupervisorAPI]("supervisor", NewSupervisorApi),
	ProvideRpcService[*CommunicationAPI]("comms", NewCommunicationApi),
)
