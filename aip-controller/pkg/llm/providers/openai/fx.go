package openai

import (
	"os"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"openai",

	fx.Provide(func() *openai.Client {
		return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	}),
)
