package storage

import (
	"context"

	"github.com/greenboxal/aip/pkg/indexing"
)

type Storage interface {
	AppendSegment(ctx context.Context, segment *indexing.MemorySegment) error
}
