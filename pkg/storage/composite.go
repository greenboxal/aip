package storage

import (
	"context"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type CompositeStorage struct {
}

func (c *CompositeStorage) List(ctx context.Context, typ forddb.ResourceTypeID) ([]forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Get(ctx context.Context, typ forddb.ResourceTypeID, id forddb.BasicResourceID) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Put(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Delete(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}
