package p2p

import (
	"context"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"network/p2p",

	fx.Provide(NewConnMgr),
	fx.Provide(NewNetwork),
	fx.Provide(NewPubSub),

	fx.Provide(func(n *Network) (res struct {
		fx.Out

		Host          host.Host
		PeerStore     peerstore.Peerstore
		PeerRouter    routing.PeerRouting
		ContentRouter routing.ContentRouting
		IpfsRouter    IpfsRouting
	}) {
		res.Host = n.host
		res.ContentRouter = n.dht
		res.PeerRouter = n.dht
		res.IpfsRouter = n.dht
		res.PeerStore = n.host.Peerstore()

		return
	}),
)

func NewConnMgr(lc fx.Lifecycle) (*connmgr.BasicConnMgr, error) {
	cm, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return cm.Close()
		},
	})

	return cm, nil
}

func NewPubSub(h host.Host) (*pubsub.PubSub, error) {
	return pubsub.NewGossipSub(context.Background(), h)
}
