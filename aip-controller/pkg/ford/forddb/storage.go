package forddb

import "context"

type Storage interface {
	List(ctx context.Context, typ ResourceTypeID) ([]BasicResource, error)
	Get(ctx context.Context, typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	Put(ctx context.Context, resource BasicResource) (BasicResource, error)
	Delete(ctx context.Context, resource BasicResource) (BasicResource, error)
	Close() error
}
