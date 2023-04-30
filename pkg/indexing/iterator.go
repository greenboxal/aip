package indexing

import "github.com/greenboxal/aip/pkg/collective"

type Iterator interface {
	RootMemoryID() collective.MemoryID
	BranchMemoryID() collective.MemoryID
	ParentMemoryID() collective.MemoryID
	CurrentMemoryID() collective.MemoryID
	MemoryAddress() collective.MemoryAbsoluteAddress

	SeekRelative(offset int) error
	SeekTo(id collective.MemoryID) error

	GetCurrentMemory() collective.Memory
	GetCurrentMemoryData() collective.MemoryData
}
