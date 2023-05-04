package inmemory

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford/forddb/inmemory",

	fx.Provide(NewInMemory),
)

func WithInMemoryDatabase() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *InMemoryDatabase) forddb.Storage {
			return db
		}),
	)
}
