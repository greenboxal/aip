package impl

import (
	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	indexing22 "github.com/greenboxal/aip/aip-langchain/pkg/indexing"
)

type Iterator struct {
	index   *Index
	context indexing22.IndexContext

	rootMemoryID    collective2.MemoryID
	branchMemoryID  collective2.MemoryID
	parentMemoryID  collective2.MemoryID
	currentMemoryID collective2.MemoryID

	currentHeight uint64
	currentClock  uint64

	current *collective2.Memory
}

func NewIterator(index *Index, ctx indexing22.IndexContext) *Iterator {
	return &Iterator{
		index:   index,
		context: ctx,
	}
}

func (s *Iterator) Index() indexing22.Index {
	return s.index
}

func (s *Iterator) MemoryAddress() collective2.MemoryAbsoluteAddress {
	return collective2.Absolute(
		collective2.Anchors(s.rootMemoryID, s.branchMemoryID, s.parentMemoryID),
		collective2.Relative(s.currentClock, s.currentHeight),
	)
}

func (s *Iterator) RootMemoryID() collective2.MemoryID {
	return s.rootMemoryID
}

func (s *Iterator) BranchMemoryID() collective2.MemoryID {
	return s.branchMemoryID
}

func (s *Iterator) ParentMemoryID() collective2.MemoryID {
	return s.parentMemoryID
}

func (s *Iterator) CurrentMemoryID() collective2.MemoryID {
	return s.currentMemoryID
}

func (s *Iterator) SeekRelative(offset int) error {
	return s.setCoordinates(s.currentHeight, s.currentClock+uint64(offset))
}

func (s *Iterator) SeekTo(id collective2.MemoryID) error {
	m, err := s.context.ResolveMemory(id)

	if err != nil {
		return err
	}

	return s.setCurrent(m)
}

func (s *Iterator) GetCurrentMemory() collective2.Memory {
	if s.current == nil {
		panic("no current memory")
	}

	return *s.current
}

func (s *Iterator) GetCurrentMemoryData() collective2.MemoryData {
	if s.current == nil {
		return collective2.MemoryData{}
	}

	return s.current.Data
}

func (s *Iterator) setCurrent(m *collective2.Memory) error {
	return s.set(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID, m.ID, m.Height, m.Clock, m)
}

func (s *Iterator) set(root, branch, parent, current collective2.MemoryID, height, clock uint64, m *collective2.Memory) error {
	s.current = m
	s.currentHeight = height
	s.currentClock = clock
	s.rootMemoryID = root
	s.branchMemoryID = branch
	s.parentMemoryID = parent
	s.currentMemoryID = current

	return s.checkCurrent()
}

func (s *Iterator) setPointers(root, branch, parent, current collective2.MemoryID) error {
	s.rootMemoryID = root
	s.branchMemoryID = branch
	s.parentMemoryID = parent
	s.currentMemoryID = current

	return s.checkCurrent()
}

func (s *Iterator) setCoordinates(height, clock uint64) error {
	s.currentHeight = height
	s.currentClock = clock

	return s.checkCurrent()
}

func (s *Iterator) checkCurrent() error {
	if !s.isValid() {
		if err := s.invalidate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Iterator) invalidate() error {
	query := indexing22.ResolveOptions{
		Context: s.context.Context(),

		BaseMemoryID:   s.currentMemoryID,
		TargetMemoryID: nil,

		BaseMemoryHint: s.current,
	}

	memory, err := s.context.ResolveMemoryRelative(query)

	if err != nil {
		return err
	}

	s.current = memory

	return nil
}

func (s *Iterator) isValid() bool {
	if s.current == nil {
		return false
	}

	if s.currentMemoryID != s.current.ID {
		return false
	}

	if s.parentMemoryID != s.current.ParentMemoryID {
		return false
	}

	if s.branchMemoryID != s.current.BranchMemoryID {
		return false
	}

	if s.rootMemoryID != s.current.RootMemoryID {
		return false
	}

	if s.currentHeight != s.current.Height {
		return false
	}

	if s.currentClock != s.current.Clock {
		return false
	}

	return true
}
