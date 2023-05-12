package tracing

import "go.uber.org/fx"

var Module = fx.Module(
	"tracing",

	fx.Provide(NewTracer),
)
