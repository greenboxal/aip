package impl

import (
	"context"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing2"
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

func (ictx *indexContext) Index() indexing2.Index {
	return ictx.index
}

func (ictx *indexContext) Context() context.Context {
	return ictx.ctx
}

func (ictx *indexContext) ResolveRelative(options indexing2.ResolveOptions) (collective2.MemoryID, error) {
	//TODO implement me
	panic("implement me")
}

func (ictx *indexContext) ResolveMemory(id collective2.MemoryID) (*collective2.Memory, error) {
	return ictx.index.storage.GetMemory(ictx.ctx, id)
}

func (ictx *indexContext) ResolveMemoryRelative(options indexing2.ResolveOptions) (*collective2.Memory, error) {
	id, err := ictx.ResolveRelative(options)

	if err != nil {
		return nil, err
	}

	return ictx.ResolveMemory(id)
}

func (ictx *indexContext) AppendSegment(segment *collective2.MemorySegment) error {
	return ictx.index.storage.AppendSegment(ictx.ctx, segment)
}
