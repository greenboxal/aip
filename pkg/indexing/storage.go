package indexing

import (
	"context"
)

type Storage interface {
	AppendSegment(ctx context.Context, segment *MemorySegment) error
}
