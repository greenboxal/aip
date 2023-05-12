package firestore

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
)

var Module = fx.Module(
	"aip-forddb/objectstore/firestore",

	fx.Provide(NewStorage),

	config.RegisterConfig[Config]("storage.firestore"),
)

func WithObjectStore() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) objectstore.ObjectStore {
			return s
		}),
	)
}
