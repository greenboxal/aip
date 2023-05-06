package p2p

import (
	"context"
	"fmt"

	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
)

type IpfsRouting interface {
	routing.ContentRouting
	routing.PeerRouting
	routing.ValueStore

	Bootstrap(context.Context) error
}

type Network struct {
	logger *zap.SugaredLogger

	nm *config.NetworkManager
	cm *connmgr.BasicConnMgr

	host          host.Host
	mdns          mdns.Service
	dht           *dual.DHT
	peerRouter    routing.PeerRouting
	providerStore *providers.ProviderManager
}

func NewNetwork(
	lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	cm *connmgr.BasicConnMgr,
	nm *config.NetworkManager,
) (*Network, error) {
	n := &Network{}

	n.logger = logger.Named("p2p-network")
	n.nm = nm
	n.cm = cm

	if err := n.initializeHost(); err != nil {
		return nil, err
	}

	if err := n.initializeMdns(); err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return n.Start(ctx)
		},

		OnStop: func(ctx context.Context) error {
			return n.Shutdown(ctx)
		},
	})

	return n, nil
}

func (n *Network) Start(ctx context.Context) error {
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/27050")

	if err != nil {
		return err
	}

	addrs := []multiaddr.Multiaddr{addr}

	if err := n.host.Network().Listen(addrs...); err != nil {
		panic(err)
	}

	if err := n.dht.Bootstrap(ctx); err != nil {
		panic(err)
	}

	if err := n.mdns.Start(); err != nil {
		panic(err)
	}

	return nil
}

func (n *Network) initializeHost() error {
	key, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)

	if err != nil {
		return err
	}

	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(key),
		// Multiple listen addresses
		libp2p.NoListenAddrs,
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultPrivateTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(n.cm),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(n.initializeDht),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.EnableRelay(),
		libp2p.EnableRelayService(),
		//libp2p.EnableAutoRelayWithPeerSource(n.findBootstrapPeers),
	)

	if err != nil {
		return err
	}

	n.host = h

	return nil
}

func (n *Network) initializeDht(h host.Host) (routing.PeerRouting, error) {
	dataStore := dssync.MutexWrap(ds.NewMapDatastore())
	providerStore, err := providers.NewProviderManager(context.Background(), h.ID(), h.Peerstore(), dataStore)

	if err != nil {
		return nil, fmt.Errorf("initializing default provider manager (%v)", err)
	}

	if n.peerRouter != nil {
		return n.peerRouter, nil
	}

	dhtOptions := []dual.Option{
		dual.DHTOption(
			dht.ProviderStore(providerStore),
		),
	}

	d, err := dual.New(context.Background(), h, dhtOptions...)

	if err != nil {
		return nil, err
	}

	n.dht = d
	n.providerStore = providerStore
	n.peerRouter = n.dht

	return n.peerRouter, nil
}

func (n *Network) initializeMdns() error {
	n.mdns = mdns.NewMdnsService(n.host, "npnsd", n)

	return nil
}

func (n *Network) Shutdown(ctx context.Context) error {
	if n.dht != nil {
		if err := n.dht.Close(); err != nil {
			return err
		}

		n.dht = nil
	}

	if n.host != nil {
		if err := n.host.Close(); err != nil {
			return err
		}

		n.host = nil
	}

	return nil
}

func (n *Network) HandlePeerFound(info peer.AddrInfo) {
	n.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)
}
