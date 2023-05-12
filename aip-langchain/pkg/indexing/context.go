package indexing

import (
	"context"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type ResolveOptions struct {
	Context context.Context

	BaseMemoryID collective2.MemoryID

	// Absolute anchors (optional)
	RootMemoryID   *collective2.MemoryID
	BranchMemoryID *collective2.MemoryID
	TargetMemoryID *collective2.MemoryID

	// Relative offset from the base memory.
	RelativeHeight int64
	RelativeClock  int64

	// Optional hints for the index to use when resolving the relative memory.
	RootMemoryHint   *collective2.Memory
	BranchMemoryHint *collective2.Memory
	ParentMemoryHint *collective2.Memory
	BaseMemoryHint   *collective2.Memory
}

type IndexContext interface {
	Index() Index
	Context() context.Context

	ResolveRelative(options ResolveOptions) (collective2.MemoryID, error)

	ResolveMemory(id collective2.MemoryID) (*collective2.Memory, error)
	ResolveMemoryRelative(options ResolveOptions) (*collective2.Memory, error)

	AppendSegment(segment *collective2.MemorySegment) error
}
