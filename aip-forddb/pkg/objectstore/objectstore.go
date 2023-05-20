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

type ObjectIterator interface {
	Len() int

	Next()
	HasNext() bool

	Slice() []Object
	Map() map[forddb.BasicResourceID]Object
}

type TypedObjectIterator[T any] interface {
	ObjectIterator

	TypedSlice() []T
	TypedMap() map[forddb.BasicResourceID]Object
}

type ObjectStoreV2 interface {
	List(
		ctx context.Context,
		typ forddb.TypeID,
		opts forddb.ListOptions,
	) (ObjectIterator, error)

	Get(
		ctx context.Context,
		id forddb.BasicResourceID,
		opts forddb.GetOptions,
	) (forddb.RawResource, error)

	Put(
		ctx context.Context,
		obj Object,
		opts forddb.PutOptions,
	) (forddb.RawResource, error)

	Delete(
		ctx context.Context,
		resource Object,
		opts forddb.DeleteOptions,
	) (Object, error)

	Close() error
}
