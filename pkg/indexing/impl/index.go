package impl

import (
	"context"

	"github.com/greenboxal/aip/pkg/indexing"
)

type Index struct {
	cfg indexing.IndexConfiguration

	storage indexing.MemoryStorage
}

func NewIndex(storage indexing.MemoryStorage, cfg indexing.IndexConfiguration) *Index {
	return &Index{
		cfg: cfg,

		storage: storage,
	}
}

func (idx *Index) Configuration() indexing.IndexConfiguration {
	return idx.cfg
}

func (idx *Index) OpenSession(ctx context.Context, options indexing.SessionOptions) (indexing.Session, error) {
	return NewSession(ctx, idx, options)
}

func (idx *Index) CreateContext(ctx context.Context) (indexing.IndexContext, error) {
	return newIndexContext(ctx, idx), nil
}
