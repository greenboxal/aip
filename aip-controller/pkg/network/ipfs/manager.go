package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	time "time"

	"github.com/cockroachdb/errors"
	"github.com/ipfs/boxo/bitswap"
	bsnet "github.com/ipfs/boxo/bitswap/network"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/gateway"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/namesys"
	pin "github.com/ipfs/boxo/pinning/pinner"
	"github.com/ipfs/boxo/pinning/pinner/dspinner"
	badger "github.com/ipfs/go-ds-badger"
	graphsync "github.com/ipfs/go-graphsync"
	gsimpl "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
	madns "github.com/multiformats/go-multiaddr-dns"
	"go.uber.org/multierr"

	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-controller/pkg/config"
	"github.com/greenboxal/aip/aip-controller/pkg/network/p2p"
)

type Manager struct {
	rsm *config.ResourceManager

	host host.Host
	dht  p2p.IpfsRouting
	psm  *PubSubManager

	contentRouter routing.ContentRouting

	bitswap        *bitswap.Bitswap
	bitswapNetwork bsnet.BitSwapNetwork

	dataStore *badger.Datastore
	gcLocker  blockstore.GCLocker

	blockStore   blockstore.Blockstore
	blockService blockservice.BlockService
	blockGateway *gateway.BlocksGateway

	graphSync        graphsync.GraphExchange
	graphSyncNetwork gsnet.GraphSyncNetwork

	pinner     pin.Pinner
	dagService format.DAGService

	ns       namesys.NameSystem
	resolver *madns.Resolver

	linkSystem ipld.LinkSystem
}

func NewManager(
	lc fx.Lifecycle,
	rsm *config.ResourceManager,
	h host.Host,
	dht p2p.IpfsRouting,
	psm *PubSubManager,
	resolver *madns.Resolver,
) (*Manager, error) {
	router := NewDhtContentRouter(dht)
	bsn := bsnet.NewFromIpfsHost(h, router)
	gsn := gsnet.NewFromLibp2pHost(h)

	m := &Manager{
		rsm:              rsm,
		host:             h,
		dht:              dht,
		psm:              psm,
		resolver:         resolver,
		bitswapNetwork:   bsn,
		graphSyncNetwork: gsn,
		contentRouter:    router,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.Shutdown(ctx)
		},
	})

	if err := m.Initialize(context.Background()); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) initializeIpns(ctx context.Context) error {
	var err error

	opts := []namesys.Option{
		namesys.WithDatastore(m.dataStore),
		namesys.WithDNSResolver(m.resolver),
		namesys.WithCache(128),
	}

	m.ns, err = namesys.NewNameSystem(m.dht, opts...)

	if err != nil {
		return err
	}

	//repub := republisher.NewRepublisher(
	//	m.ns,
	//	m.dataStore,
	//	m.kernel.PrivateKey().PrivateKey(),
	//	m.ks.Ipfs(),
	//)

	return nil
}

func (m *Manager) Initialize(ctx context.Context) error {
	var err error

	dbPath := m.rsm.GetDataDirectory("ipfs")
	options := badger.DefaultOptions

	ds, err := badger.NewDatastore(dbPath, &options)

	if err != nil {
		return err
	}

	m.dataStore = ds
	m.gcLocker = blockstore.NewGCLocker()

	bs := blockstore.NewBlockstore(ds)
	bs = blockstore.NewIdStore(bs)
	bs = blockstore.NewGCBlockstore(bs, m.gcLocker)

	m.blockStore = bs

	if err := m.initializeIpns(ctx); err != nil {
		return err
	}

	m.bitswap = bitswap.New(
		ctx,
		m.bitswapNetwork,
		bs,
		bitswap.ProvideEnabled(true),
		bitswap.ProviderSearchDelay(1000*time.Millisecond),
		bitswap.EngineBlockstoreWorkerCount(128),
		bitswap.TaskWorkerCount(8),
		bitswap.EngineTaskWorkerCount(8),
		bitswap.MaxOutstandingBytesPerPeer(1<<20),
	)

	m.blockService = blockservice.New(bs, m.bitswap)

	m.blockGateway, err = gateway.NewBlocksGateway(
		m.blockService,
		gateway.WithValueStore(m.dht),
		gateway.WithNameSystem(m.ns),
	)

	if err != nil {
		return err
	}

	m.dagService = merkledag.NewDAGService(m.blockService)

	m.linkSystem = storeutil.LinkSystemForBlockstore(bs)

	m.linkSystem.StorageReadOpener = func(
		lnkCtx linking.LinkContext,
		lnk datamodel.Link,
	) (_ io.Reader, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = errors.Errorf("%s", e)
			}
		}()

		asCidLink, ok := lnk.(cidlink.Link)

		if !ok {
			return nil, fmt.Errorf("unsupported link type")
		}

		block, err := m.blockService.GetBlock(lnkCtx.Ctx, asCidLink.Cid)

		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(block.RawData()), nil
	}

	m.graphSync = gsimpl.New(ctx, m.graphSyncNetwork, m.linkSystem)

	m.pinner, err = dspinner.New(ctx, ds, m.dagService)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
	var merr error

	if m.pinner != nil {
		if err := m.pinner.Flush(ctx); err != nil {
			merr = multierr.Append(merr, err)
		}
	}

	if m.blockService != nil {
		if err := m.blockService.Close(); err != nil {
			merr = multierr.Append(merr, err)
		}

		m.blockService = nil
	}

	if m.bitswap != nil {
		if err := m.bitswap.Close(); err != nil {
			merr = multierr.Append(merr, err)
		}

		m.bitswap = nil
	}

	return nil
}

func (m *Manager) ContentRouter() routing.ContentRouting {
	return m.contentRouter
}

func (m *Manager) BlockService() blockservice.BlockService {
	return m.blockService
}

func (m *Manager) DagService() format.DAGService {
	return m.dagService
}

func (m *Manager) GraphExchange() graphsync.GraphExchange {
	return m.graphSync
}

func (m *Manager) NamingService() namesys.NameSystem {
	return m.ns
}

func (m *Manager) BlockGateway() *gateway.BlocksGateway {
	return m.blockGateway
}

func (m *Manager) Pinner() pin.Pinner {
	return m.pinner
}

func (m *Manager) Host() host.Host {
	return m.host
}

func (m *Manager) PubSub() iface.PubSubAPI {
	return m.psm
}
