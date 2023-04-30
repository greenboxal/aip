package impl

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/indexing"
)

type indexContext struct {
	ctx   context.Context
	index *Index
}

func newIndexContext(ctx context.Context, index *Index) *indexContext {
	return &indexContext{
		ctx:   ctx,
		index: index,
	}
}

func (ictx *indexContext) Index() indexing.Index {
	return ictx.index
}

func (ictx *indexContext) Context() context.Context {
	return ictx.ctx
}

func (ictx *indexContext) ResolveRelative(options indexing.ResolveOptions) (collective.MemoryID, error) {
	//TODO implement me
	panic("implement me")
}

func (ictx *indexContext) ResolveMemory(id collective.MemoryID) (*collective.Memory, error) {
	return ictx.index.storage.GetMemory(ictx.ctx, id)
}

func (ictx *indexContext) ResolveMemoryRelative(options indexing.ResolveOptions) (*collective.Memory, error) {
	id, err := ictx.ResolveRelative(options)

	if err != nil {
		return nil, err
	}

	return ictx.ResolveMemory(id)
}

func (ictx *indexContext) AppendSegment(segment *collective.MemorySegment) error {
	return ictx.index.storage.AppendSegment(ictx.ctx, segment)
}
