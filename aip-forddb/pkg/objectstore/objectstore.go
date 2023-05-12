package objectstore

import (
	"context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type ObjectStore interface {
	List(
		ctx context.Context,
		typ forddb.TypeID,
		opts forddb.ListOptions,
	) ([]forddb.RawResource, error)

	Get(
		ctx context.Context,
		typ forddb.TypeID,
		id forddb.BasicResourceID,
		opts forddb.GetOptions,
	) (forddb.RawResource, error)

	Put(
		ctx context.Context,
		resource forddb.RawResource,
		opts forddb.PutOptions,
	) (forddb.RawResource, error)

	Delete(
		ctx context.Context,
		resource forddb.RawResource,
		opts forddb.DeleteOptions,
	) (forddb.RawResource, error)

	Close() error
}
