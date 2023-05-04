package impl

import (
	"context"

	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing"
)

type Index struct {
	cfg indexing2.IndexConfiguration

	storage indexing2.MemoryStorage
}

func NewIndex(storage indexing2.MemoryStorage, cfg indexing2.IndexConfiguration) *Index {
	return &Index{
		cfg: cfg,

		storage: storage,
	}
}

func (idx *Index) Configuration() indexing2.IndexConfiguration {
	return idx.cfg
}

func (idx *Index) OpenSession(ctx context.Context, options indexing2.SessionOptions) (indexing2.Session, error) {
	return NewSession(ctx, idx, options)
}

func (idx *Index) CreateContext(ctx context.Context) (indexing2.IndexContext, error) {
	return newIndexContext(ctx, idx), nil
}
