package milvus

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
)

var Module = fx.Module(
	"milvus",

	fx.Provide(NewStorage),
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

		fx.Provide(func(s *Storage) indexing.MemoryStorage {
			return s
		}),
	)
}
