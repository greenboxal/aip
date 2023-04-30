package ipfs

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

var Module = fx.Module(
	"network/ipfs",

	fx.Provide(NewPubSubManager),
	fx.Provide(NewManager),
	fx.Provide(NewStorage),
)

func WithIpfsStorage() fx.Option {
	return fx.Options(
		fx.Provide(func(db *Storage) forddb.Database {
			return db
		}),
	)
}
