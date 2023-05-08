package forddb

import "context"

type Storage interface {
	List(
		ctx context.Context,
		typ TypeID,
		opts ListOptions,
	) ([]RawResource, error)

	Get(
		ctx context.Context,
		typ TypeID,
		id BasicResourceID,
		opts GetOptions,
	) (RawResource, error)

	Put(
		ctx context.Context,
		resource RawResource,
		opts PutOptions,
	) (RawResource, error)

	Delete(
		ctx context.Context,
		resource RawResource,
		opts DeleteOptions,
	) (RawResource, error)

	Close() error
}
