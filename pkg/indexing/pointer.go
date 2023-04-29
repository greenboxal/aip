package indexing

import "fmt"

type HasMemoryID interface {
	GetMemoryID() MemoryID
}

type HasAnchorPointers interface {
	GetRootMemoryID() MemoryID
	GetBranchMemoryID() MemoryID
	GetParentMemoryID() MemoryID
}

type HasRelativePointers interface {
	GetClock() uint64
	GetHeight() uint64
}

type RelativePointer interface {
	HasRelativePointers

	AddClock(offset int64) RelativePointer
	AddHeight(offset int64) RelativePointer
	Add(offset MemoryRelativeOffset) RelativePointer
	Sub(other RelativePointer) MemoryRelativeOffset

	String() string
}

type MemoryRelativeOffset struct {
	Clock  int64
	Height int64
}

type MemoryAbsolutePointer interface {
	HasMemoryID
	MemoryAbsoluteAddress
}

type MemoryAbsoluteAddress interface {
	HasAnchorPointers
	HasRelativePointers
}

type memoryAbsoluteAddress struct {
	HasAnchorPointers
	HasRelativePointers
}

type memoryAnchorPointers struct {
	root   MemoryID
	branch MemoryID
	parent MemoryID
}

func (m memoryAnchorPointers) GetRootMemoryID() MemoryID {
	return m.root
}

func (m memoryAnchorPointers) GetBranchMemoryID() MemoryID {
	return m.branch
}

func (m memoryAnchorPointers) GetParentMemoryID() MemoryID {
	return m.parent
}

type memoryRelativePointer struct {
	clock  uint64
	height uint64
}

func (m memoryRelativePointer) GetClock() uint64 {
	return m.clock
}

func (m memoryRelativePointer) GetHeight() uint64 {
	return m.height
}

func (m memoryRelativePointer) AddClock(offset int64) RelativePointer {
	return memoryRelativePointer{
		clock:  m.clock + uint64(offset),
		height: m.height,
	}
}

func (m memoryRelativePointer) AddHeight(offset int64) RelativePointer {
	return memoryRelativePointer{
		clock:  m.clock,
		height: m.height + uint64(offset),
	}
}

func (m memoryRelativePointer) Add(offset MemoryRelativeOffset) RelativePointer {
	return memoryRelativePointer{
		clock:  m.clock + uint64(offset.Clock),
		height: m.height + uint64(offset.Height),
	}
}

func (m memoryRelativePointer) Sub(other RelativePointer) MemoryRelativeOffset {
	return MemoryRelativeOffset{
		Clock:  int64(m.clock - other.GetClock()),
		Height: int64(m.height - other.GetHeight()),
	}
}

func (m memoryRelativePointer) String() string {
	return fmt.Sprintf("HasRelativePointers(clock=%d, height=%d)", m.clock, m.height)
}

func Absolute(anchors HasAnchorPointers, relative HasRelativePointers) MemoryAbsoluteAddress {
	return memoryAbsoluteAddress{
		HasAnchorPointers:   anchors,
		HasRelativePointers: relative,
	}
}

func Anchors(root, branch, parent MemoryID) HasAnchorPointers {
	return memoryAnchorPointers{
		root:   root,
		branch: branch,
		parent: parent,
	}
}

func Relative(clock, height uint64) HasRelativePointers {
	return memoryRelativePointer{
		clock:  clock,
		height: height,
	}
}
