package forddb

import (
	"context"
)

type Database interface {
	HasListeners

	LogStore() LogStore

	List(ctx context.Context, typ TypeID, options ...ListOption) ([]BasicResource, error)
	Get(ctx context.Context, typ TypeID, id BasicResourceID) (BasicResource, error)
	Put(ctx context.Context, resource BasicResource, options ...PutOption) (BasicResource, error)
	Delete(ctx context.Context, resource BasicResource) (BasicResource, error)
}

type QueryOptions struct {
	ReadConsistency ReadConsistencyLevel
}

type GetOptions struct {
	QueryOptions
}

type QueryOption func(opts *QueryOptions)

func WithReadConsistency(level ReadConsistencyLevel) QueryOption {
	return func(opts *QueryOptions) {
		opts.ReadConsistency = level
	}
}

type GetOption func(opts *GetOptions)

func WithGetQueryOptions(options ...QueryOption) GetOption {
	return func(opts *GetOptions) {
		for _, opt := range options {
			opt(&opts.QueryOptions)
		}
	}
}

type DeleteOptions struct {
}

type SortField struct {
	Path  string    `json:"field"`
	Order SortOrder `json:"order"`
}

type ListOptions struct {
	QueryOptions

	Offset     int
	Limit      int
	SortFields []SortField

	ResourceIDs []BasicResourceID
}

func NewListOptions(opts ...ListOption) ListOptions {
	var result ListOptions

	for _, opt := range opts {
		opt(&result)
	}

	return result
}

type ListOption func(opts *ListOptions)

func WithListQueryOptions(options ...QueryOption) ListOption {
	return func(opts *ListOptions) {
		for _, opt := range options {
			opt(&opts.QueryOptions)
		}
	}
}

func WithSortField(field string, order SortOrder) ListOption {
	return func(opts *ListOptions) {
		opts.SortFields = append(opts.SortFields, SortField{
			Path:  field,
			Order: order,
		})
	}
}

func WithSortFields(fields ...SortField) ListOption {
	return func(opts *ListOptions) {
		opts.SortFields = append(opts.SortFields, fields...)
	}
}

func WithResourceIDs(ids ...BasicResourceID) ListOption {
	return func(opts *ListOptions) {
		opts.ResourceIDs = append(opts.ResourceIDs, ids...)
	}
}

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

type PutOptions struct {
	OnConflict OnConflict
}

type PutOption func(opts *PutOptions)

func WithOnConflict(onConflict OnConflict) PutOption {
	return func(opts *PutOptions) {
		opts.OnConflict = onConflict
	}
}

func NewPutOptions(options ...PutOption) PutOptions {
	opts := PutOptions{
		OnConflict: OnConflictOptimistic,
	}

	for _, option := range options {
		option(&opts)
	}

	return opts
}
