package indexing

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type ReducerContext struct {
	Context context.Context
	Session Session
	Segment *collective.MemorySegment

	Hint string
}

type Reducer interface {
	ReduceSegment(ctx *ReducerContext) (*collective.MemorySegment, error)
}
