package milvus

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-langchain/pkg/indexing"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"milvus",

	fx.Provide(NewStorage),
	fx.Provide(NewMilvus),
	fx.Provide(NewProvider),

	config.RegisterConfig[Config]("storage.milvus"),
)

func WithFordStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) objectstore.ObjectStore {
			return s
		}),
	)
}

func WithIndexStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) indexing.MemoryStorage {
			return s
		}),

		fx.Provide(func(s *Provider) vectorstore.Provider {
			return s
		}),
	)
}
