package ipfs

import (
	"context"

	format "github.com/ipfs/go-ipld-format"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type Storage struct {
	forddb.HasListenersBase

	dag format.DAGService
}

func NewStorage(m *Manager) *Storage {
	return &Storage{
		dag: m.DagService(),
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
