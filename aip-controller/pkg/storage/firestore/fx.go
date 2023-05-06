package firestore

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var Module = fx.Module(
	"firestore",

	fx.Provide(NewStorage),

	config.RegisterConfig[Config]("storage.firestore"),
)

func WithFordStorage() fx.Option {
	return fx.Options(
		Module,

		fx.Provide(func(s *Storage) forddb.Storage {
			return s
		}),
	)
}
