package ipfsgraph

import (
	"context"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/ipfs"
)

type Store struct {
	ipfs          *ipfs.Manager
	lsys          linking.LinkSystem
	linkPrototype ipld.LinkPrototype
}

func (s *Store) List(ctx context.Context, typ forddb.TypeID, opts forddb.ListOptions) ([]forddb.RawResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Get(ctx context.Context, typ forddb.TypeID, id forddb.BasicResourceID, opts forddb.GetOptions) (forddb.RawResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Put(ctx context.Context, resource forddb.RawResource, opts forddb.PutOptions) (forddb.RawResource, error) {
	lctx := linking.LinkContext{Ctx: ctx}

	link, err := s.lsys.Store(lctx, s.linkPrototype, resource.TypedNode)

	if err != nil {
		return nil, err
	}

}

func (s *Store) Delete(ctx context.Context, resource forddb.RawResource, opts forddb.DeleteOptions) (forddb.RawResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Close() error {
	return nil
}
