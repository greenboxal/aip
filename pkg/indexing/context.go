package indexing

import "context"

type ResolveOptions struct {
	Context context.Context

	BaseMemoryID MemoryID

	// Absolute anchors (optional)
	RootMemoryID   *MemoryID
	BranchMemoryID *MemoryID
	TargetMemoryID *MemoryID

	// Relative offset from the base memory.
	RelativeHeight int64
	RelativeClock  int64

	// Optional hints for the index to use when resolving the relative memory.
	RootMemoryHint   *Memory
	BranchMemoryHint *Memory
	ParentMemoryHint *Memory
	BaseMemoryHint   *Memory
}

type IndexContext interface {
	Index() Index
	Context() context.Context

	ResolveRelative(options ResolveOptions) (MemoryID, error)

	ResolveMemory(id MemoryID) (*Memory, error)
	ResolveMemoryRelative(options ResolveOptions) (*Memory, error)

	AppendSegment(segment *MemorySegment) error
}
