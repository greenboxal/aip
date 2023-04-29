package indexing

type Iterator interface {
	RootMemoryID() MemoryID
	BranchMemoryID() MemoryID
	ParentMemoryID() MemoryID
	CurrentMemoryID() MemoryID
	MemoryAddress() MemoryAbsoluteAddress

	SeekRelative(offset int) error
	SeekTo(id MemoryID) error

	GetCurrentMemory() Memory
	GetCurrentMemoryData() MemoryData
}
