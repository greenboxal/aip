package indexing2

import (
	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type Iterator interface {
	RootMemoryID() collective2.MemoryID
	BranchMemoryID() collective2.MemoryID
	ParentMemoryID() collective2.MemoryID
	CurrentMemoryID() collective2.MemoryID
	MemoryAddress() collective2.MemoryAbsoluteAddress

	SeekRelative(offset int) error
	SeekTo(id collective2.MemoryID) error

	GetCurrentMemory() collective2.Memory
	GetCurrentMemoryData() collective2.MemoryData
}
