package storage

import (
	"context"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type CompositeStorage struct {
}

func (c *CompositeStorage) List(ctx context.Context, typ forddb2.ResourceTypeID) ([]forddb2.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Get(ctx context.Context, typ forddb2.ResourceTypeID, id forddb2.BasicResourceID) (forddb2.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Put(ctx context.Context, resource forddb2.BasicResource) (forddb2.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CompositeStorage) Delete(ctx context.Context, resource forddb2.BasicResource) (forddb2.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}
