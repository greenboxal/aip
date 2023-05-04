package memgraph

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford/forddb/memgraph",

	config.RegisterConfig[StorageConfig]("ford.storage.memgraph"),

	fx.Provide(NewStorage),
)

func WithMemgraphStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *Storage) forddb.Storage {
			return db
		}),
	)
}
