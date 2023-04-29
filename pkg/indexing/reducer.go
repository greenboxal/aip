package indexing

import (
	"context"
)

type ReducerContext struct {
	Context context.Context
	Session Session
	Segment *MemorySegment

	Hint string
}

type Reducer interface {
	ReduceSegment(ctx *ReducerContext) (*MemorySegment, error)
}
