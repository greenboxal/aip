package impl

import (
	"context"

	"github.com/greenboxal/aip/pkg/indexing"
)

type indexContext struct {
	ctx   context.Context
	index *Index
}

func (ictx *indexContext) Index() indexing.Index {
	return ictx.index
}

func (ictx *indexContext) Context() context.Context {
	return ictx.ctx
}

func (ictx *indexContext) ResolveRelative(options indexing.ResolveOptions) (indexing.MemoryID, error) {
	//TODO implement me
	panic("implement me")
}

func (ictx *indexContext) ResolveMemory(id indexing.MemoryID) (*indexing.Memory, error) {
	//TODO implement me
	panic("implement me")
}

func (ictx *indexContext) ResolveMemoryRelative(options indexing.ResolveOptions) (*indexing.Memory, error) {
	id, err := ictx.ResolveRelative(options)

	if err != nil {
		return nil, err
	}

	return ictx.ResolveMemory(id)
}

func (ictx *indexContext) AppendSegment(segment *indexing.MemorySegment) error {
	return ictx.index.storage.AppendSegment(ictx.ctx, segment)
}

func newIndexContext(ctx context.Context, index *Index) *indexContext {
	return &indexContext{
		ctx:   ctx,
		index: index,
	}
}
