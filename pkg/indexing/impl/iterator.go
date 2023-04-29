package impl

import (
	"github.com/greenboxal/aip/pkg/indexing"
)

type Iterator struct {
	index   *Index
	context indexing.IndexContext

	rootMemoryID    indexing.MemoryID
	branchMemoryID  indexing.MemoryID
	parentMemoryID  indexing.MemoryID
	currentMemoryID indexing.MemoryID

	currentHeight uint64
	currentClock  uint64

	current *indexing.Memory
}

func (s *Iterator) Index() indexing.Index {
	return s.index
}

func (s *Iterator) RootMemoryID() indexing.MemoryID {
	return s.rootMemoryID
}

func (s *Iterator) BranchMemoryID() indexing.MemoryID {
	return s.branchMemoryID
}

func (s *Iterator) ParentMemoryID() indexing.MemoryID {
	return s.parentMemoryID
}

func (s *Iterator) CurrentMemoryID() indexing.MemoryID {
	return s.currentMemoryID
}

func (s *Iterator) SeekRelative(offset int) error {
	return s.setCoordinates(s.currentHeight, s.currentClock+uint64(offset))
}

func (s *Iterator) SeekTo(id indexing.MemoryID) error {
	m, err := s.context.ResolveMemory(id)

	if err != nil {
		return err
	}

	return s.setCurrent(m)
}

func (s *Iterator) GetCurrentMemory() indexing.Memory {
	if s.current == nil {
		panic("no current memory")
	}

	return *s.current
}

func (s *Iterator) GetCurrentMemoryData() indexing.MemoryData {
	if s.current == nil {
		return nil
	}

	return s.current.Data
}

func (s *Iterator) setCurrent(m *indexing.Memory) error {
	return s.set(m.RootMemoryID, m.BranchMemoryID, m.ParentMemoryID, m.ID, m.Height, m.Clock, m)
}

func (s *Iterator) set(root, branch, parent, current indexing.MemoryID, height, clock uint64, m *indexing.Memory) error {
	s.current = m
	s.currentHeight = height
	s.currentClock = clock
	s.rootMemoryID = root
	s.branchMemoryID = branch
	s.parentMemoryID = parent
	s.currentMemoryID = current

	return s.checkCurrent()
}

func (s *Iterator) setPointers(root, branch, parent, current indexing.MemoryID) error {
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
	query := indexing.ResolveOptions{
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
