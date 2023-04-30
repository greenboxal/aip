package ipfs

import (
	"context"

	graphsync "github.com/ipfs/go-graphsync"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipld/go-ipld-prime"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Storage struct {
	forddb.HasListenersBase

	dag format.DAGService
	gs  graphsync.GraphExchange
}

func NewStorage(m *Manager) *Storage {
	return &Storage{
		dag: m.DagService(),
		gs:  m.GraphExchange(),
	}
}

func (s *Storage) List(ctx context.Context, typ forddb.ResourceTypeID) ([]forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Get(ctx context.Context, typ forddb.ResourceTypeID, id forddb.BasicResourceID) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Put(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Delete(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) syncGraphFrom(
	ctx context.Context,
	peerId peer.ID,
	pointer forddb.BasicResourcePointer,
	selector ipld.Node,
) ([]forddb.BasicResource, error) {
	//root := pointer.AsLink()

	//var recursionLimit int64
	//ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	//sel := ssb.ExploreRecursive(
	//	selector.RecursionLimitDepth(recursionLimit),
	//	ssb.ExploreFields(
	//		func(efsb builder.ExploreFieldsSpecBuilder) {
	//			efsb.Insert(
	//				"Parents",
	//				ssb.ExploreAll(
	//					ssb.ExploreRecursiveEdge(),
	//				),
	//			)
	//		},
	//	),
	//)

	//resCh, errCh := s.gs.Request(ctx, peerId, root, selector)

	panic("not implemented")
}
