package openai

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"openai",

	fx.Provide(NewClient),
)
