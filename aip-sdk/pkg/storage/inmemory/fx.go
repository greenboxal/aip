package inmemory

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-langchain/pkg/indexing"
)

var Module = fx.Module(
	"ford/forddb/inmemory",

	fx.Provide(NewInMemory),
)

func WithInIndexingMemoryDatabase() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *InMemoryDatabase) indexing.MemoryStorage {
			return db
		}),
	)
}

func WithInMemoryDatabase() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *InMemoryDatabase) objectstore.ObjectStore {
			return db
		}),

		fx.Provide(func(db *InMemoryDatabase) indexing.MemoryStorage {
			return db
		}),
	)
}
