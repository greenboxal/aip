package forddb

import (
	"context"
)

type OnConflict int

const (
	OnConflictError OnConflict = iota
	OnConflictOptimistic
	OnConflictLatestWins
	OnConflictReplace
)

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

type Database interface {
	HasListeners

	List(ctx context.Context, typ ResourceTypeID, options ...ListOption) ([]BasicResource, error)
	Get(ctx context.Context, typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	Put(ctx context.Context, resource BasicResource, options ...PutOption) (BasicResource, error)
	Delete(ctx context.Context, resource BasicResource) (BasicResource, error)
}
