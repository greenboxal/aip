package ipfs

import (
	"github.com/ipfs/boxo/gateway"
	"github.com/libp2p/go-doh-resolver"
	madns "github.com/multiformats/go-multiaddr-dns"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"network/ipfs",

	fx.Provide(NewPubSubManager),
	fx.Provide(NewManager),
	fx.Provide(NewStorage),
	fx.Provide(NewResolver),
)

func NewResolver() (*madns.Resolver, error) {
	var dohOpts []doh.Option

	resolvers := map[string]string{}

	return gateway.NewDNSResolver(resolvers, dohOpts...)
}
