package milvus

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing2"
)

var Module = fx.Module(
	"milvus",

	fx.Provide(NewStorage),
	fx.Provide(NewMilvus),
	fx.Provide(NewIndexProvider),

	config.RegisterConfig[Config]("storage.milvus"),
)

func WithFordStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) forddb.Storage {
			return s
		}),
	)
}

func WithIndexStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) indexing2.MemoryStorage {
			return s
		}),

		fx.Provide(func(s *IndexProvider) indexing.Provider {
			return s
		}),
	)
}
