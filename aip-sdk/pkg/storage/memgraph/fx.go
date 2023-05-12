package memgraph

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"ford/forddb/memgraph",

	config.RegisterConfig[StorageConfig]("ford.storage.memgraph"),

	fx.Provide(NewStorage),
)

func WithMemgraphStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *Storage) objectstore.ObjectStore {
			return db
		}),
	)
}
