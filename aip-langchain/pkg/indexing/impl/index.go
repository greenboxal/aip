package impl

import (
	"context"

	indexing22 "github.com/greenboxal/aip/aip-langchain/pkg/indexing"
)

type Index struct {
	cfg indexing22.IndexConfiguration

	storage indexing22.MemoryStorage
}

func NewIndex(storage indexing22.MemoryStorage, cfg indexing22.IndexConfiguration) *Index {
	return &Index{
		cfg: cfg,

		storage: storage,
	}
}

func (idx *Index) Configuration() indexing22.IndexConfiguration {
	return idx.cfg
}

func (idx *Index) OpenSession(ctx context.Context, options indexing22.SessionOptions) (indexing22.Session, error) {
	return NewSession(ctx, idx, options)
}

func (idx *Index) CreateContext(ctx context.Context) (indexing22.IndexContext, error) {
	return newIndexContext(ctx, idx), nil
}
