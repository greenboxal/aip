package indexing

import (
	"context"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type MemoryStorage interface {
	AppendSegment(ctx context.Context, segment *collective2.MemorySegment) error

	GetSegment(ctx context.Context, id collective2.MemorySegmentID) (*collective2.MemorySegment, error)
	GetMemory(ctx context.Context, id collective2.MemoryID) (*collective2.Memory, error)
}
