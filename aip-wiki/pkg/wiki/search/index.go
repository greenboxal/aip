package search

import "context"

type Options struct {
}

type Option func(opts *Options)

type Hit struct {
	Score float32
}

type Result struct {
	Hits []Hit
}

type Index interface {
	Search(ctx context.Context, query string, options ...Option) (*Result, error)
}
