package ipfs

import (
	"github.com/ipfs/boxo/gateway"
	"github.com/libp2p/go-doh-resolver"
	madns "github.com/multiformats/go-multiaddr-dns"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

var Module = fx.Module(
	"network/ipfs",

	fx.Provide(NewPubSubManager),
	fx.Provide(NewManager),
	fx.Provide(NewStorage),
	fx.Provide(NewResolver),
)

func WithIpfsStorage() fx.Option {
	return fx.Options(
		fx.Provide(func(db *Storage) forddb.Storage {
			return db
		}),
	)
}

func NewResolver() (*madns.Resolver, error) {
	var dohOpts []doh.Option

	resolvers := map[string]string{}

	return gateway.NewDNSResolver(resolvers, dohOpts...)
}
