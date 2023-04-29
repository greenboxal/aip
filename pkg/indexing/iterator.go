package indexing

type Iterator interface {
	RootMemoryID() MemoryID
	BranchMemoryID() MemoryID
	ParentMemoryID() MemoryID
	CurrentMemoryID() MemoryID

	SeekRelative(offset int) error
	SeekTo(id MemoryID) error

	GetCurrentMemory() Memory
	GetCurrentMemoryData() MemoryData
}
