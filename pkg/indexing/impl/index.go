package impl

import (
	"context"

	"github.com/greenboxal/aip/pkg/indexing"
)

type Index struct {
	cfg indexing.IndexConfiguration
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
