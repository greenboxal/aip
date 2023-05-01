package badger

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/config"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

var Module = fx.Module(
	"ford/forddb/badger",

	config.RegisterConfig[StorageConfig]("ford.storage.badger"),

	fx.Provide(NewStorage),
)

func WithBadgerStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(db *Storage) forddb.Storage {
			return db
		}),
	)
}
