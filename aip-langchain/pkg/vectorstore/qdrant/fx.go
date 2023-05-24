package qdrant

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"vectorstore-qdrant",

	fx.Provide(NewProvider),

	config.RegisterConfig[Config]("storage.qdrant"),
)

func WithIndexStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Provider) vectorstore.Provider {
			return s
		}),
	)
}
