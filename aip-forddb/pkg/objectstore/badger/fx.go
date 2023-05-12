package badger

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"ford/forddb/badger",

	config.RegisterConfig[StorageConfig]("ford.storage.badger"),

	fx.Provide(NewStorage),
)

func WithObjectStore() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *Storage) objectstore.ObjectStore {
			return db
		}),
	)
}
