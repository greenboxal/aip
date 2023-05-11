package inmemory

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing2"
)

var Module = fx.Module(
	"ford/forddb/inmemory",

	fx.Provide(NewInMemory),
)

func WithInIndexingMemoryDatabase() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *InMemoryDatabase) indexing2.MemoryStorage {
			return db
		}),
	)
}

func WithInMemoryDatabase() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *InMemoryDatabase) forddb.Storage {
			return db
		}),

		fx.Provide(func(db *InMemoryDatabase) indexing2.MemoryStorage {
			return db
		}),
	)
}
