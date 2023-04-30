package indexing

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective"
)

type ResolveOptions struct {
	Context context.Context

	BaseMemoryID collective.MemoryID

	// Absolute anchors (optional)
	RootMemoryID   *collective.MemoryID
	BranchMemoryID *collective.MemoryID
	TargetMemoryID *collective.MemoryID

	// Relative offset from the base memory.
	RelativeHeight int64
	RelativeClock  int64

	// Optional hints for the index to use when resolving the relative memory.
	RootMemoryHint   *collective.Memory
	BranchMemoryHint *collective.Memory
	ParentMemoryHint *collective.Memory
	BaseMemoryHint   *collective.Memory
}

type IndexContext interface {
	Index() Index
	Context() context.Context

	ResolveRelative(options ResolveOptions) (collective.MemoryID, error)

	ResolveMemory(id collective.MemoryID) (*collective.Memory, error)
	ResolveMemoryRelative(options ResolveOptions) (*collective.Memory, error)

	AppendSegment(segment *collective.MemorySegment) error
}
