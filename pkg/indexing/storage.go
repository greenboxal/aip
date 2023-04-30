package indexing

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
)

type MemoryStorage interface {
	AppendSegment(ctx context.Context, segment *collective.MemorySegment) error

	GetSegment(ctx context.Context, id collective.MemorySegmentID) (*collective.MemorySegment, error)
	GetMemory(ctx context.Context, id collective.MemoryID) (*collective.Memory, error)
}
