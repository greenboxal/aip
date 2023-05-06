package forddb

import "context"

type GetOptions struct {
}

type DeleteOptions struct {
}

type ListOptions struct {
	Offset int
	Limit  int
}

func NewListOptions(opts ...ListOption) ListOptions {
	var result ListOptions

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type ListOption func(opts *ListOptions)

func WithOffset(offset int) ListOption {
	return func(opts *ListOptions) {
		opts.Offset = offset
	}
}

func WithLimit(limit int) ListOption {
	return func(opts *ListOptions) {
		opts.Limit = limit
	}
}

type Storage interface {
	List(
		ctx context.Context,
		typ ResourceTypeID,
		opts ListOptions,
	) ([]RawResource, error)

	Get(
		ctx context.Context,
		typ ResourceTypeID,
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
