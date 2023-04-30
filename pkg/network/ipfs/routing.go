package ipfs

import (
	"context"

	"github.com/ipfs/boxo/routing/http/client"
	"github.com/ipfs/boxo/routing/http/contentrouter"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"

	"github.com/greenboxal/aip/pkg/network/p2p"
)

type RemoteContentRouter struct {
	client contentrouter.Client
	router routing.ContentRouting
}

func NewRemoteContentRouter(endpoint string) (*RemoteContentRouter, error) {
	c, err := client.New(endpoint)

	if err != nil {
		return nil, err
	}

	router := contentrouter.NewContentRoutingClient(c)

	return &RemoteContentRouter{
		client: c,
		router: router,
	}, nil
}

func (r *RemoteContentRouter) Provide(ctx context.Context, cid cid.Cid, b bool) error {
	return r.router.Provide(ctx, cid, b)
}

func (r *RemoteContentRouter) FindProvidersAsync(ctx context.Context, cid cid.Cid, i int) <-chan peer.AddrInfo {
	return r.router.FindProvidersAsync(ctx, cid, i)
}

type DhtContentRouter struct {
	dht p2p.IpfsRouting
}

func NewDhtContentRouter(dht p2p.IpfsRouting) *DhtContentRouter {
	return &DhtContentRouter{
		dht: dht,
	}
}

func (c *DhtContentRouter) Provide(ctx context.Context, cid cid.Cid, b bool) error {
	return c.dht.Provide(ctx, cid, b)
}

func (c *DhtContentRouter) FindProvidersAsync(ctx context.Context, cid cid.Cid, i int) <-chan peer.AddrInfo {
	return c.dht.FindProvidersAsync(ctx, cid, i)
}
